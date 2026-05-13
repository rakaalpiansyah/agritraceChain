# AgriTraceChain - Setup Fabric + Caliper dari Nol

Panduan ini fokus ke dua hal:

1. Menjalankan project dari awal sampai benchmark jalan.
2. Menunjukkan file mana yang perlu diubah kalau ingin mengubah data/skenario.

## 1. Ringkasan Arsitektur

AgriTraceChain memakai **Hyperledger Fabric 2.5**:

- 5 organisasi peer: Farmer, Aggregator, Processor, Regulator, Buyer
- 4 orderer
- 1 channel: `agritracechannel`
- 4 chaincode: `sc01-registration`, `sc02-certification`, `sc03-compliance`, `sc04-lc-settlement`
- Benchmark: Hyperledger Caliper

## 2. Struktur Repo

```text
agritrace-workspace/
|-- agritrace-network/      # Konfigurasi Fabric, Docker Compose, script channel/deploy
|-- chaincode/              # Source chaincode Go
|-- caliper-workspace/      # Konfigurasi benchmark Caliper
`-- fabric-bin/             # Binary Fabric (cryptogen/configtxgen)
```

## 3. Prasyarat

| Komponen | Minimal                                     |
| -------- | ------------------------------------------- |
| OS       | Windows 10/11, Linux, macOS                 |
| Docker   | Docker Desktop aktif (Linux container mode) |
| Node.js  | 18+                                         |
| npm      | 9+                                          |
| Python   | 3.10+                                       |
| Git      | Terbaru                                     |

`agritrace-network\reset_crypto.py` membutuhkan binary:

- `cryptogen.exe`
- `configtxgen.exe`

Lokasi yang diharapkan:

```text
fabric-bin\fabric-samples\bin\
```

## 4. Setup dari Nol (Urutan Aman)

> Jalankan perintah dari root repo, kecuali disebutkan lain.

### 4.1 Generate ulang crypto + channel artifacts

```powershell
cd .\agritrace-network
python .\reset_crypto.py
```

Output utama:

- `agritrace-network\crypto-config\`
- `agritrace-network\channel-artifacts\genesis.block`
- `agritrace-network\channel-artifacts\channel.tx`

### 4.2 Nyalakan jaringan Fabric

```powershell
docker compose -f .\docker-compose.yaml up -d
docker ps
```

### 4.3 Buat channel + join semua peer

```powershell
docker exec cli bash -c "cd /opt/gopath/src/github.com/hyperledger/fabric/peer && bash ./scripts/join_channel.sh"
```

### 4.4 Deploy semua chaincode

```powershell
docker exec cli bash -c "cd /opt/gopath/src/github.com/hyperledger/fabric/peer && bash ./scripts/deploy_cc.sh"
docker exec cli bash -c "cd /opt/gopath/src/github.com/hyperledger/fabric/peer && bash ./scripts/deploy_sc02.sh"
docker exec cli bash -c "cd /opt/gopath/src/github.com/hyperledger/fabric/peer && bash ./scripts/deploy_sc03.sh"
docker exec cli bash -c "cd /opt/gopath/src/github.com/hyperledger/fabric/peer && bash ./scripts/deploy_sc04.sh"
```

Sampai tahap ini, jaringan siap untuk benchmark.

## 5. Menjalankan Benchmark Caliper

```powershell
cd ..\caliper-workspace
npm install
npm run profiles
npm run lint:workloads
npm run benchmark:paper
```

Report hasil benchmark:

```text
caliper-workspace\report.html
```

## 6. Kalau Mau Ubah Data, Edit di Mana?

| Tujuan Perubahan                                                             | File yang Diubah                                                                                                  |
| ---------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| Ubah logika bisnis smart contract                                            | `chaincode\sc01-registration\main.go`                                                                             |
|                                                                              | `chaincode\sc02-certification\main.go`                                                                            |
|                                                                              | `chaincode\sc03-compliance\main.go`                                                                               |
|                                                                              | `chaincode\sc04-lc-settlement\main.go`                                                                            |
| Ubah payload/data transaksi benchmark (contoh actor, batch, certificate, LC) | `caliper-workspace\workloads\registerActor.js`                                                                    |
|                                                                              | `caliper-workspace\workloads\registerBatch.js`                                                                    |
|                                                                              | `caliper-workspace\workloads\issueCertificate.js`                                                                 |
|                                                                              | `caliper-workspace\workloads\recordCompliance.js`                                                                 |
|                                                                              | `caliper-workspace\workloads\issueLC.js`                                                                          |
|                                                                              | `caliper-workspace\workloads\settleLC.js`                                                                         |
| Ubah jumlah transaksi, TPS, durasi round benchmark                           | `caliper-workspace\benchmarks\bench-config-paper.yaml`                                                            |
| Ubah skenario benchmark eksperimen/debug                                     | `caliper-workspace\benchmarks\bench-config.yaml` dan `caliper-workspace\benchmarks\debug-register-actor.yaml`     |
| Ubah profile koneksi Caliper ke Fabric                                       | `caliper-workspace\generate_profiles.py` dan folder `caliper-workspace\networks\`                                 |
| Ubah topologi node/port/container Fabric                                     | `agritrace-network\docker-compose.yaml`                                                                           |
| Ubah organisasi, MSP, channel profile                                        | `agritrace-network\crypto-config.yaml` dan `agritrace-network\configtx.yaml`                                      |
| Ubah langkah join/deploy                                                     | `agritrace-network\scripts\join_channel.sh`, `deploy_cc.sh`, `deploy_sc02.sh`, `deploy_sc03.sh`, `deploy_sc04.sh` |

### Setelah ubah data/logic, jalankan ulang apa?

1. Jika ubah **workload Caliper** atau **bench config**: cukup jalankan ulang benchmark (`npm run benchmark:paper`).
2. Jika ubah **chaincode**: jalankan ulang script deploy chaincode terkait.
3. Jika ubah **config inti Fabric** (`crypto-config.yaml`, `configtx.yaml`, atau materi crypto): lakukan reset dari awal (`reset_crypto.py`, lalu `docker compose up`, join channel, deploy ulang).

## 7. Command Referensi Cepat

| Command                        | Fungsi                             |
| ------------------------------ | ---------------------------------- |
| `npm run profiles`             | Generate connection profile Fabric |
| `npm run lint:workloads`       | Cek sintaks workload JavaScript    |
| `npm run test`                 | Flow-only test Caliper             |
| `npm run debug:register-actor` | Debug 1 skenario register actor    |
| `npm run benchmark`            | Benchmark eksperimen               |
| `npm run benchmark:paper`      | Benchmark final untuk pelaporan    |

## 8. Troubleshooting Singkat

| Error                                                    | Penyebab Umum                                       | Solusi                                                         |
| -------------------------------------------------------- | --------------------------------------------------- | -------------------------------------------------------------- |
| `no peer combination can satisfy the endorsement policy` | Target endorsement tidak cocok / network belum siap | Pastikan semua peer join channel dan chaincode sudah committed |
| `spawn EPERM` di Windows                                 | Permission/process lock                             | Jalankan dari PowerShell normal dan cek antivirus/file lock    |
| Query round gagal                                        | Konfigurasi query handler belum final               | Gunakan `npm run benchmark:paper` untuk hasil final            |

## 9. Shutdown dan Reset

Matikan jaringan:

```powershell
cd ..\agritrace-network
docker compose -f .\docker-compose.yaml down
```

Matikan + hapus volume:

```powershell
docker compose -f .\docker-compose.yaml down -v
```

Reset penuh artifact Fabric:

```powershell
python .\reset_crypto.py
```
