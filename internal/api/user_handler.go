package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/hereswilson/femProject/internal/store"
	"github.com/hereswilson/femProject/internal/utils"
)

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (h *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if len(req.Username) > 50 {
		return errors.New("username must be at most 50 characters long")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}

	emailregex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailregex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}
	if req.Bio == "" {
		return errors.New("bio is required")
	}
	return nil
}

func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Println("Error decoding register request:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	err = h.validateRegisterRequest(&req)
	if err != nil {
		h.logger.Println("Validation error:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		h.logger.Println("Error setting password hash:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	err = h.userStore.CreateUser(user)
	if err != nil {
		h.logger.Println("Error creating user:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": user})
}
