/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ygtoken

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// YGTokenChaincode example simple Chaincode implementation
type YGTokenChaincode struct {
}

func (t *YGTokenChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *YGTokenChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ygtoken Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "register" {
		return t.register(stub, args)
	} else if function == "recharge" {
		return t.recharge(stub, args)
	} else if function == "query" {
		return t.query(stub, args)
	} else if function == "transfer" {
		return t.transfer(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"recharge\" \"delete\" \"query\"")
}

// register a new account
func (t *YGTokenChaincode) register(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var username string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	username = args[0]

	_, err = stub.GetState(username)
	if err == nil {
		return shim.Error("username repeat")
	}

	err = stub.PutState(username, []byte(strconv.Itoa(0)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// delete an account
//func (t *YGTokenChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
//	return shim.Success(nil)
//}

// query an account balance
func (t *YGTokenChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var username string

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	username = args[0]

	// Get the state from the ledger
	value, err := stub.GetState(username)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + username + "\"}"
		return shim.Error(jsonResp)
	}

	if value == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + username + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + username + "\",\"Amount\":\"" + string(value) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)

	return shim.Success(value)
}

// recharge an account
func (t *YGTokenChaincode) recharge(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var username string
	var value int

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	username = args[0]

	usernamevalbytes, err := stub.GetState(username)
	if err != nil {
		return shim.Error("Failed to get state!")
	}
	if usernamevalbytes == nil {
		return shim.Error("Entity not found")
	}

	usernameval, _ := strconv.Atoi(string(usernamevalbytes))

	value, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}

	usernameval = usernameval + value
	fmt.Printf("usernameval = %d\n", usernameval)

	err = stub.PutState(username, []byte(strconv.Itoa(usernameval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *YGTokenChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var usernameA string
	var usernameB string
	var Aval, Bval int
	var Val int
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	usernameA = args[0]
	usernameB = args[1]

	Avalbytes, err := stub.GetState(usernameA)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(usernameB)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	Val, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	Aval = Aval - Val
	Bval = Bval + Val
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(usernameA, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(usernameB, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(YGTokenChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
