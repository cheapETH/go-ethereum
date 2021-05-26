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
	config   ethash.Config
	ethash   *ethash.Ethash
	api      *ethapi.PublicBlockChainAPI
	api_init bool
	contract *contract
}

const (
	EnforcingCheckpointing = false
	MIN_VERIFIERS          = 10
	MINIMUM_BLOCK_LOOKBACK = 100
)

func New(config ethash.Config, notify []string, noverify bool, api *ethapi.PublicBlockChainAPI) *Cheapconsensus {
	ethash := ethash.New(config, notify, noverify)

	return &Cheapconsensus{
		config:   config,
		ethash:   ethash,
		api:      api,
		api_init: false,
	}
}

func (c *Cheapconsensus) Author(header *types.Header) (common.Address, error) {
	return c.ethash.Author(header)
}
func (c *Cheapconsensus) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {

	if c.api_init && c.contract != nil {
		trusted, err := c.contract.GetTrusted()
		// TODO: find an elegant way of passing errors if enforcing checkpointing or just log warns if not
		if err != nil {
			c.config.Log.Warn("error getting trusted", "err", err)
		}

		if len(trusted) < MIN_VERIFIERS {
			c.config.Log.Warn("too little trusted addresses to work properly")
		}

		live_height := header.Number
		last_possible := big.NewInt(0)
		last_possible = last_possible.Sub(live_height, big.NewInt(MINIMUM_BLOCK_LOOKBACK))
		last_possible = last_possible.Mod(last_possible, big.NewInt(10))
		// Nice uint64
		last_possible_block := chain.GetHeaderByNumber(last_possible.Uint64())

		var matched []common.Address
		//TODO: distribute rewards
		for _, v := range trusted {
			cp, err := c.contract.GetBlockByNumber(*last_possible, v)
			if err != nil {
				c.config.Log.Warn("Error getting checkpoint", "err", err)
			}
			if cp.Hash == last_possible_block.Hash() && cp.SavedBlockNumber == last_possible_block.Number {
				matched = append(matched, v)
			}
		}

		treshold := len(trusted)/2 + 1
		if len(matched) < treshold {
			c.config.Log.Warn("Not enouhg votes, should be treatead as invalid chain")
		}

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
		c.contract = InitCheckpointerContract(c.api)
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
