type Image struct {
	ID       uuid.UUID   `gorm:"primaryKey"`
	Filename string `gorm:"not null"`
	Data     []byte `gorm:"not null"`
}
func (i *Image) BeforeCreate(d *gorm.DB) (err error) {
	i.ID = uuid.New()
	return
}


