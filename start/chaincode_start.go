package main

/* import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
) */

import (
	"errors"
	"fmt"
	//"strconv"
	//"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"encoding/json"
	//"regexp"
	//"time"
	//"crypto/md5"
	//"io"
)

var logger = shim.NewLogger("CLDChaincode")

// Participant
const	SHIPPER      =  "shipper"
const	LOGISTIC_PROVIDER   =  "logistic_provider"
const	INSURENCE_COMPANY = "insurence_company"

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Volume struct {
	NextStop								string `json: nextStop`
	Origin									Origin
	Destination								Destination
	LogisticProvider						LogisticProvider
	Volume									VolumeD
	Event									Event
}

type Origin struct {
	Name 									string `json: "name"`
	Address 								string `json: "address"`
	FederalTaxPayerId						string `json: "federalTaxPayerId"`
	AddressNumber							int	   `json: "addressNumber"`
	ZipCode									string `json: "zipCode"`
}

type EndCustomerFinal struct {
	Name 									string `json: "name"`
	FederalTaxPayerId						string `json: "federalTaxPayerId"`
	Address 								string `json: "address"`
	AddressNumber							int	   `json: "addressNumber"`
	ZipCode									string `json: "zipCode"`
	City									string `json: "city"`
	Quarter									string `json: "quarter"`
	Email									string `json: "email"`
	Phone									string `json: "phone"`
	Cellphone								string `json: "cellphone"`
}

type Destination struct {
	EndCustomer 							EndCustomerFinal
	ShipperEstimatedDeliveryDate 			string `json: "shipperEstimatedDeliveryDate"`
	LogisticProviderEstimatedDeliveryDate	string `json: "logisticProviderEstimatedDeliveryDate"`
}

type LogisticProvider struct {
	Id										string `json: "id"`
	Name 									string `json: "name"`
	Address 								string `json: "address"`
	AddressNumber							int	   `json: "addressNumber"`
	ZipCode									string `json: "zipCode"`
	City									string `json: "city"`
	Quarter									string `json: "quarter"`
}

type VolumeD struct {
	TrackId									string `json: "trackId"` 
	VolumeData								VolumeData														
}

type VolumeData struct {
	Key										string `json: key`
}

type Event struct {
	Date 									string `json: date`
	StatusCode								string `json: statusCode`
	Description								string `json: description`
	LogisticProviderProperties				string `json: logisticProviderPropertis`
}

func main() {
	fmt.Println("[IP] Start Contract")

	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("invoke is running " + function)

    // Handle different functions
    if function == "init" {
        return t.Init(stub, "init", args)
	} else if function == "createVolume" {
		return t.createVolume(stub)
	} /* else if function == "shipperToLogisticProvider" {
        return t.shipperToLogisticProvider(stub, args)
    } else if function == "LogisticProviderToCustomer" {
		return t.LogisticProviderToCustomer(stub, args)
	} else if function == "LogisticProviderToLogisticProvider" {
		return t.LogisticProviderToLogisticProvider(stub, args)
	} else if function == "LogisticProviderToShipper" {
		return t.LogisticProviderToShipper(stub, args)
	} */

    fmt.Println("invoke did not find func: " + function)

    return nil, errors.New("Received unknown function invocation")
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("query is running " + function)

    /* if function == "read" {                            //read a variable
        return t.read(stub, args)
    }
    fmt.Println("query did not find func: " + function)
 */
    return nil, errors.New("Received unknown function query")
}

// Functions to Write
func (t *SimpleChaincode) createVolume(stub shim.ChaincodeStubInterface) ([]byte, error) {
	var v Volume

	//v.NextStop       	= nil
	//v.Origin         	= nil
	//v.Destination    	= nil
	//v.LogisticProvifer	= nil
	//v.Event         	= nil
	v.Volume.TrackId	= "chave";

 	fmt.Println("[IP][Volume]: fmt" + v.Volume.TrackId)
	logger.Debug("[IP][Volume]: logger", v)

	return nil, nil
	//err = json.Unmarshal([]byte(volume_json), &v)

	// if volume already exists

	/* record, err := stub.GetState(v.id)
	if record != nil { return nil, errors.New("Volume already exists") } */

	//_, err  = t.save_changes(stub, v)

	//if err != nil { fmt.Printf("Create_Volume: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }
}

/* func (t *SimpleChaincode) save_changes(stub shim.ChaincodeStubInterface, v Volume) (bool, error) {

	bytes, err := json.Marshal(v)

	if err != nil { fmt.Printf("SAVE_CHANGES: Error converting vehicle record: %s", err); return false, errors.New("Error converting volume record") }

	err = stub.PutState(v.id, bytes)

	if err != nil { fmt.Printf("SAVE_CHANGES: Error storing volume record: %s", err); return false, errors.New("Error storing volume record") }

	return true, nil
} */

// Functions to Read