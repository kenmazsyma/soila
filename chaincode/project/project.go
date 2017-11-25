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

type Project struct {
	Peer []byte
	Id   string
}

const KEY_TYPE = "PROJECT"

// genearteKey is a function for generating key from id of PROJECT
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - arguments which contains keys
//   return :
//     - key
//     - whether error object or nil
func generateKey(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	peerid, err := peer.GetId(stub)
	if err != nil {
		return "", err
	}
	return stub.CreateCompositeKey(KEY_TYPE, []string{string(peerid), args[0]})
}

// Register is a function for registering PROJECT to ledger
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - [id]
//   return :
//     - response data
//     - error object if error occured
func Register(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	log.Debug("start:")
	key, err := cmn.VerifyForRegistration(stub, generateKey, args, 1)
	if err != nil {
		return "", err
	}
	log.Debug(key)
	peerid, err := peer.GetId(stub)
	if err != nil {
		return "", err
	}
	log.Debug(string(peerid))
	// put data into ledger
	data := Project{
		Peer: peerid,
		Id:   args[0],
	}
	err = cmn.Put(stub, key, data)
	return "", err
}
