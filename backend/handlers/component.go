package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/diovane-rinaldin/golden.path.catalog/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de dados inválido"})
		return
	}

	if comp.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O nome do componente é obrigatório"})
		return
	}
	if comp.TechnologyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A tecnologia do componente é obrigatória"})
		return
	}
	if comp.ImageURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A URL da imagem do componente é obrigatória"})
		return
	}

	exist, err := h.getByName(comp.Name, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		log.Panic(err.Error())
		return
	}

	if exist == nil {

		comp.ID = uuid.New().String()

		item := map[string]types.AttributeValue{
			"id":            &types.AttributeValueMemberS{Value: comp.ID},
			"technology_id": &types.AttributeValueMemberS{Value: comp.TechnologyID},
			"name":          &types.AttributeValueMemberS{Value: comp.Name},
			"description":   &types.AttributeValueMemberS{Value: comp.Description},
			"image_url":     &types.AttributeValueMemberS{Value: comp.ImageURL},
		}
		_, err := h.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(models.ComponentsTable),
			Item:      item,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar o componente"})
			log.Panic(err.Error())
			return
		}

		c.JSON(http.StatusCreated, comp)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O componente já existe"})
		return
	}
}

func (h *ComponentHandler) ListByTechnology(c *gin.Context) {
	techID := c.Param("id")
	if techID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A tecnologia deve ser informada"})
		return
	}
	result, err := h.db.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:        aws.String(models.ComponentsTable),
		FilterExpression: aws.String("technology_id = :technology_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":technology_id": &types.AttributeValueMemberS{Value: techID},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao buscar os componentes da tecnologias informada"})
		log.Panic(err.Error())
		return
	}

	if result.Items == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Componentes não encontrados"})
		return
	}

	var components []models.Component
	for _, item := range result.Items {
		comp := models.Component{
			ID:           item["id"].(*types.AttributeValueMemberS).Value,
			TechnologyID: item["technology_id"].(*types.AttributeValueMemberS).Value,
			Name:         item["name"].(*types.AttributeValueMemberS).Value,
			Description:  item["description"].(*types.AttributeValueMemberS).Value,
			ImageURL:     item["image_url"].(*types.AttributeValueMemberS).Value,
		}
		components = append(components, comp)
	}

	if len(components) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Nenhum componente encontrado para esta tecnologia"})
		return
	}

	c.JSON(http.StatusOK, components)
}

func (h *ComponentHandler) List(c *gin.Context) {
	result, err := h.db.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(models.ComponentsTable),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao buscar os componentes"})
		log.Panic(err.Error())
		return
	}

	var components []models.Component
	for _, item := range result.Items {
		comp := models.Component{
			ID:           item["id"].(*types.AttributeValueMemberS).Value,
			TechnologyID: item["technology_id"].(*types.AttributeValueMemberS).Value,
			Name:         item["name"].(*types.AttributeValueMemberS).Value,
			Description:  item["description"].(*types.AttributeValueMemberS).Value,
			ImageURL:     item["image_url"].(*types.AttributeValueMemberS).Value,
		}
		components = append(components, comp)
	}

	c.JSON(http.StatusOK, components)
}

func (h *ComponentHandler) GetByName(c *gin.Context) {
	comp, _ := h.getByName(c.Param("name"), c)
	if comp == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Componente não encontrado"})
		return
	}

	c.JSON(http.StatusOK, comp)
}

func (h *ComponentHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O parâmetro Id deve ser informado"})
		return
	}

	result, err := h.db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(models.ComponentsTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar o componente"})
		log.Panic(err.Error())
		return
	}

	// Verifica se o item foi encontrado
	if result.Item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Componente não encontrado"})
		return
	}

	component := models.Component{
		ID:           result.Item["id"].(*types.AttributeValueMemberS).Value,
		TechnologyID: result.Item["technology_id"].(*types.AttributeValueMemberS).Value,
		Name:         result.Item["name"].(*types.AttributeValueMemberS).Value,
		Description:  result.Item["description"].(*types.AttributeValueMemberS).Value,
		ImageURL:     result.Item["image_url"].(*types.AttributeValueMemberS).Value,
	}

	c.JSON(http.StatusOK, component)
}

func (h *ComponentHandler) getByName(name string, c *gin.Context) (*models.Component, error) {
	if name == "" {
		return nil, fmt.Errorf("o parâmetro nome deve ser informado")
	}

	result, err := h.db.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(models.ComponentsTable),
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
		component := models.Component{
			ID:           item["id"].(*types.AttributeValueMemberS).Value,
			TechnologyID: item["technology_id"].(*types.AttributeValueMemberS).Value,
			Name:         item["name"].(*types.AttributeValueMemberS).Value,
			Description:  item["description"].(*types.AttributeValueMemberS).Value,
			ImageURL:     item["image_url"].(*types.AttributeValueMemberS).Value,
		}
		return &component, nil
	}
	return nil, nil
}

func (h *ComponentHandler) Update(c *gin.Context) {
	var comp models.Component
	if err := c.BindJSON(&comp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de dados inválido"})
		return
	}

	if comp.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "É necessário informar o Id do componente a ser atualizado"})
		return
	}

	exist, err := h.getByName(comp.Name, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		log.Panic(err.Error())
		return
	}

	if exist == nil || exist.ID == comp.ID {
		updateInput := &dynamodb.UpdateItemInput{
			TableName: aws.String(models.ComponentsTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: comp.ID},
			},
			UpdateExpression: aws.String("SET #name = :name, description = :description, image_url = :image_url"),
			ExpressionAttributeNames: map[string]string{
				"#name": "name", // Alias para o atributo reservado
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":name":        &types.AttributeValueMemberS{Value: comp.Name},
				":description": &types.AttributeValueMemberS{Value: comp.Description},
				":image_url":   &types.AttributeValueMemberS{Value: comp.ImageURL},
			},
			ReturnValues: types.ReturnValueUpdatedNew,
		}

		_, err := h.db.UpdateItem(context.TODO(), updateInput)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao atualizar o componente"})
			log.Panic(err.Error())
			return
		}

		c.JSON(http.StatusOK, comp)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Um componente com esse nome já existe"})
		return
	}
}
