/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements. ...



*/

package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"time"
	//"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var bitgramIndexStr = "_bitgramindex"			//name for the key/value that will store a list of all known bitgrams
var bitgramMedClaimsStr = "_allclaims"				//name for the key/value that will store all open trades
var bitgramCountStr = "_bitgramcount"			//name for the key/value that will store a list of all known bitgrams
var bitgramMedClaimsCountStr = "_claimscount"				//name for the key/value that will store all open trades



// Name, Address, FinacialScore, GovScore, SocialScore, EmployerScore, QualificationScore, KYCScore, Namespace(the first transacting bank)


/* order of arguments
'IN_BITGRAM_ID'
'IN_SMART_ID_NAMESPACE'
'IN_SMART_ID_NAME'
'IN_SMART_ID_ADDR'
'IN_SMART_ID_FINSC'
'IN_SMART_ID_GOVSC'
'IN_SMART_ID_SOCSC'
'IN_SMART_ID_EMPSC'
'IN_SMART_ID_PUBLICSC'
*/

type Bitgram struct{
	Name string `json:"name"`					   //the fieldtags are needed to keep case from bouncing around
	Address string `json:"address"`
	FS string `json:"fs"`
    GS string `json:"gs"`
	ES string `json:"es"`
	QS string `json:"qs"`
	KS string `json:"ks"`
	NS string `json:"ns"`
}


type Claim struct{
	Bitgram string `json:"bitgram"`					//user who created the open trade order
	Timestamp int64 `json:"timestamp"`			//utc timestamp of creation
	ToHospital string `json:"toHospital"`		               //bought by
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

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	//var Aval int
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(bitgramIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	
	var empty1 []string
	jsonAsBytes, _ = json.Marshal(empty1)								//clear the open trade struct
	err = stub.PutState(bitgramMedClaimsStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	
	err = stub.PutState(bitgramMedClaimsCountStr, []byte("0"))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(bitgramCountStr, []byte("0"))
	if err != nil {
		return nil, err
	}	


	
	return nil, nil
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "medClaim" {			    //updates the bitgram identity record to the chaincode state
		return t.medClaim(stub, args)
	} else if function == "createIdentity" {							//shares a new trade order
		return t.createIdentity(stub, args)
	} 
	fmt.Println("invoke did not find func: " + function)		   //error

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {													//read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query")
}

// ============================================================================================================================
// Read - read a variable from chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name)									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil													//send it onward
}
// ============================================================================================================================
// Write - write variable into chaincode state
// ============================================================================================================================

func (t *SimpleChaincode) Write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, value string // Entities
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
	}

	name = args[0]															//rename for funsies
	value = args[1]
	err = stub.PutState(name, []byte(value))								//write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ============================================================================================================================
// writeBitgramIdentity - create/update a new bitgram identity, store into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) medClaim(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	
    //str := `{ "BITGRAM_ID":"` + args[0] + `","HOSPITAL_ID":"`+ args[1] + `","PAYLOAD":"`+ args[2] + `","TIMESTAMP":"`+ strconv.FormatInt(makeTimestamp(), 10) +`"}`

	str := args[0];
	var medClaimId string = makeNewBitgramId()

		
	err = stub.PutState(medClaimId, []byte(str))								//store bitgram with id as key
	if err != nil {
		return nil, err
	}
		
	//get the bitgram index
	bitgramsAsBytes, err := stub.GetState(bitgramMedClaimsStr)
	if err != nil {
		return nil, errors.New("Failed to get bitgram share index")
	}
	var bitgramShareIndex []string
	json.Unmarshal(bitgramsAsBytes, &bitgramShareIndex)							//un stringify it aka JSON.parse()
	
	//append "[\"111\\\"_\\\"IDFC\"]"
	bitgramShareIndex = append(bitgramShareIndex, medClaimId)								//add bitgram name to index list
	fmt.Println("! bitgram Share index: ", bitgramShareIndex)
	jsonAsBytes, _ := json.Marshal(bitgramShareIndex)
	err = stub.PutState(bitgramMedClaimsStr, jsonAsBytes)						//store name of bitgram


	//update the number of claims
	noOfIdsAsbytes, err := stub.GetState(bitgramMedClaimsCountStr)
	
	var currentNumber string = string(noOfIdsAsbytes)
	
	i, err := strconv.Atoi(currentNumber)
	
	i = i + 1;
	
	s := strconv.Itoa(i)

	err = stub.PutState(bitgramMedClaimsCountStr, []byte(s))								//store bitgram with id as key
	if err != nil {
		return nil, err
	}
	
	fmt.Println("- end share bitgram")
	return nil, nil
}

// ============================================================================================================================
// writeBitgramIdentity - create/update a new bitgram identity, store into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) createIdentity(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	var superIdentity string = makeNewBitgramId()
	var encryptedIdentityData string = args[0]

	err = stub.PutState(superIdentity, []byte(encryptedIdentityData))								//store bitgram with id as key
	if err != nil {
		return nil, err
	}
		
	//get the bitgram index
	bitgramsAsBytes, err := stub.GetState(bitgramIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get bitgram index")
	}
	var bitgramIndex []string
	json.Unmarshal(bitgramsAsBytes, &bitgramIndex)							//un stringify it aka JSON.parse()
	
	//append
	bitgramIndex = append(bitgramIndex, superIdentity)								//add bitgram name to index list
	fmt.Println("! bitgram index: ", bitgramIndex)
	jsonAsBytes, _ := json.Marshal(bitgramIndex)
	err = stub.PutState(bitgramIndexStr, jsonAsBytes)						//store name of bitgram

	//update the number of identities
	noOfIdsAsbytes, err := stub.GetState(bitgramCountStr)
	
	var currentNumber string = string(noOfIdsAsbytes)
	
	i, err := strconv.Atoi(currentNumber)
	
	i = i + 1;
	
	s := strconv.Itoa(i)

	err = stub.PutState(bitgramCountStr, []byte(s))								//store bitgram with id as key
	if err != nil {
		return nil, err
	}


	fmt.Println("- end init bitgram")
	return []byte(superIdentity), nil
}


// ============================================================================================================================
// Make Timestamp - create a timestamp in ms
// ============================================================================================================================
func makeTimestamp() int64 {
    return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}

func makeNewBitgramId() string {
    return strconv.FormatInt(time.Now().Unix(), 10)
}

