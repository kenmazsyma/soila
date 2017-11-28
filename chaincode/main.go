/*
Package main provides chaincode for soila_chain.
*/
package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/kenmazsyma/soila/chaincode/log"
	"github.com/kenmazsyma/soila/chaincode/peer"
	"github.com/kenmazsyma/soila/chaincode/person"
	"github.com/kenmazsyma/soila/chaincode/project"
	"github.com/kenmazsyma/soila/chaincode/root"
	"github.com/kenmazsyma/soila/chaincode/token"
)

var invoke_list = map[string]root.InvokeRoutineType{
	"person.register":          person.Register,
	"person.update":            person.Update,
	"person.get":               person.Get,
	"person.add_activity":      person.AddActivity,
	"person.add_reputation":    person.AddReputation,
	"person.remove_reputation": person.RemoveReputation,
	"peer.register":            peer.Register,
	"peer.update":              peer.Update,
	"peer.get":                 peer.Get,
	"peer.deregister":          peer.Deregister,
	"project.register":         project.Register,
	"project.get":              project.Get,
	"project.updatestatus":     project.UpdateStatus,
	"token.register":           token.Register,
	"token.update":             token.Update,
	"token.remove":             token.Remove,
}

// ================================================
// main
// ================================================

// main is a function for executing chaincode for soila_chain
func main() {
	log.Init("soila_chain", log.LEVEL_DEBUG)
	cc := new(root.CC)
	cc.SetInvokeMap(invoke_list)
	err := shim.Start(cc)
	if err != nil {
		log.Error("Error starting Simple chaincode: " + err.Error())
	}
}
