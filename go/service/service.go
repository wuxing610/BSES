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
	MushupPrefix = "MUSHUP_"
	SpotPrefix = "SPOT_"
	HotelPrefix = "HOT_"
	RestPrefix = "REST_"
)

type Service struct {
	Name  string `json:"name"`
	Description  string `json:"description"`
	Provider string `json:"provider"`
	Kind string `json:"kind"`
	Publicationtime string `json:"publicationtime"`
	Makeupservices string `json:"makeupservices"`
}


func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}


func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryService" {
		return s.queryService(APIstub, args)
	} else if function == "createService" {
		return s.createService(APIstub, args)
	} else if function == "createMushup" {
		return s.createMushup(APIstub, args)
	} 
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryService(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	serviceAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(serviceAsBytes)
}

func (s *SmartContract) createService(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
    kind := args[3]
    key := " "
    switch kind{
    case "Spot":
    	key = SpotPrefix + args[0]
    case "Hotel":
    	key = HotelPrefix + args[0]
    case "Restaurant":
    	key = RestPrefix + args[0]
    case "Mushup":
    	key = MushupPrefix + args[0]
    }
    tNow := time.Now()
    tString := tNow.UTC().Format(time.UnixDate)
	var service = Service{Name: args[0], Description: args[1], Provider: args[2], Kind: args[3],Publicationtime: tString}

	serviceAsBytes, _ := json.Marshal(service)
	APIstub.PutState(key, serviceAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) createMushup(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
    kind := args[3]
    key := " "
    switch kind{
    case "Spot":
    	key = SpotPrefix + args[0]
    case "Hotel":
    	key = HotelPrefix + args[0]
    case "Restaurant":
    	key = RestPrefix + args[0]
    case "Mushup":
    	key = MushupPrefix + args[0]
    }
    tNow := time.Now()
    tString := tNow.UTC().Format(time.UnixDate)
	var service = Service{Name: args[0], Description: args[1], Provider: args[2], Kind: args[3], Makeupservices: args[4],Publicationtime: tString}

	serviceAsBytes, _ := json.Marshal(service)
	APIstub.PutState(key, serviceAsBytes)

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
