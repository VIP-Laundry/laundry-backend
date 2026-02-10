# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## REPORTS MODULE SPECIFICATION

---

## Endpoint : `GET /reports/dashboard`

### Description :

Endpoint ini menyediakan ringkasan eksekutif mengenai aktivitas bisnis harian secara up-to-date (berdasarkan data terbaru saat request dikirim). Data yang disajikan mencakup total pendapatan terverifikasi, volume pesanan baru, serta pemetaan beban kerja operasional. Informasi ini dirancang agar Owner dapat memantau performa toko secara instan dan akurat tanpa perlu membuka laporan detail satu per satu.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`
- `Restriction`: `Cashier, Staff, Courier` (tidak memiliki hak akses.)

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Laporan secara default menyajikan data hari ini. Namun, sistem mendukung filter tanggal untuk peninjauan data historis.

| Key  | Type   | Location | Default | Description                                  |
| ---- | ------ | -------- | ------- | -------------------------------------------- |
| date | String | Query    | Today   | Format: `YYYY-MM-DD` (Contoh: `2026-01-21`). |

```
GET /api/reports/dashboard?date=2026-01-21
```

### 🛡️ Logic Guard (Aturan Agregasi & Integritas) :

1. **Strict Role Enforcement**: Backend wajib memverifikasi bahwa `role` dalam JWT adalah `owner`. Jika tidak, kembalikan `403 Forbidden`.
2. **Verified Revenue (No Debt Policy)**: Variabel `today_revenue` dihitung menggunakan fungsi agregasi $SUM(total\_price)$ hanya untuk pesanan dengan `payment_status = 'paid'`.
3. **Liquidity Insight**: `pending_payment_value` dihitung dari $SUM(total\_price)$ untuk pesanan berstatus `unpaid` atau `cod_pending`.
4. **Operational Counting**: Status pengerjaan dihitung menggunakan `COUNT` atomik berdasarkan `status_internal` untuk memberikan gambaran beban kerja di workshop.
5. **Data Persistence**: Jika tidak ada data pada tanggal yang dipilih, server mengembalikan nilai `0.0` (Float) atau `0` (Integer) dalam respons `200 OK`.

### Request Body :

```
None (Kosong, karena method GET tidak memerlukan body request).
```

### Response Body :

Bagian ini mengembalikan objek statistik yang dikelompokkan berdasarkan kategori finansial dan operasional.

#### ✅ 200 OK

Data berhasil dikalkulasi. Jika pada tanggal yang dipilih belum ada transaksi, sistem akan mengembalikan nilai 0.0 atau 0 (bukan error 404).

```json
{
  "success": true,
  "message": "Dashboard statistics retrieved successfully",
  "data": {
    "report_date": "2026-01-21",
    "financials": {
      "today_revenue": 2450000.0,
      "pending_payment_value": 850000.0
    },
    "order_stats": {
      "total_new_orders": 15,
      "orders_in_progress": 8,
      "orders_ready": 5,
      "orders_completed": 12
    },
    "customer_stats": {
      "new_customers": 3,
      "total_active_customers": 15
    }
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi jika format tanggal pada parameter query tidak valid atau melanggar standar ISO.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "date": "Date must be in YYYY-MM-DD format"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Terjadi jika sesi login telah berakhir atau token tidak valid.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": {
    "error_code": "UNAUTHORIZED_ACCESS",
    "errors": null
  }
}
```

#### 🚫 403 Forbidden

Terjadi jika staf operasional (Staff/Courier) mencoba mengakses dashboard manajemen.

```json
{
  "success": false,
  "message": "Your role does not have permission",
  "data": {
    "error_code": "FORBIDDEN_ACCESS",
    "errors": null
  }
}
```

#### 🚫 429 Too Many Requests

Perlindungan server dari beban kueri agregasi yang terlalu masif dalam waktu singkat.

```json
{
  "success": false,
  "message": "Too many requests, please try again later",
  "data": {
    "error_code": "RATE_LIMIT_EXCEEDED",
    "errors": null
  }
}
```

#### 🔥 500 Internal Server Error

Kegagalan sistem saat melakukan kueri kompleks SUM dan COUNT pada database.

```json
{
  "success": false,
  "message": "An unexpected server error occurred",
  "data": {
    "error_code": "INTERNAL_SERVER_ERROR",
    "errors": null
  }
}
```

---

## Endpoint : `GET /reports/revenue`

### Description :

Endpoint ini digunakan untuk menarik laporan pendapatan historis berdasarkan rentang waktu tertentu. Berbeda dengan dashboard yang hanya menampilkan ringkasan harian, endpoint ini melakukan agregasi data (penjumlahan) untuk memberikan gambaran tren keuangan mingguan, bulanan, atau tahunan. Data ini merupakan landasan bagi Owner dalam mengevaluasi pertumbuhan bisnis VIP Laundry.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`
- `Restriction`: Akses ditutup total bagi `Cashier, Staff, dan Courier`. Hanya pemegang hak akses tertinggi yang dapat melihat detail akumulasi omzet perusahaan.

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Parameter ini wajib diisi untuk menentukan batasan data yang akan diproses oleh database.

