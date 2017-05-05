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
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
)

var logger = shim.NewLogger("CLDChaincode")

//==============================================================================================================================
//	 Participant types - Each participant type is mapped to an integer which we use to compare to the value stored in a
//						 user's eCert
//==============================================================================================================================
//CURRENT WORKAROUND USES ROLES CHANGE WHEN OWN USERS CAN BE CREATED SO THAT IT READ 1, 2, 3, 4, 5
const   MANUFACTURER      =  "manufacturer"
const   DISTRIBUTOR   =  "distributor"
const   RETAILER =  "retailer"
const   CUSTOMER  =  "customer"

//==============================================================================================================================
//	 Status types - Asset lifecycle is broken down into 5 statuses, this is part of the business logic to determine what can
//					be done to the sparepart at points in it's lifecycle
//==============================================================================================================================
const   STATE_TEMPLATE  			=  0
const   STATE_MANUFACTURE  			=  1
const   STATE_DISTRIBUTOR 			=  2
const   STATE_RETAILER	 			=  3
const   STATE_CUSTOMER  			=  4

//==============================================================================================================================
//	 Structure Definitions
//==============================================================================================================================
//	Chaincode - A blank struct for use with Shim (A HyperLedger included go file used for get/put state
//				and other HyperLedger functions)
//==============================================================================================================================
type SimpleChaincode struct {
}

//==============================================================================================================================
//	SparePart - Defines the structure for a car object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON make -> Struct Make.
//==============================================================================================================================
type SparePart struct {
	PartNumber            string `json:"partnumber"`
	Description           string `json:"description"`
	UIdentifier           string `json:"uid"`
	ManufacturingDate     string `json:"MDate"`
	Owners                string `json:"owners"`
	Remarks               string `json:"remarks"`
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

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	//my code starts
	err := stub.PutState("uid", []byte(args[0]))
    if err != nil {
        return nil, err
    }
	//my code ends
	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	}else if function == "write" {//my code starts
        return t.write(stub, args)
    }else if function == "create_sparepart" {
        return t.create_sparepart(stub, args)
    } //my code ends
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "dummy_query" {											//read a variable
		fmt.Println("hi there " + function)						//error
		return nil, nil;
	}else if function == "read" { //my code starts
        return t.read(stub, args)
    }else if function == "get_sparepart_detail" { 
        return t.get_sparepart_detail(stub, args)
    }
	//my code ends
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query: " + function)
}
//my own code starts here
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0]                            //rename for fun
	value = args[1]
	err = stub.PutState(key, []byte(value))  //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    key = args[0]
    valAsbytes, err := stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }

    return valAsbytes, nil
}

//=================================================================================================================================
//	 Create Function
//=================================================================================================================================
//	 Create Sparepart - Creates the initial JSON for the sparepart and then saves it to the ledger.
//=================================================================================================================================
func (t *SimpleChaincode) create_sparepart(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	//create spare part object																	
	tempObj := &SparePart{
        PartNumber:			args[0],
        Description:		args[1],
		UIdentifier:		args[2],
		ManufacturingDate:	args[3],
		Owners:				args[4],
		Remarks:			args[5]}
    sparePartBytes, err := json.Marshal(tempObj)

	if err != nil { fmt.Print("CREATE_SPAREPART: Error creating Spare part record") }
	
	//assign spare part to ledger
	err = stub.PutState(string(tempObj.UIdentifier), sparePartBytes)

	if err != nil { fmt.Printf("SAVE_CHANGES: Error storing vehicle record: %s", err)}


	return nil, nil
}
//==============================================================================================================================
// get_sparepart_detail - Reads spare part detail from ledger.
//==============================================================================================================================
func (t *SimpleChaincode) get_sparepart_detail(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	var key, jsonResp string
	var err error
	
	if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    key = args[0]
    sparePartAsbytes, err := stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }

    return sparePartAsbytes, nil
}
//my code ends