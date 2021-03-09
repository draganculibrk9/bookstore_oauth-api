package db

import (
	"github.com/draganculibrk9/bookstore_oauth-api/src/domain/access_token"
	"github.com/draganculibrk9/bookstore_oauth-api/src/utils/errors"
)

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
}

func NewRepository() DbRepository {
	return &dbRepository{}
}

type dbRepository struct {
}

func (repository *dbRepository) GetById(id string) (*access_token.AccessToken, *errors.RestErr) {
	return nil, errors.NewInternalServerError("database connection not implemented yet")
}
