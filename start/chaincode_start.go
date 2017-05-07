package main

import (
	"errors"
	"fmt"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"regexp"
)

var logger = shim.NewLogger("CLDChaincode")

//==============================================================================================================================
//	 Participant types - Each participant type is mapped to an integer which we use to compare to the value stored in a
//						 user's eCert
//==============================================================================================================================
//CURRENT WORKAROUND USES ROLES CHANGE WHEN OWN USERS CAN BE CREATED SO THAT IT READ 1, 2, 3, 4, 5
const   AUTHORITY		=  "manufacturer"
const   MANUFACTURER	=  "manufacturer"
const   DISTRIBUTOR			=  "distributor"
const   RETAILER		=  "retailer"
const   CONSUMER		=  "consumer"
//const   SCRAP_MERCHANT =  "scrap_merchant" //TO BE DONE check if this variable is used anywhere


//==============================================================================================================================
//	 Status types - Asset lifecycle is broken down into 5 statuses, this is part of the business logic to determine what can
//					be done to the spare part at points in it's lifecycle
//==============================================================================================================================
const   STATE_TEMPLATE  			=  0
const   STATE_MANUFACTURE  			=  0
const   STATE_DISTRIBUTOR 				=  1
const   STATE_RETAILER 				=  2
const   STATE_CONSUMER  			=  3

//==============================================================================================================================
//	 Structure Definitions
//==============================================================================================================================
//	Chaincode - A blank struct for use with Shim (A HyperLedger included go file used for get/put state
//				and other HyperLedger functions)
//==============================================================================================================================
type  SimpleChaincode struct {
}

//==============================================================================================================================
//	SparePart - Defines the structure for a car object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON make -> Struct Make.
//==============================================================================================================================
type SparePart struct {
	PartName            string `json:"partname"`
	PartNumber          string `json:"partnumber"`
	Description         string `json:"description"`
	PartID              string `json:"partid"`
	ManufacturingDate   string `json:"mfgdate"`
	Owner               string `json:"owner"`
	Remarks             string `json:"remarks"`
	Status              int    `json:"status"`
}
//==============================================================================================================================
//	partID Holder - Defines the structure that holds all the PartIDs for spare parts that have been created.
//				Used as an index when querying all spare parts.
//==============================================================================================================================

type PartID_Holder struct {
	PartIds 	[]string `json:"partids"`
}




//==============================================================================================================================
//	User_and_eCert - Struct for storing the JSON of a user and their ecert
//==============================================================================================================================

type User_and_eCert struct {
	Identity string `json:"identity"`
	eCert string `json:"ecert"`
}
//==============================================================================================================================
//	Init Function - Called when the user deploys the chaincode
//==============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	//Args
	//				0
	//			peer_address

	var idHolder PartID_Holder

	bytes, err := json.Marshal(idHolder)

    if err != nil { return nil, errors.New("Error creating PartID_Holder record") }

	err = stub.PutState("PartIDs", bytes)

	for i:=0; i < len(args); i=i+2 {
		t.add_ecert(stub, args[i], args[i+1])
	}

	return nil, nil
}

//==============================================================================================================================
//	 General Functions
//==============================================================================================================================
//	 get_ecert - Takes the name passed and calls out to the REST API for HyperLedger to retrieve the ecert
//				 for that user. Returns the ecert as retrived including html encoding.
//==============================================================================================================================
func (t *SimpleChaincode) get_ecert(stub shim.ChaincodeStubInterface, name string) ([]byte, error) {

	ecert, err := stub.GetState(name)

	if err != nil { return nil, errors.New("Couldn't retrieve ecert for user " + name) }

	return ecert, nil
}

//==============================================================================================================================
//	 add_ecert - Adds a new ecert and user pair to the table of ecerts
//==============================================================================================================================

func (t *SimpleChaincode) add_ecert(stub shim.ChaincodeStubInterface, name string, ecert string) ([]byte, error) {


	err := stub.PutState(name, []byte(ecert))

	if err == nil {
		return nil, errors.New("Error storing eCert for user " + name + " identity: " + ecert)
	}

	return nil, nil

}

//==============================================================================================================================
//	 get_caller - Retrieves the username of the user who invoked the chaincode.
//				  Returns the username as a string.
//==============================================================================================================================

