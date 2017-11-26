/*
Package token provides chaincode for managing TOKEN data.
TODO:nessesary to implement authentication logic
*/

package token

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/kenmazsyma/soila/chaincode/cmn"
	"github.com/kenmazsyma/soila/chaincode/log"
	"github.com/kenmazsyma/soila/chaincode/peer"
	"github.com/kenmazsyma/soila/chaincode/project"
)

type TokenOwn struct {
	OwnerID string // owner(personid, projectid) of token[key]
	TokenID string // type of token[key]
	Count   int    // number of owning token
}

const KEY_TOKENOWN = "TOKENOWN"

// generateKeyForOwn is a function for generating key from id of PROJECT
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - arguments which contains key
//   return :
//     - key
//     - whether error object or nil
func generateKeyForOwn(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	return stub.CreateCompositeKey(KEY_TOKENTRANS, args[0:2])
}

// Publish is a function for registering a data for publishing token
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - [projectid, tokenid, count]
//   return :
//     - response data
//     - error object if error occured
func Publish(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("invalid param")
	}
	peerid, err := peer.GetId(stub)
	if err != nil {
		return "", err
	}
	token, tokenkey, err := get_and_check(stub, []string{peerid, args[0]}, 2)
	if err != nil {
		return "", err
	}
	if token == nil {
		return "", errors.New("Project is not found")
	}
	// TODO:
	return "", err

}
