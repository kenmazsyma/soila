/*
Package main provides chaincode for soila_chain.
*/

package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/kenmazsyma/soila/chaincode/log"
	"github.com/kenmazsyma/soila/chaincode/person"
)

// ================================================
// main
// ================================================

type CC struct{}

// main is a function for executing chaincode for soila_chain
func main() {
	log.Init("soila_chain", log.LEVEL_DEBUG)
	err := shim.Start(new(CC))
	if err != nil {
		log.Error("Error starting Simple chaincode: " + err.Error())
	}
}

// Init is a function for initializing chaincode for soila_chain
//   parameters :
//     stub - chaincode interface
//   return :
//     response object
func (t *CC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	log.Debug("CC.Init")
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

// Invoke is a function for executing chaincode for soila_chain
//   parameters :
//     stub - chaincode interface
//   return :
//     response object
func (t *CC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	log.Debug("CC.Invoke")
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

// Query is a deprecated in Fabric v1.0.
func (t *CC) Query(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Error("Query interface was already removed.")
}
