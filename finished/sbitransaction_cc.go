/*
Invoke Methods :
****************
CreateTransaction
UpdateSupervisorDetails
UpdateL1AuthorizerDetails
UpdateL2AuthorizerDetails
SetDateTime

Query Methods :
****************
GetTransactionInitDetailsForRefAndMaker
GetTransactionInitDetailsForRef
GetAllDetailsForRef_AuditTrial
ListRefnoForDate
ListRefnoForBranch
ListAllTransactions

Dependency Methods :
*********************

GetTransactionInitiationMap
GetSupervisorMap
GetL1AuthMap
GetL2AuthMap

SetTransactionInitiationMap
SetSupervisorMap
SetL1AuthMap
SetL2AuthMap

*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// transaction will implement the processes
type SBITransaction struct {
}

type transactionInitiation struct {
	transRefNo   		string  `json:"ref_no"`
	remAccNo     		string  `json:"rem_accno"`
	remAmtINR    		float32 `json:"rem_amtinr"`
	benAccNo	 		string  `json:"ben_accno"`		
	eventDesc    		string  `json:"event_desc"`
	makerUserID  		string  `json:"maker_id"`
	makerIPAddr  		string  `json:"maker_ipaddr"`
	makerDate    		string  `json:"maker_date"`
	amlStatus    		string  `json:"aml_status"`
	ofacStatus   		string  `json:"ofac_status"`
	rbiStatus   		string  `json:"rbi_status"`
	trans_init_branch	string  `json:"trans_init_branch"`
	maker_branch		string  `json:"maker_branch"`
}

type amlCheck struct {
	transRefNo   string  `json:"ref_no"`
	supUserID    string  `json:"sup_userid"`
	supIPAddr    string  `json:"sup_ipaddr"`
	supDate      string  `json:"sup_date"`
	supStatus	 string  `json:"sup_status"`
}

type l1Auth struct {
	transRefNo  	string  `json:"ref_no"`
	l1UserID     	string  `json:"l1_userid"`
	l1IPAddr     	string  `json:"l1_ipaddr"`
	l1Date       	string  `json:"l1_date"`
	l1Status	 	string  `json:"l1_status"`
}

type l2Auth struct {
	transRefNo   		string  `json:"ref_no"`
	l2UserID    		string  `json:"l2_userid"`
	l2IPAddr     		string  `json:"l2_ipaddr"`
	l2Date       		string  `json:"l2_date"`
	l2Status	 		string  `json:"l2_status"`
	finacleDate 	 	string  `json:"finacle_date"`
	finalcleStatus 	 	string  `json:"finacle_status"`
	TCSBancsDate 		string  `json:"tcs_bancsdate"`
	TCSBancsStatus 		string  `json:"tcs_bancsstatus"`
	PSGDate      		string  `json:"psg_date"`	
	PSGStatus    		string  `json:"psg_status"`
}

type AuditTrial struct {
	trans_init     transactionInitiation `json:"trans_init"`
	aml_check      amlCheck              `json:"aml_check"`
	l1_auth        l1Auth                `json:"l1_auth"`
	l2_auth        l2Auth                `json:"l2_auth"`
}

//Global declaration of maps
var trans_Init_map map[string]transactionInitiation
var supervisor_map map[string]amlCheck
var l1Auth_map map[string]l1Auth
var l2Auth_map map[string]l2Auth
var date_map map[time.Time]string  // key : date and time ; value : ref no array


//Invoke methods starts here 

func CreateTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var trans_obj transactionInitiation
	var err error

	fmt.Println("Entering createTransaction")

	if (len(args) < 1) {
		fmt.Println("Invalid number of args")
		return nil, errors.New("Expected atleast one arguments for initiate Transaction")
	}

	//unmarshal transaction initiation data from UI to "transactionInitiation" struct
	err = json.Unmarshal([]byte(args[1]), &trans_obj)
	if err != nil {
		fmt.Printf("Unable to unmarshal createTransaction input transaction initiation : %s\n", err)
		return nil, nil
	}

	// saving transactionInitiation and maker into map
	GetTransactionInitiationMap(stub)

	//put transaction initiation data and maker data into map
	trans_Init_map[trans_obj.transRefNo] = trans_obj

	SetTransactionInitiationMap(stub)
	//SetDateTime(stub,trans_obj.transRefNo,trans_obj.makerDate)
	
	fmt.Printf("transaction initiation map : %v \n", trans_Init_map)	
	fmt.Println("Transaction initiation Successfully saved")	
	
	return nil, nil
}

func UpdateSupervisorDetails(stub shim.ChaincodeStubInterface, args1 []string) error {
	
	var sup_obj amlCheck	
	var err error

	fmt.Println("Entering UpdateSupervisorDetails")

	if (len(args1) < 1) {
		fmt.Println("Invalid number of args")
		return errors.New("Expected atleast one arguments for UpdateSupervisor")
	}

	//unmarshal supervisor data from UI to "amlCheck" struct
	err = json.Unmarshal([]byte(args1[1]), &sup_obj)
	if err != nil {
		fmt.Printf("Unable to marshal createTransaction input UpdateSupervisor : %s\n", err)
		return nil
	}

	// saving transactionInitiation and maker into map
	GetSupervisorMap(stub)

	//put supervisor data data into map
	supervisor_map[sup_obj.transRefNo] = sup_obj
	
	SetSupervisorMap(stub)
	
	fmt.Printf("supervisor map : %v \n", supervisor_map)
	fmt.Println("supervisor details Successfully updated")	
	
	return nil
}

func UpdateL1AuthorizerDetails(stub shim.ChaincodeStubInterface, args1 []string) error {
	
	var l1_obj l1Auth	
	var err error

	fmt.Println("Entering UpdateSupervisor")

	if (len(args1) < 1) {
		fmt.Println("Invalid number of args")
		return errors.New("Expected atleast one arguments for UpdateL1AuthorizerDetails")
	}

	//unmarshal l1Authorizer data from UI to "l1Auth" struct
	err = json.Unmarshal([]byte(args1[1]), &l1_obj)
	if err != nil {
		fmt.Printf("Unable to marshal  createTransaction input UpdateL1AuthorizerDetails : %s\n", err)
		return nil
	}

	// saving l1Authorizer details into map
	GetL1AuthMap(stub)

	//put l1Authorizer data into map
	l1Auth_map[l1_obj.transRefNo] = l1_obj
	
	SetL1AuthMap(stub)
	
	fmt.Printf("L1Authorizer map : %v \n", supervisor_map)	
	fmt.Println("L1Authorizer details Successfully updated")	
	
	return nil
}

func UpdateL2AuthorizerDetails(stub shim.ChaincodeStubInterface, args1 []string) error {

	var l2Auth_obj l2Auth
	var err error

	fmt.Println("Entering UpdateL2AuthorizerDetails")

	if 	(len(args1) < 1) {
		fmt.Println("Invalid number of args")
		return errors.New("Expected atleast one arguments for UpdateL2AuthorizerDetails")
	}

	//unmarshal L2Auth data from UI to "l2Auth" struct
	err = json.Unmarshal([]byte(args1[1]), &l2Auth_obj)
	if err != nil {
		fmt.Printf("Unable to marshal the input from UpdateL2AuthorizerDetails : %s\n", err)
		return nil
	}

	// saving L2Auth data and system processed data into map
	GetL2AuthMap(stub)

	//put supervisor data and system processed data into map
	l2Auth_map[l2Auth_obj.transRefNo] = l2Auth_obj
	
	SetL2AuthMap(stub)
	
	fmt.Printf("L2Auth map : %v \n", l2Auth_map)			
	fmt.Println("L2AuthorizerDetails Successfully updated")	
	
	return nil
}

/*func SetDateTime(stub shim.ChaincodeStubInterface, refNo string, trans_date string) error {
		var err error
		var bytesread []byte
		var transDate1 time.Time 

		fmt.Printf("setDateTime\n")

		bytesread, err = stub.GetState("DateMap")
		if err != nil {
		fmt.Printf("Failed to get  DateMap for block chain :%v\n", err)
		return err
		}
		if len(bytesread) != 0 {
		fmt.Printf("DateMap map exists.\n")
		err = json.Unmarshal(bytesread, &date_map)
		if err != nil {
			fmt.Printf("Failed to initialize  dateMap for block chain :%v\n", err)
			return err
			}
		} else {
		date_map = make(map[time.Time] string)

		//logic to type conversion of date time type
		var layout = "2006-01-02 15:04:05 -0700 IST"
		transDate1, err = time.Parse(layout, trans_date)
		fmt.Println(transDate1, err)

		if err != nil {
			fmt.Printf("Failed to convert input into date time format :%v\n", err)
			return err
		} else {
			date_map[transDate1]=refNo
		}		

		bytesread, err = json.Marshal(&date_map)
		if err != nil {
			fmt.Printf("Failed to initialize DateTime for block chain :%v\n", err)
			return err
		}
		err = stub.PutState("DateMap", bytesread)
		if err != nil {
			fmt.Printf("Failed to initialize  DateTime for block chain :%v\n", err)
			return err
		}
	}
	return nil
}*/


