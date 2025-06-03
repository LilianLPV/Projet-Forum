package config

import (
	_ "github.com/go-sql-driver/mysql"
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
