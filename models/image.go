type Image struct {
	ID       uuid.UUID   `gorm:"primaryKey"`
	Filename string `gorm:"not null"`
	Data     []byte `gorm:"not null"`
}
