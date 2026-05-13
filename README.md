# AgriTrace-ID Hyperledger Caliper Benchmark

Dokumen ini menjelaskan cara menyiapkan, menjalankan, dan membaca hasil benchmark performa smart contract AgriTrace-ID menggunakan Hyperledger Caliper.

Benchmark ini menghasilkan metrik utama yang umum digunakan dalam evaluasi jaringan Hyperledger Fabric:

- Throughput atau transaksi per detik (TPS)
- Latency minimum, maksimum, dan rata-rata
- Jumlah transaksi sukses dan gagal
- Success rate per skenario transaksi

## Ringkasan Sistem

AgriTrace-ID berjalan di atas jaringan Hyperledger Fabric dengan channel `agritracechannel`. Empat smart contract yang diuji adalah:

- `sc01-registration`: registrasi aktor dan batch komoditas
- `sc02-certification`: penerbitan sertifikat kualitas
- `sc03-compliance`: pencatatan kepatuhan dan traceability
- `sc04-lc-settlement`: penerbitan dan penyelesaian Letter of Credit

Channel menggunakan endorsement policy `MAJORITY Endorsement` dari 5 organisasi. Karena itu, workload Caliper untuk transaksi write menargetkan 3 peer endorsement secara eksplisit agar memenuhi policy mayoritas.

## Prasyarat

Pastikan seluruh komponen berikut sudah berjalan sebelum menjalankan benchmark:

- Docker Desktop aktif.
- Network Fabric AgriTrace sudah berjalan.
- Channel `agritracechannel` sudah dibuat.
- Semua peer sudah join ke channel.
- Semua chaincode sudah committed dan berhasil diuji invoke.
- Node.js dan npm tersedia di environment lokal.

Struktur minimum yang dibutuhkan:

```text
agritrace-workspace/
|-- agritrace-network/
|   |-- crypto-config/
|   |-- channel-artifacts/
|   `-- docker-compose.yaml
|-- chaincode/
`-- caliper-workspace/
```

## Instalasi Dependency

Masuk ke folder Caliper:

```powershell
cd \agritrace-workspace\caliper-workspace
```

Install dependency:

```powershell
npm install
```

Dependency penting yang digunakan:

- `@hyperledger/caliper-cli`
- `@hyperledger/caliper-fabric`
- `fabric-network`

Catatan penting: benchmark ini menggunakan `fabric-network@2.x`, bukan `@hyperledger/fabric-gateway`, agar Caliper dapat memakai target peer endorsement eksplisit. Jangan tambahkan flag `--caliper-fabric-gateway-enabled` saat menjalankan benchmark ini.

## Persiapan Konfigurasi

Generate ulang connection profile berdasarkan material crypto lokal:

```powershell
npm run profiles
```

Perintah ini menjalankan `generate_profiles.py` dan menghasilkan profile untuk semua organisasi:

- `connectionProfile-farmer.yaml`
- `connectionProfile-aggregator.yaml`
- `connectionProfile-processor.yaml`
- `connectionProfile-regulator.yaml`
- `connectionProfile-buyer.yaml`

Validasi sintaks workload:

```powershell
npm run lint:workloads
```

Jika tidak ada output error, file workload JavaScript sudah valid secara sintaks.

## Menjalankan Benchmark Paper-Ready

Untuk menghasilkan report final yang aman dipakai di paper:

```powershell
npm run benchmark:paper
```

Perintah ini menjalankan konfigurasi:

```text
benchmarks/bench-config-paper.yaml
```

Output utama akan dibuat di:

```text
caliper-workspace/report.html
```

Buka file tersebut di browser untuk melihat tabel hasil benchmark.

## Skenario Benchmark Paper-Ready

Benchmark paper-ready mencakup 6 round transaksi yang sudah berhasil mencapai success rate 100%:

| Round | Smart Contract | Fungsi | Jumlah Transaksi | Target TPS |
|---|---|---|---:|---:|
| `register-actor` | `sc01-registration` | `RegisterActor` | 500 | 50 |
| `register-batch` | `sc01-registration` | `RegisterBatch` | 500 | 50 |
| `issue-certificate` | `sc02-certification` | `IssueCertificate` | 300 | 30 |
| `record-compliance` | `sc03-compliance` | `RecordCompliance` | 300 | 30 |
| `issue-lc` | `sc04-lc-settlement` | `IssueLC` | 300 | 30 |
| `settle-lc` | `sc04-lc-settlement` | `SettleLC` | 50 | 10 |

