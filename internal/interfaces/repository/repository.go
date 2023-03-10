package repository

import (
	"encoding/base64"
	"encoding/json"

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

type Pagination struct {
	Page int `json:"page"`
	Size int `json:"size"`
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

func (db *database) FindLeaksByKeyword(key string, value string, token string) ([]entities.Leak, string, string, error) {
	var paginationToken Pagination
	if token == "" {
		paginationToken = Pagination{
			Page: 1,
			Size: 10,
		}
	} else {
		decodedToken, _ := base64.StdEncoding.DecodeString(token)
		_ = json.Unmarshal(decodedToken, &paginationToken)
		if paginationToken.Page < 1 {
			paginationToken.Page = 1
		}
	}

	leaks, err := db.d.GetByKeyword(key, value, paginationToken.Page, paginationToken.Size)
	if err != nil {
		return nil, "", "", err
	}

	endToken, err := json.Marshal(Pagination{Page: paginationToken.Page + 1, Size: paginationToken.Size})
	endCursor := base64.StdEncoding.EncodeToString(endToken)

	var startCursor string
	if paginationToken.Page == 1 {
		startCursor = ""
	} else {
		startToken, _ := json.Marshal(Pagination{Page: paginationToken.Page - 1, Size: paginationToken.Size})
		startCursor = base64.StdEncoding.EncodeToString(startToken)
	}

	return leaks, startCursor, endCursor, err
}
