import os

orgs = [
    {"name": "Farmer", "domain": "farmer.agritrace.com", "port": 7051, "couch_port": 5984},
    {"name": "Aggregator", "domain": "aggregator.agritrace.com", "port": 8051, "couch_port": 6984},
    {"name": "Processor", "domain": "processor.agritrace.com", "port": 9051, "couch_port": 7984},
    {"name": "Regulator", "domain": "regulator.agritrace.com", "port": 10051, "couch_port": 8984},
    {"name": "Buyer", "domain": "buyer.agritrace.com", "port": 11051, "couch_port": 9984},
]

orderers = [
    {"name": "orderer1", "domain": "agritrace.com", "port": 7050},
    {"name": "orderer2", "domain": "agritrace.com", "port": 8050},
    {"name": "orderer3", "domain": "agritrace.com", "port": 9050},
    {"name": "orderer4", "domain": "agritrace.com", "port": 10050},
]

yaml_content = f"""version: '3.7'

networks:
  agritrace:
    name: agritrace_network

volumes:
"""

for ord in orderers:
    yaml_content += f"  {ord['name']}.{ord['domain']}:\n"
for org in orgs:
    yaml_content += f"  peer0.{org['domain']}:\n"

yaml_content += """
services:
"""

# ORDERERS
for ord in orderers:
    full_name = f"{ord['name']}.{ord['domain']}"
    yaml_content += f"""  {full_name}:
    container_name: {full_name}
    image: hyperledger/fabric-orderer:2.5
    environment:
      - FABRIC_LOGGING_SPEC=INFO
      - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
      - ORDERER_GENERAL_LISTENPORT={ord['port']}
      - ORDERER_GENERAL_LOCALMSPID=OrdererMSP
      - ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp
      - ORDERER_GENERAL_TLS_ENABLED=true
      - ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key
      - ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt
      - ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]
      - ORDERER_GENERAL_CLUSTER_CLIENTCERTIFICATE=/var/hyperledger/orderer/tls/server.crt
      - ORDERER_GENERAL_CLUSTER_CLIENTPRIVATEKEY=/var/hyperledger/orderer/tls/server.key
      - ORDERER_GENERAL_CLUSTER_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]
      - ORDERER_GENERAL_BOOTSTRAPMETHOD=file
      - ORDERER_GENERAL_BOOTSTRAPFILE=/var/hyperledger/orderer/orderer.genesis.block
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    command: orderer
    volumes:
      - ./channel-artifacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block
      - ./crypto-config/ordererOrganizations/{ord['domain']}/orderers/{full_name}/msp:/var/hyperledger/orderer/msp
      - ./crypto-config/ordererOrganizations/{ord['domain']}/orderers/{full_name}/tls/:/var/hyperledger/orderer/tls
      - {full_name}:/var/hyperledger/production/orderer
    ports:
      - "{ord['port']}:{ord['port']}"
    networks:
      - agritrace

"""

# PEERS & COUCHDB
for org in orgs:
    peer_name = f"peer0.{org['domain']}"
    couch_name = f"couchdb_{org['name'].lower()}"
    
    yaml_content += f"""  {couch_name}:
    container_name: {couch_name}
    image: couchdb:3.3.3
    environment:
      - COUCHDB_USER=admin
      - COUCHDB_PASSWORD=adminpw
    ports:
      - "{org['couch_port']}:5984"
    networks:
      - agritrace

  {peer_name}:
    container_name: {peer_name}
    image: hyperledger/fabric-peer:2.5
    environment:
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=agritrace_network
      - FABRIC_LOGGING_SPEC=INFO
      - CORE_PEER_TLS_ENABLED=true
      - CORE_PEER_PROFILE_ENABLED=false
      - CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt
      - CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key
      - CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt
      - CORE_PEER_ID={peer_name}
      - CORE_PEER_ADDRESS={peer_name}:{org['port']}
      - CORE_PEER_LISTENADDRESS=0.0.0.0:{org['port']}
      - CORE_PEER_CHAINCODEADDRESS={peer_name}:{org['port']+1}
      - CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:{org['port']+1}
      - CORE_PEER_GOSSIP_EXTERNALENDPOINT={peer_name}:{org['port']}
      - CORE_PEER_GOSSIP_BOOTSTRAP={peer_name}:{org['port']}
      - CORE_PEER_LOCALMSPID={org['name']}MSP
      - CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/msp
      - CORE_LEDGER_STATE_STATEDATABASE=CouchDB
      - CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS={couch_name}:5984
      - CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME=admin
      - CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD=adminpw
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
    command: peer node start
    volumes:
      - /var/run/docker.sock:/host/var/run/docker.sock
      - ./crypto-config/peerOrganizations/{org['domain']}/peers/{peer_name}/msp:/etc/hyperledger/fabric/msp
      - ./crypto-config/peerOrganizations/{org['domain']}/peers/{peer_name}/tls:/etc/hyperledger/fabric/tls
      - {peer_name}:/var/hyperledger/production
    ports:
      - "{org['port']}:{org['port']}"
    depends_on:
      - {couch_name}
    networks:
      - agritrace

"""

# CLI
yaml_content += """  cli:
    container_name: cli
    image: hyperledger/fabric-tools:2.5
    tty: true
    stdin_open: true
    environment:
      - GOPATH=/opt/gopath
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - FABRIC_LOGGING_SPEC=INFO
      - CORE_PEER_ID=cli
      - CORE_PEER_ADDRESS=peer0.farmer.agritrace.com:7051
      - CORE_PEER_LOCALMSPID=FarmerMSP
      - CORE_PEER_TLS_ENABLED=true
      - CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/farmer.agritrace.com/peers/peer0.farmer.agritrace.com/tls/server.crt
      - CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/farmer.agritrace.com/peers/peer0.farmer.agritrace.com/tls/server.key
      - CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/farmer.agritrace.com/peers/peer0.farmer.agritrace.com/tls/ca.crt
      - CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/farmer.agritrace.com/users/Admin@farmer.agritrace.com/msp
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
    command: /bin/bash
    volumes:
      - /var/run/docker.sock:/host/var/run/docker.sock
      - ../chaincode:/opt/gopath/src/github.com/chaincode
      - ./crypto-config:/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/
      - ./channel-artifacts:/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts
      - ./scripts:/opt/gopath/src/github.com/hyperledger/fabric/peer/scripts
    depends_on:
      - peer0.farmer.agritrace.com
      - peer0.aggregator.agritrace.com
      - peer0.processor.agritrace.com
      - peer0.regulator.agritrace.com
      - peer0.buyer.agritrace.com
    networks:
      - agritrace
"""

with open("d:/semester6/sister/agritrace-workspace/agritrace-network/docker-compose.yaml", "w") as f:
    f.write(yaml_content)

print("Successfully wrote docker-compose.yaml")
