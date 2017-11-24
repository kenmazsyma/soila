/*
Package peer provides chaincode for managing peer data.
TODO: nessesary to implement logic for verification whether peer can be trusted
*/

package peer

import (
	"bytes"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/kenmazsyma/soila/chaincode/cmn"
	"github.com/kenmazsyma/soila/chaincode/log"
)

// GetId is a function for getting ID of creator peer
//   parameters :
//     stub - object for accessing ledgers from chaincode
//   returns :
//     - ID fo creator peer
//     - whether error object or nil
func GetId(stub shim.ChaincodeStubInterface) ([]byte, error) {
	creator, err := stub.GetCreator()
	if err != nil {
		return nil, err
	}
	return cmn.Sha1Byte(creator), nil
}

// CompareId is a function for comparing with ID of sender peer
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     target - data for target
//   returns :
//     - true if matched
//     - whether error object or nil
func CompareId(stub shim.ChaincodeStubInterface, target []byte) (bool, error) {
	mypeer, err := GetId(stub)
	if err != nil {
		return false, err
	}
	log.Debug(string(mypeer))
	return (bytes.Compare(mypeer, target) == 0), nil
}
