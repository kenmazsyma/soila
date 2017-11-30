// Before testing this case, it is necessary to run testtool/prepare.sh
// for replacing MockStub module in fabric project

package person

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/kenmazsyma/soila/chaincode/peer"
	"github.com/kenmazsyma/soila/chaincode/root"
	. "github.com/kenmazsyma/soila/chaincode/test"
	"os"
	"testing"
)

// ===============================
// Test environment
// ===============================

var invoke_list = map[string]root.InvokeRoutineType{
	"person.register":          Register,
	"person.update":            Update,
	"person.get":               Get,
	"person.add_activity":      AddActivity,
	"person.add_reputation":    AddReputation,
	"person.remove_reputation": RemoveReputation,
	"peer.register":            peer.Register,
	"peer.update":              peer.Update,
	"peer.get":                 peer.Get,
	"peer.deregister":          peer.Deregister,
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
func RegPeer(stub *shim.MockStub) string {
	res := stub.MockInvoke("1", MakeParam("peer.register", "1"))
	if res.Status != 200 {
		return "failed to register PEER"
	}
	return "OK"
}

// ===================================
// Test Case
// ===================================

func Test_Register(t *testing.T) {
	peer1 := []byte("abcdef0123456789")
	peer2 := []byte("bbcdef0123456789")
	stub := CreateStub(invoke_list)
	stub.SetCreator(peer1)
	rslt := RegPeer(stub)
	if rslt != "OK" {
		t.Errorf(rslt)
	}
	CASE("a-2")
	res := stub.MockInvoke("1", MakeParam("person.register", "12345"))
	CheckStatus("a-2", t, res, 500)
	CheckMessage("a-2", t, res, "length of parameter is not valid.")
	CASE("a-3")
	res = stub.MockInvoke("1", MakeParam("person.register", "12345", "23456", "123456"))
	CheckStatus("a-3", t, res, 500)
	CheckMessage("a-3", t, res, "length of parameter is not valid.")
	CASE("a-1")
	res = stub.MockInvoke("1", MakeParam("person.register", "12345", "12345"))
	CheckStatus("a-1", t, res, 200)
	ret, _ := P2o(res.Payload)
	v, _ := EncodeAll(ret)
	CheckPayload("a-1", t, res, []interface{}{v[0]})
	CASE("a-4")
	res = stub.MockInvoke("1", MakeParam("person.register", "12345", "12345"))
	CheckStatus("a-4", t, res, 500)
	CheckMessage("a-4", t, res, "data is already exists.")

	CASE("b-2")
	res = stub.MockInvoke("1", MakeParam("person.update", v[0]))
	CheckStatus("b-2", t, res, 500)
	CheckMessage("b-2", t, res, "length of parameter is not valid.")
	CASE("b-3")
	res = stub.MockInvoke("1", MakeParam("person.update", v[0], "23456", "123456"))
	CheckStatus("b-3", t, res, 500)
	CheckMessage("b-3", t, res, "length of parameter is not valid.")
	CASE("b-1")
	res = stub.MockInvoke("1", MakeParam("person.update", v[0], "UPDATE!!"))
	CheckStatus("b-1", t, res, 200)
	CheckPayload("b-1", t, res, []interface{}{v[0]})
	CASE("c-1, b-1")
	res = stub.MockInvoke("1", MakeParam("person.get", v[0]))
	CheckStatus("c-1, b-1", t, res, 200)
	CASE("b-6")
	res = stub.MockInvoke("1", MakeParam("person.update", v[0], "UPDATE!!"))
	CheckStatus("b-6", t, res, 200)
	CheckPayload("b-6", t, res, []interface{}{v[0]})
	res = stub.MockInvoke("1", MakeParam("person.get", v[0]))
	CheckStatus("b-6", t, res, 200)
	CASE("b-4")
	res = stub.MockInvoke("1", MakeParam("person.update", "123", "UPDATE!!"))
	CheckStatus("b-4", t, res, 500)
	CheckMessage("b-4", t, res, "data not found.")
	CASE("b-5")
	stub.SetCreator(peer2)
	res = stub.MockInvoke("1", MakeParam("person.update", v[0], "UPDATE!!"))
	CheckStatus("b-5", t, res, 200)
	CheckMessage("b-5", t, res, "data not found.")
}
