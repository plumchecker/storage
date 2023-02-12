package storage

import (
	"fmt"

	"github.com/plumchecker/storage/internal/entities"
)

type repository interface {
	InsertLeak(leak entities.Leak) (bool, error)
	FindLeaksByKeyword(key string, value string) ([]entities.Leak, error)
}

type application struct {
	repo repository
}

type Controller interface {
	AddLeaks(leaks []entities.Leak) (int, error)
	FindLeaksByKeyword(key string, value string) ([]entities.Leak, error)
}

func New(repo repository) (*application, error) {
	return &application{
		repo: repo,
	}, nil
}

func (app *application) AddLeaks(leaks []entities.Leak) (int, error) {
	counter := 0
	for _, leak := range leaks {
		isAdded, err := app.repo.InsertLeak(leak)
		if err != nil {
			continue
		}
		if isAdded {
			counter++
			fmt.Sprintf("Added %d leaks", counter)
		}
	}
	return counter, nil
}

func (app *application) FindLeaksByKeyword(key string, value string) ([]entities.Leak, error) {
	leaks, err := app.repo.FindLeaksByKeyword(key, value)
	if err != nil {
		return nil, err
	}
	return leaks, nil
}
