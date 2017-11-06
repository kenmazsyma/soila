#!/bin/bash

peer chaincode query -n soila -v 0 -c '{"Args":["person.get", "test"]}' -C soila
