package service

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	v1 "nunu-eth/api/v1"
	"nunu-eth/internal/model"
	"nunu-eth/internal/repository"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type AccountInfo struct {
	HexAccount     string `json:"hexAccount" gorm:"column:hexAccount"`
	HashHexAccount string `json:"hashHexAccount" gorm:"column:hashHexAccount"`
	BytesAccount   []byte `json:"bytesAccount" gorm:"column:bytesAccount"`
}

type AccountBalance struct {
	Wei        string `json:"wei" gorm:"column:wei"`
	GWei       string `json:"gwei" gorm:"column:gwei"`
	UnDealWei  string `json:"unDealWei" gorm:"column:unDealWei"`
	UnDealGWei string `json:"unDealGWei" gorm:"column:unDealGWei"`
}

type CommonService interface {
	GetCommon(ctx context.Context, id int64) (*model.Common, error)
	Test(ctx context.Context, id int64) (*model.Common, error)

	ConnectTest(ctx context.Context, req *v1.ETHConnectRequestData) (status int, e error)

	AccountFormatInfo(ctx context.Context, req *v1.AccountAddress) (accountInfo *AccountInfo, err error)

	AccountBalance(ctx context.Context, req *v1.AccountBalanceRequest) (accountBalance *AccountBalance, err error)

	BlockQuery(ctx context.Context, req *v1.BlockQueryRequest) (accountBalance *AccountBalance, err error)
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
func (s *commonService) AccountFormatInfo(ctx context.Context, req *v1.AccountAddress) (accountInfo *AccountInfo, err error) {
	address := common.HexToAddress(req.AccountAddress)
	accountInfo = &AccountInfo{
		HexAccount:     address.Hex(),
		HashHexAccount: address.String(),
		BytesAccount:   address.Bytes(),
	}
	return
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
	log.Println("address init: ", address)
	if connect(address) {
		log.Println("init address success！！！")
	} else {
		resultStatus = -1
		log.Println("init address success Fail！！")
	}

	return
}

func (s *commonService) AccountBalance(ctx context.Context, req *v1.AccountBalanceRequest) (accountBalance *AccountBalance, err error) {
	fmt.Println("url: ", req.Url, "; address: ", req.Address, "; block： ", req.Block)
	client, err := ethclient.Dial(req.Url)
	if err != nil {
		return
	}
	accountBalance = &AccountBalance{}
	account := common.HexToAddress(req.Address)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return
	}
	accountBalance.Wei = balance.String()
	fmt.Println("可用最新余额(wei)", balance) // 25893180161173005034

	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	fmt.Println("可用最新余额(gwei)", ethValue) // 25893180161173005034
	accountBalance.GWei = ethValue.String()

	pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
	fmt.Println("最新可用余额（wei）", pendingBalance) // 25729324269165216042

	if req.Block != "" {
		blockBigInt, err2 := strconv.Atoi(req.Block)
		if err2 != nil {
			_ = blockBigInt
			err = err2
			return
		}
		blockNumber := big.NewInt(int64(blockBigInt))
		balanceAt, err3 := client.BalanceAt(context.Background(), account, blockNumber)
		if err3 != nil {
			_ = balanceAt
			err = err3
			return
		}
		fbalance2 := new(big.Float)
		fbalance2.SetString(balanceAt.String())
		ethValue2 := new(big.Float).Quo(fbalance2, big.NewFloat(math.Pow10(18)))
		fmt.Println(req.Block, "可用区块余额(wei)： ", balanceAt) // 25.729324269165216041
		accountBalance.UnDealWei = balanceAt.String()
		fmt.Println(req.Block, "可用区块余额(gwei)： ", ethValue2) // 25729324269165216042
		accountBalance.UnDealGWei = ethValue2.String()

	}
	return
}

func (s *commonService) BlockQuery(ctx context.Context, req *v1.BlockQueryRequest) (accountBalance *AccountBalance, err error) {
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
