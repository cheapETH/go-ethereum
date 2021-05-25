package contracts

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
)

type Contract struct {
	Address common.Address
	Abi string
	Code string
}
// TODO: make this read soidity files and populate abi and code
var Contracts map[string]Contract = map[string]Contract {
	"Dummy" : {
		Address: common.HexToAddress("0x1337000000000000000000000000000000000000"),
		Abi: `[
			{
				"inputs": [],
				"stateMutability": "nonpayable",
				"type": "constructor"
			},
			{
				"inputs": [],
				"name": "A",
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
		]`,
		Code: 
			"6080604052348015600f57600080fd5b506004361060285760003560e01c8063f446c1d014602d575b600080fd5b60336047565b604051603e9190605e565b60405180910390f35b6000611337905090565b6058816077565b82525050565b6000602082019050607160008301846051565b92915050565b600067ffffffffffffffff8216905091905056fea26469706673582212208eb084396988ccea5e49eca052fd74fe9f102653d5e0ddebdcfb2138428db45b64736f6c63430007060033",
		},
}


func Deploy(state *state.StateDB) {
	for n, c := range Contracts {
		code, err := hex.DecodeString(c.Code)
		if err != nil {
			panic(fmt.Sprintf("Bad code for contract %s\n", n))
		}

		state.SetCode(c.Address, code)
		// TODO: find a way to check if code was set correctly
	}
}