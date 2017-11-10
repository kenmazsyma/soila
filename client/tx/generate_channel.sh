#!/bin/bash

configtxgen -profile SoilaOrdererGenesis -channelID soila -outputBlock soila.block 
configtxgen -profile SoilaChannel -channelID soila -outputCreateChannelTx soila.tx
configtxgen -channelID soila -outputBlock soila.block -inspectBlock soila.block -profile SoilaOrdererGenesis
configtxgen -channelID soila -outputCreateChannelTx soila.tx -inspectChannelCreateTx soila.tx -profile SoilaChannel
mkdir org1
mkdir org2
cp soila.block org1/soila.block
cp soila.block org2/soila.block
cp soila.tx org1/soila.tx
cp soila.tx org2/soila.tx