func (t *SimpleChaincode) get_username(stub shim.ChaincodeStubInterface) (string, error) {

    username, err := stub.ReadCertAttribute("username");
	if err != nil { return "", errors.New("Couldn't get attribute 'username'. Error: " + err.Error()) }
	return string(username), nil
}

//==============================================================================================================================
//	 check_affiliation - Takes an ecert as a string, decodes it to remove html encoding then parses it and checks the
// 				  		certificates common name. The affiliation is stored as part of the common name.
//==============================================================================================================================

func (t *SimpleChaincode) check_affiliation(stub shim.ChaincodeStubInterface) (string, error) {
    affiliation, err := stub.ReadCertAttribute("role");
	if err != nil { return "", errors.New("Couldn't get attribute 'role'. Error: " + err.Error()) }
	return string(affiliation), nil

}

//==============================================================================================================================
//	 get_caller_data - Calls the get_ecert and check_role functions and returns the ecert and role for the
//					 name passed.
//==============================================================================================================================

func (t *SimpleChaincode) get_caller_data(stub shim.ChaincodeStubInterface) (string, string, error){

	user, err := t.get_username(stub)

    // if err != nil { return "", "", err }

	// ecert, err := t.get_ecert(stub, user);

    // if err != nil { return "", "", err }

	affiliation, err := t.check_affiliation(stub);

    if err != nil { return "", "", err }

	return user, affiliation, nil
}

//==============================================================================================================================
//	 retrieve_part - Gets the state of the data at partID in the ledger then converts it from the stored
//					JSON into the SparePart struct for use in the contract. Returns the part struct.
//					Returns empty v if it errors.
//==============================================================================================================================
func (t *SimpleChaincode) retrieve_part(stub shim.ChaincodeStubInterface, partID string) (SparePart, error) {

	var sp SparePart

	bytes, err := stub.GetState(partID);

	if err != nil {	fmt.Printf("RETRIEVE_PART: Failed to invoke sparepart_code: %s", err); return sp, errors.New("RETRIEVE_PART: Error retrieving spare part with PartID = " + partID) }

	err = json.Unmarshal(bytes, &sp);

    if err != nil {	fmt.Printf("RETRIEVE_PART: Corrupt spare part record "+string(bytes)+": %s", err); return sp, errors.New("RETRIEVE_PART: Corrupt spare part record"+string(bytes))	}

	return sp, nil
}

//==============================================================================================================================
// save_changes - Writes to the ledger the SparePart struct passed in a JSON format. Uses the shim file's
//				  method 'PutState'.
//==============================================================================================================================
func (t *SimpleChaincode) save_changes(stub shim.ChaincodeStubInterface, sp SparePart) (bool, error) {

	bytes, err := json.Marshal(sp)

	if err != nil { fmt.Printf("SAVE_CHANGES: Error converting spare part record: %s", err); return false, errors.New("Error converting spare part record") }

	err = stub.PutState(sp.PartID, bytes)

	if err != nil { fmt.Printf("SAVE_CHANGES: Error storing spare part record: %s", err); return false, errors.New("Error storing spare part record") }

	return true, nil
}

