#!/bin/bash

# Parameter Chaincode
CC_NAME="sc03-compliance"
CC_SRC_PATH="/opt/gopath/src/github.com/chaincode/sc03-compliance"
CC_VERSION="1.0"
CC_SEQUENCE="1"
CHANNEL_NAME="agritracechannel"

# Setup Lingkungan CLI
export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/agritrace.com/orderers/orderer1.agritrace.com/msp/tlscacerts/tlsca.agritrace.com-cert.pem

setGlobals() {
  local ORG=$1
  if [ "$ORG" == "Farmer" ]; then
    export CORE_PEER_LOCALMSPID="FarmerMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/farmer.agritrace.com/peers/peer0.farmer.agritrace.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/farmer.agritrace.com/users/Admin@farmer.agritrace.com/msp
    export CORE_PEER_ADDRESS=peer0.farmer.agritrace.com:7051
  elif [ "$ORG" == "Aggregator" ]; then
    export CORE_PEER_LOCALMSPID="AggregatorMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/aggregator.agritrace.com/peers/peer0.aggregator.agritrace.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/aggregator.agritrace.com/users/Admin@aggregator.agritrace.com/msp
    export CORE_PEER_ADDRESS=peer0.aggregator.agritrace.com:8051
  elif [ "$ORG" == "Processor" ]; then
    export CORE_PEER_LOCALMSPID="ProcessorMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/processor.agritrace.com/peers/peer0.processor.agritrace.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/processor.agritrace.com/users/Admin@processor.agritrace.com/msp
    export CORE_PEER_ADDRESS=peer0.processor.agritrace.com:9051
  elif [ "$ORG" == "Regulator" ]; then
    export CORE_PEER_LOCALMSPID="RegulatorMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/regulator.agritrace.com/peers/peer0.regulator.agritrace.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/regulator.agritrace.com/users/Admin@regulator.agritrace.com/msp
    export CORE_PEER_ADDRESS=peer0.regulator.agritrace.com:10051
  elif [ "$ORG" == "Buyer" ]; then
    export CORE_PEER_LOCALMSPID="BuyerMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/buyer.agritrace.com/peers/peer0.buyer.agritrace.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/buyer.agritrace.com/users/Admin@buyer.agritrace.com/msp
    export CORE_PEER_ADDRESS=peer0.buyer.agritrace.com:11051
  else
    echo "Unknown org"
  fi
}

echo "=== 1. Package Chaincode ==="
peer lifecycle chaincode package ${CC_NAME}.tar.gz --path ${CC_SRC_PATH} --lang golang --label ${CC_NAME}_${CC_VERSION}

ORGS=("Farmer" "Aggregator" "Processor" "Regulator" "Buyer")

echo "=== 2. Install Chaincode on all Peers ==="
for ORG in "${ORGS[@]}"; do
  setGlobals $ORG
  peer lifecycle chaincode install ${CC_NAME}.tar.gz
done

echo "=== 3. Query Package ID ==="
setGlobals "Farmer"
peer lifecycle chaincode queryinstalled > log.txt
PACKAGE_ID=$(sed -n "/${CC_NAME}_${CC_VERSION}/{s/^Package ID: //; s/, Label:.*$//; p;}" log.txt)
echo "PACKAGE_ID: $PACKAGE_ID"

echo "=== 4. Approve Chaincode for all Orgs ==="
for ORG in "${ORGS[@]}"; do
  setGlobals $ORG
  peer lifecycle chaincode approveformyorg -o orderer1.agritrace.com:7050 --ordererTLSHostnameOverride orderer1.agritrace.com --tls --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${CC_VERSION} --package-id ${PACKAGE_ID} --sequence ${CC_SEQUENCE}
done

echo "=== 5. Check Commit Readiness ==="
peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${CC_VERSION} --sequence ${CC_SEQUENCE} --output json

echo "=== 6. Commit Chaincode Definition ==="
setGlobals "Farmer"
peer lifecycle chaincode commit -o orderer1.agritrace.com:7050 --ordererTLSHostnameOverride orderer1.agritrace.com --tls --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${CC_VERSION} --sequence ${CC_SEQUENCE} \
  --peerAddresses peer0.farmer.agritrace.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/farmer.agritrace.com/peers/peer0.farmer.agritrace.com/tls/ca.crt \
  --peerAddresses peer0.aggregator.agritrace.com:8051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/aggregator.agritrace.com/peers/peer0.aggregator.agritrace.com/tls/ca.crt \
  --peerAddresses peer0.processor.agritrace.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/processor.agritrace.com/peers/peer0.processor.agritrace.com/tls/ca.crt \
  --peerAddresses peer0.regulator.agritrace.com:10051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/regulator.agritrace.com/peers/peer0.regulator.agritrace.com/tls/ca.crt \
  --peerAddresses peer0.buyer.agritrace.com:11051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/buyer.agritrace.com/peers/peer0.buyer.agritrace.com/tls/ca.crt

echo "=== 7. Init/Invoke Chaincode (Test Record Compliance) ==="
setGlobals "Processor"
sleep 5
peer chaincode invoke -o orderer1.agritrace.com:7050 --ordererTLSHostnameOverride orderer1.agritrace.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n ${CC_NAME} \
  --peerAddresses peer0.processor.agritrace.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/processor.agritrace.com/peers/peer0.processor.agritrace.com/tls/ca.crt \
  --peerAddresses peer0.aggregator.agritrace.com:8051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/aggregator.agritrace.com/peers/peer0.aggregator.agritrace.com/tls/ca.crt \
  -c '{"function":"RecordCompliance","Args":["RECORD_01","BATCH_01","ProcessorMSP","PROCESSING","COMPLIANT","All good"]}'

echo "=== DEPLOYMENT BERHASIL! ==="
