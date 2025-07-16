package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ledongthuc/pdf"
	"github.com/unidoc/unioffice/document"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

// Enhanced PolicyFile structure for better document management with GORM tags
// User model for authentication and authorization
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null;size:50"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null;size:100"`
	Password  string    `json:"-" gorm:"not null;size:255"` // Hidden from JSON responses
	FirstName string    `json:"first_name" gorm:"not null;size:50"`
	LastName  string    `json:"last_name" gorm:"not null;size:50"`
	Role      string    `json:"role" gorm:"not null;size:20;default:'user'"` // "user", "admin", "hr", "it_security"
	IsActive  bool      `json:"is_active" gorm:"default:true;index"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// UserRole constants
const (
	RoleUser       = "user"
	RoleAdmin      = "admin"
	RoleHR         = "hr"
	RoleITSecurity = "it_security"
)

// AuditLog model for tracking system activities
type AuditLog struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       uint      `json:"user_id" gorm:"not null;index"`
	User         User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Action       string    `json:"action" gorm:"not null;size:50;index"` // CREATE, UPDATE, DELETE, VIEW, LOGIN, etc.
	ResourceType string    `json:"resource_type" gorm:"not null;size:50;index"` // USER, DOCUMENT, SYSTEM
	ResourceID   *uint     `json:"resource_id,omitempty" gorm:"index"` // ID of the affected resource
	ResourceName string    `json:"resource_name,omitempty" gorm:"size:255"` // Name of the affected resource
	Details      string    `json:"details,omitempty" gorm:"type:text"` // Additional details about the action
	IPAddress    string    `json:"ip_address,omitempty" gorm:"size:45"` // Support IPv6
	UserAgent    string    `json:"user_agent,omitempty" gorm:"type:text"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// Audit action constants
const (
	ActionCreate = "CREATE"
	ActionUpdate = "UPDATE"
	ActionDelete = "DELETE"
	ActionView   = "VIEW"
	ActionLogin  = "LOGIN"
	ActionLogout = "LOGOUT"
)

// Resource type constants
const (
	ResourceUser     = "USER"
	ResourceDocument = "DOCUMENT"
	ResourceSystem   = "SYSTEM"
)

// PolicyFile model (updated to include user relationship)
type PolicyFile struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"not null;size:255;index"`
	Content     string    `json:"content" gorm:"type:text;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Category    string    `json:"category" gorm:"not null;size:100;index"`
	DocumentType string   `json:"document_type" gorm:"not null;size:50;index"` // "policy" or "onboarding"
	Tags        string    `json:"-" gorm:"type:text"` // Store as JSON string in DB
	TagsArray   []string  `json:"tags" gorm:"-"` // For JSON response
	FilePath    string    `json:"file_path,omitempty" gorm:"size:500"`
	CreatedBy   string    `json:"created_by" gorm:"size:100"` // Will be updated to use User ID in future
	CreatedByUserID *uint `json:"created_by_user_id,omitempty" gorm:"index"` // Foreign key to User
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	LastUpdated string    `json:"last_updated" gorm:"-"` // Computed field for compatibility
	IsActive    bool      `json:"is_active" gorm:"default:true;index"`
}

// Request structures for document management
type CreateDocumentRequest struct {
	Name         string   `json:"name" binding:"required"`
	Content      string   `json:"content" binding:"required"`
	Description  string   `json:"description"`
	Category     string   `json:"category" binding:"required"`
	DocumentType string   `json:"document_type" binding:"required"`
	Tags         []string `json:"tags"`
	CreatedBy    string   `json:"created_by"`
	FilePath     string   `json:"file_path,omitempty"` // Path to original uploaded file
}

type UpdateDocumentRequest struct {
	Name         string   `json:"name"`
	Content      string   `json:"content"`
	Description  string   `json:"description"`
	Category     string   `json:"category"`
	DocumentType string   `json:"document_type"`
	Tags         []string `json:"tags"`
	IsActive     *bool    `json:"is_active"`
}

// Authentication request structures
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required,min=2,max=50"`
	LastName  string `json:"last_name" binding:"required,min=2,max=50"`
}

type AuthResponse struct {
	Token     string    `json:"token"`
	User      UserInfo  `json:"user"`
	ExpiresAt time.Time `json:"expires_at"`
}

type UserInfo struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// JWT Claims structure
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Enhanced search structures
type DocumentMatch struct {
	Document PolicyFile `json:"document"`
	Score    float64    `json:"score"`
	Matches  []Match    `json:"matches"`
}

type Match struct {
	Field   string `json:"field"`
	Text    string `json:"text"`
	Score   float64 `json:"score"`
}

type SearchEngine struct {
	Documents []PolicyFile
	Index     map[string][]DocumentIndex // word -> document indices
}

type DocumentIndex struct {
	DocumentID int
	Field      string
	Frequency  int
	Positions  []int
}

// File upload structures
type FileUploadResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	FileName    string `json:"file_name,omitempty"`
	FileType    string `json:"file_type,omitempty"`
	FileSize    int64  `json:"file_size,omitempty"`
	FilePath    string `json:"file_path,omitempty"`
	ExtractedText string `json:"extracted_text,omitempty"`
	Error       string `json:"error,omitempty"`
}

type SupportedFileType struct {
	Extension   string   `json:"extension"`
	MimeTypes   []string `json:"mime_types"`
	Description string   `json:"description"`
	MaxSize     int64    `json:"max_size_mb"`
}

// Configuration for APIs
const (
	HF_API_URL = "https://api-inference.huggingface.co/models/microsoft/DialoGPT-medium"
)

// JWT Configuration
var jwtSecret = []byte(getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"))

// Global database instance
var db *gorm.DB

// Database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Authentication utility functions
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateJWT(user User) (string, time.Time, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours
	
	claims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "security-chatbot",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	
	return tokenString, expirationTime, err
}

func validateJWT(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	
	return claims, nil
}

func userToUserInfo(user User) UserInfo {
	return UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		IsActive:  user.IsActive,
		LastLogin: user.LastLogin,
		CreatedAt: user.CreatedAt,
	}
}

