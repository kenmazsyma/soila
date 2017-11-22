package person

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/kenmazsyma/soila/chaincode/cmn"
)

type Person struct{}

type PersonReputation struct {
	Setter  string
	Content string
	Type    string
}

type PersonRec struct {
	Ver        []string
	Activity   []string
	Reputation []PersonReputation
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
	o, err := stub.GetCreator()
	if err != nil {
		fmt.Println("Error at GetCreator")
		return "", err
	}
	js, err := cmn.ToJSON(o)
	if err != nil {
		fmt.Println("Error at ToJSON")
		return "", err
	} else {
		fmt.Println("Creator:" + js)
	}
	return stub.CreateCompositeKey(KEY_TYPE, []string{string(sv), id})
}

/***************************************************
[get_and_check]
description : generate key for Person
parameters  :
   stub - chaincode interface
   args - parameters
   validlen - valid length of args
return:
	rec - Person correspond to person id
	key - key of Person
	err - err object
***************************************************/
func get_and_check(stub shim.ChaincodeStubInterface, args []string, validlen int) (rec *PersonRec, key string, err error) {
	if len(args) != validlen {
		err = errors.New("Invalid Arguments")
		return
	}
	// get ID from blockchain
	key, err = person_generateKey(stub, args[0])
	if err != nil {
		return
	}
	// check if data is already exists.
	fmt.Printf("KEY:%s\n", key)
	val, err := stub.GetState(key)
	if err != nil {
		return
	}
	if val == nil {
		err = errors.New("data is not exists.")
		return
	}
	data := PersonRec{}
	err = json.Unmarshal(val, &data)
	rec = &data
	return
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
		Ver:        []string{cmn.Sha1(args[1])},
		Activity:   []string{},
		Reputation: []PersonReputation{},
	}
	err = cmn.Put(stub, key, data)
	return "", nil
}

/***************************************************
[Update]
description : update Person object
parameters  :
   stub - chaincode interface
   args - [personid, data]
return:
   1:response data
   2:error object if error occured
***************************************************/
func Update(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	data, key, err := get_and_check(stub, args, 2)
	if err != nil {
		return "", err
	}
	// put data into blockchain
	(*data).Ver = append((*data).Ver, cmn.Sha1(args[1]))
	err = cmn.Put(stub, key, (*data))
	return "", nil
}

/***************************************************
[Get]
description : get value correspond to the key
parameters  :
   stub - chaincode interface
   args - personid
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

/***************************************************
[AddActivity]
description : add activity hash of Person to Person object
parameters  :
   stub - chaincode interface
   args - [personid, content hash]
return:
   1:response data
   2:error object if error occured
***************************************************/
func AddActivity(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	data, key, err := get_and_check(stub, args, 2)
	if err != nil {
		return "", err
	}
	// TODO: check if sender is valid
	// put data into blockchain
	(*data).Activity = append((*data).Activity, args[1])
	err = cmn.Put(stub, key, (*data))
	return "", nil
}

/***************************************************
[AddReputation]
description : add activity hash of Person to Person object
parameters  :
   stub - chaincode interface
   args - [personid, setter, content hash, type]
return:
   1:response data
   2:error object if error occured
***************************************************/
func AddReputation(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	data, key, err := get_and_check(stub, args, 4)
	if err != nil {
		return "", err
	}
	// put data into blockchain
	rep := PersonReputation{
		Setter:  args[1],
		Content: args[2],
		Type:    args[3],
	}
	// TODO: check if content hash is valid
	(*data).Reputation = append((*data).Reputation, rep)
	err = cmn.Put(stub, key, (*data))
	return "", nil
}

/***************************************************
[RemoveReputation]
description : add activity hash of Person to Person object
parameters  :
   stub - chaincode interface
   args - [personid, setter, content hash]
return:
   1:response data
   2:error object if error occured
***************************************************/
func RemoveReputation(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	data, key, err := get_and_check(stub, args, 3)
	if err != nil {
		return "", err
	}
	for i, v := range (*data).Reputation {
		if v.Setter == args[1] && v.Content == args[2] {
			(*data).Reputation = append((*data).Reputation[0:i], (*data).Reputation[i+1:]...)
			err = cmn.Put(stub, key, (*data))
			return "", err
		}
	}
	return "", errors.New("target data is not found")
}
