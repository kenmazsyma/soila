/*
Package project provdes chaincode for managing PROJECT data.
*/
package project

import (
	"encoding/json"
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
	PeerKey string // Hash of PEER [key]
	Id      string // ID of project [key]
	Status  STATUS // status of project
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
	return stub.CreateCompositeKey(KEY_TYPE, args)
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
func Register(stub shim.ChaincodeStubInterface, args []string) (key, res string, err error) {
	// check parameter
	if err = cmn.CheckParam(args, 1); err != nil {
		return
	}
	// get PEER key of sender
	peerkey, err := peer.GetKey(stub)
	if err != nil {
		return
	}
	if key, err = cmn.VerifyForRegistration(stub, generateKey, []string{peerkey, args[0]}); err != nil {
		return
	}
	log.Debug(key)
	// put data into ledger
	data := Project{
		PeerKey: peerkey,
		Id:      args[0],
		Status:  STATUS_ACTIVE,
	}
	err = cmn.Put(stub, key, data)
	return
}

// Get is a function for getting PROJECT information from ledger
//   parameters :
//     stub - object of chaincode information
//     args - [projectkey]
//  return :
//    - response data
//    - either error object or nil
func Get(stub shim.ChaincodeStubInterface, args []string) (key, res string, err error) {
	return cmn.Get(stub, args)
}

// UpdateStatus is a function for updating PROJECT staus
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - [peerkey, id, status]
//   return :
//     - response data
//     - error object if error occured
func UpdateStatus(stub shim.ChaincodeStubInterface, args []string) (key, res string, err error) {
	data, key, err := get_and_check(stub, args, 3)
	if err != nil {
		return
	}
	val, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		log.Info("status parameter is not correct.")
		return
	}
	if int64(data.Status) == val {
		log.Info("status parameter is not different from ledger.")
		return
	}
	log.Debug(key)
	data.Status = STATUS(val)
	err = cmn.Put(stub, key, data)
	return
}
