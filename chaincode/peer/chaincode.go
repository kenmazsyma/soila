/*
Package peer provides chaincode for managing PEER data.
TODO: nessesary to implement logic for verification whether peer can be trusted
*/

package peer

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/kenmazsyma/soila/chaincode/cmn"
	"github.com/kenmazsyma/soila/chaincode/log"
)

type Peer struct {
	Hash        []byte
	Id          []byte
	SiteAddress string
	ApiAddress  string
}

const KEY_TYPE = "PEER"

// genearteKey is a function for generating key from id of PEER
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - arguments which contains keys
//   return :
//     - key
//     - whether error object or nil
func generateKey(stub shim.ChaincodeStubInterface, args []string) (ret string, err error) {
	//if id, err = stub.GetCreator(); err != nil {
	//	return "", err
	//}
	return stub.CreateCompositeKey(KEY_TYPE, []string{args[0]})
}

func generateKeyFromId(stub shim.ChaincodeStubInterface, id []byte) (ret string, err error) {
	return stub.CreateCompositeKey(KEY_TYPE, []string{string(id)})
}

// Register is a function for registering PEER informartion
//   parameters :
//     stub - object of chaincode information
//     args - [siteaddress, apiaddress]
//  return :
//    - response data
//    - either error object or nil
func Register(stub shim.ChaincodeStubInterface, args []string) (res string, err error) {
	log.Info("start:")
	info := Peer{}
	// check if data is already exists
	if info.Hash, err = stub.GetCreator(); err != nil {
		return "", err
	}
	log.Debug(string(info.Hash))
	key, err := cmn.VerifyForRegistration(stub, generateKey, []string{string(info.Hash)}, 1)
	if err != nil {
		return "", err
	}
	log.Debug(key)
	id, err := GetId(stub)
	if err != nil {
		return "", err
	}
	info.Id = id
	log.Debug(string(id))
	if _, err = stub.GetState(key); err != nil {
		return "", err
	}
	log.Info("prev:")
	info.SiteAddress = args[0]
	info.ApiAddress = args[1]
	err = cmn.Put(stub, key, info)
	return "", err
}

// Get is a function for getting PEER information from ledger
//   parameters :
//     stub - object of chaincode information
//     args - [id]
//  return :
//    - response data
//    - either error object or nil
func Get(stub shim.ChaincodeStubInterface, args []string) (res string, err error) {
	return cmn.Get(stub, args, 1)
}

// Update is a function for updating PEER information
//   parameters :
//     stub - object of chaincode information
//     args - [id, siteaddress, apiaddress]
//  return :
//    - response data
//    - either error object or nil
func Update(stub shim.ChaincodeStubInterface, args []string) (res string, err error) {
	res = ""
	if len(args) != 1 {
		err = errors.New("Invalid Arguments")
		return
	}
	log.Info("start:")
	key, err := generateKeyFromId(stub, []byte(args[0]))
	if err != nil {
		return
	}
	log.Debug(key)
	val, err := stub.GetState(key)
	if err != nil {
		return
	}
	log.DebugB(val)
	data := Peer{}
	err = json.Unmarshal(val, &data)
	if err != nil {
		return
	}
	log.DebugB(data.Id)
	valid, err := CompareId(stub, data.Id)
	if err != nil {
		return
	}
	if !valid {
		err = errors.New("Peer is not owned by sender")
		return
	}
	res, err = cmn.ToJSON(data)
	return
}

// Remove is a function for updating PEER information
//   parameters :
//     stub - object of chaincode information
//     args - [id, siteaddress, apiaddress]
//  return :
//    - response data
//    - either error object or nil
func Remove(stub shim.ChaincodeStubInterface, args []string) (res string, err error) {
	res = ""
	if len(args) != 1 {
		err = errors.New("Invalid Arguments")
		return
	}
	log.Info("start:")
	key, err := generateKeyFromId(stub, []byte(args[0]))
	if err != nil {
		return
	}
	log.Debug(key)
	val, err := stub.GetState(key)
	if err != nil {
		return
	}
	log.DebugB(val)
	data := Peer{}
	err = json.Unmarshal(val, &data)
	if err != nil {
		return
	}
	log.DebugB(data.Id)
	valid, err := CompareId(stub, data.Id)
	if err != nil {
		return
	}
	if !valid {
		err = errors.New("Peer is not owned by sender")
		return
	}
	err = cmn.Delete(stub, key)
	return
}
