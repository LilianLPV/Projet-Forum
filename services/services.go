package services

import (
	"database/sql"
	"fmt"
	"forum/models"
	"forum/repositories"
)

type AccountService struct {
	AccountRepo *repositories.AccountRepositories
}

func InitAccountService(db *sql.DB) *AccountService {
	return &AccountService{AccountRepo: repositories.InitAccountRepositories(db)}
}

func (s *AccountService) Create(Account models.Account) (int, error) {
	if Account.Username == "" || Account.Password == "" || Account.Email == "" {
		return -1, fmt.Errorf(" Erreur ajout account - Données manquantes ou invalides")
	}

	AccountId, prodAccountErr := s.AccountRepo.Createaccount(Account)
	if prodAccountErr != nil {
		return -1, prodAccountErr
	}

	return AccountId, nil
}

func (s *AccountService) ReadAll() ([]models.Account, error) {
	AccountsList, AccountsErr := s.AccountRepo.ReadAll()
	if AccountsErr != nil {
		return nil, AccountsErr
	}

	return AccountsList, nil
}

func (s *AccountService) ReadById(idAccount int) (models.Account, error) {
	if idAccount <= 0 {
		return models.Account{}, fmt.Errorf(" Erreur récupération account - identifiant invalide : %d", idAccount)
	}

	Account, AccountErr := s.AccountRepo.ReadById(idAccount)
	if AccountErr != nil {
		return models.Account{}, AccountErr
	}

	return Account, nil
}