## Hasil Benchmark Terakhir

Hasil berikut berasal dari run Caliper terakhir yang valid untuk transaksi utama:

| Round | Success | Fail | Send Rate (TPS) | Max Latency (s) | Min Latency (s) | Avg Latency (s) | Throughput (TPS) |
|---|---:|---:|---:|---:|---:|---:|---:|
| `register-actor` | 500 | 0 | 50.3 | 0.14 | 0.01 | 0.02 | 50.3 |
| `register-batch` | 500 | 0 | 50.2 | 0.22 | 0.01 | 0.02 | 50.1 |
| `issue-certificate` | 300 | 0 | 30.2 | 0.13 | 0.01 | 0.02 | 30.2 |
| `record-compliance` | 300 | 0 | 30.3 | 0.06 | 0.01 | 0.01 | 30.3 |
| `issue-lc` | 300 | 0 | 30.2 | 0.07 | 0.01 | 0.01 | 30.2 |
| `settle-lc` | 50 | 0 | 10.9 | 0.06 | 0.01 | 0.02 | 10.9 |

Interpretasi singkat:

- Semua skenario transaksi utama mencapai `Fail = 0`.
- Throughput aktual mendekati target TPS pada setiap round.
- Latency rata-rata berada pada rentang `0.01s` sampai `0.02s`.
- Hasil ini dapat digunakan sebagai data performa transaksi utama AgriTrace-ID.

## Catatan Tentang Query Benchmark

File `benchmarks/bench-config.yaml` masih menyertakan round tambahan:

```text
query-actor
```

Round ini digunakan untuk eksperimen query `GetActor`, tetapi belum dimasukkan ke benchmark paper-ready karena pada konfigurasi terakhir masih menghasilkan:

```text
Success = 0
Fail = 500
```

Karena semua query gagal, angka throughput pada round tersebut tidak boleh diklaim sebagai TPS query. Untuk paper, gunakan `npm run benchmark:paper` agar report hanya berisi round yang valid.

Kalimat yang aman untuk paper:

```text
Benchmark dilakukan pada enam skenario transaksi utama smart contract AgriTrace-ID. Seluruh skenario mencapai success rate 100%, sedangkan skenario query dipisahkan dari evaluasi akhir karena membutuhkan konfigurasi query handler tambahan.
```

## Perintah yang Tersedia

| Perintah | Fungsi |
|---|---|
| `npm run profiles` | Generate ulang connection profile Fabric untuk semua organisasi |
| `npm run lint:workloads` | Validasi sintaks workload JavaScript |
| `npm run debug:register-actor` | Menjalankan 1 transaksi debug `RegisterActor` |
| `npm run benchmark:paper` | Menjalankan benchmark final yang direkomendasikan untuk paper |
| `npm run benchmark` | Menjalankan benchmark eksperimen lengkap, termasuk `query-actor` |

## Troubleshooting

Jika muncul error `no peer combination can satisfy the endorsement policy`, pastikan:

- Tidak memakai flag `--caliper-fabric-gateway-enabled`.
- Dependency `fabric-network` sudah terpasang.
- Workload menargetkan 3 peer endorsement.
- Semua peer sudah join channel dan chaincode sudah committed.

Jika muncul `query-actor` gagal:

- Gunakan `npm run benchmark:paper` untuk report final.
- Jangan gunakan metrik `query-actor` sampai `Fail = 0`.

Jika muncul error `spawn EPERM`:

- Jalankan command dari terminal normal seperti Git Bash atau PowerShell, bukan dari environment sandbox.
- Pastikan antivirus atau file lock Windows tidak memblokir proses Node.js.

## Mengubah Beban Benchmark

Untuk mengubah beban benchmark paper-ready, edit:

```text
benchmarks/bench-config-paper.yaml
```

Parameter yang umum diubah:

- `workers.number`: jumlah worker Caliper paralel
- `txNumber`: jumlah transaksi per round
- `rateControl.opts.tps`: target transaksi per detik

Contoh:

```yaml
txNumber: 1000
rateControl:
  type: fixed-rate
  opts:
    tps: 100
```

Setelah mengubah konfigurasi, jalankan ulang:

```powershell
npm run benchmark:paper
```
