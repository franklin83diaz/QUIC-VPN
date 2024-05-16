package pkg

import "gorm.io/gorm"

type Test struct {
	gorm.Model
	TestString string
}

type User struct {
	gorm.Model
	Name         string
	Last         string
	Username     string
	HashPassword string
	Email        string
	Active       bool `gorm:"default:false"`
	Ip           string
}

type Network struct {
	gorm.Model
	CIDR string
}
