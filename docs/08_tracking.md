# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## TRACKING LAUNDRY MODULE SPECIFICATION

---

## Endpoint : `GET /track/{inv}`

### Description :

Endpoint ini adalah pintu akses publik satu-satunya yang memungkinkan pelanggan memantau progres pengerjaan cucian secara _real-time_ (Near Real-time melalui metode _refresh_) tanpa perlu melakukan _login_. Sistem hanya menampilkan informasi esensial operasional dan melakukan **Data Masking** untuk melindungi privasi pelanggan dari akses pihak luar.

### Role Based Access Control (RBAC) :

- `Permissions`: `Public` (Tidak memerlukan token akses/autentikasi).

### Headers :

- `Accept`: `application/json`

### Parameters :

Bagian ini menggunakan _Path Parameter_ untuk mencari data berdasarkan nomor invoice unik yang tercetak pada nota fisik pelanggan.

| Key | Type   | Location | Default | Description                                |
| --- | ------ | -------- | ------- | ------------------------------------------ |
| INV | String | Path     | -       | Nomor Invoice unik (e.g., INV-260121-001). |

```
GET /api/track/INV-260121-001
```

### 🛡️ Logic Guard (Integritas & Keamanan Publik) :

1. **Privacy Masking**: Nama lengkap pelanggan disensor (e.g., `Mpok Romlah` menjadi Mpok R\*\*\*) untuk mencegah penyalahgunaan identitas.
2. **Data Isolation**: Informasi sensitif seperti nomor telepon lengkap, alamat detail, dan rincian metode pembayaran (seperti nomor referensi bank) tidak ditampilkan pada respons publik ini.
3. **Read-Only Context**: Endpoint ini murni hanya untuk pembacaan status. Tidak ada data yang bisa dimodifikasi melalui jalur ini.
4. **Rate Limit Protection**: Mencegah upaya pencarian invoice secara massal (scraping) menggunakan mesin atau bot.

### Request Body :

Bagian ini tidak memerlukan data tambahan.

```
None (Kosong, karena method GET tidak memerlukan body request).
```

### Response Body

#### ✅ 200 OK

Data pesanan ditemukan. Respons memberikan informasi status terakhir, estimasi waktu selesai, dan riwayat pengerjaan agar pelanggan mendapatkan kepastian layanan.

```json
{
  "success": true,
  "message": "Tracking retrieved successfully",
  "data": {
    "invoice_number": "INV-260121-001",
    "customer_name": "Mpok R***",
    "status_internal": "finished-delivery",
    "payment_status": "paid",
    "estimated_ready_at": "2026-01-19 16:20:00",
    "total_price": 110000.0,
    "order_items": [
      {
        "service_name": "Cuci Kiloan Reguler",
        "qty_pieces": 40,
        "weight_kg": 10.0,
        "unit": "Kg"
      }
    ],
    "status_history": [
      {
        "previous_status": null,
        "new_status": "pending",
        "description": "Pesanan telah diterima oleh kasir",
        "created_at": "2026-01-21 08:52:36"
      },
      {
        "previous_status": "pending",
        "new_status": "in-progress",
        "description": "Pakaian sedang dalam proses pencucian",
        "created_at": "2026-01-21 08:56:08"
      },
      {
        "previous_status": "in-progress",
        "new_status": "ready-delivery",
        "description": "Selesai packing delivery",
        "created_at": "2026-01-21 08:58:05"
      },
      {
        "previous_status": "ready-delivery",
        "new_status": "being-delivered",
        "description": "Sedang mengantar laundry",
        "created_at": "2026-01-21 09:03:37"
      },
      {
        "previous_status": "being-delivered",
        "new_status": "finished-delivery",
        "description": "Diterima oleh Mpok Romlah (Lunas)",
        "created_at": "2026-01-21 09:06:14"
      }
    ]
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi jika format nomor invoice yang dimasukkan tidak memenuhi kriteria validasi (misal: karakter ilegal atau terlalu pendek).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "inv": "Invoice number format is invalid"
    }
  }
}
```

#### 🚫 404 Not Found

Nomor invoice valid secara format, namun tidak terdaftar di database. Ini berfungsi sebagai pelindung agar pola invoice tidak mudah ditebak.

```json
{
  "success": false,
  "message": "Tracking not found",
  "data": {
    "error_code": "RESOURCE_NOT_FOUND",
    "errors": null
  }
}
```

#### 🚫 429 Too Many Requests

Sistem mendeteksi aktivitas pencarian yang terlalu sering dari alamat IP yang sama dalam waktu singkat.

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

Kegagalan teknis pada server atau koneksi database saat memproses data pelacakan.

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
