package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Ollama API structures for Google Colab integration
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// Hugging Face API structures for free inference
type HFRequest struct {
	Inputs     string                 `json:"inputs"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

type HFResponse []struct {
	GeneratedText string `json:"generated_text"`
}

type ChatRequest struct {
	Message string `json:"message"`
	Type    string `json:"type"` // "onboarding" or "policy_search"
}

type ChatResponse struct {
	Response    string       `json:"response"`
	Type        string       `json:"type"`
	PolicyFiles []PolicyFile `json:"policy_files,omitempty"`
}

type PolicyFile struct {
	Name        string `json:"name"`
	Content     string `json:"content"`
	Category    string `json:"category"`
	LastUpdated string `json:"last_updated"`
}

// Configuration for APIs
const (
	HF_API_URL = "https://api-inference.huggingface.co/models/microsoft/DialoGPT-medium"
)

// Mock policy database
var policies = []PolicyFile{
	{
		Name:        "Password Policy",
		Content:     "Passwords must be at least 12 characters long and include uppercase, lowercase, numbers, and special characters. Passwords must be changed every 90 days.",
		Category:    "Authentication",
		LastUpdated: "2024-01-15",
	},
	{
		Name:        "Data Classification Policy",
		Content:     "All company data must be classified as Public, Internal, Confidential, or Restricted. Confidential and Restricted data requires encryption at rest and in transit.",
		Category:    "Data Protection",
		LastUpdated: "2024-01-10",
	},
	{
		Name:        "Remote Work Security Policy",
		Content:     "Remote workers must use company-approved VPN, enable device encryption, and follow secure Wi-Fi practices. Personal devices require MDM enrollment.",
		Category:    "Remote Work",
		LastUpdated: "2024-01-20",
	},
	{
		Name:        "Incident Response Policy",
		Content:     "Security incidents must be reported within 2 hours. Follow the escalation matrix: L1 (Help Desk) -> L2 (Security Team) -> L3 (CISO). Document all actions taken.",
		Category:    "Incident Response",
		LastUpdated: "2024-01-05",
	},
}

// Function to call Ollama API (preferred - from Google Colab)
func callOllamaAPI(prompt string) (string, error) {
	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		return "", fmt.Errorf("OLLAMA_URL not configured")
	}

	// Create context-aware prompt for security chatbot
	securityPrompt := fmt.Sprintf(`You are an IT security assistant for company onboarding. 

Context: You help new employees understand security policies including passwords, VPN, data protection, and incident response.

Employee Question: %s

Provide a helpful, professional response about IT security. Keep it concise and actionable.`, prompt)

	requestBody := OllamaRequest{
		Model:  "llama3.1:8b",
		Prompt: securityPrompt,
		Stream: false,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", ollamaURL+"/api/generate", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ollamaResponse OllamaResponse
	err = json.Unmarshal(body, &ollamaResponse)
	if err != nil {
		return "", err
	}

	return ollamaResponse.Response, nil
}

// Function to call Hugging Face Inference API (fallback)
func callHuggingFaceAPI(prompt string) (string, error) {
	hfToken := os.Getenv("HF_TOKEN")
	if hfToken == "" {
		return "", fmt.Errorf("HF_TOKEN not configured")
	}

	requestBody := HFRequest{
		Inputs: prompt,
		Parameters: map[string]interface{}{
			"max_length":   200,
			"temperature":  0.7,
			"do_sample":    true,
			"pad_token_id": 50256,
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", HF_API_URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+hfToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var hfResponse HFResponse
	err = json.Unmarshal(body, &hfResponse)
	if err != nil {
		return "", err
	}

	if len(hfResponse) > 0 {
		return hfResponse[0].GeneratedText, nil
	}

	return "I'm sorry, I couldn't generate a response.", nil
}

// Smart LLM caller that tries Ollama first, then HF, then mock
func callLLM(prompt string) string {
	// Try Ollama first (Google Colab)
	if response, err := callOllamaAPI(prompt); err == nil {
		log.Println("✅ Using Ollama API from Google Colab")
		return response
	}

	// Fallback to Hugging Face
	if response, err := callHuggingFaceAPI(prompt); err == nil {
		log.Println("✅ Using Hugging Face API")
		return response
	}

	// Final fallback to mock responses
	log.Println("ℹ️  Using mock responses (no API configured)")
	return generateMockLLMResponse(prompt)
}

// Mock LLM response for testing without API
func generateMockLLMResponse(prompt string) string {
	prompt = strings.ToLower(prompt)

	if strings.Contains(prompt, "password") {
		return "Our password policy requires at least 12 characters with uppercase, lowercase, numbers, and special characters. Passwords must be changed every 90 days. Would you like me to show you the complete policy document?"
	}

	if strings.Contains(prompt, "vpn") {
		return "For remote work, you must use our company VPN. Make sure your device is encrypted and follow secure Wi-Fi practices. Personal devices need MDM enrollment."
	}

	if strings.Contains(prompt, "incident") {
		return "Security incidents must be reported within 2 hours. Follow our escalation process: Level 1 (Help Desk) → Level 2 (Security Team) → Level 3 (CISO). Document all actions taken."
	}

	if strings.Contains(prompt, "data") {
		return "All company data must be classified as Public, Internal, Confidential, or Restricted. Confidential and Restricted data requires encryption at rest and in transit."
	}

	return "I can help you with IT security questions including passwords, VPN access, data protection, and incident response. What would you like to know?"
}

func main() {
	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Routes
	r.POST("/api/chat", handleChat)
	r.GET("/api/policies", getPolicies)
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	log.Println("🚀 Security Chatbot Server starting on :8080...")
	log.Println("📝 Configuration:")

	if os.Getenv("OLLAMA_URL") != "" {
		log.Printf("   ✅ Ollama URL: %s", os.Getenv("OLLAMA_URL"))
	} else {
		log.Println("   ⚠️  OLLAMA_URL not set")
	}

	if os.Getenv("HF_TOKEN") != "" {
		log.Println("   ✅ Hugging Face token configured")
	} else {
		log.Println("   ⚠️  HF_TOKEN not set")
	}

	log.Println("   ℹ️  Mock responses available as fallback")

	r.Run(":8080")
}

func handleChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var response ChatResponse

	switch req.Type {
	case "onboarding":
		response = handleOnboardingWithLLM(req.Message)
	case "policy_search":
		response = handlePolicySearch(req.Message)
	default:
		response = ChatResponse{
			Response: "I can help you with IT security onboarding or policy searches. What would you like to know?",
			Type:     "general",
		}
	}

	c.JSON(http.StatusOK, response)
}

func handleOnboardingWithLLM(message string) ChatResponse {
	llmResponse := callLLM(message)

	return ChatResponse{
		Response: llmResponse,
		Type:     "onboarding",
	}
}

func handleOnboarding(message string) ChatResponse {
	// Simulate AI response for onboarding
	onboardingTopics := map[string]string{
		"password": "Let me help you understand our password policy. Passwords must be at least 12 characters long, include uppercase, lowercase, numbers, and special characters. They must be changed every 90 days. Would you like me to show you the complete password policy?",
		"vpn":      "For secure remote access, you'll need to use our company VPN. Here's what you need to know about connecting securely from remote locations...",
		"email":    "Email security is crucial. Always verify sender identity, be cautious of links and attachments, and report suspicious emails to the security team immediately.",
		"data":     "Data protection is everyone's responsibility. All data must be classified and handled according to our data classification policy. Let me explain the different classification levels...",
	}

	message = strings.ToLower(message)

	for keyword, resp := range onboardingTopics {
		if strings.Contains(message, keyword) {
			return ChatResponse{
				Response: resp,
				Type:     "onboarding",
			}
		}
	}

	return ChatResponse{
		Response: "Welcome to IT Security onboarding! I can help you with topics like passwords, VPN access, email security, and data protection. What would you like to learn about?",
		Type:     "onboarding",
	}
}

func handlePolicySearch(query string) ChatResponse {
	query = strings.ToLower(query)
	var matchedPolicies []PolicyFile
	var responseText string

	// Search for relevant policies
	for _, policy := range policies {
		if strings.Contains(strings.ToLower(policy.Name), query) ||
			strings.Contains(strings.ToLower(policy.Content), query) ||
			strings.Contains(strings.ToLower(policy.Category), query) {
			matchedPolicies = append(matchedPolicies, policy)
		}
	}

	if len(matchedPolicies) > 0 {
		responseText = fmt.Sprintf("I found %d policy document(s) related to your search. Here are the relevant policies:", len(matchedPolicies))
	} else {
		responseText = "I couldn't find any policies matching your search. Try searching for terms like 'password', 'data', 'remote work', or 'incident response'."
	}

	return ChatResponse{
		Response:    responseText,
		Type:        "policy_search",
		PolicyFiles: matchedPolicies,
	}
}

func getPolicies(c *gin.Context) {
	c.JSON(http.StatusOK, policies)
}
