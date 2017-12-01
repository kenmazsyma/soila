/*
Package cmn provides common functions for chaincode.
*/
package cmn

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	. "github.com/kenmazsyma/soila/chaincode/log"
)

// Put is a function for put data info ledger
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     key - key of target data
//   returns :
//     - whether error object or nil
func Put(stub shim.ChaincodeStubInterface, key string, val interface{}) error {
	D("check parameter")
	if val == nil {
		return errors.New("invalid param")
	}
	D("convert parameter to json")
	jsVal, err := json.Marshal(val)
	if err != nil {
		return err
	}
	D("put data to leder:%s", key)
	err = stub.PutState(key, []byte(jsVal))
	return err
}

// Delete is a function for delete data from ledger
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     key - key of target data
//   returns :
//     - whether error object or nil
func Delete(stub shim.ChaincodeStubInterface, key string) (err error) {
	return stub.DelState(key)
}

type FuncGenKey func(shim.ChaincodeStubInterface, []string) (string, error)

// VerifyForRegistration is a function for verifying if parameters is valid before registering.
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     genkey - function for generating key
//     args - target parameters for verify
//     nofElm - expected length of args
//   returns :
//     key - generated key
//     err - whether error object or nil
func VerifyForRegistration(stub shim.ChaincodeStubInterface, genkey FuncGenKey, args []string) (key string, err error) {
	D("generate key")
	key, err = genkey(stub, args)
	if err != nil {
		return
	}
	D("check if data is already exists")
	val, err := stub.GetState(key)
	if err != nil {
		return
	}
	if val != nil {
		err = errors.New("data is already exists.")
		return
	}
	return
}

// VerifyForUpdate is a function for verifying if parameters is valid before updating.
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - target parameters for verify
//     nofElm - expected length of args
//   returns :
//     ret - data got from ledger
//     err - whether error object or nil
func VerifyForUpdate(stub shim.ChaincodeStubInterface, args []string, nofElm int) (ret []byte, err error) {
	//	D("check count of parameters")
	//	if len(args) != nofElm {
	//		err = errors.New("Invalid Arguments")
	//		return
	//	}
	//	D("generate key")
	//	key, err = genkey(stub, args)
	//	if err != nil {
	//		return
	//	}
	//	D("check if data is already exists")
	//	ret, err = stub.GetState(key)
	//	if err != nil {
	//		return
	//	}
	//	if ret == nil {
	//		err = errors.New("data is not exists.")
	//		return
	//	}
	//	return
	D("check parameter")
	if err = CheckParam(args, nofElm); err != nil {
		return
	}
	D("check if data is exists")
	ret, err = stub.GetState(args[0])
	if err != nil {
		return
	}
	if len(ret) == 0 {
		err = errors.New("data not found.")
		return
	}
	return
}

// Get is a function for getting data from ledger
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - target parameters for verify
//   returns :
//     key - key of data
//     res - data got from ledger
//     err - whether error obejct or nil
func Get(stub shim.ChaincodeStubInterface, args []string) (ret []interface{}, err error) {
	D("get data from ledger:%s", args[0])
	data, err := stub.GetState(args[0])
	if err != nil {
		return
	}
	if len(data) == 0 {
		err = errors.New("data not found.")
		return
	}
	ret = []interface{}{[]byte(args[0]), data}
	return
}

// Sha1 is a function for generate sha1 hash of target string
//   parameters :
//     stub - object for accessing ledgers from chaincode
//   returns :
//     - sha1 hash
func Sha1(v string) string {
	h := sha1.New()
	h.Write([]byte(v))
	return hex.EncodeToString(h.Sum(nil))
}

// Sha1Ar is a function for generate sha1 hash of target string
//   parameters :
//     stub - object for accessing ledgers from chaincode
//   returns :
//     - sha1 hash
func Sha1Ar(v []string) string {
	h := sha1.New()
	for ix, vv := range v {
		if ix == 0 {
			h.Write([]byte(vv))
		} else {
			h.Sum([]byte(vv))
		}
	}
	return hex.EncodeToString(h.Sum(nil))
}

// Sha1B is a function for generate sha1 hash of target binary data
//   parameters :
//     stub - object for accessing ledgers from chaincode
//   returns :
//     - sha1 hash
func Sha1B(v []byte) string {
	h := sha1.New()
	h.Write(v)
	return hex.EncodeToString(h.Sum(nil))
}

// ToJSON is a function for generating json string of target object
//   parameters :
//     o - target object
//   returns :
//     - json string
//     - whether error object or nil
func ToJSON(o interface{}) (string, error) {
	data, err := json.Marshal(o)
	return string(data), err
}

func CheckParam(prm []string, validlen int) error {
	if len(prm) != validlen {
		return errors.New("length of parameter is not valid.")
	}
	return nil
}

func DecodeBase64(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}

func EncodeBase64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}
