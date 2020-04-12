package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const FUND_OBJECT_TYPE = "fund"

const (
	BAD_REQUEST        = "BAD_REQUEST"
	SERVER_ERROR       = "SERVER_ERROR"
	FUND_CREATE_FAILED = "FUND_CREATE_FAILED"
	FUND_NOT_FOUND     = "FUND_NOT_FOUND"
	FUND_DELETE_FAILED = "FUND_DELETE_FAILED"
	FUND_SELL_FAILED   = "FUND_SELL_FAILED"
	BUY_FUND_FAILED    = "BUY_FUND_FAILED"
	INSUFFICENT_FUND   = "INSUFFICENT_FUND"
	FUND_UPDATE_FAILED = "FUND_UPDATE_FAILED"
)

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
	fmt.Println("Create new fund...", args)

	//Checking access permission
	err := s.checkAccessPermission(ctx, "CreateFund")
	if err != nil {
		fmt.Errorf("Authorization failed", err.Error())
		return err
	}

	//Paring input
	fund := Fund{}
	err = json.Unmarshal([]byte(args), &fund)
	if err != nil {
		fmt.Errorf("Failed to parse fund create input. %s", err.Error())
		return fmt.Errorf(FUND_CREATE_FAILED)
	}

	//check fund already exists before
	_, err = s.GetFund(ctx, fund.ID)
	if err != nil && err.Error() != FUND_NOT_FOUND {
		return err
	}

	// check all assets exists before creating fund
	for _, asset := range fund.Assets {
		_, err := s.GetAsset(ctx, asset)
		if err != nil {
			fmt.Errorf(err.Error())
			return fmt.Errorf(FUND_CREATE_FAILED)
		}
	}

	assetAsBytes, err := json.Marshal(fund)
	if err != nil {
		fmt.Errorf("Failed to serialize fund data. %s", err.Error())
		return fmt.Errorf(FUND_CREATE_FAILED)
	}

	key, _ := ctx.GetStub().CreateCompositeKey(FUND_OBJECT_TYPE, []string{fund.ID})
	fmt.Println("key::", key)
	return ctx.GetStub().PutState(key, assetAsBytes)
}

// delete fund
func (s *SmartContract) DeleteFund(ctx contractapi.TransactionContextInterface, fundId string) error {
	fmt.Println("Deleting fund...", fundId)

	//Checking access permission
	err := s.checkAccessPermission(ctx, "DeleteFund")
	if err != nil {
		fmt.Errorf("Authorization failed", err.Error())
		return err
	}

	key, _ := ctx.GetStub().CreateCompositeKey(FUND_OBJECT_TYPE, []string{fundId})
	fmt.Println("key::", key)

	//Check fund exists
	_, err = s.GetFund(ctx, fundId)
	if err != nil {
		return err
	}

	// Deleting fund
	err = ctx.GetStub().DelState(key)
	if err != nil {
		fmt.Errorf("Failed to delete fund from world state. %s", err.Error())
		return fmt.Errorf(FUND_DELETE_FAILED)
	}

	return nil
}

// GetFund returns the fund stored in the world state with id
func (s *SmartContract) GetFund(ctx contractapi.TransactionContextInterface, fundId string) (*Fund, error) {
	fmt.Println("FetFund", fundId)

	//Checking access permission
	err := s.checkAccessPermission(ctx, "GetFund")
	if err != nil {
		fmt.Errorf("Authorization failed", err.Error())
		return nil, err
	}

	key, _ := ctx.GetStub().CreateCompositeKey(FUND_OBJECT_TYPE, []string{fundId})
	fmt.Println("key::", key)

	fundAsBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		fmt.Errorf("Failed to read fund from world state. %s", err.Error())
		return nil, fmt.Errorf(SERVER_ERROR)
	}

	if fundAsBytes == nil {
		fmt.Errorf("fund %s does not exist", fundId)
		return nil, fmt.Errorf(FUND_NOT_FOUND)
	}

	fund := new(Fund)
	_ = json.Unmarshal(fundAsBytes, fund)

	return fund, nil
}

// QueryAllFunds returns all funds
func (s *SmartContract) QueryAllFunds(ctx contractapi.TransactionContextInterface) ([]*Fund, error) {
	fmt.Println("QueryAllFunds called...")

	//Checking access permission
	err := s.checkAccessPermission(ctx, "QueryAllFunds")
	if err != nil {
		fmt.Errorf("Authorization failed", err.Error())
		return nil, err
	}

	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(FUND_OBJECT_TYPE, []string{})
	if err != nil {
		return nil, fmt.Errorf(SERVER_ERROR)
	}
	defer resultsIterator.Close()

	results := []*Fund{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf(SERVER_ERROR)
		}

		fund := new(Fund)
		_ = json.Unmarshal(queryResponse.Value, fund)

		//queryResult := FundQueryResult{Key: queryResponse.Key, Record: fund}
		results = append(results, fund)
	}

	return results, nil
}

// increaseFundCount
func (s *SmartContract) sellFund(ctx contractapi.TransactionContextInterface, fundId string) error {
	fmt.Println("increaseFundCount ...", fundId)

	//check fund already exists before
	fund, err := s.GetFund(ctx, fundId)
	if err != nil {
		return err
	}

	//Updating fund count
	fund.count = fund.count + 1

	assetAsBytes, err := json.Marshal(fund)
	if err != nil {
		fmt.Errorf("Failed to serialize fund data. %s", err.Error())
		return fmt.Errorf(FUND_CREATE_FAILED)
	}

	key, _ := ctx.GetStub().CreateCompositeKey(FUND_OBJECT_TYPE, []string{fund.ID})
	fmt.Println("key::", key)
	return ctx.GetStub().PutState(key, assetAsBytes)
}

// buyFund
func (s *SmartContract) buyFund(ctx contractapi.TransactionContextInterface, fundId string) error {
	fmt.Println("decreaseFundCount ...", fundId)

	//check fund already exists before
	fund, err := s.GetFund(ctx, fundId)
	if err != nil {
		return err
	}

	//Updating fund count
	if fund.count < 1 {
		fmt.Errorf("Failed to serialize fund data. %s", err.Error())
		return fmt.Errorf(INSUFFICENT_FUND)
	}

	fund.count = fund.count - 1

	assetAsBytes, err := json.Marshal(fund)
	if err != nil {
		fmt.Errorf("Failed to serialize fund data. %s", err.Error())
		return fmt.Errorf(FUND_UPDATE_FAILED)
	}

	key, _ := ctx.GetStub().CreateCompositeKey(FUND_OBJECT_TYPE, []string{fund.ID})
	fmt.Println("key::", key)
	return ctx.GetStub().PutState(key, assetAsBytes)
}

func deleteFromSlice(list []string, item string) []string {
	index := sort.SearchStrings(list, item)
	list[index] = list[len(list)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return list[:len(list)-1]
}
