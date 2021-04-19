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

	testBallot := new(Ballot)
	testBallot.Value = "set value"
	ballotBytes, _ := json.Marshal(testBallot)

	ms := new(MockStub)
	ms.On("GetState", "statebad").Return(nilBytes, errors.New(getStateError))
	ms.On("GetState", "missingkey").Return(nilBytes, nil)
	ms.On("GetState", "existingkey").Return([]byte("some value"), nil)
	ms.On("GetState", "ballotkey").Return(ballotBytes, nil)
	ms.On("PutState", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
	ms.On("DelState", mock.AnythingOfType("string")).Return(nil)

	mc := new(MockContext)
	mc.On("GetStub").Return(ms)

	return mc, ms
}

func TestBallotExists(t *testing.T) {
	var exists bool
	var err error

	ctx, _ := configureStub()
	c := new(BallotContract)

	exists, err = c.BallotExists(ctx, "statebad")
	assert.EqualError(t, err, getStateError)
	assert.False(t, exists, "should return false on error")

	exists, err = c.BallotExists(ctx, "missingkey")
	assert.Nil(t, err, "should not return error when can read from world state but no value for key")
	assert.False(t, exists, "should return false when no value for key in world state")

	exists, err = c.BallotExists(ctx, "existingkey")
	assert.Nil(t, err, "should not return error when can read from world state and value exists for key")
	assert.True(t, exists, "should return true when value for key in world state")
}

func TestCreateBallot(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(BallotContract)

	err = c.CreateBallot(ctx, "statebad", "some value")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

	err = c.CreateBallot(ctx, "existingkey", "some value")
	assert.EqualError(t, err, "The asset existingkey already exists", "should error when exists returns true")

	err = c.CreateBallot(ctx, "missingkey", "some value")
	stub.AssertCalled(t, "PutState", "missingkey", []byte("{\"value\":\"some value\"}"))
}

func TestReadBallot(t *testing.T) {
	var ballot *Ballot
	var err error

	ctx, _ := configureStub()
	c := new(BallotContract)

	ballot, err = c.ReadBallot(ctx, "statebad")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when reading")
	assert.Nil(t, ballot, "should not return Ballot when exists errors when reading")

	ballot, err = c.ReadBallot(ctx, "missingkey")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when reading")
	assert.Nil(t, ballot, "should not return Ballot when key does not exist in world state when reading")

	ballot, err = c.ReadBallot(ctx, "existingkey")
	assert.EqualError(t, err, "Could not unmarshal world state data to type Ballot", "should error when data in key is not Ballot")
	assert.Nil(t, ballot, "should not return Ballot when data in key is not of type Ballot")

	ballot, err = c.ReadBallot(ctx, "ballotkey")
	expectedBallot := new(Ballot)
	expectedBallot.Value = "set value"
	assert.Nil(t, err, "should not return error when Ballot exists in world state when reading")
	assert.Equal(t, expectedBallot, ballot, "should return deserialized Ballot from world state")
}

func TestUpdateBallot(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(BallotContract)

	err = c.UpdateBallot(ctx, "statebad", "new value")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when updating")

	err = c.UpdateBallot(ctx, "missingkey", "new value")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when updating")

	err = c.UpdateBallot(ctx, "ballotkey", "new value")
	expectedBallot := new(Ballot)
	expectedBallot.Value = "new value"
	expectedBallotBytes, _ := json.Marshal(expectedBallot)
	assert.Nil(t, err, "should not return error when Ballot exists in world state when updating")
	stub.AssertCalled(t, "PutState", "ballotkey", expectedBallotBytes)
}

func TestDeleteBallot(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(BallotContract)

	err = c.DeleteBallot(ctx, "statebad")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

	err = c.DeleteBallot(ctx, "missingkey")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when deleting")

	err = c.DeleteBallot(ctx, "ballotkey")
	assert.Nil(t, err, "should not return error when Ballot exists in world state when deleting")
	stub.AssertCalled(t, "DelState", "ballotkey")
}
