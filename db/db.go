package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go-bank/account"
	"go-bank/client"
	"go-bank/user"
	"log"
	"os"
)

var db *gorm.DB

func Connect() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Erro ao carregar arquivo .env: %v", err)
	}

	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("PORT")
	dbUser := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DATABASE")

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, dbUser, database, password, dbPort)
	fmt.Println(dbURI)

	var err error
	db, err = gorm.Open(dialect, dbURI)
	if err != nil {
		log.Fatalf("Erro ao conectar com o banco de dados: %v", err)
	}
	fmt.Println("Database connection established successfully!")

	DBMigrate()
}

func GetDB() *gorm.DB {
	return db
}

func DBMigrate() {
	// Migrate schemas
	if err := db.AutoMigrate(&user.User{}).Error; err != nil {
		log.Fatalf("Fail to migrate user table: %v", err)
	}
	if err := db.AutoMigrate(&account.Account{}).Error; err != nil {
		log.Fatalf("Fail to migrate account table: %v", err)
	}
	if err := db.AutoMigrate(&client.Client{}).Error; err != nil {
		log.Fatalf("Fail to migrate client table: %v", err)
	}
	if err := db.AutoMigrate(&account.Transaction{}).Error; err != nil {
		log.Fatalf("Fail to migrate transaction table: %v", err)
	}
	fmt.Println("Database migration completed successfully!")
}
