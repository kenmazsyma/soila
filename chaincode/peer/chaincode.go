/*
Package peer provides chaincode for managing PEER data.
TODO: nessesary to implement logic for verification whether peer can be trusted
*/

package peer

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/kenmazsyma/soila/chaincode/cmn"
	"github.com/kenmazsyma/soila/chaincode/log"
)

type Peer struct {
	Hash    []byte // hash of peer's signature[key]
	Address string // url of wenapp / webapi
}

const KEY_TYPE = "PEER"

// genearteKey is a function for generating key from id of PEER
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - arguments which contains keys
//   return :
//     - key
//     - whether error object or nil
func generateKey(stub shim.ChaincodeStubInterface, args []string) (ret string, err error) {
	return stub.CreateCompositeKey(KEY_TYPE, []string{args[0]})
}

// Register is a function for registering PEER informartion
//   parameters :
//     stub - object of chaincode information
//     args - [address]
//   return :
//     key - key value
//     res - response data
//     err - either error object or nil
func Register(stub shim.ChaincodeStubInterface, args []string) (key, res string, err error) {
	// get peer's signature
	log.Info("start:")
	info := Peer{}
	sig, err := stub.GetCreator()
	if err != nil {
		return
	}
	// check parameter
	if err = cmn.CheckParam(args, 1); err != nil {
		return
	}
	// verify if peer is already registered
	info.Hash = cmn.Sha1B(sig)
	log.DebugB(info.Hash)
	key, err = cmn.VerifyForRegistration(stub, generateKey, []string{string(info.Hash)})
	if err != nil {
		return
	}
	// register peer
	log.Info("Register:" + key)
	info.Address = args[0]
	err = cmn.Put(stub, key, info)
	return
}

// Get is a function for getting PEER information from ledger
//   parameters :
//     stub - object of chaincode information
//     args - [key]
//   return :
//     key - key value
//     res - response data
//     err - either error object or nil
func Get(stub shim.ChaincodeStubInterface, args []string) (key, res string, err error) {
	return cmn.Get(stub, generateKey, args, 1)
}

// Update is a function for updating PEER information
//   parameters :
//     stub - object of chaincode information
//     args - [key, address]
//   return :
//     key - key value
//     res - response data
//     err - either error object or nil
func Update(stub shim.ChaincodeStubInterface, args []string) (key, res string, err error) {
	// check parameter
	if err = cmn.CheckParam(args, 1); err != nil {
		return
	}
	// check if data is exist
	log.Info("start:")
	key = args[0]
	val, err := stub.GetState(key)
	if err != nil {
		return
	}
	if val == nil {
		return
	}
	// check if data is owned by sender
	log.DebugB(val)
	data := Peer{}
	err = json.Unmarshal(val, &data)
	if err != nil {
		return
	}
	log.DebugB(data.Hash)
	valid, err := CompareHash(stub, data.Hash)
	if err != nil {
		return
	}
	if !valid {
		err = errors.New("Peer is not owned by sender")
		return
	}
	if data.Address == args[1] {
		return
	}
	// update data
	data.Address = args[1]
	err = cmn.Put(stub, key, data)
	return
}

// Deregister is a function for removing PEER information
// TODO:consider the condition for allowing peer to deregister
//   parameters :
//     stub - object of chaincode information
//     args - [key]
//   return :
//     key - key value
//     res - response data
//     err - either error object or nil
func Deregister(stub shim.ChaincodeStubInterface, args []string) (key, res string, err error) {
	// check parameter
	if err = cmn.CheckParam(args, 1); err != nil {
		return
	}
	// check if data is exist
	key = args[0]
	val, err := stub.GetState(key)
	if err != nil {
		return
	}
	log.DebugB(val)
	data := Peer{}
	err = json.Unmarshal(val, &data)
	if err != nil {
		return
	}
	log.DebugB(data.Hash)
	// verify if data is owned by sender
	valid, err := CompareHash(stub, data.Hash)
	if err != nil {
		return
	}
	if !valid {
		err = errors.New("Peer is not owned by sender")
		return
	}
	// delete data
	err = cmn.Delete(stub, key)
	return
}
