package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	ethClient  *ethclient.Client
	contract   *bind.BoundContract
	address    common.Address
	privateKey *ecdsa.PrivateKey
}

func InitClient(rpcURL, contractAddressHex, privateKeyHex, abiPath string) (*Client, error) {
	ethClient, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to Ethereum node: %w", err)
	}

	address := common.HexToAddress(contractAddressHex)

	abiBytes, err := os.ReadFile(abiPath)
	if err != nil {
		return nil, fmt.Errorf("error reading ABI: %w", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		return nil, fmt.Errorf("error parsing ABI: %w", err)
	}

	privKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("error loading private key: %w", err)
	}

	contract := bind.NewBoundContract(address, parsedABI, ethClient, ethClient, ethClient)

	return &Client{
		ethClient:  ethClient,
		contract:   contract,
		address:    address,
		privateKey: privKey,
	}, nil
}

func (c *Client) GetValue(ctx context.Context) (*big.Int, error) {
	var out []any
	err := c.contract.Call(&bind.CallOpts{Context: ctx}, &out, "get")
	if err != nil {
		return nil, fmt.Errorf("error executing 'get': %w", err)
	}

	val, ok := out[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected type in 'get' return")
	}
	return val, nil
}

func (c *Client) SetValue(ctx context.Context, value *big.Int) (string, error) {
	chainID, err := c.ethClient.ChainID(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting chain ID: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, chainID)
	if err != nil {
		return "", fmt.Errorf("error creating transactor: %w", err)
	}

	tx, err := c.contract.Transact(auth, "set", value)
	if err != nil {
		return "", fmt.Errorf("error sending transaction: %w", err)
	}

	slog.Info("Transaction sent", "tx", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}
