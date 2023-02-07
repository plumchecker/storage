package usecases

import "github.com/plumchecker/storage/internal/entities"

type repository interface {
	InsertLeak(leak entities.Leak) error
	// QueryLeaksByEmail(email string) ([]entities.Leak, error)
	// QueryLeaksByDomain(domain string) ([]entities.Leak, error)
	// QueryLeaksByPassword(password string) ([]entities.Leak, error)
}

type application struct {
	repo repository
}

type Controller interface {
	AddLeak(leak entities.Leak) error
}

func New(repo repository) *application {
	return &application{
		repo: repo,
	}
}

func (app *application) AddLeak(leak entities.Leak) error {
	return app.repo.InsertLeak(leak)
}
