package sheets

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"expense-tracker/domain/expense"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Repository struct {
	service       *sheets.Service
	spreadsheetID string
}

func NewRepository() (*Repository, error) {
	ctx := context.Background()
	
	// Get credentials from environment
	email := os.Getenv("GOOGLE_SERVICE_ACCOUNT_EMAIL")
	privateKey := os.Getenv("GOOGLE_PRIVATE_KEY")
	spreadsheetID := os.Getenv("GOOGLE_SHEETS_ID")
	
	if email == "" || privateKey == "" || spreadsheetID == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	// Create service account credentials JSON
	creds := map[string]interface{}{
		"type":         "service_account",
		"project_id":   "cafeshop-484807",
		"private_key_id": "dummy",
		"private_key":  privateKey,
		"client_email": email,
		"client_id":    "dummy",
		"auth_uri":     "https://accounts.google.com/o/oauth2/auth",
		"token_uri":    "https://oauth2.googleapis.com/token",
	}

	credsJSON, err := json.Marshal(creds)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal credentials: %v", err)
	}

	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON(credsJSON), option.WithScopes(sheets.SpreadsheetsScope))
	if err != nil {
		return nil, err
	}

	return &Repository{
		service:       srv,
		spreadsheetID: spreadsheetID,
	}, nil
}

func (r *Repository) Save(exp *expense.Expense) error {
	start := time.Now()
	log.Printf("[SHEETS] Starting save operation for expense: %s, amount: %d", exp.Items(), exp.Amount())

	resp, err := r.service.Spreadsheets.Values.Get(r.spreadsheetID, "cost!A:E").Do()
	if err != nil {
		log.Printf("[SHEETS] ERROR: Failed to get sheet data: %v", err)
		return err
	}
	log.Printf("[SHEETS] Retrieved %d existing rows from sheet", len(resp.Values))

	rowNumber := len(resp.Values) + 1
	if rowNumber == 1 {
		log.Printf("[SHEETS] Creating header row")
		headerValues := [][]interface{}{
			{"No.", "Items", "Amount", "Paid Date", "Paid By"},
		}
		_, err = r.service.Spreadsheets.Values.Update(r.spreadsheetID, "cost!A1:E1", &sheets.ValueRange{
			Values: headerValues,
		}).ValueInputOption("RAW").Do()
		if err != nil {
			log.Printf("[SHEETS] ERROR: Failed to create header: %v", err)
			return err
		}
		log.Printf("[SHEETS] Header created successfully")
		rowNumber = 2
	}

	values := [][]interface{}{
		{rowNumber - 1, exp.Items(), exp.Amount(), exp.PaidDate().Format("02-Jan-2006"), exp.PaidBy()},
	}

	rangeToUpdate := fmt.Sprintf("cost!A%d:E%d", rowNumber, rowNumber)
	log.Printf("[SHEETS] Updating range: %s with data: %v", rangeToUpdate, values[0])
	
	_, err = r.service.Spreadsheets.Values.Update(r.spreadsheetID, rangeToUpdate, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("RAW").Do()

	if err != nil {
		log.Printf("[SHEETS] ERROR: Failed to update sheet: %v", err)
		return err
	}

	duration := time.Since(start)
	log.Printf("[SHEETS] SUCCESS: Expense saved to row %d in %v", rowNumber, duration)
	return nil
}