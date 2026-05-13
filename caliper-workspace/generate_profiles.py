from pathlib import Path

orgs = [
    {"name": "farmer", "msp": "FarmerMSP", "peer_port": 7051},
    {"name": "aggregator", "msp": "AggregatorMSP", "peer_port": 8051},
    {"name": "processor", "msp": "ProcessorMSP", "peer_port": 9051},
    {"name": "regulator", "msp": "RegulatorMSP", "peer_port": 10051},
    {"name": "buyer", "msp": "BuyerMSP", "peer_port": 11051},
]

orderers = [
    {"name": "orderer1.agritrace.com", "port": 7050},
    {"name": "orderer2.agritrace.com", "port": 8050},
    {"name": "orderer3.agritrace.com", "port": 9050},
    {"name": "orderer4.agritrace.com", "port": 10050},
]

# Build the peers block (all peers listed)
peers_block = ""
for org in orgs:
    peers_block += f"""  peer0.{org['name']}.agritrace.com:
    url: grpcs://localhost:{org['peer_port']}
    tlsCACerts:
      path: ../agritrace-network/crypto-config/peerOrganizations/{org['name']}.agritrace.com/tlsca/tlsca.{org['name']}.agritrace.com-cert.pem
    grpcOptions:
      ssl-target-name-override: peer0.{org['name']}.agritrace.com
      hostnameOverride: peer0.{org['name']}.agritrace.com
"""

# Build the organizations block (all orgs)
orgs_block = ""
for org in orgs:
    orgs_block += f"""  {org['name'].capitalize()}Org:
    mspid: {org['msp']}
    peers:
      - peer0.{org['name']}.agritrace.com
"""

channel_peers_block = ""
for org in orgs:
    channel_peers_block += f"""    peer0.{org['name']}.agritrace.com:
      endorsingPeer: true
      chaincodeQuery: true
      ledgerQuery: true
      eventSource: true
"""

channel_orderers_block = ""
for orderer in orderers:
    channel_orderers_block += f"    - {orderer['name']}\n"

orderers_block = ""
for orderer in orderers:
    orderers_block += f"""  {orderer['name']}:
    url: grpcs://localhost:{orderer['port']}
    tlsCACerts:
      path: ../agritrace-network/crypto-config/ordererOrganizations/agritrace.com/tlsca/tlsca.agritrace.com-cert.pem
    grpcOptions:
      ssl-target-name-override: {orderer['name']}
      hostnameOverride: {orderer['name']}
"""

out_dir = Path(__file__).resolve().parent / "networks"
out_dir.mkdir(parents=True, exist_ok=True)

for target_org in orgs:
    content = f"""name: AgriTrace {target_org['name'].capitalize()} Connection Profile
version: "1.0"

client:
  organization: {target_org['name'].capitalize()}Org

organizations:
{orgs_block}
channels:
  agritracechannel:
    orderers:
{channel_orderers_block}
    peers:
{channel_peers_block}
orderers:
{orderers_block}
peers:
{peers_block}"""

    filepath = out_dir / f"connectionProfile-{target_org['name']}.yaml"
    with open(filepath, "w", newline="\n") as f:
        f.write(content)
    print(f"  Created {filepath}")

print("All connection profiles generated with all peers!")