| Key        | Type   | Location | Default | Description                                 |
| ---------- | ------ | -------- | ------- | ------------------------------------------- |
| start_date | String | Query    | -       | Tanggal awal laporan (Format: YYYY-MM-DD).  |
| end_date   | String | Query    | -       | Tanggal akhir laporan (Format: YYYY-MM-DD). |

```
GET /api/reports/revenue?start_date=2026-01-01&end_date=2026-01-31
```

### 🛡️ Logic Guard (Aturan Agregasi & Integritas) :

1. Strict Owner Check: Backend memverifikasi role pada klaim JWT. Jika bukan owner, kembalikan 403 Forbidden.
2. Date Range Validation:
   - Memastikan format tanggal adalah YYYY-MM-DD.
   - Memastikan $start\_date \le end\_date$. Jika terbalik, kembalikan 400 Bad Request.
3. Revenue Filtering (No Debt Policy): Hanya menjumlahkan pesanan yang memiliki payment_status = 'paid'.
4. SQL Aggregation Logic: Menggunakan kueri $SUM(total\_price)$ dan $GROUP BY$ tanggal agar Owner bisa melihat grafik pendapatan per hari di dalam rentang waktu yang dipilih.
5. Decimal Precision: Seluruh hasil perhitungan finansial wajib bertipe data Float untuk menghindari pembulatan yang tidak akurat.

### Request Body :

```
None (Kosong, karena method GET tidak memerlukan body request).
```

### Response Body :

#### ✅ 200 OK

Data berhasil dikalkulasi. Respons menyertakan total keseluruhan dan rincian per hari untuk mempermudah pembuatan grafik (Chart).

```json
{
  "success": true,
  "message": "Revenue report generated successfully",
  "data": {
    "period": {
      "start_date": "2026-01-01",
      "end_date": "2026-01-07",
      "total_days": 7
    },
    "summary": {
      "total_revenue": 15750000.0,
      "total_orders_paid": 124,
      "average_daily_revenue": 2250000.0
    },
    "daily_breakdown": [
      {
        "date": "2026-01-01",
        "revenue": 2000000.0,
        "order_count": 18
      },
      {
        "date": "2026-01-02",
        "revenue": 2500000.0,
        "order_count": 22
      }
    ]
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi jika parameter tanggal tidak lengkap, format salah, atau rentang tanggal tidak logis.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "date": "Date must be in YYYY-MM-DD format"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Terjadi jika token tidak ada atau sudah kedaluwarsa.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": {
    "error_code": "UNAUTHORIZED_ACCESS",
    "errors": null
  }
}
```

#### 🚫 403 Forbidden

Terjadi jika peran pengguna bukan Owner.

```json
{
  "success": false,
  "message": "Your role does not have permission",
  "data": {
    "error_code": "FORBIDDEN_ACCESS",
    "errors": null
  }
}
```

#### 🚫 429 Too Many Requests

Terjadi jika ada permintaan laporan massal dalam waktu yang sangat singkat.

```json
{
  "success": false,
  "message": "Too many requests, please try again later",
  "data": {
    "error_code": "RATE_LIMIT_EXCEEDED",
    "errors": null
  }
}
```

#### 🔥 500 Internal Server Error

Kegagalan teknis saat memproses agregasi data besar di MySQL.

```json
{
  "success": false,
  "message": "An unexpected server error occurred",
  "data": {
    "error_code": "INTERNAL_SERVER_ERROR",
    "errors": null
  }
}
```

---

## Endpoint : `GET /reports/payments`

### Description :

