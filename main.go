package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"regexp"
	"shortener/configuration"
	"shortener/driver"
	"shortener/migrator"
	"shortener/models"
	"shortener/models/options"
	"shortener/repository"
	"shortener/routes"
	"shortener/utils"
	"time"

	"github.com/gorilla/mux"
)

const Authorization = "Authorization"
const address = "127.0.0.1:8000"

func isAuthorizedRoute(r *http.Request, rm *mux.RouteMatch) bool {
	value := r.Header.Get(Authorization)

	return value != ""
}

func isAnonRoute(r *http.Request, rm *mux.RouteMatch) bool {
	return !isAuthorizedRoute(r, rm)
}

func anonUserMiddlewareGenerator(db *sql.DB) func(http.Handler) http.Handler {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	repository := repository.NewUserRepository(db)
	anon, err := repository.FindByLoginWithContext(ctx, configuration.GetConfiguration().AnonUserLogin)

	if err != nil {
		log.Fatalln(err)
	}

	anonUserMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := models.GenerateAuthToken(*anon)

			if err != nil {
				utils.RespondWithError(&w, http.StatusInternalServerError, models.Error{err.Error()})
				return
			}

			r.Header.Add("Authorization", "Bearer "+token.Value)

			next.ServeHTTP(w, r)
		})
	}

	return anonUserMiddleware
}

func withAuth(next http.Handler) http.Handler {
	var checkAuthHeader = func(header string) (*models.User, error) {
		r, _ := regexp.Compile("^Bearer (?P<token>.*)")
		matches := r.FindStringSubmatch(header)

		// whole string + token itself
		if len(matches) != 2 {
			return nil, errors.New("User is not authorized")
		}

		token := models.Token{Value: matches[1]}

		claims, err := token.GetClaims()

		if err != nil {
			return nil, err
		}

		return &claims.User, nil
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := checkAuthHeader(r.Header.Get(Authorization))

		if err != nil {
			utils.RespondWithError(&w, http.StatusUnauthorized, models.Error{err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func withOptions(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		options := options.NewOptionsFromRequest(r)

		ctx := context.WithValue(r.Context(), "options", options)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func stop(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func startServer(r *mux.Router) chan error {
	var errs chan error

	go func() {
		srv := &http.Server{
			Handler:      r,
			Addr:         address,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		err := srv.ListenAndServe()
		errs <- err
	}()

	return errs
}

func main() {
	db := driver.ConnectPostgreSQL()

	migrator.MigrateDatabase(db)

	r := mux.NewRouter()

	r.Use(withOptions)

	authorizedRouter := r.MatcherFunc(isAuthorizedRoute).Subrouter()
	anonRouter := r.MatcherFunc(isAnonRoute).Subrouter()
	anonRouter.Use(anonUserMiddlewareGenerator(db))

	for _, router := range []*mux.Router{authorizedRouter, anonRouter} {
		router.Use(withAuth)
		err := routes.AddProtectedRoutes(router, db)
		stop(err)
	}

	err := routes.AddOpenRoutes(anonRouter, db)

	stop(err)

	errs := startServer(r)
	log.Println("server started on", address)
	log.Fatalln(<-errs)
}
