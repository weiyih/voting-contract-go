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
func (s *BallotContract) BallotExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	hashId := hash(id)
	data, err := ctx.GetStub().GetState(hashId)

	if err != nil {
		return false, err
	}

	return data != nil, nil
}

// CreateBallot
// TransactionContextInterface defines the interface which TransactionContext meets. This can be taken by transacton functions on a contract
// which has not set a custom transaction context to allow transaction functions to take an interface to simplify unit testing.
// https://godoc.org/github.com/hyperledger/fabric-contract-api-go/contractapi#TransactionContextInterface
func (s *BallotContract) CreateBallot(ctx contractapi.TransactionContextInterface, id string, electionId string, districtId string, candidateId string, timestamp string) error {
	// Checks if ballot already exists
	exists, err := s.BallotExists(ctx, id)
	if err != nil {
		return fmt.Errorf("could not read from world state: %s", err)
	} else if exists {
		return fmt.Errorf("ballot already exists")
	}

	hashId := hash(id)
	// Create vote Object
	ballot := Ballot{
		Id:          hashId,
		ElectionId:  electionId,
		DistrictId:  districtId,
		CandidateId: candidateId,
		Timestamp:   timestamp,
	}

	ballotJSON, err := json.Marshal(ballot)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(hashId, ballotJSON)
}

// Returns the ballot if it exists
func (s *BallotContract) ReadBallot(ctx contractapi.TransactionContextInterface, id string) (*Ballot, error) {
	exists, err := s.BallotExists(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("the ballot does not exist")
	}

	hashId := hash(id)
	bytes, _ := ctx.GetStub().GetState(hashId)

	ballot := new(Ballot)

	err = json.Unmarshal(bytes, ballot)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal world state data to type Ballot")
	}

	return ballot, nil
}

// Returns all ballots found in world state
func (s *BallotContract) GetAllBallots(ctx contractapi.TransactionContextInterface) ([]*Ballot, error) {

	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var ballots []*Ballot

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var ballot Ballot

		err = json.Unmarshal(queryResponse.Value, &ballot)
		if err != nil {
			return nil, err
		}
		ballots = append(ballots, &ballot)
	}

	return ballots, nil
}

// https://softwareengineering.stackexchange.com/questions/49550/which-hashing-algorithm-is-best-for-uniqueness-and-speed
// https://golang.org/pkg/hash/crc32
func hash(stringId string) string {
	// Koopman's polynomial.
	// Also has better error detection characteristics than IEEE.
	// https://dx.doi.org/10.1109/DSN.2002.1028931
	koopmanTable := crc32.MakeTable(crc32.Koopman)
	b := []byte(stringId)
	checksum := crc32.Checksum(b, koopmanTable)
	hexString := fmt.Sprintf("%x", checksum)
	return hexString
}
