package handlers

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/diovane-rinaldin/golden.path.catalog/backend/models"
	"github.com/gin-gonic/gin"
)

type TechnologyHandler struct {
	db *dynamodb.Client
}

func NewTechnologyHandler(db *dynamodb.Client) *TechnologyHandler {
	return &TechnologyHandler{db: db}
}

func (h *TechnologyHandler) Create(c *gin.Context) {
	var tech models.Technology
	if err := c.BindJSON(&tech); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Save to DynamoDB
	// ... DynamoDB put logic ...

	c.JSON(201, tech)
}

func (h *TechnologyHandler) Get(c *gin.Context) {
	name := c.Param("name")

	// Get from DynamoDB
	// ... DynamoDB get logic ...

	c.JSON(200, tech)
}

func (h *TechnologyHandler) List(c *gin.Context) {
	// List from DynamoDB
	// ... DynamoDB scan logic ...

	c.JSON(200, techs)
}

func (h *TechnologyHandler) Update(c *gin.Context) {
	var tech models.Technology
	if err := c.BindJSON(&tech); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Update in DynamoDB
	// ... DynamoDB update logic ...

	c.JSON(200, tech)
}
