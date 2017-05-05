/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"time"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)


type SimpleChaincode struct {
}

var moneyIndexStr = "_moneyindex"				
var openTradesStr = "_opentrades"				

type Money struct{
	Name string `json:"name"`					
	Value string `json:"value"`
	Number int `json:"number"`
	User string `json:"user"`
}

type Description struct{
	Value string `json:"value"`
	Number int `json:"number"`
}

type AnOpenTrade struct{
	User string `json:"user"`					
	Timestamp int64 `json:"timestamp"`			
	Want Description  `json:"want"`			
	Willing []Description `json:"willing"`		
}

type AllTrades struct{
	OpenTrades []AnOpenTrade `json:"open_trades"`
}


func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}


func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var Aval int
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	
	err = stub.PutState("abc", []byte(strconv.Itoa(Aval)))				
	if err != nil {
		return nil, err
	}
	
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								
	err = stub.PutState(marbleIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	
	var trades AllTrades
	jsonAsBytes, _ = json.Marshal(trades)						
	err = stub.PutState(openTradesStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	
	return nil, nil
}


func (t *SimpleChaincode) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}


func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	
	if function == "init" {										
		return t.Init(stub, "init", args)
	} else if function == "delete" {								
		res, err := t.Delete(stub, args)
		cleanTrades(stub)													
		return res, err
	} else if function == "write" {										
		return t.Write(stub, args)
	} else if function == "init_money" {								
		return t.init_money(stub, args)
	} else if function == "set_user" {										
		res, err := t.set_user(stub, args)
		cleanTrades(stub)													
		return res, err
	} else if function == "open_trade" {									
		return t.open_trade(stub, args)
	} else if function == "perform_trade" {									
		res, err := t.perform_trade(stub, args)
		cleanTrades(stub)													
		return res, err
	} else if function == "remove_trade" {									
		return t.remove_trade(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)			

	return nil, errors.New("Received unknown function invocation")
}


func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	
	if function == "read" {												
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query")
}


func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name)									
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil													
}


func (t *SimpleChaincode) Delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	name := args[0]
	err := stub.DelState(name)													
	if err != nil 
		return nil, errors.New("Failed to delete state")
	}

	
	marblesAsBytes, err := stub.GetState(moneyIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get money index")
	}
	var moneyIndex []string
	json.Unmarshal(moneysAsBytes, &moneyIndex)								
	
	
	for i,val := range moneyIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + name)
		if val == name{															
			fmt.Println("found money")
			moneyIndex = append(moneyIndex[:i], moneyIndex[i+1:]...)	
			for x:= range moneyIndex{									
				fmt.Println(string(x) + " - " + moneyIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(moneyIndex)									
	err = stub.PutState(moneyIndexStr, jsonAsBytes)
	return nil, nil
}


func (t *SimpleChaincode) Write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, value string // Entities
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
	}

	name = args[0]															
	value = args[1]
	err = stub.PutState(name, []byte(value))								
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) init_money(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	
	
	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	
	fmt.Println("- start init money")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
	name := args[0]
	value := strings.ToLower(args[1])
	user := strings.ToLower(args[3])
	size, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("3rd argument must be a numeric string")
	}

	//check if marble already exists
	moneyAsBytes, err := stub.GetState(name)
	if err != nil {
		return nil, errors.New("Failed to get money name")
	}
	res := Money{}
	json.Unmarshal(moneyAsBytes, &res)
	if res.Name == name{
		fmt.Println("This money already exists: " + name)
		fmt.Println(res);
		return nil, errors.New("This money already exists")				
	}
	
	
	str := `{"name": "` + name + `", "value": "` + value + `", "number": ` + strconv.Itoa(size) + `, "user": "` + user + `"}`
	err = stub.PutState(name, []byte(str))									
	if err != nil {
		return nil, err
	}
		
	
	moneysAsBytes, err := stub.GetState(moneyIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get money index")
	}
	var moneyIndex []string
	json.Unmarshal(moneysAsBytes, &moneyIndex)							
	
	
	moneyIndex = append(moneyIndex, name)									
	fmt.Println("! money index: ", moneyIndex)
	jsonAsBytes, _ := json.Marshal(moneyIndex)
	err = stub.PutState(moneyIndexStr, jsonAsBytes)					

	fmt.Println("- end init money")
	return nil, nil
}


