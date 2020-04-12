package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const (
	USER_CREATE_FAILED    = "USER_CREATE_FAILED"
	USER_UPDATE_FAILED    = "USER_UPDATE_FAILED"
	USER_NOT_FOUND        = "USER_NOT_FOUND"
	USER_RETRIVAL_FAILED  = "USER_RETRIVAL_FAILED"
	USER_DO_NOT_HAVE_FUND = "USER_DO_NOT_HAVE_FUND"
)

const STATUS_APPROVED = "approved"
const STATUS_REJECTED = "rejected"
const USER_OBJECT_TYPE = "user"

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
	fmt.Println("Registering new user...", args)

	user := User{}
	err := json.Unmarshal([]byte(args), &user)
	if err != nil {
		fmt.Errorf("Failed to parse user register input. %s", err.Error())
		return fmt.Errorf(USER_CREATE_FAILED)
	}

	//Check if user already exist
	_, err = s.GetUser(ctx, user.ID)
	if err != nil && err.Error() != USER_NOT_FOUND {
		return err
	}

	//Upating funds empty
	user.Funds = []string{}

	userAsBytes, err := json.Marshal(user)
	if err != nil {
		fmt.Errorf("Failed to serialize user data. %s", err.Error())
		return fmt.Errorf(USER_CREATE_FAILED)
	}

	key, _ := ctx.GetStub().CreateCompositeKey(USER_OBJECT_TYPE, []string{user.ID})
	fmt.Println("key::", key)
	err = ctx.GetStub().PutState(key, userAsBytes)
	if err != nil {
		fmt.Errorf("Failed to record user data. %s", err.Error())
		return fmt.Errorf(USER_CREATE_FAILED)
	}

	return nil
}

// Approve user by user id
func (s *SmartContract) ApproveUser(ctx contractapi.TransactionContextInterface, userId string) error {
	fmt.Println("ApproveUser...")

	//Checking access permission
	err := s.checkAccessPermission(ctx, "ApproveUser")
	if err != nil {
		fmt.Errorf("Authorization failed", err.Error())
		return err
	}

	//Activating user
	return s.updateUserStatus(ctx, userId, STATUS_APPROVED)
}

func (s *SmartContract) RejectUser(ctx contractapi.TransactionContextInterface, userId string) error {
	fmt.Println("RejectUser...", userId)

	//Checking access permission
	err := s.checkAccessPermission(ctx, "RejectUser")
	if err != nil {
		fmt.Errorf("Authorization failed", err.Error())
		return err
	}

	//Reject user
	return s.updateUserStatus(ctx, userId, STATUS_REJECTED)
}

// update user by user id
func (s *SmartContract) updateUserStatus(ctx contractapi.TransactionContextInterface, userId string, status string) error {
	fmt.Println("UpdateUserStatus...", status, userId)

	//Fetching user
	user, err := s.GetUser(ctx, userId)
	if err != nil {
		fmt.Errorf("Failed to fetch user data. %s", err.Error())
		return err
	}

	//Activating user
	user.Status = status
	userAsBytes, err := json.Marshal(user)
	if err != nil {
		fmt.Errorf("Failed to serialize user data. %s", err.Error())
		return fmt.Errorf(USER_UPDATE_FAILED)
	}

	err = ctx.GetStub().PutState(user.ID, userAsBytes)
	if err != nil {
		fmt.Errorf("Failed to record user data. %s", err.Error())
		return fmt.Errorf(USER_UPDATE_FAILED)
	}

	return nil
}

// GetUser returns the user stored in the world state with id
func (s *SmartContract) GetUser(ctx contractapi.TransactionContextInterface, userId string) (*User, error) {
	fmt.Println("GetUser called..", userId)

	//Checking access permission
	err := s.checkAccessPermission(ctx, "GetUser")
	if err != nil {
		fmt.Errorf("Authorization failed", err.Error())
		return nil, err
	}

	key, _ := ctx.GetStub().CreateCompositeKey(USER_OBJECT_TYPE, []string{userId})
	fmt.Println("key::", key)

	userAsBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		fmt.Errorf("Failed to read user from world state. %s", err.Error())
		return nil, fmt.Errorf(USER_RETRIVAL_FAILED)
	}

	if userAsBytes == nil {
		fmt.Errorf("user %s does not exist", userId)
		return nil, fmt.Errorf(USER_NOT_FOUND)
	}

	user := new(User)
	_ = json.Unmarshal(userAsBytes, user)

	fmt.Println("user:", user)
	return user, nil
}

