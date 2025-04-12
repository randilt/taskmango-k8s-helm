package controllers

import (
	"net/http"
	"time"

	"taskmango/authsvc/internal/config"
	"taskmango/authsvc/internal/models"
	"taskmango/authsvc/internal/repositories"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthController struct {
	userRepo    *repositories.UserRepository
	cfg         *config.Config
	argonParams *argon2id.Params
}

func NewAuthController(userRepo *repositories.UserRepository, cfg *config.Config) *AuthController {
	return &AuthController{
		userRepo: userRepo,
		cfg:      cfg,
		argonParams: &argon2id.Params{
			Memory:      64 * 1024, // 64 MB
			Iterations:  3,
			Parallelism: 4,
			SaltLength:  16,
			KeyLength:   32,
		},
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Message string `json:"message,omitempty"`
	Token   string `json:"token,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, AuthResponse{Error: "Username and password are required"})
		return
	}

	// Check if username already exists
	exists, err := c.userRepo.UsernameExists(req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, AuthResponse{Error: "Database error"})
		return
	}
	if exists {
		ctx.JSON(http.StatusConflict, AuthResponse{Error: "Username already exists"})
		return
	}

	// Generate PHC-formatted hash
	hash, err := argon2id.CreateHash(req.Password, c.argonParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, AuthResponse{Error: "Failed to hash password"})
		return
	}

	user := &models.User{
		Username:     req.Username,
		PasswordHash: hash,
	}

	if err := c.userRepo.Create(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, AuthResponse{Error: "Failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, AuthResponse{Message: "User registered successfully"})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, AuthResponse{Error: "Username and password are required"})
		return
	}

	user, err := c.userRepo.FindByUsername(req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusUnauthorized, AuthResponse{Error: "Invalid credentials"})
		} else {
			ctx.JSON(http.StatusInternalServerError, AuthResponse{Error: "Database error"})
		}
		return
	}

	// Verify password against PHC-formatted hash
	match, err := argon2id.ComparePasswordAndHash(req.Password, user.PasswordHash)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, AuthResponse{Error: "Failed to verify password"})
		return
	}

	if !match {
		ctx.JSON(http.StatusUnauthorized, AuthResponse{Error: "Invalid credentials"})
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Username,
		"user_id": user.ID,
		"iat": time.Now().Unix(),
		"iss": c.cfg.JWTIssuer,
		"aud": c.cfg.JWTAudience,
		"exp": time.Now().Add(time.Duration(c.cfg.JWTValidity) * time.Second).Unix(),
	})

	tokenString, err := token.SignedString([]byte(c.cfg.JWTSigningKey))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, AuthResponse{Error: "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, AuthResponse{Token: tokenString})
}