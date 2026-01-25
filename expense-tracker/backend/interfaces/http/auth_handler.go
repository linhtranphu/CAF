package http

import (
	"net/http"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserRepository interface {
	CreateUser(username, password string) error
	GetUser(username string) (string, error)
}

type AuthHandler struct {
	userRepo UserRepository
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewAuthHandler(userRepo UserRepository) *AuthHandler {
	return &AuthHandler{userRepo: userRepo}
}

func (h *AuthHandler) LoginPage(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/admin")
		return
	}
	c.HTML(http.StatusOK, "login.html", nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
	log.Printf("[AUTH] Login request received")
	
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[AUTH] Invalid JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	log.Printf("[AUTH] Login attempt: %s", req.Username)
	
	// Check credentials from database
	password, err := h.userRepo.GetUser(req.Username)
	if err != nil || password != req.Password {
		log.Printf("[AUTH] Failed login attempt: %s", req.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Save to session
	session := sessions.Default(c)
	session.Set("user_id", req.Username)
	session.Set("username", req.Username)
	if err := session.Save(); err != nil {
		log.Printf("[AUTH] Session save error: %v", err)
	}

	log.Printf("[AUTH] User logged in successfully: %s", req.Username)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func (h *AuthHandler) Register(c *gin.Context) {
	log.Printf("[AUTH] Register request received")
	
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[AUTH] Invalid JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	log.Printf("[AUTH] Register attempt: %s", req.Username)
	
	// Create user in database
	err := h.userRepo.CreateUser(req.Username, req.Password)
	if err != nil {
		log.Printf("[AUTH] User already exists: %s", req.Username)
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	log.Printf("[AUTH] User registered successfully: %s", req.Username)
	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
}

func (h *AuthHandler) Logout(c *gin.Context) {session := sessions.Default(c)
	username := session.Get("username")
	session.Clear()
	session.Save()
	log.Printf("[AUTH] User logged out: %s", username)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip auth for OPTIONS requests (CORS preflight)
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}
		
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			log.Printf("[AUTH] Unauthorized access attempt to %s", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}
		c.Next()
	}
}