func (t *SimpleChaincode) set_user(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	
	
	
	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	
	fmt.Println("- start set user")
	fmt.Println(args[0] + " - " + args[1])
	marbleAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get thing")
	}
	res := Marble{}
	json.Unmarshal(marbleAsBytes, &res)										
	res.User = args[1]													
	
	jsonAsBytes, _ := json.Marshal(res)
	err = stub.PutState(args[0], jsonAsBytes)								
	if err != nil {
		return nil, err
	}
	
	fmt.Println("- end set user")
	return nil, nil
}


func (t *SimpleChaincode) open_trade(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var will_size int
	var trade_away Description
	
	
	
	if len(args) < 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting like 5?")
	}
	if len(args)%2 == 0{
		return nil, errors.New("Incorrect number of arguments. Expecting an odd number")
	}

	size1, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("3rd argument must be a numeric string")
	}

	open := AnOpenTrade{}
	open.User = args[0]
	open.Timestamp = makeTimestamp()										
	open.Want.Value = args[1]
	open.Want.Number =  size1
	fmt.Println("- start open trade")
	jsonAsBytes, _ := json.Marshal(open)
	err = stub.PutState("_debug1", jsonAsBytes)

	for i:=3; i < len(args); i++ {												
		will_size, err = strconv.Atoi(args[i + 1])
		if err != nil {
			msg := "is not a numeric string " + args[i + 1]
			fmt.Println(msg)
			return nil, errors.New(msg)
		}
		
		trade_away = Description{}
		trade_away.Color = args[i]
		trade_away.Size =  will_size
		fmt.Println("! created trade_away: " + args[i])
		jsonAsBytes, _ = json.Marshal(trade_away)
		err = stub.PutState("_debug2", jsonAsBytes)
		
		open.Willing = append(open.Willing, trade_away)
		fmt.Println("! appended willing to open")
		i++;
	}
	
	
	tradesAsBytes, err := stub.GetState(openTradesStr)
	if err != nil {
		return nil, errors.New("Failed to get opentrades")
	}
	var trades AllTrades
	json.Unmarshal(tradesAsBytes, &trades)										
	
	trades.OpenTrades = append(trades.OpenTrades, open);				
	fmt.Println("! appended open to trades")
	jsonAsBytes, _ = json.Marshal(trades)
	err = stub.PutState(openTradesStr, jsonAsBytes)					
	if err != nil {
		return nil, err
	}
	fmt.Println("- end open trade")
	return nil, nil
}


func (t *SimpleChaincode) perform_trade(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	
	
	
	if len(args) < 6 {
		return nil, errors.New("Incorrect number of arguments. Expecting 6")
	}
	
	fmt.Println("- start close trade")
	timestamp, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return nil, errors.New("1st argument must be a numeric string")
	}
	
	size, err := strconv.Atoi(args[5])
	if err != nil {
		return nil, errors.New("6th argument must be a numeric string")
	}
	
	
	tradesAsBytes, err := stub.GetState(openTradesStr)
	if err != nil {
		return nil, errors.New("Failed to get opentrades")
	}
	var trades AllTrades
	json.Unmarshal(tradesAsBytes, &trades)															
	
	for i := range trades.OpenTrades{															
		fmt.Println("looking at " + strconv.FormatInt(trades.OpenTrades[i].Timestamp, 10) + " for " + strconv.FormatInt(timestamp, 10))
		if trades.OpenTrades[i].Timestamp == timestamp{
			fmt.Println("found the trade");
			
			
			moneyAsBytes, err := stub.GetState(args[2])
			if err != nil {
				return nil, errors.New("Failed to get thing")
			}
			closersMoney := Money{}
			json.Unmarshal(moneyAsBytes, &closersMoney)									
			
			
			if closersMarble.Value != trades.OpenTrades[i].Want.Value || closersMoney.Number != trades.OpenTrades[i].Want.Number {
				msg := "money in input does not meet trade requriements"
				fmt.Println(msg)
				return nil, errors.New(msg)
			}
			
			marble, e := findMoney4Trade(stub, trades.OpenTrades[i].User, args[4], number)		
			if(e == nil){
				fmt.Println("! no errors, proceeding")

				t.set_user(stub, []string{args[2], trades.OpenTrades[i].User})						
				t.set_user(stub, []string{money.Name, args[1]})						
			
				trades.OpenTrades = append(trades.OpenTrades[:i], trades.OpenTrades[i+1:]...)	
				jsonAsBytes, _ := json.Marshal(trades)
				err = stub.PutState(openTradesStr, jsonAsBytes)						
				if err != nil {
					return nil, err
				}
			}
		}
	}
	fmt.Println("- end close trade")
	return nil, nil
}