// Audit logging helper functions
func logAuditActivity(c *gin.Context, userID uint, action, resourceType string, resourceID *uint, resourceName, details string) {
	// Get client IP and User-Agent
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	auditLog := AuditLog{
		UserID:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		ResourceName: resourceName,
		Details:      details,
		IPAddress:    clientIP,
		UserAgent:    userAgent,
	}

	// Log to database (non-blocking)
	go func() {
		if err := db.Create(&auditLog).Error; err != nil {
			log.Printf("Failed to create audit log: %v", err)
		}
	}()
}

// Convenience functions for specific audit types
func logDocumentActivity(c *gin.Context, userID uint, action string, document *PolicyFile, details string) {
	var resourceID *uint
	var resourceName string
	
	if document != nil {
		resourceID = &document.ID
		resourceName = document.Name
	}
	
	logAuditActivity(c, userID, action, ResourceDocument, resourceID, resourceName, details)
}

func logUserActivity(c *gin.Context, actorUserID uint, action string, targetUser *User, details string) {
	var resourceID *uint
	var resourceName string
	
	if targetUser != nil {
		resourceID = &targetUser.ID
		resourceName = targetUser.Username
	}
	
	logAuditActivity(c, actorUserID, action, ResourceUser, resourceID, resourceName, details)
}

func logSystemActivity(c *gin.Context, userID uint, action, details string) {
	logAuditActivity(c, userID, action, ResourceSystem, nil, "", details)
}

// Helper methods for PolicyFile
func (p *PolicyFile) BeforeSave(tx *gorm.DB) error {
	// Convert TagsArray to JSON string for database storage
	if len(p.TagsArray) > 0 {
		tagsJSON, err := json.Marshal(p.TagsArray)
		if err != nil {
			return err
		}
		p.Tags = string(tagsJSON)
	}
	return nil
}

func (p *PolicyFile) AfterFind(tx *gorm.DB) error {
	// Convert JSON string back to TagsArray
	if p.Tags != "" {
		err := json.Unmarshal([]byte(p.Tags), &p.TagsArray)
		if err != nil {
			p.TagsArray = []string{}
		}
	}
	// Set LastUpdated for compatibility
	p.LastUpdated = p.UpdatedAt.Format("2006-01-02")
	return nil
}

// Database connection function
func connectDB() (*gorm.DB, error) {
	config := DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "chatbot_user"),
		Password: getEnv("DB_PASSWORD", "chatbot_password"),
		DBName:   getEnv("DB_NAME", "chatbot_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}

	// Auto-migrate the schema
	err = database.AutoMigrate(&User{}, &PolicyFile{}, &AuditLog{})
	if err != nil {
		return nil, err
	}

	return database, nil
}

// Get environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// File upload and processing utilities

// Save uploaded file to disk and return the file path
func saveUploadedFile(fileHeader *multipart.FileHeader) (string, error) {
	// Create uploads directory if it doesn't exist
	uploadsDir := "uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create uploads directory: %v", err)
	}

	// Generate unique filename to avoid conflicts
	ext := filepath.Ext(fileHeader.Filename)
	baseName := strings.TrimSuffix(fileHeader.Filename, ext)
	timestamp := time.Now().Format("20060102_150405")
	uniqueFilename := fmt.Sprintf("%s_%s%s", baseName, timestamp, ext)
	
	// Clean filename to avoid path traversal issues
	uniqueFilename = filepath.Base(uniqueFilename)
	filePath := filepath.Join(uploadsDir, uniqueFilename)

	// Open uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return filePath, nil
}

// Get supported file types
func getSupportedFileTypes() []SupportedFileType {
	return []SupportedFileType{
		{
			Extension:   ".pdf",
			MimeTypes:   []string{"application/pdf"},
			Description: "PDF Document",
			MaxSize:     10, // 10MB
		},
		{
			Extension:   ".docx",
			MimeTypes:   []string{"application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
			Description: "Microsoft Word Document",
			MaxSize:     10, // 10MB
		},
		{
			Extension:   ".txt",
			MimeTypes:   []string{"text/plain"},
			Description: "Text File",
			MaxSize:     5, // 5MB
		},
		{
			Extension:   ".md",
			MimeTypes:   []string{"text/markdown", "text/plain"},
			Description: "Markdown File",
			MaxSize:     5, // 5MB
		},
	}
}

// Extract text from uploaded file based on file type
func extractTextFromFile(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Read file content into memory
	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// Detect MIME type
	mtype := mimetype.Detect(content)
	
	switch mtype.String() {
	case "application/pdf":
		return extractTextFromPDF(content)
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return extractTextFromDocx(content)
	case "text/plain":
		return cleanText(string(content)), nil
	default:
		// Try to treat as plain text if it's readable
		if isValidUTF8(content) {
			return cleanText(string(content)), nil
		}
		return "", fmt.Errorf("unsupported file type: %s", mtype.String())
	}
}

// cleanText removes null bytes and other problematic characters for database storage
func cleanText(text string) string {
	// Remove null bytes and other control characters except newlines, tabs, and carriage returns
	cleaned := strings.Map(func(r rune) rune {
		if r == 0 || (r < 32 && r != '\n' && r != '\t' && r != '\r') {
			return -1 // Remove character
		}
		return r
	}, text)
	
	// Ensure the text is valid UTF-8
	if !utf8.ValidString(cleaned) {
		// Convert to valid UTF-8, replacing invalid sequences
		cleaned = strings.ToValidUTF8(cleaned, "")
	}
	
	return cleaned
}

// Extract text from PDF file
func extractTextFromPDF(content []byte) (string, error) {
	reader := bytes.NewReader(content)
	pdfReader, err := pdf.NewReader(reader, int64(len(content)))
	if err != nil {
		return "", fmt.Errorf("failed to parse PDF: %v", err)
	}

	var text strings.Builder
	totalPages := pdfReader.NumPage()

	for i := 1; i <= totalPages; i++ {
		page := pdfReader.Page(i)
		if page.V.IsNull() {
			continue
		}

		// GetPlainText requires a font map parameter
		fonts := make(map[string]*pdf.Font)
		pageText, err := page.GetPlainText(fonts)
		if err != nil {
			log.Printf("Warning: failed to extract text from page %d: %v", i, err)
			continue
		}

		// Clean the extracted text before adding it
		cleanPageText := cleanText(pageText)
		text.WriteString(cleanPageText)
		text.WriteString("\n")
	}

	result := text.String()
	if strings.TrimSpace(result) == "" {
		return "", fmt.Errorf("no text content found in PDF")
	}

	return result, nil
}

// Extract text from DOCX file
func extractTextFromDocx(content []byte) (string, error) {
	reader := bytes.NewReader(content)
	doc, err := document.Read(reader, int64(len(content)))
	if err != nil {
		return "", fmt.Errorf("failed to parse DOCX: %v", err)
	}
	defer doc.Close()

	var text strings.Builder
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			// Clean the text before adding it
			cleanRunText := cleanText(run.Text())
			text.WriteString(cleanRunText)
		}
		text.WriteString("\n")
	}

	result := text.String()
	if strings.TrimSpace(result) == "" {
		return "", fmt.Errorf("no text content found in DOCX")
	}

	return result, nil
}

