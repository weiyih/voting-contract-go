// Based on sample from
// https://github.com/hyperledger/fabric-samples/blob/master/chaincode/fabcar/go/fabcar.go

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"encoding/json"
	"fmt"
	"hash/crc32"

	"github.com/hyperledger/fabric-contract-api-go/contractapi" // https://godoc.org/github.com/hyperledger/fabric-contract-api-go
)

// SmartContract provides functions for CRUD
type SmartContract struct {
	contractapi.Contract
}

// Vote Object
type Ballot struct {
	Id				string `json:"id"`
	ElectionId 		string `json:"election_id"`
	DistrictId      string `json:"district_id"`
	CandidateId   	string `json:"candidate_id"`
	Timestamp  		string `json:"timestamp"`
}

// QueryResult handle result of Vote query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Vote
}

// CreateBallot
// TransactionContextInterface defines the interface which TransactionContext meets. This can be taken by transacton functions on a contract
// which has not set a custom transaction context to allow transaction functions to take an interface to simplify unit testing.
// https://godoc.org/github.com/hyperledger/fabric-contract-api-go/contractapi#TransactionContextInterface
func (s *SmartContract) CreateBallot(ctx contractapi.TransactionContextInterface, id string, electionId string, districtId string, candidateId string, timestamp string) error {
	// Checks if ballot already exists
	exists, err := s.BallotExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("Error - ballot %s already exists", id)
	}
	
	hashId := hash(id)
	// Create vote Object
	vote := Vote{
		Id: 			hashId,
		ElectionId: 	electionId,
		DistrictId:     districtId,
		CandidateId:   	candidateId,
		Timestamp:  	timestamp,
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// Returns all ballots found in world state
func (s *SmartContract) GetAllBallots(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {

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

func (s *SmartContract) BallotExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("Failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
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

func main() {

	smartContract := new(SmartContract)

	chaincode, err := contractapi.NewChaincode(smartContract)

	if err != nil {
		fmt.Printf("Error create chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s", err.Error())
	}
}
