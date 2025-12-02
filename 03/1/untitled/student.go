package main

type Student struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"type:varchar(50);not null"`
	Age   int    `gorm:"not null"`
	Grade string `gorm:"type:varchar(20);not null"`
}
