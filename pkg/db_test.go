package pkg_test

import (
	. "QUIC-VPN/pkg"
	"reflect"
	"testing"
)

func TestConnect(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"TestConnect", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Connect(); reflect.TypeOf(got).String() != "*gorm.DB" {
				t.Errorf("Connect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateTable(t *testing.T) {

	test := Test{
		TestString: "test",
	}

	var test2 Test

	CreateTable(&test, Db)

	Db.Create(&test)

	Db.First(&test2)

	if test2.TestString != "test" {
		t.Errorf("CreateTable() = %v, want %v", test2.TestString, "test")
	}

}
