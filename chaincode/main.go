/*
Package main provides chaincode for soila_chain.
*/
package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/kenmazsyma/soila/chaincode/log"
	"github.com/kenmazsyma/soila/chaincode/peer"
	"github.com/kenmazsyma/soila/chaincode/person"
	"github.com/kenmazsyma/soila/chaincode/project"
	"github.com/kenmazsyma/soila/chaincode/token"
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

type invokeRoutineType func(shim.ChaincodeStubInterface, []string) (string, string, error)

var invoke_list = map[string]invokeRoutineType{
	"person.register":          person.Register,
	"person.update":            person.Update,
	"person.get":               person.Get,
	"person.add_activity":      person.AddActivity,
	"person.add_reputation":    person.AddReputation,
	"person.remove_reputation": person.RemoveReputation,
	"peer.register":            peer.Register,
	"peer.update":              peer.Update,
	"peer.get":                 peer.Get,
	"peer.deregister":          peer.Deregister,
	"project.register":         project.Register,
	"project.get":              project.Get,
	"project.updatestatus":     project.UpdateStatus,
	"token.register":           token.Register,
	"token.update":             token.Update,
	"token.remove":             token.Remove,
}

// Invoke is a function for executing chaincode for soila_chain
//   parameters :
//     stub - chaincode interface
//   return :
//     response object
func (t *CC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	log.Debug("CC.Invoke")
	funcname, args := stub.GetFunctionAndParameters()
	log.Debug(funcname)
	m := invoke_list[funcname]
	if m == nil {
		return shim.Error("Invalid function name.")
	}
	ret1, ret2, err := m(stub, args)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("['" + ret1 + "'," + ret2 + "]"))
}

// ================================================
//  Query
// ================================================

// Query is a deprecated in Fabric v1.0.
func (t *CC) Query(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Error("Query interface was already removed.")
}
