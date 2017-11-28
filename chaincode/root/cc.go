/*
Package root provides root routine for chaincode.
*/
package root

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/kenmazsyma/soila/chaincode/log"
)

type InvokeRoutineType func(shim.ChaincodeStubInterface, []string) (string, string, error)
type CC struct {
	sub map[string]InvokeRoutineType
}

func (t *CC) SetInvokeMap(sub map[string]InvokeRoutineType) {
	t.sub = sub
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

// Invoke is a function for executing chaincode for soila_chain
//   parameters :
//     stub - chaincode interface
//   return :
//     response object
func (t *CC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	log.Debug("CC.Invoke")
	funcname, args := stub.GetFunctionAndParameters()
	log.Debug(funcname)
	m := t.sub[funcname]
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
