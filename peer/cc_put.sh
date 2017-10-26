#!/bin/bash

peer chaincode invoke -n soila -v 0 -c '{"Args":["person.put", "test", "aaaaaaa"]}' -C myc
