// Before testing this case, it is necessary to run testtool/prepare.sh
// for replacing MockStub module in fabric project

package project

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/kenmazsyma/soila/chaincode/peer"
	"github.com/kenmazsyma/soila/chaincode/person"
	"github.com/kenmazsyma/soila/chaincode/root"
	. "github.com/kenmazsyma/soila/chaincode/test"
	"os"
	"testing"
)

// ===============================
// Test environment
// ===============================

var invoke_list = map[string]root.InvokeRoutineType{
	"project.register":         Register,
	"project.get":              Get,
	"project.updatestatus":     UpdateStatus,
	"peer.register":            peer.Register,
	"peer.update":              peer.Update,
	"peer.get":                 peer.Get,
	"peer.deregister":          peer.Deregister,
	"person.register":          person.Register,
	"person.update":            person.Update,
	"person.get":               person.Get,
	"person.add_activity":      person.AddActivity,
	"person.add_reputation":    person.AddReputation,
	"person.remove_reputation": person.RemoveReputation,
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

func prepare(stub *shim.MockStub) (key string, personkey string, err error) {
	res := stub.MockInvoke("1", MakeParam("peer.register", "1"))
	if res.Status != 200 {
		err = errors.New("failed to register PEER")
		return
	}
	if key, err = getKeyFromPayload(res); err != nil {
		return
	}
	fmt.Printf("PEERKEY:%s", key)
	res = stub.MockInvoke("1", MakeParam("person.register", "12345", "12345"))
	if personkey, err = getKeyFromPayload(res); err != nil {
		return
	}
	fmt.Printf("PERSONKEY:%s", personkey)
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
	_, _, err := prepare(stub)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	CASE("a-1")
	res := stub.MockInvoke("1", MakeParam("project.register", "12345"))
	CheckStatus("a-1", t, res, 200)
	prjkey, err := getKeyFromPayload(res)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	fmt.Printf("\nPROJECTKEY:%s\n", prjkey)
}
