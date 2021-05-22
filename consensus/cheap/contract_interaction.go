package cheap

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/net/context"
)

const caddress = "0xdf224098536510991780072E5A4d4EEb1CAD7730"
const ABI = `
[
	{
		"inputs": [],
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "num",
				"type": "uint256"
			}
		],
		"name": "retrieve",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "num",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "val",
				"type": "uint256"
			}
		],
		"name": "store",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]`
const method = "retrieve"

func loadAbi() (abi.ABI, error) {
	return abi.JSON(strings.NewReader(ABI))
}
func Try(a interface{}, e error) interface{} {
	if e != nil {
		panic(e)
	}
	return a
}


func makeData(Method abi.Method, args ...interface{}) hexutil.Bytes {
	d, err := Method.Inputs.Pack(args...)
	if err != nil {
		panic("pack error")
	}
	return (hexutil.Bytes)(append(Method.ID, d...))
}
func contract_call(block_hash common.Hash, api *ethapi.PublicBlockChainAPI) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	contract_abi, err := loadAbi()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	Method := contract_abi.Methods[method]

	bytes_data := makeData(Method, big.NewInt(123))
	fmt.Println(bytes_data)
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
	final, err := Method.Outputs.Unpack(res)
	if err != nil {
		fmt.Println("Unpack err", err)
		return "", err
	}
	fmt.Println(final)
	return "", nil
}