Endpoint ini menyediakan laporan rincian penerimaan dana berdasarkan Metode Pembayaran (Cash, Transfer, QRIS, dll) dalam rentang waktu tertentu. Laporan ini merupakan instrumen utama bagi Owner untuk melakukan audit harian guna memastikan uang tunai yang terkumpul di lapangan (termasuk hasil penagihan COD oleh kurir) cocok dengan data yang tercatat di sistem.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`
- `Restriction`: `Cashier, Staff, dan Courier` dilarang mengakses untuk mencegah manipulasi data saat proses serah terima uang tunai.

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Parameter ini wajib diisi untuk menentukan batasan data yang akan diproses oleh database.

| Key        | Type   | Location | Default | Description                               |
| ---------- | ------ | -------- | ------- | ----------------------------------------- |
| start_date | String | Query    | -       | Tanggal awal audit (Format: YYYY-MM-DD).  |
| end_date   | String | Query    | -       | Tanggal akhir audit (Format: YYYY-MM-DD). |

```
GET /api/reports/payments?start_date=2026-01-21&end_date=2026-01-21
```

### 🛡️ Logic Guard (Aturan Agregasi & Integritas) :

1. Strict Owner Check: Backend wajib menolak akses jika role user bukan owner.
2. Paid-Only Filter: Hanya data dengan payment_status = 'paid' yang dihitung. Transaksi yang masih unpaid atau cod_pending tidak boleh muncul dalam laporan audit dana masuk.
3. Method Aggregation: Menggunakan fungsi SQL GROUP BY payment_method untuk memisahkan total dana yang masuk lewat jalur fisik (Cash) dan jalur digital (QRIS/Transfer).
4. Audit Formula: Sistem menghitung total keseluruhan dengan rumus:$$Total\_Collected = \sum Cash + \sum Transfer + \sum QRIS$$
5. Data Precision: Menggunakan tipe data Float untuk seluruh nilai nominal uang.

### Request Body :

```
None (Kosong, karena method GET tidak memerlukan body request).
```

### Response Body :

#### ✅ 200 OK

Respons memberikan ringkasan total dan rincian per metode untuk mempermudah pemeriksaan fisik uang.

```json
{
  "success": true,
  "message": "Payment report generated successfully",
  "data": {
    "audit_period": {
      "start_date": "2026-01-21",
      "end_date": "2026-01-21",
      "total_days": 1
    },
    "summary": {
      "total_collected": 3500000.0,
      "total_transactions": 25
    },
    "breakdown": [
      {
        "payment_method": "cash",
        "total_amount": 2000000.0,
        "transaction_count": 15,
        "description": "Total uang tunai yang harus ada di kasir/kurir"
      },
      {
        "payment_method": "qris",
        "total_amount": 1000000.0,
        "transaction_count": 7,
        "description": "Total dana masuk melalui sistem QRIS"
      },
      {
        "payment_method": "transfer",
        "total_amount": 500000.0,
        "transaction_count": 3,
        "description": "Total dana masuk melalui mutasi bank"
      }
    ]
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi jika format tanggal salah atau rentang tanggal tidak valid.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "start_date": "Invalid date format, use YYYY-MM-DD"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Terjadi jika token tidak valid atau sesi telah berakhir.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": {
    "error_code": "UNAUTHORIZED_ACCESS",
    "errors": null
  }
}
```

#### 🚫 403 Forbidden

Terjadi jika Kasir/Staf mencoba mengintip laporan audit keuangan.

```json
{
  "success": false,
  "message": "Your role does not have permission",
  "data": {
    "error_code": "FORBIDDEN_ACCESS",
    "errors": null
  }
}
```

#### 🚫 429 Too Many Requests

Terjadi jika ada permintaan laporan massal dalam waktu yang sangat singkat.

```json
{
  "success": false,
  "message": "Too many requests, please try again later",
  "data": {
    "error_code": "RATE_LIMIT_EXCEEDED",
    "errors": null
  }
}
```

#### 🔥 500 Internal Server Error

Kegagalan kueri agregasi saat memproses data besar.

```json
{
  "success": false,
  "message": "An unexpected server error occurred",
  "data": {
    "error_code": "INTERNAL_SERVER_ERROR",
    "errors": null
  }
}
```

---

## Endpoint : `GET /reports/employees`

### Description :

Endpoint ini menyediakan laporan produktivitas karyawan berdasarkan jumlah tugas yang diselesaikan dalam rentang waktu tertentu. Laporan ini membedakan kontribusi antara Staff (Workshop) yang melakukan pemrosesan pakaian dan Courier (Logistik) yang menyelesaikan pengantaran. Data ini digunakan oleh Owner sebagai dasar evaluasi kinerja, pemberian bonus, atau penyesuaian beban kerja.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`
- `Restriction`: `Cashier, Staff, dan Courier` dilarang mengakses untuk menjaga kerahasiaan data performa antar rekan kerja dan menghindari persaingan tidak sehat di level operasional.

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Parameter ini wajib diisi untuk menentukan batasan data yang akan diproses oleh database.

