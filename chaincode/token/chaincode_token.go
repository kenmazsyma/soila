/*
Package token provides chaincode for managing TOKEN data.
*/

package token

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/kenmazsyma/soila/chaincode/cmn"
	"github.com/kenmazsyma/soila/chaincode/log"
	"github.com/kenmazsyma/soila/chaincode/peer"
)

type Token struct {
	CreatorPeerID    []byte // peer id of creator project[key]
	CreatorProjectID string // id of creator project[key]

	Name     string // name of token
	DescHash string // hash of description of token
	// rule // TODO:
}

const KEY_TOKEN = "TOKEN"

// generateKey is a function for generating key from id of PROJECT
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - arguments which contains key
//   return :
//     - key
//     - whether error object or nil
func generateKey(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	return stub.CreateCompositeKey(KEY_TOKEN, args[0:2])
}

// get_and_check is a function for getting data of PERSON
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - parameters received from client
//     nofElm - valid length of args
//   return :
//     - PERSON object
//     - key
//     - whether error object or nil
func get_and_check(stub shim.ChaincodeStubInterface, args []string, nofElm int) (rec *Token, key string, err error) {
	rec = nil
	js, key, err := cmn.VerifyForUpdate(stub, generateKey, args, nofElm)
	if err != nil {
		return
	}
	*rec = Token{}
	err = json.Unmarshal(js, rec)
	return
}

// Register is a function for registering PROJECT to ledger
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - [projectid, name, desc]
//   return :
//     - response data
//     - error object if error occured
func Register(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	log.Debug("start:")
	peerid, err := peer.GetId(stub)
	if err != nil {
		return "", err
	}
	key, err := cmn.VerifyForRegistration(stub, generateKey, []string{string(peerid), args[0]}, 3)
	if err != nil {
		return "", err
	}
	log.Debug(key)
	// hash of description
	deschash := cmn.Sha1(args[2])
	// put data into ledger
	data := Token{
		CreatorPeerID:    peerid,
		CreatorProjectID: args[0],
		Name:             args[1],
		DescHash:         deschash,
	}
	err = cmn.Put(stub, key, data)
	return "", err
}

// Get is a function for getting TOKEN information from ledger
//   parameters :
//     stub - object of chaincode information
//     args - [peerid, projectid]
//  return :
//    - response data
//    - either error object or nil
func Get(stub shim.ChaincodeStubInterface, args []string) (res string, err error) {
	return cmn.Get(stub, generateKey, args, 2)
}

// Update is a function for updating TOKEN information
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - [projectid, name, description]
//   return :
//     - response data
//     - error object if error occured
func Update(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 3 {
		return "", errors.New("invalid param")
	}
	peerid, err := peer.GetId(stub)
	if err != nil {
		return "", err
	}
	data, key, err := get_and_check(stub, []string{string(peerid), args[0]}, 2)
	if err != nil {
		return "", err
	}
	log.Debug(key)
	// peer id is different from sender id
	if data == nil {
		return "", errors.New("Data is not found in ledger.")
	}
	//TODO:Token information can be updated only in the case that token is not issued yet
	data.Name = args[1]
	data.DescHash = cmn.Sha1(args[2])
	err = cmn.Put(stub, key, data)
	return "", err
}

// Remove is a function for updating PEER information
//   parameters :
//     stub - object of chaincode information
//     args - [projectid]
//  return :
//    - response data
//    - either error object or nil
func Remove(stub shim.ChaincodeStubInterface, args []string) (res string, err error) {
	//TODO:Token information can be updated only in the case that token is not issued yet
	res = ""
	if len(args) != 1 {
		err = errors.New("Invalid Arguments")
		return
	}
	peerid, err := peer.GetId(stub)
	if err != nil {
		return "", err
	}
	val, key, err := get_and_check(stub, []string{string(peerid), args[0]}, 2)
	if err != nil {
		return "", err
	}
	log.Debug(key)
	// peer id is different from sender id
	if val == nil {
		return "", errors.New("Data is not found in ledger.")
	}
	err = cmn.Delete(stub, key)
	return "", err
}
