package person

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/kenmazsyma/soila/chaincode/cmn"
)

type Person struct{}

type PersonRec struct {
	Ver        []string
	Activity   []string
	Reputation []string
}

var KEY_TYPE = "PERSON"

/***************************************************
[person_genearteKey]
description : generate key for Person
parameters  :
   stub - chaincode interface
   id - person id
return: response object
***************************************************/
func person_generateKey(stub shim.ChaincodeStubInterface, id string) (string, error) {
	//	sv, err := stub.GetCreator()
	//	if err != nil {
	//		return "", err
	//	}
	//	fmt.Printf("EE:%x(%s)\n", sv, string(sv))
	// TODO:set creator id
	sv := "TEST"
	return stub.CreateCompositeKey(KEY_TYPE, []string{string(sv), id})
}

/***************************************************
[Put]
description : add Person object to blockchain
parameters  :
   stub - chaincode interface
   args - [personid, data]
return:
   1:response data
   2:error object if error occured
***************************************************/
func Put(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("Invalid Arguments")
	}
	// get ID from blockchain
	key, err := person_generateKey(stub, args[0])
	if err != nil {
		return "", err
	}
	fmt.Printf("KEY:%s\n", key)
	// check if data is already exists.
	val, err := stub.GetState(key)
	if err != nil {
		return "", err
	}
	if val != nil {
		return "", errors.New("data is already exists.")
	}
	// put data into blockchain
	data := PersonRec{
		Ver:        []string{args[1]},
		Activity:   []string{},
		Reputation: []string{},
	}
	err = cmn.Put(stub, key, data)
	return "", nil
}

/***************************************************
[Update]
description : update Person object to blockchain
parameters  :
   stub - chaincode interface
   args - [personid, data]
return:
   1:response data
   2:error object if error occured
***************************************************/
func Update(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("Invalid Arguments")
	}
	// get ID from blockchain
	key, err := person_generateKey(stub, args[0])
	if err != nil {
		return "", err
	}
	// check if data is already exists.
	fmt.Printf("KEY:%s\n", key)
	val, err := stub.GetState(key)
	if err != nil {
		return "", err
	}
	if val == nil {
		return "", errors.New("data is not exists.")
	}
	data := PersonRec{}
	err = json.Unmarshal(val, &data)
	// put data into blockchain
	data.Ver = append(data.Ver, args[1])
	err = cmn.Put(stub, key, data)
	return "", nil
}

/***************************************************
[Get]
description : get value correspond to the key
parameters  :
   stub - chaincode interface
   args - key-value pair
return: error string if error occured
***************************************************/
func Get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("Invalid parameter")
	}
	key, err := person_generateKey(stub, args[0])
	val, err := stub.GetState(key)
	fmt.Printf("GetState:%s\n", val)
	return string(val), err
}