// Check if content is valid UTF-8 text
func isValidUTF8(content []byte) bool {
	// Check if content is valid UTF-8
	if !isUTF8(content) {
		return false
	}

	// Check if content contains mostly printable characters
	printableCount := 0
	totalCount := 0
	
	for _, r := range string(content) {
		totalCount++
		if r >= 32 && r <= 126 || r == '\n' || r == '\r' || r == '\t' {
			printableCount++
		}
	}

	// Consider it text if at least 80% of characters are printable
	return totalCount > 0 && float64(printableCount)/float64(totalCount) >= 0.8
}

// Simple UTF-8 validation
func isUTF8(data []byte) bool {
	i := 0
	for i < len(data) {
		r := data[i]
		
		if r < 128 {
			i++
			continue
		}
		
		// Multi-byte character
		var size int
		if r>>5 == 0b110 {
			size = 2
		} else if r>>4 == 0b1110 {
			size = 3
		} else if r>>3 == 0b11110 {
			size = 4
		} else {
			return false
		}
		
		if i+size > len(data) {
			return false
		}
		
		for j := 1; j < size; j++ {
			if data[i+j]>>6 != 0b10 {
				return false
			}
		}
		
		i += size
	}
	
	return true
}

// Validate file upload
func validateFileUpload(fileHeader *multipart.FileHeader) error {
	// Check file size (max 10MB)
	maxSize := int64(10 * 1024 * 1024) // 10MB in bytes
	if fileHeader.Size > maxSize {
		return fmt.Errorf("file size (%d bytes) exceeds maximum allowed size (%d bytes)", fileHeader.Size, maxSize)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	supportedTypes := getSupportedFileTypes()
	
	for _, fileType := range supportedTypes {
		if fileType.Extension == ext {
			return nil
		}
	}

	return fmt.Errorf("unsupported file type: %s", ext)
}

// Initialize database with sample data
func initializeDatabase() error {
	// Check if users already exist
	var userCount int64
	db.Model(&User{}).Count(&userCount)
	
	// Create default admin user if no users exist
	var adminUser User
	if userCount == 0 {
		hashedPassword, err := hashPassword("SecureAdmin123!")
		if err != nil {
			return fmt.Errorf("failed to hash admin password: %v", err)
		}
		
		adminUser = User{
			Username:  "admin",
			Email:     "admin@company.com",
			Password:  hashedPassword,
			FirstName: "System",
			LastName:  "Administrator",
			Role:      RoleAdmin,
			IsActive:  true,
		}
		
		if err := db.Create(&adminUser).Error; err != nil {
			return fmt.Errorf("failed to create admin user: %v", err)
		}
		
		log.Println("✅ Created default admin user:")
		log.Println("   Username: admin")
		log.Println("   Password: SecureAdmin123!")
		log.Println("   ⚠️  Please change the default password after first login!")
	} else {
		// Get the first admin user
		if err := db.Where("role = ?", RoleAdmin).First(&adminUser).Error; err != nil {
			// If no admin exists, get any user
			if err := db.First(&adminUser).Error; err != nil {
				return fmt.Errorf("failed to find any user: %v", err)
			}
		}
	}

	// Check if policy data already exists
	var policyCount int64
	db.Model(&PolicyFile{}).Count(&policyCount)
	if policyCount > 0 {
		log.Println("Database already contains policy data, skipping policy initialization")
		return nil
	}

	// Sample policy data
	samplePolicies := []PolicyFile{
	{
		Name:         "Password Policy",
		Content:      "Passwords must be at least 12 characters long and include uppercase, lowercase, numbers, and special characters. Passwords must be changed every 90 days.",
		Description:  "Comprehensive password requirements for all company accounts",
		Category:     "Authentication",
		DocumentType: "policy",
		TagsArray:    []string{"password", "security", "authentication", "compliance"},
		CreatedBy:    "IT Security Team",
		CreatedByUserID: &adminUser.ID,
		IsActive:     true,
	},
	{
		Name:         "Data Classification Policy",
		Content:      "All company data must be classified as Public, Internal, Confidential, or Restricted. Confidential and Restricted data requires encryption at rest and in transit.",
		Description:  "Guidelines for classifying and protecting company data",
		Category:     "Data Protection",
		DocumentType: "policy",
		TagsArray:    []string{"data", "classification", "encryption", "confidential"},
		CreatedBy:    "Data Protection Officer",
		CreatedByUserID: &adminUser.ID,
		IsActive:     true,
	},
	{
		Name:         "Remote Work Security Policy",
		Content:      "Remote workers must use company-approved VPN, enable device encryption, and follow secure Wi-Fi practices. Personal devices require MDM enrollment.",
		Description:  "Security requirements for remote work arrangements",
		Category:     "Remote Work",
		DocumentType: "policy",
		TagsArray:    []string{"remote", "vpn", "encryption", "mdm", "wifi"},
		CreatedBy:    "IT Operations",
		CreatedByUserID: &adminUser.ID,
		IsActive:     true,
	},
	{
		Name:         "Incident Response Policy",
		Content:      "Security incidents must be reported within 2 hours. Follow the escalation matrix: L1 (Help Desk) -> L2 (Security Team) -> L3 (CISO). Document all actions taken.",
		Description:  "Procedures for reporting and handling security incidents",
		Category:     "Incident Response",
		DocumentType: "policy",
		TagsArray:    []string{"incident", "response", "escalation", "security", "reporting"},
		CreatedBy:    "CISO Office",
		CreatedByUserID: &adminUser.ID,
		IsActive:     true,
	},
	{
		Name:         "New Employee Security Onboarding",
		Content:      "Welcome to the company! This guide covers essential security practices including password setup, VPN configuration, email security awareness, and device encryption. Please complete all steps within your first week.",
		Description:  "Complete security onboarding checklist for new employees",
		Category:     "Onboarding",
		DocumentType: "onboarding",
		TagsArray:    []string{"onboarding", "new-employee", "checklist", "setup"},
		CreatedBy:    "HR Security Team",
		CreatedByUserID: &adminUser.ID,
		IsActive:     true,
	},
	{
		Name:         "VPN Setup Guide",
		Content:      "Step-by-step instructions for configuring the company VPN on Windows, Mac, and mobile devices. Includes troubleshooting common connection issues.",
		Description:  "Technical guide for VPN setup and configuration",
		Category:     "Technical Guides",
		DocumentType: "onboarding",
		TagsArray:    []string{"vpn", "setup", "configuration", "troubleshooting", "guide"},
		CreatedBy:    "IT Help Desk",
		CreatedByUserID: &adminUser.ID,
		IsActive:     true,
	},
	}

	// Insert sample data into database
	for _, policy := range samplePolicies {
		if err := db.Create(&policy).Error; err != nil {
			log.Printf("Error inserting policy %s: %v", policy.Name, err)
			return err
		}
	}

	log.Printf("Successfully initialized database with %d documents", len(samplePolicies))
	return nil
}

// Authentication middleware
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Validate JWT token
		claims, err := validateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if user still exists and is active
		var user User
		if err := db.Where("id = ? AND is_active = ?", claims.UserID, true).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found or inactive"})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)
		c.Set("user", user)

		c.Next()
	}
}

