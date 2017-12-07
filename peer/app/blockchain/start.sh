#!/bin/bash

export NODE_CONFIG_DIR=../config
node joinchannel.js org1
node joinchannel.js org2
node install_cc.js org1
node install_cc.js org2
node instantiate_cc.js org1