//==============================================================================================================================
//	 Router Functions
//==============================================================================================================================
//	Invoke - Called on chaincode invoke. Takes a function name passed and calls that function. Converts some
//		  initial arguments passed to other things for use in the called function e.g. name -> ecert
//==============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	caller, caller_affiliation, err := t.get_caller_data(stub)

	if err != nil { return nil, errors.New("Error retrieving caller information")}


	if function == "create_part" {
        return t.create_part(stub, caller, caller_affiliation, args[0])
	} else if function == "ping" {
        return t.ping(stub)
    } else { 																				// If the function is not a create then there must be a car so we need to retrieve the car.
		argPos := 1

		if function == "return_part" {																// If its a scrap spare part then only two arguments are passed (no update value) all others have three arguments and the PartID is expected in the last argument
			argPos = 0
		}

		v, err := t.retrieve_part(stub, args[argPos])

        if err != nil { fmt.Printf("INVOKE: Error retrieving v5c: %s", err); return nil, errors.New("Error retrieving v5c") }


        if strings.Contains(function, "update") == false && function != "return_part"    { 									// If the function is not an update or a scrappage it must be a transfer so we need to get the ecert of the recipient.


				if 		   function == "manufacturer_to_distributor"		{ return t.manufacturer_to_distributor(stub, v, caller, caller_affiliation, args[0], DISTRIBUTOR) 
				} else if  function == "manufacturer_to_private"	{ return t.manufacturer_to_private(stub, v, caller, caller_affiliation, args[0], "private")
				} else if  function == "distributor_to_distributor"			{ return t.distributor_to_distributor(stub, v, caller, caller_affiliation, args[0], DISTRIBUTOR)
				} else if  function == "distributor_to_retailer"			{ return t.distributor_to_retailer(stub, v, caller, caller_affiliation, args[0], RETAILER)
				} else if  function == "retailer_to_consumer"		{ return t.retailer_to_consumer(stub, v, caller, caller_affiliation, args[0], CONSUMER)
				} else if  function == "private_to_scrap_merchant"	{ return t.private_to_scrap_merchant(stub, v, caller, caller_affiliation, args[0], "scrap_merchant")
				}

		} else if function == "update_description"		{ return t.update_description(stub, v, caller, caller_affiliation, args[0])
		} else if function == "update_remarks"			{ return t.update_remarks(stub, v, caller, caller_affiliation, args[0])
		} else if function == "update_partnumber"		{ return t.update_partnumber(stub, v, caller, caller_affiliation, args[0])
		} else if function == "update_partname"			{ return t.update_partname(stub, v, caller, caller_affiliation, args[0])
        } else if function == "update_mfgdate"			{ return t.update_mfgdate(stub, v, caller, caller_affiliation, args[0])
		} else if function == "return_part"				{ return t.return_part(stub, v, caller, caller_affiliation) }

		return nil, errors.New("Function of the name "+ function +" doesn't exist.")

	}
}
//=================================================================================================================================
//	Query - Called on chaincode query. Takes a function name passed and calls that function. Passes the
//  		initial arguments passed are passed on to the called function.
//=================================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	caller, caller_affiliation, err := t.get_caller_data(stub)
	if err != nil { fmt.Printf("QUERY: Error retrieving caller details", err); return nil, errors.New("QUERY: Error retrieving caller details: "+err.Error()) }

    logger.Debug("function: ", function)
    logger.Debug("caller: ", caller)
    logger.Debug("affiliation: ", caller_affiliation)

	if function == "get_part_details" {
		if len(args) != 1 { fmt.Printf("Incorrect number of arguments passed"); return nil, errors.New("QUERY: Incorrect number of arguments passed") }
		v, err := t.retrieve_part(stub, args[0])
		if err != nil { fmt.Printf("QUERY: Error retrieving v5c: %s", err); return nil, errors.New("QUERY: Error retrieving v5c "+err.Error()) }
		return t.get_part_details(stub, v, caller, caller_affiliation)
	} else if function == "check_unique_partId" {
		return t.check_unique_partId(stub, args[0], caller, caller_affiliation)
	} else if function == "get_parts" {
		return t.get_parts(stub, caller, caller_affiliation)
	} else if function == "get_ecert" {
		return t.get_ecert(stub, args[0])
	} else if function == "ping" {
		return t.ping(stub)
	}

	return nil, errors.New("Received unknown function invocation " + function)

}

//=================================================================================================================================
//	 Ping Function
//=================================================================================================================================
//	 Pings the peer to keep the connection alive
//=================================================================================================================================
func (t *SimpleChaincode) ping(stub shim.ChaincodeStubInterface) ([]byte, error) {
	return []byte("Hello, world!"), nil
}

