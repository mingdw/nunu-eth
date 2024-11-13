package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("eth etst")
	//testInitClient()
	//testAccount()
	testAccountAccountBanlance()
}

func testInitClient() {
	fmt.Println("eth etst")
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("Could not connect to Infura with ethclient: %s", err)
	}
	_ = client
	fmt.Println("we have a connection: ", "success")
}

func testAccount() {
	address := common.HexToAddress("0x71c7656ec7ab88b098defb751b7401b5f6d8976f")
	fmt.Println(address.Hex())    // 0x71C7656EC7ab88b098defB751B7401B5f6d8976F
	fmt.Println(address.String()) // 0x00000000000000000000000071c7656ec7ab88b098defb751b7401b5f6d8976f
	fmt.Println(address.Bytes())  // [113 199 101 110 199 171 136 176 152 222 251 117 27 116 1 181 246 216 151 111]
}

func testAccountAccountBanlance() {
	fmt.Println("************测试获取账户余额**************")

	client, err := ethclient.Dial("http://localhost:1234")
	if err != nil {
		fmt.Println("连接失败")
		log.Fatal(err)
	}
	fmt.Println("connect ---> ", "http://localhost:8545", " 连接成功！！")
	account := common.HexToAddress("0x2461e4B64b78D4f64A127c3B7A5E0A57E5D098Ab")
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	//最小单位下的账户余额
	fmt.Println("default get account balance: ", balance) // 25893180161173005034

	//最高单位转换后的余额
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	fmt.Println("after format balance: ", ethValue) // 25.729324269165216041

	//可使用的余额
	pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
	fmt.Println("peiding balance : ", pendingBalance) // 25729324269165216042

}
