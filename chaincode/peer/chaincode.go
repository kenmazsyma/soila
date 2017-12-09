/*
Package peer provides chaincode for managing PEER data.
  TODO: nessesary to implement logic for verification whether peer can be trusted
*/
package peer

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/kenmazsyma/soila/chaincode/cmn"
	. "github.com/kenmazsyma/soila/chaincode/log"
)

type Peer struct {
	Hash    string `json:"hash"`    // hash of peer's signature[key]
	Address string `json:"address"` // url of wenapp / webapi
}

const KEY_TYPE = "PEER"

// genearteKey is a function for generating key from id of PEER
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - not use
//   return :
//     - key
//     - whether error object or nil
func generateKey(stub shim.ChaincodeStubInterface, args interface{}) (ret string, err error) {
	sig, err := stub.GetCreator()
	if err != nil {
		return
	}
	hash := cmn.Sha512B(sig)
	return stub.CreateCompositeKey(KEY_TYPE, []string{hash})
}

// Register is a function for registering PEER informartion
//   parameters :
//     stub - object of chaincode information
//     args - [address]
//   return :
//     ret - return value
//     err - either error object or nil
func Register(stub shim.ChaincodeStubInterface, args []string) (ret []interface{}, err error) {
	D("get peer's signature")
	info := Peer{}
	D("check parameter")
	if err = cmn.CheckParam(args, 1); err != nil {
		return
	}
	D("verify if peer is already registered")
	key, err := cmn.VerifyForRegistration(stub, generateKey, nil)
	if err != nil {
		return
	}
	sig, err := stub.GetCreator()
	if err != nil {
		return
	}
	info.Hash = cmn.Sha512B(sig)
	Info("register peer: %s", key)

	info.Address = args[0]
	err = cmn.Put(stub, key, info)
	ret = []interface{}{[]byte(key)}
	return
}

// Get is a function for getting PEER information from ledger
//   parameters :
//     stub - object of chaincode information
//     args - [key]
//   return :
//     - return value
//     - either error object or nil
func Get(stub shim.ChaincodeStubInterface, args []string) (data []interface{}, err error) {
	if err = cmn.CheckParam(args, 1); err != nil {
		return
	}
	return cmn.Get(stub, args[0])
}

// Update is a function for updating PEER information
//   parameters :
//     stub - object of chaincode information
//     args - [address]
//   return :
//     ret - return value
//     err - either error object or nil
func Update(stub shim.ChaincodeStubInterface, args []string) (ret []interface{}, err error) {
	D("check parameter")
	if err = cmn.CheckParam(args, 1); err != nil {
		return
	}
	D("generate key from peer's signature")
	key, err := generateKey(stub, nil)
	if err != nil {
		return
	}
	D("get data from ledger")
	js, err := cmn.VerifyForUpdate(stub, []string{key}, 1)
	if err != nil {
		return
	}
	D("get:%s", string(js))
	data := Peer{}
	err = json.Unmarshal(js, &data)
	if err != nil {
		return
	}
	D("udpate data")
	data.Address = args[0]
	err = cmn.Put(stub, key, data)
	ret = []interface{}{[]byte(key)}
	return
}

// Deregister is a function for removing PEER information
// TODO:consider the condition for allowing peer to deregister
//   parameters :
//     stub - object of chaincode information
//     args - []
//   return :
//     ret - return value
//     err - either error object or nil
func Deregister(stub shim.ChaincodeStubInterface, args []string) (ret []interface{}, err error) {
	D("check parameter")
	if err = cmn.CheckParam(args, 0); err != nil {
		return
	}
	D("generate key from peer's signature")
	key, err := generateKey(stub, nil)
	if err != nil {
		return
	}
	D("check if data is exists:%s", key)
	_, err = cmn.Get(stub, key)
	if err != nil {
		return
	}
	D("delete data")
	err = cmn.Delete(stub, key)
	return
}
