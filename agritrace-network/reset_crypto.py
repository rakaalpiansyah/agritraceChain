"""
Full reset: regenerate crypto, fix for Linux containers, regenerate genesis/channel.
Run this from agritrace-network/ directory.
"""
import os
import glob
import shutil
import subprocess
import sys

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
FABRIC_BIN = os.path.join(SCRIPT_DIR, "..", "fabric-bin", "fabric-samples", "bin")
CRYPTOGEN = os.path.join(FABRIC_BIN, "cryptogen.exe")
CONFIGTXGEN = os.path.join(FABRIC_BIN, "configtxgen.exe")

os.chdir(SCRIPT_DIR)

# 1. Remove old crypto
print("=== Step 1: Removing old crypto-config ===")
if os.path.exists("crypto-config"):
    shutil.rmtree("crypto-config")

# 2. Generate fresh crypto
print("=== Step 2: Generating fresh crypto material ===")
result = subprocess.run([CRYPTOGEN, "generate", "--config=crypto-config.yaml"], capture_output=True, text=True)
print(result.stdout)
if result.returncode != 0:
    print(f"ERROR: {result.stderr}")
    sys.exit(1)

# 3. Fix Windows line endings and backslashes in all config.yaml and .pem files
print("=== Step 3: Fixing config.yaml and .pem files for Linux containers ===")
count = 0
for filepath in glob.iglob("crypto-config/**/config.yaml", recursive=True):
    with open(filepath, "rb") as f:
        content = f.read()
    new_content = content.replace(b"\\", b"/").replace(b"\r\n", b"\n")
    with open(filepath, "wb") as f:
        f.write(new_content)
    count += 1

pem_count = 0
for filepath in glob.iglob("crypto-config/**/*.pem", recursive=True):
    with open(filepath, "rb") as f:
        content = f.read()
    new_content = content.replace(b"\r\n", b"\n")
    with open(filepath, "wb") as f:
        f.write(new_content)
    pem_count += 1
print(f"  Fixed {count} config.yaml files and {pem_count} .pem files")

# 4. Remove old channel artifacts
print("=== Step 4: Removing old channel artifacts ===")
for fname in ["genesis.block", "channel.tx", "agritracechannel.block"]:
    fpath = os.path.join("channel-artifacts", fname)
    if os.path.exists(fpath):
        os.remove(fpath)

# 5. Generate genesis block and channel tx
print("=== Step 5: Generating genesis block ===")
env = os.environ.copy()
env["FABRIC_CFG_PATH"] = SCRIPT_DIR

result = subprocess.run(
    [CONFIGTXGEN, "-profile", "AgriTraceGenesis", "-channelID", "system-channel",
     "-outputBlock", "./channel-artifacts/genesis.block"],
    capture_output=True, text=True, env=env
)
print(result.stdout)
if result.returncode != 0:
    print(f"ERROR: {result.stderr}")
    sys.exit(1)

print("=== Step 6: Generating channel transaction ===")
result = subprocess.run(
    [CONFIGTXGEN, "-profile", "AgriTraceChannel",
     "-outputCreateChannelTx", "./channel-artifacts/channel.tx",
     "-channelID", "agritracechannel"],
    capture_output=True, text=True, env=env
)
print(result.stdout)
if result.returncode != 0:
    print(f"ERROR: {result.stderr}")
    sys.exit(1)

print("\n=== DONE! All artifacts generated successfully. ===")
print("Next: run 'docker-compose up -d' to start the network.")
