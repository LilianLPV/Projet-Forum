package config

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

/*
func test() {

	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		log.Fatalf("Impossible de lancer le projet - Erreur chargement fichier d'environnement \n%s", errEnv.Error())
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatalf("Impossible de lancer le projet - Erreur chargement fichier d'environnement \n%s", dbName)
	}
	dbport, dbPortCheck := os.LookupEnv("DB_PORT")
	if !dbPortCheck {
		log.Fatalf("Impossible de lancer le projet - Erreur chargement fichier d'environnement \n%s", dbport)
	}
	println(dbport, dbName)
}
*/

var DbContext *sql.DB

func InitDB() {
	user := os.Getenv("DB_USER")
	pwd := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

}
