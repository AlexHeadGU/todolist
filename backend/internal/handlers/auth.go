package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AlexHeadGU/todolist/internal/logger"
	"github.com/AlexHeadGU/todolist/internal/models"
	"github.com/AlexHeadGU/todolist/internal/service"
	"github.com/AlexHeadGU/todolist/internal/utils"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	logger.Info("Registration attempt", "ip", r.RemoteAddr)

	var req models.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Invalid registration request", "error", err)
		utils.SendValidationError(w, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		logger.Warn("Missing registration fields", "email", req.Email != "")
		utils.SendValidationError(w, "Email and password are required")
		return
	}

	user, err := h.authService.Register(req.Email, req.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			logger.Warn("Registration failed - user exists", "email", req.Email)
			utils.SendError(w, http.StatusConflict, "User already exists")
		} else {
			logger.Error("Registration failed", "error", err, "email", req.Email)
			utils.SendInternalError(w, err, "Failed to create user")
		}
		return
	}

	logger.Info("User registered successfully", "user_id", user.ID, "email", user.Email)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	logger.Info("Login attempt", "ip", r.RemoteAddr)

	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Invalid login request", "error", err)
		utils.SendValidationError(w, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		logger.Warn("Missing login fields", "email", req.Email != "")
		utils.SendValidationError(w, "Email and password are required")
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		logger.Warn("Login failed - invalid credentials", "email", req.Email)
		utils.SendUnauthorizedError(w, "Invalid credentials")
		return
	}

	logger.Info("User logged in successfully", "user_id", user.ID, "email", user.Email)

	response := models.LoginResponse{
		Token: token,
		User:  *user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
