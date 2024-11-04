package repository

import (
	"context"
	"nunu-eth/internal/model"
)

type CommonRepository interface {
	GetCommon(ctx context.Context, id int64) (*model.Common, error)
	Test(ctx context.Context, id int64) (*model.Common, error)
}

func NewCommonRepository(
	repository *Repository,
) CommonRepository {
	return &commonRepository{
		Repository: repository,
	}
}

type commonRepository struct {
	*Repository
}

func (r *commonRepository) GetCommon(ctx context.Context, id int64) (*model.Common, error) {
	var common model.Common

	return &common, nil
}

func (r *commonRepository) Test(ctx context.Context, id int64) (*model.Common, error) {
	var common model.Common

	return &common, nil
}