func UpdateEventDesc(stub shim.ChaincodeStubInterface, args []string) error {
		var err error
		var bytesread []byte		
		var object transactionInitiation

		var refNo string
		var event_desc string

		fmt.Printf("updateEventDesc\n")

		if 	(len(args) < 2) {
		fmt.Println("Invalid number of args")
		return errors.New("Expected atleast one arguments for updateEventDesc" + args[0])
		}
		refNo=args[0]
		event_desc=args[1]

		fmt.Printf("Ref NO sent :%v\n", refNo)
		fmt.Printf("Event Desc sent :%v\n", event_desc)

		bytesread, err = stub.GetState("TransactionInitiationMap")
		if err != nil {
		fmt.Printf("Failed to get  TransactionInitiationMap for block chain :%v\n", err)
		return err
		}
		if len(bytesread) != 0 {

		fmt.Printf("TransactionInitiationMap map exists.\n")
		//Fetch existing values from blockchain
		err = json.Unmarshal(bytesread, &trans_Init_map)

		if err != nil {
			fmt.Printf("Failed to initialize TransactionInitiationMap for block chain :%v\n", err)
			return err
		}

		object=trans_Init_map[refNo]
		//update event desc 
		object.eventDesc=event_desc

		//update new struct values
		trans_Init_map[refNo]=object

		//update new values in blockchain
		bytesread, err = json.Marshal(&trans_Init_map)

		if err != nil {
			fmt.Printf("Failed to initialize TransactionInitiationMap for block chain :%v\n", err)
			return err
		}
		err = stub.PutState("TransactionInitiationMap", bytesread)
		if err != nil {
			fmt.Printf("Failed to initialize  TransactionInitiationMap for block chain :%v\n", err)
			return err
			}

		}	
		return nil	
}

