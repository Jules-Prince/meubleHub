package models

type Room struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Name   string `json:"name" gorm:"not null"`
	HomeID uint   `json:"home_id"`  // Foreign key linking to a Home
}