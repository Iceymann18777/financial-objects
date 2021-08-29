/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// CryptoMyAssetContract contract for managing CRUD for CryptoMyAsset
type CryptoMyAssetContract struct {
	contractapi.Contract
}

// CryptoMyAssetExists returns true when asset with given ID exists in world state
func (c *CryptoMyAssetContract) CryptoMyAssetExists(ctx contractapi.TransactionContextInterface, cryptoMyAssetID string) (bool, error) {
	data, err := ctx.GetStub().GetState(cryptoMyAssetID)

	if err != nil {
		return false, err
	}

	return data != nil, nil
}

// CreateCryptoMyAsset creates a new instance of CryptoMyAsset
func (c *CryptoMyAssetContract) CreateCryptoMyAsset(ctx contractapi.TransactionContextInterface, cryptoMyAssetID string, value string) error {
	exists, err := c.CryptoMyAssetExists(ctx, cryptoMyAssetID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if exists {
		return fmt.Errorf("The asset %s already exists", cryptoMyAssetID)
	}

	cryptoMyAsset := new(CryptoMyAsset)
	cryptoMyAsset.Value = value

	bytes, _ := json.Marshal(cryptoMyAsset)

	return ctx.GetStub().PutState(cryptoMyAssetID, bytes)
}

// ReadCryptoMyAsset retrieves an instance of CryptoMyAsset from the world state
func (c *CryptoMyAssetContract) ReadCryptoMyAsset(ctx contractapi.TransactionContextInterface, cryptoMyAssetID string) (*CryptoMyAsset, error) {
	exists, err := c.CryptoMyAssetExists(ctx, cryptoMyAssetID)
	if err != nil {
		return nil, fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("The asset %s does not exist", cryptoMyAssetID)
	}

	bytes, _ := ctx.GetStub().GetState(cryptoMyAssetID)

	cryptoMyAsset := new(CryptoMyAsset)

	err = json.Unmarshal(bytes, cryptoMyAsset)

	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal world state data to type CryptoMyAsset")
	}

	return cryptoMyAsset, nil
}

// UpdateCryptoMyAsset retrieves an instance of CryptoMyAsset from the world state and updates its value
func (c *CryptoMyAssetContract) UpdateCryptoMyAsset(ctx contractapi.TransactionContextInterface, cryptoMyAssetID string, newValue string) error {
	exists, err := c.CryptoMyAssetExists(ctx, cryptoMyAssetID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("The asset %s does not exist", cryptoMyAssetID)
	}

	cryptoMyAsset := new(CryptoMyAsset)
	cryptoMyAsset.Value = newValue

	bytes, _ := json.Marshal(cryptoMyAsset)

	return ctx.GetStub().PutState(cryptoMyAssetID, bytes)
}

// DeleteCryptoMyAsset deletes an instance of CryptoMyAsset from the world state
func (c *CryptoMyAssetContract) DeleteCryptoMyAsset(ctx contractapi.TransactionContextInterface, cryptoMyAssetID string) error {
	exists, err := c.CryptoMyAssetExists(ctx, cryptoMyAssetID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("The asset %s does not exist", cryptoMyAssetID)
	}

	return ctx.GetStub().DelState(cryptoMyAssetID)
}
