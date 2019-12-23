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
	Prefix = "COMMENT_"
)

type Comment struct {
	Commentname  string `json:"commentname"`
	User  string `json:"user"`
	Service string `json:"service"`
	Contents string `json:"contents"`
	Score string `json:"score"`
	Commenttime string `json:"commenttime"`
}


func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}


func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryComment" {
		return s.queryComment(APIstub, args)
	} else if function == "makeComment" {
		return s.makeComment(APIstub, args)
	} 
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryComment(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	commentAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(commentAsBytes)
}

func (s *SmartContract) makeComment(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	key := Prefix + args[0]
    tNow := time.Now()
    tString := tNow.UTC().Format(time.UnixDate)
	var comment = Comment{Commentname: args[0], User: args[1], Service: args[2], Contents: args[3], Score: args[4], Commenttime: tString}

	commentAsBytes, _ := json.Marshal(comment)
	APIstub.PutState(key, commentAsBytes)

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
