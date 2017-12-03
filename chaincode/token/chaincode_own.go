/*
Package token provides chaincode for managing TOKEN data.
TODO:nessesary to implement authentication logic
*/
package token

import (
	//	"encoding/json"
	//	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//	"github.com/kenmazsyma/soila/chaincode/cmn"
	//	"github.com/kenmazsyma/soila/chaincode/log"
	//	"github.com/kenmazsyma/soila/chaincode/peer"
	//	"github.com/kenmazsyma/soila/chaincode/project"
)

type TokenOwn struct {
	Owner string // owner(personkey, projectkey)
	Token string // tokenkey
	Count int    // number of token
}

const KEY_TOKENOWN = "TOKENOWN"

// Publish is a function for registering a data for publishing token
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     args - [projectid, tokenid, count]
//   return :
//     ret - return value
//     err - either error object or nil
func Publish(stub shim.ChaincodeStubInterface, args []string) (ret []interface{}, err error) {
	return
}
