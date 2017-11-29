// Before testing this case, it is necessary to run testtool/prepare.sh
// for replacing MockStub module in fabric project

package peer

import (
	"fmt"
	"github.com/kenmazsyma/soila/chaincode/cmn"
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

// number of parameters is invalid
func Test_Register1(t *testing.T) {
	valid := "length of parameter is not valid."
	stub := CreateStub(invoke_list)
	stub.SetCreator([]byte("abcdef0123456789"))
	res := stub.MockInvoke("1", MakeParam("peer.register"))
	CheckMessage(t, res, valid)
	res = stub.MockInvoke("1", MakeParam("peer.register", "1", "2"))
	CheckMessage(t, res, valid)
}

// check if data duplicate
func Test_Register2(t *testing.T) {
	valid := "data is already exists."
	stub := CreateStub(invoke_list)
	stub.SetCreator([]byte("abcdef0123456789"))
	res := stub.MockInvoke("1", MakeParam("peer.register", "1"))
	res = stub.MockInvoke("1", MakeParam("peer.register", "2"))
	CheckMessage(t, res, valid)
}

// success
func Test_Register3(t *testing.T) {
	stub := CreateStub(invoke_list)
	stub.SetCreator([]byte("abcdef0123456789"))
	res := stub.MockInvoke("1", MakeParam("peer.register", "1"))
	stub.SetCreator([]byte("2bcdef0123456789"))
	res = stub.MockInvoke("1", MakeParam("peer.register", "2"))
	CheckStatus(t, res, 200)
}

// check if function can get data correctly
func Test_Get1(t *testing.T) {
	peer1 := []byte("abcdef0123456789")
	peer2 := []byte("bbcdef0123456789")
	stub := CreateStub(invoke_list)
	// register 1st peer
	stub.SetCreator(peer1)
	res := stub.MockInvoke("1", MakeParam("peer.register", "1"))
	ret, err := P2o(res.Payload)
	bt := []byte{}
	if err != nil {
		t.Errorf(err.Error())
	} else {
		bt, _ = cmn.DecodeBase64(ret[0].(string))
		fmt.Println(string(bt))
	}
	// register 2nd peer
	stub.SetCreator(peer2)
	res = stub.MockInvoke("1", MakeParam("peer.register", "2"))
	// get 1st data
	stub.SetCreator(peer1)
	res = stub.MockInvoke("1", MakeParam("peer.get", string(bt)))
	CheckStatus(t, res, 200)
	ret, err = P2o(res.Payload)
	expect := "{\"Hash\":\"d80e5e55dd4128844827a53d7363045485f08751\",\"Address\":\"1\"}"
	CheckPayload(t, res, []string{string(bt), expect})
}

// number of parameters is not correct(1)
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
	// deregister peer
	res = stub.MockInvoke("1", MakeParam("peer.deregister"))
	CheckStatus(t, res, 500)
	CheckMessage(t, res, "length of parameter is not valid.")
}

// number of parameters is not correct(2)
func Test_Deregister2(t *testing.T) {
	peer1 := []byte("abcdef0123456789")
	stub := CreateStub(invoke_list)
	// register peer
	stub.SetCreator(peer1)
	res := stub.MockInvoke("1", MakeParam("peer.register", "1"))
	_, err := P2o(res.Payload)
	if err != nil {
		t.Errorf("register failed : %s", err.Error())
	}
	// deregister peer
	res = stub.MockInvoke("1", MakeParam("peer.deregister", "1", "2"))
	CheckStatus(t, res, 500)
	CheckMessage(t, res, "length of parameter is not valid.")
}

// deregister from another peer
func Test_Deregister3(t *testing.T) {
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
	// deregister peer
	stub.SetCreator(peer2)
	res = stub.MockInvoke("1", MakeParam("peer.deregister", v[0]))
	CheckStatus(t, res, 500)
	CheckMessage(t, res, "Peer is not owned by sender")
}
