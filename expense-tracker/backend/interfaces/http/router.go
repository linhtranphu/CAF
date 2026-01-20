package http

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		
		log.Printf("[%s] %s %s - %d (%v) from %s", 
			method, path, 
			getStatusText(statusCode), statusCode, 
			latency, clientIP)
	}
}

func getStatusText(code int) string {
	if code >= 200 && code < 300 {
		return "SUCCESS"
	} else if code >= 400 && code < 500 {
		return "CLIENT_ERROR"
	} else if code >= 500 {
		return "SERVER_ERROR"
	}
	return "INFO"
}

func NewRouter(expenseHandler *ExpenseHandler, adminHandler *AdminHandler) *gin.Engine {
	r := gin.Default()
	
	// Load HTML templates
	r.LoadHTMLGlob("templates/*")
	
	// Sessions
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "default-secret-key-change-in-production"
	}
	store := cookie.NewStore([]byte(sessionSecret))
	r.Use(sessions.Sessions("expense-session", store))
	
	r.Use(LoggerMiddleware())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081", "http://localhost:3000", "https://your-domain.com"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Auth handler
	authHandler := NewAuthHandler()

	// Public routes
	r.GET("/", func(c *gin.Context) {
		log.Printf("[AUTH] Root route accessed")
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID != nil {
			log.Printf("[AUTH] User already logged in, redirecting to admin")
			c.Redirect(http.StatusTemporaryRedirect, "/admin")
			return
		}
		log.Printf("[AUTH] Showing login page")
		c.HTML(http.StatusOK, "login.html", nil)
	})
	r.POST("/auth/login", authHandler.Login)
	r.GET("/auth/logout", authHandler.Logout)

	// Protected routes
	protected := r.Group("/")
	protected.Use(AuthRequired())
	{
		protected.GET("/admin", adminHandler.AdminPage)
		protected.DELETE("/admin/expense/:id", adminHandler.DeleteExpense)
	}

	api := r.Group("/api")
	api.Use(AuthRequired())
	{
		api.POST("/expense", expenseHandler.CreateExpense)
		api.GET("/expenses", expenseHandler.GetExpenses)
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "timestamp": time.Now()})
		})
	}

	return r
}