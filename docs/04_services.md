# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## SERVICES MODULE SPECIFICATION

---

## Endpoint : `POST /services`

### Description :

Endpoint ini digunakan oleh **Owner** untuk mendaftarkan item layanan laundry baru. Setiap layanan wajib dikaitkan dengan `category_id` yang valid untuk menentukan klasifikasinya (**Kiloan/Satuan**). Sistem juga mewajibkan adanya kode unik (`code`) sebagai identitas SKU yang mempermudah pencarian dan integrasi dengan sistem inventaris atau barcode di masa depan.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

Parameter di bawah ini dikirimkan melalui Request Body untuk membentuk entitas layanan baru.

| Key            | Tipe   | Location | Default | Description                                       |
| -------------- | ------ | -------- | ------- | ------------------------------------------------- |
| category_id    | Int    | Body     |         | ID Kategori yang valid (Foreign Key).             |
| code           | String | Body     |         | Kode unik layanan (Contoh: `SVC-LKR-1`).          |
| service_name   | String | Body     |         | Nama layanan (Contoh: `Kiloan Regular`).          |
| unit           | Enum   | Body     |         | Satuan layanan (`kg` atau `pcs`).                 |
| price          | Int    | Body     |         | Harga per unit (Contoh: 7000)                     |
| duration_hours | Int    | Body     |         | Estimasi waktu selesai dalam jam. (Contoh: `72`). |

```
{
  "category_id": 1,
  "code": "SVC-LKR-1",
  "service_name": "Layanan Cuci Kiloan Regular",
  "unit": "kg",
  "price": 7000,
  "duration_hours": 72
}
```

### Request Body :

Objek JSON berisi detail spesifikasi layanan yang akan disimpan.

```json
{
  "category_id": 1,
  "code": "SVC-LKR-1",
  "service_name": "Kiloan Regular",
  "unit": "kg",
  "price": 7000,
  "duration_hours": 72
}
```

### Responses Body :

#### ✅ 201 Created

Layanan baru berhasil didaftarkan. Status default adalah aktif (`is_active: 1`).

```json
{
  "success": true,
  "message": "Service created successfully",
  "data": {
    "id": 1,
    "category_id": 1,
    "code": "SVC-LKR-1",
    "service_name": "Kiloan Regular",
    "unit": "kg",
    "price": 7000,
    "duration_hours": 72,
    "is_active": 1,
    "created_at": "2026-01-20 18:00:00",
    "updated_at": null
  }
}
```

#### ⚠️ 400 Bad Request

Input tidak valid (misal: `price` negatif, durasi 0, atau `unit` bukan kg/pcs).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "price": "Price cannot be negative",
      "duration_hours": "Duration must be at least 1 hour"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Respons ketika token akses tidak valid atau tidak disertakan.

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

Akses ditolak karena role Anda bukan **Owner**.

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

Terjadi jika `category_id` yang dikirimkan tidak terdaftar di database.

```json
{
  "success": false,
  "message": "Category not found",
  "data": {
    "error_code": "RESOURCE_NOT_FOUND",
    "errors": null
  }
}
```

#### 🚫 409 Conflict

Kode layanan (`code`) atau nama layanan sudah digunakan oleh layanan lain.

```json
{
  "success": false,
  "message": "Data already exists",
  "data": {
    "error_code": "DUPLICATE_DATA",
    "errors": {
      "code": "Service code 'SVC-LKR-1' already exists"
    }
  }
}
```

#### 🚫 429 Too Many Requests

Terlalu banyak permintaan pembuatan layanan dalam waktu singkat.

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

Kegagalan sistem atau database saat memproses penyimpanan data layanan.

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

## Endpoint : `GET /services`

#### Description :

Endpoint ini digunakan untuk mengambil daftar seluruh layanan laundry dalam format **Brief View**. Sistem menggunakan teknik **Pagination** untuk efisiensi beban data dan mendukung **Filtering** serta **Sorting** yang mendalam. Data dikirimkan secara **Nested** (menyertakan informasi dasar kategori) untuk mempermudah identifikasi layanan di sisi UI tanpa beban tambahan pada database.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner, cashier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Pengaturan pengambilan data dilakukan melalui Query String.

| Key         | Type   | Location | Default      | Description                                     |
| ----------- | ------ | -------- | ------------ | ----------------------------------------------- |
| page        | Int    | Query    | 1            | Nomor halaman data.                             |
| per_page    | Int    | Query    | 10           | Jumlah data per halaman (Maks. 100).            |
| category_id | Int    | Query    | -            | Filter berdasarkan ID Kategori.                 |
| search      | String | Query    | -            | Cari berdasarkan Kode SVC atau Nama Layanan.    |
| status      | Int    | Query    | 1            | Filter status: 1 (Aktif), 0 (Non-aktif/Arsip).  |
| sort_by     | String | Query    | service_name | Kolom pengurutan (contoh: price, service_name). |
| order       | String | Query    | asc          | Arah urutan: asc atau desc.                     |

