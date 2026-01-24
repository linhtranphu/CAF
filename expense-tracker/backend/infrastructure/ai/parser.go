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

type MessageParser struct {
	client   *genai.Client
	cache    map[string]ExpenseData
	lastCall time.Time
}

type ExpenseData struct {
	Items    string `json:"items"`
	Amount   int64  `json:"amount"`
	Quantity string `json:"quantity,omitempty"`
	Unit     string `json:"unit,omitempty"`
	PaidDate string `json:"paidDate,omitempty"`
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

func NewMessageParser() *MessageParser {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Printf("[AI] No Gemini API key found in environment")
		return &MessageParser{
			client: nil,
			cache:  make(map[string]ExpenseData),
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
		}
	}
	
	log.Printf("[AI] Gemini client created successfully")
	return &MessageParser{
		client: client,
		cache:  make(map[string]ExpenseData),
	}
}

func (p *MessageParser) Parse(message string) (string, int64, string, string, time.Time, error) {
	log.Printf("[AI] Parsing message: %s", message)
	
	if p.client == nil {
		log.Printf("[AI] No Gemini client available, returning basic parse")
		return strings.Title(strings.ToLower(message)), 1, "", "", time.Now(), nil
	}
	
	// Check cache first
	messageKey := strings.ToLower(strings.TrimSpace(message))
	if cached, exists := p.cache[messageKey]; exists {
		log.Printf("[AI] Cache hit for: %s", message)
		parsedDate := parseDate(cached.PaidDate)
		return cached.Items, cached.Amount, cached.Quantity, cached.Unit, parsedDate, nil
	}

	currentDate := time.Now().Format("2006-01-02")
	prompt := `Parse Vietnamese expense message to JSON with quantity/unit extraction:

Current date: ` + currentDate + `
Message: "` + message + `"

Return only JSON:
{"items": "description", "amount": number_in_VND, "quantity": "number", "unit": "unit", "paidDate": "YYYY-MM-DD"}

Rules:
- "triệu" = x1,000,000
- "k"/"nghìn" = x1,000
- "tỷ" = x1,000,000,000
- Extract quantity and unit separately:
  * "2 cái bánh" → quantity: "2", unit: "cái"
  * "50kg gạo" → quantity: "50", unit: "kg"
  * "3 lon bia" → quantity: "3", unit: "lon"
  * "1 bịch kẹo" → quantity: "1", unit: "bịch"
  * "5 chai nước" → quantity: "5", unit: "chai"
- Remove amount, quantity, unit and date from items
- Date parsing:
  * "hôm nay" = current date
  * "hôm qua" = current date - 1 day
  * "hôm kia" = current date - 2 days
  * "tuần trước" = current date - 7 days
  * "tháng trước" = current date - 30 days
  * Specific dates: "22/1", "22 tháng 1", etc.
  * If no date mentioned, use current date

Examples:
"hôm qua mua 2 cái bánh 50k" → {"items": "Bánh", "amount": 50000, "quantity": "2", "unit": "cái", "paidDate": "2026-01-21"}
"50kg gạo 1.2 triệu" → {"items": "Gạo", "amount": 1200000, "quantity": "50", "unit": "kg", "paidDate": "2026-01-22"}
"ăn trưa 150k" → {"items": "Ăn trưa", "amount": 150000, "quantity": "", "unit": "", "paidDate": "2026-01-22"}`

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
		"gemini-2.5-flash-lite",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		log.Printf("[AI] Gemini API error: %v, using fallback", err)
		return strings.Title(strings.ToLower(message)), 1, "", "", time.Now(), nil
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
		return strings.Title(strings.ToLower(message)), 1, "", "", time.Now(), nil
	}

	// Cache the result
	p.cache[messageKey] = resultData
	log.Printf("[AI] Cached result for: %s", message)

	parsedDate := parseDate(resultData.PaidDate)
	log.Printf("[AI] Gemini result: items=%s, amount=%d, quantity=%s, unit=%s, date=%s", 
		resultData.Items, resultData.Amount, resultData.Quantity, resultData.Unit, parsedDate.Format("2006-01-02"))
	return resultData.Items, resultData.Amount, resultData.Quantity, resultData.Unit, parsedDate, nil
}