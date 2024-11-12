package service

import (
	"context"
	"fmt"
	"log"
	v1 "nunu-eth/api/v1"
	"nunu-eth/internal/model"
	"nunu-eth/internal/repository"

	"github.com/ethereum/go-ethereum/ethclient"
)

type CommonService interface {
	GetCommon(ctx context.Context, id int64) (*model.Common, error)
	Test(ctx context.Context, id int64) (*model.Common, error)

	ConnectTest(ctx context.Context, req *v1.ETHConnectRequestData) (status int, e error)
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

func (s *commonService) ConnectTest(ctx context.Context, req *v1.ETHConnectRequestData) (resultStatus int, e error) {
	fmt.Println("url: ", req.Url, "; port: ", req.Port)
	resultStatus = 0
	if req.Url == "" || req.Port == "" {
		resultStatus = -1
		return
	}
	address := "http://" + req.Url + ":" + req.Port
	log.Fatal("address: ", address)
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatal(err)
	}
	_ = client
	fmt.Println("we have a connection: ", client)

	return
}
