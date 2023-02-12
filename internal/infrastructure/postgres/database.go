package postgres

import (
	"fmt"
	"log"

	"github.com/google/uuid"
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

type Leak struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Email    string    `gorm:"type:varchar(255); not null"`
	Password string    `gorm:"type:varchar(255); not null"`
	Domain   string    `gorm:"type:varchar(255); not null"`
}

func (l Leak) toEntitiesLeak() entities.Leak {
	return entities.Leak{
		Email:    l.Email,
		Password: l.Password,
		Domain:   l.Domain,
	}
}

func New(cnfg *Config) (*dbClient, error) {
	dsn := fmt.Sprintf("host=%s user=%s  password=%s  dbname=%s  port=%s  sslmode=disable", cnfg.IP, cnfg.User, cnfg.Password, cnfg.Database, cnfg.Port)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	err = conn.AutoMigrate(&Leak{})
	if err != nil {
	}
	return &dbClient{
		client: conn,
	}, nil
}

func (postgres *dbClient) Create(leak entities.Leak) error {
	result := postgres.client.Create(&leak)
	return result.Error
}

func (postgres *dbClient) GetByKeyword(key string, value string) ([]entities.Leak, error) {
	var leaks []Leak
	err := postgres.client.Find(&leaks, fmt.Sprintf("%s = ?", key), value).Error
	if err != nil {
		return nil, err
	}
	result := make([]entities.Leak, 0, len(leaks))
	for _, leak := range leaks {
		result = append(result, leak.toEntitiesLeak())
	}
	return result, nil
}
