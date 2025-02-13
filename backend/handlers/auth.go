package handlers

import (
	"net/http"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Credenciais em formato inválido"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível gerar o token"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{Token: tokenString})
}
