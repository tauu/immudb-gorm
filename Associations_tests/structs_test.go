package test_associations

import "gorm.io/gorm"

// Belongs structs
type Employee struct {
	gorm.Model
	Name      string
	CompanyID int
	Company   Company
}

type Company struct {
	ID   int
	Name string
}

// HasMany structs
type Owner struct {
	gorm.Model
	Name        string
	Restaurants []Restaurant
}

type Restaurant struct {
	gorm.Model
	Name    string
	OwnerID uint
}

// HasOne structs
type User struct {
	gorm.Model
	Name       string
	CreditCard CreditCard
}

type CreditCard struct {
	gorm.Model
	Number string
	UserID uint
}

// ManyToMany structs
type Student struct {
	gorm.Model
	Name      string
	Age       int
	Languages []Language `gorm:"many2many:student_languages;"`
}

type Language struct {
	gorm.Model
	Name string
}
