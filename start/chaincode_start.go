/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var accountPrefix = "acct:"
var securityToken = "D44867B6ADB93F15D3DD77C323BF6"

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// comment
type Account struct {
	ID          string  `json:"id"`
	CashBalance float64 `json:"cashBalance"`
}

type Transaction struct {
	FromId string `json:"fromId"`
	ToId   string `json:"toId"`
	amount int    `json:"amount"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Initialize the collection of commercial paper keys
	fmt.Println("Initializing accountIds collection")
	var blank map[string]int
	blankBytes, _ := json.Marshal(&blank)
	err := stub.PutState("AccountIds", blankBytes)
	if err != nil {
		fmt.Println("Failed to initialize paper key collection")
	}

	fmt.Println("Initialization complete")
	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "transfer" {
		return t.transfer(stub, args)
	} else if function == "registerAccounts" {
		return t.registerAccounts(stub, args)
	}
	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *SimpleChaincode) registerAccounts(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Creating accounts")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting account numbers")
	}

	//var err error
	var ids []string
	err := json.Unmarshal([]byte(args[0]), &ids)
	if err != nil {
		fmt.Println("error creating accounts with input")
		return nil, errors.New("registerAccounts accepts an array of account ids")
	}

	registeredIds, err := GetAllAccountIds(stub)
	var newIds []string

	//create a bunch of accounts
	for _, id := range ids {
		if registeredIds[id] != 1 {
			newIds = append(newIds, id)
		}
		registeredIds[id] = 1
	}

	registeredIdsBytes, err := json.Marshal(&registeredIds)
	if err != nil {
		fmt.Println("error marshaling accounts")
		return nil, errors.New("Error marshaling accounts")
	}

	err = stub.PutState("AccountIds", registeredIdsBytes)
	if err != nil {
		fmt.Println("error putting accounts")
		return nil, errors.New("Error putting accounts")
	}

	var balance = 10000

	for _, newId := range newIds {
		balanceBytes, err := json.Marshal(&balance)
		err = stub.PutState(accountPrefix+newId, balanceBytes)
		if err != nil {
			fmt.Println("error putting account balance")
			return nil, errors.New("Error putting account balance")
		}
	}

	fmt.Println("Accounts created")
	return nil, nil

}

func (t *SimpleChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Handling a transfer")

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting from, to, amount and token")
	}

	if args[3] != securityToken {
		fmt.Println("security token does not match")
		return nil, errors.New("security token does not match")
	}

	registeredIds, err := GetAllAccountIds(stub)
	if err != nil {
		fmt.Println("error gettng account ids")
		return nil, errors.New("Error gettng account ids")
	}

	var fromId = args[0]
	var toId = args[1]

	amount, err := strconv.ParseInt(args[2], 10, 0)
	if err != nil {
		fmt.Println("Amount is not a number")
		return nil, errors.New("Amount is not a number")
	}

	if registeredIds[fromId] != 1 {
		fmt.Println("error: from account is not registered")
		return nil, errors.New("Erro: from account is not registered")
	}

	if registeredIds[toId] != 1 {
		fmt.Println("error: to account is not registered")
		return nil, errors.New("Error: to account is not registered")
	}

	fromBalance, err := GetAccountBalance(stub, fromId)
	if err != nil {
		fmt.Println("error gettng from account balance")
		return nil, errors.New("Error getting from account balance")
	}

	var fromBalance64 = int64(fromBalance)

	if fromBalance64 < amount {
		fmt.Println("error not enough resources on from account")
		return nil, errors.New("error not enough resources on from account")
	}

	toBalance, err := GetAccountBalance(stub, toId)
	if err != nil {
		fmt.Println("error gettng from account balance")
		return nil, errors.New("Error getting from account balance")
	}
	var toBalance64 = int64(toBalance)

	fromBalance64 = fromBalance64 - amount
	toBalance64 = toBalance64 + amount

	fromBalanceBytes, err := json.Marshal(&fromBalance64)
	err = stub.PutState(accountPrefix+fromId, fromBalanceBytes)
	if err != nil {
		fmt.Println("error putting from account balance")
		return nil, errors.New("Error putting from account balance")
	}

	toBalanceBytes, err := json.Marshal(&toBalance64)
	err = stub.PutState(accountPrefix+toId, toBalanceBytes)
	if err != nil {
		fmt.Println("error putting to account balance")
		return nil, errors.New("Error putting to account balance")
	}

	return nil, nil
}

func GetAllAccountIds(stub shim.ChaincodeStubInterface) (map[string]int, error) {

	var accountIds map[string]int

	// Get list of all the keys
	idsBytes, err := stub.GetState("AccountIds")
	if err != nil {
		fmt.Println("Error retrieving account Ids")
		return nil, errors.New("Error retrieving account Ids")
	}

	err = json.Unmarshal(idsBytes, &accountIds)
	if err != nil {
		fmt.Println("Error unmarshalling account Ids")
		return nil, errors.New("Error unmarshalling account Ids")
	}

	return accountIds, nil
}

func GetAccountBalance(stub shim.ChaincodeStubInterface, id string) (int, error) {

	balanceBytes, err := stub.GetState(accountPrefix + id)
	if err != nil {
		fmt.Println("Error retrieving account Ids")
		return 0, errors.New("Error retrieving account Ids")
	}
	var balance int
	err = json.Unmarshal(balanceBytes, &balance)
	if err != nil {
		fmt.Println("Error unmarshalling account Ids")
		return 0, errors.New("Error unmarshalling account Ids")
	}

	return balance, nil
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "GetAllAccountIds" { //read a variable
		fmt.Println("Getting all accounts")
		registeredIds, err := GetAllAccountIds(stub)
		if err != nil {
			fmt.Println("error gettng account ids")
			return nil, err
		}

		registeredIdsBytes, err1 := json.Marshal(&registeredIds)
		if err1 != nil {
			fmt.Println("Error marshalling registeredIds")
			return nil, err1
		}

		fmt.Println("All success, returning accounts")
		return registeredIdsBytes, nil
	} else if function == "GetAccountBalance" { //read a variable
		fmt.Println("Getting account balance")

		balance, err := GetAccountBalance(stub, args[0])
		if err != nil {
			fmt.Println("error gettng account balance")
			return nil, err
		}

		balanceBytes, err1 := json.Marshal(&balance)
		if err1 != nil {
			fmt.Println("Error marshalling balanceBytes")
			return nil, err1
		}

		fmt.Println("All success, returning account balance")
		return balanceBytes, nil
	}

	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query: " + function)
}