//Invoke methods ends here 


//Query methods starts here 
func ListRefnoForBranch(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var object transactionInitiation 
	var bytes []byte
	var err error

	fmt.Println("Entering ListRefnoForBranch")

	if (len(args) < 1) {
		fmt.Println("Invalid number of arguments")
		return nil, errors.New("Missing Ref no")
	}

	// Getting values from blockchain
	GetTransactionInitiationMap(stub)

	fmt.Printf("Entering GetTransactionInitiation : %v\n", args[0])
	//var refNo = args[1]

	for _, value := range trans_Init_map {
			if value.trans_init_branch  == args[0] {				
					object = value				
			}					
	}

	bytes, err = json.Marshal(&object)
	if err != nil {
		fmt.Printf("Unable to marshal the object array %s\n", err)
		return nil, err
	}

	fmt.Printf(" Transaction initiation details  for particular ref no : %v\n", bytes)
	return bytes, nil
}



func GetTransactionInitDetailsForRefAndMaker(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var object transactionInitiation 
	var bytes []byte
	var err error

	fmt.Println("Entering getTransactionInitDetailsForRefstub")

	if (len(args) < 1) {
		fmt.Println("Invalid number of arguments")
		return nil, errors.New("Missing Ref no")
	}

	// Getting values from blockchain
	GetTransactionInitiationMap(stub)

	fmt.Printf("Entering GetTransactionInitiation : %v\n", args[1])
	var refNo = args[1]

	for _, value := range trans_Init_map {
		if value.transRefNo  == refNo {
			if value.trans_init_branch  == value.maker_branch {				
					object = value				
			}			
		}
	}

	bytes, err = json.Marshal(&object)
	if err != nil {
		fmt.Printf("Unable to marshal the object array %s\n", err)
		return nil, err
	}

	fmt.Printf(" Transaction initiation details  for particular ref no : %v\n", bytes)
	return bytes, nil
}


