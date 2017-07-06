package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var loginPrefix = "LOGIN"

type SimpleChaincode struct {

}

type user struct {
	LoginID 		string 	`json:"user_id"`
	UserType		string 	`json:"user_type"`
	CompanyName 	string	`json:"company_name"`
	CompanyLocation	string	`json:"company_location"`
	BankAccountNum		int		`json:"bank_account_num"`
	BankBalance		float64	`json:"bank_balance"`
}

type userLogin struct {
	LoginName		string 	`json:"login_name"`
	Password 		string	`json:"password"`
}


func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
	
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var testUserName, testUserType, testCompName, testCompLoc, testPassword string
	var testBankAccountNum int
	var testBankBalance float64

	var testUser user
	var testUserLogin userLogin

	testUserName = "testUser1"
	testPassword = "testUser1"
	testUserType = "Buyer"
	testCompName = "testCompany1"
	testCompLoc = "Vancouver"
	testBankAccountNum = 123
	testBankBalance = 1000

	testUser = user{LoginID: testUserName, UserType: testUserType, CompanyName: testCompName, 
	CompanyLocation: testCompLoc, BankAccountNum: testBankAccountNum, BankBalance: testBankBalance}
	userObjBytes, err := json.Marshal(&testUser)
	if err != nil {
		return nil, err
	}

	err1 := stub.PutState(testUserName, userObjBytes)
	if err1 != nil {
		fmt.Println("Failed to save User Details. UserObj")
	}

	testUserLogin =	userLogin{LoginName: testUserName, Password: testPassword} 
	userObjLoginBytes, err := json.Marshal(&testUserLogin)
	err2 := stub.PutState(loginPrefix + testUserName, userObjLoginBytes)
	if err2 != nil {
		fmt.Println("Failed to save user credentials. UserLoginObj")
	}

	return nil, nil

}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Running Invoke function")

	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "register" {
		return t.register(stub, args)
	}

	fmt.Println("Invoke did not find func:" + function)

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
	} else if function == "getUserCredentials" {
		return t.getUserCredentials(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
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

func (t *SimpleChaincode) register(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var userName, userType, compName, compLoc, password string
	var bankAccountNum int
	var bankBalance float64

	var userObj user	
	var userLoginObj userLogin
	fmt.Println("Running function Register")

	if len (args) != 7 {
		return nil, errors.New("Incorrect number of argumets. Expecting 7")
	}
	
	userName = args[0]
	userType = args[1]
	compName = args[2]
	compLoc = args[3]
	bankAccountNum,_ = strconv.Atoi(args[4])
	bankBalance,_ = strconv.ParseFloat(args[5], 64)
	password = args[6]

	userObj = user{LoginID: userName, UserType: userType, 
	CompanyName: compName, CompanyLocation: compLoc, BankAccountNum: bankAccountNum, 
	BankBalance: bankBalance}
	userObjBytes, err := json.Marshal(&userObj)
	if err != nil {
		fmt.Println("Failed to save user credentials. UserObj")
	}
	err3 := stub.PutState(userName, userObjBytes)
	if err3 != nil {
		return nil, errors.New("Failed to save User credentials")
	}

	userLoginObj = userLogin{LoginName: userName, Password: password}
	userLoginBytes, err1 := json.Marshal(&userLoginObj)
	if err1 != nil {
		fmt.Println("Failed to save user credentials. UserObj")
	}

	err2 := stub.PutState(loginPrefix + userName, userLoginBytes)
	if err2 != nil {
		fmt.Println("Failed to save user credentials. UserLoginObj")
	}
	return nil, nil

}

func (t *SimpleChaincode) getUserCredentials(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var userNameGuess, returnMessage string
	var userSample user
	fmt.Println("Getting User Credentials")
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2.")
	}

	userNameGuess = args[0]
	//passwordGuess = args[1]
	
	verifyBytes, err3 := t.verifyUser(stub, args)
	if err3 != nil {
		return nil, err3
	}
	if testEqualSlice(verifyBytes, []byte("Valid")) {
		userInfo, err := stub.GetState(userNameGuess)
		if err != nil {
			return nil, errors.New("User was not properly registered")
		}
		err1 := json.Unmarshal(userInfo, &userSample)
		if err1 != nil {
			return nil, err1
		}
		//more can be added
		returnMessage = "Retrieved Credentials are " + userSample.LoginID + " " + userSample.UserType
		return []byte(returnMessage), nil
	} else {
		returnMessage = "Not authorized to get access"
		return []byte(returnMessage), nil
	}
	return nil, nil

}

func (t *SimpleChaincode) verifyUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var userNameGuess, returnMessage, passwordGuess string
	var loginObj userLogin

	fmt.Println("Verifying User")
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2.")
	}

	userNameGuess = args[0]
	passwordGuess = args[1]

	userLoginInfo, err := stub.GetState(loginPrefix + userNameGuess)
	if userLoginInfo == nil {
		returnMessage = "Invalid Username"
		return []byte(returnMessage), nil
	}

	err1 := json.Unmarshal(userLoginInfo, &loginObj)
	if err1 != nil {
		return nil, err
	}

	if passwordGuess == loginObj.Password {
		returnMessage = "Valid"
		return []byte(returnMessage), nil
	} else {
		returnMessage = "Invalid Password"
		return []byte(returnMessage), nil
	}
	return nil, nil
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
