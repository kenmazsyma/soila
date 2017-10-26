package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/kenmazsyma/soila/chaincode/person"
)

// ================================================
// main
// ================================================
func main() {
	err := shim.Start(new(CC))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

type CC struct{}

/***************************************************
[Init]
description : initialize
parameters  :
   stub - chaincode interface
return: response object
***************************************************/
func (t *CC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("CC.Init")
	return shim.Success(nil)
}

// ================================================
//  Invoke
// ================================================

type invokeRoutineType func(shim.ChaincodeStubInterface, []string) (string, error)

var invoke_list = map[string]invokeRoutineType{
	"person.put":               person.Put,
	"person.update":            person.Update,
	"person.get":               person.Get,
	"person.add_activity":      person.AddActivity,
	"person.add_reputation":    person.AddReputation,
	"person.remove_reputation": person.RemoveReputation,
}

/***************************************************
[Invoke]
description : invoke
parameters  :
   stub - chaincode interface
return: response object
***************************************************/
func (t *CC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("CC.Invoke")
	funcname, args := stub.GetFunctionAndParameters()
	fmt.Printf("NAME:%s\n", funcname)
	m := invoke_list[funcname]
	if m == nil {
		return shim.Error("Invalid function name.")
	}
	ret, err := m(stub, args)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(ret))
}

// ================================================
//  Query
// ================================================

/***************************************************
[Query]
description : query
parameters  :
   stub - chaincode interface
return: response object
***************************************************/
func (t *CC) Query(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Error("Query interface was already removed.")
}
