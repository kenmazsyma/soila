/*
Package test provides common module for testing.
*/
package test

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/kenmazsyma/soila/chaincode/cmn"
	"github.com/kenmazsyma/soila/chaincode/root"
	"reflect"
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

func UnmarshalPayload(payload []byte) (ret []interface{}, err error) {
	if len(payload) == 0 {
		return
	}
	val, err := P2o(payload)
	if err != nil {
		return
	}
	enc, err := EncodeAll(val)
	if err != nil {
		return
	}
	for _, v := range enc {
		elm := map[string]interface{}{}
		if v[0] == '[' || v[0] == '{' {
			if err = json.Unmarshal([]byte(v), &elm); err != nil {
				return
			}
			ret = append(ret, elm)
		} else {
			ret = append(ret, v)
		}
	}
	return
}

func EncodeAll(src []interface{}) (ret []string, err error) {
	ret = []string{}
	for i := 0; i < len(src); i++ {
		bt, _ := cmn.DecodeBase64(src[i].(string))
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

func CheckMessage(cs string, t *testing.T, res pb.Response, expect string) {
	if res.Message != expect {
		t.Errorf("\n##%s##\nexpect:%s\nactual:%s", cs, expect, res.Message)
	}
}

func CheckStatus(cs string, t *testing.T, res pb.Response, expect int32) {
	if res.Status != expect {
		t.Errorf("\n##%s##\nexpect:%d\nactual:%d", cs, expect, res.Status)
	}
}

func CheckPayload(cs string, t *testing.T, res pb.Response, expect []interface{}) {
	if len(res.Payload) == 0 {
		if len(expect) != 0 {
			t.Errorf("\n##%s##\nPayload is nil.", cs)
			return
		}
	}
	ret, err := P2o(res.Payload)
	if err != nil {
		t.Errorf("\n##%s##\nerrored when converting json to object\n%s\n", cs, err.Error())
	}
	if len(ret) != len(expect) {
		t.Errorf("\n##%s##\nnumber of payload members\nexpect:%s\nactual:%s", cs, len(expect), len(ret))
	}
	for i := 0; i < len(expect); i++ {
		if ret[i] != nil {
			decode, _ := cmn.DecodeBase64(ret[i].(string))
			actual := string(decode)
			if actual != expect[i] {
				t.Errorf("\n##%s##\nindex:%d\nexpect:%s\nactual:%s", cs, i, expect[i], actual)
			}
		} else {
			if ret[i] != expect[i] {
				t.Errorf("\n##%s##\nindex:%d\nexpect:%s\nactual:%s", cs, i, expect[i], ret[i])
			}
		}
	}
}

func CheckPayloadMember(cs string, t *testing.T, v []interface{}, idx int, expect string, place []interface{}) {
	if len(v) <= idx {
		t.Errorf("\n##%s##\nnumber of payload members is not as expected", cs)
		return
	}
	vv := v[idx]
	for ix, o := range place {
		switch reflect.TypeOf(o).Kind() {
		case reflect.Int:
			{
				conv, ok := vv.([]interface{})
				pos := o.(int)
				if !ok {
					t.Errorf("\n##%s##\n(%d)failed to cast payload data", cs, ix)
					return
				}
				if len(conv) <= pos {
					t.Errorf("\n##%s##\n(%d)size of parameter is not as expected\nsize:%d, index:%d", cs, ix, pos, len(conv))
					return
				}
				vv = conv[pos]
			}
			break
		case reflect.String:
			{
				conv, ok := vv.(map[string]interface{})
				key := o.(string)
				if !ok {
					t.Errorf("\n##%s##\n(%d)failed to cast payload data", cs, ix)
					return
				}
				vv, ok = conv[key]
				if !ok {
					t.Errorf("\n##%s##\n(%d)member not found:%s", cs, ix, key)
					return
				}
			}
			break
		default:
			t.Errorf("\n##%s##parameter of place is not valid.", cs)
			return
		}
	}
	if vv.(string) != expect {
		t.Errorf("\n##%s##value is not as expexetd\nexpect:%s\nactual:%s", cs, expect, vv.(string))
	}
}

func CASE(val string) {
	fmt.Println("### " + val + " ###")
}