//=================================================================================================================================
//	 Create Function
//=================================================================================================================================
//	 Create SparePart - Creates the initial JSON for the vehcile and then saves it to the ledger.
//=================================================================================================================================
func (t *SimpleChaincode) create_part(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, partID string) ([]byte, error) {
	var v SparePart

	part_ID				:= "\"partID\":\""+partID+"\", "							// Variables to define the JSON
	partname			:= "\"PartName\":\"UNDEFINED\", "
	partnumber			:= "\"PartNumber\":0, "
	description			:= "\"Description\":\"UNDEFINED\", "
	mfgdate				:= "\"MfgDt\":\"UNDEFINED\", "
	owner				:= "\"Owner\":\""+caller+"\", "
	status				:= "\"Status\":0, "
	remarks       		:= "\"Remarks\":\"UNDEFINED\" "

	spare_part_json := "{"+part_ID+partname+partnumber+description+mfgdate+owner+status+remarks+"}" 	// Concatenates the variables to create the total JSON object
	matched, err := regexp.Match("^[A-z][A-z][0-9]{7}", []byte(partID))  				// matched = true if the PartID passed fits format of two letters followed by seven digits

												if err != nil { fmt.Printf("CREATE_PART: Invalid PartID: %s", err); return nil, errors.New("Invalid PartID") }

	if 				part_ID  == "" 	 ||
					matched == false    {
																		fmt.Printf("CREATE_PART: Invalid PartID provided");
																		return nil, errors.New("Invalid PartID provided")
	}

	err = json.Unmarshal([]byte(spare_part_json), &v)							// Convert the JSON defined above into a spare part object for go

																		if err != nil { return nil, errors.New("CREATE_PART:Invalid JSON object") }

	record, err := stub.GetState(v.PartID) 								// If not an error then a record exists so cant create a new car with this PartID as it must be unique

																		if record != nil { return nil, errors.New("CREATE_PART:SparePart already exists") }

	if 	caller_affiliation != AUTHORITY {							// Only the manufacturer can create a new part id

		return nil, errors.New(fmt.Sprintf("Permission Denied. create_part. %v === %v", caller_affiliation, AUTHORITY))

	}

	_, err  = t.save_changes(stub, v)

																		if err != nil { fmt.Printf("CREATE_PART: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	bytes, err := stub.GetState("PartIDs")

																		if err != nil { return nil, errors.New("CREATE_PART:Unable to get PartIds") }

	var idHolder PartID_Holder

	err = json.Unmarshal(bytes, &idHolder)

																		if err != nil {	return nil, errors.New("CREATE_PART:Corrupt PartID_Holder record") }

	idHolder.PartIds = append(idHolder.PartIds, partID)


	bytes, err = json.Marshal(idHolder)

															if err != nil { fmt.Print("CREATE_PART:Error creating PartID_Holder record") }

	err = stub.PutState("PartIDs", bytes)

															if err != nil { return nil, errors.New("CREATE_PART:Unable to put the state") }

	return nil, nil

}

//=================================================================================================================================
//	 Transfer Functions
//=================================================================================================================================
//	 manufacturer_to_distributor
//=================================================================================================================================
func (t *SimpleChaincode) manufacturer_to_distributor(stub shim.ChaincodeStubInterface, sp SparePart, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {

	if     	sp.Status				== STATE_MANUFACTURE	&&
			sp.Owner				== caller				&&
			caller_affiliation		== MANUFACTURER			&&
			recipient_affiliation	== DISTRIBUTOR{		

					sp.Owner  = recipient_name		// then make the owner the new owner
					sp.Status = STATE_DISTRIBUTOR			// and mark it in the state of distributor

	} else {									// Otherwise if there is an error
					fmt.Printf("MANUFACTURER_TO_DISTRIBUTOR: Permission Denied");
                    return nil, errors.New(fmt.Sprintf("Permission Denied. manufacturer_to_distributor. %sp %sp === %sp, %sp === %sp, %sp === %sp, %sp === %sp", sp, sp.Status, STATE_MANUFACTURE, sp.Owner, caller, caller_affiliation, DISTRIBUTOR, recipient_affiliation,DISTRIBUTOR))

	}

	_, err := t.save_changes(stub, sp)						// Write new state

															if err != nil {	fmt.Printf("MANUFACTURER_TO_DISTRIBUTOR: Error saving changes: %s", err); return nil, errors.New("Error saving changes")	}

	return nil, nil									// We are Done

}

//=================================================================================================================================
//	 manufacturer_to_private //TO BE DONE - NOT REQUIRED SO COMMENTING
//=================================================================================================================================
func (t *SimpleChaincode) manufacturer_to_private(stub shim.ChaincodeStubInterface, v SparePart, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {

	/*if 		v.Make 	 == "UNDEFINED" ||
			v.Model  == "UNDEFINED" ||
			v.Reg 	 == "UNDEFINED" ||
			v.Colour == "UNDEFINED" ||
			v.VIN == 0				{					//If any part of the car is undefined it has not been fully manufacturered so cannot be sent
															fmt.Printf("MANUFACTURER_TO_PRIVATE: Car not fully defined")
															return nil, errors.New(fmt.Sprintf("Car not fully defined. %v", v))
	}

	if 		v.Status				== STATE_MANUFACTURE	&&
			v.Owner					== caller				&&
			caller_affiliation		== MANUFACTURER			&&
			recipient_affiliation	== PRIVATE_ENTITY		&&
			v.Scrapped     == false							{

					v.Owner = recipient_name
					v.Status = STATE_PRIVATE_OWNERSHIP

	} else {
        return nil, errors.New(fmt.Sprintf("Permission Denied. manufacturer_to_private. %v %v === %v, %v === %v, %v === %v, %v === %v, %v === %v", v, v.Status, STATE_PRIVATE_OWNERSHIP, v.Owner, caller, caller_affiliation, PRIVATE_ENTITY, recipient_affiliation, SCRAP_MERCHANT, v.Scrapped, false))
    }

	_, err := t.save_changes(stub, v)

	if err != nil { fmt.Printf("MANUFACTURER_TO_PRIVATE: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }
	return nil, nil
*/
	return nil, errors.New(fmt.Sprint("Permission denied. manufacturer_to_private "))
	

}

//=================================================================================================================================
//	 distributor_to_distributor
//=================================================================================================================================
func (t *SimpleChaincode) distributor_to_distributor(stub shim.ChaincodeStubInterface, v SparePart, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {

	if 		v.Status				== STATE_DISTRIBUTOR		&&
			v.Owner					== caller			&&
			caller_affiliation		== DISTRIBUTOR			&&
			recipient_affiliation	== DISTRIBUTOR{

					v.Owner = recipient_name

	} else {
        return nil, errors.New(fmt.Sprintf("Permission Denied. distributor_to_distributor. %v %v === %v, %v === %v, %v === %v, %v === %v", v, v.Status, STATE_DISTRIBUTOR, v.Owner, caller, caller_affiliation, DISTRIBUTOR, recipient_affiliation,DISTRIBUTOR))
	}

	_, err := t.save_changes(stub, v)

															if err != nil { fmt.Printf("DISTRIBUTOR_TO_DISTRIBUTOR: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}

//=================================================================================================================================
//	 distributor_to_retailer
//=================================================================================================================================
func (t *SimpleChaincode) distributor_to_retailer(stub shim.ChaincodeStubInterface, v SparePart, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {

	if 		v.Status				== STATE_DISTRIBUTOR	&&
			v.Owner					== caller		&&
			caller_affiliation		== DISTRIBUTOR		&&
			recipient_affiliation	== RETAILER{

					v.Owner = recipient_name
					v.Status = STATE_RETAILER

	} else {
        return nil, errors.New(fmt.Sprintf("Permission denied. distributor_to_retailer. %v === %v, %v === %v, %v === %v, %v === %v", v.Status, STATE_DISTRIBUTOR, v.Owner, caller, caller_affiliation, DISTRIBUTOR, recipient_affiliation, RETAILER))

	}

	_, err := t.save_changes(stub, v)
															if err != nil { fmt.Printf("distributor_to_retailer: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}

//=================================================================================================================================
//	 retailer_to_consumer
//=================================================================================================================================
func (t *SimpleChaincode) retailer_to_consumer(stub shim.ChaincodeStubInterface, v SparePart, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {

	if		v.Status				== STATE_RETAILER	&&
			v.Owner  				== caller			&&
			caller_affiliation		== RETAILER			&&
			recipient_affiliation	== CONSUMER{

				v.Owner = recipient_name
				v.Status = STATE_CONSUMER

	} else {
		return nil, errors.New(fmt.Sprintf("Permission Denied. retailer_to_consumer. %v %v === %v, %v === %v, %v === %v, %v === %v", v, v.Status, STATE_RETAILER, v.Owner, caller, caller_affiliation, RETAILER, recipient_affiliation, CONSUMER))
	}

	_, err := t.save_changes(stub, v)
	if err != nil { fmt.Printf("RETAILER_TO_CONSUMER: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}

//=================================================================================================================================
//	 retailer_to_consumer
//=================================================================================================================================
func (t *SimpleChaincode) consumer_to_retailer(stub shim.ChaincodeStubInterface, v SparePart, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {

	if		v.Status				== STATE_CONSUMER	&&
			v.Owner  				== caller			&&
			caller_affiliation		== CONSUMER		&&
			recipient_affiliation	== RETAILER{

				v.Owner = recipient_name
				v.Status = STATE_RETAILER

	} else {
		return nil, errors.New(fmt.Sprintf("Permission Denied. retailer_to_consumer. %v %v === %v, %v === %v, %v === %v, %v === %v", v, v.Status, STATE_CONSUMER, v.Owner, caller, caller_affiliation,CONSUMER, recipient_affiliation, RETAILER))
	}

	_, err := t.save_changes(stub, v)
	if err != nil { fmt.Printf("CONSUMER_TO_RETAILER: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}

//=================================================================================================================================
//	 private_to_scrap_merchant
//=================================================================================================================================
func (t *SimpleChaincode) private_to_scrap_merchant(stub shim.ChaincodeStubInterface, v SparePart, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {

	/*if		v.Status				== STATE_PRIVATE_OWNERSHIP	&&
			v.Owner					== caller					&&
			caller_affiliation		== PRIVATE_ENTITY			&&
			recipient_affiliation	== SCRAP_MERCHANT			&&
			v.Scrapped				== false					{

					v.Owner = recipient_name
					v.Status = STATE_BEING_SCRAPPED

	} else {
        return nil, errors.New(fmt.Sprintf("Permission Denied. private_to_scrap_merchant. %v %v === %v, %v === %v, %v === %v, %v === %v, %v === %v", v, v.Status, STATE_PRIVATE_OWNERSHIP, v.Owner, caller, caller_affiliation, PRIVATE_ENTITY, recipient_affiliation, SCRAP_MERCHANT, v.Scrapped, false))
	}

	_, err := t.save_changes(stub, v)

															if err != nil { fmt.Printf("PRIVATE_TO_SCRAP_MERCHANT: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil
	*/
	return nil, errors.New(fmt.Sprint("Permission denied. "))

}

//=================================================================================================================================
//	 Update Functions
//=================================================================================================================================
//	 update_partname
//=================================================================================================================================
func (t *SimpleChaincode) update_partname(stub shim.ChaincodeStubInterface, v SparePart, caller string, caller_affiliation string, new_value string) ([]byte, error) {

	//new_vin, err := strconv.Atoi(string(new_value)) 		                // will return an error if the new vin contains non numerical chars

															//if err != nil || len(string(new_value)) != 15 { return nil, errors.New("Invalid value passed for new VIN") }

	if 		v.Status			== STATE_MANUFACTURE	&&
			v.Owner				== caller				&&
			caller_affiliation	== MANUFACTURER{

					v.PartName = new_value					// Update to the new value
	} else {

        return nil, errors.New(fmt.Sprintf("Permission denied. update_partname %v %v %v %v %v", v.Status, STATE_MANUFACTURE, v.Owner, caller, v.PartName))

	}

	_, err  := t.save_changes(stub, v)						// Save the changes in the blockchain

	if err != nil { fmt.Printf("UPDATE_PARTNAME: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}


//=================================================================================================================================
//	 update_partnumber
//=================================================================================================================================
func (t *SimpleChaincode) update_partnumber(stub shim.ChaincodeStubInterface, v SparePart, caller string, caller_affiliation string, new_value string) ([]byte, error) {


	if		v.Status			== STATE_MANUFACTURE	&&
			v.Owner				== caller				&&
			caller_affiliation	== MANUFACTURER{

					v.PartNumber = new_value

	} else {
        return nil, errors.New(fmt.Sprint("Permission denied. update_partnumber"))
	}

	_, err := t.save_changes(stub, v)

															if err != nil { fmt.Printf("UPDATE_PARTNUMBER: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}

//=================================================================================================================================
//	 update_mfgdate
//=================================================================================================================================
func (t *SimpleChaincode) update_mfgdate(stub shim.ChaincodeStubInterface, v SparePart, caller string, caller_affiliation string, new_value string) ([]byte, error) {

	if 		v.Status			== STATE_MANUFACTURE	&&
			v.Owner				== caller				&&
			caller_affiliation	== MANUFACTURER{

					v.ManufacturingDate = new_value
	} else {

		return nil, errors.New(fmt.Sprint("Permission denied. update_mfgdate"))
	}

	_, err := t.save_changes(stub, v)

		if err != nil { fmt.Printf("UPDATE_MFGDATE: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}

//=================================================================================================================================
//	 update_description
//=================================================================================================================================
func (t *SimpleChaincode) update_description(stub shim.ChaincodeStubInterface, v SparePart, caller string, caller_affiliation string, new_value string) ([]byte, error) {

	if 		v.Status			== STATE_MANUFACTURE	&&
			v.Owner				== caller				&&
			caller_affiliation	== MANUFACTURER{

					v.Description = new_value
	}else{

        return nil, errors.New(fmt.Sprint("Permission denied. update_description "))


	}

	_, err := t.save_changes(stub, v)

															if err != nil { fmt.Printf("UPDATE_DESCRIPTION: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}

//=================================================================================================================================
//	 update_remarks
//=================================================================================================================================
func (t *SimpleChaincode) update_remarks(stub shim.ChaincodeStubInterface, v SparePart, caller string, caller_affiliation string, new_value string) ([]byte, error) {

	if 		v.Status			== STATE_CONSUMER	&&
			v.Owner				== caller				&&
			caller_affiliation	== CONSUMER{

					v.Remarks = new_value

	}else {
        return nil, errors.New(fmt.Sprint("Permission denied. update_model "))

	}

	_, err := t.save_changes(stub, v)

	if err != nil { fmt.Printf("UPDATE_MODEL: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}

//=================================================================================================================================
//	 return_part
//=================================================================================================================================
func (t *SimpleChaincode) return_part(stub shim.ChaincodeStubInterface, v SparePart, caller string, caller_affiliation string) ([]byte, error) {

	/*if		v.Status			== STATE_BEING_SCRAPPED	&&
			v.Owner				== caller				&&
			caller_affiliation	== SCRAP_MERCHANT		&&
			v.Scrapped			== false				{

					v.Scrapped = true

	} else {
		return nil, errors.New("Permission denied. return_part")
	}

	_, err := t.save_changes(stub, v)

															if err != nil { fmt.Printf("RETURN_PART: Error saving changes: %s", err); return nil, errors.New("RETURN_PART Error saving changes") }

	return nil, nil*/
	return nil, errors.New(fmt.Sprint("Permission denied. return_part "))

}

//=================================================================================================================================
//	 Read Functions
//=================================================================================================================================
//	 get_part_details
//=================================================================================================================================
func (t *SimpleChaincode) get_part_details(stub shim.ChaincodeStubInterface, v SparePart, caller string, caller_affiliation string) ([]byte, error) {

	bytes, err := json.Marshal(v)

																if err != nil { return nil, errors.New("GET_PART_DETAILS: Invalid spare part object") }

	if 		v.Owner	== caller || caller_affiliation	== AUTHORITY{
				return bytes, nil
	} else {
			return nil, errors.New("Permission Denied. get_part_details")
	}

}

//=================================================================================================================================
//	 get_parts
//=================================================================================================================================

func (t *SimpleChaincode) get_parts(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string) ([]byte, error) {
	bytes, err := stub.GetState("PartIDs")

																			if err != nil { return nil, errors.New("Unable to get partIDs") }

	var idHolder PartID_Holder

	err = json.Unmarshal(bytes, &idHolder)

																			if err != nil {	return nil, errors.New("Corrupt PartID_Holder") }

	result := "["

	var temp []byte
	var v SparePart

	for _, partId := range idHolder.PartIds {

		v, err = t.retrieve_part(stub, partId)

		if err != nil {return nil, errors.New("Failed to retrieve partId")}

		temp, err = t.get_part_details(stub, v, caller, caller_affiliation)

		if err == nil {
			result += string(temp) + ","
		}
	}

	if len(result) == 1 {
		result = "[]"
	} else {
		result = result[:len(result)-1] + "]"
	}

	return []byte(result), nil
}

//=================================================================================================================================
//	 check_unique_partId
//=================================================================================================================================
func (t *SimpleChaincode) check_unique_partId(stub shim.ChaincodeStubInterface, partId string, caller string, caller_affiliation string) ([]byte, error) {
	_, err := t.retrieve_part(stub, partId)
	if err == nil {
		return []byte("false"), errors.New("PartID is not unique")
	} else {
		return []byte("true"), nil
	}
}

//=================================================================================================================================
//	 Main - main - Starts up the chaincode
//=================================================================================================================================
func main() {

	err := shim.Start(new(SimpleChaincode))

															if err != nil { fmt.Printf("Error starting Chaincode: %s", err) }
}