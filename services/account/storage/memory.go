package storage

import (
	"sync"
	"github.com/twinj/uuid"
)

type (
	Map map[string]Account

	memoryStorage struct {
		mtx      sync.RWMutex
		accounts Map
	}
)

func NewMemoryStorage() *memoryStorage {
	accounts := make(Map, 0)
	account := Account{
		ID:    uuid.NewV4().String(),
		Name:  "ycrxun",
		Email: "ycrxun@163.com",
	}
	accounts[account.ID] = account
	return &memoryStorage{
		accounts: accounts,
	}
}

func (s *memoryStorage) List(count32 int32, token string) (accounts []*Account, next string, err error) {

	for _, v := range s.accounts {
		accounts = append(accounts, &v)
	}

	return accounts, next, nil
}

func (s *memoryStorage) ReadByID(ID string) (account *Account, err error) {
	a := s.accounts[ID]

	return &a, nil
}

func (s *memoryStorage) ReadByEmail(email string) (*Account, error) {
	for _, v := range s.accounts {
		if v.Email == email {
			return &v, nil
		}
	}
	return nil, ErrAccountNotFound
}
func (s *memoryStorage) Create(a *Account, password string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	for _, v := range s.accounts {
		if v.Email == a.Email {
			return ErrEmailExists
		}
	}
	a.ID = uuid.NewV4().String()
	s.accounts[a.ID] = *a
	return nil
}
func (s *memoryStorage) Update(a *Account) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.accounts[a.ID] = *a
	return nil
}
func (s *memoryStorage) Delete(ID string) error {
	for k, v := range s.accounts {
		if v.ID == ID {
			delete(s.accounts, k)
		}
	}
	return nil
}
func (s *memoryStorage) Confirm(token string) (*Account, error) {
	return nil, nil
}
func (s *memoryStorage) GeneratePasswordToken(email string) (*Account, error) {
	return nil, nil
}
func (s *memoryStorage) UpdatePassword(string, string) (*Account, error) {
	return nil, nil
}
func (s *memoryStorage) Migrate() error {
	return nil
}
func (s *memoryStorage) Truncate() error {
	return nil
}
func (s *memoryStorage) Close() error {
	return nil
}
