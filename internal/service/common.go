package service

import (
	"context"
	"fmt"
	"log"
	v1 "nunu-eth/api/v1"
	"nunu-eth/internal/model"
	"nunu-eth/internal/repository"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type AccountInfo struct {
	HexAccount     string `json:"hexAccount" gorm:"column:hexAccount"`
	HashHexAccount string `json:"hashHexAccount" gorm:"column:hashHexAccount"`
	BytesAccount   string `json:"bytesAccount" gorm:"column:bytesAccount"`
}

type CommonService interface {
	GetCommon(ctx context.Context, id int64) (*model.Common, error)
	Test(ctx context.Context, id int64) (*model.Common, error)

	ConnectTest(ctx context.Context, req *v1.ETHConnectRequestData) (status int, e error)

	AccountFormatInfo(ctx context.Context, req *v1.AccountAddress) (accountInfo AccountInfo, err error)
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

// AccountFormatInfo implements CommonService.
func (s *commonService) AccountFormatInfo(ctx context.Context, req *v1.AccountAddress) (accountInfo AccountInfo, err error) {
	panic("unimplemented")
}

func (s *commonService) GetCommon(ctx context.Context, id int64) (*model.Common, error) {
	return s.commonRepository.GetCommon(ctx, id)
}

func (s *commonService) Test(ctx context.Context, id int64) (*model.Common, error) {
	return s.commonRepository.Test(ctx, id)
}

func AccountFormatInfo(ctx context.Context, req *v1.AccountAddress) (account AccountInfo, err error) {
	address := common.HexToAddress(req.AccountAddress)
	if len(address) > 0 {
		account := &AccountInfo{
			HexAccount:     address.Hex(),
			HashHexAccount: address.String(),
			BytesAccount:   string(address.Bytes()),
		}
		_ = account
	}
	return
}

func (s *commonService) ConnectTest(ctx context.Context, req *v1.ETHConnectRequestData) (resultStatus int, e error) {
	fmt.Println("url: ", req.Url, "; port: ", req.Port)
	resultStatus = 0
	if req.Url == "" || req.Port == "" {
		resultStatus = -1
		return
	}
	address := "http://" + req.Url + ":" + req.Port
	log.Println("address init: ", address)
	if connect(address) {
		log.Println("init address success！！！")
	} else {
		resultStatus = -1
		log.Println("init address success Fail！！")
	}

	return
}

func connect(url string) bool {
	client, err := ethclient.Dial(url)
	if err != nil {
		fmt.Println("Could not connect to Infura with ethclient: fail")
		return false
	}
	_ = client
	return true
}
