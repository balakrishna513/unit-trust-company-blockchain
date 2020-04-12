package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const ASSET_OBJECT_TYPE = "asset"

//Asset struct
type Asset struct {
	DocType   string `json:"docType"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	AssetType string `json:"assetType"`
}

const (
	ASSET_NOT_FOUND     = "ASSET_NOT_FOUND"
	ASSET_CREATE_FAILED = "ASSET_CREATE_FAILED"
)

// create new asset
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, args string) error {
	fmt.Println("Create new asset...", args)

	//Checking access permission
	err := s.checkAccessPermission(ctx, "CreateAsset")
	if err != nil {
		fmt.Errorf("Authorization failed", err.Error())
		return err
	}

	var asset Asset
	err = json.Unmarshal([]byte(args), &asset)
	if err != nil {
		fmt.Errorf("Failed to parse asset create input. %s", err.Error())
		return fmt.Errorf(ASSET_CREATE_FAILED)
	}

	//check asset exists before
	_, err = s.GetAsset(ctx, asset.ID)
	if err != nil && err.Error() != ASSET_NOT_FOUND {
		return err
	}

	assetAsBytes, err := json.Marshal(asset)
	if err != nil {
		fmt.Errorf("Failed to serialize asset data. %s", err.Error())
		return fmt.Errorf(ASSET_CREATE_FAILED)
	}

	key, _ := ctx.GetStub().CreateCompositeKey(ASSET_OBJECT_TYPE, []string{asset.ID})
	fmt.Println("key::", key)
	err = ctx.GetStub().PutState(key, assetAsBytes)
	if err != nil {
		fmt.Errorf("Failed to record asset data. %s", err.Error())
		return fmt.Errorf(ASSET_CREATE_FAILED)
	}

	return nil
}

// GetAsset returns the asset stored in the world state with id
func (s *SmartContract) GetAsset(ctx contractapi.TransactionContextInterface, assetId string) (*Asset, error) {
	fmt.Println("GetAsset started...", assetId)

	//Checking access permission
	err := s.checkAccessPermission(ctx, "GetAsset")
	if err != nil {
		fmt.Errorf("Authorization failed", err.Error())
		return nil, err
	}

	key, _ := ctx.GetStub().CreateCompositeKey(ASSET_OBJECT_TYPE, []string{assetId})
	fmt.Println("key::", key)

	// Reading asset
	assetAsBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		fmt.Errorf("Failed to read asset from world state. %s", err.Error())
		return nil, fmt.Errorf(SERVER_ERROR)
	}

	if assetAsBytes == nil {
		fmt.Errorf("asset %s does not exist", assetId)
		return nil, fmt.Errorf(ASSET_NOT_FOUND)
	}

	asset := new(Asset)
	_ = json.Unmarshal(assetAsBytes, asset)

	return asset, nil
}

// QueryAllAssets returns all assets
func (s *SmartContract) QueryAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	fmt.Println("key::", ASSET_OBJECT_TYPE, []string{})

	//Checking access permission
	err := s.checkAccessPermission(ctx, "QueryAllAssets")
	if err != nil {
		fmt.Errorf("Authorization failed", err.Error())
		return nil, err
	}

	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(ASSET_OBJECT_TYPE, []string{})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []*Asset{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			fmt.Println("failed to iterate over resultset", err)
			return nil, fmt.Errorf(SERVER_ERROR)
		}

		asset := new(Asset)
		_ = json.Unmarshal(queryResponse.Value, asset)

		//queryResult := AssetQueryResult{Key: queryResponse.Key, Record: asset}
		results = append(results, asset)
	}

	return results, nil
}
