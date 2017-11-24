/*
Package peer provides chaincode for managing peer data.
TODO: nessesary to implement logic for verification whether peer can be trusted
*/

package peer

import (
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

var KEY_TYPE = "PEER"

// genearteKey is a function for generating key from id of PEER
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     id - id of PERSON
//   return :
//     - key
//     - whether error object or nil
func generateKey(stub shim.ChaincodeStubInterface) (ret string, err error) {
	id := []byte{}
	if id, err = stub.GetCreator(); err != nil {
		return "", err
	}
	return stub.CreateCompositeKey(KEY_TYPE, []string{string(id)})
}

// Register is a function for registering peer informartion
//   parameters :
//     stub - object of chaincode information
//     args - [siteaddress, apiaddress]
//  return :
//    1:response data
//    2:either error object or nil
func Register(stub shim.ChaincodeStubInterface, args []string) (res string, err error) {
	if len(args) != 2 {
		return "", errors.New("Invalid Arguments")
	}
	log.Info("start:")
	key, err := generateKey(stub)
	if err != nil {
		return "", err
	}
	info := Peer{}
	// check if data is already exists
	if info.Hash, err = stub.GetCreator(); err != nil {
		return "", err
	}
	log.Debug(string(info.Hash))
	id, err := GetId(stub)
	if err != nil {
		return "", err
	}
	info.Id = id
	log.Debug(string(id))
	if _, err = stub.GetState(string(id)); err != nil {
		return "", err
	}
	log.Info("prev:")
	info.SiteAddress = args[0]
	info.ApiAddress = args[1]
	err = cmn.Put(stub, key, info)
	return "", err
}
