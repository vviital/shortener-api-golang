package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"shortener/models"
	"shortener/models/crypto"
	"shortener/repository"
	"shortener/utils"
)

// UserController struct represents user controller
type UserController struct {
	userRepository repository.UserRepository
}

// NewUserController func returns UserController object
func NewUserController(db *sql.DB) UserController {
	controller := UserController{
		userRepository: repository.NewUserRepository(db),
	}

	return controller
}

// Create is a handler for create user action
func (controller *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var ctx context.Context
	user, err := models.NewUserFromRequest(r)

	if err != nil {
		utils.RespondWithError(&w, http.StatusBadRequest, models.NewError(err.Error()))
		return
	}

	ctx, cancel := context.WithCancel(r.Context())

	defer cancel()

	user.Password, err = crypto.CreatePassword(user.Password)

	user, err = controller.userRepository.Create(ctx, user)

	if err != nil {
		utils.RespondWithError(&w, http.StatusBadRequest, models.NewError(err.Error()))
		return
	}

	utils.RespondWithJSON(&w, http.StatusCreated, user)
}

func (controller *UserController) Authorize(w http.ResponseWriter, r *http.Request) {
	var ctx context.Context

	user, err := models.NewUserFromRequest(r)
	plainTextPassword := user.Password

	if err != nil {
		utils.RespondWithError(&w, http.StatusBadRequest, models.NewError(err.Error()))
		return
	}

	ctx, cancel := context.WithCancel(r.Context())

	defer cancel()

	user, err = controller.userRepository.FindByLogin(ctx, user.Login)

	if err != nil {
		utils.RespondWithError(&w, http.StatusBadRequest, models.NewError(err.Error()))
		return
	}

	ok := crypto.ValidatePassword(plainTextPassword, user.Password)

	if !ok {
		utils.RespondWithError(&w, http.StatusUnauthorized, models.Error{"User is not authorized"})
		return
	}

	user.CleanPrivateFields()

	token, err := models.GenerateAuthToken(user)

	utils.RespondWithJSON(&w, http.StatusOK, token)
}
