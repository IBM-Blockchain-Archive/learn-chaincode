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
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

/*type registeredUser struct {
	Name string
	Type string
	BankBalance float64
	Username string
	Password string
}*/


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

	err := stub.PutState("hello_world", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "enroll" {
		return t.enroll(stub, args)
	} 
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "verifyUser" {
		return t.verifyUser(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) hello(stub shim.ChaincodeStubInterface) ([]byte, error){
	b,err := stub.GetState("hi")
	if b != nil {
		return []byte("b is not nil"), nil
	} else {
		return []byte("b is nil"), nil
	}

	if err != nil {
		return []byte("err is not nil"), nil
	} else {
		return []byte("err is nil"), nil
	}
}

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

// write - invoke function to write key/value pair 
/*func (t *SimpleChaincode) enroll(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) { 
	var name, type, username, password string 
	var bankBalance float64 
	var newUser registeredUser
	fmt.Println("running write()") 
	
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5. name of the key and value to set")
	}

	name = args[0] //rename for fun
	type = args[1]
	bankBalance = args[2]
	username = args[3]
	password = args[4]

	newUser = registeredUser{name, type, bankBalance, username, password}
	jsonUserInfo, err := json.Marshal(newUser)
	if err != nil {

	}



	err := stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}*/

func (t *SimpleChaincode) verifyUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var keyGuess string
	var valGuess string
	var returnMessage string
	var err error
	fmt.Println("running read")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expected 2")
	}

	keyGuess = args[0]
	valGuess = args[1]
	
	valActual, err := stub.GetState(keyGuess) 
	if err != nil { 
		returnMessage = "Username Incorrect. Login Failed"
		return []byte(returnMessage), nil
	}

	if testEqualSlice([]byte(valGuess), valActual) {
		returnMessage = "Login Succesful"	
		return []byte(returnMessage), nil
	} else {
		returnMessage = "Password Incorrect. Login Failed"
		return []byte(returnMessage), nil
	}
}



// read - query function to read key/value pair
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

func testEqualSlice (a []byte, b []byte) bool {

	if a == nil && b == nil { 
        return true; 
    } else if a == nil || b == nil { 
        return false; 
    } 
	
	if len(a) != len(b) {
        return false
    }

    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}

