package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"shortener/models"
	"shortener/models/options"
	"shortener/repository"
	"shortener/utils"

	"github.com/davecgh/go-spew/spew"

	"github.com/gorilla/mux"
)

// LinkController represent link repository
type LinkController struct {
	linkRepository  repository.LinkRepository
	usageRepository repository.UsageRepository
}

// NewUserController func returns UserController object
func NewLinkController(db *sql.DB) LinkController {
	return LinkController{
		linkRepository:  repository.NewLinkRepository(db),
		usageRepository: repository.NewUsageRepository(db),
	}
}

// List returns user links using offset and limit
func (controller *LinkController) List(w http.ResponseWriter, r *http.Request) {
	opts := options.NewOptionsFromContext(r.Context())
	user, _ := models.NewUserFromContext(r.Context())

	spew.Dump("r.Context()", r.Context())

	links, err := controller.linkRepository.GetUserLinks(r.Context(), *user, *opts)

	if err != nil {
		utils.RespondWithError(&w, http.StatusBadRequest, models.NewError("Id is not provided"))
		return
	}

	utils.RespondWithJSON(&w, http.StatusOK, links)
}

// FetchByID redirects user to the link
func (controller *LinkController) FetchByID(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]

	if !ok {
		utils.RespondWithError(&w, http.StatusBadRequest, models.NewError("Id is not provided"))
		return
	}

	link, err := controller.linkRepository.FindByID(r.Context(), models.Link{
		ID: id,
	})

	if err != nil {
		utils.RespondWithError(&w, http.StatusNotFound, models.NewError(err.Error()))
		return
	}

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		usage, err := controller.usageRepository.Create(ctx, link.ID)
		fmt.Println("usage", usage)
		if err != nil {
			log.Println("--- error ---", err)
		}
	}()

	utils.RedirectToAnotherResource(&w, link.URL)
}

// Create saves link to the database
func (controller *LinkController) Create(w http.ResponseWriter, r *http.Request) {
	var link models.Link
	user, err := models.NewUserFromContext(r.Context())

	if err != nil {
		utils.RespondWithError(&w, http.StatusBadRequest, models.NewError(err.Error()))
		return
	}

	err = link.Populate(r)
	link.UserID = user.ID

	if err != nil {
		utils.RespondWithError(&w, http.StatusBadRequest, models.NewError(err.Error()))
		return
	}

	link, err = controller.linkRepository.Create(r.Context(), link)

	if err != nil {
		utils.RespondWithError(&w, http.StatusBadRequest, models.NewError(err.Error()))
		return
	}

	utils.RespondWithJSON(&w, http.StatusCreated, link)
}
