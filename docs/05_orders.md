# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## ORDERS MODULE SPECIFICATION

---

## Endpoint : `POST /orders`

### Description :

Endpoint ini digunakan untuk mencatat pesanan laundry baru ke dalam sistem. Akses dibatasi hanya untuk **Cashier** dan **Owner**. Backend akan menjalankan transaksi atomik (semua sukses atau semua gagal) untuk mengisi tabel-tabel utama secara otomatis:

1. customers: Mencari data berdasarkan `customer_id` atau otomatis membuat data baru jika data pelanggan belum ada di database.
2. orders: Menyimpan data induk pesanan (Status pengerjaan awal: `pending`).
3. order_items: Menyimpan rincian layanan. Backend menghitung subtotal tiap item berdasarkan harga terbaru di database.
4. deliveries: Menyimpan data logistik/pengiriman jika `is_delivery` bernilai 1.
5. payments: Menginisiasi catatan tagihan. Jika `amount_received` mencukupi, status pembayaran menjadi `paid`, jika tidak maka `unpaid`.
6. status history: Mencatat log pembuatan pesanan sebagai langkah awal pelacakan (_tracking_).

### Role Based Access Control (RBAC) :

- `Permissions`: `owner, cashier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

Bagian ini merinci data yang harus dikirimkan dalam _Body Request_. Kolom yang ditandai opsional akan diisi dengan nilai _default_ oleh _Backend_.

| Key              | Type   | Location | Default | Description                                                  |
| ---------------- | ------ | -------- | ------- | ------------------------------------------------------------ |
| customer_id      | Int    | Body     | null    | ID unik pelanggan (opsional jika pelanggan sudah terdaftar). |
| customer_name    | String | Body     | -       | Nama lengkap pelanggan (wajib jika customer_id null).        |
| customer_phone   | String | Body     | -       | Nomor telepon pelanggan (wajib jika customer_id null).       |
| customer_address | String | Body     | -       | Alamat lengkap pelanggan (wajib jika customer_id null).      |
| is_delivery      | Int    | Body     | 0       | Indikator pengiriman (0 = ambil sendiri, 1 = antar).         |
| notes            | String | Body     | -       | Catatan khusus untuk pesanan ini (opsional).                 |
| deliveries       | Object | Body     | -       | Objek berisi shipping_cost.                                  |
| order_items      | Array  | Body     | -       | Daftar objek service_id, weight_kg, atau quantity.           |
| payment          | Object | Body     | -       | Objek berisi method, amount_received, reference_no.          |

```
  {
  "customer_id": Integer | null,
  "customer_name": String,
  "customer_phone": String,
  "customer_address": String,
  "is_delivery": Integer,
  "notes": String,
  "deliveries": {
    "shipping_cost": "Float"
  },
  "order_items": [
    {
      "service_id": Integer,
      "quantity": Integer | null,
      "weight_kg": Float | null,
      "qty_pieces": Integer | null,
      "item_notes": String
    }
  ],
  "payment": {
    "method": String,
    "amount_received": Float,
    "reference_no": String
  }
}
```

### 🛡️ Logic Guard (Aturan Bisnis & Integritas) :

1. Customer Lookup: Jika customer_id diisi, sistem akan memverifikasi keberadaannya. Jika null, sistem wajib membuat data di tabel customers terlebih dahulu.
2. Price Protection: Harga satuan (unit_price) diambil langsung dari tabel services saat transaksi dibuat untuk menghindari manipulasi harga dari sisi klien.
3. Automatic Estimation: estimated_ready_at dihitung otomatis: created_at + MAX(duration_hours) dari seluruh item layanan yang dipilih.
4. Payment Status:
   - Jika amount_received >= total_price, status payment = paid.
   - Jika amount_received == 0, status payment = unpaid.

### Request Body :

Berisi detail pesanan, pelanggan, dan item cucian. Pada tahap ini, fokus utama adalah validasi fisik cucian (berat/jumlah) dan data pengiriman.

```json
{
  "customer_id": null,
  "customer_name": "Mpok Romlah",
  "customer_phone": "081234567890",
  "customer_address": "Jl. Merpati No. 12",
  "is_delivery": 1,
  "notes": "Jangan dicampur dengan baju luntur",
  "deliveries": {
    "shipping_cost": 10000.0
  },
  "order_items": [
    {
      "service_id": 1,
      "weight_kg": 5.0,
      "qty_pieces": 20,
      "item_notes": "Pisahkan warna putih"
    }
  ],
  "payment": {
    "method": null,
    "amount_received": 0.0,
    "reference_no": null
  }
}
```

### Responses Body :

#### ✅ 201 Created

Deskripsi: Pesanan berhasil dibuat dan disimpan ke seluruh tabel terkait dalam satu transaksi. Response ini mengembalikan `invoice_number` yang digunakan sebagai referensi utama dan objek `payment` untuk memantau status tagihan. Pelunasan di masa mendatang dilakukan melalui endpoint `PATCH /payments/{id}`.

```json
{
  "success": true,
  "message": "Order created successfully",
  "data": {
    "id": 45,
    "invoice_number": "INV-260105-001",
    "is_delivery": 1,
    "total_price": 60000.0,
    "payment_status": "cod_pending",
    "status_internal": "pending",
    "estimated_ready_at": "2026-01-08 13:00:00",
    "notes": "Jangan dicampur dengan baju luntur",
    "created_by": 2,
    "created_by_name": "Siti Aminah",
    "created_at": "2026-01-05 13:00:00",
    "updated_at": null,
    "customer": {
      "id": 101,
      "name": "Mpok Romlah",
      "phone": "081234567890",
      "address": "Jl. Merpati No. 12"
    },
    "order_items": [
      {
        "id": 12,
        "service_id": 1,
        "service_name": "Cuci Kiloan Reguler",
        "item_notes": "Pisahkan warna putih",
        "quantity": null,
        "qty_pieces": 20, // <--- Kolom Baru untuk Tracking Jumlah Helai
        "weight_kg": 5.0,
        "unit": "Kg",
        "unit_price": 10000.0,
        "subtotal": 50000.0
      }
    ],
    "payment": {
      "id": 1,
      "method": null,
      "amount": 60000.0,
      "amount_received": 0.0,
      "amount_change": 0.0,
      "reference_no": null,
      "status": "pending",
      "created_by": 2,
      "collected_by": null
    },
    "delivery": {
      "id": 12,
      "shipping_cost": 10000.0,
      "courier_id": null,
      "courier_name": null,
      "courier_phone": null,
      "courier_departed_at": null,
      "courier_arrived_at": null,
      "cod_collected_amount": 0.0
    },
    "status_history": [
      {
        "id": 1,
        "previous_status": null,
        "new_status": "pending",
        "actor_name": "Siti Aminah",
        "actor_role": "cashier",
        "notes": "Initial order creation",
        "created_at": "2026-01-05 13:00:00"
      }
    ]
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi ketika validasi input gagal memenuhi kriteria bisnis, seperti format data yang salah, ID layanan yang tidak aktif, atau parameter wajib yang tidak diisi.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "order_items": "At least one service item is required",
      "deliveries": "Shipping cost is required when is_delivery is 1"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Terjadi ketika permintaan tidak menyertakan token akses yang valid, token telah kedaluwarsa, atau format token di dalam _header_ `Authorization` salah.

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

Terjadi ketika pengguna berhasil terautentikasi tetapi tidak memiliki hak akses (Role) yang diizinkan untuk melakukan aksi ini (selain Cashier atau Owner).

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

Terjadi jika ada pertentangan dengan data yang sudah ada di server, misalnya upaya membuat pesanan dengan nomor invoice yang sudah terdaftar secara tidak sengaja.

```json
{
  "success": false,
  "message": "Data already exists",
  "data": {
    "error_code": "DUPLICATE_DATA",
    "errors": {
      "invoice_number": "Invoice number already in use"
    }
  }
}
```

#### 🚫 429 Too Many Requests

Sistem membatasi frekuensi pembuatan pesanan dalam waktu singkat untuk melindungi database dari beban berlebih atau percobaan serangan spamming.

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

Terjadi kegagalan sistem yang tidak terduga pada server atau kegagalan transaksi pada database (misalnya koneksi terputus saat proses penyimpanan).

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

## Endpoint : `GET /orders`

### Description :

Endpoint ini digunakan untuk mengambil daftar seluruh pesanan laundry dalam format ringkasan (_Summary_). Dirancang khusus untuk kebutuhan operasional harian pada dashboard **Kasir** (antrean masuk) dan **Kurir** (antrean kirim). Sistem mendukung **Pagination** untuk efisiensi beban kerja server, serta **Search** dan **Filtering** multi-parameter.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner, cashier, staff, courier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Bagian ini digunakan untuk mengontrol aliran data yang keluar dari database. Dengan menggunakan parameter ini, Anda memastikan server Golang tidak menarik data yang tidak diperlukan (Query Optimization).

| Key             | Type   | Location | Default    | Description                                                                |
| --------------- | ------ | -------- | ---------- | -------------------------------------------------------------------------- |
| page            | Int    | Query    | 1          | Nomor halaman data.                                                        |
| per_page        | Int    | Query    | 10         | Jumlah data per halaman (Maks. 100).                                       |
| search          | String | Query    | -          | Partial search berdasarkan No. Invoice atau Nama Pelanggan.                |
| status_internal | String | Query    | -          | Filter status internal proses (pending, in-progress, ready-delivery, dll). |
| payment_status  | String | Query    | -          | Filter pembayaran (paid dan unpaid).                                       |
| sort_by         | String | Query    | created_at | Mengurutkan berdasarkan kolom tertentu.                                    |
| order           | String | Query    | desc       | asc (Terlama/A-Z) atau desc (Terbaru/Z-A).                                 |

```
GET /api/orders?page=1&per_page=10&status_internal=pending&search=Romlah&sort_by=created_at&order=desc
```

### Request Body :

```json
None (Kosong, karena method GET tidak boleh punya body).
```

### Responses Body :

#### ✅ 200 OK

Data berhasil diambil. Struktur ini memisahkan antara data (daftar pesanan) dan meta (informasi halaman) agar Frontend dapat membuat komponen pagination dengan mudah.

```json
{
  "success": true,
  "message": "Orders retrieved successfully",
  "data": [
    {
      "id": 45,
      "invoice_number": "INV-260105-001",
      "is_delivery": 1,
      "total_price": 60000.0,
      "payment_status": "cod_pending",
      "status_internal": "pending",
      "estimated_ready_at": "2026-01-08 13:00:00",
      "created_by": 2,
      "created_by_name": "Siti Aminah",
      "created_at": "2026-01-05 13:00:00",
      "updated_at": null,
      "customer": {
        "id": 101,
        "name": "Mpok Romlah",
        "phone": "081234567890"
      },
      "delivery": {
        "id": 12,
        "shipping_cost": 10000.0
      }
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total_items": 50,
    "total_pages": 5
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi jika parameter filter salah format.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "page": "Page must be a positive integer"
    }
  }
}
```

#### 🚫 403 Forbidden

Audit: Menjaga insting keamanan Anda agar Customer tidak bisa bypass data.

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

Terjadi kegagalan pada sistem internal atau kesalahan dari database.

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

## Endpoint : `GET /orders/{id}`

### Description :

Endpoint ini digunakan untuk mengambil data detail lengkap dari satu pesanan tertentu. Backend akan melakukan operasi **JOIN** atau pemanggilan data dari **6 tabel** (`orders`, `customers`, `order_items`, `deliveries`, `payments`, dan `status_history`) untuk memberikan gambaran utuh mengenai satu transaksi secara transparan.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner, cashier, staff, courier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Bagian ini menggunakan _Path Parameter_ untuk mengidentifikasi sumber data secara spesifik di database.

| Key | Type    | Location | Default | Description                                 |
| --- | ------- | -------- | ------- | ------------------------------------------- |
| id  | Integer | Path     | -       | ID unik pesanan yang ingin diambil datanya. |

```
GET /api/orders/45
```

### 🛡️ Logic Guard (Integritas Data) :

1. Full View Consistency: Nominal uang (total_price, shipping_cost, subtotal) wajib menggunakan tipe Float untuk menjaga presisi desimal sesuai kesepakatan database.
2. No Debt Policy: Karena sistem tidak mengenal hutang, payment_status pada level order harus sinkron dengan status di objek payment (Hanya paid atau unpaid).
3. State Visibility: qty_pieces disajikan untuk membantu Staff melakukan verifikasi jumlah helai fisik saat proses pencucian agar tidak ada pakaian yang tertukar atau hilang.

### Request Body :

```json
None (Kosong, karena method GET tidak boleh punya body).
```

### Responses Body :

#### ✅ 200 OK

Data ditemukan dan dikembalikan secara lengkap. Struktur data tetap konsisten menggunakan objek singular untuk hubungan 1-ke-1 dan array untuk riwayat/item.

```json
{
  "success": true,
  "message": "Order detail retrieved successfully",
  "data": {
    "id": 45,
    "invoice_number": "INV-260105-001",
    "is_delivery": 1,
    "total_price": 60000.0,
    "payment_status": "cod_pending",
    "status_internal": "pending",
    "estimated_ready_at": "2026-01-08 13:00:00",
    "notes": "Jangan dicampur dengan baju luntur",
    "created_by": 2,
    "created_by_name": "Siti Aminah",
    "created_at": "2026-01-05 13:00:00",
    "updated_at": null,
    "customer": {
      "id": 101,
      "name": "Mpok Romlah",
      "phone": "081234567890",
      "address": "Jl. Merpati No. 12"
    },
    "order_items": [
      {
        "id": 12,
        "service_id": 1,
        "service_name": "Cuci Kiloan Reguler",
        "item_notes": "Pisahkan warna putih",
        "quantity": null,
        "qty_pieces": 20,
        "weight_kg": 5.0,
        "unit": "Kg",
        "unit_price": 10000.0,
        "subtotal": 50000.0
      }
    ],
    "payment": {
      "id": 1,
      "method": null,
      "amount": 60000.0,
      "amount_received": 0.0,
      "amount_change": 0.0,
      "reference_no": null,
      "status": "pending",
      "created_by": 2,
      "collected_by": null
    },
    "delivery": {
      "id": 12,
      "shipping_cost": 10000.0,
      "courier_id": null,
      "courier_name": null,
      "courier_phone": null,
      "courier_departed_at": null,
      "courier_arrived_at": null,
      "cod_collected_amount": 0.0
    },
    "status_history": [
      {
        "id": 1,
        "previous_status": null,
        "new_status": "pending",
        "actor_name": "Siti Aminah",
        "actor_role": "cashier",
        "notes": "Initial order creation",
        "created_at": "2026-01-05 13:00:00"
      }
    ]
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi jika format ID yang dikirimkan pada URL tidak valid (bukan angka).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "id": "Order ID must be a valid integer"
    }
  }
}
```

#### 🚫 403 Forbidden

Terjadi jika role user (misal: Customer) mencoba mengakses detail pesanan milik orang lain (IDOR Protection).

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

Terjadi jika ID pesanan yang dicari tidak ada di database MySQL.

```json
{
  "success": false,
  "message": "Order not found",
  "data": {
    "error_code": "RESOURCE_NOT_FOUND",
    "errors": null
  }
}
```

#### 🚫 429 Too Many Requests

Mencegah upaya pemindaian ID secara otomatis (ID Scanning/Brute Force).

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

Terjadi kegagalan koneksi database atau kesalahan logika JOIN pada query Golang.

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

## Endpoint : `PUT /orders/{id}`

### Description :

Endpoint ini digunakan oleh **Owner/Cashier** untuk melakukan pembaruan data pesanan secara menyeluruh. **Aturan Bisnis Utama**: Perubahan hanya diizinkan jika `status_internal` masih bernilai `pending`. Jika pesanan sudah mulai diproses (`in-progress`), data dikunci untuk menjaga integritas laporan. Sistem akan menghitung ulang `total_price` dan memperbarui tagihan pada tabel `payments` secara otomatis.

### Role Based Access Control (RBAC) :

- `Permissions`: `cashier, owner`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

Bagian ini menggunakan _Path Parameter_ untuk mengunci ID pesanan yang akan direvisi di database.

| Key | Type | Location | Default | Description                            |
| --- | ---- | -------- | ------- | -------------------------------------- |
| id  | Int  | Path     | -       | ID unik pesanan yang ingin diperbarui. |

```
PUT /api/orders/45
```

### 🛡️ Logic Guard (Aturan Bisnis & Integritas) :

1. Status Restriction: Permintaan wajib ditolak (400 Bad Request) jika pesanan sudah melewati tahap pending di database.
2. Financial Integrity: Seluruh nominal menggunakan tipe Float. Sistem akan menghitung ulang total_price berdasarkan harga layanan terbaru.
3. Payment Synchronization: Jika pesanan direvisi dan harga berubah, transaksi pembayaran lama di tabel payments yang masih pending akan disesuaikan nilainya. Jika sudah confirmed, maka Admin harus melakukan penyesuaian manual melalui endpoint pembayaran.
4. No Debt Policy: Meskipun ada status transaksi pembayaran, sistem tetap memastikan pesanan tidak bisa dianggap lunas (paid) sebelum transaksi di tabel payments mencapai status confirmed

### Request Body :

Mengirimkan struktur data lengkap untuk menggantikan data lama. Gunakan customer_id jika pelanggan sudah terdaftar.

```json
{
  "customer_id": 101,
  "customer_name": "Mpok Romlah",
  "customer_phone": "081234567890",
  "customer_address": "Jl. Merpati No. 12",
  "is_delivery": 1,
  "notes": "Jangan dicampur dengan baju luntur",
  "order_items": [
    {
      "service_id": 1,
      "quantity": null,
      "qty_pieces": 40, // update jumlah helai pakaian
      "weight_kg": 10.0, // update berat laundry kiloan
      "item_notes": "Pisahkan warna putih"
    }
  ]
}
```

### Responses Body :

#### ✅ 200 OK

Pesanan berhasil diperbarui. Response mengembalikan objek data terbaru yang sudah dikalkulasi ulang oleh server.

```json
{
  "success": true,
  "message": "Order updated successfully",
  "data": {
    "id": 45,
    "invoice_number": "INV-260105-001",
    "is_delivery": 1,
    "total_price": 110000.0,
    "payment_status": "cod_pending",
    "status_internal": "pending",
    "estimated_ready_at": "2026-01-08 13:00:00",
    "notes": "Jangan dicampur dengan baju luntur",
    "created_by": 2,
    "created_by_name": "Siti Aminah",
    "created_at": "2026-01-05 13:00:00",
    "updated_at": "2026-01-05 13:10:00",
    "customer": {
      "id": 101,
      "name": "Mpok Romlah",
      "phone": "081234567890",
      "address": "Jl. Merpati No. 12"
    },
    "order_items": [
      {
        "id": 12,
        "service_id": 1,
        "service_name": "Cuci Kiloan Reguler",
        "item_notes": "Pisahkan warna putih",
        "qty_pieces": 40, // update jumlah helai pakaian
        "weight_kg": 10.0, // update tambah berat layanan kiloan
        "unit": "Kg",
        "unit_price": 10000.0,
        "subtotal": 100000.0
      }
    ],
    "payment": {
      "id": 1,
      "method": null,
      "amount": 110000.0,
      "amount_received": 0.0,
      "amount_change": 0.0,
      "reference_no": null,
      "status": "pending",
      "created_by": 2,
      "collected_by": null
    },
    "delivery": {
      "id": 12,
      "shipping_cost": 10000.0,
      "courier_id": null,
      "courier_name": null,
      "courier_phone": null,
      "courier_departed_at": null,
      "courier_arrived_at": null,
      "cod_collected_amount": 0.0
    },
    "status_history": [
      {
        "id": 1,
        "previous_status": null,
        "new_status": "pending",
        "actor_name": "Siti Aminah",
        "actor_role": "cashier",
        "notes": "Initial order creation",
        "created_at": "2026-01-05 13:00:00"
      }
    ]
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi jika format input salah atau melanggar aturan bisnis (mencoba edit pesanan yang sudah diproses).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "status": "Order can only be edited when status is pending"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Terjadi jika token akses tidak valid atau sudah kadaluwarsa.

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

Terjadi jika role user tidak memiliki izin (misal: Staff atau Courier mencoba melakukan edit total).

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

Terjadi jika ID pesanan yang ingin di-update tidak ditemukan di database.

```json
{
  "success": false,
  "message": "Order not found",
  "data": {
    "error_code": "RESOURCE_NOT_FOUND",
    "errors": null
  }
}
```

#### 🚫 409 Conflict

Terjadi jika terdapat duplikasi data unik yang tidak sengaja tercipta saat proses update.

```json
{
  "success": false,
  "message": "Data already exists",
  "data": {
    "error_code": "DUPLICATE_DATA",
    "errors": {
      "invoice_number": "Invoice number 'INV-260105-001' is already in use"
    }
  }
}
```

#### 🚫 429 Too Many Requests

Perlindungan server dari upaya perubahan data yang terlalu masif dalam waktu singkat.

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

Terjadi kesalahan pada transaksi database atau kegagalan sistem internal saat pemrosesan.

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

## Endpoint : `PATCH /orders/{id}`

### Description :

Endpoint ini digunakan untuk memperbarui status operasional pesanan secara parsial (State Transition). Berbeda dengan PUT yang melakukan update data fisik, PATCH berfokus pada pergerakan workflow. Backend akan menjalankan transaksi atomik untuk memperbarui dua tabel utama:

1. orders: Memperbarui kolom status_internal dan updated_at.
2. status_history: Menambahkan baris riwayat baru untuk melacak siapa yang mengubah status, kapan, dan alasan perubahannya (audit trail).

### Role Based Access Control (RBAC) :

- `Permissions`: `owner, cashier, staff, courier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

Bagian ini menggunakan Path Parameter untuk menentukan pesanan mana yang akan diproses transisi statusnya.

| Key | Type | Location | Default | Description                                  |
| --- | ---- | -------- | ------- | -------------------------------------------- |
| id  | Int  | Path     | -       | ID unik pesanan yang ingin diubah statusnya. |

```
PATCH /api/orders/45
```

### 🛡️ Logic Guard (Aturan State Machine) :

Untuk menjaga integritas operasional, sistem wajib memvalidasi transisi status berdasarkan Role:

1. Sequence Validation: Status tidak boleh mundur (misal: dari ready kembali ke pending) kecuali dilakukan oleh owner.
2. Role Restriction:
   - staff: Hanya bisa mengubah status ke in-progress atau ready.
   - courier: Hanya bisa mengubah status ke being-delivered atau completed (jika antar-jemput).
   - cashier/owner: Memiliki otoritas penuh untuk semua status, termasuk cancelled.
3. Payment Requirement: Status completed (pesanan diambil/diterima pelanggan) hanya bisa dipicu jika payment_status sudah paid. Sesuai kebijakan Tidak Bisa Hutang, barang tidak boleh keluar jika belum lunas.

### Request Body :

Hanya mengirimkan informasi transisi status. Field notes bersifat opsional namun sangat disarankan untuk audit internal.

```json
{
  "new_status": "in-progress",
  "notes": "Pakaian mulai dimasukkan ke mesin cuci nomor 03"
}
```

### Responses Body :

#### ✅ 200 OK

Status berhasil diperbarui. Mengikuti prinsip "The Finished Plate", response mengembalikan data lengkap agar UI dapat langsung memindahkan card pesanan ke kolom yang sesuai tanpa re-fetch.

```json
{
  "success": true,
  "message": "Order updated successfully",
  "data": {
    "id": 45,
    "invoice_number": "INV-260105-001",
    "is_delivery": 1,
    "total_price": 110000.0,
    "payment_status": "cod_pending",
    "status_internal": "in-progress",
    "estimated_ready_at": "2026-01-08 13:00:00",
    "notes": "Jangan dicampur dengan baju luntur",
    "created_by": 2,
    "created_by_name": "Siti Aminah",
    "created_at": "2026-01-05 13:00:00",
    "updated_at": "2026-01-06 10:00:00",
    "customer": {
      "id": 101,
      "name": "Mpok Romlah",
      "phone": "081234567890",
      "address": "Jl. Merpati No. 12"
    },
    "order_items": [
      {
        "id": 12,
        "service_id": 1,
        "service_name": "Cuci Kiloan Reguler",
        "item_notes": "Pisahkan warna putih",
        "quantity": null,
        "qty_pieces": 40,
        "weight_kg": 10.0,
        "unit": "Kg",
        "unit_price": 10000.0,
        "subtotal": 100000.0
      }
    ],
    "payment": {
      "id": 1,
      "method": null,
      "amount": 110000.0,
      "amount_received": 0.0,
      "amount_change": 0.0,
      "reference_no": null,
      "status": "pending",
      "created_by": 2,
      "collected_by": null
    },
    "delivery": {
      "id": 12,
      "shipping_cost": 10000.0,
      "courier_name": null,
      "courier_phone": null,
      "courier_departed_at": null,
      "courier_arrived_at": null,
      "cod_collected_amount": 0.0
    },
    "status_history": [
      {
        "id": 1,
        "previous_status": null,
        "new_status": "pending",
        "actor_name": "Siti Aminah",
        "actor_role": "cashier",
        "notes": "Initial order creation",
        "created_at": "2026-01-05 13:00:00"
      },
      {
        "id": 2,
        "previous_status": "pending",
        "new_status": "in-progress",
        "actor_name": "Fadhillah Kurnia",
        "actor_role": "staff",
        "notes": "Pakaian mulai dimasukkan ke mesin cuci nomor 03",
        "created_at": "2026-01-06 10:00:00"
      }
    ]
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi jika transisi status melanggar aturan bisnis (misal: pesanan cancelled tidak bisa diubah ke in-progress).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "status": "Cannot change status from 'ready' back to 'pending'"
    }
  }
}
```

#### 🚫 403 Forbidden

Audit: Menjamin bahwa hanya user dengan role tertentu yang bisa mengubah status spesifik (misal: hanya Courier yang bisa mengubah ke being-delivered).

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

Terjadi jika ID pesanan yang ingin di-patch tidak ditemukan di database.

```json
{
  "success": false,
  "message": "Order not found",
  "data": {
    "error_code": "RESOURCE_NOT_FOUND",
    "errors": null
  }
}
```

#### 🚫 409 Conflict

Terjadi jika status yang ada di database sudah berubah sejak user terakhir kali mengambil data, sehingga transisi status yang diminta menjadi tidak valid.

```json
{
  "success": false,
  "message": "The order has been updated by another user",
  "data": {
    "error_code": "STATE_CONFLICT",
    "errors": {
      "current_status": "Status has changed to 'washing', please refresh your data."
    }
  }
}
```

#### 🚫 429 Too Many Requests

Mencegah spamming aksi pada tombol operasional yang berakibat pada penulisan log berlebihan di database.

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

Terjadi kegagalan transaksi pada database saat mencoba menulis riwayat status baru.

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
