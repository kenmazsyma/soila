/*
Package token provides chaincode for managing TOKEN data.
*/
package token

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/kenmazsyma/soila/chaincode/cmn"
	. "github.com/kenmazsyma/soila/chaincode/log"
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
	return stub.CreateCompositeKey(KEY_TOKEN, []string{cmn.Sha1(args[0]), args[1]})
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
	D("get data from ledger")
	rec = Token{}
	js, err := stub.GetState(pkey)
	if err != nil {
		return
	}
	if js == nil {
		err = errors.New("Data is not found in ledger")
		return
	}
	D("convert to TOKEN object")
	err = json.Unmarshal(js, &rec)
	if err != nil {
		return
	}
	own, err := project.IsOwn(stub, rec.Creator)
	if err != nil {
		return
	}
	if !own {
		err = errors.New("Sender not own specifying token.")
		return
	}
	return
}

// Register is a function for registering PROJECT to ledger
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - [projectkey, id, name, data]
//   return :
//     ret - return value
//     err - either error object or nil
func Register(stub shim.ChaincodeStubInterface, args []string) (ret []interface{}, err error) {
	D("check parameter")
	if err = cmn.CheckParam(args, 4); err != nil {
		return
	}
	D("check if sender own specifying project")
	own, err := project.IsOwn(stub, args[0])
	if err != nil {
		return
	}
	if !own {
		err = errors.New("specifying project is not owned by sender.")
		return
	}
	D("check if data already exists")
	key, err := cmn.VerifyForRegistration(stub, generateKey, []string{args[0], args[1]})
	if err != nil {
		return
	}
	D("generate hash of 'data' member:%s", key)
	datahash := cmn.Sha1(args[3])
	D("put data into ledger")
	data := Token{
		Creator:  args[0],
		Id:       args[1],
		Name:     args[2],
		DataHash: datahash,
	}
	err = cmn.Put(stub, key, data)
	D("!!!!KEY:%s", key)
	ret = []interface{}{[]byte(key)}
	return
}

// Get is a function for getting TOKEN information from ledger
//   parameters :
//     stub - object of chaincode information
//     args - [tokenkey]
//  return :
//     ret - return value
//     err - either error object or nil
func Get(stub shim.ChaincodeStubInterface, args []string) (ret []interface{}, err error) {
	return cmn.Get(stub, args)
}

// Update is a function for updating TOKEN information
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - [tokenkey, name, data]
//   return :
//     ret - return value
//     err - either error object or nil
func Update(stub shim.ChaincodeStubInterface, args []string) (ret []interface{}, err error) {
	D("check parameter")
	if err = cmn.CheckParam(args, 3); err != nil {
		return
	}
	D("check if project which manages token is owned by sender")
	data, key, err := get_and_check(stub, args[0])
	if err != nil {
		return
	}
	D("put data into ledger:%s", key)
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
//     ret - return value
//     err - either error object or nil
func Remove(stub shim.ChaincodeStubInterface, args []string) (ret []interface{}, err error) {
	//TODO:Token information can be updated only in the case that token is not issued yet
	D("check parameter")
	if err = cmn.CheckParam(args, 3); err != nil {
		return
	}
	_, key, err := get_and_check(stub, args[0])
	if err != nil {
		return
	}
	D("delete data from ledger:%s", key)
	err = cmn.Delete(stub, key)
	return
}
