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
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
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

	err := stub.PutState ("laponia_key", []byte(args[0]))
	if err != nil {
		fmt.Println("\nPopulada a laponia_key: %v", err)
		return nil, err
	}
	return nil, nil
}

// **************************************************************************
// Feb 18th 2017
// @itamar-m
// function write as part blockchain tutorial
	func ( t *SimpleChaincode) write ( stub shim.ChaincodeStubInterface, args []string) ([]byte, error)  {
		var key, value string
		var err error

		fmt.Println ("running func write().")

		if len(args) != 2{
			return nil, errors.New("Numero incorreto de argumentos. Esperado 2 argumentos, nome e valor. Ao menos foi o Daniel Skeff que disse..")
		}

		key = args[0]
		value = args[1]

		err = stub.PutState(key, []byte (value))

		if err != nil {
			fmt.Println("A func write() deu certo.")
			return nil, err
			}
		return nil, nil
	}

// **************************************************************************

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "write" {
			fmt.Println("\nentrou no write!\n")
			return t.write(stub, args)
		}
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}


// **************************************************************************
// Feb 18th 2017 - minor change
// @itamar-m
// function write as part blockchain tutorial

func (t *SimpleChaincode) read (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var key, jsonResp string
	var err error

	if len (args) != 1 {
		return nil, errors.New("Numero incorreto de argumentos. Esperado 1 argumento, nome e valor. Ao menos foi o Daniel Skeff que disse..")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)

	if err != nil {
	  jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}
	return valAsbytes, nil
}

// **************************************************************************

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	// if function == "dummy_query" {											//read a variable
	// 	fmt.Println("hi there " + function)						//error
	//	return nil, nil;
	// }

	if function == "read" {
		fmt.Println("chamando a func read().")
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query: " + function)
}
