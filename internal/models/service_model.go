package models

import "time"

// Service merepresentasikan struktur tabel 'services' di database
type Service struct {
	ID            int64      `db:"id"`
	Code          string     `db:"code"`
	ServiceName   string     `db:"service_name"`
	Unit          string     `db:"unit"`           // Enum: 'kg' atau 'pcs'
	Price         float64    `db:"price"`          // Menyimpan DECIMAL(15,2)
	IsActive      bool       `db:"is_active"`      // TINYINT(1) -> true/false
	CreatedAt     time.Time  `db:"created_at"`     // Tanpa pointer (Selalu ada isinya)
	UpdatedAt     *time.Time `db:"updated_at"`     // Pakai pointer (Awalnya NULL)
	CategoryID    int64      `db:"category_id"`    // Foreign Key
	DurationHours int        `db:"duration_hours"` // Estimasi pengerjaan (jam)
}

// ServiceWithCategory digunakan untuk menampung hasil query JOIN dengan tabel 'service_categories'
// Ini akan dipakai nanti di layer Repository untuk endpoint GET List dan GET Detail
type ServiceWithCategory struct {
	Service // Embed (Menyisipkan) struct Service utama di atas

	// Tambahan kolom dari tabel service_categories hasil operasi JOIN (AS)
	CategoryName        string  `db:"category_name"`
	CategoryDescription *string `db:"category_description"` // Pakai pointer karena di DB bisa bernilai NULL
}
