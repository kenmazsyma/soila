// Before testing this case, it is necessary to run testtool/prepare.sh
// for replacing MockStub module in fabric project

package peer

import (
	"fmt"
	"github.com/kenmazsyma/soila/chaincode/root"
	. "github.com/kenmazsyma/soila/chaincode/test"
	"os"
	"testing"
)

// ===============================
// Test environment
// ===============================

var invoke_list = map[string]root.InvokeRoutineType{
	"peer.register":   Register,
	"peer.update":     Update,
	"peer.get":        Get,
	"peer.deregister": Deregister,
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
// Test Case
// ===================================

func Test_Register1(t *testing.T) {
	valid := "number of parameter is not valid."
	stub := CreateStub(invoke_list)
	CASE("a-2")
	stub.SetCreator([]byte("abcdef0123456789"))
	res := stub.MockInvoke("1", MakeParam("peer.register"))
	CheckMessage("a-2", t, res, valid)
	CASE("a-3")
	res = stub.MockInvoke("1", MakeParam("peer.register", "1", "2"))
	CheckMessage("a-3", t, res, valid)
}

func Test_Register2(t *testing.T) {
	valid := "data already exists."
	stub := CreateStub(invoke_list)
	CASE("a-4")
	stub.SetCreator([]byte("abcdef0123456789"))
	res := stub.MockInvoke("1", MakeParam("peer.register", "1"))
	res = stub.MockInvoke("1", MakeParam("peer.register", "2"))
	CheckMessage("a-4", t, res, valid)
}

func Test_Get(t *testing.T) {
	peer1 := []byte("abcdef0123456789")
	peer2 := []byte("bbcdef0123456789")
	stub := CreateStub(invoke_list)
	// register 1st peer
	CASE("a-1")
	stub.SetCreator(peer1)
	res := stub.MockInvoke("1", MakeParam("peer.register", "1"))
	CheckStatus("a-1", t, res, 200)
	ret, _ := P2o(res.Payload)
	v, _ := EncodeAll(ret)
	// register 2nd peer
	CASE("b-1")
	stub.SetCreator(peer2)
	res = stub.MockInvoke("1", MakeParam("peer.register", "2"))
	// get 1st data
	stub.SetCreator(peer1)
	res = stub.MockInvoke("1", MakeParam("peer.get", v[0]))
	CheckStatus("b-1", t, res, 200)
	ret, _ = P2o(res.Payload)
	fmt.Printf("payload:%s\n", string(res.Payload))
	expect := "{\"hash\":\"cc6e41f3c71832a71d88c8bd833d1790e158c360698de61d9f957afa816b32048d8d1a6a2fb105f09130975ef3b596cadcdec8fed4f0a7793141b375460141c3\",\"address\":\"1\"}"
	CheckPayload("b-1", t, res, []interface{}{v[0], expect})
	// get data by key which is not registered
	CASE("b-2")
	res = stub.MockInvoke("1", MakeParam("peer.get", "123"))
	CheckStatus("b-2", t, res, 500)
}

func Test_Deregister1(t *testing.T) {
	peer1 := []byte("abcdef0123456789")
	stub := CreateStub(invoke_list)
	// register peer
	stub.SetCreator(peer1)
	res := stub.MockInvoke("1", MakeParam("peer.register", "1"))
	_, err := P2o(res.Payload)
	if err != nil {
		t.Errorf("register failed : %s", err.Error())
	}
	// deregister peer1
	CASE("d-2")
	res = stub.MockInvoke("1", MakeParam("peer.deregister", "1"))
	CheckStatus("d-2", t, res, 500)                                  // d-2
	CheckMessage("d-2", t, res, "number of parameter is not valid.") // d-2
}

func Test_Deregister2(t *testing.T) {
	peer1 := []byte("abcdef0123456789")
	stub := CreateStub(invoke_list)
	// register peer
	stub.SetCreator(peer1)
	res := stub.MockInvoke("1", MakeParam("peer.register", "1"))
	o, err := P2o(res.Payload)
	if err != nil {
		t.Errorf("register failed : %s", err.Error())
	}
	v, _ := EncodeAll(o)
	// deregister peer2
	CASE("d-1")
	stub.SetCreator(peer1)
	res = stub.MockInvoke("1", MakeParam("peer.deregister"))
	CheckStatus("d-1", t, res, 200) // d-1
	// verify if data is successfully deleted
	res = stub.MockInvoke("1", MakeParam("peer.get", v[0]))
	CheckStatus("d-1", t, res, 500)                // d-1
	CheckMessage("d-1", t, res, "data not found.") // d-1
	// deregister peer3
	CASE("d-3")
	res = stub.MockInvoke("1", MakeParam("peer.deregister"))
	CheckStatus("d-3", t, res, 500) // d-4
}

func Test_Update(t *testing.T) {
	peer1 := []byte("abcdef0123456789")
	peer2 := []byte("bbcdef0123456789")
	stub := CreateStub(invoke_list)
	// register peer
	stub.SetCreator(peer1)
	res := stub.MockInvoke("1", MakeParam("peer.register", "1"))
	o, err := P2o(res.Payload)
	if err != nil {
		t.Errorf("register failed : %s", err.Error())
	}
	v, _ := EncodeAll(o)
	// update1
	CASE("c-2")
	res = stub.MockInvoke("1", MakeParam("peer.update"))
	CheckStatus("c-2", t, res, 500)                                  // c-2
	CheckMessage("c-2", t, res, "number of parameter is not valid.") // c-2
	// update2
	CASE("c-3")
	res = stub.MockInvoke("1", MakeParam("peer.update", "127.0.0.1", "255.255.255.0"))
	CheckStatus("c-3", t, res, 500)                                  // c-3
	CheckMessage("c-3", t, res, "number of parameter is not valid.") // c-3
	// update3
	//CASE("c-4")
	//res = stub.MockInvoke("1", MakeParam("peer.update", "127.0.0.1"))
	//CheckStatus("c-4", t, res, 500)                // c-4
	//CheckMessage("c-4", t, res, "data not found.") // c-4
	// update4
	CASE("c-5")
	stub.SetCreator(peer2)
	res = stub.MockInvoke("1", MakeParam("peer.update", "127.0.0.1"))
	CheckStatus("c-5", t, res, 500)                // c-5
	CheckMessage("c-5", t, res, "data not found.") // c-5
	// update5
	CASE("c-1")
	stub.SetCreator(peer1)
	res = stub.MockInvoke("1", MakeParam("peer.update", "127.0.0.1"))
	CheckStatus("c-1(1)", t, res, 200)
	res = stub.MockInvoke("1", MakeParam("peer.get", v[0]))
	CheckStatus("c-1(2)", t, res, 200)
	fmt.Printf("payload:%s\n", string(res.Payload))
	expect := "{\"hash\":\"cc6e41f3c71832a71d88c8bd833d1790e158c360698de61d9f957afa816b32048d8d1a6a2fb105f09130975ef3b596cadcdec8fed4f0a7793141b375460141c3\",\"address\":\"127.0.0.1\"}"
	CheckPayload("c-1", t, res, []interface{}{v[0], expect}) // c-1
	// update6
	CASE("c-6")
	res = stub.MockInvoke("1", MakeParam("peer.update", "127.0.0.1"))
	CheckStatus("c-6", t, res, 200) // c-6
}
