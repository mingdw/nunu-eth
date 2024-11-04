package service

import (
	"context"
	"nunu-eth/internal/model"
	"nunu-eth/internal/repository"
)

type CommonService interface {
	GetCommon(ctx context.Context, id int64) (*model.Common, error)
	Test(ctx context.Context, id int64) (*model.Common, error)
}

func NewCommonService(
	service *Service,
	commonRepository repository.CommonRepository,
) CommonService {
	return &commonService{
		Service:          service,
		commonRepository: commonRepository,
	}
}

type commonService struct {
	*Service
	commonRepository repository.CommonRepository
}

func (s *commonService) GetCommon(ctx context.Context, id int64) (*model.Common, error) {
	return s.commonRepository.GetCommon(ctx, id)
}

func (s *commonService) Test(ctx context.Context, id int64) (*model.Common, error) {
	return s.commonRepository.Test(ctx, id)
}
