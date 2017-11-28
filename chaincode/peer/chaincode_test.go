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

// the case number of parameters is invalid
func Test_Register1(t *testing.T) {
	valid := "length of parameter is not valid."
	stub := CreateStub(invoke_list)
	stub.SetCreator([]byte("abcdef0123456789"))
	res := stub.MockInvoke("1", S2b([]string{"peer.register"}))
	CheckMessage(t, res, valid)
	res = stub.MockInvoke("1", S2b([]string{"peer.register", "1", "2"}))
	CheckMessage(t, res, valid)
}

// check if data duplicate
func Test_Register2(t *testing.T) {
	valid := "data is already exists."
	stub := CreateStub(invoke_list)
	stub.SetCreator([]byte("abcdef0123456789"))
	res := stub.MockInvoke("1", S2b([]string{"peer.register", "1"}))
	res = stub.MockInvoke("1", S2b([]string{"peer.register", "2"}))
	CheckMessage(t, res, valid)
}

// success
func Test_Register3(t *testing.T) {
	stub := CreateStub(invoke_list)
	stub.SetCreator([]byte("abcdef0123456789"))
	res := stub.MockInvoke("1", S2b([]string{"peer.register", "1"}))
	stub.SetCreator([]byte("2bcdef0123456789"))
	res = stub.MockInvoke("1", S2b([]string{"peer.register", "2"}))
	CheckStatus(t, res, 200)
}

// check if function can get data correctly
func Test_Get1(t *testing.T) {
	peer1 := []byte("abcdef0123456789")
	peer2 := []byte("bbcdef0123456789")
	stub := CreateStub(invoke_list)
	stub.SetCreator(peer1)
	res := stub.MockInvoke("1", S2b([]string{"peer.register", "1"}))
	ret, err := P2o(res.Payload)
	if err != nil {
		t.Errorf(err.Error())
	} else {
		bt, _ := DecodeBase64(ret[0].(string))
		fmt.Println(string(bt))
	}
	fmt.Println(string(res.Payload))
	stub.SetCreator(peer2)
	res = stub.MockInvoke("1", S2b([]string{"peer.register", "2"}))
	stub.SetCreator(peer1)

	// pending
	res = stub.MockInvoke("1", S2b([]string{"peer.get", "2"}))
	CheckStatus(t, res, 200)

}
