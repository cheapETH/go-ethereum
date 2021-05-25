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

type MoreComplex struct {
	A *big.Int  "json:\"a\""
	B [32]uint8 "json:\"b\""
	C string    "json:\"c\""
}

func MoreComplexFromInterface(i []interface{}) *MoreComplex {
	return abi.ConvertType(i[0], new(MoreComplex)).(*MoreComplex)
}

func contract_call(api *ethapi.PublicBlockChainAPI) {
	contract, err := NewContract(api, contracts.Contracts["Dummy"].Abi, contracts.Contracts["Dummy"].Address)
	if err != nil {
		panic(err)
	}

	data, err := contract.Call("A")

	if err != nil {
		fmt.Println("Call faliled\n", err)
	}

	unpack, err := contract.UnpackResult(data, "A")
	asuint := abi.ConvertType(unpack[0], new(uint64)).(*uint64)
	fmt.Printf("---- %x\n", *asuint)
	//decode := MoreComplexFromInterface(unpack)
	//fmt.Printf(" -- %v %v %v\n", decode.A, decode.B, decode.C)
	if err != nil {
		fmt.Printf("Unpack err %s\n", err)
	}
}
