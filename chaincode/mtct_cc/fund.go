package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const FUND_OBJECT_TYPE = "fund"

// QueryResult structure used for handling result of query
type FundQueryResult struct {
	Key    string `json:"Key"`
	Record *Fund
}

//Fund struct
type Fund struct {
	DocType string   `json:"docType"`
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	count   int      `json:"count"`
	Assets  []string `json:"assets"`
}

// create new fund
func (s *SmartContract) CreateFund(ctx contractapi.TransactionContextInterface, args string) error {
	fmt.Println("Create new fund...")
	var fund Fund
	err := json.Unmarshal([]byte(args), &fund)
	if err != nil {
		return fmt.Errorf("Failed to parse fund create input. %s", err.Error())
	}

	// check all assets exists before creating fund
	for _, asset := range fund.Assets {
		_, err := s.GetAsset(ctx, asset)
		if err != nil {
			return fmt.Errorf("Asset %s not found %s", asset, err.Error())
		}
	}

	assetAsBytes, err := json.Marshal(fund)
	if err != nil {
		return fmt.Errorf("Failed to serialize fund data. %s", err.Error())
	}

	key, _ := ctx.GetStub().CreateCompositeKey(FUND_OBJECT_TYPE, []string{fund.ID})

	return ctx.GetStub().PutState(key, assetAsBytes)
}

// delete fund
func (s *SmartContract) DeleteFund(ctx contractapi.TransactionContextInterface, fundId string) error {
	fmt.Println("Deleting fund...", fundId)
	key, _ := ctx.GetStub().CreateCompositeKey(FUND_OBJECT_TYPE, []string{fundId})

	// Deleting fund
	err := ctx.GetStub().DelState(key)
	if err != nil {
		return fmt.Errorf("Failed to delete fund from world state. %s", err.Error())
	}

	return nil
}

// GetFund returns the fund stored in the world state with id
func (s *SmartContract) GetFund(ctx contractapi.TransactionContextInterface, fundId string) (*Fund, error) {
	key, _ := ctx.GetStub().CreateCompositeKey(FUND_OBJECT_TYPE, []string{fundId})
	assetAsBytes, err := ctx.GetStub().GetState(key)

	if err != nil {
		return nil, fmt.Errorf("Failed to read fund from world state. %s", err.Error())
	}

	if assetAsBytes == nil {
		return nil, fmt.Errorf("fund %s does not exist", fundId)
	}

	fund := new(Fund)
	_ = json.Unmarshal(assetAsBytes, fund)

	return fund, nil
}

// QueryAllFunds returns all funds
func (s *SmartContract) QueryAllFunds(ctx contractapi.TransactionContextInterface) ([]FundQueryResult, error) {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(FUND_OBJECT_TYPE, []string{})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []FundQueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		fund := new(Fund)
		_ = json.Unmarshal(queryResponse.Value, fund)

		queryResult := FundQueryResult{Key: queryResponse.Key, Record: fund}
		results = append(results, queryResult)
	}

	return results, nil
}

func deleteFromSlice(list []string, item string) []string {
	index := sort.SearchStrings(list, item)
	list[index] = list[len(list)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return list[:len(list)-1]
}
