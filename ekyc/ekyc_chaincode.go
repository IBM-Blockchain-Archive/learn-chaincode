/*
ekyc Smart Contract (Chain Code) is intended to
manage transactions on EKYC Block chain. It implements the
chaincode shim interface and contains logic for the three
core functions - Init, Invoke and Query.
*/
package main

//Import Dependencies.
import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode struct declaration.
type SimpleChaincode struct {
}

var financialIndexStr = "_allFinancialInstitutes" //name for the key/value that will store a list of all ADDED Financial Institutes

// Ekyc struct - Holds attributes for storing EKYC block.
type Ekyc struct {
	UUID      string `json:"uuid"`      //The Unique Identifier (Aadhar number of the Customer
	Timestamp int64  `json:"timestamp"` //utc timestamp of creation
	AddedBy   string `json:"addedBy"`   //Financial Institute which added this Block.
}

//FinancialInstitute struct - Holds attributes for storing Financial Institue block.
type FinancialInstitute struct {
	Name    string `json:"name"`    //Name of the Financial Institute.
	Address string `json:"address"` //Address of the Financial Institute.
}

// ============================================================================================================================
// main function is invoked when the instance of Chain Code is deployed.
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode ::: %s", err)
	}
}

// Init function is intended to intialize the Chain Code.
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Running " + function)
	var err error
	//Do we need to cross check if base entry is already present?
	listAsBytes, _ := json.Marshal("") //marshal an emtpy string to clear the index
	err = stub.PutState(financialIndexStr, listAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke - Our entry point for Invocations. This is primarily responsible for
//writing to the Ledger.
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "writeKYC" { //writes a KYC block to the Blockchain.
		return t.WriteKYC(stub, args)
	} else if function == "writePeer" { //writes a block entry encapsulating Financial Institution.
		return t.WritePeer(stub, args)
	}
	fmt.Println("Invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation")
}

// Query - Our entry point for Queries. This is primarily responsible for
//querying chaincode state.
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Running " + function)

	// Handle different functions
	if function == "readKYC" {
		return t.readKYC(stub, args) //read a KYC
	} else if function == "readPeer" {
		return t.readPeer(stub, args) //read a Financial Institute details
	} else if function == "readAllPeers" {
		return t.readAllPeers(stub, args) //read list of registered Financial Institutes
	}
	fmt.Println("Query did not find func: " + function) //error

	return nil, errors.New("Received unknown function " + function)
}

// readKYC - reads a KYC Block from Block Chain.
func (t *SimpleChaincode) readKYC(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string

	//There will be one arg passed.
	//Arg 1 - The UUID.

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting UUID")
	}

	valAsbytes, err := stub.GetState(args[0]) //get the value from chaincode state

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}
	return valAsbytes, nil

}

// readPeer - reads a Financial Institute Peer from Block Chain.
func (t *SimpleChaincode) readPeer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	//There will be one arg passed.
	//Arg 1 - The Name of the Financial Institute.

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting Financial Institute Name to be queried")
	}

	valAsbytes, err := stub.GetState(args[0]) //get the var from chaincode state

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}
	return valAsbytes, nil

}

// readAllPeers - reads all Financial Institute Names from Block Chain.
func (t *SimpleChaincode) readAllPeers(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string

	valAsbytes, err := stub.GetState(financialIndexStr)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for list of registered Financial Institutes \"}"
		return nil, errors.New(jsonResp)
	}
	return valAsbytes, nil

}

// WriteKYC - write variable into chaincode state
func (t *SimpleChaincode) WriteKYC(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//There will be two args passed.
	//Arg 1 - The UUID of Customer.
	//Arg 2 - The Finacial Institute who approved this UUID.
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	var err error
	kyc := Ekyc{
		UUID:      args[0],
		Timestamp: makeTimestamp(),
		AddedBy:   args[1],
	}
	bytes, err := json.Marshal(kyc)
	if err != nil {
		fmt.Println("Error marshaling KYC")
		return nil, errors.New("Error marshaling KYC")
	}

	err = stub.PutState(kyc.UUID, bytes)

	if err != nil {
		return nil, err
	}
	return nil, nil
}

// WritePeer - write Bank realated data into chaincode state
func (t *SimpleChaincode) WritePeer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//There will be two args passed.
	//Arg 1 - The Name of the Financial Institute.
	//Arg 2 - The Address of the Financial Institute.
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	var err error

	fi := FinancialInstitute{
		Name:    args[0],
		Address: args[1],
	}

	bytes, err := json.Marshal(fi)
	if err != nil {
		fmt.Println("Error marshaling Financial Institute")
		return nil, errors.New("Error marshaling Financial Institute")
	}

	err = stub.PutState(fi.Name, bytes)

	if err != nil {
		return nil, err
	}

	valAsbytes, err := t.readAllPeers(stub, args)
	if err != nil {
		return nil, err
	}
	var valueAll string
	valueAll = string(valAsbytes)

	if valueAll == "" {
		valueAll = args[0]
	} else {
		valueAll = string(valAsbytes) + ";" + args[0]
	}

	err = stub.PutState(financialIndexStr, []byte(valueAll))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ============================================================================================================================
// Make Timestamp - create a timestamp in ms
// ============================================================================================================================
func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}