| Key        | Type   | Location | Default | Description                                         |
| ---------- | ------ | -------- | ------- | --------------------------------------------------- |
| start_date | String | Query    | -       | Tanggal awal periode laporan (Format: YYYY-MM-DD).  |
| end_date   | String | Query    | -       | Tanggal akhir periode laporan (Format: YYYY-MM-DD). |

```
GET /api/reports/employees?start_date=2026-01-01&end_date=2026-01-21
```

### 🛡️ Logic Guard (Aturan Agregasi & Integritas) :

1. Strict Owner Policy: Hanya pengguna dengan klaim role: owner pada JWT yang diizinkan memproses permintaan.
2. Activity Tracking Logic:
   - Staff Performance: Dihitung dari jumlah transaksi di mana karyawan tersebut tercatat melakukan perubahan status ke in-progress, washing, atau ready.
   - Courier Performance: Dihitung dari jumlah transaksi di mana karyawan tersebut tercatat sebagai pengantar pada status finished-delivery.
3. Cross-Reference Integrity: Sistem melakukan JOIN antara tabel users dan tabel order_history (atau kolom updated_by pada log status) untuk memastikan data akurat per individu.
4. Date Range Validation: Memastikan format tanggal benar dan rentang waktu logis (tidak mencari data masa depan).

### Request Body :

```
None (Kosong, karena method GET tidak memerlukan body request).
```

### Response Body :

#### ✅ 200 OK

Respons menyertakan peringkat produktivitas karyawan yang dikelompokkan berdasarkan peran mereka.

```json
{
  "success": true,
  "message": "Employee productivity report generated successfully",
  "data": {
    "report_period": {
      "start_date": "2026-01-01",
      "end_date": "2026-01-07",
      "total_days": 7
    },
    "cashier_performance": [
      {
        "employee_id": 101,
        "name": "Budi Kasir",
        "role": "cashier",
        "total_activity": 210,
        "activity": "Orders Created",
        "average_per_day": 30.0
      }
    ],
    "staff_performance": [
      {
        "employee_id": 102,
        "name": "Siti Aminah",
        "role": "staff",
        "total_activity": 70,
        "activity": "Orders Processed",
        "average_per_day": 10.0
      }
    ],
    "courier_performance": [
      {
        "employee_id": 201,
        "name": "Andi Kurir",
        "role": "courier",
        "total_activity": 105,
        "activity": "Deliveries Completed",
        "average_per_day": 15.0
      }
    ]
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi jika parameter filter tidak lengkap atau format tanggal salah.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "end_date": "Invalid date format, use YYYY-MM-DD"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Terjadi jika token tidak valid atau sesi telah berakhir.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": {
    "error_code": "UNAUTHORIZED_ACCESS",
    "errors": null
  }
}
```

#### 🚫 403 Forbidden

Terjadi jika Kasir/Staf mencoba mengintip laporan audit keuangan.

```json
{
  "success": false,
  "message": "Your role does not have permission",
  "data": {
    "error_code": "FORBIDDEN_ACCESS",
    "errors": null
  }
}
```

#### 🚫 429 Too Many Requests

Terjadi jika ada permintaan laporan massal dalam waktu yang sangat singkat.

```json
{
  "success": false,
  "message": "Too many requests, please try again later",
  "data": {
    "error_code": "RATE_LIMIT_EXCEEDED",
    "errors": null
  }
}
```

#### 🔥 500 Internal Server Error

Kegagalan saat melakukan kueri JOIN dan COUNT yang berat pada tabel riwayat transaksi.

```json
{
  "success": false,
  "message": "An unexpected server error occurred",
  "data": {
    "error_code": "INTERNAL_SERVER_ERROR",
    "errors": null
  }
}
```

---
