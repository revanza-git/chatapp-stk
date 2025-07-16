// Database initialization script
// Run this once to populate your Railway PostgreSQL database
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Copy the models from your main.go
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null;size:50"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null;size:100"`
	Password  string    `json:"-" gorm:"not null;size:255"`
	FirstName string    `json:"first_name" gorm:"not null;size:50"`
	LastName  string    `json:"last_name" gorm:"not null;size:50"`
	Role      string    `json:"role" gorm:"not null;size:20;default:'user'"`
	IsActive  bool      `json:"is_active" gorm:"default:true;index"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type PolicyFile struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"not null;size:255;index"`
	Content     string    `json:"content" gorm:"type:text;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Category    string    `json:"category" gorm:"not null;size:100;index"`
	DocumentType string   `json:"document_type" gorm:"not null;size:50;index"`
	Tags        string    `json:"-" gorm:"type:text"`
	TagsArray   []string  `json:"tags" gorm:"-"`
	FilePath    string    `json:"file_path,omitempty" gorm:"size:500"`
	CreatedBy   string    `json:"created_by" gorm:"size:100"`
	CreatedByUserID *uint `json:"created_by_user_id,omitempty" gorm:"index"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	IsActive    bool      `json:"is_active" gorm:"default:true;index"`
}

type AuditLog struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       uint      `json:"user_id" gorm:"not null;index"`
	User         User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Action       string    `json:"action" gorm:"not null;size:50;index"`
	ResourceType string    `json:"resource_type" gorm:"not null;size:50;index"`
	ResourceID   *uint     `json:"resource_id,omitempty" gorm:"index"`
	ResourceName string    `json:"resource_name,omitempty" gorm:"size:255"`
	Details      string    `json:"details,omitempty" gorm:"type:text"`
	IPAddress    string    `json:"ip_address,omitempty" gorm:"size:45"`
	UserAgent    string    `json:"user_agent,omitempty" gorm:"type:text"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (p *PolicyFile) BeforeSave(tx *gorm.DB) error {
	if len(p.TagsArray) > 0 {
		tagsJSON, err := json.Marshal(p.TagsArray)
		if err != nil {
			return err
		}
		p.Tags = string(tagsJSON)
	}
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func main() {
	// Get database connection from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if host == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("Missing required environment variables: DB_HOST, DB_USER, DB_PASSWORD, DB_NAME")
	}

	if port == "" {
		port = "5432"
	}
	if sslmode == "" {
		sslmode = "require"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&User{}, &PolicyFile{}, &AuditLog{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("✅ Database schema migrated successfully")

	// Check if users already exist
	var userCount int64
	db.Model(&User{}).Count(&userCount)
	
	var adminUser User
	if userCount == 0 {
		hashedPassword, err := hashPassword("SecureAdmin123!")
		if err != nil {
			log.Fatal("Failed to hash admin password:", err)
		}
		
		adminUser = User{
			Username:  "admin",
			Email:     "admin@company.com",
			Password:  hashedPassword,
			FirstName: "System",
			LastName:  "Administrator",
			Role:      "admin",
			IsActive:  true,
		}
		
		if err := db.Create(&adminUser).Error; err != nil {
			log.Fatal("Failed to create admin user:", err)
		}
		
		log.Println("✅ Created default admin user:")
		log.Println("   Username: admin")
		log.Println("   Password: SecureAdmin123!")
		log.Println("   ⚠️  Please change the default password after first login!")
	} else {
		if err := db.Where("role = ?", "admin").First(&adminUser).Error; err != nil {
			if err := db.First(&adminUser).Error; err != nil {
				log.Fatal("Failed to find any user:", err)
			}
		}
		log.Println("✅ Admin user already exists")
	}

	// Check if policy data already exists
	var policyCount int64
	db.Model(&PolicyFile{}).Count(&policyCount)
	if policyCount > 0 {
		log.Println("✅ Policy data already exists")
		return
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
			continue
		}
	}

	log.Printf("✅ Successfully initialized database with %d documents", len(samplePolicies))
} 