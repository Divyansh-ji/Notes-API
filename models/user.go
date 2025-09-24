package models

type User struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Email    string `gorm:"uniqueIndex"`
	Password string
	Notes    []Note // has-many relation; exported field name
}