// Role-based access control middleware
func requireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		roleStr := userRole.(string)
		for _, role := range roles {
			if roleStr == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}

// Initialize search engine with database
func NewSearchEngine() *SearchEngine {
	var documents []PolicyFile
	db.Find(&documents)
	
	engine := &SearchEngine{
		Documents: documents,
		Index:     make(map[string][]DocumentIndex),
	}
	engine.BuildIndex()
	return engine
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
	// Check if AI features are enabled
	aiEnabled := os.Getenv("AI_ENABLED")
	if aiEnabled != "true" {
		log.Println("🤖 AI features disabled, using mock responses")
		return generateMockLLMResponse(prompt)
	}

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

// Authentication handlers
func handleRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username or email already exists
	var existingUser User
	if err := db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	// Validate password strength (following their own security policy!)
	if len(req.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
		return
	}

	// Hash password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      RoleUser, // Default role
		IsActive:  true,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Log user registration activity
	logSystemActivity(c, user.ID, ActionCreate, fmt.Sprintf("New user registered: %s %s (%s)", user.FirstName, user.LastName, user.Username))

	// Generate JWT token
	token, expiresAt, err := generateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return auth response
	response := AuthResponse{
		Token:     token,
		User:      userToUserInfo(user),
		ExpiresAt: expiresAt,
	}

	c.JSON(http.StatusCreated, response)
}

func handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by username or email
	var user User
	if err := db.Where("(username = ? OR email = ?) AND is_active = ?", req.Username, req.Username, true).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if !checkPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	db.Save(&user)

	// Generate JWT token
	token, expiresAt, err := generateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Log successful login
	logSystemActivity(c, user.ID, ActionLogin, fmt.Sprintf("User %s logged in successfully", user.Username))

	// Return auth response
	response := AuthResponse{
		Token:     token,
		User:      userToUserInfo(user),
		ExpiresAt: expiresAt,
	}

	c.JSON(http.StatusOK, response)
}

func handleProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
		return
	}

	userInfo := userToUserInfo(user.(User))
	c.JSON(http.StatusOK, gin.H{"user": userInfo})
}

func handleUpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update user
	updates := make(map[string]interface{})
	if req.FirstName != "" {
		updates["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		updates["last_name"] = req.LastName
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}

	var user User
	if err := db.Model(&user).Where("id = ?", userID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	// Fetch updated user
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": userToUserInfo(user)})
}

func handleChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current user
	var user User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Verify current password
	if !checkPassword(req.CurrentPassword, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is incorrect"})
		return
	}

	// Hash new password
	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	// Update password
	if err := db.Model(&user).Update("password", hashedPassword).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// User Management Handlers (Admin only)

func handleGetAllUsers(c *gin.Context) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	var userInfos []UserInfo
	for _, user := range users {
		userInfos = append(userInfos, userToUserInfo(user))
	}

	// Log user list access activity
	currentUserID, _ := c.Get("user_id")
	logSystemActivity(c, currentUserID.(uint), ActionView, fmt.Sprintf("Accessed user management list (%d users)", len(users)))

	c.JSON(http.StatusOK, gin.H{"users": userInfos})
}

