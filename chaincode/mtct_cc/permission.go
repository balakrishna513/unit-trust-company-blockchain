package main

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

var ACL = map[string][]string{
	"admin": []string{
		"CreateAsset",
		"GetAsset",
		"QueryAllAssets",
		"CreateFund",
		"DeleteFund",
		"GetFund",
		"QueryAllFunds",
		"RegisterUser",
		"ApproveUser",
		"RejectUser",
		"GetUser",
		"QueryAllUsers",
	},
	"agent": []string{
		"sellFund",
	},
	"investor": []string{
		"GetFund",
		"QueryAllFunds",
		"buyFund",
		"sellFund",
		"GetAsset",
		"QueryAllAssets",
	},
}

const ACL1 = `{
   "admin":[
      "CreateAsset",
      "GetAsset",
      "QueryAllAssets",
      "CreateFund",
      "DeleteFund",
      "GetFund",
      "QueryAllFunds",
      "RegisterUser",
      "ApproveUser",
      "RejectUser",
      "GetUser",
      "QueryAllUsers"
   ],
   "agent":[
      "sellFund"
   ],
   "investor":[
      "GetFund",
      "QueryAllFunds",
      "buyFund",
      "sellFund",
      "GetAsset",
      "QueryAllAssets"
   ]
}`

func (s *SmartContract) checkAccessPermission(ctx contractapi.TransactionContextInterface, callingFunc string) error {
	//Reading user's role attribute from user context
	var errorMsg = fmt.Errorf("User not authorized.")

	// Reading role attribute from context
	role, ok, err := cid.GetAttributeValue(ctx.GetStub(), "role")
	if err != nil {
		fmt.Errorf("Failed to read role attribute from user context. %s", err.Error())
		return errorMsg
	}

	if !ok {
		//The client identity does not possess the attribute
		fmt.Println("User role attribute not found in the user context.")
		return errorMsg
	}

	var authorizedFunctions = ACL[role]
	isExist := StringInSlice(callingFunc, authorizedFunctions)
	if !isExist {
		fmt.Println("User not authorized.")
		return errorMsg
	}

	return nil
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
