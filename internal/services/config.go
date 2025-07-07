package services

import (
	"github.com/koki-develop/gorocket/internal/models"
	"github.com/koki-develop/gorocket/internal/providers"
)

type ConfigService interface {
	ConfigExists() bool
	CreateDefaultConfig() error
	LoadConfig(templateData *models.TemplateData) (*models.Config, error)
	GetDefaultConfigData() []byte
}

type configService struct {
	configProvider providers.ConfigProvider
}

func NewConfigService(fsProvider providers.FileSystemProvider) ConfigService {
	return &configService{
		configProvider: providers.NewConfigProvider(fsProvider),
	}
}

func (c *configService) ConfigExists() bool {
	return c.configProvider.ConfigExists()
}

func (c *configService) CreateDefaultConfig() error {
	return c.configProvider.CreateDefaultConfig()
}

func (c *configService) LoadConfig(templateData *models.TemplateData) (*models.Config, error) {
	return c.configProvider.LoadConfig(templateData)
}

func (c *configService) GetDefaultConfigData() []byte {
	return c.configProvider.GetDefaultConfigData()
}
