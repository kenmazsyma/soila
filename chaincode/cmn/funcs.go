package cmn

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func Put(stub shim.ChaincodeStubInterface, key string, val interface{}) error {
	if val == nil {
		return errors.New("invalid param")
	}
	jsVal, err := json.Marshal(val)
	if err != nil {
		return err
	}
	fmt.Printf("KEY:%s\n", key)
	err = stub.PutState(key, []byte(jsVal))
	return err
}

func Get(stub shim.ChaincodeStubInterface, key string) (val interface{}, err error) {
	var jsVal []byte
	jsVal, err = stub.GetState(key)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsVal, &val)
	return
}

//package main
//
//import (
//)
//
////updating tradeBlock once it is passed in, creating new TradeBlock if one does not exist
////adding on the contract value if it does exist
////input is ledger map and returning an updated tradeblock or ledger
//func (a *ArgsMap) updateTradeBlock(regCompany bool, tradeCredits string, tradePrice string, tradetimestamp string, tradeCompany string, tradeType string) (map[string]interface {}){
//    var tradeBlockMap map[string]interface{}
//    //get the object from ledger if tradeHistory already exists
//    tbytes, found := getObject(*a, "tradeHistory")
//    //if found is false, then a tradeHistory does not exists and new struct needs to be created
//    if found == false {
//        tradeBlockMap = make(map[string]interface{})
//        tradeBlockMap["credits"] = []string{tradeCredits}
//        tradeBlockMap["price"] = []string{tradePrice}
//        tradeBlockMap["timestamp"] = []string{tradetimestamp}
//        if regCompany {
//            tradeBlockMap["company"] = []string{tradeCompany}
//            tradeBlockMap["buysell"] = []string{tradeType}
//        }
//    } else {
//        tradeBlockMap = tbytes.(map[string]interface{})
//        //appending all the new attributes
//        tradeBlockMap["credits"] = append(tradeBlockMap["credits"].([]interface{}), tradeCredits)
//        tradeBlockMap["price"] = append(tradeBlockMap["price"].([]interface{}), tradePrice)
//        tradeBlockMap["timestamp"] = append(tradeBlockMap["timestamp"].([]interface{}), tradetimestamp)
//        if regCompany {
//            tradeBlockMap["company"] = append(tradeBlockMap["price"].([]interface{}), tradeCompany)
//            tradeBlockMap["buysell"] = append(tradeBlockMap["buysell"].([]interface{}), tradeType)
//        }
//    }
//    return tradeBlockMap
//}
