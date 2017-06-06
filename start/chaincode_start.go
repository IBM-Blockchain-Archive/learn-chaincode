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
    "strconv"
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

	err := stub.PutState("aValue",[]byte(args[0]))
    err1:= stub.PutState("bvalue",[]byte(args[1]))

	if err != nil{
		return nil,err
	}
	
     if err1 != nil {
		return nil,err1
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	  if function == "init" {
        return t.Init(stub, "init", args)
    }
     if function == "changeValue" {
    	return t.changeValue(stub,args)
    } 
     if function == "write" {
        return t.write(stub, args)
    }
   
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
	}

	if function == "read"{
		return t.read(stub,args)
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query: " + function)
}

//Custom write function
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface,args []string) ([]byte,error) {

	var err error
    var key,value,jsonResp string

	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}   

	key   = args[0]
	value = args[1]

        valAsbytes, err1 := stub.GetState(key)

        fmt.Println(valAsbytes)
        
        if err1 != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
          return nil, errors.New(jsonResp)
       }

        err = stub.PutState(key,[]byte(value))

      if err!= nil {
      	return nil,err
      }

      return nil,nil
}

func (t *SimpleChaincode) changeValue(stub shim.ChaincodeStubInterface,args []string) ([]byte,error) {

	var err,err1 error
	var changeValue,oldValue int
	var jsonResp string

	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}   
       
      //--------get old value of key 

        key := args[0] 

        valAsbytes, err2 := stub.GetState(key)
        //convert byte array to string
        oldValString:= string(valAsbytes[:])
       
       if err2 != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
       }
        
        //convert old value of key from string to int 
         oldValue ,err = strconv.Atoi(oldValString)

          if err!= nil {
      	return nil,err
      }

       //convert new value of key from string to int
        changeValue,err1 = strconv.Atoi(args[1])
        
        if err1!= nil {
      	return nil,err1
      }

      //Compute new value
       
       newValue := strconv.Itoa(oldValue + changeValue)

    //write new value to stub
   err4:= stub.PutState(key,[]byte(newValue))

    if err4!= nil {
      	return nil,err4
      }
 
      return nil,nil
}

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var key, jsonResp string

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