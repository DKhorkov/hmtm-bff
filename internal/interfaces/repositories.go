package interfaces

import (
	"context"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
)

type SsoRepository interface {
	GetAllUsers(ctx context.Context) ([]entities.User, error)
	GetUserByID(ctx context.Context, id uint64) (*entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	RegisterUser(ctx context.Context, userData entities.RegisterUserDTO) (userID uint64, err error)
	LoginUser(ctx context.Context, userData entities.LoginUserDTO) (*entities.TokensDTO, error)
	LogoutUser(ctx context.Context, accessToken string) error
	GetMe(ctx context.Context, accessToken string) (*entities.User, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*entities.TokensDTO, error)
	VerifyUserEmail(ctx context.Context, verifyEmailToken string) error
	ForgetPassword(ctx context.Context, accessToken string) error
	ChangePassword(ctx context.Context, accessToken, oldPassword, newPassword string) error
	SendVerifyEmailMessage(ctx context.Context, email string) error
	UpdateUserProfile(ctx context.Context, userProfileData entities.UpdateUserProfileDTO) error
}

type ToysRepository interface {
	AddToy(ctx context.Context, toyData entities.AddToyDTO) (toyID uint64, err error)
	GetAllToys(ctx context.Context) ([]entities.Toy, error)
	GetToyByID(ctx context.Context, id uint64) (*entities.Toy, error)
	GetMasterToys(ctx context.Context, masterID uint64) ([]entities.Toy, error)
	GetUserToys(ctx context.Context, userID uint64) ([]entities.Toy, error)
	GetAllMasters(ctx context.Context) ([]entities.Master, error)
	GetMasterByID(ctx context.Context, id uint64) (*entities.Master, error)
	RegisterMaster(
		ctx context.Context,
		masterData entities.RegisterMasterDTO,
	) (masterID uint64, err error)
	GetAllCategories(ctx context.Context) ([]entities.Category, error)
	GetCategoryByID(ctx context.Context, id uint32) (*entities.Category, error)
	GetAllTags(ctx context.Context) ([]entities.Tag, error)
	GetTagByID(ctx context.Context, id uint32) (*entities.Tag, error)
	CreateTags(ctx context.Context, tagsData []entities.CreateTagDTO) ([]uint32, error)
	UpdateToy(ctx context.Context, toyData entities.UpdateToyDTO) error
	DeleteToy(ctx context.Context, id uint64) error
	GetMasterByUser(ctx context.Context, userID uint64) (*entities.Master, error)
	UpdateMaster(ctx context.Context, masterData entities.UpdateMasterDTO) error
}

type FileStorageRepository interface {
	Upload(ctx context.Context, key string, file []byte) (string, error)
	Delete(ctx context.Context, key string) error
	DeleteMany(ctx context.Context, keys []string) []error
}

type TicketsRepository interface {
	CreateTicket(
		ctx context.Context,
		ticketData entities.CreateTicketDTO,
	) (ticketID uint64, err error)
	GetTicketByID(ctx context.Context, id uint64) (*entities.RawTicket, error)
	GetAllTickets(ctx context.Context) ([]entities.RawTicket, error)
	GetUserTickets(ctx context.Context, userID uint64) ([]entities.RawTicket, error)
	RespondToTicket(
		ctx context.Context,
		respondData entities.RespondToTicketDTO,
	) (respondID uint64, err error)
	GetRespondByID(ctx context.Context, id uint64) (*entities.Respond, error)
	GetTicketResponds(ctx context.Context, ticketID uint64) ([]entities.Respond, error)
	GetUserResponds(ctx context.Context, userID uint64) ([]entities.Respond, error)
	UpdateRespond(ctx context.Context, respondData entities.UpdateRespondDTO) error
	DeleteRespond(ctx context.Context, id uint64) error
	UpdateTicket(ctx context.Context, ticketData entities.UpdateTicketDTO) error
	DeleteTicket(ctx context.Context, id uint64) error
}

type NotificationsRepository interface {
	GetUserEmailCommunications(ctx context.Context, userID uint64) ([]entities.Email, error)
}
