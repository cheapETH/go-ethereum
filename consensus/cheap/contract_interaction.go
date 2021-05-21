package cheap

import (
	"encoding/hex"
	"fmt"
	"math"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/net/context"
)

const caddress = "0x411D7Dd3A717fD95e808c21A347E174eD6aE78bc"
const ABI = `
[
	{
		"inputs": [],
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"inputs": [],
		"name": "a",
		"outputs": [
			{
				"internalType": "uint64",
				"name": "",
				"type": "uint64"
			}
		],
		"stateMutability": "pure",
		"type": "function"
	}
]`
const method = "retrieve"

func loadAbi() (abi.ABI, error) {
	return abi.JSON(strings.NewReader(ABI))
}

func contract_call(block_hash common.Hash, api *ethapi.PublicBlockChainAPI) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	contract_abi, err := loadAbi()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	tx := "0x0dbe671f"
	decodedSig, _ := hex.DecodeString(tx[2:10])

	method, err := contract_abi.MethodById(decodedSig)
	if err != nil {
		fmt.Println("methodid err", err)
		return "", err
	}
	
	// data, err := method.Inputs.Pack()
	// if err != nil {
	// 	fmt.Println("Pack err", err)
	// 	return "", err
	// }

	bytes_data := (hexutil.Bytes)(decodedSig)
	to := common.HexToAddress(caddress)
	gas := (hexutil.Uint64)(uint64(math.MaxUint16 / 2))
	callArgs := ethapi.CallArgs{
		Data: &bytes_data,
		To:   &to,
		Gas:  &gas,
	}
	res, err := api.Call(ctx, callArgs, rpc.BlockNumberOrHashWithHash(block_hash, false), nil)
	if err != nil {
		fmt.Println("Call err", err.Error())
		return "", err
	}
	final, err := method.Outputs.Unpack(res)
	if err != nil {
		fmt.Println("Unpack err", err)
		return "", err
	}
	fmt.Println(final)
	return "", nil
}
