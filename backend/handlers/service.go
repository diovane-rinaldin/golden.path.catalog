package handlers

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/diovane-rinaldin/golden.path.catalog/backend/models"
	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	db *dynamodb.Client
}

func NewServiceHandler(db *dynamodb.Client) *ServiceHandler {
	return &ServiceHandler{db: db}
}

func (h *ServiceHandler) Create(c *gin.Context) {
	var svc models.Service
	if err := c.BindJSON(&svc); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Save to DynamoDB
	// ... DynamoDB put logic ...

	c.JSON(201, svc)
}

func (h *ServiceHandler) ListByComponent(c *gin.Context) {
	compID := c.Param("id")

	// Query DynamoDB by component ID
	// ... DynamoDB query logic ...

	c.JSON(200, services)
}

func (h *ServiceHandler) List(c *gin.Context) {
	// Query DynamoDB by component ID
	// ... DynamoDB query logic ...

	c.JSON(200, services)
}

func (h *ServiceHandler) Get(c *gin.Context) {
	name := c.Param("name")
	// Query DynamoDB by component ID
	// ... DynamoDB query logic ...

	c.JSON(200, services)
}

func (h *ServiceHandler) Update(c *gin.Context) {
	id := c.Param("id")

	// Query DynamoDB by component ID
	// ... DynamoDB query logic ...

	c.JSON(200, services)
}
