package ai

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/genai"
)

type APIKeyRepository interface {
	GetAPIKey() (string, error)
}

type MessageParser struct {
	client   *genai.Client
	cache    map[string]ExpenseData
	lastCall time.Time
	repo     APIKeyRepository
}

type ExpenseData struct {
	Items           string `json:"items"`
	Amount          int64  `json:"amount"`
	Quantity        string `json:"quantity,omitempty"`
	Unit            string `json:"unit,omitempty"`
	BaseQuantity    string `json:"baseQuantity,omitempty"`
	BaseUnit        string `json:"baseUnit,omitempty"`
	PaidDate        string `json:"paidDate,omitempty"`
	OriginalMessage string `json:"originalMessage,omitempty"`
}

// parseDate parses date string to time.Time
func parseDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Now()
	}
	
	// Try parsing YYYY-MM-DD format
	if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
		return parsed
	}
	
	// Try parsing DD/MM format (assume current year)
	if parsed, err := time.Parse("02/01", dateStr); err == nil {
		return time.Date(time.Now().Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.Local)
	}
	
	// Default to current time if parsing fails
	return time.Now()
}

func NewMessageParser(repo APIKeyRepository) *MessageParser {
	apiKey, _ := repo.GetAPIKey()
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}
	
	if apiKey == "" {
		log.Printf("[AI] No Gemini API key found")
		return &MessageParser{
			client: nil,
			cache:  make(map[string]ExpenseData),
			repo:   repo,
		}
	}
	
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Printf("[AI] Failed to create Gemini client: %v", err)
		return &MessageParser{
			client: nil,
			cache:  make(map[string]ExpenseData),
			repo:   repo,
		}
	}
	
	log.Printf("[AI] Gemini client created successfully")
	return &MessageParser{
		client: client,
		cache:  make(map[string]ExpenseData),
		repo:   repo,
	}
}

