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
	userRepository repository.UserRepositoryInterface
}

// UserRequest represents object of user request
type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
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

	createdUser, err := controller.userRepository.CreateWithContext(ctx, user)

	if err != nil {
		utils.RespondWithError(&w, http.StatusBadRequest, models.NewError(err.Error()))
		return
	}

	utils.RespondWithJSON(&w, http.StatusCreated, *createdUser)
}

// Authorize user to get access
func (controller *UserController) Authorize(w http.ResponseWriter, r *http.Request) {
	user, err := models.NewUserFromRequest(r)
	plainTextPassword := user.Password

	if err != nil {
		utils.RespondWithError(&w, http.StatusBadRequest, models.NewError(err.Error()))
		return
	}

	foundUser, err := controller.userRepository.FindByLoginWithContext(r.Context(), user.Login)

	if err != nil {
		utils.RespondWithError(&w, http.StatusBadRequest, models.NewError(err.Error()))
		return
	}

	ok := crypto.ValidatePassword(plainTextPassword, foundUser.Password)

	if !ok {
		utils.RespondWithError(&w, http.StatusUnauthorized, models.Error{"User is not authorized"})
		return
	}

	foundUser.CleanPrivateFields()

	token, err := models.GenerateAuthToken(*foundUser)

	utils.RespondWithJSON(&w, http.StatusOK, token)
}
