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
	AdPrefix = "AD_"
	PurchasePrefix = "PURCHASE_"
)

type Ad struct {
	Publisher string `json:"publisher"`
	Adname  string `json:"adname"`
	Adcontent string `json:"adcontent"`
	Token string `json:"token"`
	Adtime string `json:"adtime"`
}

type Purchase struct {
	Purchasename string `json:"purchasename"`
	Purchaser string `json:"purchaser"`
	Seller  string `json:"seller"`
	Purchasecontent string `json:"purchasecontent"`
	Token string `json:"token"`
	Purchasetime string `json:"purchasetime"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}


func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryCommercial" {
		return s.queryCommercial(APIstub, args)
	} else if function == "publishAd" {
		return s.publishAd(APIstub, args)
	} else if function == "purchaseData" {
		return s.purchaseData(APIstub, args)
	} 
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryCommercial(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	commercialAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(commercialAsBytes)
}

func (s *SmartContract) publishAd(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	key := AdPrefix + args[1]
    tNow := time.Now()
    tString := tNow.UTC().Format(time.UnixDate)
	var ad = Ad{Publisher: args[0], Adname: args[1], Adcontent: args[2], Token: args[3], Adtime: tString}

	adAsBytes, _ := json.Marshal(ad)
	APIstub.PutState(key, adAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) purchaseData(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	key := PurchasePrefix + args[0]
    tNow := time.Now()
    tString := tNow.UTC().Format(time.UnixDate)
	var purchase = Purchase{Purchasename: args[0], Purchaser: args[1], Seller: args[2], Purchasecontent: args[3], Token: args[4], Purchasetime: tString}

	purchaseAsBytes, _ := json.Marshal(purchase)
	APIstub.PutState(key, purchaseAsBytes)

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
