/*
Package peer provides chaincode for managing peer data.
TODO: nessesary to implement logic for verification whether peer can be trusted
*/
package peer

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/kenmazsyma/soila/chaincode/cmn"
	. "github.com/kenmazsyma/soila/chaincode/log"
)

// GetHash is a function for getting hash of sender's signature
//   parameters :
//     stub - object for accessing ledgers from chaincode
//   returns :
//     - hash of sender's signature
//     - whether error object or nil
func GetHash(stub shim.ChaincodeStubInterface) (string, error) {
	creator, err := stub.GetCreator()
	if err != nil {
		return "", err
	}
	return cmn.Sha1B(creator), nil
}

// GetKey is a function for getting key of sender's PEER data
//   parameters :
//     stub - object for accessing ledgers from chaincode
//   returns :
//     - hash of sender's signature
//     - whether error object or nil
func GetKey(stub shim.ChaincodeStubInterface) (ret string, err error) {
	hash, err := GetHash(stub)
	if err != nil {
		return
	}
	ret, err = generateKey(stub, []string{string(hash)})
	return
}

// CompareHash is a function for comparing with hash of sender peer
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     target - target hash
//   returns :
//     - true if matched
//     - whether error object or nil
func CompareHash(stub shim.ChaincodeStubInterface, target string) (bool, error) {
	mypeer, err := GetHash(stub)
	if err != nil {
		return false, err
	}
	D("compare hash:sender(%s), target(%s)", mypeer, target)
	return (mypeer == target), nil
}
