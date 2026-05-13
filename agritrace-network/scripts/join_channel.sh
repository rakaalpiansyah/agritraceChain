#!/bin/bash

# Pastikan script berhenti jika ada error
set -e

echo "=== 1. Membuat Channel ==="
# Membuat channel (secara default CLI menggunakan environment FarmerMSP)
peer channel create -o orderer1.agritrace.com:7050 -c agritracechannel -f ./channel-artifacts/channel.tx --outputBlock ./channel-artifacts/agritracechannel.block --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/agritrace.com/orderers/orderer1.agritrace.com/msp/tlscacerts/tlsca.agritrace.com-cert.pem

echo "=== 2. Join Peer Farmer ==="
peer channel join -b ./channel-artifacts/agritracechannel.block

echo "=== 3. Join Peer Aggregator ==="
CORE_PEER_LOCALMSPID="AggregatorMSP" \
CORE_PEER_ADDRESS="peer0.aggregator.agritrace.com:8051" \
CORE_PEER_TLS_ROOTCERT_FILE="/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/aggregator.agritrace.com/peers/peer0.aggregator.agritrace.com/tls/ca.crt" \
CORE_PEER_MSPCONFIGPATH="/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/aggregator.agritrace.com/users/Admin@aggregator.agritrace.com/msp" \
peer channel join -b ./channel-artifacts/agritracechannel.block

echo "=== 4. Join Peer Processor ==="
CORE_PEER_LOCALMSPID="ProcessorMSP" \
CORE_PEER_ADDRESS="peer0.processor.agritrace.com:9051" \
CORE_PEER_TLS_ROOTCERT_FILE="/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/processor.agritrace.com/peers/peer0.processor.agritrace.com/tls/ca.crt" \
CORE_PEER_MSPCONFIGPATH="/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/processor.agritrace.com/users/Admin@processor.agritrace.com/msp" \
peer channel join -b ./channel-artifacts/agritracechannel.block

echo "=== 5. Join Peer Regulator ==="
CORE_PEER_LOCALMSPID="RegulatorMSP" \
CORE_PEER_ADDRESS="peer0.regulator.agritrace.com:10051" \
CORE_PEER_TLS_ROOTCERT_FILE="/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/regulator.agritrace.com/peers/peer0.regulator.agritrace.com/tls/ca.crt" \
CORE_PEER_MSPCONFIGPATH="/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/regulator.agritrace.com/users/Admin@regulator.agritrace.com/msp" \
peer channel join -b ./channel-artifacts/agritracechannel.block

echo "=== 6. Join Peer Buyer ==="
CORE_PEER_LOCALMSPID="BuyerMSP" \
CORE_PEER_ADDRESS="peer0.buyer.agritrace.com:11051" \
CORE_PEER_TLS_ROOTCERT_FILE="/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/buyer.agritrace.com/peers/peer0.buyer.agritrace.com/tls/ca.crt" \
CORE_PEER_MSPCONFIGPATH="/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/buyer.agritrace.com/users/Admin@buyer.agritrace.com/msp" \
peer channel join -b ./channel-artifacts/agritracechannel.block

echo "=== Channel Join Berhasil! ==="
