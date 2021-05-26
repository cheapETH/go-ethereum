package cheap

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/cheap/contracts"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/net/context"
)

func loadAbi(s string) (abi.ABI, error) {
	return abi.JSON(strings.NewReader(s))
}

func makeData(Method abi.Method, args ...interface{}) hexutil.Bytes {

	if args == nil {
		return (hexutil.Bytes)(Method.ID)
	}

	d, err := Method.Inputs.Pack(args...)
	if err != nil {
		panic("pack error")
	}

	return (hexutil.Bytes)(append(Method.ID, d...))
}

type contract struct {
	api  *ethapi.PublicBlockChainAPI
	abi  abi.ABI
	addr common.Address
}

func NewContract(api *ethapi.PublicBlockChainAPI, abi_json string, addr common.Address) (*contract, error) {
	contract_abi, err := loadAbi(abi_json)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &contract{
		api:  api,
		abi:  contract_abi,
		addr: addr,
	}, nil
}

func (c *contract) Call(method_name string, args ...interface{}) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	method := c.abi.Methods[method_name]
	bData := makeData(method, args...)
	to := c.addr
	//TOOD: Use reasonble amoutn of gas, get rid of the warn
	gas := (hexutil.Uint64)(math.MaxUint32)
	callData := ethapi.CallArgs{
		Data: &bData,
		To:   &to,
		Gas:  &gas,
	}

	res, err := c.api.Call(
		ctx,
		callData,
		rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(c.api.BlockNumber())),
		nil,
	)

	if err != nil {
		return make([]byte, 0), err
	}

	return res, nil
}
func (c *contract) UnpackResult(data []byte, method_name string) ([]interface{}, error) {
	method := c.abi.Methods[method_name]
	res, err := method.Outputs.Unpack(data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type Checkpoint struct {
	SavedBlockNumber *big.Int  "json:\"savedBlockNumber\""
	Hash [32]uint8 "json:\"hash\""
}

func TrustedFromInterface(i []interface{}) []common.Address {
	return *abi.ConvertType(i[0], new([]common.Address)).(*[]common.Address)
}

func CheckpointFromInterface(i []interface{}) *Checkpoint {
	return abi.ConvertType(i[0], new(Checkpoint)).(*Checkpoint)
}

func (c *contract) GetTrusted() ([]common.Address, error) {
	data, err := c.Call("getTrusted")

	if err != nil {
		return nil, fmt.Errorf("call faliled with error: %s", err)
	}

	unpack, err := c.UnpackResult(data, "getTrusted")
	trusted := TrustedFromInterface(unpack)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack data: %s", err)
	}

	return trusted, nil
}

func (c *contract) GetBlockByNumber(number big.Int, verifier common.Address) (*Checkpoint, error){
	data, err := c.Call("getBlockByNumber", number, verifier)

	if err != nil {
		return nil, fmt.Errorf("call faliled with error: %s", err)
	}

	unpack, err := c.UnpackResult(data, "getBlockByNumber")
	point := CheckpointFromInterface(unpack)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack data: %s", err)
	}

	return point, nil

}

func InitCheckpointerContract(api *ethapi.PublicBlockChainAPI) (*contract) {
	contract, err := NewContract(api, contracts.Contracts["Checkpointer"].Abi, contracts.Contracts["Checkpointer"].Address)
	if err != nil {
		panic(err)
	}
	return contract
}
