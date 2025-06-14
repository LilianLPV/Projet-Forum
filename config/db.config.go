package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() (*sql.DB, error) {

	user := GetEnvWithDefault("DB_USER", "")
	pwd := GetEnvWithDefault("DB_PWD", "")
	host := GetEnvWithDefault("DB_HOST", "")
	port := GetEnvWithDefault("DB_PORT", "")
	name := GetEnvWithDefault("DB_NAME", "")

	if user == "" || host == "" || port == "" || name == "" {
		return nil, fmt.Errorf(" Erreur connection base de données - Donneés de connexions manquantes")
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pwd, host, port, name)

	dbContext, dbContextErr := sql.Open("mysql", connectionString)
	if dbContextErr != nil {
		return nil, fmt.Errorf(" Erreur connection base de données - Erreur : \n\t %s", dbContextErr.Error())
	}

	pingErr := dbContext.Ping()
	if pingErr != nil {
		dbContext.Close()
		return nil, fmt.Errorf(" Erreur ping base de données - Erreur : \n\t %s", pingErr.Error())
	}

	return dbContext, nil
}
