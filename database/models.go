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

	ID          uint `gorm:"primaryKey"`
	Activated   bool
	Name        string
	WebSiteURL  string
	Icon        string
	Offres      []*Offre `gorm:"many2many:offer_sites;"`
	AccountID   uint
	AccessRules []AccessRule
	Deleted     bool
}

type AccessRule struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	SiteID    uint
	Name      string
	Offres    []*Offre `gorm:"many2many:accessrules_offre;"`
	IsDefault bool
	Deleted   bool
}

type Offre struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	Activated bool
	Name      string
	Price     float64
	Frequency string

	Title       string
	Description string
	Icon        string

	AccessRules []*AccessRule `gorm:"many2many:accessrules_offre;"`
	Sites       []*Site       `gorm:"many2many:offer_sites;"`
	AccountID   uint
	Deleted     bool
}

type Account struct {
	gorm.Model
	ID      uint `gorm:"primaryKey"`
	Name    string
	Sites   []Site
	Offres  []Offre
	Deleted bool
}

func InitDatabases() {
	dsn := "root:1234@tcp(127.0.0.1:3306)/paywall?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(
		// &Rule{},
		&AccessRule{},
		&User{},
		&Client{},
		&Site{},
		&Offre{},
		&Account{},
	)
	if err != nil {
		panic(err)
	}
	Db = db

}
