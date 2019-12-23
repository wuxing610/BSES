package main

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"time"
	// "strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)


type SmartContract struct {
}

const(
	Prefix = "EVENT_"
)

type Event struct {
	Eventname  string `json:"eventname"`
	User  string `json:"user"`
	Service string `json:"service"`
	Eventcontents string `json:"eventcontents"`
	Eventtime string `json:"eventtime"`
}


func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}


func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryEvent" {
		return s.queryEvent(APIstub, args)
	} else if function == "happenEvent" {
		return s.happenEvent(APIstub, args)
	} 
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryEvent(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	eventAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(eventAsBytes)
}

func (s *SmartContract) happenEvent(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	key := Prefix + args[0]
    tNow := time.Now()
    tString := tNow.UTC().Format(time.UnixDate)
	var event = Event{Eventname: args[0], User: args[1], Service: args[2], Eventcontents: args[3], Eventtime: tString}

	eventAsBytes, _ := json.Marshal(event)
	APIstub.PutState(key, eventAsBytes)

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
