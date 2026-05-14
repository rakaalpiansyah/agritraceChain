# AgriTrace-ID Project Guide

Dokumen ini menjelaskan secara rinci cara membuat, menjalankan, memahami struktur file, melakukan deploy smart contract, dan menjalankan benchmark performa proyek AgriTrace-ID dari awal.

AgriTrace-ID adalah sistem traceability komoditas pertanian berbasis Hyperledger Fabric. Sistem ini memodelkan alur supply chain pertanian dari registrasi aktor dan batch, sertifikasi, compliance, sampai penyelesaian pembayaran berbasis Letter of Credit.

## 1. Gambaran Umum

AgriTrace-ID terdiri dari tiga bagian utama:

- `agritrace-network`: konfigurasi jaringan Hyperledger Fabric, Docker Compose, crypto material, channel artifact, dan script deploy.
- `chaincode`: source code smart contract dalam bahasa Go.
- `caliper-workspace`: konfigurasi Hyperledger Caliper untuk benchmark throughput, latency, dan success rate.

Komponen live network:

| Komponen | Jumlah | Keterangan |
|---|---:|---|
| Orderer | 4 | Menggunakan Raft/etcdraft |
| Peer | 5 | Satu peer aktif per organisasi pada Docker Compose |
| CouchDB | 5 | State database untuk setiap peer |
| Channel | 1 | `agritracechannel` |
| Smart Contract | 4 | Registration, Certification, Compliance, LC Settlement |
| Benchmark Tool | 1 | Hyperledger Caliper |

Organisasi dalam jaringan:

| Organisasi | MSP ID | Domain |
|---|---|---|
| Farmer | `FarmerMSP` | `farmer.agritrace.com` |
| Aggregator | `AggregatorMSP` | `aggregator.agritrace.com` |
| Processor | `ProcessorMSP` | `processor.agritrace.com` |
| Regulator | `RegulatorMSP` | `regulator.agritrace.com` |
| Buyer | `BuyerMSP` | `buyer.agritrace.com` |

## 2. Alur Bisnis Sistem

Secara konseptual, alur AgriTrace-ID adalah:

1. Aktor supply chain didaftarkan ke jaringan.
2. Petani atau entitas terkait mendaftarkan batch komoditas.
3. Regulator menerbitkan sertifikat kualitas untuk batch.
4. Aggregator atau Processor mencatat data compliance/traceability.
5. Buyer menerbitkan Letter of Credit.
6. Buyer atau Regulator menyelesaikan Letter of Credit.

Alur ini dibagi menjadi empat smart contract agar modular dan mudah diuji.

## 3. Prasyarat

Sebelum menjalankan project, pastikan perangkat sudah memiliki:

| Komponen | Rekomendasi |
|---|---|
| OS | Windows 10/11 dengan PowerShell atau Git Bash |
| Docker | Docker Desktop dengan Linux container mode |
| Git | Versi terbaru |
| Python | 3.10 atau lebih baru |
| Node.js | 18 atau lebih baru |
| npm | 9 atau lebih baru |
| Go | 1.20 atau kompatibel dengan Fabric 2.5 |
| Hyperledger Fabric binaries | `cryptogen` dan `configtxgen` |

Binary Fabric yang dibutuhkan oleh script:

```text
fabric-bin/fabric-samples/bin/cryptogen.exe
fabric-bin/fabric-samples/bin/configtxgen.exe
```

Jika menggunakan Linux/macOS, sesuaikan ekstensi executable dan path binary Fabric.

## 4. Struktur Repository

Struktur utama repository:

