# This script copies source code for testing module to fabric project.
# We need to use GetCreator function for testing chaincode, 
# but this function of testing module in fabric project is not implemented.
# We can't implement original testing module
# because ChaincodeStubInterface depends on vendoring module of fabric project.
# We had no choice but to replace testing module of fabric project to customized one.

#!/bin/bash

pushd $GOPATH/src/github.com/hyperledger/fabric/core/chaincode/shim
cp mockstub.go mockstub.go_
popd
cp mockstub.go $GOPATH/src/github.com/hyperledger/fabric/core/chaincode/shim/