```
GET /api/services?page=1&per_page=10&status=1&category_id=1&sort_by=price&order=asc
```

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK

Daftar layanan berhasil ditarik dengan informasi kategori ringkas (_Optimized Nested Payload_).

```json
{
  "success": true,
  "message": "Services retrieved successfully",
  "data": [
    {
      "id": 1,
      "code": "SVC-LKR-1",
      "service_name": "Layanan Cuci Kiloan Regular",
      "unit": "kg",
      "price": 7000,
      "duration_hours": 72,
      "is_active": 1,
      "category": {
        "id": 1,
        "category_name": "Layanan Kiloan"
      }
    },
    {
      "id": 2,
      "code": "SVC-LJSR-1",
      "service_name": "Layanan Cuci Jas Satuan Regular",
      "unit": "pcs",
      "price": 10000,
      "duration_hours": 72,
      "is_active": 1,
      "category": {
        "id": 2,
        "category_name": "Layanan Satuan"
      }
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total_items": 12,
    "total_pages": 2
  }
}
```

#### ⚠️ 400 Bad Request

Parameter query tidak valid (format angka salah atau melebihi batas).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "per_page": "Max per_page is 100",
      "category_id": "Category ID must be a number"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Token tidak valid, kadaluarsa, atau tidak disertakan.

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

Akses ditolak karena peran pengguna tidak memiliki izin (**Kurir/Staff**).

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

Terlalu banyak permintaan dalam waktu singkat (Rate Limit).

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

Kegagalan teknis pada server atau database saat pengambilan data.

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

## Endpoint : `GET /services/{id}`

#### Description :

Endpoint ini digunakan untuk mengambil profil lengkap dari satu layanan spesifik. Berbeda dengan endpoint _List_, di sini sistem menyajikan data secara **Full Detail**, termasuk deskripsi lengkap kategori yang menaunginya serta rekam jejak waktu (_timestamps_). Data ini menjadi referensi utama bagi **Owner** saat ingin melakukan validasi data sebelum proses pembaruan atau audit internal.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner, cashier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Parameter ini disisipkan langsung di URL (contoh: /services/1).

| Key | Type | Location | Default | Description                                            |
| --- | ---- | -------- | ------- | ------------------------------------------------------ |
| id  | Int  | Path     | -       | ID Unik (Primary Key) dari layanan yang ingin diakses. |

```
GET /api/services/1
```

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK

Detail layanan berhasil ditarik beserta informasi lengkap kategori terkait.

