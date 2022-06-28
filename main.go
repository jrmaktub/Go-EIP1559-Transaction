package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
)

func main() {
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/2ea67de6c0fe4745b93a7bf99ba89c86")

	if err != nil {
		log.Fatal("Unable to reach client")
	}

	privateKey, err := crypto.HexToECDSA("157cb4bc16536b5b2f459ae9cde6100bdbd9b41814396f65207308748ba321be")

	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)

	if !ok {
		log.Fatal("Unable to cast public key to ECDSA")
	}
	//converting the public key into an address golang understands

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)

	if err != nil {
		log.Fatal("unable to get nonce")
	}

	fmt.Println("nonce", nonce)

	toAddress := common.HexToAddress("0x9289b45f7B5BC4E426E397C8695DF303A03670B3")

	value := big.NewInt(1000)

	gasFeeCap, gasTipCap, gas := big.NewInt(38694000460), big.NewInt(3869400046), uint64(22012)

	var data []byte
	//passing a pointer to a new dynamic fee transaction.
	tx := types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Gas:       gas,
		To:        &toAddress,
		Value:     value,
		Data:      data,
	})

	config, block := params.RinkebyChainConfig, params.RinkebyChainConfig.LondonBlock

	//signer
	signer := types.MakeSigner(config, block)

	signedTx, err := types.SignTx(tx, signer, privateKey)

	if err != nil {
		log.Fatal("Unable to sign TX")
	}

	//examining transaction Hash
	hash := signedTx.Hash().Bytes()

	//raw representation
	raw, err := rlp.EncodeToBytes(signedTx)

	if err != nil {
		log.Fatal("Unable to cast to raw txn")
	}

	fmt.Printf("hash: %x\n0x%x", hash, raw)

	client.SendTransaction(context.Background(), signedTx)

	if err != nil {
		log.Fatal("Unable to submit transaction")
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())

}
