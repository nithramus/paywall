package database

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// var ctx context
var Db *gorm.DB

type User struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	AccountID uint
	Account   Account
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Client struct {
	gorm.Model

	ID               uint `gorm:"primaryKey"`
	FirstName        string
	LastName         string
	SubscribedOffres []Offre   `gorm:"many2many:offer_clients;"`
	CreatedAt        time.Time // Set to current time if it is zero on creating
	UpdatedAt        int       // Set to current unix seconds on updaing or if it is zero on creating
	Deleted          bool
}
type Site struct {
	gorm.Model

	ID          uint     `gorm:"primaryKey"`
	WebOffreUrl string   `json: webOffreUrl`
	Offres      []*Offre `gorm:"many2many:offer_sites;"`
	AccountID   uint
}

type Offre struct {
	gorm.Model
	ID      uint `gorm:"primaryKey"`
	Price   float64
	Deleted bool    `json: deleted`
	Sites   []*Site `gorm:"many2many:offer_sites;"`

	AccountID uint
}

type Account struct {
	gorm.Model
	ID      uint `gorm:"primaryKey"`
	Sites   []Site
	Offres  []Offre
	Deleted bool `json: deleted`
}

func InitDatabases() {
	dsn := "root:1234@tcp(127.0.0.1:3306)/paywall?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}

}