```json
{
  "success": true,
  "message": "Service detail retrieved successfully",
  "data": {
    "id": 1,
    "code": "SVC-LKR-1",
    "service_name": "Layanan Cuci Kiloan Regular",
    "unit": "kg",
    "price": 7000,
    "duration_hours": 72,
    "is_active": 1,
    "created_at": "2026-01-20 07:24:03",
    "updated_at": null,
    "category": {
      "id": 1,
      "category_name": "Layanan Kiloan",
      "description": "Cuci pakaian sehari-hari dengan sistem kiloan"
    }
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi jika format ID pada URL bukan merupakan angka yang valid.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "id": "id must be a valid number"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Respons ketika token akses tidak valid, kedaluwarsa, atau tidak disertakan dalam header.

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

Akses ditolak karena peran pengguna (`Staff/Courier`) tidak memiliki otoritas untuk melihat master data layanan.

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

ID layanan valid secara format (angka), namun data tersebut tidak ditemukan di database.

```json
{
  "success": false,
  "message": "Service not found",
  "data": {
    "error_code": "RESOURCE_NOT_FOUND",
    "errors": null
  }
}
```

#### 🚫 429 Too Many Requests

Terlalu banyak permintaan akses detail dalam waktu singkat (Rate Limiting).

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

Kegagalan teknis pada server atau database saat proses pencarian data layanan.

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

## Endpoint : `PUT /services/{id}`

#### Description :

Endpoint ini digunakan secara eksklusif oleh **Owner** untuk memperbarui data layanan laundry yang sudah terdaftar. Sistem menerapkan logika **Partial Update** (Pembaruan Sebagian), di mana field yang tidak disertakan dalam _Request Body_ akan tetap menggunakan nilai lama di database. Jika terdapat perubahan pada `code` (SKU), sistem akan memvalidasi keunikannya, serta memastikan `category_id` tujuan benar-benar tersedia di sistem jika terjadi perpindahan kategori.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

bagian ini mendefinisikan parameter path yang diperlukan untuk mengidentifikasi layanan yang akan diperbarui.

| Key | Type | Location | Default | Description                                               |
| --- | ---- | -------- | ------- | --------------------------------------------------------- |
| id  | Int  | Path     | -       | ID Unik (Primary Key) dari layanan yang ingin diperbarui. |

```
PUT /api/services/1
```

### Request Body :

Kirimkan objek JSON berisi field yang ingin diubah. Field bersifat opsional untuk mendukung pembaruan sebagian.

```json
{
  "category_id": 1,
  "code": "SVC-LKR-1",
  "service_name": "Kiloan Regular",
  "unit": "kg",
  "price": 7500,
  "duration_hours": 72,
  "is_active": 1
}
```

### Responses Body :

#### ✅ 200 OK

Data layanan berhasil diperbarui secara sukses. Objek `data` mengembalikan profil terbaru beserta timestamp `updated_at`.

```json
{
  "success": true,
  "message": "Service updated successfully",
  "data": {
    "id": 1,
    "category_id": 1,
    "code": "SVC-LKR-1",
    "service_name": "Kiloan Regular",
    "unit": "kg",
    "price": 7500,
    "duration_hours": 72,
    "is_active": 1,
    "created_at": "2025-12-28 07:24:03",
    "updated_at": "2026-01-20 23:24:00"
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi jika input tidak memenuhi validasi skema atau format ID pada URL salah.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "price": "Price cannot be negative",
      "duration_hours": "Duration must be greater than 0"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Respons ketika token akses tidak valid, kedaluwarsa, atau tidak disertakan dalam header.

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

Akses ditolak karena peran pengguna saat ini bukan `owner`.

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

ID layanan tidak ditemukan, atau `category_id` baru yang dikirimkan tidak terdaftar di database.

```json
{
  "success": false,
  "message": "Service not found",
  "data": {
    "error_code": "RESOURCE_NOT_FOUND",
    "errors": null
  }
}
```

#### 🚫 409 Conflict

Terjadi jika pembaruan `code` layanan baru sudah digunakan oleh layanan lain.

```json
{
  "success": false,
  "message": "Data already exists",
  "data": {
    "error_code": "DUPLICATE_DATA",
    "errors": {
      "code": "Service code 'SVC-LKR-1' already exists"
    }
  }
}
```

#### 🚫 429 Too Many Requests

Terlalu banyak permintaan pembaruan dalam waktu singkat (Rate Limiting).

```json
{
  "success": false,
  "message": "Too many requests, please try again later.",
  "data": {
    "error_code": "RATE_LIMIT_EXCEEDED",
    "errors": null
  }
}
```

#### 🔥 500 Internal Server Error

Kegagalan teknis pada server atau database saat memproses pembaruan data layanan.

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

## Endpoint : `DELETE /services/{id}`

### Description :

Endpoint ini digunakan secara eksklusif oleh **Owner** untuk menonaktifkan layanan laundry secara logika (**Soft Delete**). Sistem tidak akan menghapus data dari baris database, melainkan mengubah nilai status `is_active` menjadi `0`. Hal ini menjamin bahwa seluruh riwayat transaksi di masa lalu tetap akurat, namun layanan tersebut tidak akan tersedia lagi untuk dipilih saat pembuatan pesanan baru.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Identifikasi layanan yang akan dinonaktifkan melalui Path Parameter.

| Key | Type | Location | Default | Description                                                  |
| --- | ---- | -------- | ------- | ------------------------------------------------------------ |
| id  | Int  | Path     | -       | ID Unik (Primary Key) dari layanan yang ingin dinonaktifkan. |

```
DELETE /api/services/1
```

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK

Layanan berhasil dinonaktifkan secara sukses dari sistem operasional.

```json
{
  "success": true,
  "message": "Service deleted successfully",
  "data": {
    "id": 1
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi jika format ID pada URL bukan merupakan angka yang valid.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "id": "id must be a valid number"
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

Akses ditolak karena peran pengguna saat ini bukan `owner`.

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

ID layanan valid secara format (angka), namun data tersebut tidak ditemukan di database.

```json
{
  "success": false,
  "message": "Service not found",
  "data": {
    "error_code": "RESOURCE_NOT_FOUND",
    "errors": null
  }
}
```

#### 🚫 429 Too Many Requests

Terlalu banyak permintaan penghapusan dalam waktu singkat (Rate Limiting).

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

Kegagalan teknis pada server atau database saat proses pembaruan status layanan.

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
