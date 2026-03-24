package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// Config holds application configuration
type Config struct {
	Port      string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	JWTSecret string
}

// Document represents a PDF document in the system
type Document struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Owner    string `json:"owner"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
}

var db *sql.DB

func loadConfig() Config {
	return Config{
		Port:      getEnv("APP_PORT", "8080"),
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBUser:    getEnv("DB_USER", "app"),
		DBPass:    getEnv("DB_PASS", "changeme"),
		DBName:    getEnv("DB_NAME", "trivy_demo"),
		JWTSecret: getEnv("JWT_SECRET", "default-secret"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func initDB(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName,
	)
	return sql.Open("postgres", dsn)
}

func main() {
	cfg := loadConfig()

	log.SetFormatter(&log.JSONFormatter{})
	log.Info("Starting Trivy Demo Application")

	var err error
	db, err = initDB(cfg)
	if err != nil {
		log.WithError(err).Warn("Could not connect to database, running without persistence")
	} else {
		defer db.Close()
	}

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "version": "1.0.0"})
	})

	// Document endpoints
	r.GET("/api/documents", listDocuments)
	r.GET("/api/documents/:id", getDocument)
	r.POST("/api/documents", createDocument)

	// Auth endpoint
	r.POST("/api/login", login)

	log.Infof("Server listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func listDocuments(c *gin.Context) {
	// Demo: return static data when DB is not available
	docs := []Document{
		{ID: 1, Name: "report.pdf", Owner: "alice", Size: 102400, MimeType: "application/pdf"},
		{ID: 2, Name: "invoice.pdf", Owner: "bob", Size: 51200, MimeType: "application/pdf"},
	}
	c.JSON(http.StatusOK, docs)
}

func getDocument(c *gin.Context) {
	id := c.Param("id")
	doc := Document{ID: 1, Name: fmt.Sprintf("document-%s.pdf", id), Owner: "demo", Size: 1024, MimeType: "application/pdf"}
	c.JSON(http.StatusOK, doc)
}

func createDocument(c *gin.Context) {
	var doc Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	doc.ID = 100
	c.JSON(http.StatusCreated, doc)
}

func login(c *gin.Context) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Demo: accept any credentials
	if creds.Username == "" || creds.Password == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing credentials"})
		return
	}

	cfg := loadConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  creds.Username,
		"role": "user",
	})
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
