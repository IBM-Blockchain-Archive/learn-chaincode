

package main

import (
	//"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Init create tables for tests
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Create table one
	err := createTableOne(stub)
	if err != nil {
		return nil, fmt.Errorf("Error creating table one during init. %s", err)
	}

	

	return nil, nil
}


func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	switch function {

	case "insertRowTableOne":
		if len(args) < 6 {
			return nil, errors.New("insertTableOne failed. Must include 6 column values")
		}

		//col1Val := args[0]
		//col2Int, err := strconv.ParseInt(args[1], "10", "32")
		//if err != nil {
		//	return nil, errors.New("insertTableOne failed. arg[1] must be convertable to int32")
		//}
		//col2Val := int32(col2Int)
//		col3Int, err := strconv.ParseInt(args[2], "10", "32")
		//if err != nil {
		//	return nil, errors.New("insertTableOne failed. arg[2] must be convertable to int32")
		//}
		//col3Val := int32(col3Int)

		//COMMENTED FOR hard coded values
		col1Val := args[0]
		col2Val := args[1]
		col3Val := args[2]
		col4Val := args[3]
		col5Val := args[4]
		col6Val := args[5]

		


		var columns []*shim.Column
		//col1 := shim.Column{Value: &shim.Column_String_{String_: col1Val}}
		
		col1 := shim.Column{Value: &shim.Column_String_{String_: col1Val}}
		col2 := shim.Column{Value: &shim.Column_String_{String_: col2Val}}
		col3 := shim.Column{Value: &shim.Column_String_{String_: col3Val}}
		col4 := shim.Column{Value: &shim.Column_String_{String_: col4Val}}
		col5 := shim.Column{Value: &shim.Column_String_{String_: col5Val}}
		col6 := shim.Column{Value: &shim.Column_String_{String_: col6Val}}
		columns = append(columns, &col1)
		columns = append(columns, &col2)
		columns = append(columns, &col3)
		columns = append(columns, &col4)
		columns = append(columns, &col5)
		columns = append(columns, &col6)

		row := shim.Row{Columns: columns}
		ok, err := stub.InsertRow("tableOne", row)
		if err != nil {
			return nil, fmt.Errorf("insertTableOne operation failed. %s", err)
		}
		if !ok {
			return nil, errors.New("insertTableOne operation failed. Row with given key already exists")
		}


	case "replaceRowTableOne":
		if len(args) < 3 {
			return nil, errors.New("replaceRowTableOne failed. Must include 3 column values")
		}

		col1Val := args[0]
		col2Int, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil {
			return nil, errors.New("replaceRowTableOne failed. arg[1] must be convertable to int32")
		}
		col2Val := int32(col2Int)
		col3Int, err := strconv.ParseInt(args[2], 10, 32)
		if err != nil {
			return nil, errors.New("replaceRowTableOne failed. arg[2] must be convertable to int32")
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
		}

	case "deleteAndRecreateTableOne":

		err := stub.DeleteTable("tableOne")
		if err != nil {
			return nil, fmt.Errorf("deleteAndRecreateTableOne operation failed. Error deleting table. %s", err)
		}

		err = createTableOne(stub)
		if err != nil {
			return nil, fmt.Errorf("deleteAndRecreateTableOne operation failed. Error creating table. %s", err)
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

	case "getRowTableOne":
		if len(args) < 1 {
			return nil, errors.New("getRowTableOne failed. Must include 1 key value")
		}

		col1Val := args[0]
		var columns []shim.Column
		col1 := shim.Column{Value: &shim.Column_String_{String_: col1Val}}
		columns = append(columns, col1)

		row, err := stub.GetRow("tableOne", columns)
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

func createTableOne(stub shim.ChaincodeStubInterface) error {
	// Create table one
	var columnDefsTableOne []*shim.ColumnDefinition
	columnOneTableOneDef := shim.ColumnDefinition{Name: "Vehicle_VIN",
		Type: shim.ColumnDefinition_STRING, Key: true}
	columnTwoTableOneDef := shim.ColumnDefinition{Name: "Vehicle_Make",
		Type: shim.ColumnDefinition_STRING, Key: false}
	columnThreeTableOneDef := shim.ColumnDefinition{Name: "Vehicle_Model",
		Type: shim.ColumnDefinition_STRING, Key: false}
	columnFourTableOneDef := shim.ColumnDefinition{Name: "Vehicle_Colour",
		Type: shim.ColumnDefinition_STRING, Key: false}
	columnFiveTableOneDef := shim.ColumnDefinition{Name: "Vehicle_Reg",
		Type: shim.ColumnDefinition_STRING, Key: false}
	columnsSixTableOneDef := shim.ColumnDefinition{Name: "Vehicle_Owner",
		Type: shim.ColumnDefinition_STRING, Key: false}
	columnDefsTableOne = append(columnDefsTableOne, &columnOneTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnTwoTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnThreeTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnFourTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnFiveTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnsSixTableOneDef)
	return stub.CreateTable("tableOne", columnDefsTableOne)
}