```text
agritrace-workspace/
|-- .github/
|-- agritrace-network/
|-- caliper-workspace/
|-- chaincode/
|-- fabric-bin/
|-- .gitignore
|-- README.md
`-- PROJECT_GUIDE.md
```

## 5. Fungsi Folder dan File Utama

### 5.1 Root Project

| Path | Fungsi |
|---|---|
| `README.md` | Ringkasan project dan quick start |
| `PROJECT_GUIDE.md` | Dokumentasi lengkap project dari awal |
| `.gitignore` | Mengabaikan dependency, generated artifacts, report, dan file lokal |
| `.github/` | Konfigurasi GitHub jika digunakan untuk CI/CD atau metadata repository |
| `fabric-bin/` | Tempat binary dan sample Fabric lokal |

Catatan: folder `fabric-bin/` biasanya tidak perlu di-commit karena berisi binary dan sample bawaan Fabric.

### 5.2 Folder `agritrace-network`

Folder ini berisi konfigurasi jaringan Hyperledger Fabric.

| Path | Fungsi |
|---|---|
| `agritrace-network/docker-compose.yaml` | Mendefinisikan container orderer, peer, CouchDB, CLI, volume, port, dan network Docker |
| `agritrace-network/crypto-config.yaml` | Definisi organisasi untuk `cryptogen`, termasuk orderer org dan peer org |
| `agritrace-network/configtx.yaml` | Definisi channel, consortium, orderer profile, application policy, dan endorsement policy |
| `agritrace-network/reset_crypto.py` | Script otomatis untuk menghapus crypto lama dan generate ulang crypto serta channel artifact |
| `agritrace-network/scripts/join_channel.sh` | Script membuat channel dan join semua peer ke `agritracechannel` |
| `agritrace-network/scripts/deploy_cc.sh` | Deploy smart contract `sc01-registration` |
| `agritrace-network/scripts/deploy_sc02.sh` | Deploy smart contract `sc02-certification` |
| `agritrace-network/scripts/deploy_sc03.sh` | Deploy smart contract `sc03-compliance` |
| `agritrace-network/scripts/deploy_sc04.sh` | Deploy smart contract `sc04-lc-settlement` |
| `agritrace-network/crypto-config/` | Generated crypto material dari Fabric, tidak perlu dibuat manual |
| `agritrace-network/channel-artifacts/` | Generated genesis block, channel transaction, dan channel block |

### 5.3 Folder `chaincode`

Folder ini berisi source smart contract.

| Path | Fungsi |
|---|---|
| `chaincode/sc01-registration/main.go` | Smart contract registrasi aktor dan batch |
| `chaincode/sc02-certification/main.go` | Smart contract sertifikasi kualitas |
| `chaincode/sc03-compliance/main.go` | Smart contract compliance dan traceability |
| `chaincode/sc04-lc-settlement/main.go` | Smart contract Letter of Credit dan settlement |
| `go.mod` | Definisi module Go untuk masing-masing chaincode |
| `go.sum` | Lock checksum dependency Go |
| `vendor/` | Dependency Go yang sudah di-vendor agar packaging chaincode lebih stabil |

### 5.4 Folder `caliper-workspace`

Folder ini berisi konfigurasi benchmark Hyperledger Caliper.

| Path | Fungsi |
|---|---|
| `caliper-workspace/package.json` | Script npm dan dependency Caliper |
| `caliper-workspace/package-lock.json` | Lock dependency npm |
| `caliper-workspace/generate_profiles.py` | Generate connection profile untuk semua organisasi |
| `caliper-workspace/README.md` | Dokumentasi khusus benchmark Caliper |
| `caliper-workspace/benchmarks/bench-config-paper.yaml` | Konfigurasi benchmark final untuk paper |
| `caliper-workspace/benchmarks/bench-config.yaml` | Konfigurasi benchmark eksperimen lengkap |
| `caliper-workspace/benchmarks/debug-register-actor.yaml` | Benchmark minimal untuk debug koneksi |
| `caliper-workspace/networks/networkConfig.yaml` | Network config utama Caliper |
| `caliper-workspace/networks/connectionProfile-*.yaml` | Connection profile per organisasi |
| `caliper-workspace/workloads/registerActor.js` | Workload untuk `RegisterActor` |
| `caliper-workspace/workloads/registerBatch.js` | Workload untuk `RegisterBatch` |
| `caliper-workspace/workloads/issueCertificate.js` | Workload untuk `IssueCertificate` |
| `caliper-workspace/workloads/recordCompliance.js` | Workload untuk `RecordCompliance` |
| `caliper-workspace/workloads/issueLC.js` | Workload untuk `IssueLC` |
| `caliper-workspace/workloads/settleLC.js` | Workload untuk `SettleLC` |
| `caliper-workspace/workloads/queryActor.js` | Workload eksperimen untuk query `GetActor` |
| `caliper-workspace/report.html` | Output report Caliper, generated file |

## 6. Smart Contract dan Fungsi

### 6.1 `sc01-registration`

Fungsi:

| Fungsi | Deskripsi |
|---|---|
| `RegisterActor` | Mendaftarkan aktor supply chain |
| `GetActor` | Mengambil data aktor berdasarkan ID |
| `RegisterBatch` | Mendaftarkan batch komoditas |
| `GetBatch` | Mengambil data batch berdasarkan ID |

Data utama:

- `Actor`: ID, name, role, location, createdAt
- `Batch`: batchId, ownerId, cropType, quantity, status, createdAt

Catatan penting: `RegisterBatch` mengecek apakah `ownerId` sudah ada. Karena itu workload benchmark membuat owner actor terlebih dahulu.

### 6.2 `sc02-certification`

Fungsi:

| Fungsi | Deskripsi | Akses |
|---|---|---|
| `IssueCertificate` | Menerbitkan sertifikat kualitas | `RegulatorMSP` |
| `GetCertificate` | Mengambil data sertifikat | Semua organisasi |
| `RevokeCertificate` | Mencabut sertifikat | `RegulatorMSP` |

Smart contract ini menggunakan pengecekan MSP ID agar hanya Regulator yang dapat menerbitkan dan mencabut sertifikat.

### 6.3 `sc03-compliance`

Fungsi:

| Fungsi | Deskripsi | Akses |
|---|---|---|
| `RecordCompliance` | Mencatat kepatuhan atau traceability | `AggregatorMSP` atau `ProcessorMSP` |
| `GetComplianceRecord` | Mengambil data compliance | Semua organisasi |

Smart contract ini cocok untuk mencatat tahap supply chain seperti collection, processing, packaging, dan shipping.

### 6.4 `sc04-lc-settlement`

Fungsi:

| Fungsi | Deskripsi | Akses |
|---|---|---|
| `IssueLC` | Menerbitkan Letter of Credit | `BuyerMSP` |
| `SettleLC` | Menyelesaikan Letter of Credit | `BuyerMSP` atau `RegulatorMSP` |
| `GetLC` | Mengambil data Letter of Credit | Semua organisasi |

Smart contract ini memodelkan pembayaran antara Buyer dan Supplier.

## 7. Cara Membuat Project dari Awal

Bagian ini menjelaskan urutan pembuatan project secara konseptual dan praktis.

### 7.1 Buat Struktur Folder

```powershell
mkdir agritrace-workspace
cd agritrace-workspace
mkdir agritrace-network
mkdir chaincode
mkdir caliper-workspace
mkdir fabric-bin
```

### 7.2 Siapkan Konfigurasi Fabric

Di dalam `agritrace-network`, buat:

- `crypto-config.yaml` untuk mendefinisikan orderer dan peer organization.
- `configtx.yaml` untuk mendefinisikan genesis block, channel, consortium, dan policy.
- `docker-compose.yaml` untuk menjalankan container Fabric.
- Folder `scripts/` untuk command channel dan deploy.

Konfigurasi organisasi dibuat dengan 5 peer org dan 1 orderer org.

### 7.3 Siapkan Chaincode

Setiap smart contract dibuat sebagai module Go terpisah:

```text
chaincode/sc01-registration/
chaincode/sc02-certification/
chaincode/sc03-compliance/
chaincode/sc04-lc-settlement/
```

Di setiap folder chaincode, minimal ada:

```text
main.go
go.mod
go.sum
vendor/
```

Setelah menulis chaincode, jalankan vendor dependency:

```powershell
cd chaincode\sc01-registration
go mod tidy
go mod vendor
```

Ulangi untuk chaincode lainnya.

### 7.4 Generate Crypto dan Channel Artifact

Masuk ke folder network:

```powershell
cd agritrace-network
python .\reset_crypto.py
```

Script ini melakukan:

1. Menghapus crypto lama.
2. Menghapus channel artifact lama.
3. Menjalankan `cryptogen`.
4. Menjalankan `configtxgen` untuk genesis block.
5. Menjalankan `configtxgen` untuk channel transaction.

Output penting:

```text
agritrace-network/crypto-config/
agritrace-network/channel-artifacts/genesis.block
agritrace-network/channel-artifacts/channel.tx
```

### 7.5 Jalankan Network

```powershell
docker compose -f .\docker-compose.yaml up -d
docker ps
```

Pastikan container berikut aktif:

- `orderer1.agritrace.com`
- `orderer2.agritrace.com`
- `orderer3.agritrace.com`
- `orderer4.agritrace.com`
- `peer0.farmer.agritrace.com`
- `peer0.aggregator.agritrace.com`
- `peer0.processor.agritrace.com`
- `peer0.regulator.agritrace.com`
- `peer0.buyer.agritrace.com`
- `cli`
- semua CouchDB terkait

### 7.6 Buat Channel dan Join Peer

```powershell
docker exec cli bash -c "cd /opt/gopath/src/github.com/hyperledger/fabric/peer && bash ./scripts/join_channel.sh"
```

Script ini:

1. Membuat channel `agritracechannel`.
2. Join peer Farmer.
3. Join peer Aggregator.
4. Join peer Processor.
5. Join peer Regulator.
6. Join peer Buyer.

### 7.7 Deploy Smart Contract

Deploy semua chaincode:

```powershell
docker exec cli bash -c "cd /opt/gopath/src/github.com/hyperledger/fabric/peer && bash ./scripts/deploy_cc.sh"
docker exec cli bash -c "cd /opt/gopath/src/github.com/hyperledger/fabric/peer && bash ./scripts/deploy_sc02.sh"
docker exec cli bash -c "cd /opt/gopath/src/github.com/hyperledger/fabric/peer && bash ./scripts/deploy_sc03.sh"
docker exec cli bash -c "cd /opt/gopath/src/github.com/hyperledger/fabric/peer && bash ./scripts/deploy_sc04.sh"
```

Setiap script deploy melakukan:

1. Package chaincode.
2. Install chaincode ke semua peer.
3. Query package ID.
4. Approve chaincode untuk semua organisasi.
5. Check commit readiness.
6. Commit chaincode definition.
7. Invoke transaksi uji.

Cek chaincode committed:

```powershell
docker exec cli bash -c "peer lifecycle chaincode querycommitted -C agritracechannel"
```

## 8. Menjalankan Benchmark Caliper

Masuk ke workspace Caliper:

```powershell
cd D:\semester6\sister\agritrace-workspace\caliper-workspace
```

Install dependency:

```powershell
npm install
```

Generate connection profile:

```powershell
npm run profiles
```

Validasi workload:

```powershell
npm run lint:workloads
```

Jalankan benchmark final untuk paper:

```powershell
npm run benchmark:paper
```

Report akan dibuat di:

```text
caliper-workspace/report.html
```

## 9. Command NPM Caliper

| Command | Fungsi |
|---|---|
| `npm run profiles` | Generate ulang connection profile Fabric |
| `npm run lint:workloads` | Validasi sintaks workload JavaScript |
| `npm run test` | Flow-only test Caliper tanpa benchmark penuh |
| `npm run debug:register-actor` | Debug satu transaksi `RegisterActor` |
| `npm run benchmark` | Benchmark eksperimen lengkap termasuk query |
| `npm run benchmark:paper` | Benchmark final untuk paper |

Gunakan `benchmark:paper` untuk laporan akhir karena konfigurasi ini hanya berisi round yang sudah valid.

## 10. Skenario Benchmark Final

| Round | Chaincode | Function | Jumlah Transaksi | Target TPS |
|---|---|---|---:|---:|
| `register-actor` | `sc01-registration` | `RegisterActor` | 500 | 50 |
| `register-batch` | `sc01-registration` | `RegisterBatch` | 500 | 50 |
| `issue-certificate` | `sc02-certification` | `IssueCertificate` | 300 | 30 |
| `record-compliance` | `sc03-compliance` | `RecordCompliance` | 300 | 30 |
| `issue-lc` | `sc04-lc-settlement` | `IssueLC` | 300 | 30 |
| `settle-lc` | `sc04-lc-settlement` | `SettleLC` | 50 | 10 |

## 11. Hasil Benchmark Valid Terakhir

Hasil valid terakhir untuk transaksi utama:

| Round | Success | Fail | Send Rate (TPS) | Max Latency (s) | Min Latency (s) | Avg Latency (s) | Throughput (TPS) |
|---|---:|---:|---:|---:|---:|---:|---:|
| `register-actor` | 500 | 0 | 50.3 | 0.14 | 0.01 | 0.02 | 50.3 |
| `register-batch` | 500 | 0 | 50.2 | 0.22 | 0.01 | 0.02 | 50.1 |
| `issue-certificate` | 300 | 0 | 30.2 | 0.13 | 0.01 | 0.02 | 30.2 |
| `record-compliance` | 300 | 0 | 30.3 | 0.06 | 0.01 | 0.01 | 30.3 |
| `issue-lc` | 300 | 0 | 30.2 | 0.07 | 0.01 | 0.01 | 30.2 |
| `settle-lc` | 50 | 0 | 10.9 | 0.06 | 0.01 | 0.02 | 10.9 |

Interpretasi:

- Semua transaksi utama memiliki `Fail = 0`.
- Throughput aktual mendekati target TPS.
- Latency rata-rata berada pada rentang `0.01` sampai `0.02` detik.
- Hasil ini dapat digunakan untuk tabel evaluasi performa pada paper.

## 12. Catatan Tentang Query Actor

Konfigurasi eksperimen `bench-config.yaml` masih memiliki round `query-actor`. Namun pada pengujian terakhir, round ini menghasilkan:

```text
Success = 0
Fail = 500
```

Karena itu, jangan gunakan metrik `query-actor` untuk paper sampai konfigurasi query handler diperbaiki dan hasilnya menunjukkan `Fail = 0`.

Kalimat aman untuk paper:

```text
Benchmark dilakukan pada enam skenario transaksi utama smart contract AgriTrace-ID. Seluruh skenario mencapai success rate 100%, sedangkan skenario query dipisahkan dari evaluasi akhir karena membutuhkan konfigurasi query handler tambahan.
```

## 13. Cara Mengubah Project

### 13.1 Mengubah Smart Contract

Edit file:

```text
chaincode/sc01-registration/main.go
chaincode/sc02-certification/main.go
chaincode/sc03-compliance/main.go
chaincode/sc04-lc-settlement/main.go
```

Setelah mengubah chaincode:

1. Update version atau sequence di script deploy terkait.
2. Package dan deploy ulang chaincode.
3. Jalankan benchmark ulang jika perlu.

### 13.2 Mengubah Topologi Network

Edit:

```text
agritrace-network/docker-compose.yaml
agritrace-network/crypto-config.yaml
agritrace-network/configtx.yaml
```

Jika mengubah organisasi, MSP, channel policy, atau orderer, lakukan reset ulang dari awal:

```powershell
cd agritrace-network
docker compose -f .\docker-compose.yaml down -v
python .\reset_crypto.py
docker compose -f .\docker-compose.yaml up -d
```

Lalu join channel dan deploy chaincode ulang.

### 13.3 Mengubah Benchmark

Edit konfigurasi benchmark final:

```text
caliper-workspace/benchmarks/bench-config-paper.yaml
```

Parameter umum:

- `workers.number`: jumlah worker paralel.
- `txNumber`: jumlah transaksi.
- `rateControl.opts.tps`: target TPS.

Edit workload payload:

```text
caliper-workspace/workloads/*.js
```

Setelah mengubah workload:

```powershell
npm run lint:workloads
npm run benchmark:paper
```

## 14. Troubleshooting

| Error | Penyebab Umum | Solusi |
|---|---|---|
| `no peer combination can satisfy the endorsement policy` | Endorsement peer tidak memenuhi policy mayoritas | Pastikan workload menargetkan 3 peer dan semua chaincode sudah committed |
| `spawn EPERM` | Windows memblokir spawn process Node.js | Jalankan dari PowerShell/Git Bash normal, cek antivirus atau file lock |
| `chaincode response 500 ... does not exist` | Data setup belum tersedia saat transaksi lanjutan dijalankan | Tambahkan setup/polling/jeda commit pada workload |
| `query-actor` gagal | Query handler Caliper belum stabil untuk konfigurasi ini | Gunakan `benchmark:paper`, jangan pakai metrik query |
| `Cannot find module` | Dependency npm belum lengkap | Jalankan `npm install` dari `caliper-workspace` |
| `cryptogen not found` | Fabric binary tidak ada di path yang diharapkan | Pastikan `fabric-bin/fabric-samples/bin/` tersedia |

## 15. Shutdown dan Reset

Matikan jaringan tanpa hapus volume:

```powershell
cd agritrace-network
docker compose -f .\docker-compose.yaml down
```

Matikan jaringan dan hapus volume:

```powershell
docker compose -f .\docker-compose.yaml down -v
```

Reset crypto dan channel artifact:

```powershell
python .\reset_crypto.py
```
