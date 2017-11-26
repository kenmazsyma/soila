/*
Package token provides chaincode for managing TOKEN data.
*/

package token

import ()

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
