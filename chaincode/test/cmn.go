/*
Package test provides common module for testing.
*/
package test

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/kenmazsyma/soila/chaincode/cmn"
	"github.com/kenmazsyma/soila/chaincode/root"
	"testing"
)

func S2b(src []string) (ret [][]byte) {
	for _, v := range src {
		ret = append(ret, []byte(v))
	}
	return
}

func MakeParam(params ...string) (ret [][]byte) {
	ret = [][]byte{}
	for _, v := range params {
		if len(ret) == 0 {
			ret = append(ret, []byte(params[0])) // funcname
		} else {
			ret = append(ret, []byte(cmn.EncodeBase64([]byte(v))))
		}
	}
	return
}

func P2o(payload []byte) (ret []interface{}, err error) {
	err = json.Unmarshal(payload, &ret)
	return
}

func EncodeAll(src []interface{}) (ret []string, err error) {
	ret = []string{}
	for i := 0; i < len(src); i++ {
		bt, _ := cmn.DecodeBase64(src[0].(string))
		ret = append(ret, string(bt))
	}
	return
}

func CreateStub(invoke_list map[string]root.InvokeRoutineType) *shim.MockStub {
	cc := new(root.CC)
	cc.SetInvokeMap(invoke_list)
	stub := shim.NewMockStub("soila_test", cc)
	return stub
}

func CheckMessage(t *testing.T, res pb.Response, expect string) {
	if res.Message != expect {
		t.Errorf("\nexpect:%s\nactual:%s", expect, res.Message)
	}
}

func CheckStatus(t *testing.T, res pb.Response, expect int32) {
	if res.Status != expect {
		t.Errorf("\nexpect:%d\nactual:%d", expect, res.Status)
	}
}

func CheckPayload(t *testing.T, res pb.Response, expect []string) {
	ret, err := P2o(res.Payload)
	if err != nil {
		t.Errorf("\nerror raised when converting json to object\n%s\n", err.Error())
	}
	if len(ret) != len(expect) {
		t.Errorf("\nlength of payload\nexpect:%s\nactual:%s", len(expect), len(ret))
	}
	for i := 0; i < len(expect); i++ {
		decode, _ := cmn.DecodeBase64(ret[i].(string))
		actual := string(decode)
		if actual != expect[i] {
			t.Errorf("\nindex:%d\nexpect:%s\nactual:%s", i, expect[i], actual)
		}
	}
}
