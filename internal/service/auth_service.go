package service

import (
	"challenge2/internal/models"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidUser  = errors.New("username must be at least 4 characters")
	ErrInvalidToken = errors.New("invalid session token")
)

// AuthService is the in-memory source of truth for users and sessions.
type AuthService struct {
	users    map[string]models.User // username -> user
	sessions map[string]string      // token -> username
	mu       sync.RWMutex
}

// NewAuthService constructs a new AuthService and seeds the default admin user.
func NewAuthService() *AuthService {
	s := &AuthService{
		users:    make(map[string]models.User),
		sessions: make(map[string]string),
	}

	// seed admin user: admin / password123
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	admin := models.User{
		ID:        generateID(),
		Username:  "admin",
		Password:  string(hashed),
		CreatedAt: time.Now().UTC(),
	}

	s.users["admin"] = admin
	return s
}

// generateID creates a short random identifier for users.
func generateID() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// generateToken creates a random session token.
func generateToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// Register creates a new user if it passes validation rules.
func (s *AuthService) Register(username, password string) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(username) < 4 {
		return models.User{}, ErrInvalidUser
	}

	if _, exists := s.users[username]; exists {
		return models.User{}, ErrUserExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		ID:        generateID(),
		Username:  username,
		Password:  string(hashed),
		CreatedAt: time.Now().UTC(),
	}

	s.users[username] = user
	return user, nil
}

// Login validates credentials and returns a new session token.
func (s *AuthService) Login(username, password string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[username]
	if !exists {
		return "", ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrUnauthorized
	}

	token := generateToken()
	s.sessions[token] = username

	return token, nil
}

// GetProfileByToken returns the user associated with a given session token.
func (s *AuthService) GetProfileByToken(token string) (models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	username, exists := s.sessions[token]
	if !exists {
		return models.User{}, ErrInvalidToken
	}

	user, exists := s.users[username]
	if !exists {
		return models.User{}, ErrInvalidToken
	}

	return user, nil
}
