package http

import (
	"bufio"
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/genai"
)

type SettingsHandler struct{}

func NewSettingsHandler() *SettingsHandler {
	return &SettingsHandler{}
}

type SettingsData struct {
	Message    string
	Success    bool
	HasAPIKey  bool
	CurrentKey string
}

func (h *SettingsHandler) ShowSettings(c *gin.Context) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	
	// Mask API key for display
	maskedKey := ""
	if apiKey != "" {
		if len(apiKey) > 10 {
			maskedKey = apiKey[:4] + "..." + apiKey[len(apiKey)-4:]
		} else {
			maskedKey = "***"
		}
	}

	data := SettingsData{
		HasAPIKey:  apiKey != "",
		CurrentKey: maskedKey,
	}

	tmpl, err := template.ParseFiles("templates/settings.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading template: "+err.Error())
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(c.Writer, data)
}

func (h *SettingsHandler) SaveSettings(c *gin.Context) {
	apiKey := strings.TrimSpace(c.PostForm("gemini_api_key"))

	if apiKey == "" {
		h.renderSettings(c, "⚠️ API key không được để trống", false)
		return
	}

	// Update .env file
	if err := updateEnvFile("GEMINI_API_KEY", apiKey); err != nil {
		h.renderSettings(c, "❌ Lỗi lưu file: "+err.Error(), false)
		return
	}

	// Update environment variable
	os.Setenv("GEMINI_API_KEY", apiKey)

	h.renderSettings(c, "✅ Đã lưu API key thành công! Vui lòng restart server để áp dụng.", true)
}

func (h *SettingsHandler) TestAPI(c *gin.Context) {
	var req struct {
		APIKey string `json:"api_key"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request"})
		return
	}

	// Test Gemini API
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  req.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error()})
		return
	}

	// Simple test
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-lite",
		genai.Text("Parse: ăn trưa 50k. Return JSON: {\"items\": \"Ăn trưa\", \"amount\": 50000}"),
		nil,
	)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"test_result": result.Text(),
	})
}

func (h *SettingsHandler) renderSettings(c *gin.Context, message string, success bool) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	
	maskedKey := ""
	if apiKey != "" {
		if len(apiKey) > 10 {
			maskedKey = apiKey[:4] + "..." + apiKey[len(apiKey)-4:]
		} else {
			maskedKey = "***"
		}
	}

	data := SettingsData{
		Message:    message,
		Success:    success,
		HasAPIKey:  apiKey != "",
		CurrentKey: maskedKey,
	}

	tmpl, err := template.ParseFiles("templates/settings.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading template: "+err.Error())
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(c.Writer, data)
}

func updateEnvFile(key, value string) error {
	envPath := ".env"
	
	// Read existing .env
	lines := []string{}
	file, err := os.Open(envPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		keyFound := false
		
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, key+"=") {
				lines = append(lines, key+"="+value)
				keyFound = true
			} else {
				lines = append(lines, line)
			}
		}
		
		if !keyFound {
			lines = append(lines, key+"="+value)
		}
		
		if err := scanner.Err(); err != nil {
			return err
		}
	}

	// Write back to .env
	content := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
		return err
	}

	log.Printf("[Settings] Updated %s in .env file", key)
	return nil
}
