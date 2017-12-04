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
//     ret - return value
//     err - either error object or nil
func Register(stub shim.ChaincodeStubInterface, args []string) (ret []interface{}, err error) {
	D("get peer's signature")
	info := Peer{}
	sig, err := stub.GetCreator()
	if err != nil {
		return
	}
	D("check parameter")
	if err = cmn.CheckParam(args, 1); err != nil {
		return
	}
	D("verify if peer is already registered")
	info.Hash = cmn.Sha512B(sig)
	D("hash:%s", info.Hash)
	key, err := cmn.VerifyForRegistration(stub, generateKey, []string{string(info.Hash)})
	if err != nil {
		return
	}
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
func Get(stub shim.ChaincodeStubInterface, args []string) ([]interface{}, error) {
	return cmn.Get(stub, args)
}

// Update is a function for updating PEER information
//   parameters :
//     stub - object of chaincode information
//     args - [key, address]
//   return :
//     ret - return value
//     err - either error object or nil
func Update(stub shim.ChaincodeStubInterface, args []string) (ret []interface{}, err error) {
	//	D("check parameter")
	//	if err = cmn.CheckParam(args, 2); err != nil {
	//		return
	//	}
	//	D("check if data is exists")
	//	val, err := stub.GetState(args[0])
	//	if err != nil {
	//		return
	//	}
	//	if len(val) == 0 {
	//		err = errors.New("data not found.")
	//		return
	//	}
	js, err := cmn.VerifyForUpdate(stub, args, 2)
	if err != nil {
		return
	}
	D("check if data is owned by sender")
	data := Peer{}
	err = json.Unmarshal(js, &data)
	if err != nil {
		return
	}
	D("hash:%s", data.Hash)
	valid, err := CompareHash(stub, data.Hash)
	if err != nil {
		return
	}
	if !valid {
		err = errors.New("peer is not owned by sender.")
		return
	}
	if data.Address == args[1] {
		return
	}
	D("udpate data")
	data.Address = args[1]
	err = cmn.Put(stub, args[0], data)
	ret = []interface{}{[]byte(args[0])}
	return
}

// Deregister is a function for removing PEER information
// TODO:consider the condition for allowing peer to deregister
//   parameters :
//     stub - object of chaincode information
//     args - [key]
//   return :
//     ret - return value
//     err - either error object or nil
func Deregister(stub shim.ChaincodeStubInterface, args []string) (ret []interface{}, err error) {
	D("check parameter")
	if err = cmn.CheckParam(args, 1); err != nil {
		return
	}
	D("check if data is exists:%s", args[0])
	val, err := stub.GetState(args[0])
	if err != nil {
		return
	}
	D("val:%s", val)
	data := Peer{}
	err = json.Unmarshal(val, &data)
	if err != nil {
		return
	}
	D("hash:%s", data.Hash)
	D("verify if data is owned by sender")
	valid, err := CompareHash(stub, data.Hash)
	if err != nil {
		return
	}
	if !valid {
		err = errors.New("Peer is not owned by sender")
		return
	}
	D("delete data")
	err = cmn.Delete(stub, args[0])
	return
}
