package entity

type Appointment struct {
	ID         int   `gorm:"primaryKey"`
	BeginsAt   int64 `gorm:"not null"`
	EndsAt     int64 `gorm:"not null"`
	CustomerID int   `gorm:"not null"` // References: users(id)
	IsDeleted  bool  `gorm:"not null"`
	CreatedAt  int64 `gorm:"not null"`
	UpdatedAt  int64 `gorm:"not null"`

	// Relations
	CreatedBy User `gorm:"foreignKey:CustomerID;references:ID"`
}
