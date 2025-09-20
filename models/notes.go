package models

type Note struct {
	ID      uint `gorm:"primaryKey"`
	Title   string
	Content string
	UserID  uint // foreign key that links to User.ID
	//User    //User // optional backref
}
