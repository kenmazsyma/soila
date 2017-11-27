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
	"github.com/kenmazsyma/soila/chaincode/project"
)

type Token struct {
	Creator  string // key of creator's PROJECT[key]
	Id       string // id of token[key]
	Name     string // name of token
	DataHash string // hash of data of token
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

// get_and_check is a function for getting data of TOKEN
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     key - key of TOKEN
//   return :
//     - TOKEN object
//     - key
//     - whether error object or nil
func get_and_check(stub shim.ChaincodeStubInterface, pkey string) (rec Token, key string, err error) {
	// get data from ledger
	rec = Token{}
	js, err := stub.GetState(pkey)
	if err != nil {
		return
	}
	if js == nil {
		err = errors.New("Data is not found in ledger")
		return
	}
	// convert to TOKEN object
	err = json.Unmarshal(js, &rec)
	if err != nil {
		return
	}
	if own := project.IsOwn(stub, rec.Creator); !own {
		err = errors.New("Token is not owned by sender.")
		return
	}
	return
}

// Register is a function for registering PROJECT to ledger
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - [projectid, id, name, data]
//   return :
//     key - key value
//     res - response data
//     err - either error object or nil
func Register(stub shim.ChaincodeStubInterface, args []string) (key, res string, err error) {
	// check parameter
	if err = cmn.CheckParam(args, 4); err != nil {
		return
	}
	// get key of PROJECT
	projectkey := project.GetKeyInPeer(stub, args[0])
	if len(projectkey) == 0 {
		err = errors.New("project is not exist in sender's peer")
		return
	}
	// check if data is already exists
	key, err = cmn.VerifyForRegistration(stub, generateKey, []string{projectkey, args[1]})
	if err != nil {
		return
	}
	log.Debug(key)
	// hash of description
	datahash := cmn.Sha1(args[3])
	// put data into ledger
	data := Token{
		Creator:  projectkey,
		Id:       args[1],
		Name:     args[2],
		DataHash: datahash,
	}
	err = cmn.Put(stub, key, data)
	return
}

// Get is a function for getting TOKEN information from ledger
//   parameters :
//     stub - object of chaincode information
//     args - [tokenkey]
//  return :
//    - response data
//    - either error object or nil
func Get(stub shim.ChaincodeStubInterface, args []string) (key, res string, err error) {
	return cmn.Get(stub, args)
}

// Update is a function for updating TOKEN information
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - [tokenkey, name, data]
//   return :
//     key - key value
//     res - response data
//     err - either error object or nil
func Update(stub shim.ChaincodeStubInterface, args []string) (key, res string, err error) {
	// check parameter
	if err = cmn.CheckParam(args, 3); err != nil {
		return
	}
	// check if project which manages token is owned by sender
	data, key, err := get_and_check(stub, args[0])
	if err != nil {
		return
	}
	log.Debug(key)
	//TODO:Token information can be updated only in the case that token is not issued yet
	data.Name = args[1]
	data.DataHash = cmn.Sha1(args[2])
	err = cmn.Put(stub, key, data)
	return
}

// Remove is a function for removing TOKEN information
//   parameters :
//     stub - object of chaincode information
//     args - [tokenkey]
//  return :
//     key - key value
//     res - response data
//     err - either error object or nil
func Remove(stub shim.ChaincodeStubInterface, args []string) (key, res string, err error) {
	//TODO:Token information can be updated only in the case that token is not issued yet
	// check parameter
	if err = cmn.CheckParam(args, 3); err != nil {
		return
	}
	_, key, err = get_and_check(stub, args[0])
	if err != nil {
		return
	}
	log.Debug(key)
	err = cmn.Delete(stub, key)
	return
}
