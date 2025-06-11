package services

import (
	"database/sql"
	"fmt"
	"forum/hash"
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
		return -1, fmt.Errorf("Erreur ajout account - Données manquantes ou invalides")
	}

	if len(Account.Password) < 8 && len(Account.Password) > 30 && len(Account.Username) < 3 && len(Account.Username) > 30 {
		return -1, fmt.Errorf("Erreur ajout account - Longueur des champs invalide")
	}

	_, usernameErr := s.AccountRepo.ReadByUsername(Account.Username)
	if usernameErr == nil {
		return -1, fmt.Errorf("Erreur ajout account - Nom d'utilisateur déjà utilisé")
	}

	_, emailErr := s.AccountRepo.ReadByEmail(Account.Email)
	if emailErr == nil {
		return -1, fmt.Errorf("Erreur ajout account - Email déjà utilisé")
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

func (s *AccountService) Login(email, password string) (models.Account, error) {
	if email == "" || password == "" {
		return models.Account{}, fmt.Errorf("email et mot de passe requis")
	}

	account, err := s.AccountRepo.ReadByEmail(email)
	if err != nil {
		return models.Account{}, fmt.Errorf("identifiants invalides-email", err)
	}

	if !hash.ComparePasswords(account.Password, password) {
		return models.Account{}, fmt.Errorf("identifiants invalides-password")
	}

	return account, nil
}
