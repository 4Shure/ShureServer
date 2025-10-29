package entity

type Appointment struct {
	ID        int   `gorm:"primaryKey"`
	BeginsAt  int64 `gorm:"not null"`
	EndsAt    int64 `gorm:"not null"`
	UserID    int   `gorm:"not null"` // References: users(id)
	IsDeleted bool  `gorm:"not null"`
	CreatedAt int64 `gorm:"not null"`
	UpdatedAt int64 `gorm:"not null"`
	Title     *string

	// Relations
	CreatedBy User `gorm:"foreignKey:UserID;references:ID"`
}
