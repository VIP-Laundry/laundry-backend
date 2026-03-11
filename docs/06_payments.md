# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## PAYMENTS MODULE SPECIFICATION

---

## Endpoint : `GET /payments`

### Description :

Endpoint ini digunakan untuk mengambil daftar seluruh transaksi pembayaran (tagihan dan pelunasan). Memberikan visibilitas penuh terhadap piutang yang masih menggantung (`pending`) maupun uang yang sudah dikonfirmasi masuk ke kas (`confirmed`).

### Role Based Access Control (RBAC) :

- `Permissions`: `owner, cashier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Bagian ini mendefinisikan parameter query opsional untuk memfilter hasil pembayaran.

| Key      | Type   | Location | Default    | Description                                                     |
| -------- | ------ | -------- | ---------- | --------------------------------------------------------------- |
| page     | Int    | Query    | 1          | Nomor halaman (Pagination)                                      |
| per_page | Int    | Query    | 10         | Jumlah data per halaman                                         |
| search   | String | Query    | -          | Cari berdasarkan Order ID atau Nomor Referensi                  |
| status   | String | Query    | -          | Filter berdasarkan status pembayaran (pending, confirmed, void) |
| method   | String | Query    | -          | Filter berdasarkan metode pembayaran (cash, transfer, etc)      |
| sort_by  | String | Query    | created_at | Pengurutan (contoh: amount, created_at).                        |
| order    | String | Query    | desc       | Arah urutan: asc atau desc.                                     |

```
GET /api/payments?page=1&per_page=10&status=pending&sort_by=created_at&order=desc
```

### Request Body :

```
None (Kosong, karena method GET tidak memerlukan body request).
```

### Responses Body :

#### ✅ 200 OK

Bagian ini berisi contoh respons sukses ketika daftar pembayaran berhasil diambil.

```json
{
  "success": true,
  "message": "Payments retrieved successfully",
  "data": [
    {
      "id": 1,
      "order_id": 45,
      "method": null,
      "amount": 110000.0,
      // "amount_received": 0.0,
      // "amount_change": 0.0,
      "reference_no": null,
      "status": "pending"
      // "created_by": 2,
      // "collected_by": null,
      // "collected_at": null,
      // "created_at": "2026-01-05 13:00:00"
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

Bagian ini berisi contoh respons ketika terjadi kesalahan validasi pada parameter query yang dikirimkan.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "page": "The page must be a positive integer.",
      "per_page": "The per_page must be between 1 and 100."
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Bagian ini berisi contoh respons ketika token otentikasi tidak valid atau tidak disertakan dalam header permintaan.

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

Bagian ini berisi contoh respons ketika peran pengguna tidak memiliki izin untuk mengakses endpoint ini.

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

Bagian ini berisi contoh respons ketika terjadi kelebihan permintaan dari client.

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

Bagian ini berisi contoh respons ketika terjadi kesalahan server saat memproses permintaan.

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

## Endpoint : `GET /payments/{id}`

### Description :

Endpoint ini digunakan untuk mengambil informasi mendalam mengenai satu transaksi pembayaran berdasarkan ID uniknya. Sesuai prinsip Pemisahan Modul, respons ini hanya mengembalikan data yang tersimpan di tabel `payments`. Untuk melihat detail pesanan terkait, klien harus melakukan permintaan terpisah ke modul Orders menggunakan `order_id` yang tersedia di sini.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner, cashier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Bagian ini menggunakan _Path Parameter_ untuk menentukan data pembayaran spesifik yang ingin diakses.

| Key | Type | Location | Default | Description                                          |
| --- | ---- | -------- | ------- | ---------------------------------------------------- |
| id  | Int  | Path     | -       | ID unik (Primary Key) pembayaran yang ingin diambil. |

```
GET /api/payments/1
```

### Request Body :

```
None (Kosong, karena method GET tidak memerlukan body request).
```

### Responses Body :

#### ✅ 200 OK

Bagian ini berisi contoh respons sukses ketika detail pembayaran berhasil diambil.

```json
{
  "success": true,
  "message": "Payment retrieved successfully",
  "data": {
    "id": 1,
    "order_id": 45,
    "method": null,
    "amount": 110000.0,
    "amount_received": 0.0,
    "amount_change": 0.0,
    "reference_no": null,
    "status": "pending",
    "created_by": 1,
    "collected_by": null,
    "collected_at": null,
    "created_at": "2026-01-05 13:00:00"
  }
}
```

#### ⚠️ 400 Bad Request

Bagian ini berisi contoh respons ketika terjadi kesalahan validasi pada parameter path yang dikirimkan.

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

Bagian ini berisi contoh respons ketika token otentikasi tidak valid atau tidak disertakan dalam header permintaan.

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

Bagian ini berisi contoh respons ketika peran pengguna tidak memiliki izin untuk mengakses endpoint ini.

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

Bagian ini berisi contoh respons ketika pembayaran dengan ID yang diminta tidak ditemukan.

```json
{
  "success": false,
  "message": "Payment not found",
  "data": {
    "error_code": "RESOURCE_NOT_FOUND",
    "errors": null
  }
}
```

#### 🚫 429 Too Many Requests

Bagian ini berisi contoh respons ketika terjadi kelebihan permintaan dari client.

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

Bagian ini berisi contoh respons ketika terjadi kesalahan server saat memproses permintaan.

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

## Endpoint : `PATCH /payments/{id}`

### Description :

Endpoint ini digunakan untuk memproses pelunasan transaksi (Settlement). Kasir menginput nominal uang yang diterima dan metode pembayaran. Backend akan memvalidasi jumlah uang, menghitung kembalian, dan mencatat waktu pelunasan secara otomatis.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner, cashier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

Bagian ini mendefinisikan parameter path yang wajib disertakan dalam permintaan.

| Key | Type | Location | Default | Description                            |
| --- | ---- | -------- | ------- | -------------------------------------- |
| id  | Int  | Path     | -       | ID unik pembayaran yang akan dilunasi. |

```
PATCH /api/payments/1
```

### Request Body :

Bagian ini berisi contoh objek JSON yang dikirimkan dalam Request Body untuk memperbarui informasi pembayaran.

```json
{
  "method": "cash",
  "amount_receive": 150000.0,
  // "reference_no": null,
  "status": "confirmed"
}
```

### Responses Body :

#### ✅ 200 OK

Bagian ini berisi contoh respons sukses ketika informasi pembayaran berhasil diperbarui.

```json
{
  "success": true,
  "message": "Payment updated successfully",
  "data": {
    "id": 1,
    "order_id": 45,
    "method": "cash",
    "amount": 110000.0,
    "amount_received": 150000.0,
    "amount_change": 40000.0,
    "reference_no": null,
    "status": "confirmed",
    "created_by": 1,
    "collected_by": 2,
    "collected_at": "2026-01-19 14:00:00",
    "created_at": "2026-01-19 13:00:00"
  }
}
```

#### ⚠️ 400 Bad Request

Bagian ini berisi contoh respons ketika terjadi kesalahan validasi pada data yang dikirimkan dalam Request Body.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "amount_received": "The amount_received must be greater than or equal to amount."
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Bagian ini berisi contoh respons ketika token otentikasi tidak valid atau tidak disertakan dalam header permintaan.

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

Bagian ini berisi contoh respons ketika peran pengguna tidak memiliki izin untuk mengakses endpoint ini.

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

Bagian ini berisi contoh respons ketika pembayaran dengan ID yang diminta tidak ditemukan.

```json
{
  "success": false,
  "message": "Payment not found",
  "data": {
    "error_code": "RESOURCE_NOT_FOUND",
    "errors": null
  }
}
```

#### 🚫 409 Conflict

Bagian ini berisi contoh respons ketika terjadi konflik data, misalnya nomor referensi transaksi sudah digunakan.

```json
{
  "success": false,
  "message": "Data conflict occurred",
  "data": {
    "error_code": "DATA_CONFLICT",
    "errors": {
      "reference_no": "Reference number already in use"
    }
  }
}
```

#### 🚫 429 Too Many Requests

Bagian ini berisi contoh respons ketika terjadi kelebihan permintaan dari client.

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

Bagian ini berisi contoh respons ketika terjadi kesalahan server saat memproses permintaan.

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
