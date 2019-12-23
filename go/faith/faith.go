package main

import (
	// "bytes"
	"encoding/json"
	"fmt"
	// "time"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)


type SmartContract struct {
}

const(
	Prefix = "FAITH_"
)

type Faith struct {
	Name  string `json:"name"`
	Faithvalue  string `json:"faithvalue"`
}


func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}


func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryFaith" {
		return s.queryFaith(APIstub, args)
	} else if function == "initFaith" {
		return s.initFaith(APIstub, args)
	} else if function == "changeFaith" {
		return s.changeFaith(APIstub, args)
	} 
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryFaith(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	faithAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(faithAsBytes)
}

func (s *SmartContract) initFaith(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
    key := Prefix + args[0]
    value := 100
	var faith = Faith{Name: args[0], Faithvalue: strconv.Itoa(value)}

	faithAsBytes,_ := json.Marshal(faith)
	APIstub.PutState(key, faithAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) changeFaith(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	faithAsBytes,_ := APIstub.GetState(args[0])
	faith := Faith{}
	json.Unmarshal(faithAsBytes, &faith)
	value,_ := strconv.Atoi(faith.Faithvalue)
	changevalue,_ := strconv.Atoi(args[2])
	new_value := 0
	switch args[1]{
	case "increase":
		new_value = value + changevalue
	case "decrease":
		new_value = value - changevalue
	}
	faith.Faithvalue = strconv.Itoa(new_value)
	faithAsBytes,_ = json.Marshal(faith)
	APIstub.PutState(args[0], faithAsBytes)
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
