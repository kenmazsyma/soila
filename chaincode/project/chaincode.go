/*
Package project provdes chaincode for managing PROJECT data.
*/

package project

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/kenmazsyma/soila/chaincode/cmn"
	"github.com/kenmazsyma/soila/chaincode/log"
	"github.com/kenmazsyma/soila/chaincode/peer"
	"strconv"
)

type STATUS int

const (
	STATUS_ACTIVE STATUS = iota
	STATUS_SLEEP
)

type Project struct {
	Peer   []byte // Hash of peer id [key]
	Id     string // ID of prokect [key]
	Status STATUS // status of project
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

// get_and_check is a function for getting data of PERSON
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - parameters received from client
//     nofElm - valid length of args
//   return :
//     - PERSON object
//     - key
//     - whether error object or nil
func get_and_check(stub shim.ChaincodeStubInterface, args []string, nofElm int) (rec *Project, key string, err error) {
	rec = nil
	js, key, err := cmn.VerifyForUpdate(stub, generateKey, args, nofElm)
	if err != nil {
		return
	}
	*rec = Project{}
	err = json.Unmarshal(js, rec)
	return
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
		Peer:   peerid,
		Id:     args[0],
		Status: STATUS_ACTIVE,
	}
	err = cmn.Put(stub, key, data)
	return "", err
}

// Get is a function for getting PROJECT information from ledger
//   parameters :
//     stub - object of chaincode information
//     args - [id]
//  return :
//    - response data
//    - either error object or nil
func Get(stub shim.ChaincodeStubInterface, args []string) (res string, err error) {
	return cmn.Get(stub, generateKey, args, 1)
}

// UpdateStatus is a function for updating PROJECT staus
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - [peerid, id, status]
//   return :
//     - response data
//     - error object if error occured
func UpdateStatus(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	data, key, err := get_and_check(stub, args, 3)
	if err != nil {
		return "", err
	}
	val, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		log.Info("status parameter is not correct.")
		return "", err
	}
	if int64(data.Status) == val {
		log.Info("status parameter is not different from ledger.")
		return "", nil
	}
	log.Debug(key)
	valid, err := peer.CompareId(stub, data.Peer)
	if err != nil {
		return "", err
	}
	// peer id is different from sender id
	if !valid {
		return "", errors.New("Project is not found in ledger.")
	}
	return "", err
}
