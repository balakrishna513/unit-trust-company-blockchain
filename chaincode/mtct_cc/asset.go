package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const ASSET_OBJECT_TYPE = "asset"

// QueryResult structure used for handling result of query
type AssetQueryResult struct {
	Key    string `json:"Key"`
	Record *Asset
}

//Asset struct
type Asset struct {
	DocType   string `json:"docType"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	AssetType string `json:"assetType"`
}

// create new asset
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, args string) error {
	fmt.Println("Create new asset...")
	var asset Asset
	err := json.Unmarshal([]byte(args), &asset)
	if err != nil {
		return fmt.Errorf("Failed to parse asset create input. %s", err.Error())
	}

	assetAsBytes, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("Failed to serialize asset data. %s", err.Error())
	}

	key, _ := ctx.GetStub().CreateCompositeKey(ASSET_OBJECT_TYPE, []string{asset.ID})

	return ctx.GetStub().PutState(key, assetAsBytes)
}

// GetAsset returns the asset stored in the world state with id
func (s *SmartContract) GetAsset(ctx contractapi.TransactionContextInterface, assetId string) (*Asset, error) {
	key, _ := ctx.GetStub().CreateCompositeKey(ASSET_OBJECT_TYPE, []string{assetId})
	assetAsBytes, err := ctx.GetStub().GetState(key)

	if err != nil {
		return nil, fmt.Errorf("Failed to read asset from world state. %s", err.Error())
	}

	if assetAsBytes == nil {
		return nil, fmt.Errorf("asset %s does not exist", assetId)
	}

	asset := new(Asset)
	_ = json.Unmarshal(assetAsBytes, asset)

	return asset, nil
}

// QueryAllAssets returns all assets
func (s *SmartContract) QueryAllAssets(ctx contractapi.TransactionContextInterface) ([]AssetQueryResult, error) {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(ASSET_OBJECT_TYPE, []string{})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []AssetQueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		asset := new(Asset)
		_ = json.Unmarshal(queryResponse.Value, asset)

		queryResult := AssetQueryResult{Key: queryResponse.Key, Record: asset}
		results = append(results, queryResult)
	}

	return results, nil
}
