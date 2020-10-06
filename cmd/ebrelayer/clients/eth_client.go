package clients

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"io"
	"log"
	"math/big"
	"strconv"
	"time"

	sdkContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/peggy/cmd/ebrelayer/txs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	amino "github.com/tendermint/go-amino"
	tmLog "github.com/tendermint/tendermint/libs/log"
	"golang.org/x/crypto/sha3"
)

// const (
// LockMethodHash is hash of method signature
// LockMethodHash = sha3.Sum256([]byte(""))
// BurnMethodHash is hash of method signature
// BurnMethodHash = sha3("")
// )

// Client is a structure to sign and broadcast tx to Ethereum chain used by signer mostly
type Client struct {
	// rpcURL string
	Cdc                     *codec.Codec
	ChainID                 string
	RegistryContractAddress common.Address
	EthClient               *ethclient.Client
	ValidatorName           string
	ValidatorAddress        sdk.ValAddress
	CliCtx                  sdkContext.CLIContext
	TxBldr                  authtypes.TxBuilder
	PrivateKey              *ecdsa.PrivateKey
	EipSigner               types.EIP155Signer
	BridgeBankAddress       common.Address
	Logger                  tmLog.Logger
}

// NewClient init an Ethereum client
func NewClient(inBuf io.Reader, rpcURL string, cdc *codec.Codec, validatorMoniker string, chainID string,
	registryContractAddress common.Address, privateKey *ecdsa.PrivateKey, log tmLog.Logger) (Client, error) {
	ethClient, err := ethclient.Dial(rpcURL)
	if ethClient == nil {
		fmt.Println("Failed to new a client.")
	}
	if err != nil {
		return Client{}, err
	}

	chainIDInt, err := strconv.ParseInt(chainID, 10, 64)
	if err != nil {
		return Client{}, err
	}

	eipSigner := types.NewEIP155Signer(big.NewInt(chainIDInt))

	LockMethodHash := sha3.Sum256([]byte(""))

	fmt.Println(LockMethodHash)

	// Load validator details
	validatorAddress, validatorName, err := LoadValidatorCredentials(validatorMoniker, inBuf)
	if err != nil {
		return Client{}, err
	}

	// Load CLI context and Tx builder
	cliCtx := LoadTendermintCLIContext(cdc, validatorAddress, validatorName, rpcURL, chainID)
	txBldr := authtypes.NewTxBuilderFromCLI(nil).
		WithTxEncoder(utils.GetTxEncoder(cdc)).
		WithChainID(chainID)

	bridgeBankAddress, err := txs.GetAddressFromBridgeRegistry(ethClient, registryContractAddress, txs.BridgeBank)
	// cosmosBridgeAddress, err := txs.GetAddressFromBridgeRegistry(ethClient, registryContractAddress, txs.CosmosBridge)

	fmt.Println(bridgeBankAddress, err)
	return Client{
		Cdc:                     cdc,
		ChainID:                 chainID,
		RegistryContractAddress: registryContractAddress,
		EthClient:               ethClient,
		ValidatorName:           validatorName,
		ValidatorAddress:        validatorAddress,
		CliCtx:                  cliCtx,
		TxBldr:                  txBldr,
		PrivateKey:              privateKey,
		EipSigner:               eipSigner,
		BridgeBankAddress:       bridgeBankAddress,
		Logger:                  log,
	}, nil
}

// Start the scan
func (c *Client) Start() {
	for {
		select {
		default:
			block, _ := c.GetBlock()
			// blockNumber := block.
			fmt.Println(block.Hash)
			time.Sleep(5 * time.Second)
		}
	}
}

// GetBlock to get block from chain
func (c *Client) GetBlock() (*types.Block, error) {
	block, err := c.EthClient.BlockByNumber(context.Background(), nil)
	if err != nil {
		c.Logger.Error("he")
	}

	for _, tx := range block.Transactions() {
		// sender, _ := c.EipSigner.Sender(tx)
		//fmt.Println(sender)
		// 1 get the to address from tx
		// 2 get the method signature and compare with hash of lock and burn
		// 3 get the parameter from tx accordingly
		// wrap up the tx and sign then send to cosmos

	}
	return block, nil
}

