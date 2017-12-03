/*
Package project provdes chaincode for managing PROJECT data.
*/
package project

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	. "github.com/kenmazsyma/soila/chaincode/log"
	"github.com/kenmazsyma/soila/chaincode/peer"
)

// IsOwn is a function for checking if project exists in sender's peer
//   parameters :
//     stub - object for accessing ledgers from chaincode
//     prjkey - projectkey
//   return :
//     rslt - true if project is not owned by sender's peer
//     error - whether error object or nil
func IsOwn(stub shim.ChaincodeStubInterface, prjkey string) (rslt bool, err error) {
	D("get PEER key of sender")
	rslt = false
	peerkey, err := peer.GetKey(stub)
	if err != nil {
		return
	}
	D("get PROJECT data from ledger")
	data, err := Get(stub, []string{prjkey})
	if err != nil {
		return
	}
	if data == nil || len(data) < 2 {
		err = errors.New("data not found.")
		return
	}
	prj := Project{}
	if err = json.Unmarshal(data[1].([]byte), &prj); err != nil {
		return
	}
	rslt = (peerkey == prj.PeerKey)
	return
}
