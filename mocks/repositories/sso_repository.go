package mocks

import (
	"errors"
	"time"

	ssoentities "github.com/DKhorkov/hmtm-sso/pkg/entities"
)

type MockedSsoRepository struct {
	UsersStorage map[uint64]*ssoentities.User
}

func (repo *MockedSsoRepository) RegisterUser(userData ssoentities.RegisterUserDTO) (uint64, error) {
	for _, user := range repo.UsersStorage {
		if user.Email == userData.Credentials.Email {
			return 0, errors.New("user already exists")
		}
	}

	var user ssoentities.User
	user.Email = userData.Credentials.Email
	user.ID = uint64(len(repo.UsersStorage) + 1)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	repo.UsersStorage[user.ID] = &user
	return user.ID, nil
}

func (repo *MockedSsoRepository) LoginUser(userData ssoentities.LoginUserDTO) (*ssoentities.TokensDTO, error) {
	for _, user := range repo.UsersStorage {
		if user.Email == userData.Email {
			if user.Password != userData.Password {
				return nil, errors.New("invalid password")
			}

			return &ssoentities.TokensDTO{
				AccessToken:  "AccessToken",
				RefreshToken: "RefreshToken",
			}, nil
		}
	}

	return nil, errors.New("user not found")
}

func (repo *MockedSsoRepository) GetUserByID(id uint64) (*ssoentities.User, error) {
	user := repo.UsersStorage[id]
	if user != nil {
		return user, nil
	}

	return nil, errors.New("user not found")
}

func (repo *MockedSsoRepository) GetAllUsers() ([]*ssoentities.User, error) {
	var users []*ssoentities.User
	for _, user := range repo.UsersStorage {
		users = append(users, user)
	}

	return users, nil
}

func (repo *MockedSsoRepository) RefreshTokens(
	refreshTokensData ssoentities.TokensDTO,
) (*ssoentities.TokensDTO, error) {
	return &ssoentities.TokensDTO{
		AccessToken:  refreshTokensData.AccessToken,
		RefreshToken: refreshTokensData.RefreshToken,
	}, nil
}

func (repo *MockedSsoRepository) GetMe(accessToken string) (*ssoentities.User, error) {
	if len(repo.UsersStorage) == 0 {
		return nil, errors.New("user not found for " + accessToken)
	}

	return repo.UsersStorage[0], nil
}