// GetBlockHeight return current block height
// func (c *Client) GetBlockHeight() (*big.Int, error) {
// 	number, err := c.EthClient.BlockNumber(context.Background())
// 	if err != nil {
// 		c.Logger.Error("he")
// 	}

// 	return number, nil
// }

// GetTxs fetch transactions from block
// func (c *Client) GetTxs(block *types.Block) {
// 	txs := block.transactions
// for _, tx range txs {
// 	tx.Recipient
// }
// filter the tx according to address. just remain lock and burn
// both the destination address is bridge bank

// }

// func (c *Client) getTransactionsFromBlock(block *etypes.Block) ([]string, error) {
// 	txs := make([]string, 0)
// 	for _, tx := range block.Transactions() {
// 		bytes, err := tx.MarshalJSON()
// 		if err != nil {
// 			return nil, fmt.Errorf("fail to marshal tx from block: %w", err)
// 		}
// 		txs = append(txs, string(bytes))
// 	}
// 	return txs, nil
// }

// func (c *Client) fromTxToTxIn(tx *etypes.Transaction) (*stypes.TxInItem, error) {
// 	eipSigner = etypes.NewEIP155Signer(big.NewInt(int64(c.chainID)))

// 	txInItem := &stypes.TxInItem{
// 		Tx: tx.Hash().Hex()[2:],
// 	}
// 	// tx data field bytes should be hex encoded byres string as outboud or yggradsil- or migrate or yggdrasil+, etc
// 	txInItem.Memo = string(tx.Data())

// 	sender, err := eipSigner.Sender(tx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	txInItem.Sender = strings.ToLower(sender.String())
// 	if tx.To() == nil {
// 		return nil, err
// 	}
// 	txInItem.To = strings.ToLower(tx.To().String())

// 	asset, err := common.NewAsset("ETH.ETH")
// 	if err != nil {
// 		e.errCounter.WithLabelValues("fail_create_ticker", "ETH").Inc()
// 		return nil, fmt.Errorf("fail to create asset, ETH is not valid: %w", err)
// 	}
// 	txInItem.Coins = append(txInItem.Coins, common.NewCoin(asset, cosmos.NewUintFromBigInt(tx.Value())))
// 	txInItem.Gas = e.getGasUsed(tx.Hash().Hex())
// 	return txInItem, nil
// }

// LoadValidatorCredentials : loads validator's credentials (address, moniker, and passphrase)
func LoadValidatorCredentials(validatorFrom string, inBuf io.Reader) (sdk.ValAddress, string, error) {
	// Get the validator's name and account address using their moniker
	validatorAccAddress, validatorName, err := sdkContext.GetFromFields(inBuf, validatorFrom, false)
	if err != nil {
		return sdk.ValAddress{}, "", err
	}
	validatorAddress := sdk.ValAddress(validatorAccAddress)

	// Confirm that the key is valid
	_, err = authtxb.MakeSignature(nil, validatorName, keys.DefaultKeyPass, authtxb.StdSignMsg{})
	if err != nil {
		return sdk.ValAddress{}, "", err
	}

	return validatorAddress, validatorName, nil
}

// LoadTendermintCLIContext : loads CLI context for tendermint txs
func LoadTendermintCLIContext(appCodec *amino.Codec, validatorAddress sdk.ValAddress, validatorName string,
	rpcURL string, chainID string) sdkContext.CLIContext {
	// Create the new CLI context
	cliCtx := sdkContext.NewCLIContext().
		WithCodec(appCodec).
		WithFromAddress(sdk.AccAddress(validatorAddress)).
		WithFromName(validatorName)

	if rpcURL != "" {
		cliCtx = cliCtx.WithNodeURI(rpcURL)
	}
	cliCtx.SkipConfirm = true

	// Confirm that the validator's address exists
	accountRetriever := authtypes.NewAccountRetriever(cliCtx)
	err := accountRetriever.EnsureExists((sdk.AccAddress(validatorAddress)))
	if err != nil {
		log.Fatal(err)
	}
	return cliCtx
}
