package ygDigitalAssets

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type YGDAChaincode struct {
}

type digitalassets struct {
	AssetsName  string `json:name`
	AssetsType  string `json:type`
	AssetsValue string `json:value`
	AssetsDesc  string `json:desc`
	Owner       string `json:owner`
}

func (t *YGDAChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *YGDAChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("YGDA Invoke")
	function, args := stub.GetFunctionAndParameters()

	if function == "initAssets" {
		return t.initAssets(stub, args)
	} else if function == "queryAssets" {
		return t.queryAssets(stub, args)
	} else if function == "transferAssets" {
		return t.transferAssets(stub, args)
	} else if function == "modifyAssets" {
		return t.modifyAssets(stub, args)
	} else if function == "deleteAssets" {
		return t.deleteAssets(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"query\" \"addAssets\" \"transferAssets\" \"modifyAssets\" \"deleteAssets\" ")
}

func (t *YGDAChaincode) queryAssets(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	return shim.Success(nil)
}

func (t *YGDAChaincode) initAssets(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var UserName string
	var Type string
	var Value string
	var Desc string
	var Name string
	var err error

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	UserName = args[0]
	Name = args[1]
	Type = args[2]
	Value = args[3]
	Desc = args[4]

	return shim.Success(nil)
}

func (t *YGDAChaincode) transferAssets(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	return shim.Success(nil)
}

func (t *YGDAChaincode) modifyAssets(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	return shim.Success(nil)
}

func (t *YGDAChaincode) deleteAssets(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(YGDAChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
