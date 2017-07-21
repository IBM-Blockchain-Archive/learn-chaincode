package main
 
import (
	//"bytes"
	"errors"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
    PEER1 = "PHI"
    PEER2 = "SHELL"
    PEER3 = "CVXG"
	PEER4 = "BP"
	CONTRACT_KEY = "_Contract"
)
	
type SmartContract struct {
}

type FlightShrContract struct {
	OwnerCompany	string `json:"ownerCompany"`
	PercSeatAlloc	uint8 `json:"percSeatAlloc"`
}

type Flight struct {
	FlightKey   	string `json:"flightKey"`
	FlightName  	string `json:"flightName"`
	OwnerCompany  	string `json:"ownerCompany"`
	FlightType 		string `json:"flightType"`
	SlNo  			string `json:"slNo"`
	Origin			string `json:"origin"`
	Destination		string `json:"destination"`
	DeptDate		string `json:"deptDate"`
	DeptTime		string `json:"deptTime"`
	ArrDate			string `json:"arrDate"`
	ArrTime			string `json:"arrTime"`
	NoOfSeats		uint8 `json:"noOfSeats"`
	NoOfStops		uint8 `json:"noOfStops"`
	LegDetails		[]FlightLeg `json:"legDetails"`
}

type FlightLeg struct {
	Origin		string `json:"origin"`
	Destination	string `json:"destination"`
	DeptDate	string `json:"deptDate"`
	DeptTime	string `json:"deptTime"`
	ArrDate		string `json:"arrDate"`
	ArrTime		string `json:"arrTime"`
	TravelMode	string `json:"travelMode"`
	LegNo		uint8 `json:"legNo"`
	AvailSeats	uint8 `json:"availSeats"`
}

func (t *SmartContract) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("--> shareflight initialized...")
	return nil, nil
}
 
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Retrieve the requested Smart Contract function and arguments
	//function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryFlight" {
		return s.queryFlight(APIstub, args)
	}else if function == "createFlight" {
		return s.createFlight(APIstub, args)
	} else if function == "queryAllFlights" {
		return s.queryAllFlights(APIstub,args)
	} else if function == "changeFlightOwner" {
		return s.changeFlightOwner(APIstub, args)
	}
	return nil, errors.New("Invalid Smart Contract function name.")
}

func (s *SmartContract) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return s.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

func (s *SmartContract) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (s *SmartContract) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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

func (s *SmartContract) queryFlight(APIstub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	flightAsBytes, _ := APIstub.GetState(args[0])
	return flightAsBytes, nil
}

func (s *SmartContract) createFlight(APIstub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//legDetails1 := []FlightLeg{FlightLeg{Origin: "LOC1", Destination: "LOC2", DeptDate:"19-07-2017", DeptTime:"10:00", ArrDate:"19-07-2017", ArrTime:"11:00", TravelMode: "Fixed Wing", LegNo: 1, AvailSeats: 100},FlightLeg{Origin: "LOC2", Destination: "LOC3", DeptDate:"19-07-2017", DeptTime:"11:10", ArrDate:"19-07-2017", ArrTime:"12:30", TravelMode: "Fixed Wing", LegNo: 1, AvailSeats: 100}}
	
	//flight := Flight{FlightKey: "Flight#", FlightName: "TEST_FLT", OwnerCompany: "PHI", FlightType: "FTYPE1", SlNo: "SL01", Origin: "LOC1", Destination: "LOC3", DeptDate: "19-07-2017", DeptTime: "10:00", ArrDate: "19-07-2017", ArrTime: "12:30", NoOfSeats: 100, NoOfStops: 1, LegDetails: legDetails1}
	
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
    flight := Flight{}
	flight_json := args[1]
	flightByteArray := []byte(flight_json)
    err := json.Unmarshal(flightByteArray, &flight)
    if err != nil {
        fmt.Println("Error while parsing file")
        return nil, errors.New("Incorrect number of arguments. Expecting 2")
    }
	
	createSharedFlights(APIstub,flight)
	return nil, nil
}

func createSharedFlights(APIstub shim.ChaincodeStubInterface,flight Flight) ([]Flight, error){
	fltShrContracts := []FlightShrContract{FlightShrContract{OwnerCompany: PEER2, PercSeatAlloc: 20},FlightShrContract{OwnerCompany: PEER3, PercSeatAlloc: 30},FlightShrContract{OwnerCompany: PEER4, PercSeatAlloc: 10}}
	totalSeats := flight.NoOfSeats
	availSeat  := totalSeats
	var noOfSeats uint8
	var key string
	i := 0
	for i < len(fltShrContracts) {
		fltShrContract := fltShrContracts[i]
		if(fltShrContract.PercSeatAlloc > 0){
			noOfSeats = totalSeats * (fltShrContract.PercSeatAlloc/100);
			if(availSeat>=noOfSeats){
				newFlight := prepareFlight(flight,noOfSeats,&availSeat,fltShrContract.OwnerCompany)
				key = fmt.Sprintf("%s%s%d", newFlight.OwnerCompany, "_F", i)
				addFlightToLedger(APIstub,key,newFlight)
			}
		}
		i++;
	}
	if(availSeat>0){
		newFlight := prepareFlight(flight,availSeat,&availSeat,flight.OwnerCompany)
		key = fmt.Sprintf("%s%s%d", newFlight.OwnerCompany, "_F", i)
		addFlightToLedger(APIstub,key,newFlight)
	}
	return nil, nil
}

func prepareFlight(flight Flight, noOfSeats uint8, availSeat *uint8, ownerCompany string) Flight {
	newFlight := Flight{}
	*(&newFlight)= *(&flight)
	newFlight.LegDetails = copyLegDetails(flight.LegDetails,noOfSeats)
	newFlight.OwnerCompany = ownerCompany
	newFlight.NoOfSeats = noOfSeats
	*availSeat = *availSeat - noOfSeats
	return newFlight
}

func copyLegDetails(flightLegs []FlightLeg, noOfSeats uint8) []FlightLeg{
	var newFlightLegs []FlightLeg
	var flightLeg FlightLeg
	i := 0
	for i < len(flightLegs) {
		flightLeg = FlightLeg{}
		*(&flightLeg) = *(&flightLegs[i])
		flightLeg.AvailSeats = noOfSeats
		newFlightLegs = append(newFlightLegs, flightLeg)
		i++
	}
	return newFlightLegs
}

func addFlightToLedger(APIstub shim.ChaincodeStubInterface, key string, flight Flight) ([]byte, error) {

	fmt.Printf(">> start writing flight to ledger - Key:",key)
	flightAsBytes, _ := json.Marshal(flight)
	APIstub.PutState(key, flightAsBytes)
	fmt.Printf(">> writing flight to ledger completed- Key:",key)
	return nil, nil
}

func (s *SmartContract) queryAllFlights(APIstub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	//startKey :=  args[0]
	//endKey :=  args[1]

	/*resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, errors.New(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf(">> queryAllFlights:\n%s\n", buffer.String())

	return buffer.Bytes(), nil*/
	return nil, nil
}

func (s *SmartContract) changeFlightOwner(APIstub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	flightAsBytes, _ := APIstub.GetState(args[0])
	flight := Flight{}

	json.Unmarshal(flightAsBytes, &flight)
	flight.OwnerCompany = args[1]

	flightAsBytes, _ = json.Marshal(flight)
	APIstub.PutState(args[0], flightAsBytes)

	return nil, nil
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}