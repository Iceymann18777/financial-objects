/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/metadata"
)

func main() {
	cryptoMyAssetContract := new(CryptoMyAssetContract)
	cryptoMyAssetContract.Info.Version = "0.0.1"
	cryptoMyAssetContract.Info.Description = "My Smart Contract"
	cryptoMyAssetContract.Info.License = new(metadata.LicenseMetadata)
	cryptoMyAssetContract.Info.License.Name = "Apache-2.0"
	cryptoMyAssetContract.Info.Contact = new(metadata.ContactMetadata)
	cryptoMyAssetContract.Info.Contact.Name = "John Doe"

	chaincode, err := contractapi.NewChaincode(cryptoMyAssetContract)
	chaincode.Info.Title = "pypi-codegen chaincode"
	chaincode.Info.Version = "0.0.1"

	if err != nil {
		panic("Could not create chaincode from CryptoMyAssetContract." + err.Error())
	}

	err = chaincode.Start()

	if err != nil {
		panic("Failed to start chaincode. " + err.Error())
	}
}
