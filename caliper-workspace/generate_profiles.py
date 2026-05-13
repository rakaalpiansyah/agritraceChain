import os

orgs = [
    {"name": "farmer", "msp": "FarmerMSP", "peer_port": 7051},
    {"name": "aggregator", "msp": "AggregatorMSP", "peer_port": 8051},
    {"name": "processor", "msp": "ProcessorMSP", "peer_port": 9051},
    {"name": "regulator", "msp": "RegulatorMSP", "peer_port": 10051},
    {"name": "buyer", "msp": "BuyerMSP", "peer_port": 11051},
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

out_dir = "d:/semester6/sister/agritrace-workspace/caliper-workspace/networks"

for target_org in orgs:
    content = f"""name: AgriTrace {target_org['name'].capitalize()} Connection Profile
version: "1.0"

client:
  organization: {target_org['name'].capitalize()}Org

organizations:
{orgs_block}
peers:
{peers_block}"""

    filepath = os.path.join(out_dir, f"connectionProfile-{target_org['name']}.yaml")
    with open(filepath, "w", newline="\\n") as f:
        f.write(content)
    print(f"  Created {filepath}")

print("All connection profiles generated with all peers!")
