package postgres

import (
	"fmt"
	"log"

	"github.com/plumchecker/storage/internal/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dbClient struct {
	client *gorm.DB
}

type Config struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func New(cnfg *Config) (*dbClient, error) {
	dsn := fmt.Sprintf("host=%s user=%s  password=%s  dbname=%s  port=%s  sslmode=disable", cnfg.IP, cnfg.User, cnfg.Password, cnfg.Database, cnfg.Port)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	return &dbClient{
		client: conn,
	}, nil
}

func (postgres *dbClient) Create(leak entities.Leak) error {
	result := postgres.client.Create(&leak)
	return result.Error
}
