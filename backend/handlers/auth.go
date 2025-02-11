package handlers

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/dgrijalva/jwt-go"
	"github.com/diovane-rinaldin/golden.path.catalog/backend/models"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	kmsClient *kms.Client
}

func NewAuthHandler(kmsClient *kms.Client) *AuthHandler {
	return &AuthHandler{kmsClient: kmsClient}
}

func (h *AuthHandler) Authenticate(c *gin.Context) {
	var creds struct {
		Key    string `json:"key"`
		Secret string `json:"secret"`
	}

	if err := c.BindJSON(&creds); err != nil {
		c.JSON(400, gin.H{"error": "Invalid credentials format"})
		return
	}

	// Verify credentials using KMS
	// ... KMS verification logic ...

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 4).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(200, models.AuthResponse{Token: tokenString})
}
