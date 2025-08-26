package database

type FailedLogin struct {
	ID        uint `gorm:"primarykey"`
	Email string
	Password string
}