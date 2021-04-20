/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

// Ballot stores a value
type Ballot struct {
	Id				string 	`json:"id"`
	ElectionId 		string 	`json:"election_id"`
	DistrictId      int		`json:"district_id"`
	CandidateId   	string 	`json:"candidate_id"`
	Timestamp  		int 	`json:"timestamp"`
}

