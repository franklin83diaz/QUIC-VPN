package pkg

import "gorm.io/gorm"

type Test struct {
	gorm.Model
	TestString string
}

type User struct {
	gorm.Model
	Name   string
	Last   int
	Active bool `gorm:"default:false"`
	Ip     string
}

type Network struct {
	gorm.Model
	CIDR string
}
