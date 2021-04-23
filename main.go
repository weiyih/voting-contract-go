/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/metadata"
)

func main() {
	ballotContract := new(BallotContract)
	ballotContract.Info.Version = "1.0.0"
	ballotContract.Info.Description = "Ballot Contract"
	ballotContract.Info.License = new(metadata.LicenseMetadata)
	ballotContract.Info.License.Name = "Apache-2.0"
	ballotContract.Info.Contact = new(metadata.ContactMetadata)
	ballotContract.Info.Contact.Name = "Kevin Wei"

	chaincode, err := contractapi.NewChaincode(ballotContract)
	chaincode.Info.Title = "vote-contract-go chaincode"
	chaincode.Info.Version = "1.0.0"

	if err != nil {
		panic("Could not create chaincode from BallotContract." + err.Error())
	}

	err = chaincode.Start()

	if err != nil {
		panic("Failed to start chaincode. " + err.Error())
	}
}
