# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## DELIVERIES MODULE SPECIFICATION

---

## Endpoint : `GET /deliveries`

### Description :

Mengambil daftar pesanan yang memerlukan pengiriman. Secara default (**Task Pool**), endpoint ini menyaring data agar hanya menampilkan pesanan dengan `delivery_status = 'ready-delivery'` yang belum memiliki kurir. Data yang ditampilkan adalah ringkasan informasi untuk memudahkan kurir melakukan pemindaian tugas secara cepat.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner, cashier, courier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Bagian ini mendefinisikan parameter query opsional untuk menyaring, mengurutkan, dan membatasi daftar pengiriman yang ditampilkan.

| Key      | Type   | Location | Default    | Description                                                                |
| -------- | ------ | -------- | ---------- | -------------------------------------------------------------------------- |
| page     | Int    | Query    | 1          | Nomor halaman (Pagination)                                                 |
| per_page | Int    | Query    | 10         | Jumlah data per halaman                                                    |
| search   | String | Query    | -          | Cari                                                                       |
| status   | String | Query    | null       | Filter status (ready-delivery, being-delivered, finished-delivery).        |
| sort_by  | String | Query    | created_at | Pengurutan berdasarkan kolom tertentu (contoh: shipping_cost, created_at). |
| order    | String | Query    | -          | Arah urutan: asc (terlama ke terbaru) atau desc (terbaru ke terlama).      |

```
GET /api/deliveries?page=1&per_page=10&status=ready-delivery&sort_by=created_at&order=asc
```

### Request Body :

```
None (Kosong, karena method GET tidak memerlukan body request).
```

### Response Body :

#### ✅ 200 OK

Data berhasil diambil. Response menyertakan informasi detail dari tabel `orders` agar kurir dapat mengetahui tujuan pengiriman.