func (p *MessageParser) Parse(message string) (string, int64, string, string, string, string, string, time.Time, error) {
	log.Printf("[AI] Parsing message: %s", message)
	
	// Refresh client if needed
	if p.client == nil {
		apiKey, _ := p.repo.GetAPIKey()
		if apiKey != "" {
			ctx := context.Background()
			client, err := genai.NewClient(ctx, &genai.ClientConfig{
				APIKey:  apiKey,
				Backend: genai.BackendGeminiAPI,
			})
			if err == nil {
				p.client = client
				log.Printf("[AI] Gemini client created from MongoDB key")
			}
		}
	}
	
	if p.client == nil {
		log.Printf("[AI] No Gemini client available, returning basic parse")
		return strings.Title(strings.ToLower(message)), 1, "", "", "", "", message, time.Now(), nil
	}
	
	// Check cache first
	messageKey := strings.ToLower(strings.TrimSpace(message))
	if cached, exists := p.cache[messageKey]; exists {
		log.Printf("[AI] Cache hit for: %s", message)
		parsedDate := parseDate(cached.PaidDate)
		return cached.Items, cached.Amount, cached.Quantity, cached.Unit, cached.BaseQuantity, cached.BaseUnit, cached.OriginalMessage, parsedDate, nil
	}

	currentDate := time.Now().Format("2006-01-02")
	prompt := `Parse Vietnamese expense message to JSON with base unit conversion (ISO standard):

Current date: ` + currentDate + `
Message: "` + message + `"

Return ONLY valid JSON with this exact structure:
{"items": "description", "amount": number_in_VND, "quantity": "display_number", "unit": "display_unit", "baseQuantity": "base_number", "baseUnit": "iso_unit", "paidDate": "YYYY-MM-DD"}

IMPORTANT: You MUST include baseQuantity and baseUnit fields in your response!

Rules:
- "triệu" = x1,000,000
- "k"/"nghìn" = x1,000  
- "tỷ" = x1,000,000,000
- ALWAYS extract quantity/unit AND convert to base unit:
  * "2 bao cà phê 0.5kg" → quantity: "2", unit: "bao", baseQuantity: "1", baseUnit: "kg"
  * "50kg gạo" → quantity: "50", unit: "kg", baseQuantity: "50", baseUnit: "kg"
  * "500g thịt" → quantity: "500", unit: "g", baseQuantity: "0.5", baseUnit: "kg"
  * "2 lít dầu" → quantity: "2", unit: "lít", baseQuantity: "2", baseUnit: "L"
  * "3kg gạo" → quantity: "3", unit: "kg", baseQuantity: "3", baseUnit: "kg"

Base units (ISO): kg (mass), L (volume), m (length), pcs (count)
Conversions: 1000g=1kg, 1000ml=1L, 100cm=1m

Examples:
"2 bao cà phê 0.5kg 200k" → {"items": "Cà phê", "amount": 200000, "quantity": "2", "unit": "bao", "baseQuantity": "1", "baseUnit": "kg", "paidDate": "` + currentDate + `"}
"500g thịt 150k" → {"items": "Thịt", "amount": 150000, "quantity": "500", "unit": "g", "baseQuantity": "0.5", "baseUnit": "kg", "paidDate": "` + currentDate + `"}
"3kg gạo 180k" → {"items": "Gạo", "amount": 180000, "quantity": "3", "unit": "kg", "baseQuantity": "3", "baseUnit": "kg", "paidDate": "` + currentDate + `"}`

	log.Printf("[AI] Calling Gemini API...")
	
	// Rate limiting
	if time.Since(p.lastCall) < 1*time.Second {
		waitTime := 1*time.Second - time.Since(p.lastCall)
		log.Printf("[AI] Rate limiting: waiting %v", waitTime)
		time.Sleep(waitTime)
	}
	p.lastCall = time.Now()
	
	// Call Gemini API using SDK
	ctx := context.Background()
	result, err := p.client.Models.GenerateContent(
		ctx,
		"models/gemini-2.5-flash-lite",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		log.Printf("[AI] Gemini API error: %v, using fallback", err)
		return strings.Title(strings.ToLower(message)), 1, "", "", "", "", message, time.Now(), nil
	}

	responseText := result.Text()
	log.Printf("[AI] Gemini response: %s", responseText)

	// Clean response - remove markdown code blocks
	cleanResponse := strings.TrimSpace(responseText)
	if strings.HasPrefix(cleanResponse, "```json") {
		cleanResponse = strings.TrimPrefix(cleanResponse, "```json")
	}
	if strings.HasPrefix(cleanResponse, "```") {
		cleanResponse = strings.TrimPrefix(cleanResponse, "```")
	}
	if strings.HasSuffix(cleanResponse, "```") {
		cleanResponse = strings.TrimSuffix(cleanResponse, "```")
	}
	cleanResponse = strings.TrimSpace(cleanResponse)

	var resultData ExpenseData
	if err := json.Unmarshal([]byte(cleanResponse), &resultData); err != nil {
		log.Printf("[AI] JSON parse error: %v, response: %s, using fallback", err, cleanResponse)
		return strings.Title(strings.ToLower(message)), 1, "", "", "", "", message, time.Now(), nil
	}

	// Store original message
	resultData.OriginalMessage = message

	log.Printf("[AI] ===== AFTER JSON UNMARSHAL =====")
	log.Printf("[AI] resultData struct: %+v", resultData)
	log.Printf("[AI] Items: %s", resultData.Items)
	log.Printf("[AI] BaseQuantity: '%s'", resultData.BaseQuantity)
	log.Printf("[AI] BaseUnit: '%s'", resultData.BaseUnit)
	log.Printf("[AI] ====================================")

	// Cache the result
	p.cache[messageKey] = resultData
	log.Printf("[AI] Cached result for: %s", message)

	parsedDate := parseDate(resultData.PaidDate)
	log.Printf("[AI] Gemini result: items=%s, amount=%d, quantity=%s, unit=%s, baseQuantity=%s, baseUnit=%s, date=%s", 
		resultData.Items, resultData.Amount, resultData.Quantity, resultData.Unit, resultData.BaseQuantity, resultData.BaseUnit, parsedDate.Format("2006-01-02"))
	log.Printf("[AI] Full parsed data: %+v", resultData)
	return resultData.Items, resultData.Amount, resultData.Quantity, resultData.Unit, resultData.BaseQuantity, resultData.BaseUnit, resultData.OriginalMessage, parsedDate, nil
}