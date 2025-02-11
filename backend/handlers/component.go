package handlers

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/diovane-rinaldin/golden.path.catalog/backend/models"
	"github.com/gin-gonic/gin"
)

type ComponentHandler struct {
	db *dynamodb.Client
}

func NewComponentHandler(db *dynamodb.Client) *ComponentHandler {
	return &ComponentHandler{db: db}
}

func (h *ComponentHandler) Create(c *gin.Context) {
	var comp models.Component
	if err := c.BindJSON(&comp); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Save to DynamoDB
	// ... DynamoDB put logic ...

	c.JSON(201, comp)
}

func (h *ComponentHandler) ListByTechnology(c *gin.Context) {
	techID := c.Param("id")

	// Query DynamoDB by technology ID
	// ... DynamoDB query logic ...

	c.JSON(200, components)
}

func (h *ComponentHandler) List(c *gin.Context) {
	// Query DynamoDB by component ID
	// ... DynamoDB query logic ...

	c.JSON(200, components)
}

func (h *ComponentHandler) Get(c *gin.Context) {
	name := c.Param("name")
	// Query DynamoDB by component ID
	// ... DynamoDB query logic ...

	c.JSON(200, components)
}

func (h *ComponentHandler) Update(c *gin.Context) {
	id := c.Param("id")

	// Query DynamoDB by component ID
	// ... DynamoDB query logic ...

	c.JSON(200, components)
}
