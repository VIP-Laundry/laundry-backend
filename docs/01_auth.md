# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## AUTH MODULE SPECIFICATION

---

## Endpoint : `POST /api/v1/auth/login`

### Description :

Endpoint ini digunakan untuk memverifikasi identitas pengguna (Otentikasi). Sistem akan melakukan pencocokan `username` pada tabel `users` dan memverifikasi hash password menggunakan algoritma aman. Jika valid, server akan menghasilkan **Access Token (JWT)** yang berisi _claim_ identitas dan peran pengguna untuk mengakses _endpoint_ lainnya.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner, cashier, staff, courier`

### Headers :

- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

Bagian ini mendefinisikan field yang wajib dikirimkan dalam Request Body untuk proses otentikasi. Password dikirim dalam bentuk teks biasa (plain text) melalui koneksi aman (HTTPS).

| Key      | Type   | Location | Description                                  |
| -------- | ------ | -------- | -------------------------------------------- |
| username | String | Body     | Username unik pengguna yang sudah terdaftar. |
| password | String | Body     | Kata sandi pengguna (Min. 8 karakter).       |

```
{
  "username": "farhanrizkimln",
  "password": "rahasia123"
}
```

### Request Body :

Objek JSON yang berisi kredensial pengguna. Password dikirim dalam bentuk teks biasa (plain text) melalui koneksi aman (HTTPS).

```json
{
  "username": "farhanrizkimln",
  "password": "rahasia123"
}
```

### Responses Body :

#### ✅ 200 OK

Otentikasi berhasil. Mengembalikan token akses dan profil singkat pengguna.

```json
{
  "success": true,
  "message": "Login successfully",
  "data": {
    "token": {
      "token_type": "Bearer",
      "access_token": "eyJhbGciOiJIUzI1NiIsInR...",
      "refresh_token": "RT-XYZ-999...",
      "expires_in": 900
    },
    "user": {
      "id": 1,
      "username": "farhanrizkimln",
      "role": "owner"
    }
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi jika format input tidak valid atau ada field yang kosong.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "username": "Username is required",
      "password": "Password must be at least 8 characters"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Kombinasi username dan password tidak cocok. Sistem memberikan pesan umum demi alasan keamanan (_Security Best Practice_).

```json
{
  "success": false,
  "message": "Invalid username or password",
  "data": {
    "error_code": "INVALID_CREDENTIALS",
    "errors": null
  }
}
```

#### 🚫 403 Forbidden

Kredensial benar, tetapi akun pengguna dalam status dinonaktifkan (`is_active = 0`).

```json
{
  "success": false,
  "message": "Your account is inactive. Please contact the administrator.",
  "data": {
    "error_code": "ACCOUNT_INACTIVE",
    "errors": null
  }
}
```

#### 🚫 429 Too Many Requests

Terlalu banyak percobaan login dalam waktu singkat, memicu mekanisme rate limiting.

```json
{
  "success": false,
  "message": "Too many login attempts, please try again later.",
  "data": {
    "error_code": "RATE_LIMIT_EXCEEDED",
    "errors": null
  }
}
```

#### 🔥 500 Internal Server Error

Kegagalan teknis pada server, seperti database timeout atau gagal melakukan signing pada JWT.

```json
{
  "success": false,
  "message": "An unexpected server error occurred during authentication",
  "data": {
    "error_code": "INTERNAL_SERVER_ERROR",
    "errors": null
  }
}
```

---

## Endpoint : `POST /api/v1/auth/refresh-token`

## Endpoint : `POST /api/v1/auth/logout`

### Description :

Endpoint ini berfungsi sebagai sinyal formal pengakhiran akses pengguna. Mengingat sistem menggunakan **Stateless JWT**, server akan memvalidasi keabsahan token yang dikirimkan. Setelah mendapatkan respons sukses, klien (_Frontend/Mobile_) wajib menghapus Access Token dari penyimpanan lokal (_Local Storage/Secure Storage_) untuk mengamankan akun.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner, cashier, staff, courier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

Bagian ini tidak memerlukan parameter pada _Path_ maupun _Query_ karena identitas akses sepenuhnya diambil dari _Authorization Header_.

| Key  | Type | Location | Default | Description                                  |
| ---- | ---- | -------- | ------- | -------------------------------------------- |
| None | -    | -        | -       | Tidak ada parameter yang diperlukan di body. |

```
POST /api/auth/logout
```

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK

Sinyal logout berhasil diterima. Server mengonfirmasi bahwa token valid dan memberikan instruksi bagi klien untuk membersihkan akses di sisi lokal.

```json
{
  "success": true,
  "message": "Logout successfully",
  "data": {
    "status": "access_terminated"
  }
}
```

#### ⚠️ 401 Unauthorized

Respons ketika token akses tidak valid, sudah kedaluwarsa, atau tidak disertakan dalam header.

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

#### 🚫 429 Too Many Requests

Terlalu banyak permintaan logout dalam waktu singkat dari satu identitas atau IP untuk mencegah penyalahgunaan resource.

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

Terjadi kesalahan teknis yang tidak terduga pada server saat memproses validasi logout.

```json
{
  "success": false,
  "message": "An unexpected server error occurred during logout",
  "data": {
    "error_code": "INTERNAL_SERVER_ERROR",
    "errors": null
  }
}
```

---

## Endpoint : `GET /api/v1/auth/me`

### Description :

Endpoint ini digunakan untuk mengambil data profil lengkap pengguna yang sedang aktif (**pemilik token**). Backend akan mendekripsi Access Token dari header untuk mendapatkan identitas pengguna (`user_id`), lalu melakukan query ke database untuk mengambil informasi terbaru. Endpoint ini krusial untuk sinkronisasi state profil dan hak akses saat aplikasi pertama kali dimuat (_Initial Load_).

### Role Based Access Control (RBAC) :

- `Permissions`: `owner, cashier, staff, courier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

Identifikasi identitas dilakukan sepenuhnya melalui dekripsi pada Authorization Header. Tidak ada parameter tambahan yang diperlukan.

| Key  | Type | Location | Default | Description                                                   |
| ---- | ---- | -------- | ------- | ------------------------------------------------------------- |
| None | -    | -        |         | Identifikasi pengguna dilakukan melalui Authorization Header. |

```
GET /api/auth/me
```

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK

Profil pengguna berhasil diambil secara sukses.

```json
{
  "success": true,
  "message": "User profile retrieved successfully",
  "data": {
    "id": 1,
    "full_name": "Farhan Rizki Maulana",
    "username": "farhanrizkimln",
    "role": "owner",
    "is_active": 1,
    "created_at": "2026-01-12 08:12:36"
  }
}
```

#### ⚠️ 401 Unauthorized

Respons ketika token akses tidak valid, sudah kedaluwarsa, atau tidak disertakan dalam header permintaan.

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

Token valid secara teknis, namun akun pengguna dalam status dinonaktifkan (`is_active = 0`).

```json
{
  "success": false,
  "message": "Access denied: your account is currently inactive",
  "data": {
    "error_code": "ACCOUNT_INACTIVE",
    "errors": null
  }
}
```

#### 🚫 429 Too Many Requests

Terlalu banyak permintaan pengambilan profil dalam waktu singkat untuk mencegah eksploitasi _resource server_.

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

Terjadi kesalahan teknis yang tidak terduga pada server (seperti kegagalan query database).

```json
{
  "success": false,
  "message": "An unexpected server error occurred during profile retrieval",
  "data": {
    "error_code": "INTERNAL_SERVER_ERROR",
    "errors": null
  }
}
```
