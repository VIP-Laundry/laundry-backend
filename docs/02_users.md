# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## USERS MODULE SPECIFICATION

---

## Endpoint : `POST /api/v1/users`

### Description :

Endpoint ini digunakan secara eksklusif oleh **Owner** untuk mendaftarkan akun karyawan baru (Kasir, Staff, atau Kurir) ke dalam sistem. Backend melakukan validasi ganda untuk memastikan `email` dan `username` belum pernah digunakan sebelumnya (_unique constraint_). Kata sandi akan diproses menggunakan algoritma _hashing_ aman sebelum disimpan ke database untuk menjamin privasi pengguna.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

Seluruh field di bawah ini wajib disertakan dalam Request Body.

| Key          | Type   | Location | Default | Description                                      |
| ------------ | ------ | -------- | ------- | ------------------------------------------------ |
| full_name    | String | Body     | -       | Nama lengkap user (Max. 150 karakter).           |
| username     | String | Body     | -       | Username unik, tanpa spasi (Max. 100 karakter).  |
| email        | String | Body     | -       | Email unik dan format valid (Max. 150 karakter). |
| password     | String | Body     | -       | Kata sandi minimal 8 karakter.                   |
| phone_number | String | Body     | -       | Nomor telepon aktif (Max. 30 karakter).          |
| role         | Enum   | Body     | -       | Pilihan: `owner`, `cashier`, `staff`, `courier`. |

```
{
  "full_name": "Siti Aminah",
  "username": "sitiaminah",
  "email": "sitiaminah@gmail.com",
  "password": "rahasia123",
  "phone_number": "082345678901",
  "role": "cashier"
}
```

### Request Body :

Objek JSON yang dikirimkan untuk proses registrasi karyawan baru.

```json
{
  "full_name": "Siti Aminah",
  "username": "sitiaminah",
  "email": "sitiaminah@gmail.com",
  "password": "rahasia123",
  "phone_number": "082345678901",
  "role": "cashier"
}
```

### Responses Body :

#### ✅ 201 Created

Karyawan baru berhasil didaftarkan. Akun secara otomatis diatur dalam status aktif (`is_active: 1`).

