/*
Package cmn provides common functions for chaincode.
*/

package cmn

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Put is a function for put data info ledger
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     key - key of target data
//   returns :
//     - whether error object or nil
func Put(stub shim.ChaincodeStubInterface, key string, val interface{}) error {
	if val == nil {
		return errors.New("invalid param")
	}
	jsVal, err := json.Marshal(val)
	if err != nil {
		return err
	}
	fmt.Printf("KEY:%s\n", key)
	err = stub.PutState(key, []byte(jsVal))
	return err
}

// Get is a functioin for get data from ledger
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     key - key of target data
//   returns :
//     - target data
//     - whether error object or nil
func Get(stub shim.ChaincodeStubInterface, key string) (val interface{}, err error) {
	var jsVal []byte
	jsVal, err = stub.GetState(key)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsVal, &val)
	return
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

// Sha1Byte is a function for generate sha1 hash of target binary data
//   parameters :
//     stub - object for accessing ledgers from chaincode
//   returns :
//     - sha1 hash
func Sha1Byte(v []byte) []byte {
	h := sha1.New()
	h.Write(v)
	return h.Sum(nil)
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
