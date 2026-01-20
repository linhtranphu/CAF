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
	PaidDate string `json:"paidDate,omitempty"`
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

func (p *MessageParser) Parse(message string) (string, int64, error) {
	log.Printf("[AI] Parsing message: %s", message)
	
	if p.client == nil {
		log.Printf("[AI] No Gemini client available, returning basic parse")
		return strings.Title(strings.ToLower(message)), 1, nil
	}
	
	// Check cache first
	messageKey := strings.ToLower(strings.TrimSpace(message))
	if cached, exists := p.cache[messageKey]; exists {
		log.Printf("[AI] Cache hit for: %s", message)
		return cached.Items, cached.Amount, nil
	}

	prompt := `Parse Vietnamese expense message to JSON:

Message: "` + message + `"

Return only JSON:
{"items": "description", "amount": number_in_VND}

Rules:
- "triệu" = x1,000,000
- "k"/"nghìn" = x1,000
- "tỷ" = x1,000,000,000
- Remove amount from items

Examples:
"cọc nhà 34 triệu" → {"items": "Cọc nhà", "amount": 34000000}
"ăn trưa 150k" → {"items": "Ăn trưa", "amount": 150000}`

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
		return strings.Title(strings.ToLower(message)), 1, nil
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
		return strings.Title(strings.ToLower(message)), 1, nil
	}

	// Cache the result
	p.cache[messageKey] = resultData
	log.Printf("[AI] Cached result for: %s", message)

	log.Printf("[AI] Gemini result: items=%s, amount=%d", resultData.Items, resultData.Amount)
	return resultData.Items, resultData.Amount, nil
}