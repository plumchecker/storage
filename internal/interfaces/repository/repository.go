package repository

import (
	"github.com/plumchecker/storage/internal/entities"
)

type driver interface {
	Create(leak entities.Leak) error
	GetByKeyword(key string, value string, page int, size int) ([]entities.Leak, error)
	FindByEmail(email string) ([]entities.Leak, error)
}

type database struct {
	d driver
}

func New(dbHandler driver) *database {
	return &database{
		d: dbHandler,
	}
}

func (db *database) InsertLeak(leak entities.Leak) (bool, error) {
	availableLeaks, err := db.d.FindByEmail(leak.Email)
	if err != nil {
		return false, err
	}

	for _, fLeak := range availableLeaks {
		if leak.Password == fLeak.Password && leak.Domain == fLeak.Domain {
			return false, nil
		}
	}

	err = db.d.Create(leak)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *database) FindLeaksByKeyword(key string, value string, token string) ([]entities.Leak, string, error) {
	page := 1
	size := 10
	leaks, err := db.d.GetByKeyword(key, value, page, size)
	if err != nil {
		return nil, "", err
	}
	return leaks, "", err
}
