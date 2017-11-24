/*
Package person provides chaincode for managing person data.
*/

package person

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/kenmazsyma/soila/chaincode/cmn"
	"github.com/kenmazsyma/soila/chaincode/log"
	"github.com/kenmazsyma/soila/chaincode/peer"
)

type PersonReputation struct {
	Setter  string
	Content string
	Type    string
}

type Person struct {
	Id         string
	Peer       []byte
	Ver        []string
	Activity   []string
	Reputation []PersonReputation
}

const KEY_TYPE = "PERSON"

// genearteKey is a function for generating key from id of PERSON
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     id - id of PERSON
//   return :
//     - key
//     - whether error object or nil
func generateKey(stub shim.ChaincodeStubInterface, id string) (string, error) {
	return stub.CreateCompositeKey(KEY_TYPE, []string{id})
}

// get_and_check is a function for getting data of PERSON
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - parameters received from client
//     validlen - valid length of args
//   return :
//     - PERSON object
//     - key
//     - whether error object or nil
func get_and_check(stub shim.ChaincodeStubInterface, args []string, validlen int) (rec *Person, key string, err error) {
	if len(args) != validlen {
		err = errors.New("Invalid Arguments")
		return
	}
	// get ID from blockchain
	key, err = generateKey(stub, args[0])
	if err != nil {
		return
	}
	// check if data is already exists.
	val, err := stub.GetState(key)
	if err != nil {
		return
	}
	if val == nil {
		err = errors.New("data is not exists.")
		return
	}
	data := Person{}
	err = json.Unmarshal(val, &data)
	rec = &data
	return
}

// Put is a function for registering PERSON to ledger
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - [personid, data]
//   return :
//     - response data
//     - error object if error occured
func Put(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("Invalid Arguments")
	}
	log.Info("start:")
	key, err := generateKey(stub, args[0])
	if err != nil {
		return "", err
	}
	log.Debug("KEY:" + key)
	// check if data is already exists.
	val, err := stub.GetState(key)
	if err != nil {
		return "", err
	}
	if val != nil {
		return "", errors.New("data is already exists.")
	}
	log.Debug(string(val))
	peerid, err := peer.GetId(stub)
	if err != nil {
		return "", err
	}
	log.Debug(string(peerid))
	// put data into ledger
	data := Person{
		Id:         args[0],
		Peer:       peerid,
		Ver:        []string{cmn.Sha1(args[1])},
		Activity:   []string{},
		Reputation: []PersonReputation{},
	}
	err = cmn.Put(stub, key, data)
	return "", err
}

// Update is a function for updating PERSON object
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - [personid, data]
//   return :
//     - response data
//     - error object if error occured
func Update(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	data, key, err := get_and_check(stub, args, 2)
	if err != nil {
		return "", err
	}
	log.Debug(key)
	valid, err := peer.CompareId(stub, data.Peer)
	if err != nil {
		return "", err
	}
	// peer id is different from sender id
	if !valid {
		return "", errors.New("Person is owned by another peer.")
	}
	// put data into ledger
	(*data).Ver = append((*data).Ver, cmn.Sha1(args[1]))
	err = cmn.Put(stub, key, (*data))
	return "", err
}

// Get is a function for getting PERSON object
//   parameters  :
//     stub - object for accessing ledgers from chaincode
//     args - personid
//   return :
//     - json data of PERSON data
//     - error string if error occured
func Get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("Invalid parameter")
	}
	log.Info("start")
	key, err := generateKey(stub, args[0])
	if err != nil {
		return "", err
	}
	log.Debug(key)
	val, err := stub.GetState(key)
	log.Debug(string(val))
	return string(val), err
}

// AddActivity is a function for append hash of activity information for PERSON
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - [personid, content hash]
//   returns :
//     - response data
//     - whether error object or nil
func AddActivity(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	data, key, err := get_and_check(stub, args, 2)
	if err != nil {
		return "", err
	}
	log.Debug(key)
	// check whether data is owned by sender peer
	valid, err := peer.CompareId(stub, data.Peer)
	if err != nil {
		return "", err
	}
	if !valid {
		return "", errors.New("Person is owned by another peer.")
	}
	// put data into ledger
	(*data).Activity = append((*data).Activity, args[1])
	err = cmn.Put(stub, key, (*data))
	return "", err
}

// AddReputation is a function for append hash of reputation information for PERSON
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - [personid, content hash]
//   returns :
//     - response data
//     - whether error object or nil
func AddReputation(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	data, key, err := get_and_check(stub, args, 4)
	if err != nil {
		return "", err
	}
	log.Debug(key)
	// check whether data is owned by sender peer
	valid, err := peer.CompareId(stub, data.Peer)
	if err != nil {
		return "", err
	}
	if !valid {
		return "", errors.New("Person is owned by another peer.")
	}
	// put data into ledger
	rep := PersonReputation{
		Setter:  args[1],
		Content: args[2],
		Type:    args[3],
	}
	(*data).Reputation = append((*data).Reputation, rep)
	err = cmn.Put(stub, key, (*data))
	return "", err
}

/***************************************************
[RemoveReputation]
description : add activity hash of Person to Person object
parameters  :
   stub - chaincode interface
   args - [personid, setter, content hash]
return:
   1:response data
   2:error object if error occured
***************************************************/
// RemoveReputation is a function for remove hash of reputation information for PERSON
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - [personid, content hash]
//   returns :
//     - response data
//     - whether error object or nil
func RemoveReputation(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	data, key, err := get_and_check(stub, args, 3)
	if err != nil {
		return "", err
	}
	log.Debug(key)
	// check whether data is owned by sender peer
	valid, err := peer.CompareId(stub, data.Peer)
	if err != nil {
		return "", err
	}
	if !valid {
		return "", errors.New("Person is owned by another peer.")
	}
	// put object removed target reputation data into ledger
	for i, v := range (*data).Reputation {
		log.Debug(v.Setter + "," + args[1] + "," + v.Content + "," + args[2])
		if v.Setter == args[1] && v.Content == args[2] {
			(*data).Reputation = append((*data).Reputation[0:i], (*data).Reputation[i+1:]...)
			err = cmn.Put(stub, key, (*data))
			return "", err
		}
	}
	return "", errors.New("target data is not found")
}
