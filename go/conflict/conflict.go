package main

import (
	// "bytes"
	"math"
	"encoding/json"
	"fmt"
	"time"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)


type SmartContract struct {
}

const(
	ConflictPrefix = "CONF_"
	ResultPrefix = "RES_"
	N = 5
)

type Conflict struct {
	Conflictname  string `json:"conflictname"`
	User  string `json:"user"`
	Service string `json:"service"`
	Conflictcontents string `json:"conflictcontents"`
	Conflicttime string `json:"conflicttime"`
}

type Result struct {
	Resultname string `json:"resultname"`
	User string `json:"user"`
	Service string `json:"service"`
	Resultcontents string `json:"resultcontents"`
	Trialtime string `json:"trialtime"`
	Winner string `json:"winner"`
}

type Judger struct{
	Name string
	Faith string
	Result string
	Content string
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}


func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryConflict" {
		return s.queryConflict(APIstub, args)
	} else if function == "happenConflict" {
		return s.happenConflict(APIstub, args)
	} else if function == "queryResult" {
		return s.queryResult(APIstub, args)
	} else if function == "makeTrial" {
		return s.makeTrial(APIstub, args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryConflict(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	conflictAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(conflictAsBytes)
}

func (s *SmartContract) queryResult(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	conflictAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(conflictAsBytes)
}

func (s *SmartContract) happenConflict(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	key := ConflictPrefix + args[0]
    tNow := time.Now()
    tString := tNow.UTC().Format(time.UnixDate)
	var conflict = Conflict{Conflictname: args[0], User: args[1], Service: args[2], Conflictcontents: args[3], Conflicttime: tString}

	conflictAsBytes, _ := json.Marshal(conflict)
	APIstub.PutState(key, conflictAsBytes)

	return shim.Success(nil)
}



func (s *SmartContract) makeTrial(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 * N + 3 {
		return shim.Error("Incorrect number of arguments. Expecting 9")
	}
	var judger [N]Judger
	var faith [N]int
	var results [N]int
	key := ResultPrefix + args[0]
	user := args[1]
	service := args[2]
	winner := user
	vote := float64(0)
	for i := 0; i < N; i++ {
             judger[i] = Judger{Name:args[4 * i + 3], Faith:args[4 * i + 4], Result:args[4 * i + 5], Content:args[4 * i + 6]}
             faith[i], _= strconv.Atoi(args[4 * i + 4])
             if args[4 * i + 5] == user{
             	results[i] = 1
             } else{
             	results[i] = -1
             }

        }
    for i := 0; i < N; i++ {
    	vote = vote + math.Log2(float64(faith[i] + 1)) * float64(results[i])
    }
	if vote < 0{
		winner = service
	}
	result_contents := ""
	for i := 0; i < N; i++ {
    	result_contents = result_contents + "{ " + judger[i].Name + ", " + judger[i].Faith + ": " + judger[i].Content + "; " + judger[i].Result + "}  "
    }
    tNow := time.Now()
    tString := tNow.UTC().Format(time.UnixDate)
	var result = Result{Resultname: args[0], User: user, Service: service, Resultcontents: result_contents, Winner: winner, Trialtime: tString}

	resultAsBytes, _ := json.Marshal(result)
	APIstub.PutState(key, resultAsBytes)

	return shim.Success(nil)
}


// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