```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": 2,
    "full_name": "Siti Aminah",
    "username": "sitiaminah",
    "email": "sitiaminah@gmail.com",
    "role": "cashier",
    "phone_number": "082345678901",
    "is_active": true,
    "created_at": "2026-01-20 07:24:03",
    "updated_at": null
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi kegagalan validasi skema data (field kosong atau format tidak sesuai).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "email": "Invalid email format",
      "password": "Password must be at least 8 characters"
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

Pengguna memiliki token valid, namun tidak memiliki hak akses `owner` untuk melakukan pendaftaran.

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

Username atau Email yang dikirimkan sudah terdaftar di database.

```json
{
  "success": false,
  "message": "Data already exists",
  "data": {
    "error_code": "DUPLICATE_DATA",
    "errors": {
      "username": "Username 'sitiaminah' is already taken"
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

Kegagalan teknis pada server atau database.

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

## Endpoint : `GET /api/v1/users`

### Description :

Endpoint ini digunakan oleh **Owner** untuk mendapatkan daftar seluruh akun karyawan secara terorganisir. Menggunakan teknik **Pagination** untuk efisiensi bandwidth dan mendukung pencarian dinamis (Filtering). Data ditarik langsung dari tabel `users`.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Daftar filter pencarian dan pengaturan data melalui Query String.

| Key      | Type   | Location | Default    | Description                                              |
| -------- | ------ | -------- | ---------- | -------------------------------------------------------- |
| page     | Int    | Query    | 1          | Nomor halaman.                                           |
| per_page | Int    | Query    | 10         | Jumlah data per halaman.                                 |
| search   | String | Query    | -          | Cari berdasarkan nama atau username.                     |
| role     | Enum   | Query    | -          | Filter peran: owner, cashier, staff, courier.            |
| status   | Int    | Query    | -          | Filter status akun: 1 (Aktif/true), 0 (Non-aktif/false). |
| sort_by  | String | Query    | created_at | Kolom pengurutan (contoh: full_name, created_at).        |
| order    | String | Query    | desc       | Arah: asc (A-Z/Lama) atau desc (Z-A/Baru).               |

```
GET /api/v1/users?page=1&per_page=10&status=1&sort_by=full_name&order=asc
```

### Request Body :

```
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK

Daftar pengguna berhasil diambil secara sukses beserta informasi metadata halaman untuk kebutuhan navigasi di Frontend.

```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": 1,
      "full_name": "Farhan Rizki Maulana",
      "username": "farhanrizkimln",
      "role": "owner",
      "is_active": true
    },
    {
      "id": 2,
      "full_name": "Siti Aminah",
      "username": "sitiaminah",
      "role": "cashier",
      "is_active": true
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total_items": 48,
    "total_pages": 5
  }
}
```

#### ⚠️ 400 Bad Request

Parameter query tidak valid (format salah).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "page": "Must be a number",
      "limit": "Max limit is 100"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Token tidak valid atau tidak disertakan.

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

Akses ditolak karena pengguna bukan Owner.

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

Kegagalan sistem saat pengambilan data.

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

## Endpoint : `GET /api/v1/users/{id}`

### Description :

Endpoint ini digunakan oleh **Owner** untuk mendapatkan informasi profil lengkap dan mendalam dari seorang karyawan. Berbeda dengan endpoint list yang hanya memberikan informasi ringkas, endpoint detail ini menyajikan seluruh atribut pengguna (_kecuali hash password_) termasuk kontak, status akun, dan rekam jejak waktu (_timestamps_) untuk keperluan administrasi dan audit.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Identifikasi pengguna dilakukan melalui jalur URL menggunakan ID unik yang terdaftar di database.

| Key | Type | Location | Default | Description                                                  |
| --- | ---- | -------- | ------- | ------------------------------------------------------------ |
| id  | Int  | Path     | -       | ID Unik (Primary Key) karyawan yang ingin dilihat detailnya. |

```
GET /api/v1/users/1
```

### Request Body :

```
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK (Success)

Data pengguna ditemukan. Informasi disajikan dalam bentuk objek tunggal yang berisi profil lengkap.

```json
{
  "success": true,
  "message": "User detail retrieved successfully",
  "data": {
    "id": 1,
    "full_name": "Farhan Rizki Maulana",
    "username": "farhanrizkimln",
    "email": "farhanrizki@gmail.com",
    "role": "owner",
    "phone_number": "081234567890",
    "is_active": true,
    "last_login_at": "2025-12-28 05:12:36",
    "created_at": "2025-12-28 03:12:36",
    "updated_at": null
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi jika format ID yang dikirimkan tidak valid (bukan angka).

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

Akses ditolak karena peran pengguna saat ini tidak memiliki izin untuk melihat detail karyawan lain.

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

Permintaan ID valid secara format, namun data karyawan tersebut tidak ditemukan di database.

```json
{
  "success": false,
  "message": "User not found",
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
  "message": "Too many requests, please try again later.",
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
  "message": "An unexpected server error occurred",
  "data": {
    "error_code": "INTERNAL_SERVER_ERROR",
    "errors": null
  }
}
```

---

## Endpoint : `PUT /api/v1/users/{id}`

### Description :

Endpoint ini digunakan untuk memperbarui data profil karyawan. Sistem menerapkan **Multi-layered Authorization** di mana Owner memiliki kendali mutlak, sedangkan karyawan lain hanya diperbolehkan mengelola data pribadi mereka sendiri. Perubahan pada hak akses (_role_) dan status akun dikunci secara ketat agar tidak bisa dimanipulasi oleh entitas non-owner.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner, cashier, staff, courier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

Bagian ini menggunakan Path Parameter untuk menentukan target pengguna yang akan diperbarui.

| Key | Type | Location | Default | Description                                    |
| --- | ---- | -------- | ------- | ---------------------------------------------- |
| id  | Int  | Path     | -       | ID Unik karyawan yang akan diperbarui datanya. |

```
PUT /api/v1/users/1
```

### Logic Guard :

Sebagai **Architect**, Anda harus memastikan logika ini tertanam di dalam _middleware_ atau _service layer_:

| Skenario                   | Izin Akses    | Field yang Boleh Diubah            |
| -------------------------- | ------------- | ---------------------------------- |
| Owner akses ID siapa saja  | **DIIZINKAN** | Semua field tanpa pengecualian     |
| User akses ID diri sendiri | **DIIZINKAN** | Semua kecuali `role` & `is_active` |
| User akses ID orang lain   | **DITOLAK**   | Tidak ada (Respons 403).           |

**Catatan Teknis**: Jika non-owner mengirimkan field `role` atau `is_active`, sistem harus mengabaikan field tersebut dan tetap mempertahankan nilai lama di database tanpa memberikan error (Silent Ignore).

### Request Body :

Gunakan format JSON. Field yang tidak dikirimkan akan tetap menggunakan nilai yang sudah ada di database.

```json
{
  "full_name": "Farhan Rizki Maulana",
  "username": "farhanrizkimln",
  "email": "farhanrizki@gmail.com",
  "password": "barurahasia123",
  "phone_number": "081234567890",
  "role": "owner", // Opsional, hanya berlaku jika pengirim adalah Owner
  "is_active": 1 // Opsional, hanya berlaku jika pengirim adalah Owner
}
```

### Responses Body :

#### ✅ 200 OK

Data pengguna berhasil diperbarui secara sukses.

```json
{
  "success": true,
  "message": "User updated successfully",
  "data": {
    "id": 5,
    "full_name": "Farhan Rizki Maulana",
    "username": "farhanrizkimln",
    "email": "farhanrizki@gmail.com",
    "role": "owner",
    "phone_number": "081234567890",
    "is_active": 1,
    "last_login_at": "2026-01-12 08:12:36",
    "created_at": "2025-12-28 03:12:36",
    "updated_at": "2026-01-13 15:30:00"
  }
}
```

#### ⚠️ 400 Bad Request

Gagal validasi input atau format ID salah.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "password": "Password must be at least 8 characters"
    }
  }
}
```

#### 🚫 403 Forbidden

Terjadi ketika karyawan mencoba mengedit profil karyawan lain.

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

ID user yang dituju tidak ada di database.

```json
{
  "success": false,
  "message": "User not found",
  "data": {
    "error_code": "RESOURCE_NOT_FOUND",
    "errors": null
  }
}
```

#### 🚫 409 Conflict

Username atau Email baru yang dimasukkan sudah digunakan oleh orang lain.

```json
{
  "success": false,
  "message": "Data already exists",
  "data": {
    "error_code": "DUPLICATE_DATA",
    "errors": {
      "username": "Username already taken"
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

Kegagalan teknis pada server atau database.

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

## Endpoint : `DELETE /api/v1/users/{id}`

### Description :

Endpoint ini digunakan untuk menghapus pengguna secara logika (**Soft Delete**). Sistem tidak akan menghapus baris data dari database, melainkan mengubah nilai `is_active` menjadi `0`. Hal ini bertujuan untuk menjaga integritas data pada riwayat transaksi laundry yang pernah ditangani oleh pengguna tersebut.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Identifikasi target dilakukan melalui jalur URL menggunakan ID unik pengguna.

| Key | Type    | Location | Default | Description                                             |
| --- | ------- | -------- | ------- | ------------------------------------------------------- |
| id  | Integer | Path     | -       | ID Unik (Primary Key) karyawan yang akan dinonaktifkan. |

```
DELETE /api/v1/users/1
```

### 🛡️ Logic Guard (Aturan Keamanan) :

1. **Authorization**: Hanya pengguna dengan peran `owner` yang diizinkan memanggil endpoint ini.
2. **Anti Self-Deletion**: Sistem wajib menolak permintaan jika `id` pada _path_ sama dengan `user_id` yang sedang login (Owner tidak boleh menonaktifkan akunnya sendiri untuk mencegah sistem tanpa admin).

### Request Body :

```
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK

User berhasil dinonaktifkan. Kita tetap mengirimkan objek ID sebagai konfirmasi.

```json
{
  "success": true,
  "message": "User deleted successfully",
  "data": {
    "id": 1
  }
}
```

#### ⚠️ 400 Bad Request

Format ID pada URL tidak valid.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "id": "ID must be a valid number"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Token tidak valid, expired, atau tidak disertakan.

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

Terjadi jika pengirim bukan Owner, atau Owner mencoba menghapus ID-nya sendiri.

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

ID user tidak ditemukan di database.

```json
{
  "success": false,
  "message": "User not found",
  "data": {
    "error_code": "RESOURCE_NOT_FOUND",
    "errors": null
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

Kesalahan teknis saat proses update di database.

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
