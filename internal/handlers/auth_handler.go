package handlers

import (
	"challenge2/internal/middleware"
	"challenge2/internal/service"
	"encoding/json"
	"net/http"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: s}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// RegisterHandler handles POST /api/user/register.
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid or missing JSON body", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(req.Username, req.Password)
	if err != nil {
		switch err {
		case service.ErrInvalidUser:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case service.ErrUserExists:
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

// LoginHandler handles POST /api/user/login.
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid or missing JSON body", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		switch err {
		case service.ErrUnauthorized:
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, http.StatusOK, LoginResponse{Token: token})
}

// ProfileHandler handles GET /api/user/profile.
func (h *AuthHandler) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	writeJSON(w, http.StatusOK, user)
}
