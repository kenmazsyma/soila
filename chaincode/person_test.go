package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"testing"
)

func s2b(src []string) (ret [][]byte) {
	for _, v := range src {
		ret = append(ret, []byte(v))
	}
	return
}

func TestPerson_Init(t *testing.T) {
	scc := new(CC)
	stub := shim.NewMockStub("soila_test", scc)
	checkInit(t, stub, []string{})

}

func TestPerson_Put(t *testing.T) {
	scc := new(CC)
	stub := shim.NewMockStub("soila_test", scc)
	checkInvoke(t, stub, []string{"person.put", "data1", "{\"a\":1}"})
}

func TestPerson_Update(t *testing.T) {
	scc := new(CC)
	stub := shim.NewMockStub("soila_test", scc)
	stub.MockInvoke("1", s2b([]string{"person.put", "data1", "{\"a\":1}"}))
	checkInvoke(t, stub, []string{"person.update", "data1", "{\"a\":2}"})
}

func TestPerson_Get(t *testing.T) {
	scc := new(CC)
	stub := shim.NewMockStub("soila_test", scc)
	stub.MockInvoke("1", s2b([]string{"person.put", "data1", "{\"a\":1}"}))
	stub.MockInvoke("1", s2b([]string{"person.update", "data1", "{\"a\":2}"}))
	checkQuery(t, stub, []string{"person.get", "data1"}, "test")
}

/** Sub Routine **/
func checkInit(t *testing.T, stub *shim.MockStub, args []string) {
	res := stub.MockInit("1", s2b(args))
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args []string) {
	res := stub.MockInvoke("1", s2b(args))
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, args []string, expect string) {
	//	args = append([]string{"query"}, args...)
	res := stub.MockInvoke("1", s2b(args))
	if res.Status != shim.OK {
		fmt.Println("Query", args[0], "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", args[0], "failed to get value")
		t.FailNow()
	}
	if string(res.Payload) != expect {
		fmt.Println("Query value", string(res.Payload), "was not", expect, "as expected")
		t.FailNow()
	}
}
