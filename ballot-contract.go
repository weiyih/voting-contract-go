/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"
	"hash/crc32"

	"github.com/hyperledger/fabric-contract-api-go/contractapi" // https://godoc.org/github.com/hyperledger/fabric-contract-api-go
)

// BallotContract contract for managing CRUD for Ballot
type BallotContract struct {
	contractapi.Contract
}

// BallotExists returns true when asset with given ID exists in world state
func (s *SmartContract) BallotExists(ctx contractapi.TransactionContextInterface, ballotID string) (bool, error) {
	hashId = hash(ballotID)
	assetJSON, err := ctx.GetStub().GetState(hashId)
	if err != nil {
		return false, fmt.Errorf("Failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// CreateBallot
// TransactionContextInterface defines the interface which TransactionContext meets. This can be taken by transacton functions on a contract
// which has not set a custom transaction context to allow transaction functions to take an interface to simplify unit testing.
// https://godoc.org/github.com/hyperledger/fabric-contract-api-go/contractapi#TransactionContextInterface
func (s *SmartContract) CreateBallot(ctx contractapi.TransactionContextInterface, id string, electionId string, districtId string, candidateId string, timestamp string) error {
	// Checks if ballot already exists
	exists, err := s.BallotExists(ctx, id)
	if err != nil {
		return fmt.Errorf("Error - Could not read from world state. %s", err)
	}
	if exists {
		return fmt.Errorf("Error - Ballot %s already exists", id)
	}
	
	hashId := hash(id)
	// Create vote Object
	ballot := Ballot{
		Id: 			hashId,
		ElectionId: 	electionId,
		DistrictId:     districtId,
		CandidateId:   	candidateId,
		Timestamp:  	timestamp,
	}

	ballotJSON, err := json.Marshal(ballot)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(hashId, ballotJSON)
}

// Returns all ballots found in world state
func (s *SmartContract) GetAllBallots(ctx contractapi.TransactionContextInterface) ([]*Ballot, error) {

	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []*Ballot

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var ballot Ballot

		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return results, nil
}

// https://softwareengineering.stackexchange.com/questions/49550/which-hashing-algorithm-is-best-for-uniqueness-and-speed
// https://golang.org/pkg/hash/crc32
func hash(id string) string {
	// Koopman's polynomial.
    // Also has better error detection characteristics than IEEE.
    // https://dx.doi.org/10.1109/DSN.2002.1028931
	koopmanTable := crc32.MakeTable(crc32.Koopman)

	b := []byte(id)
	hash := crc32.checksum(b, koopmanTable)
	return hash
}