//Query methods starts here 
func GetTransactionInitDetailsForRef(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var object transactionInitiation 
	var bytes []byte
	var err error

	fmt.Println("Entering getTransactionInitDetailsForRefstub")

	if (len(args) < 1) {
		fmt.Println("Invalid number of arguments")
		return nil, errors.New("Missing Ref no")
	}

	// Getting values from blockchain
	GetTransactionInitiationMap(stub)

	fmt.Printf("Entering GetTransactionInitiation : %v\n", args[1])
	var refNo = args[1]

	for _, value := range trans_Init_map {
		if value.transRefNo  == refNo {							
					object = value								
		}
	}

	bytes, err = json.Marshal(&object)
	if err != nil {
		fmt.Printf("Unable to marshal the object array %s\n", err)
		return nil, err
	}

	fmt.Printf(" Transaction initiation details  for particular ref no : %v\n", bytes)
	return bytes, nil
}


func GetAllDetailsForRef_AuditTrial(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var object transactionInitiation 
	var object1 amlCheck 
	var object2 l1Auth 
	var object3 l2Auth 
	var bytes []byte
	var fin_object AuditTrial
	var err error

	fmt.Println("Entering getTransactionInitDetailsForRefstub")

	if (len(args) < 1) {
		fmt.Println("Invalid number of arguments")
		return nil, errors.New("Missing Ref no")
	}

	// 1) Transaction Initiation

	fmt.Printf("Entering GetTransactionInitiation : %v\n", args[0])

	fmt.Printf("Transaction initiation details")
	GetTransactionInitiationMap(stub)

	var refNo = args[0]

	for _, value := range trans_Init_map {
		if value.transRefNo  == refNo {
			object = value
		}
	}

	// Return transaction initiation details
	bytes, err = json.Marshal(&object)
	if err != nil {
		fmt.Printf("Unable to marshal the transaction initiation  array %s\n", err)
		return nil, err
	}
	fmt.Printf(" Transaction initiation details  for particular ref no : %v\n", object)
	

	// 2) Supervisor details
	
	fmt.Printf("Superviosor details")
	GetSupervisorMap(stub)

	for _, value := range supervisor_map {
		if value.transRefNo  == refNo {
			object1 = value
		}
	}

	// Return Supervisor details
	bytes, err = json.Marshal(&object1)
	if err != nil {
		fmt.Printf("Unable to marshal the Supervisor array %s\n", err)
		return nil, err
	}

	fmt.Printf(" Transaction initiation details  for particular ref no : %v\n", object)

	
	// 3) L1Auth Details

	fmt.Printf("L1Authorizer details")
	GetL1AuthMap(stub)

	for _, value := range l1Auth_map {
		if value.transRefNo  == refNo {
			object2 = value
		}
	}

		// Return l1Auth details
	bytes, err = json.Marshal(&object2)
	if err != nil {
		fmt.Printf("Unable to marshal the l1Auth array %s\n", err)
		return nil, err
	}

	fmt.Printf(" Transaction initiation details  for particular ref no : %v\n", object)


	// 4) L2Auth Details

	fmt.Printf("L2Authorizer details")
	GetL2AuthMap(stub)

	for _, value := range l2Auth_map {
		if value.transRefNo  == refNo {
			object3 = value
		}
	}

	// Return l2Auth Details
	bytes, err = json.Marshal(&object3)
	if err != nil {
		fmt.Printf("Unable to marshal the l2Auth array %s\n", err)
		return nil, err
	}

	fin_object.trans_init=object
	fin_object.aml_check=object1
	fin_object.l1_auth=object2
	fin_object.l2_auth=object3

	// Return AuditTrial Details
	bytes, err = json.Marshal(&fin_object)
	if err != nil {
		fmt.Printf("Unable to marshal audit trial %s\n", err)
		return nil, err
	}


	fmt.Printf(" audit trial details for particular ref no : %v\n", bytes)

	return bytes, nil
}

