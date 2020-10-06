package clients

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	tmLog "github.com/tendermint/tendermint/libs/log"
	// "github.com/hashicorp/go-multierror"
	// "github.com/rs/zerolog"
)

// Client is a structure to sign and broadcast tx to Ethereum chain used by signer mostly
type Client struct {
	logger tmLog.Logger
	// cfg             config.ChainConfiguration
	// chainID types.ChainID
	// pk              common.PubKey
	client  *ethclient.Client
	rpcHost string
	// kw              *KeySignWrapper
	// ethScanner      *BlockScanner
	// thorchainBridge *thorclient.ThorchainBridge
	// blockScanner    *blockscanner.BlockScanner
	// keySignPartyMgr *thorclient.KeySignPartyMgr
}

// NewClient init an Ethereum client
func NewClient(log tmLog.Logger) (*Client, error) {
	ethClient, err := ethclient.Dial("http://127.0.0.1:7545")
	if ethClient == nil {
		fmt.Println("Failed to new a client.")
	}
	if err != nil {
		return nil, err
	}

	return &Client{
		logger: log,
		client: ethClient,
	}, nil
}

// Start the scan
func (c *Client) Start() {
	for {
		select {
		default:
			block, _ := c.GetBlock()
			fmt.Println(block.Hash)
			time.Sleep(5 * time.Second)
		}
	}
}

// GetBlock to get block from chain
func (c *Client) GetBlock() (*types.Block, error) {
	block, err := c.client.BlockByNumber(context.Background(), big.NewInt(11))
	if err != nil {
		c.logger.Error("he")
	}

	fmt.Println(block)
	return block, nil
}

// GetBlockHeight return current block height
func (c *Client) GetBlockHeight() (*big.Int, error) {
	number, err := c.client.BlockNumber(context.Background())
	if err != nil {
		c.logger.Error("he")
	}

	fmt.Println(block)
	return number, nil

}
