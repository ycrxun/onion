package storage

import (
	"strconv"
	"sync"
)

type (
	Map map[string]Account

	memoryStorage struct {
		mtx      sync.RWMutex
		accounts Map
	}
)

func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{
		accounts: Map{},
	}
}

func (s *memoryStorage) MapToSlice() []Account {
	output := make([]Account, 0)
	for _, v := range s.accounts {
		output = append(output, v)
	}
	return output
}

func (s *memoryStorage) List(count32 int32, token string) (accounts []*Account, next string, err error) {
	count := int(count32)
	if token == "" {
		token = "0"
	}

	offset, err := strconv.Atoi(token)
	if err != nil {
		return accounts, next, err
	}

	max := len(s.accounts)

	if offset > max {
	}

	for _, v := range s.MapToSlice() {
		accounts = append(accounts, &v)
	}

	return accounts[count:offset], next, nil
}

func (s *memoryStorage) ReadByID(ID string) (account *Account, err error) {
	a := s.accounts[ID]

	return &a, nil
}

func (s *memoryStorage) ReadByEmail(email string) (*Account, error) {
	return nil, nil
}
func (s *memoryStorage) Create(a *Account, password string) error {
	return nil
}
func (s *memoryStorage) Update(a *Account) error {
	return nil
}
func (s *memoryStorage) Delete(ID string) error {
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
