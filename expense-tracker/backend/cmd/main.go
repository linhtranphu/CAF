package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"expense-tracker/application/services"
	"expense-tracker/infrastructure/ai"
	"expense-tracker/infrastructure/mongodb"
	"expense-tracker/interfaces/http"
)

func loadEnv() error {
	file, err := os.Open(".env")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.Trim(strings.TrimSpace(parts[1]), "\"")
			os.Setenv(key, value)
			log.Printf("[ENV] Loaded: %s", key)
		}
	}

	return scanner.Err()
}

func main() {
	// Load environment variables from file if exists
	if err := loadEnv(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
		log.Println("Using environment variables from system")
	}

	// Set default values if not provided
	if os.Getenv("MONGODB_URI") == "" {
		os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	}
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8081")
	}

	log.Printf("MongoDB URI: %s", os.Getenv("MONGODB_URI"))
	log.Printf("Port: %s", os.Getenv("PORT"))

	// Infrastructure
	mongoRepo, err := mongodb.NewRepository()
	if err != nil {
		log.Fatal("Failed to create mongodb repository:", err)
	}
	defer mongoRepo.Close()

	parser := ai.NewMessageParser(mongoRepo)

	// Initialize default users
	if err := mongoRepo.InitDefaultUsers(); err != nil {
		log.Printf("Warning: Failed to initialize default users: %v", err)
	}

	// Application
	expenseService := services.NewExpenseService(mongoRepo, parser)

	// Interface
	expenseHandler := http.NewExpenseHandler(expenseService)
	adminHandler := http.NewAdminHandler(expenseService)
	authHandler := http.NewAuthHandler(mongoRepo)
	settingsHandler := http.NewSettingsHandler(mongoRepo)
	router := http.NewRouter(expenseHandler, adminHandler, authHandler, settingsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Server starting on :%s", port)
	log.Printf("Database: %s", os.Getenv("MONGODB_URI"))
	log.Printf("Admin panel: http://localhost:%s/admin", port)
	if os.Getenv("GEMINI_API_KEY") != "" {
		log.Println("AI Parser: Gemini ✅")
	} else if os.Getenv("OPENAI_API_KEY") != "" {
		log.Println("AI Parser: OpenAI GPT-3.5 ✅")
	} else {
		log.Println("AI Parser: Fallback (simple parsing)")
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}