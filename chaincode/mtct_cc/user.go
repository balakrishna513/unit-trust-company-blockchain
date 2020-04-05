package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const STATUS_APPROVED = "approved"
const STATUS_REJECTED = "rejected"
const USER_OBJECT_TYPE = "user"

// QueryResult structure used for handling result of query
type UserQueryResult struct {
	Key    string `json:"Key"`
	Record *User
}

//user struct
type User struct {
	DocType string   `json:"docType"`
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	OrgID   string   `json:"orgId"`
	OrgName string   `json:"orgName"`
	OrgType string   `json:"orgType"`
	Role    string   `json:"role"`
	Funds   []string `json:"funds"`
	Status  string   `json:"status"`
}

// Register new user
func (s *SmartContract) RegisterUser(ctx contractapi.TransactionContextInterface, args string) error {
	fmt.Println("Registering new user...")
	var user User
	err := json.Unmarshal([]byte(args), &user)
	if err != nil {
		return fmt.Errorf("Failed to parse user register input. %s", err.Error())
	}

	userAsBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("Failed to serialize user data. %s", err.Error())
	}

	key, _ := ctx.GetStub().CreateCompositeKey(USER_OBJECT_TYPE, []string{user.ID})

	return ctx.GetStub().PutState(key, userAsBytes)
}

// Approve user by user id
func (s *SmartContract) ApproveUser(ctx contractapi.TransactionContextInterface, userId string) error {
	fmt.Println("ApproveUser...")
	//Activating user
	return updateUserStatus(ctx, userId, STATUS_APPROVED)
}

func (s *SmartContract) RejectUser(ctx contractapi.TransactionContextInterface, userId string) error {
	fmt.Println("RejectUser...")
	//Activating user
	return updateUserStatus(ctx, userId, STATUS_REJECTED)
}

// update user by user id
func updateUserStatus(ctx contractapi.TransactionContextInterface, userId string, status string) error {
	fmt.Println("UpdateUserStatus...")
	var user User
	key, _ := ctx.GetStub().CreateCompositeKey(USER_OBJECT_TYPE, []string{userId})
	userAsBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("Failed to read user from world state. %s", err.Error())
	}

	err = json.Unmarshal(userAsBytes, &user)
	if err != nil {
		return fmt.Errorf("Failed to parse user register input. %s", err.Error())
	}

	//Activating user
	user.Status = status
	userAsBytes, err = json.Marshal(user)
	if err != nil {
		return fmt.Errorf("Failed to serialize user data. %s", err.Error())
	}

	return ctx.GetStub().PutState(user.ID, userAsBytes)
}

// GetUser returns the user stored in the world state with id
func (s *SmartContract) GetUser(ctx contractapi.TransactionContextInterface, userId string) (*User, error) {
	key, _ := ctx.GetStub().CreateCompositeKey(USER_OBJECT_TYPE, []string{userId})
	userAsBytes, err := ctx.GetStub().GetState(key)

	if err != nil {
		return nil, fmt.Errorf("Failed to read user from world state. %s", err.Error())
	}

	if userAsBytes == nil {
		return nil, fmt.Errorf("user %s does not exist", userId)
	}

	user := new(User)
	_ = json.Unmarshal(userAsBytes, user)

	return user, nil
}

// QueryAllUsers returns all user
func (s *SmartContract) QueryAllUsers(ctx contractapi.TransactionContextInterface) ([]UserQueryResult, error) {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(USER_OBJECT_TYPE, []string{})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []UserQueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		user := new(User)
		_ = json.Unmarshal(queryResponse.Value, user)

		queryResult := UserQueryResult{Key: queryResponse.Key, Record: user}
		results = append(results, queryResult)
	}

	return results, nil
}

// buy existing fund
func (s *SmartContract) BuyFund(ctx contractapi.TransactionContextInterface, fundId string, userId string) error {
	fmt.Println("Buying fund...", fundId, userId)

	//Check if fund exists
	_, err := s.GetFund(ctx, fundId)
	if err != nil {
		return fmt.Errorf("Failed to retrieve Fund with id: %s. Error: %s", fundId, err.Error())
	}

	// Retrieving user
	user, err := s.GetUser(ctx, userId)
	if err != nil {
		return fmt.Errorf("Failed to retrieve user with id: %s. Error: %s", userId, err.Error())
	}

	//Updating user funds list with new fund id
	user.Funds = append(user.Funds, fundId)

	// Serializing user object
	userAsBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("Failed to serialize user data. %s", err.Error())
	}

	key, _ := ctx.GetStub().CreateCompositeKey(USER_OBJECT_TYPE, []string{userId})

	return ctx.GetStub().PutState(key, userAsBytes)
}

// sell existing fund
func (s *SmartContract) SellFund(ctx contractapi.TransactionContextInterface, fundId string, userId string) error {
	fmt.Println("Selling fund...", fundId, userId)

	//Check if fund exists
	_, err := s.GetFund(ctx, fundId)
	if err != nil {
		return fmt.Errorf("Failed to retrieve Fund with id: %s. Error: %s", fundId, err.Error())
	}

	// Retrieving user
	user, err := s.GetUser(ctx, userId)
	if err != nil {
		return fmt.Errorf("Failed to retrieve user with id: %s. Error: %s", userId, err.Error())
	}

	fundExists := StringInSlice(fundId, user.Funds)
	if !fundExists {
		return fmt.Errorf("User %s does not own fund %s. Hence cannot sell.", err.Error())
	}

	//Deleting fund from fund list
	user.Funds = deleteFromSlice(user.Funds, fundId)

	// Serializing user object
	userAsBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("Failed to serialize user data. %s", err.Error())
	}

	key, _ := ctx.GetStub().CreateCompositeKey(USER_OBJECT_TYPE, []string{userId})

	return ctx.GetStub().PutState(key, userAsBytes)
}
