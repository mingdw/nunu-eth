package service

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math"
	"math/big"
	v1 "nunu-eth/api/v1"
	"nunu-eth/api/variable"
	"nunu-eth/internal/model"
	"nunu-eth/internal/repository"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
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

	BlockQuery(ctx context.Context, req *v1.BlockQueryRequest) (header *types.Header, err error)

	TransactionQuery(ctx context.Context, req *v1.TransactionsQueryRequest) (mapData map[string]interface{}, err error)

	CreateAccount(ctx context.Context) (mapData map[string]interface{}, err error)

	TxQuery(ctx context.Context, txHash string) (mapData map[string]interface{}, err error)
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
	log.Println("url: ", req.Url, "; port: ", req.Port)
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
	log.Println("url: ", req.Url, "; address: ", req.Address, "; block： ", req.Block)
	client, err := ethclient.Dial(getRealUrl(req.Url))
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
	log.Println("可用最新余额(wei)", balance) // 25893180161173005034

	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	fmt.Println("可用最新余额(gwei)", ethValue) // 25893180161173005034
	accountBalance.GWei = ethValue.String()

	pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
	log.Println("最新可用余额（wei）", pendingBalance) // 25729324269165216042

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
		log.Println(req.Block, "可用区块余额(wei)： ", balanceAt) // 25.729324269165216041
		accountBalance.UnDealWei = balanceAt.String()
		log.Println(req.Block, "可用区块余额(gwei)： ", ethValue2) // 25729324269165216042
		accountBalance.UnDealGWei = ethValue2.String()

	}
	return
}

func (s *commonService) BlockQuery(ctx context.Context, req *v1.BlockQueryRequest) (header *types.Header, err error) {
	log.Println("blockNum: ", req.BlockNum, "; address: ", req.Url)
	client, err := ethclient.Dial(getRealUrl(req.Url))
	if err != nil {
		return
	}
	bn, err := parseBlock(req.BlockNum)
	if err != nil {
		return
	}
	header, err = client.HeaderByNumber(context.Background(), bn) //查询区块信息，如果区块号为空则查询最新的区块信息
	if err != nil {
		return
	}
	// v, _ := json.Marshal(header)
	// jsonStr := string(v)
	// log.Println("eth header info: ", jsonStr)
	return
}

func (s *commonService) TransactionQuery(ctx context.Context, req *v1.TransactionsQueryRequest) (mapData map[string]interface{}, err error) {
	url := getRealUrl(req.Url)
	client, err := ethclient.Dial(url)
	if err != nil {
		return
	}
	mapData = make(map[string]interface{})
	blockHash := common.HexToHash(req.BlockHash)
	block, err := client.BlockByHash(context.Background(), blockHash) //查询区块信息，如果区块号为空则查询最新的区块信息
	if err != nil {
		return
	}
	mapData["total"] = block.Transactions().Len()
	mapData["data"] = block.Transactions()
	return
}

func (s *commonService) CreateAccount(ctx context.Context) (mapData map[string]interface{}, err error) {
	mapData = make(map[string]interface{})
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	accountPrivateKey := hexutil.Encode(privateKeyBytes)[2:]
	log.Println("privateKey: ", accountPrivateKey) // 0xfad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19
	mapData["accountPrivateKey"] = accountPrivateKey

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	accountPublic := hexutil.Encode(publicKeyBytes)[4:]
	log.Println("accountPunlicKey: ", accountPublic) // 0x049a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05
	mapData["accountPublicKey"] = accountPublic

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	log.Println("accountAddress", address) // 0x96216849c49358B10257cb55b28eA603c874b05E
	mapData["accountAddress"] = address

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	log.Println("手动生成账户地址：", hexutil.Encode(hash.Sum(nil)[12:])) // 0x96216849c49358b10257cb55b28ea603c874b05e
	mapData["accountAddress2"] = hexutil.Encode(hash.Sum(nil)[12:])
	return
}

func (s *commonService) TxQuery(ctx context.Context, txHash string) (mapData map[string]interface{}, err error) {

	url := getRealUrl("")
	client, err := ethclient.Dial(url)
	if err != nil {
		return
	}

	hash := common.HexToHash(txHash)
	tx, isPending, err := client.TransactionByHash(context.Background(), hash)
	if err != nil {
		return
	}
	mapData = make(map[string]interface{})
	mapData["tx"] = tx
	mapData["isPending"] = isPending
	return
}

func connect(url string) bool {
	client, err := ethclient.Dial(getRealUrl(url))
	if err != nil {
		fmt.Println("Could not connect to Infura with ethclient: fail")
		return false
	}
	_ = client
	return true
}

func parseBlock(blockNum string) (num *big.Int, err error) {
	if blockNum == "" {
		num = nil
		return
	}
	n, err := strconv.ParseInt(blockNum, 10, 64)
	if err != nil {
		_ = n
		return
	}
	num = big.NewInt(int64(n))
	return
}

func getRealUrl(url string) string {
	returnUrl := variable.EthClientAddress
	if url != "" && variable.ClientSwitch {
		return variable.EthClientAddress
	}
	return returnUrl
}
