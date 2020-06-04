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

	"github.com/hyperledger/fabric-contract-api-go/contractapi" // https://godoc.org/github.com/hyperledger/fabric-contract-api-go
)

// SmartContract provides functions for CRUD
type SmartContract struct {
	contractapi.Contract
}

// Vote Object
type Vote struct {
	ElectionID string `json:"election_id"`
	Ward       string `json:"ward"`
	Selected   string `json:"selected_candidate"`
	Timestamp  string `json:"timestamp"`
}

// QueryResult handle result of Vote query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Vote
}

// CreateVote s
// TransactionContextInterface defines the interface which TransactionContext meets. This can be taken by transacton functions on a contract
// which has not set a custom transaction context to allow transaction functions to take an interface to simplify unit testing.
// https://godoc.org/github.com/hyperledger/fabric-contract-api-go/contractapi#TransactionContextInterface
func (s *SmartContract) CreateVote(ctx contractapi.TransactionContextInterface, voteKey string, electionID string, wardID string, candidate string, timestamp string) error {
	// Create vote Object
	vote := Vote{
		ElectionID: electionID,
		Ward:       wardID,
		Selected:   candidate,
		Timestamp:  timestamp,
	}

	// Note _ replaces the error
	voteAsBytes, _ := json.Marshal(vote)

	return ctx.GetStub().PutState(voteKey, voteAsBytes)
}

// QueryAllVotes returns all votes found in world state
func (s *SmartContract) QueryAllVotes(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {

	startKey := "VOTE000"
	endKey := "VOTE999"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		vote := new(Vote)
		_ = json.Unmarshal(queryResponse.Value, vote)

		queryResult := QueryResult{Key: queryResponse.Key, Record: vote}
		results = append(results, queryResult)
	}

	return results, nil
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
