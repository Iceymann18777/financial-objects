/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const getStateError = "world state get error"

type MockStub struct {
	shim.ChaincodeStubInterface
	mock.Mock
}

func (ms *MockStub) GetState(key string) ([]byte, error) {
	args := ms.Called(key)

	return args.Get(0).([]byte), args.Error(1)
}

func (ms *MockStub) PutState(key string, value []byte) error {
	args := ms.Called(key, value)

	return args.Error(0)
}

func (ms *MockStub) DelState(key string) error {
	args := ms.Called(key)

	return args.Error(0)
}

type MockContext struct {
	contractapi.TransactionContextInterface
	mock.Mock
}

func (mc *MockContext) GetStub() shim.ChaincodeStubInterface {
	args := mc.Called()

	return args.Get(0).(*MockStub)
}

func configureStub() (*MockContext, *MockStub) {
	var nilBytes []byte

	testCryptoMyAsset := new(CryptoMyAsset)
	testCryptoMyAsset.Value = "set value"
	cryptoMyAssetBytes, _ := json.Marshal(testCryptoMyAsset)

	ms := new(MockStub)
	ms.On("GetState", "statebad").Return(nilBytes, errors.New(getStateError))
	ms.On("GetState", "missingkey").Return(nilBytes, nil)
	ms.On("GetState", "existingkey").Return([]byte("some value"), nil)
	ms.On("GetState", "cryptoMyAssetkey").Return(cryptoMyAssetBytes, nil)
	ms.On("PutState", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
	ms.On("DelState", mock.AnythingOfType("string")).Return(nil)

	mc := new(MockContext)
	mc.On("GetStub").Return(ms)

	return mc, ms
}

func TestCryptoMyAssetExists(t *testing.T) {
	var exists bool
	var err error

	ctx, _ := configureStub()
	c := new(CryptoMyAssetContract)

	exists, err = c.CryptoMyAssetExists(ctx, "statebad")
	assert.EqualError(t, err, getStateError)
	assert.False(t, exists, "should return false on error")

	exists, err = c.CryptoMyAssetExists(ctx, "missingkey")
	assert.Nil(t, err, "should not return error when can read from world state but no value for key")
	assert.False(t, exists, "should return false when no value for key in world state")

	exists, err = c.CryptoMyAssetExists(ctx, "existingkey")
	assert.Nil(t, err, "should not return error when can read from world state and value exists for key")
	assert.True(t, exists, "should return true when value for key in world state")
}

func TestCreateCryptoMyAsset(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(CryptoMyAssetContract)

	err = c.CreateCryptoMyAsset(ctx, "statebad", "some value")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

	err = c.CreateCryptoMyAsset(ctx, "existingkey", "some value")
	assert.EqualError(t, err, "The asset existingkey already exists", "should error when exists returns true")

	err = c.CreateCryptoMyAsset(ctx, "missingkey", "some value")
	stub.AssertCalled(t, "PutState", "missingkey", []byte("{\"value\":\"some value\"}"))
}

func TestReadCryptoMyAsset(t *testing.T) {
	var cryptoMyAsset *CryptoMyAsset
	var err error

	ctx, _ := configureStub()
	c := new(CryptoMyAssetContract)

	cryptoMyAsset, err = c.ReadCryptoMyAsset(ctx, "statebad")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when reading")
	assert.Nil(t, cryptoMyAsset, "should not return CryptoMyAsset when exists errors when reading")

	cryptoMyAsset, err = c.ReadCryptoMyAsset(ctx, "missingkey")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when reading")
	assert.Nil(t, cryptoMyAsset, "should not return CryptoMyAsset when key does not exist in world state when reading")

	cryptoMyAsset, err = c.ReadCryptoMyAsset(ctx, "existingkey")
	assert.EqualError(t, err, "Could not unmarshal world state data to type CryptoMyAsset", "should error when data in key is not CryptoMyAsset")
	assert.Nil(t, cryptoMyAsset, "should not return CryptoMyAsset when data in key is not of type CryptoMyAsset")

	cryptoMyAsset, err = c.ReadCryptoMyAsset(ctx, "cryptoMyAssetkey")
	expectedCryptoMyAsset := new(CryptoMyAsset)
	expectedCryptoMyAsset.Value = "set value"
	assert.Nil(t, err, "should not return error when CryptoMyAsset exists in world state when reading")
	assert.Equal(t, expectedCryptoMyAsset, cryptoMyAsset, "should return deserialized CryptoMyAsset from world state")
}

func TestUpdateCryptoMyAsset(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(CryptoMyAssetContract)

	err = c.UpdateCryptoMyAsset(ctx, "statebad", "new value")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when updating")

	err = c.UpdateCryptoMyAsset(ctx, "missingkey", "new value")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when updating")

	err = c.UpdateCryptoMyAsset(ctx, "cryptoMyAssetkey", "new value")
	expectedCryptoMyAsset := new(CryptoMyAsset)
	expectedCryptoMyAsset.Value = "new value"
	expectedCryptoMyAssetBytes, _ := json.Marshal(expectedCryptoMyAsset)
	assert.Nil(t, err, "should not return error when CryptoMyAsset exists in world state when updating")
	stub.AssertCalled(t, "PutState", "cryptoMyAssetkey", expectedCryptoMyAssetBytes)
}

func TestDeleteCryptoMyAsset(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(CryptoMyAssetContract)

	err = c.DeleteCryptoMyAsset(ctx, "statebad")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

	err = c.DeleteCryptoMyAsset(ctx, "missingkey")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when deleting")

	err = c.DeleteCryptoMyAsset(ctx, "cryptoMyAssetkey")
	assert.Nil(t, err, "should not return error when CryptoMyAsset exists in world state when deleting")
	stub.AssertCalled(t, "DelState", "cryptoMyAssetkey")
}