// QueryAllUsers returns all user
func (s *SmartContract) QueryAllUsers(ctx contractapi.TransactionContextInterface) ([]*User, error) {
	fmt.Println("USER_OBJECT_TYPE, []string{}:", USER_OBJECT_TYPE, []string{})

	//Checking access permission
	err := s.checkAccessPermission(ctx, "QueryAllUsers")
	if err != nil {
		fmt.Errorf("Authorization failed", err.Error())
		return nil, err
	}

	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(USER_OBJECT_TYPE, []string{})
	if err != nil {
		return nil, fmt.Errorf(SERVER_ERROR)
	}
	defer resultsIterator.Close()

	results := []*User{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf(SERVER_ERROR)
		}

		user := new(User)
		_ = json.Unmarshal(queryResponse.Value, user)

		//queryResult := UserQueryResult{Key: queryResponse.Key, Record: user}
		results = append(results, user)
	}

	fmt.Println("results:", results)
	return results, nil
}

// buy existing fund
func (s *SmartContract) BuyFund(ctx contractapi.TransactionContextInterface, fundId string, userId string) error {
	fmt.Println("Buying fund...", fundId, userId)

	//Checking access permission
	err := s.checkAccessPermission(ctx, "BuyFund")
	if err != nil {
		fmt.Errorf("Authorization failed", err.Error())
		return err
	}

	//Check if fund exists
	_, err = s.GetFund(ctx, fundId)
	if err != nil {
		fmt.Errorf("Failed to retrieve Fund with id: %s. Error: %s", fundId, err.Error())
		return err
	}

	// Retrieving user
	user, err := s.GetUser(ctx, userId)
	if err != nil {
		fmt.Errorf("Failed to retrieve user with id: %s. Error: %s", userId, err.Error())
		return err
	}

	//Updating user funds list with new fund id
	user.Funds = append(user.Funds, fundId)

	// Serializing user object
	userAsBytes, err := json.Marshal(user)
	if err != nil {
		fmt.Errorf("Failed to serialize user data. %s", err.Error())
		return fmt.Errorf(BUY_FUND_FAILED)
	}

	key, _ := ctx.GetStub().CreateCompositeKey(USER_OBJECT_TYPE, []string{userId})
	fmt.Println("key::", key)
	err = ctx.GetStub().PutState(key, userAsBytes)
	if err != nil {
		fmt.Errorf("Failed to record user data. %s", err.Error())
		return fmt.Errorf(BUY_FUND_FAILED)
	}

	return nil
}

// sell existing fund
func (s *SmartContract) SellFund(ctx contractapi.TransactionContextInterface, fundId string, userId string) error {
	fmt.Println("Selling fund...", fundId, userId)

	//Checking access permission
	err := s.checkAccessPermission(ctx, "SellFund")
	if err != nil {
		fmt.Errorf("Authorization failed", err.Error())
		return err
	}

	//Check if fund exists
	_, err = s.GetFund(ctx, fundId)
	if err != nil {
		fmt.Errorf("Failed to retrieve Fund with id: %s. Error: %s", fundId, err.Error())
		return err
	}

	// Retrieving user
	user, err := s.GetUser(ctx, userId)
	if err != nil {
		fmt.Errorf("Failed to retrieve user with id: %s. Error: %s", userId, err.Error())
		return err
	}

	fundExists := StringInSlice(fundId, user.Funds)
	if !fundExists {
		fmt.Errorf("User %s does not own fund %s. Hence cannot sell.", err.Error())
		return fmt.Errorf(USER_DO_NOT_HAVE_FUND)
	}

	//Deleting fund from fund list
	user.Funds = deleteFromSlice(user.Funds, fundId)

	// Serializing user object
	userAsBytes, err := json.Marshal(user)
	if err != nil {
		fmt.Errorf("Failed to serialize user data. %s", err.Error())
		return fmt.Errorf(FUND_SELL_FAILED)
	}

	key, _ := ctx.GetStub().CreateCompositeKey(USER_OBJECT_TYPE, []string{userId})
	fmt.Println("key::", key)
	err = ctx.GetStub().PutState(key, userAsBytes)
	if err != nil {
		fmt.Errorf("Failed to update user data. %s", err.Error())
		return fmt.Errorf(FUND_SELL_FAILED)
	}

	return nil
}
