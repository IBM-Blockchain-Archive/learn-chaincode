/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"regexp"
	"strconv"
	// "strings"
)

const (
	Pattern_Author       = "^[\u4e00-\u9fa5_a-zA-Z, ]+$"
	Pattern_Title        = "^[\u4e00-\u9fa5_a-zA-Z0-9\\[\\]\\(\\)\":\\- ]+$"
	Pattern_Journal      = "^[\u4e00-\u9fa5_a-zA-Z ]+$"
	Pattern_Year         = "^\\d{4}$"
	Pattern_Volume_Issue = "^\\d+\\(\\d+\\)$"
	Pattern_Page         = "^(\\d+)-(\\d+)$|^\\d+$"
)

// ArchiveChaincode example simple Chaincode implementation
type ArchiveChaincode struct {
	total uint64
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(ArchiveChaincode))
	if err != nil {
		fmt.Printf("Error starting Archive chaincode: %s", err)
	}
}

// Init resets all the things
func (t *ArchiveChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function != "init" {
		return nil, errors.New("error function name!")
	}

	t.total = 0
	fmt.Printf("Init called success!")
	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *ArchiveChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// append item of archive
	if function == "append" {
		// 1.Author, 2.Title, 3.Journal, 4.Year, 5.Volume(Issue), 6.Page
		if len(args) < 4 {
			fmt.Println("[error] missing parameter")
			return nil, errors.New("missing parameter")
		}
		if len(args) > 6 {
			fmt.Println("[error] too many parameter")
			return nil, errors.New("too many parameter")
		}
		// define a regexp
		var reg *regexp.Regexp

		author := args[0]
		reg = regexp.MustCompile(Pattern_Author)
		if reg.MatchString(author) != true {
			fmt.Println("[error] format error by author")
			return nil, errors.New("format error by author")
		}

		title := args[1]
		reg = regexp.MustCompile(Pattern_Title)
		if reg.MatchString(title) != true {
			fmt.Println("[error] format error by title")
			return nil, errors.New("format error by title")
		}

		journal := args[2]
		reg = regexp.MustCompile(Pattern_Journal)
		if reg.MatchString(journal) != true {
			fmt.Println("[error] format error by journal")
			return nil, errors.New("format error by journal")
		}

		year := args[3]
		reg = regexp.MustCompile(Pattern_Year)
		if reg.MatchString(year) != true {
			fmt.Println("[error] format error by year")
			return nil, errors.New("format error by year")
		}

		if len(args) == 5 {
			page := args[4]
			reg = regexp.MustCompile(Pattern_Page)
			str := reg.FindStringSubmatch(page)
			if len(str) == 0 {
				fmt.Println("[error] format error by page")
				return nil, errors.New("format error by page")
			}
			p1, _ := strconv.Atoi(str[1])
			p2, _ := strconv.Atoi(str[2])
			if p1 > p2 {
				fmt.Println("[error] page must be from min to max")
				return nil, errors.New("page must be from min to max")
			}
		}

		if len(args) == 6 {
			volume_issue := args[4]
			reg = regexp.MustCompile(Pattern_Volume_Issue)
			if reg.MatchString(volume_issue) != true {
				fmt.Println("[error] format error by volume_issue")
				return nil, errors.New("format error by volume_issue")
			}

			page := args[5]
			reg = regexp.MustCompile(Pattern_Page)
			str := reg.FindStringSubmatch(page)
			if len(str) == 0 {
				fmt.Println("[error] format error by page")
				return nil, errors.New("format error by page")
			}
			p1, _ := strconv.Atoi(str[1])
			p2, _ := strconv.Atoi(str[2])
			if p1 > p2 {
				fmt.Println("[error] page must be from min to max")
				return nil, errors.New("page must be from min to max")
			}
		}

		t.total += 1
		fmt.Printf("append: %d\n", t.total)
		return nil, nil
	}

	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *ArchiveChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// get the total number of archives
	if function == "total" {
		fmt.Printf("total: %d\n", t.total)
		str := strconv.FormatUint(t.total, 10)
		return []byte(str), nil
	}

	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query: " + function)
}
