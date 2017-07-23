package main
 
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

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

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) ([]byte, error) {
	return shim.Success(nil)
}
 
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) ([]byte, error) {
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryFlight" {
		return s.queryFlight(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createFlight" {
		return s.createFlight(APIstub, args)
	} else if function == "queryAllFlights" {
		return s.queryAllFlights(APIstub)
	} else if function == "changeFlightOwner" {
		return s.changeFlightOwner(APIstub, args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryFlight(APIstub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	flightAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(flightAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) ([]byte, error) {
	legDetails1 := []FlightLeg{FlightLeg{Origin: "LOC1", Destination: "LOC2", DeptDate:"19-07-2017", DeptTime:"10:00", ArrDate:"19-07-2017", ArrTime:"11:00", TravelMode: "Fixed Wing", LegNo: 1, AvailSeats: 100},FlightLeg{Origin: "LOC2", Destination: "LOC3", DeptDate:"19-07-2017", DeptTime:"11:10", ArrDate:"19-07-2017", ArrTime:"12:30", TravelMode: "Fixed Wing", LegNo: 1, AvailSeats: 100}}
	
	flight := Flight{FlightKey: "Flight#", FlightName: "TEST_FLT", OwnerCompany: "PHI", FlightType: "FTYPE1", SlNo: "SL01", Origin: "LOC1", Destination: "LOC3", DeptDate: "19-07-2017", DeptTime: "10:00", ArrDate: "19-07-2017", ArrTime: "12:30", NoOfSeats: 100, NoOfStops: 1, LegDetails: legDetails1}
	flights := getSharedFlights(flight)
	i := 0
	for i < len(flights) {
		fmt.Println("i is ", i)
		flights[i].FlightKey = "Flight#"+strconv.Itoa(i)
		flightAsBytes, _ := json.Marshal(flights[i])
		APIstub.PutState(flights[i].FlightKey, flightAsBytes)
		fmt.Println("Added", flights[i])
		i = i + 1
	}
	return shim.Success(nil)
}

func getSharedFlights(flight Flight) ([]Flight, error){
	fltShrContracts := []FlightShrContract{FlightShrContract{OwnerCompany: PEER2, PercSeatAlloc: 20},FlightShrContract{OwnerCompany: PEER3, PercSeatAlloc: 30},FlightShrContract{OwnerCompany: PEER4, PercSeatAlloc: 10}}
	totalSeats := flight.NoOfSeats
	availSeat  := totalSeats
	var noOfSeats uint8
	i := 0
	for i < len(fltShrContracts) {
		fltShrContract := fltShrContracts[i]
		if(fltShrContract.PercSeatAlloc > 0){
			noOfSeats = totalSeats * (fltShrContract.PercSeatAlloc/100);
			if(availSeat>=noOfSeats){
				newFlight := prepareFlight(flight,noOfSeats,&availSeat)
				createFlight(newFlight)
			}
		}		
	}
	if(availSeat>0){
		newFlight = prepareFlight(flight,availSeat,&availSeat)
		createFlight(newFlight)
	}
}

func prepareFlight(flight Flight, noOfSeats uint8, availSeat *uint8) (Flight,error) {
	newFlight := Flight{}
	copy(&newFlight, &flight)
	newFlight.LegDetails = copyLegDetails(flight.LegDetails,noOfSeats)
	newFlight.OwnerCompany = fltShrContract.OwnerCompany
	newFlight.NoOfSeats = noOfSeats
	createFlight(newFlight)
	availSeat = availSeat - noOfSeats
	return newFlight
}

func copyLegDetails(flightLegs []FlightLeg, noOfSeats uint8) ([]Flight, error){
	var newFlightLegs = make([]FlightLeg,3,5)
	var flightLeg FlightLeg
	for i < len(flightLegs) {
		flightLeg = FlightLeg{}
		copy(&flightLeg, &flightLegs[0])
		flightLeg.AvailSeats = noOfSeats
		newFlightLegs.append(flightLeg)		
	}
	return newFlightLegs
}

func (s *SmartContract) createFlight(APIstub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	
	legDetails := make([]FlightLeg, 0)
    flight := jsonData{legDetails}
	flight_json := args[1];
	flightByteArray := []byte(myString)
    err = json.Unmarshal(flightByteArray, &flight)
    if err != nil {
        fmt.Println("Error while parsing file")
        return shim.Error("Incorrect number of arguments. Expecting 2")
    }
	newflights := getSharedFlights(flight)
	if len(newflights) == 0 {
		return shim.Error("No flights to create")
	}
	
	for i < len(newflights) {
		fmt.Println("---> Adding new flight for ",flight[i].OwnerCompany)
		flightAsBytes, _ := json.Marshal(flight[i])
		APIstub.PutState(args[0], flightAsBytes)
		fmt.Println("---> Flight added successfully for ",flight[i].OwnerCompany)
	}
	return shim.Success(nil)
}

func (s *SmartContract) queryAllFlights(APIstub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	startKey :=  args[0]
	endKey :=  args[1]

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
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

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeFlightOwner(APIstub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	flightAsBytes, _ := APIstub.GetState(args[0])
	flight := Flight{}

	json.Unmarshal(flightAsBytes, &flight)
	flight.Owner = args[1]

	flightAsBytes, _ = json.Marshal(flight)
	APIstub.PutState(args[0], flightAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
/*func main() {
	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}*/