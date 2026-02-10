# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## SERVICE CATEGORIES MODULE SPECIFICATION

---

## Endpoint : `POST /api/v1/categories`

### Description :

Endpoint ini digunakan secara eksklusif oleh **Owner** untuk mendefinisikan klasifikasi layanan baru dalam sistem (`Contoh: "Kiloan", "Satuan", "Dry Clean"`). Backend akan melakukan pengecekan keunikan nama kategori guna mencegah duplikasi data yang dapat membingungkan pelanggan dan kasir saat proses pembuatan pesanan.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

Parameter di bawah ini dikirimkan melalui Request Body untuk membentuk entitas kategori baru.

| Key           | Type   | Location | Default | Description                                            |
| ------------- | ------ | -------- | ------- | ------------------------------------------------------ |
| category_name | String | Body     | -       | Nama unik kategori layanan (Contoh: "Layanan Kiloan"). |
| description   | String | Body     | -       | Deskripsi singkat tentang kategori layanan (Opsional). |

```
{
  "category_name": "Layanan Kiloan",
  "description": "Cuci pakaian sehari-hari dihitung per kilogram"
}
```

### Request Body :

Objek JSON berisi identitas kategori yang akan disimpan.

```json
{
  "category_name": "Layanan Kiloan",
  "description": "Cuci pakaian sehari-hari dihitung per kilogram"
}
```

### Responses Body :

#### ✅ 201 Created

Kategori layanan berhasil didaftarkan. Secara default, kategori baru akan berstatus aktif (`is_active: 1`).

