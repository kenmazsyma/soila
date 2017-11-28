/*
Package test provides common module for testing.
*/
package test

import (
	"encoding/base64"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/kenmazsyma/soila/chaincode/root"
	"testing"
)

func S2b(src []string) (ret [][]byte) {
	for _, v := range src {
		ret = append(ret, []byte(v))
	}
	return
}

func P2o(payload []byte) (ret []interface{}, err error) {
	err = json.Unmarshal(payload, &ret)
	return
}

func DecodeBase64(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}

func CreateStub(invoke_list map[string]root.InvokeRoutineType) *shim.MockStub {
	cc := new(root.CC)
	cc.SetInvokeMap(invoke_list)
	stub := shim.NewMockStub("soila_test", cc)
	return stub
}

func CheckMessage(t *testing.T, res pb.Response, valid string) {
	if res.Message != valid {
		t.Errorf("\nexpect:%s\nactual:%s", valid, res.Message)
	}
}

func CheckStatus(t *testing.T, res pb.Response, valid int32) {
	if res.Status != valid {
		t.Errorf("\nexpect:%d\nactual:%d", valid, res.Status)
	}
}
