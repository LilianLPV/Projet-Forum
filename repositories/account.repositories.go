package repositories

import (
	"database/sql"
	"fmt"
	"forum/models"
	"time"
)

type AccountRepositories struct {
	db *sql.DB
}

func InitAccountRepositories(db *sql.DB) *AccountRepositories {
	return &AccountRepositories{db}
}

func (r *AccountRepositories) Createaccount(Account models.Account) (int, error) {
	query := "INSERT INTO `account`(`user_id`,`user_name`,`user_password`,`user_email`,`user_profile_picture`,`user_bio`,`user_last_connection`) VALUES (?,?,?,?,?,?,?);"

	sqlResult, sqlErr := r.db.Exec(query,
		Account.Id,
		Account.Password,
		Account.Email,
		Account.Profilepicture,
		Account.Bio,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	if sqlErr != nil {
		return -1, fmt.Errorf(" Erreur ajout account - Erreur : \n\t %s", sqlErr.Error())
	}

	id, idErr := sqlResult.LastInsertId()
	if idErr != nil {
		return -1, fmt.Errorf(" Erreur ajout account - Erreur récupération identifiant : \n\t %s", idErr.Error())
	}

	return int(id), nil
}

func (r *AccountRepositories) ReadAll() ([]models.Account, error) {
	var listaccounts []models.Account
	sqlResult, sqlErr := r.db.Query("SELECT * FROM `account`;")
	if sqlErr != nil {
		return listaccounts, fmt.Errorf(" Erreur récupération account - Erreur : \n\t %s", sqlErr.Error())
	}

	defer sqlResult.Close()

	for sqlResult.Next() {
		var Account models.Account
		errScan := sqlResult.Scan(&Account.Id, &Account.Username, &Account.Password, &Account.Email, &Account.Profilepicture, &Account.Bio, &Account.Lastconnection)
		if errScan != nil {
			continue
		}
		listaccounts = append(listaccounts, Account)
	}

	return listaccounts, nil
}

func (r *AccountRepositories) ReadById(id int) (models.Account, error) {
	var Account models.Account
	sqlErr := r.db.QueryRow("SELECT * FROM `account` WHERE `account`.id = ?;", id).
		Scan(&Account.Id, &Account.Username, &Account.Password, &Account.Email, &Account.Profilepicture, &Account.Bio, &Account.Lastconnection)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return models.Account{}, nil
		}
		return models.Account{}, fmt.Errorf(" Erreur récupération account - Erreur : \n\t %s", sqlErr.Error())
	}

	return Account, nil
}
func (r *AccountRepositories) ReadByEmail(Email string) (models.Account, error) {
	var Account models.Account
	sqlErr := r.db.QueryRow("SELECT * FROM `account` WHERE `account`.Email = ?;", Email).
		Scan(&Account.Id, &Account.Username, &Account.Password, &Account.Email, &Account.Profilepicture, &Account.Bio, &Account.Lastconnection)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return models.Account{}, nil
		}
		return models.Account{}, fmt.Errorf(" Erreur récupération account - Erreur : \n\t %s", sqlErr.Error())
	}

	return Account, nil
}
func (r *AccountRepositories) ReadByUsername(Username string) (models.Account, error) {
	var Account models.Account
	sqlErr := r.db.QueryRow("SELECT * FROM `account` WHERE `account`.Username = ?;", Username).
		Scan(&Account.Id, &Account.Username, &Account.Password, &Account.Email, &Account.Profilepicture, &Account.Bio, &Account.Lastconnection)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return models.Account{}, nil
		}
		return models.Account{}, fmt.Errorf(" Erreur récupération account - Erreur : \n\t %s", sqlErr.Error())
	}

	return Account, nil
}

/*
func (r *AccountRepositories) UpdateaccountById(Account models.Account) error {
	query := "UPDATE `account` SET `Id`=?,`d`=?,`prix`=?,`categorie_id`=? WHERE `account`.id=?;"

	sqlResult, sqlErr := r.db.Exec(query,
		Account.Username,
		Account.Description,
		Account.Price,
		Account.CategorieId,
		Account.Id)

	if sqlErr != nil {
		return fmt.Errorf(" Erreur modification account - Erreur : \n\t %s", sqlErr.Error())
	}

	if nbrRow, _ := sqlResult.RowsAffected(); nbrRow <= 0 {
		return fmt.Errorf(" Erreur modification account - Aucune ligne modifiée")
	}

	return nil
}
*/
