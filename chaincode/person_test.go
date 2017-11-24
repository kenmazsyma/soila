package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/kenmazsyma/soila/chaincode/cmn"
	"testing"
)

func s2b(src []string) (ret [][]byte) {
	for _, v := range src {
		ret = append(ret, []byte(v))
	}
	return
}

var testval = []string{
	"{\"a\":1}",
	"{\"a\":2}",
}

const code = "data1"

func TestPerson_Init(t *testing.T) {
	scc := new(CC)
	stub := shim.NewMockStub("soila_test", scc)
	checkInit(t, stub, []string{})

}

func TestPerson_Put(t *testing.T) {
	scc := new(CC)
	stub := shim.NewMockStub("soila_test", scc)
	checkInvoke(t, stub, []string{"person.register", code, testval[0]})
}

func TestPerson_Update(t *testing.T) {
	scc := new(CC)
	stub := shim.NewMockStub("soila_test", scc)
	stub.MockInvoke("1", s2b([]string{"person.register", code, testval[0]}))
	checkInvoke(t, stub, []string{"person.update", code, testval[1]})
}

func TestPerson_Get(t *testing.T) {
	scc := new(CC)
	stub := shim.NewMockStub("soila_test", scc)
	stub.MockInvoke("1", s2b([]string{"person.register", code, testval[0]}))
	stub.MockInvoke("1", s2b([]string{"person.update", code, testval[1]}))
	expect := fmt.Sprintf("{\"Ver\":[\"%s\",\"%s\"],\"Activity\":[],\"Reputation\":[]}",
		cmn.Sha1(testval[0]), cmn.Sha1(testval[1]))
	checkQuery(t, stub, []string{"person.get", code}, expect)
}

func TestPerson_AddActivity(t *testing.T) {
	scc := new(CC)
	stub := shim.NewMockStub("soila_test", scc)
	stub.MockInvoke("1", s2b([]string{"person.register", code, testval[0]}))
	checkInvoke(t, stub, []string{"person.add_activity", code, "a"})
	checkInvoke(t, stub, []string{"person.add_activity", code, "b"})
	expect := fmt.Sprintf("{\"Ver\":[\"%s\"],\"Activity\":[\"a\",\"b\"],\"Reputation\":[]}",
		cmn.Sha1(testval[0]))
	checkQuery(t, stub, []string{"person.get", code}, expect)
}

func TestPerson_AddReputation(t *testing.T) {
	scc := new(CC)
	stub := shim.NewMockStub("soila_test", scc)
	stub.MockInvoke("1", s2b([]string{"person.register", code, testval[0]}))
	checkInvoke(t, stub, []string{"person.add_reputation", code, "a", "b", "1"})
	checkInvoke(t, stub, []string{"person.add_reputation", code, "ccc", "d", "2"})
	expect := fmt.Sprintf("{\"Ver\":[\"%s\"],\"Activity\":[],\"Reputation\":[{\"Setter\":\"a\",\"Content\":\"b\",\"Type\":\"1\"},{\"Setter\":\"ccc\",\"Content\":\"d\",\"Type\":\"2\"}]}",
		cmn.Sha1(testval[0]))
	checkQuery(t, stub, []string{"person.get", code}, expect)
}

func TestPerson_RemoveReputation(t *testing.T) {
	scc := new(CC)
	stub := shim.NewMockStub("soila_test", scc)
	stub.MockInvoke("1", s2b([]string{"person.register", code, testval[0]}))
	checkInvoke(t, stub, []string{"person.add_reputation", code, "a", "b", "1"})
	checkInvoke(t, stub, []string{"person.add_reputation", code, "c", "d", "2"})
	checkInvoke(t, stub, []string{"person.add_reputation", code, "e", "f", "3"})
	checkInvoke(t, stub, []string{"person.add_reputation", code, "g", "h", "4"})
	checkInvoke(t, stub, []string{"person.remove_reputation", code, "e", "f"})
	expect := fmt.Sprintf("{\"Ver\":[\"%s\"],\"Activity\":[],\"Reputation\":[{\"Setter\":\"a\",\"Content\":\"b\",\"Type\":\"1\"},{\"Setter\":\"c\",\"Content\":\"d\",\"Type\":\"2\"},{\"Setter\":\"g\",\"Content\":\"h\",\"Type\":\"4\"}]}",
		cmn.Sha1(testval[0]))
	checkQuery(t, stub, []string{"person.get", code}, expect)
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
