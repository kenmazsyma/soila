/*
Package person provides chaincode for managing PERSON data.
TODO:nessesary to implement authentication logic
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
	Setter  string // key of PERSON who generate this data
	Content string // key of CONTENT
	Type    string // type of reputation(TODO:under consideration)
}

type Person struct {
	PeerKey    string             // key of peer to which person belong
	Id         string             // id of person
	Ver        []string           // list of hash for information
	Activity   []string           // list of key of CONTENT which PERSON acts
	Reputation []PersonReputation // list of reputation for PERSON
}

const KEY_TYPE = "PERSON"

// generateKey is a function for generating key from id of PERSON
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
func get_and_check(stub shim.ChaincodeStubInterface, args []string, nofElm int) (rec *Person, key string, err error) {
	rec = nil
	js, key, err := cmn.VerifyForUpdate(stub, generateKey, args, nofElm)
	if err != nil {
		return
	}
	*rec = Person{}
	err = json.Unmarshal(js, rec)
	return
}

// Register is a function for registering PERSON to ledger
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - [personid, data]
//   return :
//     - response data
//     - error object if error occured
func Register(stub shim.ChaincodeStubInterface, args []string) (key, res string, err error) {
	// check parameter
	if err = cmn.CheckParam(args, 2); err != nil {
		return
	}
	// get PEER key for sender
	peerkey, err := peer.GetKey(stub)
	if err != nil {
		return
	}
	log.Debug(string(peerkey))
	// check if data is already exist
	key, err = cmn.VerifyForRegistration(stub, generateKey, []string{peerkey, args[0]})
	if err != nil {
		return
	}
	// put data into ledger
	data := Person{
		PeerKey:    peerkey,
		Id:         args[0],
		Ver:        []string{cmn.Sha1(args[1])},
		Activity:   []string{},
		Reputation: []PersonReputation{},
	}
	err = cmn.Put(stub, key, data)
	return
}

// Update is a function for updating PERSON object
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - [personid, data]
//   return :
//     - response data
//     - error object if error occured
func Update(stub shim.ChaincodeStubInterface, args []string) (key, res string, err error) {
	// check parameter
	if len(args) != 2 {
		err = errors.New("Invalid Parameter")
		return
	}
	// get PEER key for sender
	peerkey, err := peer.GetKey(stub)
	if err != nil {
		return
	}
	log.Debug(string(peerkey))
	// check if data is exists
	data, key, err := get_and_check(stub, []string{peerkey, args[0]}, 2)
	if err != nil {
		return
	}
	if data == nil {
		err = errors.New("data is not exists")
		return
	}
	log.Debug(key)
	// put data into ledger
	(*data).Ver = append((*data).Ver, cmn.Sha1(args[1]))
	err = cmn.Put(stub, key, (*data))
	return
}

// Get is a function for getting PERSON object
//   parameters  :
//     stub - object for accessing ledgers from chaincode
//     args - [personkey]
//   return :
//     - json data of PERSON data
//     - error string if error occured
func Get(stub shim.ChaincodeStubInterface, args []string) (string, string, error) {
	return cmn.Get(stub, args)
}

// AddActivity is a function for append hash of activity information for PERSON
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - [peerkey, personid, contentkey]
//   returns :
//     - response data
//     - whether error object or nil
func AddActivity(stub shim.ChaincodeStubInterface, args []string) (key, res string, err error) {
	// check parameter
	if err = cmn.CheckParam(args, 3); err != nil {
		return
	}
	// check if data is exists
	data, key, err := get_and_check(stub, args, 2)
	if err != nil {
		return
	}
	if data == nil {
		err = errors.New("data is not exists")
		return
	}
	log.Debug(key)
	// put data into ledger
	(*data).Activity = append((*data).Activity, args[2])
	err = cmn.Put(stub, key, (*data))
	return
}

// AddReputation is a function for append hash of reputation information for PERSON
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - [peerkey, personid, PERSON key of setter, contentkey, type]
//   returns :
//     - response data
//     - whether error object or nil
func AddReputation(stub shim.ChaincodeStubInterface, args []string) (key, res string, err error) {
	// check parameter
	if err = cmn.CheckParam(args, 5); err != nil {
		return
	}
	// check if data is exists
	data, key, err := get_and_check(stub, args, 2)
	if err != nil {
		return
	}
	log.Debug(key)
	if data == nil {
		err = errors.New("data is not exists")
		return
	}
	// TODO:add check if setter is belong to sender peer
	// TODO:add check if reputation is already appended
	// put data into ledger
	rep := PersonReputation{
		Setter:  args[2],
		Content: args[3],
		Type:    args[4],
	}
	(*data).Reputation = append((*data).Reputation, rep)
	err = cmn.Put(stub, key, (*data))
	return
}

// RemoveReputation is a function for remove hash of reputation information for PERSON
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - [peerkey, personid, setter, content, type]
//   returns :
//     - response data
//     - whether error object or nil
func RemoveReputation(stub shim.ChaincodeStubInterface, args []string) (key, res string, err error) {
	// check parameter
	if err = cmn.CheckParam(args, 5); err != nil {
		return
	}
	// check if data is exists
	data, key, err := get_and_check(stub, args, 2)
	if err != nil {
		return
	}
	if data == nil {
		err = errors.New("data is not exists")
		return
	}
	log.Debug(key)
	// put object which is removed target reputation data into ledger
	for i, v := range (*data).Reputation {
		log.Debug(v.Setter + "," + args[2] + "," + v.Content + "," + args[3])
		if v.Setter == args[2] && v.Content == args[3] {
			(*data).Reputation = append((*data).Reputation[0:i], (*data).Reputation[i+1:]...)
			err = cmn.Put(stub, key, (*data))
			return
		}
	}
	err = errors.New("reputation is not exists in target person")
	return
}