func handleUpdateUser(c *gin.Context) {
	userID := c.Param("id")
	
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		IsActive  *bool  `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.FirstName != "" {
		updates["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		updates["last_name"] = req.LastName
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	// Update user
	var user User
	if err := db.Model(&user).Where("id = ?", userID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Fetch updated user
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated user"})
		return
	}

	// Log user update activity
	currentUserID, _ := c.Get("user_id")
	logUserActivity(c, currentUserID.(uint), ActionUpdate, &user, fmt.Sprintf("Updated user details for %s %s (ID: %d)", user.FirstName, user.LastName, user.ID))

	c.JSON(http.StatusOK, gin.H{"user": userToUserInfo(user)})
}

func handleUpdateUserRole(c *gin.Context) {
	userID := c.Param("id")
	
	var req struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate role
	validRoles := []string{RoleUser, RoleAdmin, RoleHR, RoleITSecurity}
	roleValid := false
	for _, validRole := range validRoles {
		if req.Role == validRole {
			roleValid = true
			break
		}
	}

	if !roleValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	// Update user role
	var user User
	if err := db.Model(&user).Where("id = ?", userID).Update("role", req.Role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}

	// Fetch updated user
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated user"})
		return
	}

	// Log user role update activity
	currentUserID, _ := c.Get("user_id")
	logUserActivity(c, currentUserID.(uint), ActionUpdate, &user, fmt.Sprintf("Changed role to '%s' for user %s %s (ID: %d)", req.Role, user.FirstName, user.LastName, user.ID))

	c.JSON(http.StatusOK, gin.H{"user": userToUserInfo(user)})
}

func handleDeleteUser(c *gin.Context) {
	userID := c.Param("id")
	
	// Check if trying to delete self
	currentUserID, _ := c.Get("user_id")
	if userID == fmt.Sprintf("%v", currentUserID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete your own account"})
		return
	}

	// Fetch user details before deactivation for audit log
	var user User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Soft delete by setting is_active = false
	if err := db.Model(&user).Where("id = ?", userID).Update("is_active", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate user"})
		return
	}

	// Log user deactivation activity
	logUserActivity(c, currentUserID.(uint), ActionDelete, &user, fmt.Sprintf("Deactivated user %s %s (ID: %d)", user.FirstName, user.LastName, user.ID))

	c.JSON(http.StatusOK, gin.H{"message": "User deactivated successfully"})
}

// Audit Log Handlers (Admin only)

func handleGetAuditLogs(c *gin.Context) {
	// Parse query parameters for filtering
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 200 {
			limit = l
		}
	}

	offset := (page - 1) * limit

	// Build query
	query := db.Model(&AuditLog{}).Preload("User")

	// Apply filters
	if action := c.Query("action"); action != "" {
		query = query.Where("action = ?", action)
	}
	if resourceType := c.Query("resource_type"); resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if from := c.Query("from"); from != "" {
		if fromTime, err := time.Parse("2006-01-02", from); err == nil {
			query = query.Where("created_at >= ?", fromTime)
		}
	}
	if to := c.Query("to"); to != "" {
		if toTime, err := time.Parse("2006-01-02", to); err == nil {
			query = query.Where("created_at <= ?", toTime.Add(24*time.Hour))
		}
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get audit logs
	var auditLogs []AuditLog
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&auditLogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch audit logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"audit_logs": auditLogs,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func main() {
	// Initialize database connection
	var err error
	db, err = connectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	log.Println("Successfully connected to PostgreSQL database")
	
	// Initialize database with sample data
	if err := initializeDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	
	// Initialize search engine with database data
	_ = NewSearchEngine() // Initialize for testing, search engines are created fresh for each request
	
	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	
	// Default allowed origins for development
	allowedOrigins := []string{
		"http://localhost:3000",     // Development frontend
		"http://frontend:3000",      // Docker container frontend
		"http://127.0.0.1:3000",     // Alternative localhost
	}
	
	// Add production frontend URL if specified
	if frontendURL := getEnv("FRONTEND_URL", ""); frontendURL != "" {
		allowedOrigins = append(allowedOrigins, frontendURL)
	}
	
	config.AllowOrigins = allowedOrigins
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.ExposeHeaders = []string{"Content-Disposition", "Content-Type", "Content-Length"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Public routes (no authentication required)
	r.POST("/api/auth/login", handleLogin)
	r.POST("/api/auth/register", handleRegister)
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Protected routes (authentication required)
	authenticated := r.Group("/api")
	authenticated.Use(authMiddleware())
	{
		// User profile routes
		authenticated.GET("/profile", handleProfile)
		authenticated.PUT("/profile", handleUpdateProfile)
		authenticated.POST("/change-password", handleChangePassword)

		// Chat routes (all users can chat)
		authenticated.POST("/chat", handleChat)
		authenticated.GET("/policies", getPolicies)

		// Document viewing (all authenticated users)
		authenticated.GET("/documents", getDocuments)
		authenticated.GET("/documents/:id", getDocumentByID)
		authenticated.GET("/documents/:id/download", downloadDocument)
		authenticated.GET("/documents/search", searchDocuments)
	}

	// Admin-only routes (document and user management)
	adminOnly := r.Group("/api")
	adminOnly.Use(authMiddleware(), requireRole(RoleAdmin, RoleITSecurity))
	{
		// Document management
		adminOnly.POST("/documents", createDocument)
		adminOnly.PUT("/documents/:id", updateDocument)
		adminOnly.DELETE("/documents/:id", deleteDocument)
		
		// File upload endpoints
		adminOnly.POST("/upload", handleFileUpload)
		adminOnly.GET("/upload/supported-types", handleGetSupportedFileTypes)

		// User management (admin only)
		adminOnly.GET("/users", handleGetAllUsers)
		adminOnly.PUT("/users/:id", handleUpdateUser)
		adminOnly.PUT("/users/:id/role", handleUpdateUserRole)
		adminOnly.DELETE("/users/:id", handleDeleteUser)

		// Audit logs (admin only)
		adminOnly.GET("/audit-logs", handleGetAuditLogs)
	}

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

	if os.Getenv("JWT_SECRET") != "" {
		log.Println("   ✅ JWT secret configured")
	} else {
		log.Println("   ⚠️  JWT_SECRET not set (using default - change in production!)")
	}

	log.Println("   ℹ️  Mock responses available as fallback")
	log.Println("🔐 Authentication enabled:")
	log.Println("   📋 Public routes: /api/auth/login, /api/auth/register, /api/health")
	log.Println("   🔒 Protected routes: /api/chat, /api/policies, /api/documents (view)")
	log.Println("   👑 Admin-only routes: /api/documents (create/edit/delete)")

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

	// Log chat activity with document access
	userID, _ := c.Get("user_id")
	if len(response.PolicyFiles) > 0 {
		var docNames []string
		for _, doc := range response.PolicyFiles {
			docNames = append(docNames, doc.Name)
		}
		logSystemActivity(c, userID.(uint), ActionView, fmt.Sprintf("Chat search '%s' returned %d documents: %s", req.Message, len(response.PolicyFiles), strings.Join(docNames, ", ")))
	} else {
		logSystemActivity(c, userID.(uint), ActionView, fmt.Sprintf("Chat search '%s' (type: %s) - no documents returned", req.Message, req.Type))
	}

	c.JSON(http.StatusOK, response)
}

func handleOnboardingWithLLM(message string) ChatResponse {
	llmResponse := callLLM(message)

	// Use enhanced search engine to find relevant documents from database
	searchEngine := NewSearchEngine()
	matches := searchEngine.Search(message, 5) // Limit to top 5 for onboarding
	
	var matchedPolicies []PolicyFile
	for _, match := range matches {
		// Prioritize onboarding documents, but include relevant policies too
		matchedPolicies = append(matchedPolicies, match.Document)
	}

	return ChatResponse{
		Response:    llmResponse,
		Type:        "onboarding",
		PolicyFiles: matchedPolicies,
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
	// Use enhanced search engine with database data
	searchEngine := NewSearchEngine()
	matches := searchEngine.Search(query, 10)
	
	var matchedPolicies []PolicyFile
	for _, match := range matches {
		matchedPolicies = append(matchedPolicies, match.Document)
	}

	var responseText string
	if len(matchedPolicies) > 0 {
		responseText = fmt.Sprintf("I found %d document(s) related to your search with relevance scoring. Here are the most relevant documents:", len(matchedPolicies))
	} else {
		responseText = "I couldn't find any documents matching your search. Try searching for terms like 'password', 'data', 'remote work', 'onboarding', or 'incident response'."
	}

	return ChatResponse{
		Response:    responseText,
		Type:        "policy_search",
		PolicyFiles: matchedPolicies,
	}
}



func getPolicies(c *gin.Context) {
	var documents []PolicyFile
	if err := db.Find(&documents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch policies"})
		return
	}

	// Log policy access activity
	userID, _ := c.Get("user_id")
	logSystemActivity(c, userID.(uint), ActionView, fmt.Sprintf("Accessed all policies (%d documents)", len(documents)))

	c.JSON(http.StatusOK, documents)
}

// Document CRUD handlers

// Get all documents with optional filtering
func getDocuments(c *gin.Context) {
	documentType := c.Query("type")     // "policy" or "onboarding"
	category := c.Query("category")     // filter by category
	activeOnly := c.Query("active")     // "true" to show only active documents

	var documents []PolicyFile
	query := db.Model(&PolicyFile{})

	// Apply filters
	if activeOnly == "true" {
		query = query.Where("is_active = ?", true)
	}
	if documentType != "" {
		query = query.Where("document_type = ?", documentType)
	}
	if category != "" {
		query = query.Where("LOWER(category) = LOWER(?)", category)
	}

	// Execute query
	if err := query.Find(&documents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch documents"})
		return
	}

	// Log bulk document access activity
	userID, _ := c.Get("user_id")
	var filterDesc string
	if documentType != "" {
		filterDesc += fmt.Sprintf(" type=%s", documentType)
	}
	if category != "" {
		filterDesc += fmt.Sprintf(" category=%s", category)
	}
	if activeOnly == "true" {
		filterDesc += " active-only"
	}
	logSystemActivity(c, userID.(uint), ActionView, fmt.Sprintf("Accessed %d documents%s", len(documents), filterDesc))

	c.JSON(http.StatusOK, documents)
}

// Get document by ID
func getDocumentByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var document PolicyFile
	if err := db.First(&document, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch document"})
		}
		return
	}

	// Log document view activity
	userID, _ := c.Get("user_id")
	logDocumentActivity(c, userID.(uint), ActionView, &document, fmt.Sprintf("Viewed %s document: %s", document.DocumentType, document.Name))

	c.JSON(http.StatusOK, document)
}

// Download original file for a document
func downloadDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("❌ Download failed: Invalid document ID: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	log.Printf("🔍 Attempting to download document ID: %d", id)

	var document PolicyFile
	if err := db.First(&document, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("❌ Download failed: Document %d not found", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		} else {
			log.Printf("❌ Download failed: Database error for document %d: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch document"})
		}
		return
	}

	log.Printf("📄 Found document: Name='%s', FilePath='%s'", document.Name, document.FilePath)

	// Check if document has an original file
	if document.FilePath == "" {
		log.Printf("❌ Download failed: Document %d has no file path", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "No original file available for this document"})
		return
	}

	// Check if file exists on disk
	if _, err := os.Stat(document.FilePath); os.IsNotExist(err) {
		log.Printf("❌ Download failed: File not found on disk: %s", document.FilePath)
		c.JSON(http.StatusNotFound, gin.H{"error": "Original file not found on server"})
		return
	}

	// Log download activity
	userID, _ := c.Get("user_id")
	logDocumentActivity(c, userID.(uint), ActionView, &document, fmt.Sprintf("Downloaded original file: %s", document.FilePath))

	// Set appropriate headers and serve file
	storedFilename := filepath.Base(document.FilePath)
	
	// Extract original filename by removing timestamp (format: name_YYYYMMDD_HHMMSS.ext)
	ext := filepath.Ext(storedFilename)
	nameWithoutExt := strings.TrimSuffix(storedFilename, ext)
	
	// Initialize with the document name as fallback
	downloadFilename := "document"
	
	// Try to extract original filename from stored filename by removing timestamp pattern
	// The timestamp pattern is: _YYYYMMDD_HHMMSS (15 characters total)
	if len(nameWithoutExt) > 15 {
		// Look for the last occurrence of the timestamp pattern: _YYYYMMDD_HHMMSS
		// Pattern: underscore + 8 digits + underscore + 6 digits
		timestampPattern := regexp.MustCompile(`_\d{8}_\d{6}$`)
		if timestampPattern.MatchString(nameWithoutExt) {
			// Remove the timestamp part to get original filename
			downloadFilename = timestampPattern.ReplaceAllString(nameWithoutExt, "")
		} else {
			// If no timestamp pattern found, use the stored filename without extension
			downloadFilename = nameWithoutExt
		}
	} else {
		// If filename is too short to contain timestamp, use as-is
		downloadFilename = nameWithoutExt
	}
	
	// If we couldn't extract a good filename, fall back to document name
	if downloadFilename == "" || downloadFilename == "document" {
		if document.Name != "" {
			// Clean the document name to remove any file extension it might have
			documentNameWithoutExt := strings.TrimSuffix(document.Name, filepath.Ext(document.Name))
			
			// Also clean timestamp from document name if it exists (for legacy documents)
			timestampPattern := regexp.MustCompile(`_\d{8}_\d{6}$`)
			if timestampPattern.MatchString(documentNameWithoutExt) {
				documentNameWithoutExt = timestampPattern.ReplaceAllString(documentNameWithoutExt, "")
			}
			
			downloadFilename = documentNameWithoutExt
		} else {
			downloadFilename = "document"
		}
	}
	
	// Ensure we have the proper extension
	if !strings.HasSuffix(downloadFilename, ext) {
		downloadFilename = downloadFilename + ext
	}

	log.Printf("✅ Generated download filename: '%s'", downloadFilename)
	
	headerValue := fmt.Sprintf("attachment; filename=\"%s\"", downloadFilename)
	log.Printf("📋 Setting Content-Disposition header: %s", headerValue)
	
	c.Header("Content-Disposition", headerValue)
	c.Header("Content-Type", "application/octet-stream")
	
	log.Printf("🚀 Serving file: %s", document.FilePath)
	c.File(document.FilePath)
}

// Create new document
func createDocument(c *gin.Context) {
	var req CreateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate document type
	if req.DocumentType != "policy" && req.DocumentType != "onboarding" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document type must be 'policy' or 'onboarding'"})
		return
	}

	// Create new document
	newDoc := PolicyFile{
		Name:         req.Name,
		Content:      req.Content,
		Description:  req.Description,
		Category:     req.Category,
		DocumentType: req.DocumentType,
		TagsArray:    req.Tags,
		CreatedBy:    req.CreatedBy,
		FilePath:     req.FilePath,
		IsActive:     true,
	}

	// Save to database
	if err := db.Create(&newDoc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create document"})
		return
	}

	// Log document creation
	userID, _ := c.Get("user_id")
	logDocumentActivity(c, userID.(uint), ActionCreate, &newDoc, fmt.Sprintf("Created %s document: %s", newDoc.DocumentType, newDoc.Name))

	// Update search engine with fresh database data
	searchEngine := NewSearchEngine()
	_ = searchEngine // Update global reference if needed

	c.JSON(http.StatusCreated, newDoc)
}

// Update existing document
func updateDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var req UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find existing document
	var document PolicyFile
	if err := db.First(&document, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find document"})
		}
		return
	}

	// Update only provided fields
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.DocumentType != "" {
		if req.DocumentType != "policy" && req.DocumentType != "onboarding" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Document type must be 'policy' or 'onboarding'"})
			return
		}
		updates["document_type"] = req.DocumentType
	}
	if req.Tags != nil {
		// Update TagsArray field and let BeforeSave handle JSON conversion
		document.TagsArray = req.Tags
		tagsJSON, _ := json.Marshal(req.Tags)
		updates["tags"] = string(tagsJSON)
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	// Update in database
	if err := db.Model(&document).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update document"})
		return
	}

	// Fetch updated document
	if err := db.First(&document, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated document"})
		return
	}

	// Log document update
	userID, _ := c.Get("user_id")
	logDocumentActivity(c, userID.(uint), ActionUpdate, &document, fmt.Sprintf("Updated %s document: %s", document.DocumentType, document.Name))

	// Update search engine with fresh database data
	searchEngine := NewSearchEngine()
	_ = searchEngine // Update global reference if needed

	c.JSON(http.StatusOK, document)
}

// Delete document (soft delete by setting IsActive to false)
func deleteDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	// Find existing document
	var document PolicyFile
	if err := db.First(&document, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find document"})
		}
		return
	}

	// Soft delete by setting IsActive to false
	if err := db.Model(&document).Update("is_active", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document"})
		return
	}

	// Log document deletion
	userID, _ := c.Get("user_id")
	logDocumentActivity(c, userID.(uint), ActionDelete, &document, fmt.Sprintf("Deleted %s document: %s", document.DocumentType, document.Name))
	
	// Update search engine with fresh database data
	searchEngine := NewSearchEngine()
	_ = searchEngine // Update global reference if needed
	
	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}

// Advanced search for documents
func searchDocuments(c *gin.Context) {
	query := c.Query("q")
	documentType := c.Query("type")
	category := c.Query("category")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	// Use enhanced search engine with fresh database data
	searchEngine := NewSearchEngine()
	matches := searchEngine.Search(query, 20) // Allow more results for dashboard
	
	var filteredDocuments []PolicyFile
	for _, match := range matches {
		doc := match.Document
		
		// Apply type and category filters
		if documentType != "" && doc.DocumentType != documentType {
			continue
		}
		if category != "" && !strings.EqualFold(doc.Category, category) {
			continue
		}
		
		filteredDocuments = append(filteredDocuments, doc)
	}

	// Log document search activity
	userID, _ := c.Get("user_id")
	var filterDesc string
	if documentType != "" {
		filterDesc += fmt.Sprintf(" type=%s", documentType)
	}
	if category != "" {
		filterDesc += fmt.Sprintf(" category=%s", category)
	}
	logSystemActivity(c, userID.(uint), ActionView, fmt.Sprintf("Searched documents: '%s' returned %d results%s", query, len(filteredDocuments), filterDesc))

	response := gin.H{
		"documents": filteredDocuments,
		"total":     len(filteredDocuments),
		"query":     c.Query("q"),
		"matches":   matches[:minInt(len(matches), len(filteredDocuments))], // Include match details
	}

	c.JSON(http.StatusOK, response)
}

// File upload handlers

// Handle file upload and text extraction
func handleFileUpload(c *gin.Context) {
	// Parse multipart form
	err := c.Request.ParseMultipartForm(10 << 20) // 10MB max memory
	if err != nil {
		c.JSON(http.StatusBadRequest, FileUploadResponse{
			Success: false,
			Error:   "Failed to parse form data",
		})
		return
	}

	// Get uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, FileUploadResponse{
			Success: false,
			Error:   "No file provided or invalid file field name. Use 'file' as field name.",
		})
		return
	}

	// Validate file
	if err := validateFileUpload(fileHeader); err != nil {
		c.JSON(http.StatusBadRequest, FileUploadResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Save original file to disk
	filePath, err := saveUploadedFile(fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, FileUploadResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to save file: %v", err),
		})
		return
	}

	// Extract text from file
	extractedText, err := extractTextFromFile(fileHeader)
	if err != nil {
		// If text extraction fails, clean up the saved file
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, FileUploadResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to extract text: %v", err),
		})
		return
	}

	// Detect MIME type for response
	file, _ := fileHeader.Open()
	defer file.Close()
	content, _ := io.ReadAll(file)
	mtype := mimetype.Detect(content)

	// Log file upload activity
	userID, _ := c.Get("user_id")
	logSystemActivity(c, userID.(uint), ActionCreate, fmt.Sprintf("Uploaded and processed file: %s (%s, %d bytes) saved to %s", fileHeader.Filename, mtype.String(), fileHeader.Size, filePath))

	// Return successful response with extracted text and file path
	c.JSON(http.StatusOK, FileUploadResponse{
		Success:       true,
		Message:       "File uploaded and processed successfully",
		FileName:      fileHeader.Filename,
		FileType:      mtype.String(),
		FileSize:      fileHeader.Size,
		FilePath:      filePath,
		ExtractedText: extractedText,
	})
}

// Get supported file types
func handleGetSupportedFileTypes(c *gin.Context) {
	supportedTypes := getSupportedFileTypes()
	
	c.JSON(http.StatusOK, gin.H{
		"supported_types": supportedTypes,
		"max_file_size":   "10MB",
		"accepted_extensions": []string{".pdf", ".docx", ".txt", ".md"},
	})
}

// Text processing utilities
func tokenize(text string) []string {
	// Remove punctuation and split by whitespace
	reg := regexp.MustCompile(`[^\p{L}\p{N}]+`)
	text = reg.ReplaceAllString(text, " ")
	
	words := strings.Fields(strings.ToLower(text))
	
	// Remove common stop words
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true,
		"but": true, "in": true, "on": true, "at": true, "to": true,
		"for": true, "of": true, "with": true, "by": true, "is": true,
		"are": true, "was": true, "were": true, "be": true, "been": true,
		"have": true, "has": true, "had": true, "do": true, "does": true,
		"did": true, "will": true, "would": true, "could": true, "should": true,
	}
	
	var filtered []string
	for _, word := range words {
		if !stopWords[word] && len(word) > 1 {
			filtered = append(filtered, word)
		}
	}
	
	return filtered
}

func normalizeWord(word string) string {
	// Convert to lowercase and remove diacritics if needed
	word = strings.ToLower(word)
	
	// Basic stemming - remove common suffixes
	suffixes := []string{"ing", "ed", "er", "est", "ly", "ion", "tion", "sion", "ness", "ment"}
	for _, suffix := range suffixes {
		if strings.HasSuffix(word, suffix) && len(word) > len(suffix)+2 {
			word = word[:len(word)-len(suffix)]
			break
		}
	}
	
	return word
}

// Levenshtein distance for fuzzy matching
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

func min(a, b, c int) int {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Build full-text search index
func (se *SearchEngine) BuildIndex() {
	se.Index = make(map[string][]DocumentIndex)
	
	for _, doc := range se.Documents {
		if !doc.IsActive {
			continue
		}
		
		// Index different fields with different weights
		se.indexField(int(doc.ID), "name", doc.Name, 3.0)
		se.indexField(int(doc.ID), "description", doc.Description, 2.0)
		se.indexField(int(doc.ID), "content", doc.Content, 1.0)
		se.indexField(int(doc.ID), "category", doc.Category, 2.5)
		se.indexField(int(doc.ID), "tags", strings.Join(doc.TagsArray, " "), 2.0)
	}
}

func (se *SearchEngine) indexField(docID int, field, text string, weight float64) {
	words := tokenize(text)
	
	for pos, word := range words {
		if len(word) < 2 { // Skip very short words
			continue
		}
		
		word = normalizeWord(word)
		
		// Find existing index entry
		found := false
		for i := range se.Index[word] {
			if se.Index[word][i].DocumentID == docID && se.Index[word][i].Field == field {
				se.Index[word][i].Frequency++
				se.Index[word][i].Positions = append(se.Index[word][i].Positions, pos)
				found = true
				break
			}
		}
		
		if !found {
			se.Index[word] = append(se.Index[word], DocumentIndex{
				DocumentID: docID,
				Field:      field,
				Frequency:  1,
				Positions:  []int{pos},
			})
		}
	}
}

// Enhanced search with relevance scoring
func (se *SearchEngine) Search(query string, limit int) []DocumentMatch {
	if limit == 0 {
		limit = 10
	}
	
	queryWords := tokenize(query)
	if len(queryWords) == 0 {
		return []DocumentMatch{}
	}
	
	// Calculate document scores
	docScores := make(map[int]float64)
	docMatches := make(map[int][]Match)
	
	for _, queryWord := range queryWords {
		queryWord = normalizeWord(queryWord)
		
		// Try exact match first
		matches := se.findMatches(queryWord)
		
		// If no exact matches, try fuzzy matching
		if len(matches) == 0 {
			fuzzyMatches := se.findFuzzyMatches(queryWord, 2) // max 2 edits
			matches = append(matches, fuzzyMatches...)
		}
		
		// Calculate TF-IDF scores
		idf := se.calculateIDF(queryWord)
		
		for _, match := range matches {
			tf := float64(match.Frequency)
			fieldWeight := se.getFieldWeight(match.Field)
			score := tf * idf * fieldWeight
			
			docScores[match.DocumentID] += score
			
			docMatches[match.DocumentID] = append(docMatches[match.DocumentID], Match{
				Field: match.Field,
				Text:  queryWord,
				Score: score,
			})
		}
	}
	
	// Convert to sorted results
	var results []DocumentMatch
	for docID, score := range docScores {
		doc := se.getDocumentByID(uint(docID))
		if doc != nil {
			results = append(results, DocumentMatch{
				Document: *doc,
				Score:    score,
				Matches:  docMatches[docID],
			})
		}
	}
	
	// Sort by relevance score (descending)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	
	// Apply limit
	if len(results) > limit {
		results = results[:limit]
	}
	
	return results
}

func (se *SearchEngine) findMatches(word string) []DocumentIndex {
	return se.Index[word]
}

func (se *SearchEngine) findFuzzyMatches(word string, maxDistance int) []DocumentIndex {
	var matches []DocumentIndex
	
	for indexWord, docIndices := range se.Index {
		if levenshteinDistance(word, indexWord) <= maxDistance {
			matches = append(matches, docIndices...)
		}
	}
	
	return matches
}

func (se *SearchEngine) calculateIDF(word string) float64 {
	totalDocs := len(se.Documents)
	docsWithWord := len(se.Index[word])
	
	if docsWithWord == 0 {
		return 0
	}
	
	return math.Log(float64(totalDocs) / float64(docsWithWord))
}

func (se *SearchEngine) getFieldWeight(field string) float64 {
	weights := map[string]float64{
		"name":        3.0,
		"description": 2.0,
		"category":    2.5,
		"tags":        2.0,
		"content":     1.0,
	}
	
	if weight, exists := weights[field]; exists {
		return weight
	}
	return 1.0
}

func (se *SearchEngine) getDocumentByID(id uint) *PolicyFile {
	for _, doc := range se.Documents {
		if doc.ID == id {
			return &doc
		}
	}
	return nil
}

// Update documents in search engine
func (se *SearchEngine) UpdateDocuments(docs []PolicyFile) {
	se.Documents = docs
	se.BuildIndex()
}
