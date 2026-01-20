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
	// Load environment variables
	if err := loadEnv(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Infrastructure
	mongoRepo, err := mongodb.NewRepository()
	if err != nil {
		log.Fatal("Failed to create mongodb repository:", err)
	}
	defer mongoRepo.Close()

	parser := ai.NewMessageParser()

	// Application
	expenseService := services.NewExpenseService(mongoRepo, parser)

	// Interface
	expenseHandler := http.NewExpenseHandler(expenseService)
	adminHandler := http.NewAdminHandler(expenseService)
	router := http.NewRouter(expenseHandler, adminHandler)

	log.Println("Server starting on :8081")
	log.Println("Database: MongoDB (localhost:27017)")
	log.Println("Admin panel: http://localhost:8081/admin")
	if os.Getenv("OPENAI_API_KEY") != "" {
		log.Println("AI Parser: OpenAI GPT-3.5 ✅")
	} else {
		log.Println("AI Parser: Fallback (simple parsing)")
	}
	if err := router.Run(":8081"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}