func ListRefnoForDate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var bytesRead []byte
	var refNo string
	var t time.Time

	fmt.Println("Entering getTransactionInitDetailsForRefstub")

	if (len(args) < 1) {
		fmt.Println("Invalid number of arguments")
		return nil, errors.New("Missing Ref no")
	}

	var trans_date = args[1]
	bytesRead, err = stub.GetState("DateMap")

	if err != nil {
		fmt.Printf("Failed to get  SupervisorMap for block chain :%v\n", err)
		return nil, err
	}
	if (len(bytesRead) != 0) {
		fmt.Printf("DateMap exists.\n")
		err = json.Unmarshal(bytesRead, &date_map)
		if err != nil {
			fmt.Printf("Failed to initialize  DateMap for block chain :%v\n", err)
			return nil, err
		}
	}

	//logic to type conversion of date time type
	var layout = "2006-01-02 15:04:05 -0700 MST"
	t, err = time.Parse(layout, trans_date)
	fmt.Println(t, err)

	if err != nil {
		fmt.Printf("Failed to convert input into date time format :%v\n", err)
		return nil, err
	} 
		
	refNo=date_map[t]
	bytesRead, err = json.Marshal(&refNo)
	return bytesRead, err

}


func ListAllTransactions(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var bytesRead []byte
	var trans_list []transactionInitiation	

	fmt.Println("Entering listAllTransactions")

	err = GetTransactionInitiationMap(stub)

	if err != nil {
		fmt.Printf("Unable to read the list of AllTransactions : %s\n", err)
		return nil, err
	}

	for _, value := range trans_Init_map {
		trans_list = append(trans_list, value)
	}
	fmt.Printf("list of AllTransactions : %v\n", trans_list)
	bytesRead, err = json.Marshal(&trans_list)
	fmt.Printf("list of AllTransactions after Marshal : %v\n", bytesRead)
	if err != nil {
		fmt.Printf("Unable to return the list of AllTransactions : %s\n", err)
		return nil, err
	}

	return bytesRead, nil
}

//Query methods ends here 

func GetTransactionInitiationMap(stub shim.ChaincodeStubInterface) error {
	var err error
	var bytesread []byte

	bytesread, err = stub.GetState("TransactionInitiationMap")
	if err != nil {
		fmt.Printf("Failed to get  Transaction initiation for block chain :%v\n", err)
		return err
	}
	if len(bytesread) != 0 {
		fmt.Printf("TransactionInitiationMap map exists.\n")
		err = json.Unmarshal(bytesread, &trans_Init_map)
		if err != nil {
			fmt.Printf("Failed to initialize  TransactionInitiationMap for block chain :%v\n", err)
			return err
		}
	} else {
		fmt.Printf("TransactionInitiationMap map does not exist. To be created. \n")
		trans_Init_map = make(map[string]transactionInitiation)
		bytesread, err = json.Marshal(&trans_Init_map)
		if err != nil {
			fmt.Printf("Failed to initialize  TransactionInitiationMap for block chain :%v\n", err)
			return err
		}
		err = stub.PutState("TransactionInitiationMap", bytesread)
		if err != nil {
			fmt.Printf("Failed to initialize  TransactionInitiationMap for block chain :%v\n", err)
			return err
		}
	}
	return nil
}


