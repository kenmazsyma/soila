/*
Package root provides root routine for chaincode.
*/
package root

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/kenmazsyma/soila/chaincode/cmn"
	. "github.com/kenmazsyma/soila/chaincode/log"
)

type InvokeRoutineType func(shim.ChaincodeStubInterface, []string) ([]interface{}, error)
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
	D("CC.Init")
	return shim.Success(nil)
}

// Invoke is a function for executing chaincode for soila_chain
//   parameters :
//     stub - chaincode interface
//   return :
//     response object
func (t *CC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	D("CC.Invoke")
	funcname, args := stub.GetFunctionAndParameters()
	D(funcname)
	m := t.sub[funcname]
	if m == nil {
		return shim.Error("Invalid function name.")
	}
	decoded := []string{}
	for i := 0; i < len(args); i++ {
		if len(args[i]) > 0 {
			key, err := cmn.DecodeBase64(args[i])
			if err != nil {
				return shim.Error(err.Error())
			}
			decoded = append(decoded, string(key))
		} else {
			decoded = append(decoded, args[i])
		}
	}
	ret, err := m(stub, decoded)
	if err != nil {
		return shim.Error(err.Error())
	}
	//	rest := []string{}
	//	for _, v := range ret {
	//		rest = append(rest, string(v.([]byte)))
	//		D("%s => %s", funcname, v.([]byte))
	//	}
	//	js, err := json.Marshal(rest)
	//	D("%s => %s", funcname, js)
	js, err := json.Marshal(ret)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(js)
}

// ================================================
//  Query
// ================================================

// Query is a deprecated in Fabric v1.0.
func (t *CC) Query(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Error("Query interface was already removed.")
}
