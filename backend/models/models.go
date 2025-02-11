package models

type Technology struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

type Component struct {
	ID           string `json:"id"`
	TechnologyID string `json:"technology_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ImageURL     string `json:"image_url"`
}

type Service struct {
	ID              string `json:"id"`
	ComponentID     string `json:"component_id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	ImageURL        string `json:"image_url"`
	CloudProvider   string `json:"cloud_provider"`
	ServiceCloudURL string `json:"service_cloud_url"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
