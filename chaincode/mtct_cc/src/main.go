package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// const (
// 	BAD_REQUEST     = "BAD_REQUEST"
// 	USER_NOT_FOUND  = "USER_NOT_FOUND"
// 	ASSET_NOT_FOUND = "ASSET_NOT_FOUND"
// 	FUND_NOT_FOUND  = "FUND_NOT_FOUND"
// )

// SmartContract provides functions for managing mctc operations
type SmartContract struct {
	contractapi.Contract
}

// InitLedger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create mtct chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting mtct chaincode: %s", err.Error())
	}
}
