

package main

import (
	//"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Init create tables for tests
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Create table one
	err := createTableAuction(stub)
	if err != nil {
		return nil, fmt.Errorf("Error creating table one during init. %s", err)
	}

	

	return nil, nil
}


func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	switch function {

	case "insertRowAuction":
		if len(args) < 8 {
			return nil, errors.New("insertTableOne failed. Must include 8 column values")
		}


		//COMMENTED FOR hard coded values
		col1Val := args[0]
		col2Val := args[1]
		col3Val := args[2]
		col4Val := args[3]
		col5Val := args[4]
		col6Val := args[5]
		col7Val := args[6]
		col8Val := args[7]
		col9Val := args[8]
		col10Val := args[9]
		col11Val := args[10]
		col12Val := args[11]
		col13Val := args[12]
		col14Val := args[13]
		col15Val := args[14]
		col16Val := args[15]
		col17Val := args[16]
		col18Val := args[17]

		


		var columns []*shim.Column
		//col1 := shim.Column{Value: &shim.Column_String_{String_: col1Val}}
		
		col1 := shim.Column{Value: &shim.Column_String_{String_: col1Val}}
		col2 := shim.Column{Value: &shim.Column_String_{String_: col2Val}}
		col3 := shim.Column{Value: &shim.Column_String_{String_: col3Val}}
		col4 := shim.Column{Value: &shim.Column_String_{String_: col4Val}}
		col5 := shim.Column{Value: &shim.Column_String_{String_: col5Val}}
		col6 := shim.Column{Value: &shim.Column_String_{String_: col6Val}}
		col7 := shim.Column{Value: &shim.Column_String_{String_: col7Val}}
		col8 := shim.Column{Value: &shim.Column_String_{String_: col8Val}}
		col9 := shim.Column{Value: &shim.Column_String_{String_: col9Val}}
		col10 := shim.Column{Value: &shim.Column_String_{String_: col10Val}}
		col11 := shim.Column{Value: &shim.Column_String_{String_: col11Val}}
		col12 := shim.Column{Value: &shim.Column_String_{String_: col12Val}}
		col13 := shim.Column{Value: &shim.Column_String_{String_: col13Val}}
		col14 := shim.Column{Value: &shim.Column_String_{String_: col14Val}}
		col15 := shim.Column{Value: &shim.Column_String_{String_: col15Val}}
		col16 := shim.Column{Value: &shim.Column_String_{String_: col16Val}}
		col17 := shim.Column{Value: &shim.Column_String_{String_: col17Val}}
		col18 := shim.Column{Value: &shim.Column_String_{String_: col18Val}}
		columns = append(columns, &col1)
		columns = append(columns, &col2)
		columns = append(columns, &col3)
		columns = append(columns, &col4)
		columns = append(columns, &col5)
		columns = append(columns, &col6)
		columns = append(columns, &col7)
		columns = append(columns, &col8)
		columns = append(columns, &col9)
		columns = append(columns, &col10)
		columns = append(columns, &col11)
		columns = append(columns, &col12)
		columns = append(columns, &col13)
		columns = append(columns, &col14)
		columns = append(columns, &col15)
		columns = append(columns, &col16)
		columns = append(columns, &col17)
		columns = append(columns, &col18)

		row := shim.Row{Columns: columns}
		ok, err := stub.InsertRow("auction", row)
		if err != nil {
			return nil, fmt.Errorf("insertauction operation failed. %s", err)
		}
		if !ok {
			return nil, errors.New("insertauction operation failed. Row with given key already exists")
		}


	/*case "replaceRowAuction":
		if len(args) < 3 {
			return nil, errors.New("replaceRowAuction failed. Must include 3 column values")
		}

		col1Val := args[0]
		col2Int, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil {
			return nil, errors.New("replaceRowAuction failed. arg[1] must be convertable to int32")
		}
		col2Val := int32(col2Int)
		col3Int, err := strconv.ParseInt(args[2], 10, 32)
		if err != nil {
			return nil, errors.New("replaceRowAuction failed. arg[2] must be convertable to int32")
		}
		col3Val := int32(col3Int)

		var columns []*shim.Column
		col1 := shim.Column{Value: &shim.Column_String_{String_: col1Val}}
		col2 := shim.Column{Value: &shim.Column_Int32{Int32: col2Val}}
		col3 := shim.Column{Value: &shim.Column_Int32{Int32: col3Val}}
		columns = append(columns, &col1)
		columns = append(columns, &col2)
		columns = append(columns, &col3)

		row := shim.Row{Columns: columns}
		ok, err := stub.ReplaceRow("tableOne", row)
		if err != nil {
			return nil, fmt.Errorf("replaceRowTableOne operation failed. %s", err)
		}
		if !ok {
			return nil, errors.New("replaceRowTableOne operation failed. Row with given key does not exist")
		}*/

	case "deleteAndRecreateAuction":

		err := stub.DeleteTable("auction")
		if err != nil {
			return nil, fmt.Errorf("deleteAndRecreateAuction operation failed. Error deleting table. %s", err)
		}

		err = createTableAuction(stub)
		if err != nil {
			return nil, fmt.Errorf("deleteAndRecreateAuction operation failed. Error creating table. %s", err)
		}

		return nil, nil

	default:
		return nil, errors.New("Unsupported operation")
	}
	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {

	case "getRowAuction":
		if len(args) < 1 {
			return nil, errors.New("getRowAuction failed. Must include 1 key value")
		}

		col1Val := args[0]
		var columns []shim.Column
		col1 := shim.Column{Value: &shim.Column_String_{String_: col1Val}}
		columns = append(columns, col1)

		row, err := stub.GetRow("auction", columns)
		if err != nil {
			return nil, fmt.Errorf("getRowTableOne operation failed. %s", err)
		}

		rowString := fmt.Sprintf("%s", row)
		return []byte(rowString), nil

	default:
		return nil, errors.New("Unsupported operation")
	}
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func createTableAuction(stub shim.ChaincodeStubInterface) error {
	// Create table one
	var columnDefsTableOne []*shim.ColumnDefinition
	
	columnOneTableOneDef := shim.ColumnDefinition{Name: "auction ID",Type: shim.ColumnDefinition_STRING, Key: true}
	columnTwoTableOneDef := shim.ColumnDefinition{Name: "Consignor ID",Type: shim.ColumnDefinition_STRING, Key: false}
	columnThreeTableOneDef := shim.ColumnDefinition{Name: "Sale Year",Type: shim.ColumnDefinition_INT32, Key: false}
	columnFourTableOneDef := shim.ColumnDefinition{Name: "Sale Number",Type: shim.ColumnDefinition_INT32, Key: false}
	columnFiveTableOneDef := shim.ColumnDefinition{Name: "Lane Number",Type: shim.ColumnDefinition_INT32, Key: false}
	columnSixTableOneDef := shim.ColumnDefinition{Name: "RUN Number",Type: shim.ColumnDefinition_INT32, Key: false}
	columnSevenTableOneDef := shim.ColumnDefinition{Name: "Lease Account Number",Type: shim.ColumnDefinition_STRING, Key: false}
	columnEightTableOneDef := shim.ColumnDefinition{Name: "Work order Number",Type: shim.ColumnDefinition_INT32, Key: false}
	columnNineTableOneDef := shim.ColumnDefinition{Name: "Mileage",Type: shim.ColumnDefinition_INT32, Key: false}
	columnTenTableOneDef := shim.ColumnDefinition{Name: "Buyer ID",Type: shim.ColumnDefinition_INT32, Key: false}
	columnElevenTableOneDef := shim.ColumnDefinition{Name: "Buyer Name",Type: shim.ColumnDefinition_STRING, Key: false}
	columnTwelveTableOneDef := shim.ColumnDefinition{Name: "Sale Price",Type: shim.ColumnDefinition_INT64, Key: false}
	columnThirteenTableOneDef := shim.ColumnDefinition{Name: "Check Amt",Type: shim.ColumnDefinition_INT64, Key: false}
	columnFourteenTableOneDef := shim.ColumnDefinition{Name: "Payment Mode",Type: shim.ColumnDefinition_STRING, Key: false}
	columnFifteenTableOneDef := shim.ColumnDefinition{Name: "Title Status",Type: shim.ColumnDefinition_STRING, Key: false}
	columnSixteenTableOneDef := shim.ColumnDefinition{Name: "Vehicle Status",Type: shim.ColumnDefinition_STRING, Key: false}
	columnSeventeenTableOneDef := shim.ColumnDefinition{Name: "vehicle ID",Type: shim.ColumnDefinition_STRING, Key: true}
	columnDefsTableOne = append(columnDefsTableOne, &columnOneTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnTwoTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnThreeTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnFourTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnFiveTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnSixTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnSevenTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnEightTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnNineTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnTenTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnElevenTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnTwelveTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnThirteenTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnFourteenTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnFifteenTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnSixteenTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnSeventeenTableOneDef)
	
	return stub.CreateTable("auction", columnDefsTableOne)
}