```json
{
  "success": true,
  "message": "Deliveries retrieved successfully",
  "data": [
    {
      "id": 1,
      "order_id": 45,
      "invoice_number": "INV-20260120-001",
      "customer_name": "Mpok Romlah",
      "delivery_status": "ready-delivery",
      "shipping_cost": 10000.0,
      "created_at": "2026-01-19 16:00:00",
      "updated_at": "2026-01-20 16:00:00"
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total_items": 1,
    "total_pages": 1
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi kesalahan pada parameter input (misal: `page` diisi teks).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "page": "page must be a number"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Token expired atau tidak menyertakan Header Authorization.

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

User mencoba mengakses tapi tidak memiliki role yang diizinkan.

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

User melakukan request terlalu cepat dalam waktu singkat.

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

Terjadi kesalahan pada query database atau logika server.

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

## Endpoint : `GET /deliveries/{id}`

### Description :

Mengambil informasi detail untuk satu data pengiriman tertentu berdasarkan ID. Endpoint ini menggabungkan data dari tabel `deliveries`, `orders`, dan `customers` untuk menyediakan informasi lengkap bagi kurir, termasuk nomor telepon pelanggan dan catatan khusus pesanan.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner, cashier, courier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Bagian ini menggunakan _Path Parameter_ untuk menentukan data pengiriman spesifik yang ingin diakses.

| Key | Type | Location | Default | Description                                       |
| --- | ---- | -------- | ------- | ------------------------------------------------- |
| id  | Int  | Path     | -       | ID unik dari tabel deliveries yang ingin diambil. |

```
GET /api/deliveries/1
```

### Request Body :

```
None (Kosong, karena method GET tidak memerlukan body request).
```

### Responses Body :

#### ✅ 200 OK

Data berhasil ditemukan. Response ini menyediakan informasi lengkap untuk eksekusi pengiriman oleh kurir.

```json
{
  "success": true,
  "message": "Delivery retrieved successfully",
  "data": {
    "id": 1,
    "order_id": 45,
    "invoice_number": "INV-20260120-001",
    "customer_name": "Mpok Romlah",
    "customer_phone": "081234567890",
    "customer_address": "Jl. Merpati 12",
    "delivery_status": "ready-delivery",
    "shipping_cost": 10000.0,
    "courier_id": null,
    "courier_departed_at": null,
    "courier_arrived_at": null,
    "receiver_name": null,
    "cod_collected_amount": 0.0,
    "created_at": "2026-01-19 16:00:00",
    "updated_at": "2026-01-20 16:00:00"
  }
}
```

#### ⚠️ 400 Bad Request

Format ID yang dikirimkan tidak valid atau bukan angka.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "id": "The id must be a positive integer."
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Respons ketika token otentikasi tidak valid atau tidak disertakan dalam header.

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

Respons ketika peran (role) pengguna tidak memiliki izin untuk mengakses detail pengiriman ini.

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

#### 🚫 404 Not Found

Respons ketika data pengiriman dengan ID yang diminta tidak ditemukan di database.

```json
{
  "success": false,
  "message": "Delivery not found",
  "data": {
    "error_code": "RESOURCE_NOT_FOUND",
    "errors": null
  }
}
```

#### 🚫 429 Too Many Requests

Respons ketika terjadi lonjakan permintaan (rate limit) dari sisi klien.

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

Respons ketika terjadi kegagalan sistem atau database saat memproses data pengiriman.

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

## Endpoint : `GET /deliveries/my-tasks`

### Description :

Mengambil daftar pengiriman yang sedang ditangani atau telah diselesaikan oleh kurir yang sedang login. Endpoint ini secara otomatis menyaring data berdasarkan `courier_id` yang diambil dari **Access Token (JWT)**, sehingga kurir hanya dapat melihat tugas milik mereka sendiri.

### Role Based Access Control (RBAC) :

- `Permissions`: `courier`
- _Catatan_ : `owner` dan `cashier` menggunakan endpoint `GET /deliveries` umum untuk memantau seluruh kurir.

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Bagian ini menggunakan _Query Parameter_ untuk membantu kurir menemukan tugas spesifik atau memfilter riwayat tugas.

| Key      | Type   | Location | Default         | Description                                                       |
| -------- | ------ | -------- | --------------- | ----------------------------------------------------------------- |
| page     | Int    | Query    | 1               | Nomor halaman untuk pagination.                                   |
| per_page | Int    | Query    | 10              | Jumlah data per halaman.                                          |
| status   | String | Query    | being-delivered | Filter: being-delivered (aktif) atau finished-delivery (riwayat). |
| search   | String | Query    | -               | Cari berdasarkan invoice_number atau customer_name.               |
| sort_by  | String | Query    | updated_at      | Pengurutan (disarankan berdasarkan update status terakhir).       |
| order    | String | Query    | desc            | Urutan: desc (terbaru) atau asc (terlama).                        |

```
GET /api/deliveries/my-tasks?status=being-delivered&page=1&per_page=10
```

### Request Body :

```
None (Kosong, karena method GET tidak memerlukan body request).
```

### Response Body :

#### ✅ 200 OK

Data berhasil diambil. Response menggunakan strategi **Brief List** untuk menjaga kecepatan aplikasi mobile kurir.

```json
{
  "success": true,
  "message": "Deliveries retrieved successfully",
  "data": [
    {
      "id": 1,
      "order_id": 45,
      "invoice_number": "INV-20260120-001",
      "customer_name": "Mpok Romlah",
      "delivery_status": "being-delivered",
      "shipping_cost": 10000.0,
      "created_at": "2026-01-19 16:00:00",
      "updated_at": "2026-01-20 16:10:00"
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total_items": 1,
    "total_pages": 1
  }
}
```

#### ⚠️ 400 Bad Request

Input parameter tidak valid (misal: status yang tidak terdaftar).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "status": "The status field must be one of: being-delivered, finished-delivery."
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Token tidak valid atau kurir belum melakukan login.

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

User mencoba mengakses namun tidak memiliki role courier.

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

User melakukan request terlalu cepat dalam waktu singkat.

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

Terjadi kegagalan koneksi database atau error pada logika filter courier_id.

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

## Endpoint : `PATCH /deliveries/{id}`

### Description :

Memperbarui status operasional pengiriman. Endpoint ini berfungsi sebagai **State Machine** yang mengelola transisi status dari posisi awal `ready-delivery` (setelah dipicu modul Orders) hingga penyelesaian. Sistem secara otomatis melakukan validasi alur dan mencatat timestamp logistik untuk akurasi data.

1. Transisi ke `being-delivered`: Digunakan saat kurir mengambil tugas (_Pick Up_). Sistem mencatat `courier_departed_at` dan mengikat `courier_id` dengan user yang sedang login.
2. Transisi ke `finished-delivery`: Digunakan saat pesanan sampai. Sistem mencatat `courier_arrived_at` dan melakukan **Double Update** pada tabel `orders` (menyelesaikan status pesanan secara global).

### Role Based Access Control (RBAC) :

- `Permissions`: `owner, cashier, courier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

Bagian ini menggunakan _Path Parameter_ untuk menentukan data pengiriman spesifik yang akan diperbarui.

| Key | Type | Location | Default | Description                                        |
| --- | ---- | -------- | ------- | -------------------------------------------------- |
| id  | Int  | Path     | -       | ID unik dari tabel deliveries yang ingin diupdate. |

```
PATCH /api/deliveries/1
```

### Request Body :

Field bersifat kondisional. `delivery_status` adalah satu-satunya field mandatori.

| Key                    | Type   | Mandatory   | Description                                                     |
| ---------------------- | ------ | ----------- | --------------------------------------------------------------- |
| `delivery_status`      | String | Ya          | Target status baru (`being-delivered`atau `finished-delivery`). |
| `receiver_name`        | String | Kondisional | Wajib diisi jika status = `finished-delivery`.                  |
| `cod_collected_amount` | Float  | Opsional    | Jumlah uang COD yang diterima (0.0 jika sudah lunas/non-COD).   |

```json
{
  "delivery_status": "finished-delivery",
  "receiver_name": "Mpok Romlah (Penerima Langsung)",
  "cod_collected_amount": 110000.0
}
```

### Responses Body :

#### ✅ 200 OK

Status pengiriman berhasil diperbarui dan status pesanan pada tabel `orders` telah disinkronkan secara otomatis

```json
{
  "success": true,
  "message": "Delivery updated successfully",
  "data": {
    "id": 1,
    "order_id": 45,
    "delivery_status": "finished-delivery",
    "shipping_cost": 10000.0,
    "courier_id": 4,
    "courier_departed_at": "2026-01-19 16:10:00",
    "courier_arrived_at": "2026-01-19 16:20:00",
    "receiver_name": "Mpok Romlah (Penerima Langsung)",
    "cod_collected_amount": 110000.0,
    "created_at": "2026-01-19 16:00:00",
    "updated_at": "2026-01-19 16:20:00"
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi kesalahan validasi data atau pelanggaran aturan transisi status (misal: melompat langsung ke finished tanpa melalui tahap being-delivered).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "receiver_name": "receiver_name is required when status is finished-delivery"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Respons ketika token akses tidak valid, kedaluwarsa, atau tidak disertakan dalam header permintaan.

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

Respons ketika peran (role) pengguna tidak memiliki otoritas untuk memperbarui status pengiriman.

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

#### 🚫 404 Not Found

Respons ketika data pengiriman dengan ID yang diminta tidak ditemukan di dalam sistem.

```json
{
  "success": false,
  "message": "Delivery not found",
  "data": {
    "error_code": "RESOURCE_NOT_FOUND",
    "errors": null
  }
}
```

#### 🚫 409 Conflict

Terjadi konflik logika bisnis, seperti mencoba mengambil tugas yang sudah diklaim oleh kurir lain.

```json
{
  "success": false,
  "message": "Data conflict occurred",
  "data": {
    "error_code": "DATA_CONFLICT",
    "errors": {
      "courier_id": "Delivery task already taken by another courier"
    }
  }
}
```

#### 🚫 429 Too Many Requests

Respons ketika klien mengirimkan permintaan dalam jumlah yang melampaui batas (rate limit) dalam periode waktu tertentu.

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

Terjadi kesalahan internal pada server, seperti kegagalan transaksi database saat melakukan sinkronisasi data ke tabel orders.

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
