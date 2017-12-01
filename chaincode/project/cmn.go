/*
Package project provdes chaincode for managing PROJECT data.
*/
package project

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/kenmazsyma/soila/chaincode/cmn"
	"github.com/kenmazsyma/soila/chaincode/log"
	"github.com/kenmazsyma/soila/chaincode/peer"
)

// GetKeyInPeer is a function for checking if project exists in sender's peer
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     id - projectid
//   return :
//     - either key or empty string
func GetKeyInPeer(stub shim.ChaincodeStubInterface, id string) string {
	// get PEER key of sender
	peerkey, err := peer.GetKey(stub)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	js, key, err := cmn.VerifyForUpdate(stub, generateKey, []string{peerkey, id}, 2)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	if js != nil {
		return key
	}
	return ""
}

// IsOwn is a function for checking if project is owned by sender
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     key - projectkey
//   returns :
//     - true if sender own
func IsOwn(stub shim.ChaincodeStubInterface, key string) bool {
	_, val, err := stub.SplitCompositeKey(key)
	if err != nil {
		log.Error(err.Error())
		return false
	}
	if len(val) != 3 {
		log.Error("key is not correct:" + key)
		return false
	}
	valid, err := peer.CompareHash(stub, []byte(val[1]))
	if err != nil {
		log.Error(err.Error())
		return false
	}
	return valid
}
