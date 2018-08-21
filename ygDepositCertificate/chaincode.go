package ygdepositcertificate

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type YGDCChaincode struct {
}

type ygdc struct {
	Owner string `json:"owner"`
	Data  string `json:"Data"`
}

func (t *YGDCChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *YGDCChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("YGDC Invoke")
	function, args := stub.GetFunctionAndParameters()

	if function == "invoke" {
		return t.invoke(stub, args)
	} else if function == "query" {
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"query\"")
}

func (t *YGDCChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var Owner string
	var Desc string
	var Data string
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	Owner = args[0]
	Desc = args[1]
	Data = args[2]

	key := tomd5(Owner + Desc)
	_, err = stub.GetState(key)
	if err != nil {
		return shim.Error("data repeat")
	}

	dataToTransfer := ygdc{}
	dataToTransfer.Owner = Owner
	dataToTransfer.Data = Data

	dataJSONasBytes, _ := json.Marshal(dataToTransfer)
	err = stub.PutState(key, dataJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *YGDCChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var Owner string
	var Desc string
	var dataJson ygdc
	var jsonResp string
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	Owner = args[0]
	Desc = args[1]

	dataJSONasBytes, err := stub.GetState(tomd5(Owner + Desc))
	if err != nil {
		return shim.Error("Failed to get deposit certificate")
	}

	err = json.Unmarshal([]byte(dataJSONasBytes), &dataJson)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + Owner + Desc + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success([]byte(dataJson.Data))
}

func tomd5(data string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(data))
	res := hex.EncodeToString(md5Ctx.Sum(nil))
	return res
}

func main() {
	err := shim.Start(new(YGDCChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
