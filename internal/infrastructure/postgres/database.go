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
	ID       uuid.UUID `gorm:"type:uuid;primary_key"`
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
	dbLeak := Leak{
		ID:       uuid.New(),
		Email:    leak.Email,
		Password: leak.Password,
		Domain:   leak.Domain,
	}
	result := postgres.client.Create(&dbLeak)
	return result.Error
}

func (postgres *dbClient) GetByKeyword(key string, value string, page int, size int) ([]entities.Leak, error) {
	var leaks []Leak
	offset, limit := postgres.Paginate(page, size)
	err := postgres.client.Offset(offset).Limit(limit).Where(fmt.Sprintf("%s = ?", key), value).Find(&leaks).Error
	if err != nil {
		return nil, err
	}
	result := make([]entities.Leak, 0, len(leaks))
	for _, leak := range leaks {
		result = append(result, leak.toEntitiesLeak())
	}
	return result, nil
}

func (postgres *dbClient) FindByEmail(email string) ([]entities.Leak, error) {
	var leaks []Leak

	err := postgres.client.Find(&leaks, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	result := make([]entities.Leak, 0, len(leaks))
	for _, leak := range leaks {
		result = append(result, leak.toEntitiesLeak())
	}
	return result, nil
}

func (postgres *dbClient) Paginate(page int, size int) (int, int) {
	offset := (page - 1) * size
	return offset, size
}
