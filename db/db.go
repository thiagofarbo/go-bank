package db

import (
	"fmt"
	"github.com/joho/godotenv"
	"go-bank/account"
	"go-bank/user"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
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

	// String de conexão com o banco de dados
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, dbUser, database, password, dbPort)
	fmt.Println(dbURI)

	// Abrir conexão com o banco de dados
	var err error
	db, err = gorm.Open(dialect, dbURI)
	if err != nil {
		log.Fatalf("Erro ao conectar com o banco de dados: %v", err)
	}
	fmt.Println("Conexão com o banco de dados estabelecida com sucesso!")

	// Migrar esquemas
	if err := db.AutoMigrate(&account.Account{}).Error; err != nil {
		log.Fatalf("Falha ao migrar banco de dados: %v", err)
	}
	if err := db.AutoMigrate(&user.User{}).Error; err != nil {
		log.Fatalf("Falha ao migrar banco de dados: %v", err)
	}
	if err := db.AutoMigrate(&account.Transaction{}).Error; err != nil {
		log.Fatalf("Falha ao migrar tabela transaction: %v", err)
	}

	fmt.Println("Migração de banco de dados concluída com sucesso!")
}

func GetDB() *gorm.DB {
	return db
}
