package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de dados inválido"})
		return
	}

	if tech.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O nome da tecnologia é obrigatório"})
		return
	}
	if tech.ImageURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A URL da imagem da tecnologia é obrigatória"})
		return
	}

	exist, err := h.getByName(tech.Name, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		log.Panic(err.Error())
		return
	}

	if exist == nil {
		tech.ID = uuid.New().String()

		item := map[string]types.AttributeValue{
			"id":          &types.AttributeValueMemberS{Value: tech.ID},
			"name":        &types.AttributeValueMemberS{Value: tech.Name},
			"description": &types.AttributeValueMemberS{Value: tech.Description},
			"image_url":   &types.AttributeValueMemberS{Value: tech.ImageURL},
		}

		_, err := h.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(models.TechnologyTable),
			Item:      item,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao salvar a tecnologia"})
			log.Panic(err.Error())
			return
		}

		c.JSON(http.StatusCreated, tech)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A tecnologia já existe"})
		return
	}
}

func (h *TechnologyHandler) GetByName(c *gin.Context) {
	tech, _ := h.getByName(c.Param("name"), c)
	if tech == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tecnologia não encontrada"})
		return
	}
	c.JSON(http.StatusOK, tech)
}

func (h *TechnologyHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O parâmetro Id deve ser informado"})
		return
	}

	result, err := h.db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(models.TechnologyTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar tecnologia"})
		log.Panic(err.Error())
		return
	}

	// Verifica se o item foi encontrado
	if result.Item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tecnologia não encontrada"})
		return
	}

	technology := models.Technology{
		ID:          result.Item["id"].(*types.AttributeValueMemberS).Value,
		Name:        result.Item["name"].(*types.AttributeValueMemberS).Value,
		Description: result.Item["description"].(*types.AttributeValueMemberS).Value,
		ImageURL:    result.Item["image_url"].(*types.AttributeValueMemberS).Value,
	}

	c.JSON(http.StatusOK, technology)
}

func (h *TechnologyHandler) getByName(name string, c *gin.Context) (*models.Technology, error) {
	if name == "" {
		return nil, fmt.Errorf("O parâmetro nome deve ser informado")
	}

	result, err := h.db.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(models.TechnologyTable),
		IndexName:              aws.String("name-index"),
		KeyConditionExpression: aws.String("#name = :nameValue"),
		ExpressionAttributeNames: map[string]string{
			"#name": "name",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":nameValue": &types.AttributeValueMemberS{Value: name},
		},
		Limit: aws.Int32(1),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ""})
		log.Panic(err.Error())
		return nil, fmt.Errorf("Ocorreu um erro ao buscar a tecnologia")
	}

	if len(result.Items) > 0 {
		item := result.Items[0]
		technology := models.Technology{
			ID:          item["id"].(*types.AttributeValueMemberS).Value,
			Name:        item["name"].(*types.AttributeValueMemberS).Value,
			Description: item["description"].(*types.AttributeValueMemberS).Value,
			ImageURL:    item["image_url"].(*types.AttributeValueMemberS).Value,
		}
		return &technology, nil
	}
	return nil, nil
}

func (h *TechnologyHandler) List(c *gin.Context) {
	result, err := h.db.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(models.TechnologyTable),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao buscar as tecnologias"})
		log.Panic(err.Error())
		return
	}

	var technologies []models.Technology
	for _, item := range result.Items {
		tech := models.Technology{
			ID:          item["id"].(*types.AttributeValueMemberS).Value,
			Name:        item["name"].(*types.AttributeValueMemberS).Value,
			Description: item["description"].(*types.AttributeValueMemberS).Value,
			ImageURL:    item["image_url"].(*types.AttributeValueMemberS).Value,
		}
		technologies = append(technologies, tech)
	}

	c.JSON(http.StatusOK, technologies)
}

func (h *TechnologyHandler) Update(c *gin.Context) {
	var tech models.Technology
	if err := c.BindJSON(&tech); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de dados inválido"})
		return
	}

	if tech.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "É necessário informar o Id da tecnologia a ser atualizada"})
		return
	}
	exist, err := h.getByName(tech.Name, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		log.Panic(err.Error())
		return
	}

	if exist == nil || exist.ID == tech.ID {
		updateInput := &dynamodb.UpdateItemInput{
			TableName: aws.String(models.TechnologyTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: tech.ID},
			},
			UpdateExpression: aws.String("SET #name = :name, description = :description, image_url = :image_url"),
			ExpressionAttributeNames: map[string]string{
				"#name": "name", // Alias para o atributo reservado
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":name":        &types.AttributeValueMemberS{Value: tech.Name},
				":description": &types.AttributeValueMemberS{Value: tech.Description},
				":image_url":   &types.AttributeValueMemberS{Value: tech.ImageURL},
			},
			ReturnValues: types.ReturnValueUpdatedNew,
		}

		_, err := h.db.UpdateItem(context.TODO(), updateInput)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao atualizar a tecnologia"})
			log.Panic(err.Error())
			return
		}

		c.JSON(http.StatusOK, tech)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Uma tecnologia com esse nome já existe"})
		return
	}
}
