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
	Prefix = "COIN_"
)

type Travelcoin struct {
	Name  string `json:"name"`
	Balance  string `json:"balance"`
}


func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}


func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryTravelcoin" {
		return s.queryTravelcoin(APIstub, args)
	} else if function == "initTravelcoin" {
		return s.initTravelcoin(APIstub, args)
	} else if function == "changeTravelcoin" {
		return s.changeTravelcoin(APIstub, args)
	} else if function == "transferTravelcoin"{
		return s.transferTravelcoin(APIstub, args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryTravelcoin(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	travelcoinAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(travelcoinAsBytes)
}

func (s *SmartContract) initTravelcoin(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
    key := Prefix + args[0]
    value := 100
	var travelcoin = Travelcoin{Name: args[0], Balance: strconv.Itoa(value)}

	travelcoinAsBytes,_ := json.Marshal(travelcoin)
	APIstub.PutState(key, travelcoinAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) changeTravelcoin(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	travelcoinAsBytes,_ := APIstub.GetState(args[0])
	travelcoin := Travelcoin{}
	json.Unmarshal(travelcoinAsBytes, &travelcoin)
	value,_ := strconv.Atoi(travelcoin.Balance)
	changevalue,_ := strconv.Atoi(args[2])
	new_value := 0
	switch args[1]{
	case "increase":
		new_value = value + changevalue
	case "decrease":
		new_value = value - changevalue
	}
	travelcoin.Balance = strconv.Itoa(new_value)
	travelcoinAsBytes,_ = json.Marshal(travelcoin)
	APIstub.PutState(args[0], travelcoinAsBytes)
	return shim.Success(nil)
}

func (s *SmartContract) transferTravelcoin(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	A_travelcoinAsBytes,_ := APIstub.GetState(args[0])
	A_travelcoin := Travelcoin{}
	json.Unmarshal(A_travelcoinAsBytes, &A_travelcoin)
	B_travelcoinAsBytes,_ := APIstub.GetState(args[1])
	B_travelcoin := Travelcoin{}
	json.Unmarshal(B_travelcoinAsBytes, &B_travelcoin)
	A_value,_ := strconv.Atoi(A_travelcoin.Balance)
	B_value,_ := strconv.Atoi(B_travelcoin.Balance)
	changevalue,_ := strconv.Atoi(args[2])
	A_value_new := A_value - changevalue
	B_value_new := B_value + changevalue
	A_travelcoin.Balance = strconv.Itoa(A_value_new)
	B_travelcoin.Balance = strconv.Itoa(B_value_new)
	A_travelcoinAsBytes,_ = json.Marshal(A_travelcoin)
	APIstub.PutState(args[0], A_travelcoinAsBytes)
	B_travelcoinAsBytes,_ = json.Marshal(B_travelcoin)
	APIstub.PutState(args[1], B_travelcoinAsBytes)
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