func findMarble4Trade(stub shim.ChaincodeStubInterface, user string, value string, number int )(m Money, err error){
	var fail Marble;
	fmt.Println("- start find money 4 trade")
	fmt.Println("looking for " + user + ", " + value + ", " + strconv.Itoa(number));


	moneysAsBytes, err := stub.GetState(moneyIndexStr)
	if err != nil {
		return fail, errors.New("Failed to get money index")
	}
	var moneyIndex []string
	json.Unmarshal(moneysAsBytes, &moneyIndex)								
	
	for i:= range moneyIndex{												
		

		moneyAsBytes, err := stub.GetState(moneyIndex[i])						
		if err != nil {
			return fail, errors.New("Failed to get money")
		}
		res := Money{}
		json.Unmarshal(moneyAsBytes, &res)										
		
		
	
		if strings.ToLower(res.User) == strings.ToLower(user) && strings.ToLower(res.Value) == strings.ToLower(value) && res.Number == number{
			fmt.Println("found a money: " + res.Name)
			fmt.Println("! end find money 4 trade")
			return res, nil
		}
	}
	
	fmt.Println("- end find money 4 trade - error")
	return fail, errors.New("Did not find money to use in this trade")
}


func makeTimestamp() int64 {
    return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}


func (t *SimpleChaincode) remove_trade(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	
	
	if len(args) < 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	fmt.Println("- start remove trade")
	timestamp, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return nil, errors.New("1st argument must be a numeric string")
	}
	
	//get the open trade struct
	tradesAsBytes, err := stub.GetState(openTradesStr)
	if err != nil {
		return nil, errors.New("Failed to get opentrades")
	}
	var trades AllTrades
	json.Unmarshal(tradesAsBytes, &trades)															
	
	for i := range trades.OpenTrades{																	
		
		if trades.OpenTrades[i].Timestamp == timestamp{
			fmt.Println("found the trade");
			trades.OpenTrades = append(trades.OpenTrades[:i], trades.OpenTrades[i+1:]...)		
			jsonAsBytes, _ := json.Marshal(trades)
			err = stub.PutState(openTradesStr, jsonAsBytes)											
			if err != nil {
				return nil, err
			}
			break
		}
	}
	
	fmt.Println("- end remove trade")
	return nil, nil
}


func cleanTrades(stub shim.ChaincodeStubInterface)(err error){
	var didWork = false
	fmt.Println("- start clean trades")
	
	
	tradesAsBytes, err := stub.GetState(openTradesStr)
	if err != nil {
		return errors.New("Failed to get opentrades")
	}
	var trades AllTrades
	json.Unmarshal(tradesAsBytes, &trades)															
	
	fmt.Println("# trades " + strconv.Itoa(len(trades.OpenTrades)))
	for i:=0; i<len(trades.OpenTrades); {																	
		fmt.Println(strconv.Itoa(i) + ": looking at trade " + strconv.FormatInt(trades.OpenTrades[i].Timestamp, 10))
		
		fmt.Println("# options " + strconv.Itoa(len(trades.OpenTrades[i].Willing)))
		for x:=0; x<len(trades.OpenTrades[i].Willing); {														
			fmt.Println("! on next option " + strconv.Itoa(i) + ":" + strconv.Itoa(x))
			_, e := findMarble4Trade(stub, trades.OpenTrades[i].User, trades.OpenTrades[i].Willing[x].Value, trades.OpenTrades[i].Willing[x].Number)
			if(e != nil){
				fmt.Println("! errors with this option, removing option")
				didWork = true
				trades.OpenTrades[i].Willing = append(trades.OpenTrades[i].Willing[:x], trades.OpenTrades[i].Willing[x+1:]...)	
				x--;
			}else{
				fmt.Println("! this option is fine")
			}
			
			x++
			fmt.Println("! x:" + strconv.Itoa(x))
			if x >= len(trades.OpenTrades[i].Willing) {												
				break
			}
		}
		
		if len(trades.OpenTrades[i].Willing) == 0 {
			fmt.Println("! no more options for this trade, removing trade")
			didWork = true
			trades.OpenTrades = append(trades.OpenTrades[:i], trades.OpenTrades[i+1:]...
			i--;
		}
		
		i++
		fmt.Println("! i:" + strconv.Itoa(i))
		if i >= len(trades.OpenTrades) {															
			break
		}
	}

	if(didWork){
		fmt.Println("! saving open trade changes")
		jsonAsBytes, _ := json.Marshal(trades)
		err = stub.PutState(openTradesStr, jsonAsBytes)													
		if err != nil {
			return err
		}
	}else{
		fmt.Println("! all open trades are fine")
	}

	fmt.Println("- end clean trades")
	return nil
}