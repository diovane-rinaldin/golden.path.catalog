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

type ServiceHandler struct {
	db *dynamodb.Client
}

func NewServiceHandler(db *dynamodb.Client) *ServiceHandler {
	return &ServiceHandler{db: db}
}

func (h *ServiceHandler) Create(c *gin.Context) {
	var svc models.Service
	if err := c.BindJSON(&svc); err != nil {
		c.JSON(400, gin.H{"error": gin.H{"error": "Formato de dados inválido"}})
		return
	}

	if svc.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O nome do serviço é obrigatório"})
		return
	}
	if svc.ComponentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O componente do serviço é obrigatório"})
		return
	}
	if svc.ImageURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A URL da imagem do serviço é obrigatória"})
		return
	}
	if svc.CloudProvider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O cloud provider do serviço é obrigatório"})
		return
	}

	exist, err := h.getByName(svc.Name, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		log.Panic(err.Error())
		return
	}

	if exist == nil {

		svc.ID = uuid.New().String()

		item := map[string]types.AttributeValue{
			"id":                &types.AttributeValueMemberS{Value: svc.ID},
			"component_id":      &types.AttributeValueMemberS{Value: svc.ComponentID},
			"name":              &types.AttributeValueMemberS{Value: svc.Name},
			"description":       &types.AttributeValueMemberS{Value: svc.Description},
			"cloud_provider":    &types.AttributeValueMemberS{Value: svc.CloudProvider},
			"image_url":         &types.AttributeValueMemberS{Value: svc.ImageURL},
			"service_cloud_url": &types.AttributeValueMemberS{Value: svc.ServiceCloudURL},
		}
		_, err := h.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(models.ServicesTable),
			Item:      item,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar o serviço"})
			log.Panic(err.Error())
			return
		}

		c.JSON(http.StatusCreated, svc)
	}
}

func (h *ServiceHandler) ListByComponent(c *gin.Context) {
	compID := c.Param("id")
	if compID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O componente deve ser informado"})
		return
	}
	result, err := h.db.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:        aws.String(models.ServicesTable),
		FilterExpression: aws.String("component_id = :component"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":component": &types.AttributeValueMemberS{Value: compID},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao buscar os serviços do componente informado"})
		log.Panic(err.Error())
		return
	}
	if result.Items == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Serviços não encontrados"})
		return
	}

	var services []models.Service
	for _, item := range result.Items {
		service := models.Service{
			ID:              item["id"].(*types.AttributeValueMemberS).Value,
			ComponentID:     item["component_id"].(*types.AttributeValueMemberS).Value,
			Name:            item["name"].(*types.AttributeValueMemberS).Value,
			Description:     item["description"].(*types.AttributeValueMemberS).Value,
			ImageURL:        item["image_url"].(*types.AttributeValueMemberS).Value,
			ServiceCloudURL: item["service_cloud_url"].(*types.AttributeValueMemberS).Value,
		}
		services = append(services, service)
	}

	if len(services) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Nenhum serviço encontrado para este componente"})
		return
	}

	c.JSON(http.StatusOK, services)
}

func (h *ServiceHandler) List(c *gin.Context) {
	result, err := h.db.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(models.ServicesTable),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao buscar os serviços"})
		log.Panic(err.Error())
		return
	}

	var services []models.Service
	for _, item := range result.Items {
		service := models.Service{
			ID:              item["id"].(*types.AttributeValueMemberS).Value,
			ComponentID:     item["component_id"].(*types.AttributeValueMemberS).Value,
			Name:            item["name"].(*types.AttributeValueMemberS).Value,
			Description:     item["description"].(*types.AttributeValueMemberS).Value,
			ImageURL:        item["image_url"].(*types.AttributeValueMemberS).Value,
			CloudProvider:   item["cloud_provider"].(*types.AttributeValueMemberS).Value,
			ServiceCloudURL: item["service_cloud_url"].(*types.AttributeValueMemberS).Value,
		}
		services = append(services, service)
	}

	c.JSON(http.StatusOK, services)
}

func (h *ServiceHandler) Get(c *gin.Context) {
	svc, _ := h.getByName(c.Param("name"), c)
	if svc == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Serviço não encontrado"})
		return
	}

	c.JSON(http.StatusOK, svc)
}

func (h *ServiceHandler) getByName(name string, c *gin.Context) (*models.Service, error) {
	if name == "" {
		return nil, fmt.Errorf("o parâmetro nome deve ser informado")
	}

	result, err := h.db.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(models.ServicesTable),
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
		return nil, fmt.Errorf("Ocorreu um erro ao buscar o serviço")
	}

	if len(result.Items) > 0 {
		item := result.Items[0]
		service := models.Service{
			ID:          item["id"].(*types.AttributeValueMemberS).Value,
			ComponentID: item["component_id"].(*types.AttributeValueMemberS).Value,
			Name:        item["name"].(*types.AttributeValueMemberS).Value,
			Description: item["description"].(*types.AttributeValueMemberS).Value,
			ImageURL:    item["image_url"].(*types.AttributeValueMemberS).Value,
		}
		return &service, nil
	}
	return nil, nil
}

func (h *ServiceHandler) Update(c *gin.Context) {
	var service models.Service
	if err := c.BindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"error": "Formato de dados inválido"}})
		return
	}

	if service.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "É necessário informar o Id do serviço a ser atualizado"})
		return
	}

	exist, err := h.getByName(service.Name, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		log.Panic(err.Error())
		return
	}

	if exist == nil || exist.ID == service.ID {
		updateInput := &dynamodb.UpdateItemInput{
			TableName: aws.String(models.ServicesTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: service.ID},
			},
			UpdateExpression: aws.String("SET #name = :name, description = :description, cloud_provider = :cloud_provider, image_url = :image_url, service_cloud_url = :service_cloud_url"),
			ExpressionAttributeNames: map[string]string{
				"#name": "name", // Alias para o atributo reservado
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":name":              &types.AttributeValueMemberS{Value: service.Name},
				":description":       &types.AttributeValueMemberS{Value: service.Description},
				":cloud_provider":    &types.AttributeValueMemberS{Value: service.CloudProvider},
				":image_url":         &types.AttributeValueMemberS{Value: service.ImageURL},
				":service_cloud_url": &types.AttributeValueMemberS{Value: service.ServiceCloudURL},
			},
			ReturnValues: types.ReturnValueUpdatedNew,
		}

		_, err := h.db.UpdateItem(context.TODO(), updateInput)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao atualizar o componente"})
			log.Panic(err.Error())
			return
		}

		c.JSON(http.StatusOK, service)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Um componente com esse nome já existe"})
		return
	}
}
