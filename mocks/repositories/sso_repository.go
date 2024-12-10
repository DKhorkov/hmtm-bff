package mocks

import (
	"errors"
	"time"

	"github.com/DKhorkov/hmtm-bff/internal/models"
)

type MockedSsoRepository struct {
	UsersStorage map[uint64]*models.User
}

func (repo *MockedSsoRepository) RegisterUser(userData models.RegisterUserDTO) (uint64, error) {
	for _, user := range repo.UsersStorage {
		if user.Email == userData.Email {
			return 0, errors.New("user already exists")
		}
	}

	var user models.User
	user.Email = userData.Email
	user.ID = uint64(len(repo.UsersStorage) + 1)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	repo.UsersStorage[user.ID] = &user
	return user.ID, nil
}

func (repo *MockedSsoRepository) LoginUser(userData models.LoginUserDTO) (*models.TokensDTO, error) {
	for _, user := range repo.UsersStorage {
		if user.Email == userData.Email {
			if user.Password != userData.Password {
				return nil, errors.New("invalid password")
			}

			return &models.TokensDTO{
				AccessToken:  "AccessToken",
				RefreshToken: "RefreshToken",
			}, nil
		}
	}

	return nil, errors.New("user not found")
}

func (repo *MockedSsoRepository) GetUserByID(id uint64) (*models.User, error) {
	user := repo.UsersStorage[id]
	if user != nil {
		return user, nil
	}

	return nil, errors.New("user not found")
}

func (repo *MockedSsoRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	for _, user := range repo.UsersStorage {
		users = append(users, *user)
	}

	return users, nil
}

func (repo *MockedSsoRepository) RefreshTokens(
	refreshTokensData models.TokensDTO,
) (*models.TokensDTO, error) {
	return &models.TokensDTO{
		AccessToken:  refreshTokensData.AccessToken,
		RefreshToken: refreshTokensData.RefreshToken,
	}, nil
}

func (repo *MockedSsoRepository) GetMe(accessToken string) (*models.User, error) {
	if len(repo.UsersStorage) == 0 {
		return nil, errors.New("user not found for " + accessToken)
	}

	return repo.UsersStorage[0], nil
}
