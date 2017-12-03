// Before testing this case, it is necessary to run testtool/prepare.sh
// for replacing MockStub module in fabric project

package token

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/kenmazsyma/soila/chaincode/peer"
	"github.com/kenmazsyma/soila/chaincode/project"
	"github.com/kenmazsyma/soila/chaincode/root"
	. "github.com/kenmazsyma/soila/chaincode/test"
	"os"
	"testing"
)

// ===============================
// Test environment
// ===============================

var invoke_list = map[string]root.InvokeRoutineType{
	"token.register":       Register,
	"token.update":         Update,
	"token.remove":         Remove,
	"project.register":     project.Register,
	"project.get":          project.Get,
	"project.updatestatus": project.UpdateStatus,
	"peer.register":        peer.Register,
	"peer.update":          peer.Update,
	"peer.get":             peer.Get,
	"peer.deregister":      peer.Deregister,
}

func initialize() {
	fmt.Println("init")
}

func terminate() {
	fmt.Println("term")
}

func TestMain(m *testing.M) {
	initialize()
	retCode := m.Run()
	terminate()
	os.Exit(retCode)
}

// ===================================
// sub routine
// ===================================
func getKeyFromPayload(res pb.Response) (key string, err error) {
	fmt.Printf("\nPAYLOAD:%s\n", res.Payload)
	o, err := UnmarshalPayload(res.Payload)
	if err != nil {
		return
	}
	if len(o) < 1 {
		err = errors.New("number of response is not correct.")
		return
	}
	key = o[0].(string)
	return
}

func prepare(stub *shim.MockStub) (key string, prjkey string, err error) {
	res := stub.MockInvoke("1", MakeParam("peer.register", "1"))
	if res.Status != 200 {
		err = errors.New("failed to register PEER")
		return
	}
	if key, err = getKeyFromPayload(res); err != nil {
		return
	}
	fmt.Printf("PEERKEY:%s\n", key)
	res = stub.MockInvoke("1", MakeParam("project.register", "12345"))
	if prjkey, err = getKeyFromPayload(res); err != nil {
		return
	}
	fmt.Printf("PROJECTKEY:%s\n", prjkey)
	return
}

// ===================================
// Test Case
// ===================================
func Test_Register(t *testing.T) {
	peer1 := []byte("abcdef0123456789")
	stub := CreateStub(invoke_list)
	// prepare a premise data
	stub.SetCreator(peer1)
	_, prjkey, err := prepare(stub)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	res := stub.MockInvoke("1", MakeParam("token.register", prjkey, "12345", "test", "aaaaa"))
	CheckStatus("a-1", t, res, 200)
	tokenkey, err := getKeyFromPayload(res)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	fmt.Printf("\nTOKENKEY:%s\n", tokenkey)
}
