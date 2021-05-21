package cheap

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/rpc"
)

type Cheapconsensus struct {
	ethash *ethash.Ethash
	api *ethapi.PublicBlockChainAPI
	api_init bool
}

func New(config ethash.Config, notify []string, noverify bool, api *ethapi.PublicBlockChainAPI) *Cheapconsensus {
	ethash := ethash.New(config, notify, noverify)

	return &Cheapconsensus{
		ethash: ethash,
		api: api,
		api_init: false,
	}
}

func (c *Cheapconsensus) Author(header *types.Header) (common.Address, error) {
	return c.ethash.Author(header)
}
func (c *Cheapconsensus) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	
	if c.api_init {
		fmt.Printf("\n\nEthapi is %p\n", c.api)
		fmt.Printf("Chain ID: %d\n\n\n", c.api.ChainId().ToInt())
	} else {
		fmt.Printf("Api not ready yet...\n")
	}
	return c.ethash.VerifyHeader(chain, header, seal)
}
func (c *Cheapconsensus) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	return c.ethash.VerifyHeaders(chain, headers, seals)
}
func (c *Cheapconsensus) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	return c.ethash.VerifyUncles(chain, block)
}
func (c *Cheapconsensus) VerifySeal(chain consensus.ChainHeaderReader, header *types.Header) error {
	return c.ethash.VerifySeal(chain, header)
}
func (c *Cheapconsensus) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	return c.ethash.Prepare(chain, header)
}
func (c *Cheapconsensus) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header) {
	c.ethash.Finalize(chain, header, state, txs, uncles)
}
func (c *Cheapconsensus) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	// FIXME: this is bad, we need a better way to track api ready state
	if c.api != nil {
		fmt.Printf("\n\nEthapi is %p\n", c.api)
		fmt.Printf("Chain ID: %d\n\n\n", c.api.ChainId().ToInt())

		fmt.Printf("%s\n", string(state.Dump(false, false, false)))
		contract_call(header.ParentHash, c.api)
		c.api_init = true	
	}
	return c.ethash.FinalizeAndAssemble(chain, header, state, txs, uncles, receipts)
}
func (c *Cheapconsensus) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	return c.ethash.Seal(chain, block, results, stop)
}
func (c *Cheapconsensus) SealHash(header *types.Header) common.Hash {
	return c.ethash.SealHash(header)
}
func (c *Cheapconsensus) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	return c.ethash.CalcDifficulty(chain, time, parent)
}
func (c *Cheapconsensus) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	return c.ethash.APIs(chain)
}
func (c *Cheapconsensus) Close() error {
	return c.ethash.Close()
}
func (c *Cheapconsensus) SetThreads(threads int) {
	c.ethash.SetThreads(threads)
}

func (c *Cheapconsensus) Threads() int {
	return c.ethash.Threads()
}
