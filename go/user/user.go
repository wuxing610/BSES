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
	ConsumerPrefix = "CON_"
	ProviderPrefix = "PRO_"
	ComposerPrefix = "COM_"
)

type User struct {
	Name  string `json:"name"`
	Introduction  string `json:"introduction"`
	Registrationtime string `json:"registrationtime"`
	Kind string `json:"kind"`
}


func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}


func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryUser" {
		return s.queryUser(APIstub, args)
	} else if function == "registerUser" {
		return s.registerUser(APIstub, args)
	} 
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	userAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(userAsBytes)
}

func (s *SmartContract) registerUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
    kind := args[2]
    key := " "
    switch kind{
    case "Consumer":
    	key = ConsumerPrefix + args[0]
    case "Provider":
    	key = ProviderPrefix + args[0]
    case "Composer":
    	key = ComposerPrefix + args[0]
    }
    tNow := time.Now()
    tString := tNow.UTC().Format(time.UnixDate)
	var user = User{Name: args[0], Introduction: args[1], Kind: args[2], Registrationtime: tString}

	userAsBytes, _ := json.Marshal(user)
	APIstub.PutState(key, userAsBytes)

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
