package interfaces

import (
	"context"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
)

type UseCases interface {
	// SSO cases:
	SsoService

	// Files cases:
	UploadFile(ctx context.Context, filename string, file []byte) (string, error)

	// Toys cases:
	AddToy(ctx context.Context, rawToyData entities.RawAddToyDTO) (toyID uint64, err error)
	GetAllToys(ctx context.Context) ([]entities.Toy, error)
	GetToyByID(ctx context.Context, id uint64) (*entities.Toy, error)
	GetMasterToys(ctx context.Context, masterID uint64) ([]entities.Toy, error)
	GetAllMasters(ctx context.Context) ([]entities.Master, error)
	GetMasterByID(ctx context.Context, id uint64) (*entities.Master, error)
	RegisterMaster(ctx context.Context, rawMasterData entities.RawRegisterMasterDTO) (masterID uint64, err error)
	GetAllCategories(ctx context.Context) ([]entities.Category, error)
	GetCategoryByID(ctx context.Context, id uint32) (*entities.Category, error)
	GetAllTags(ctx context.Context) ([]entities.Tag, error)
	GetTagByID(ctx context.Context, id uint32) (*entities.Tag, error)

	// Tickets cases:
	CreateTicket(ctx context.Context, rawTicketData entities.RawCreateTicketDTO) (ticketID uint64, err error)
	GetTicketByID(ctx context.Context, id uint64) (*entities.Ticket, error)
	GetAllTickets(ctx context.Context) ([]entities.Ticket, error)
	GetUserTickets(ctx context.Context, userID uint64) ([]entities.Ticket, error)
	GetMyTickets(ctx context.Context, accessToken string) ([]entities.Ticket, error)
	RespondToTicket(ctx context.Context, rawRespondData entities.RawRespondToTicketDTO) (respondID uint64, err error)
	GetRespondByID(ctx context.Context, id uint64, accessToken string) (*entities.Respond, error)
	GetTicketResponds(ctx context.Context, ticketID uint64, accessToken string) ([]entities.Respond, error)
	GetMyResponds(ctx context.Context, accessToken string) ([]entities.Respond, error)
}
