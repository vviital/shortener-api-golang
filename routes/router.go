package routes

import (
	"database/sql"
	"errors"
	"log"
	"shortener/controllers"

	"github.com/gorilla/mux"
)

// AddOpenRoutes adds open routes to the router (gorilla mux)
// Basically we need it only for the sign in and sign up pages
func AddOpenRoutes(router *mux.Router, args ...interface{}) error {
	if len(args) == 0 {
		return errors.New("Database connection is missing")
	}

	db, ok := args[0].(*sql.DB)

	if !ok {
		log.Fatalln("Wrong parameters in the AddOpenRoutes function")
	}

	userController := controllers.NewUserController(db)

	router.HandleFunc("/users", userController.Create).Methods("POST")
	router.HandleFunc("/users/token", userController.Authorize).Methods("POST")

	return nil
}

// AddProtectedRoutes adds protected routes to the router (gorilla mux)
func AddProtectedRoutes(router *mux.Router, args ...interface{}) error {
	if len(args) == 0 {
		return errors.New("Database connection is missing")
	}

	db, ok := args[0].(*sql.DB)

	if !ok {
		log.Fatalln("Wrong parameters in the AddOpenRoutes function")
	}

	linkController := controllers.NewLinkController(db)
	router.HandleFunc("/l", linkController.Create).Methods("POST")
	router.HandleFunc("/l/{id}", linkController.FetchByID).Methods("GET")
	router.HandleFunc("/l", linkController.List).Methods("GET")
	return nil
}