func GetSupervisorMap(stub shim.ChaincodeStubInterface) error {
	var err error
	var bytesread []byte

	bytesread, err = stub.GetState("SupervisorMap")
	if err != nil {
		fmt.Printf("Failed to get  SupervisorMap for block chain :%v\n", err)
		return err
	}
	if len(bytesread) != 0 {
		fmt.Printf("SupervisorMap map exists.\n")
		err = json.Unmarshal(bytesread, &supervisor_map)
		if err != nil {
			fmt.Printf("Failed to initialize  SupervisorMap for block chain :%v\n", err)
			return err
		}
	} else {
		fmt.Printf("SupervisorMap map does not exist. To be created. \n")
		supervisor_map = make(map[string]amlCheck)
		bytesread, err = json.Marshal(&supervisor_map)
		if err != nil {
			fmt.Printf("Failed to initialize  SupervisorMap for block chain :%v\n", err)
			return err
		}
		err = stub.PutState("SupervisorMap", bytesread)
		if err != nil {
			fmt.Printf("Failed to initialize  SupervisorMap for block chain :%v\n", err)
			return err
		}
	}
	return nil
}

func GetL1AuthMap(stub shim.ChaincodeStubInterface) error {
	var err error
	var bytesread []byte

	bytesread, err = stub.GetState("L1AuthMap")
	if err != nil {
		fmt.Printf("Failed to get  Transaction initiation for block chain :%v\n", err)
		return err
	}
	if len(bytesread) != 0 {
		fmt.Printf("L1AuthMap map exists.\n")
		err = json.Unmarshal(bytesread, &l1Auth_map)
		if err != nil {
			fmt.Printf("Failed to initialize  L1AuthMap for block chain :%v\n", err)
			return err
		}
	} else {
		fmt.Printf("L1AuthMap map does not exist. To be created. \n")
		l1Auth_map = make(map[string]l1Auth)
		bytesread, err = json.Marshal(&l1Auth_map)
		if err != nil {
			fmt.Printf("Failed to initialize  L1AuthMap for block chain :%v\n", err)
			return err
		}
		err = stub.PutState("L1AuthMap", bytesread)
		if err != nil {
			fmt.Printf("Failed to initialize  L1AuthMap for block chain :%v\n", err)
			return err
		}
	}
	return nil
}

func GetL2AuthMap(stub shim.ChaincodeStubInterface) error {
	var err error
	var bytesread []byte

	bytesread, err = stub.GetState("L2AuthMap")
	if err != nil {
		fmt.Printf("Failed to get  Transaction initiation for block chain :%v\n", err)
		return err
	}
	if len(bytesread) != 0 {
		fmt.Printf("L2AuthMap map exists.\n")
		err = json.Unmarshal(bytesread, &l2Auth_map)
		if err != nil {
			fmt.Printf("Failed to initialize L2AuthMap for block chain :%v\n", err)
			return err
		}
	} else {
		fmt.Printf("L2AuthMapL2AuthMapL2AuthMap map does not exist. To be created. \n")
		l2Auth_map = make(map[string]l2Auth)
		bytesread, err = json.Marshal(&l2Auth_map)
		if err != nil {
			fmt.Printf("Failed to initialize  L2AuthMapL2AuthMap for block chain :%v\n", err)
			return err
		}
		err = stub.PutState("L2AuthMap", bytesread)
		if err != nil {
			fmt.Printf("Failed to initialize  L2AuthMap for block chain :%v\n", err)
			return err
		}
	}
	return nil
}

