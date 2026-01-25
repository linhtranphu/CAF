package http

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/genai"
)

type SettingsRepository interface {
	SaveAPIKey(apiKey string) error
	GetAPIKey() (string, error)
}

type SettingsHandler struct {
	repo SettingsRepository
}

func NewSettingsHandler(repo SettingsRepository) *SettingsHandler {
	return &SettingsHandler{repo: repo}
}

type SettingsData struct {
	Message    string
	Success    bool
	HasAPIKey  bool
	CurrentKey string
}

func (h *SettingsHandler) ShowSettings(c *gin.Context) {
	apiKey, _ := h.repo.GetAPIKey()

	data := SettingsData{
		HasAPIKey:  apiKey != "",
		CurrentKey: apiKey,
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
	apiKey := c.PostForm("gemini_api_key")

	if apiKey == "" {
		h.renderSettings(c, "⚠️ API key không được để trống", false)
		return
	}

	// Save to MongoDB
	if err := h.repo.SaveAPIKey(apiKey); err != nil {
		h.renderSettings(c, "❌ Lỗi lưu: "+err.Error(), false)
		return
	}

	h.renderSettings(c, "✅ Đã lưu API key thành công! Không cần restart server.", true)
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
		"models/gemini-2.5-flash-lite",
		genai.Text("Say hello"),
		nil,
	)
	if err != nil {
		log.Printf("[Settings] Test API error: %v", err)
		c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"test_result": result.Text(),
	})
}

func (h *SettingsHandler) renderSettings(c *gin.Context, message string, success bool) {
	apiKey, _ := h.repo.GetAPIKey()

	data := SettingsData{
		Message:    message,
		Success:    success,
		HasAPIKey:  apiKey != "",
		CurrentKey: apiKey,
	}

	tmpl, err := template.ParseFiles("templates/settings.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading template: "+err.Error())
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(c.Writer, data)
}