```json
{
  "success": true,
  "message": "Category created successfully",
  "data": {
    "id": 1,
    "category_name": "Layanan Kiloan",
    "description": "Cuci pakaian sehari-hari dihitung per kilogram",
    "is_active": 1,
    "created_at": "2026-01-20 16:00:00",
    "updated_at": null
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi jika input tidak memenuhi validasi skema (misal: `category_name` kosong atau terlalu panjang).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "category_name": "Category name is required"
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

Akses ditolak karena pengguna bukan Role **Owner**.

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

#### 🚫 409 Conflict

Nama kategori yang dikirimkan sudah terdaftar sebelumnya di database.

```json
{
  "success": false,
  "message": "Data already exists",
  "data": {
    "error_code": "DUPLICATE_DATA",
    "errors": {
      "category_name": "Category 'Layanan Kiloan' already exists"
    }
  }
}
```

#### 🚫 429 Too Many Requests

Terjadi jika terlalu banyak permintaan yang dikirim dalam waktu singkat, memicu mekanisme rate limiting.

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

Kegagalan teknis pada server atau koneksi database saat memproses data.

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

## Endpoint : `GET /api/v1/categories`

### Description :

Endpoint ini digunakan untuk mengambil daftar kategori layanan dalam format **Brief Data**. Sistem menggunakan teknik **Pagination** untuk efisiensi beban data dan mendukung Filtering serta **Sorting** untuk memudahkan pencarian kategori tertentu saat manajemen layanan atau pembuatan transaksi.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`, `cashier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Bagian ini mendefinisikan filter pencarian melalui Query String.

| Key      | Type   | Location | Default       | Description                                           |
| -------- | ------ | -------- | ------------- | ----------------------------------------------------- |
| page     | Int    | Query    | 1             | Nomor halaman data yang ingin diambil.                |
| per_page | Int    | Query    | 10            | Jumlah data per halaman (Maks. 100).                  |
| search   | String | Query    | -             | Cari berdasarkan nama kategori.                       |
| status   | Int    | Query    | -             | Filter status: 1 (Aktif), 0 (Non-aktif).              |
| sort_by  | String | Query    | category_name | Kolom pengurutan (contoh: category_name, created_at). |
| order    | String | Query    | asc           | Arah urutan: asc (A-Z) atau desc (Z-A).               |

```
GET /api/categories?page=1&per_page=10&status=1&sort_by=category_name&order=asc
```

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK

Daftar kategori berhasil diambil dalam format ringkas (_Optimized Payload_).

```json
{
  "success": true,
  "message": "Categories retrieved successfully",
  "data": [
    {
      "id": 1,
      "category_name": "Layanan Kiloan",
      "is_active": 1
    },
    {
      "id": 2,
      "category_name": "Layanan Satuan",
      "is_active": 1
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

Terjadi jika parameter query melanggar validasi (misal: `per_page` di atas 100).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "page": "must be a number",
      "per_page": "max per_page is 100"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Respons ketika token akses tidak valid, kedaluwarsa, atau tidak disertakan.

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

Akses ditolak karena peran pengguna tidak memiliki izin (misal: Role `staff` atau `courier`).

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

Terlalu banyak permintaan pengambilan data dalam waktu singkat.

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

## Endpoint : `GET /api/v1/categories/{id}`

#### Description :

Endpoint ini digunakan untuk mengambil informasi mendalam dari satu kategori layanan tertentu. Berbeda dengan endpoint _List_, di sini sistem menyajikan seluruh atribut kategori termasuk deskripsi lengkap dan rekam jejak waktu (_timestamps_). Data ini biasanya diperlukan oleh Frontend untuk menampilkan detail di UI atau mengisi data awal pada formulir pembaruan (_Update Form_).

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`, `cashier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Identifikasi kategori dilakukan secara spesifik melalui jalur URL menggunakan ID unik.

| Key | Type | Location | Default | Description                                                     |
| --- | ---- | -------- | ------- | --------------------------------------------------------------- |
| id  | Int  | Path     | -       | ID Unik (Primary Key) dari kategori layanan yang ingin diakses. |

```
GET /api/categories/1
```

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK (Success)

Data kategori ditemukan. Informasi disajikan dalam bentuk objek tunggal yang berisi profil kategori secara lengkap.

```json
{
  "success": true,
  "message": "Category detail retrieved successfully",
  "data": {
    "id": 1,
    "category_name": "Layanan Kiloan",
    "description": "Cuci pakaian sehari-hari dihitung per kilogram",
    "is_active": 1,
    "created_at": "2026-01-20 10:00:00",
    "updated_at": null
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

Akses ditolak karena peran pengguna (`Staff/Courier`) tidak memiliki otoritas untuk melihat master data kategori.

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

ID kategori valid secara format (angka), namun data tersebut tidak ada di database.

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

#### 🚫 429 Too Many Requests

Terjadi jika terlalu banyak permintaan dalam waktu singkat, memicu mekanisme rate limiting.

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

Kegagalan teknis pada server atau database saat melakukan pencarian data.

```json
{
  "success": false,
  "message": "An unexpected error occurred",
  "data": {
    "error_code": "INTERNAL_SERVER_ERROR",
    "errors": null
  }
}
```

---

## Endpoint : `PUT /api/v1/categories/{id}`

### Description :

Endpoint ini digunakan secara eksklusif oleh **Owner** untuk memperbarui informasi pada kategori layanan yang sudah ada. Sistem mendukung pembaruan parsial (_Partial Update_), di mana field yang tidak disertakan dalam _Request Body_ akan tetap mempertahankan nilai lamanya di database. Jika terdapat perubahan pada `category_name`, sistem akan melakukan validasi keunikan untuk mencegah duplikasi.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

Bagian ini mendefinisikan ID kategori yang ingin diubah, disisipkan langsung di URL.

| Key | Type | Location | Default | Description                                                   |
| --- | ---- | -------- | ------- | ------------------------------------------------------------- |
| id  | Int  | Path     | -       | ID Unik (Primary Key) kategori yang ingin diperbarui datanya. |

```
PUT /api/categories/1
```

### Request Body :

Kirimkan objek JSON berisi field yang ingin diubah. Field bersifat opsional untuk mendukung pembaruan sebagian.

```json
{
  "category_name": "Layanan Satuan",
  "description": "Kategori untuk layanan cuci per item",
  "is_active": 1
}
```

### Responses Body :

#### ✅ 200 OK

Data kategori berhasil diperbarui. Objek data mengembalikan profil terbaru beserta timestamp updated_at.

```json
{
  "success": true,
  "message": "Category updated successfully",
  "data": {
    "id": 1,
    "category_name": "Layanan Satuan",
    "description": "Kategori untuk layanan cuci per item",
    "is_active": 1,
    "created_at": "2025-12-28 07:24:03",
    "updated_at": "2026-01-20 18:15:00"
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi jika format ID pada URL salah atau input melanggar validasi skema.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "category_name": "category name is required"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Respons ketika token akses tidak valid, kedaluwarsa, atau tidak disertakan.

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

Akses ditolak karena pengguna bukan Role **Owner**.

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

ID kategori tidak ditemukan di database.

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

Terjadi jika `category_name` yang baru sudah digunakan oleh kategori lain.

```json
{
  "success": false,
  "message": "Data already exists",
  "data": {
    "error_code": "DUPLICATE_DATA",
    "errors": {
      "category_name": "Category name 'Layanan Satuan' already taken"
    }
  }
}
```

#### 🚫 429 Too Many Requests

Terlalu banyak permintaan dalam waktu singkat.

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

Kegagalan teknis pada server atau database saat memproses pembaruan data.

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

## Endpoint : `DELETE /api/v1/categories/{id}`

### Description :

Endpoint ini digunakan oleh **Owner** untuk menonaktifkan kategori layanan secara logika (**Soft Delete**). Sistem tidak akan menghapus baris data dari database, melainkan mengubah status `is_active` menjadi `0`. Hal ini menjamin bahwa seluruh pesanan (_orders_) yang pernah tercatat menggunakan kategori ini tetap memiliki referensi data yang valid untuk kebutuhan laporan keuangan, namun kategori tersebut tidak akan muncul lagi sebagai pilihan saat pembuatan transaksi baru.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Identitas kategori yang akan dinonaktifkan ditentukan melalui jalur URL menggunakan ID unik.

| Key | Type | Location | Default | Description                                              |
| --- | ---- | -------- | ------- | -------------------------------------------------------- |
| id  | Int  | Path     | -       | ID unik (Primary Key) kategori yang ingin dinonaktifkan. |

```
DELETE /api/categories/1
```

### Request Body :

```
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK

Kategori berhasil dinonaktifkan secara sukses.

```json
{
  "success": true,
  "message": "Category deleted successfully",
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

Akses ditolak karena Anda bukan **Owner**.

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

Data kategori tidak ditemukan di database.

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

Kegagalan teknis pada server atau database saat memproses pembaruan status kategori.

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
