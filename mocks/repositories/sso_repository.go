package mocks

import (
	"errors"
	"time"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
)

type MockedSsoRepository struct {
	UsersStorage map[uint64]*entities.User
}

func (repo *MockedSsoRepository) RegisterUser(userData entities.RegisterUserDTO) (uint64, error) {
	for _, user := range repo.UsersStorage {
		if user.Email == userData.Email {
			return 0, errors.New("user already exists")
		}
	}

	var user entities.User
	user.Email = userData.Email
	user.ID = uint64(len(repo.UsersStorage) + 1)
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()

	repo.UsersStorage[user.ID] = &user
	return user.ID, nil
}

func (repo *MockedSsoRepository) LoginUser(userData entities.LoginUserDTO) (*entities.TokensDTO, error) {
	for _, user := range repo.UsersStorage {
		if user.Email == userData.Email {
			if user.Password != userData.Password {
				return nil, errors.New("invalid password")
			}

			return &entities.TokensDTO{
				AccessToken:  "AccessToken",
				RefreshToken: "RefreshToken",
			}, nil
		}
	}

	return nil, errors.New("user not found")
}

func (repo *MockedSsoRepository) GetUserByID(id uint64) (*entities.User, error) {
	user := repo.UsersStorage[id]
	if user != nil {
		return user, nil
	}

	return nil, errors.New("user not found")
}

func (repo *MockedSsoRepository) GetAllUsers() ([]entities.User, error) {
	var users []entities.User
	for _, user := range repo.UsersStorage {
		users = append(users, *user)
	}

	return users, nil
}

func (repo *MockedSsoRepository) RefreshTokens(
	refreshTokensData entities.TokensDTO,
) (*entities.TokensDTO, error) {
	return &entities.TokensDTO{
		AccessToken:  refreshTokensData.AccessToken,
		RefreshToken: refreshTokensData.RefreshToken,
	}, nil
}

func (repo *MockedSsoRepository) GetMe(accessToken string) (*entities.User, error) {
	if len(repo.UsersStorage) == 0 {
		return nil, errors.New("user not found for " + accessToken)
	}

	return repo.UsersStorage[0], nil
}
