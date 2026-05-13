# AgriTraceChain - Hyperledger Fabric Network & Caliper Benchmark

Dokumen ini adalah panduan resmi untuk menyiapkan environment pengembangan AgriTraceChain secara rapi, mulai dari persiapan Docker, bootstrap jaringan Hyperledger Fabric, deployment chaincode, hingga eksekusi benchmark dengan Hyperledger Caliper.

## 1. Ringkasan Proyek

AgriTraceChain menjalankan jaringan **Hyperledger Fabric 2.5** dengan:

- **5 organisasi peer**: Farmer, Aggregator, Processor, Regulator, Buyer
- **4 orderer node**
- **1 channel utama**: `agritracechannel`
- **4 smart contract**:
  - `sc01-registration`
  - `sc02-certification`
  - `sc03-compliance`
  - `sc04-lc-settlement`
- **Benchmark engine**: Hyperledger Caliper

## 2. Struktur Repository

```text
agritrace-workspace/
|-- agritrace-network/      # Konfigurasi Fabric, Docker Compose, script channel/deploy
|-- chaincode/              # Source chaincode Go
|-- caliper-workspace/      # Konfigurasi dan workload benchmark Caliper
`-- fabric-bin/             # Binary Fabric (cryptogen/configtxgen) - siapkan manual
```

## 3. Prasyarat

Pastikan komponen berikut tersedia:

| Komponen | Kebutuhan |
|---|---|
| OS | Windows 10/11, Linux, atau macOS |
| Docker | Docker Desktop aktif, mode Linux container |
| Node.js | v18+ |
| npm | v9+ |
| Python | v3.10+ |
| Git | Versi terbaru |

### Prasyarat tambahan Fabric

Script `agritrace-network\reset_crypto.py` membutuhkan binary:

- `cryptogen`
- `configtxgen`

Lokasi yang dipakai script:

```text
fabric-bin\fabric-samples\bin\
```

Jika folder tersebut belum berisi binary Fabric, siapkan terlebih dahulu sebelum menjalankan reset.

## 4. Alur Setup End-to-End (Dari Nol)

> Jalankan perintah dari root repo, kecuali jika disebutkan folder spesifik.

### 4.1 Generate ulang crypto material dan channel artifacts

```powershell
cd .\agritrace-network
python .\reset_crypto.py
```

Output utama:

- `agritrace-network\crypto-config\`
- `agritrace-network\channel-artifacts\genesis.block`
- `agritrace-network\channel-artifacts\channel.tx`

### 4.2 Jalankan seluruh container Fabric

```powershell
docker compose -f .\docker-compose.yaml up -d
```

Cek status:

```powershell
docker ps
```

### 4.3 Buat channel dan join seluruh peer

```powershell
docker exec cli bash -c "cd /opt/gopath/src/github.com/hyperledger/fabric/peer && bash ./scripts/join_channel.sh"
```

### 4.4 Deploy chaincode

Jalankan dari host (masing-masing satu kali):

```powershell
docker exec cli bash -c "cd /opt/gopath/src/github.com/hyperledger/fabric/peer && bash ./scripts/deploy_cc.sh"
docker exec cli bash -c "cd /opt/gopath/src/github.com/hyperledger/fabric/peer && bash ./scripts/deploy_sc02.sh"
docker exec cli bash -c "cd /opt/gopath/src/github.com/hyperledger/fabric/peer && bash ./scripts/deploy_sc03.sh"
docker exec cli bash -c "cd /opt/gopath/src/github.com/hyperledger/fabric/peer && bash ./scripts/deploy_sc04.sh"
```

Setelah tahap ini, seluruh chaincode sudah terpasang dan siap dipakai benchmark.

## 5. Setup dan Eksekusi Benchmark (Caliper)

Masuk ke workspace benchmark:

```powershell
cd ..\caliper-workspace
npm install
```

Generate connection profile:

```powershell
npm run profiles
```

Validasi sintaks workload:

```powershell
npm run lint:workloads
```

Jalankan benchmark final (direkomendasikan):

```powershell
npm run benchmark:paper
```

Hasil report:

```text
caliper-workspace\report.html
```

## 6. Penjelasan Command Caliper

| Perintah | Fungsi |
|---|---|
| `npm run profiles` | Generate ulang connection profile untuk seluruh organisasi |
| `npm run lint:workloads` | Validasi sintaks file workload JavaScript |
| `npm run test` | Menjalankan flow-only test Caliper |
| `npm run debug:register-actor` | Debug 1 skenario transaksi `RegisterActor` |
| `npm run benchmark` | Menjalankan benchmark eksperimen lengkap |
| `npm run benchmark:paper` | Menjalankan benchmark final untuk pelaporan |

## 7. Catatan Teknis Penting

1. Benchmark ini menggunakan `fabric-network@2.x` (bukan Fabric Gateway SDK baru).
2. Jangan menambahkan flag `--caliper-fabric-gateway-enabled` untuk flow benchmark ini.
3. Endorsement policy menggunakan mayoritas organisasi, sehingga transaksi write perlu target endorsement multi-peer.

## 8. Troubleshooting Cepat

| Masalah | Penyebab Umum | Tindakan |
|---|---|---|
| `no peer combination can satisfy the endorsement policy` | Peer target endorsement tidak sesuai policy / jaringan belum siap | Pastikan semua peer join channel, chaincode committed, dan script deploy selesai |
| `spawn EPERM` (Windows) | Konflik permission/process lock | Jalankan dari PowerShell/Git Bash normal dan cek antivirus/file lock |
| Round query gagal di benchmark eksperimen | Konfigurasi query handler belum final | Gunakan `npm run benchmark:paper` untuk hasil final |

## 9. Shutdown dan Reset

Matikan jaringan:

```powershell
cd ..\agritrace-network
docker compose -f .\docker-compose.yaml down
```

Matikan jaringan + hapus volume:

```powershell
docker compose -f .\docker-compose.yaml down -v
```

Untuk reset penuh artifact Fabric, jalankan lagi:

```powershell
python .\reset_crypto.py
```
