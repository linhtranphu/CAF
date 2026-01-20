package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"expense-tracker/application/services"
	"expense-tracker/infrastructure/ai"
	"expense-tracker/infrastructure/sheets"
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
	sheetsRepo, err := sheets.NewRepository()
	if err != nil {
		log.Fatal("Failed to create sheets repository:", err)
	}

	parser := ai.NewMessageParser()

	// Application
	expenseService := services.NewExpenseService(sheetsRepo, parser)

	// Interface
	expenseHandler := http.NewExpenseHandler(expenseService)
	router := http.NewRouter(expenseHandler)

	log.Println("Server starting on :8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}