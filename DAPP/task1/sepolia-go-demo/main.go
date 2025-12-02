package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// ==================== 改下面这三行 ====================
	infuraURL := "https://sepolia.infura.io/v3/0648fa24cdba4be693f7b224dc63037f"        // ← 改这里
	privateKeyHex := "492508ec3d819967cbe6c00a163165e07280f755fbabec6a9a9d71ac1e3d2a2a" // ← 改这里（去掉0x）
	toAddress := "0xDA9dfA130Df4dE4673b89022EE50ff26f6EA73Cf"

	// ==================== 下面代码不要动 ====================
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatal("连接失败：", err)
	}
	defer client.Close()
	fmt.Println("成功连接 Sepolia 测试网！")

	// 查询最新区块
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal("查询区块失败：", err)
	}
	fmt.Printf("\n最新区块：%d | 哈希：%s | 时间：%s | 交易数：%d\n",
		block.Number(), block.Hash().Hex(),
		time.Unix(int64(block.Time()), 0).UTC().Format("2006-01-02 15:04:05"),
		len(block.Transactions()))

	// 发送交易
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal("私钥错误：", err)
	}
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	nonce, _ := client.PendingNonceAt(context.Background(), fromAddress)
	value := new(big.Int).Mul(big.NewInt(1e15), big.NewInt(1)) // 0.001 ETH
	gasPrice, _ := client.SuggestGasPrice(context.Background())
	tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), value, 21000, gasPrice, nil)
	chainID, _ := client.ChainID(context.Background())
	signedTx, _ := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	client.SendTransaction(context.Background(), signedTx)

	fmt.Printf("\n交易已发出！哈希：https://sepolia.etherscan.io/tx/%s\n", signedTx.Hash().Hex())
}