//setTransactionInitiationMap
func SetTransactionInitiationMap(stub shim.ChaincodeStubInterface) error {
	var err error
	var bytesread []byte

	bytesread, err = json.Marshal(&trans_Init_map)
	if err != nil {
		fmt.Printf("Failed to set the TransactionItemMap for block chain :%v\n", err)
		return err
	}
	err = stub.PutState("TransactionItemMap", bytesread)
	if err != nil {
		fmt.Printf("Failed to set the TransactionItemMap %v\n", err)
		return errors.New("Failed to set the TransactionItemMap")
	}

	return nil
}

//setSupervisorMap
func SetSupervisorMap(stub shim.ChaincodeStubInterface) error {
	var err error
	var bytesread []byte

	bytesread, err = json.Marshal(&supervisor_map)
	if err != nil {
		fmt.Printf("Failed to set the SupervisorMap for block chain :%v\n", err)
		return err
	}
	err = stub.PutState("SupervisorMap", bytesread)
	if err != nil {
		fmt.Printf("Failed to set the SupervisorMap %v\n", err)
		return errors.New("Failed to set the SupervisorMap")
	}

	return nil
}

//setL1AuthMap
func SetL1AuthMap(stub shim.ChaincodeStubInterface) error {
	var err error
	var bytesread []byte

	bytesread, err = json.Marshal(&l1Auth_map)
	if err != nil {
		fmt.Printf("Failed to set the L1AuthMap for block chain :%v\n", err)
		return err
	}
	err = stub.PutState("L1AuthMap", bytesread)
	if err != nil {
		fmt.Printf("Failed to set the L1AuthMap %v\n", err)
		return errors.New("Failed to set the L1AuthMap")
	}

	return nil
}

//setL2AuthMap
func SetL2AuthMap(stub shim.ChaincodeStubInterface) error {
	var err error
	var bytesread []byte

	bytesread, err = json.Marshal(&l2Auth_map)
	if err != nil {
		fmt.Printf("Failed to set the L2AuthMap for block chain :%v\n", err)
		return err
	}
	err = stub.PutState("L2AuthMap", bytesread)
	if err != nil {
		fmt.Printf("Failed to set the L2AuthMap %v\n", err)
		return errors.New("Failed to set the L2AuthMap")
	}

	return nil
}

// Init sets up the chaincode
func (t *SBITransaction) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Inside INIT for test chaincode")
	return nil, nil
}

// Query the chaincode
func (t *SBITransaction) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	if function == "GetTransactionInitDetailsForRefAndMaker" {
		return GetTransactionInitDetailsForRefAndMaker(stub, args)
	} else if function == "ListRefnoForDate" {
		return ListRefnoForDate(stub, args)
	} else if function == "GetAllDetailsForRef_AuditTrial" {
		return GetAllDetailsForRef_AuditTrial(stub, args)
	} else if function == "GetTransactionInitDetailsForRef" {
		return GetTransactionInitDetailsForRef(stub, args)
	} else if function == "ListRefnoForBranch" {
		return ListRefnoForBranch(stub, args)
	} else if function == "ListAllTransactions" {
		return ListAllTransactions(stub, args)
	}
	return nil, nil
}

// Invoke the function in the chaincode
func (t *SBITransaction) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "CreateTransaction" {
		return CreateTransaction(stub,args)
	} else if function == "UpdateSupervisorDetails" {
		return nil, UpdateSupervisorDetails(stub, args)
	} else if function == "UpdateL1AuthorizerDetails" {
		return nil, UpdateL1AuthorizerDetails(stub, args)
	} else if function == "UpdateL2AuthorizerDetails" {
		return nil, UpdateL2AuthorizerDetails(stub,args)
	} else if function == "UpdateEventDesc" {
		return nil, UpdateEventDesc(stub,args)
	} 
	fmt.Println("Function not found")
	return nil, nil
}

func main() {
	err := shim.Start(new(SBITransaction))
	if err != nil {
		fmt.Println("Could not start SBITransaction")
	} else {
		fmt.Println("SBITransaction successfully started")
	}

}