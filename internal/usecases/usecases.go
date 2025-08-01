package usecases

import (
	"context"
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/security"
	"github.com/DKhorkov/libs/tracing"

	"github.com/DKhorkov/hmtm-bff/internal/config"
	"github.com/DKhorkov/hmtm-bff/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-bff/internal/errors"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

const (
	tagsLimit              = 10
	attachmentsLimit       = 5
	searchQueryLengthLimit = 50
)

func New(
	ssoService interfaces.SsoService,
	toysService interfaces.ToysService,
	fileStorageService interfaces.FileStorageService,
	ticketsService interfaces.TicketsService,
	notificationsService interfaces.NotificationsService,
	validationConfig config.ValidationConfig,
	logger logging.Logger,
	traceProvider tracing.Provider,
) *UseCases {
	return &UseCases{
		ssoService:           ssoService,
		toysService:          toysService,
		fileStorageService:   fileStorageService,
		ticketsService:       ticketsService,
		notificationsService: notificationsService,
		validationConfig:     validationConfig,
		logger:               logger,
		traceProvider:        traceProvider,
	}
}

type UseCases struct {
	ssoService           interfaces.SsoService
	toysService          interfaces.ToysService
	fileStorageService   interfaces.FileStorageService
	ticketsService       interfaces.TicketsService
	notificationsService interfaces.NotificationsService
	validationConfig     config.ValidationConfig
	logger               logging.Logger
	traceProvider        tracing.Provider
}

func (useCases *UseCases) RegisterUser(
	ctx context.Context,
	userData entities.RegisterUserDTO,
) (uint64, error) {
	return useCases.ssoService.RegisterUser(ctx, userData)
}

func (useCases *UseCases) LoginUser(
	ctx context.Context,
	userData entities.LoginUserDTO,
) (*entities.TokensDTO, error) {
	return useCases.ssoService.LoginUser(ctx, userData)
}

func (useCases *UseCases) LogoutUser(ctx context.Context, accessToken string) error {
	return useCases.ssoService.LogoutUser(ctx, accessToken)
}

func (useCases *UseCases) VerifyUserEmail(ctx context.Context, verifyEmailToken string) error {
	return useCases.ssoService.VerifyUserEmail(ctx, verifyEmailToken)
}

func (useCases *UseCases) SendVerifyEmailMessage(ctx context.Context, email string) error {
	return useCases.ssoService.SendVerifyEmailMessage(ctx, email)
}

func (useCases *UseCases) SendForgetPasswordMessage(ctx context.Context, email string) error {
	return useCases.ssoService.SendForgetPasswordMessage(ctx, email)
}

func (useCases *UseCases) ChangePassword(
	ctx context.Context,
	accessToken string,
	oldPassword string,
	newPassword string,
) error {
	return useCases.ssoService.ChangePassword(ctx, accessToken, oldPassword, newPassword)
}

func (useCases *UseCases) ForgetPassword(ctx context.Context, forgetPasswordToken, newPassword string) error {
	return useCases.ssoService.ForgetPassword(ctx, forgetPasswordToken, newPassword)
}

func (useCases *UseCases) GetMe(ctx context.Context, accessToken string) (*entities.User, error) {
	return useCases.ssoService.GetMe(ctx, accessToken)
}

func (useCases *UseCases) RefreshTokens(
	ctx context.Context,
	refreshToken string,
) (*entities.TokensDTO, error) {
	return useCases.ssoService.RefreshTokens(ctx, refreshToken)
}

func (useCases *UseCases) GetUserByID(ctx context.Context, id uint64) (*entities.User, error) {
	return useCases.ssoService.GetUserByID(ctx, id)
}

func (useCases *UseCases) GetUserByEmail(
	ctx context.Context,
	email string,
) (*entities.User, error) {
	return useCases.ssoService.GetUserByEmail(ctx, email)
}

func (useCases *UseCases) GetUsers(ctx context.Context, pagination *entities.Pagination) ([]entities.User, error) {
	return useCases.ssoService.GetUsers(ctx, pagination)
}

func (useCases *UseCases) AddToy(
	ctx context.Context,
	rawToyData entities.RawAddToyDTO,
) (uint64, error) {
	if len(rawToyData.Tags) > tagsLimit {
		return 0, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too many Tags. Limit is %d", tagsLimit),
		}
	}

	if len(rawToyData.Attachments) > attachmentsLimit {
		return 0, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too many Attachments. Limit is %d", attachmentsLimit),
		}
	}

	user, err := useCases.GetMe(ctx, rawToyData.AccessToken)
	if err != nil {
		return 0, err
	}

	uploadedFiles, err := useCases.UploadFiles(ctx, user.ID, rawToyData.Attachments)
	if err != nil {
		return 0, err
	}

	tagIDs, err := useCases.toysService.CreateTags(ctx, useCases.prepareTagsForCreation(rawToyData.Tags))
	if err != nil {
		return 0, err
	}

	toyData := entities.AddToyDTO{
		UserID:      user.ID,
		CategoryID:  rawToyData.CategoryID,
		Name:        rawToyData.Name,
		Description: rawToyData.Description,
		Price:       rawToyData.Price,
		Quantity:    rawToyData.Quantity,
		TagIDs:      tagIDs,
		Attachments: uploadedFiles,
	}

	return useCases.toysService.AddToy(ctx, toyData)
}

func (useCases *UseCases) GetToys(
	ctx context.Context,
	pagination *entities.Pagination,
	filters *entities.ToysFilters,
) ([]entities.Toy, error) {
	if filters != nil && filters.Search != nil && len(*filters.Search) > searchQueryLengthLimit {
		return nil, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too long search query. Limit is %d symbols", searchQueryLengthLimit),
		}
	}

	if filters != nil && len(filters.TagIDs) > tagsLimit {
		return nil, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too many Tags. Limit is %d", tagsLimit),
		}
	}

	return useCases.toysService.GetToys(ctx, pagination, filters)
}

func (useCases *UseCases) CountToys(ctx context.Context, filters *entities.ToysFilters) (uint64, error) {
	if filters != nil && filters.Search != nil && len(*filters.Search) > searchQueryLengthLimit {
		return 0, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too long search query. Limit is %d symbols", searchQueryLengthLimit),
		}
	}

	if filters != nil && len(filters.TagIDs) > tagsLimit {
		return 0, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too many Tags. Limit is %d", tagsLimit),
		}
	}

	return useCases.toysService.CountToys(ctx, filters)
}

func (useCases *UseCases) CountMasterToys(
	ctx context.Context,
	masterID uint64,
	filters *entities.ToysFilters,
) (uint64, error) {
	if filters != nil && filters.Search != nil && len(*filters.Search) > searchQueryLengthLimit {
		return 0, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too long search query. Limit is %d symbols", searchQueryLengthLimit),
		}
	}

	if filters != nil && len(filters.TagIDs) > tagsLimit {
		return 0, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too many Tags. Limit is %d", tagsLimit),
		}
	}

	return useCases.toysService.CountMasterToys(ctx, masterID, filters)
}

func (useCases *UseCases) CountMyToys(
	ctx context.Context,
	accessToken string,
	filters *entities.ToysFilters,
) (uint64, error) {
	if filters != nil && filters.Search != nil && len(*filters.Search) > searchQueryLengthLimit {
		return 0, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too long search query. Limit is %d symbols", searchQueryLengthLimit),
		}
	}

	if filters != nil && len(filters.TagIDs) > tagsLimit {
		return 0, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too many Tags. Limit is %d", tagsLimit),
		}
	}

	user, err := useCases.GetMe(ctx, accessToken)
	if err != nil {
		return 0, err
	}

	return useCases.toysService.CountUserToys(ctx, user.ID, filters)
}

func (useCases *UseCases) GetMasterToys(
	ctx context.Context,
	masterID uint64,
	pagination *entities.Pagination,
	filters *entities.ToysFilters,
) ([]entities.Toy, error) {
	if filters != nil && filters.Search != nil && len(*filters.Search) > searchQueryLengthLimit {
		return nil, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too long search query. Limit is %d symbols", searchQueryLengthLimit),
		}
	}

	if filters != nil && len(filters.TagIDs) > tagsLimit {
		return nil, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too many Tags. Limit is %d", tagsLimit),
		}
	}

	return useCases.toysService.GetMasterToys(ctx, masterID, pagination, filters)
}

func (useCases *UseCases) GetMyToys(
	ctx context.Context,
	accessToken string,
	pagination *entities.Pagination,
	filters *entities.ToysFilters,
) ([]entities.Toy, error) {
	if filters != nil && filters.Search != nil && len(*filters.Search) > searchQueryLengthLimit {
		return nil, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too long search query. Limit is %d symbols", searchQueryLengthLimit),
		}
	}

	if filters != nil && len(filters.TagIDs) > tagsLimit {
		return nil, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too many Tags. Limit is %d", tagsLimit),
		}
	}

	user, err := useCases.GetMe(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return useCases.toysService.GetUserToys(ctx, user.ID, pagination, filters)
}

func (useCases *UseCases) GetToyByID(ctx context.Context, id uint64) (*entities.Toy, error) {
	return useCases.toysService.GetToyByID(ctx, id)
}

func (useCases *UseCases) GetMasters(
	ctx context.Context,
	pagination *entities.Pagination,
	filters *entities.MastersFilters,
) ([]entities.Master, error) {
	return useCases.toysService.GetMasters(ctx, pagination, filters)
}

func (useCases *UseCases) CountMasters(ctx context.Context, filters *entities.MastersFilters) (uint64, error) {
	if filters != nil && filters.Search != nil && len(*filters.Search) > searchQueryLengthLimit {
		return 0, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too long search query. Limit is %d symbols", searchQueryLengthLimit),
		}
	}

	return useCases.toysService.CountMasters(ctx, filters)
}

func (useCases *UseCases) GetMasterByID(ctx context.Context, id uint64) (*entities.Master, error) {
	return useCases.toysService.GetMasterByID(ctx, id)
}

func (useCases *UseCases) GetMasterByUserID(ctx context.Context, userID uint64) (*entities.Master, error) {
	return useCases.toysService.GetMasterByUserID(ctx, userID)
}

func (useCases *UseCases) RegisterMaster(
	ctx context.Context,
	rawMasterData entities.RawRegisterMasterDTO,
) (uint64, error) {
	user, err := useCases.GetMe(ctx, rawMasterData.AccessToken)
	if err != nil {
		return 0, err
	}

	masterData := entities.RegisterMasterDTO{
		UserID: user.ID,
		Info:   rawMasterData.Info,
	}

	return useCases.toysService.RegisterMaster(ctx, masterData)
}

func (useCases *UseCases) GetAllCategories(ctx context.Context) ([]entities.Category, error) {
	return useCases.toysService.GetAllCategories(ctx)
}

func (useCases *UseCases) GetCategoryByID(ctx context.Context, id uint32) (*entities.Category, error) {
	return useCases.toysService.GetCategoryByID(ctx, id)
}

func (useCases *UseCases) GetAllTags(ctx context.Context) ([]entities.Tag, error) {
	return useCases.toysService.GetAllTags(ctx)
}

func (useCases *UseCases) GetTagByID(ctx context.Context, id uint32) (*entities.Tag, error) {
	return useCases.toysService.GetTagByID(ctx, id)
}

func (useCases *UseCases) UploadFile(
	ctx context.Context,
	userID uint64,
	file *graphql.Upload,
) (string, error) {
	filename, err := useCases.createFilename(userID, file)
	if err != nil {
		return "", err
	}

	binaryFile, err := io.ReadAll(file.File)
	if err != nil {
		return "", err
	}

	return useCases.fileStorageService.Upload(ctx, filename, binaryFile)
}

func (useCases *UseCases) UploadFiles(
	ctx context.Context,
	userID uint64,
	files []*graphql.Upload,
) ([]string, error) {
	uploadedFiles := make([]string, 0, len(files))
	uploadingErrors := make([]error, 0, len(files))

	for _, file := range files {
		filename, err := useCases.UploadFile(ctx, userID, file)
		if err != nil {
			uploadingErrors = append(uploadingErrors, err)
		} else {
			uploadedFiles = append(uploadedFiles, filename)
		}
	}

	concatenatedErrBuilder := strings.Builder{}
	concatenatedErrBuilder.WriteString("Failed to upload files:\n")

	for i, err := range uploadingErrors {
		// i + 1 due to index starts from zero
		concatenatedErrBuilder.WriteString(fmt.Sprintf("%d) %v\n", i+1, err))
	}

	if len(uploadingErrors) > 0 {
		// Logging errors for further investigate:
		logging.LogErrorContext(
			ctx,
			useCases.logger,
			concatenatedErrBuilder.String(),
			&customerrors.UploadFileError{Message: concatenatedErrBuilder.String()},
		)
	}

	if len(uploadedFiles) == 0 && len(uploadingErrors) > 0 {
		return nil, &customerrors.UploadFileError{Message: concatenatedErrBuilder.String()}
	}

	// Return no err and any amount of uploaded files, if exists:
	return uploadedFiles, nil
}

func (useCases *UseCases) CreateTicket(
	ctx context.Context,
	rawTicketData entities.RawCreateTicketDTO,
) (uint64, error) {
	if len(rawTicketData.Tags) > tagsLimit {
		return 0, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too many Tags. Limit is %d", tagsLimit),
		}
	}

	if len(rawTicketData.Attachments) > attachmentsLimit {
		return 0, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too many Attachments. Limit is %d", attachmentsLimit),
		}
	}

	user, err := useCases.GetMe(ctx, rawTicketData.AccessToken)
	if err != nil {
		return 0, err
	}

	uploadedFiles, err := useCases.UploadFiles(ctx, user.ID, rawTicketData.Attachments)
	if err != nil {
		return 0, err
	}

	tagIDs, err := useCases.toysService.CreateTags(ctx, useCases.prepareTagsForCreation(rawTicketData.Tags))
	if err != nil {
		return 0, err
	}

	ticketData := entities.CreateTicketDTO{
		UserID:      user.ID,
		CategoryID:  rawTicketData.CategoryID,
		Name:        rawTicketData.Name,
		Description: rawTicketData.Description,
		Price:       rawTicketData.Price,
		Quantity:    rawTicketData.Quantity,
		TagIDs:      tagIDs,
		Attachments: uploadedFiles,
	}

	return useCases.ticketsService.CreateTicket(ctx, ticketData)
}

func (useCases *UseCases) GetTicketByID(ctx context.Context, id uint64) (*entities.Ticket, error) {
	rawTicket, err := useCases.ticketsService.GetTicketByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Soft processing if tags were received not to have distributed monolith antipattern:
	tags, _ := useCases.GetAllTags(ctx)

	return useCases.processRawTicket(*rawTicket, tags), nil
}

func (useCases *UseCases) GetTickets(
	ctx context.Context,
	pagination *entities.Pagination,
	filters *entities.TicketsFilters,
) ([]entities.Ticket, error) {
	rawTickets, err := useCases.ticketsService.GetTickets(ctx, pagination, filters)
	if err != nil {
		return nil, err
	}

	// Soft processing if tags were received not to have distributed monolith antipattern:
	tags, _ := useCases.GetAllTags(ctx)

	tickets := make([]entities.Ticket, len(rawTickets))
	for i, rawTicket := range rawTickets {
		tickets[i] = *useCases.processRawTicket(rawTicket, tags)
	}

	return tickets, err
}

func (useCases *UseCases) GetUserTickets(
	ctx context.Context,
	userID uint64,
	pagination *entities.Pagination,
	filters *entities.TicketsFilters,
) ([]entities.Ticket, error) {
	rawTickets, err := useCases.ticketsService.GetUserTickets(ctx, userID, pagination, filters)
	if err != nil {
		return nil, err
	}

	// Soft processing if tags were received not to have distributed monolith antipattern:
	tags, _ := useCases.GetAllTags(ctx)

	tickets := make([]entities.Ticket, len(rawTickets))
	for i, rawTicket := range rawTickets {
		tickets[i] = *useCases.processRawTicket(rawTicket, tags)
	}

	return tickets, err
}

func (useCases *UseCases) GetMyTickets(
	ctx context.Context,
	accessToken string,
	pagination *entities.Pagination,
	filters *entities.TicketsFilters,
) ([]entities.Ticket, error) {
	user, err := useCases.GetMe(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return useCases.GetUserTickets(ctx, user.ID, pagination, filters)
}

func (useCases *UseCases) CountTickets(ctx context.Context, filters *entities.TicketsFilters) (uint64, error) {
	if filters != nil && filters.Search != nil && len(*filters.Search) > searchQueryLengthLimit {
		return 0, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too long search query. Limit is %d symbols", searchQueryLengthLimit),
		}
	}

	if filters != nil && len(filters.TagIDs) > tagsLimit {
		return 0, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too many Tags. Limit is %d", tagsLimit),
		}
	}

	return useCases.ticketsService.CountTickets(ctx, filters)
}

func (useCases *UseCases) CountUserTickets(
	ctx context.Context,
	userID uint64,
	filters *entities.TicketsFilters,
) (uint64, error) {
	if filters != nil && filters.Search != nil && len(*filters.Search) > searchQueryLengthLimit {
		return 0, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too long search query. Limit is %d symbols", searchQueryLengthLimit),
		}
	}

	if filters != nil && len(filters.TagIDs) > tagsLimit {
		return 0, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too many Tags. Limit is %d", tagsLimit),
		}
	}

	return useCases.ticketsService.CountUserTickets(ctx, userID, filters)
}

func (useCases *UseCases) CountMyTickets(
	ctx context.Context,
	accessToken string,
	filters *entities.TicketsFilters,
) (uint64, error) {
	if filters != nil && filters.Search != nil && len(*filters.Search) > searchQueryLengthLimit {
		return 0, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too long search query. Limit is %d symbols", searchQueryLengthLimit),
		}
	}

	if filters != nil && len(filters.TagIDs) > tagsLimit {
		return 0, &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too many Tags. Limit is %d", tagsLimit),
		}
	}

	user, err := useCases.GetMe(ctx, accessToken)
	if err != nil {
		return 0, err
	}

	return useCases.ticketsService.CountUserTickets(ctx, user.ID, filters)
}

func (useCases *UseCases) RespondToTicket(
	ctx context.Context,
	rawRespondData entities.RawRespondToTicketDTO,
) (uint64, error) {
	user, err := useCases.GetMe(ctx, rawRespondData.AccessToken)
	if err != nil {
		return 0, err
	}

	respondData := entities.RespondToTicketDTO{
		UserID:   user.ID,
		TicketID: rawRespondData.TicketID,
		Price:    rawRespondData.Price,
		Comment:  rawRespondData.Comment,
	}

	return useCases.ticketsService.RespondToTicket(ctx, respondData)
}

func (useCases *UseCases) GetRespondByID(
	ctx context.Context,
	id uint64,
	accessToken string,
) (*entities.Respond, error) {
	respond, err := useCases.ticketsService.GetRespondByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user, err := useCases.GetMe(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	ticket, err := useCases.GetTicketByID(ctx, respond.TicketID)
	if err != nil {
		return nil, err
	}

	master, err := useCases.GetMasterByID(ctx, respond.MasterID)
	if err != nil {
		return nil, err
	}

	// Check if Respond belongs to Ticket owner or to Master, which responded to Ticket:
	if ticket.UserID != user.ID && master.UserID != user.ID {
		return nil, &customerrors.PermissionDeniedError{
			Message: fmt.Sprintf(
				"User with ID=%d is not rather owner of Respond with ID=%d, or owner of Ticket with ID=%d",
				user.ID,
				id,
				ticket.ID,
			),
		}
	}

	return respond, nil
}

func (useCases *UseCases) GetTicketResponds(
	ctx context.Context,
	ticketID uint64,
	accessToken string,
) ([]entities.Respond, error) {
	user, err := useCases.GetMe(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	ticket, err := useCases.GetTicketByID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	// Check if Ticket belongs to current User:
	if ticket.UserID != user.ID {
		return nil, &customerrors.PermissionDeniedError{
			Message: fmt.Sprintf(
				"Ticket with ID=%d does not belong to current User with ID=%d",
				ticketID,
				user.ID,
			),
		}
	}

	return useCases.ticketsService.GetTicketResponds(ctx, ticketID)
}

func (useCases *UseCases) GetMyResponds(
	ctx context.Context,
	accessToken string,
) ([]entities.Respond, error) {
	user, err := useCases.GetMe(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return useCases.ticketsService.GetUserResponds(ctx, user.ID)
}

func (useCases *UseCases) GetMyEmailCommunications(
	ctx context.Context,
	accessToken string,
	pagination *entities.Pagination,
) ([]entities.Email, error) {
	user, err := useCases.GetMe(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return useCases.notificationsService.GetUserEmailCommunications(ctx, user.ID, pagination)
}

func (useCases *UseCases) CountMyEmailCommunications(ctx context.Context, accessToken string) (uint64, error) {
	user, err := useCases.GetMe(ctx, accessToken)
	if err != nil {
		return 0, err
	}

	return useCases.notificationsService.CountUserEmailCommunications(ctx, user.ID)
}

func (useCases *UseCases) UpdateUserProfile(
	ctx context.Context,
	rawUserProfileData entities.RawUpdateUserProfileDTO,
) error {
	user, err := useCases.GetMe(ctx, rawUserProfileData.AccessToken)
	if err != nil {
		return err
	}

	// Check old avatar existence and necessary to upload or delete files:
	var avatar *string

	if rawUserProfileData.Avatar != nil {
		var newAvatarFilename string

		newAvatarFilename, err = useCases.createFilename(user.ID, rawUserProfileData.Avatar)
		if err != nil {
			return err
		}

		var oldAvatarFilename string

		if user.Avatar != nil && *user.Avatar != "" {
			split := strings.Split(*user.Avatar, "/")
			oldAvatarFilename = split[len(split)-1]

			if newAvatarFilename != oldAvatarFilename {
				if err = useCases.fileStorageService.Delete(ctx, oldAvatarFilename); err != nil {
					return err
				}
			}
		}

		if newAvatarFilename != oldAvatarFilename {
			var avatarURL string

			if avatarURL, err = useCases.UploadFile(ctx, user.ID, rawUserProfileData.Avatar); err != nil {
				return err
			}

			if avatarURL != "" {
				avatar = &avatarURL
			}
		}
	}

	userProfileData := entities.UpdateUserProfileDTO{
		AccessToken: rawUserProfileData.AccessToken,
		DisplayName: rawUserProfileData.DisplayName,
		Phone:       rawUserProfileData.Phone,
		Telegram:    rawUserProfileData.Telegram,
		Avatar:      avatar,
	}

	return useCases.ssoService.UpdateUserProfile(ctx, userProfileData)
}

func (useCases *UseCases) UpdateToy(
	ctx context.Context,
	rawToyData entities.RawUpdateToyDTO,
) error {
	if len(rawToyData.Tags) > tagsLimit {
		return &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too many Tags. Limit is %d", tagsLimit),
		}
	}

	if len(rawToyData.Attachments) > attachmentsLimit {
		return &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too many Attachments. Limit is %d", attachmentsLimit),
		}
	}

	user, err := useCases.GetMe(ctx, rawToyData.AccessToken)
	if err != nil {
		return err
	}

	master, err := useCases.GetMasterByUserID(ctx, user.ID)
	if err != nil {
		return err
	}

	toy, err := useCases.GetToyByID(ctx, rawToyData.ID)
	if err != nil {
		return err
	}

	// Check if Toy belongs to User:
	if toy.MasterID != master.ID {
		return &customerrors.PermissionDeniedError{
			Message: fmt.Sprintf(
				"User with ID=%d is not owner of Toy with ID=%d",
				user.ID,
				rawToyData.ID,
			),
		}
	}

	tagIDs, err := useCases.toysService.CreateTags(ctx, useCases.prepareTagsForCreation(rawToyData.Tags))
	if err != nil {
		return err
	}

	// Old Toy Attachments set:
	oldAttachmentsSet := make(map[string]struct{}, len(toy.Attachments))

	for _, attachment := range toy.Attachments {
		split := strings.Split(attachment.Link, "/")
		oldAttachmentFilename := split[len(split)-1]
		oldAttachmentsSet[oldAttachmentFilename] = struct{}{}
	}

	// New Toy Attachments set:
	newAttachmentsSet := make(map[string]struct{}, len(rawToyData.Attachments))

	for _, attachment := range rawToyData.Attachments {
		filename, err := useCases.createFilename(user.ID, attachment)
		if err != nil {
			return err
		}

		newAttachmentsSet[filename] = struct{}{}
	}

	// Add new Attachments if it is not already exists:
	attachmentsToAdd := make([]*graphql.Upload, 0)

	for _, attachment := range rawToyData.Attachments {
		filename, err := useCases.createFilename(user.ID, attachment)
		if err != nil {
			return err
		}

		if _, ok := oldAttachmentsSet[filename]; !ok {
			attachmentsToAdd = append(attachmentsToAdd, attachment)
		}
	}

	// Still used Attachments:
	stillUsedAttachments := make([]string, 0)

	for _, attachment := range toy.Attachments {
		split := strings.Split(attachment.Link, "/")

		oldAttachmentFilename := split[len(split)-1]
		if _, ok := newAttachmentsSet[oldAttachmentFilename]; ok {
			stillUsedAttachments = append(stillUsedAttachments, attachment.Link)
		}
	}

	// Delete old Attachments if it is not used by Toy now:
	attachmentsToDelete := make([]string, 0)

	for _, attachment := range toy.Attachments {
		split := strings.Split(attachment.Link, "/")

		oldAttachmentFilename := split[len(split)-1]
		if _, ok := newAttachmentsSet[oldAttachmentFilename]; !ok {
			attachmentsToDelete = append(attachmentsToDelete, oldAttachmentFilename)
		}
	}

	if len(attachmentsToDelete) > 0 {
		deleteAttachmentErrors := useCases.fileStorageService.DeleteMany(ctx, attachmentsToDelete)
		if len(deleteAttachmentErrors) > 0 {
			concatenatedErrBuilder := strings.Builder{}
			concatenatedErrBuilder.WriteString("Failed to delete files:\n")

			for i, err := range deleteAttachmentErrors {
				// i + 1 due to index starts from zero
				concatenatedErrBuilder.WriteString(fmt.Sprintf("%d) %v\n", i+1, err))
			}

			// Logging errors for further investigate:
			logging.LogErrorContext(
				ctx,
				useCases.logger,
				concatenatedErrBuilder.String(),
				&customerrors.DeleteFileError{Message: concatenatedErrBuilder.String()},
			)
		}
	}

	var uploadedFiles []string
	if len(attachmentsToAdd) > 0 {
		uploadedFiles, err = useCases.UploadFiles(ctx, user.ID, attachmentsToAdd)
		if err != nil {
			return err
		}
	}

	updatedAttachments := append(stillUsedAttachments, uploadedFiles...)
	toyData := entities.UpdateToyDTO{
		ID:          rawToyData.ID,
		CategoryID:  rawToyData.CategoryID,
		Name:        rawToyData.Name,
		Description: rawToyData.Description,
		Price:       rawToyData.Price,
		Quantity:    rawToyData.Quantity,
		TagIDs:      tagIDs,
		Attachments: updatedAttachments,
	}

	return useCases.toysService.UpdateToy(ctx, toyData)
}

func (useCases *UseCases) DeleteToy(ctx context.Context, accessToken string, id uint64) error {
	user, err := useCases.GetMe(ctx, accessToken)
	if err != nil {
		return err
	}

	master, err := useCases.GetMasterByUserID(ctx, user.ID)
	if err != nil {
		return err
	}

	toy, err := useCases.GetToyByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if Toy belongs to User:
	if toy.MasterID != master.ID {
		return &customerrors.PermissionDeniedError{
			Message: fmt.Sprintf(
				"User with ID=%d is not owner of Toy with ID=%d",
				user.ID,
				toy.ID,
			),
		}
	}

	attachmentsToDelete := make([]string, 0, len(toy.Attachments))

	for _, attachment := range toy.Attachments {
		split := strings.Split(attachment.Link, "/")
		oldAttachmentFilename := split[len(split)-1]
		attachmentsToDelete = append(attachmentsToDelete, oldAttachmentFilename)
	}

	if len(attachmentsToDelete) > 0 {
		deleteAttachmentErrors := useCases.fileStorageService.DeleteMany(ctx, attachmentsToDelete)
		if len(deleteAttachmentErrors) > 0 {
			concatenatedErrBuilder := strings.Builder{}
			concatenatedErrBuilder.WriteString("Failed to delete files:\n")

			for i, err := range deleteAttachmentErrors {
				// i + 1 due to index starts from zero
				concatenatedErrBuilder.WriteString(fmt.Sprintf("%d) %v\n", i+1, err))
			}

			// Logging errors for further investigate:
			logging.LogErrorContext(
				ctx,
				useCases.logger,
				concatenatedErrBuilder.String(),
				&customerrors.DeleteFileError{Message: concatenatedErrBuilder.String()},
			)
		}
	}

	return useCases.toysService.DeleteToy(ctx, toy.ID)
}

func (useCases *UseCases) UpdateRespond(
	ctx context.Context,
	rawRespondData entities.RawUpdateRespondDTO,
) error {
	user, err := useCases.GetMe(ctx, rawRespondData.AccessToken)
	if err != nil {
		return err
	}

	master, err := useCases.GetMasterByUserID(ctx, user.ID)
	if err != nil {
		return err
	}

	respond, err := useCases.ticketsService.GetRespondByID(ctx, rawRespondData.ID)
	if err != nil {
		return err
	}

	// Check if Respond belongs to User:
	if respond.MasterID != master.ID {
		return &customerrors.PermissionDeniedError{
			Message: fmt.Sprintf(
				"User with ID=%d is not owner of Respond with ID=%d",
				user.ID,
				respond.ID,
			),
		}
	}

	respondData := entities.UpdateRespondDTO{
		ID:      rawRespondData.ID,
		Price:   rawRespondData.Price,
		Comment: rawRespondData.Comment,
	}

	return useCases.ticketsService.UpdateRespond(ctx, respondData)
}

func (useCases *UseCases) DeleteRespond(ctx context.Context, accessToken string, id uint64) error {
	user, err := useCases.GetMe(ctx, accessToken)
	if err != nil {
		return err
	}

	master, err := useCases.GetMasterByUserID(ctx, user.ID)
	if err != nil {
		return err
	}

	respond, err := useCases.ticketsService.GetRespondByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if Respond belongs to User:
	if respond.MasterID != master.ID {
		return &customerrors.PermissionDeniedError{
			Message: fmt.Sprintf(
				"User with ID=%d is not owner of Respond with ID=%d",
				user.ID,
				respond.ID,
			),
		}
	}

	return useCases.ticketsService.DeleteRespond(ctx, respond.ID)
}

func (useCases *UseCases) UpdateMaster(
	ctx context.Context,
	rawMasterData entities.RawUpdateMasterDTO,
) error {
	user, err := useCases.GetMe(ctx, rawMasterData.AccessToken)
	if err != nil {
		return err
	}

	master, err := useCases.GetMasterByID(ctx, rawMasterData.ID)
	if err != nil {
		return err
	}

	// Check if Master belongs to User:
	if master.UserID != user.ID {
		return &customerrors.PermissionDeniedError{
			Message: fmt.Sprintf(
				"User with ID=%d is not owner of Master with ID=%d",
				user.ID,
				master.ID,
			),
		}
	}

	masterData := entities.UpdateMasterDTO{
		ID:   rawMasterData.ID,
		Info: rawMasterData.Info,
	}

	return useCases.toysService.UpdateMaster(ctx, masterData)
}

func (useCases *UseCases) UpdateTicket(
	ctx context.Context,
	rawTicketData entities.RawUpdateTicketDTO,
) error {
	if len(rawTicketData.Tags) > tagsLimit {
		return &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too many Tags. Limit is %d", tagsLimit),
		}
	}

	if len(rawTicketData.Attachments) > attachmentsLimit {
		return &customerrors.LimitExceededError{
			Message: fmt.Sprintf("Too many Attachments. Limit is %d", attachmentsLimit),
		}
	}

	user, err := useCases.GetMe(ctx, rawTicketData.AccessToken)
	if err != nil {
		return err
	}

	ticket, err := useCases.GetTicketByID(ctx, rawTicketData.ID)
	if err != nil {
		return err
	}

	// Check if Ticket belongs to User:
	if ticket.UserID != user.ID {
		return &customerrors.PermissionDeniedError{
			Message: fmt.Sprintf(
				"User with ID=%d is not owner of Ticket with ID=%d",
				user.ID,
				ticket.ID,
			),
		}
	}

	tagIDs, err := useCases.toysService.CreateTags(ctx, useCases.prepareTagsForCreation(rawTicketData.Tags))
	if err != nil {
		return err
	}

	// Old Ticket Attachments set:
	oldAttachmentsSet := make(map[string]struct{}, len(ticket.Attachments))

	for _, attachment := range ticket.Attachments {
		split := strings.Split(attachment.Link, "/")
		oldAttachmentFilename := split[len(split)-1]
		oldAttachmentsSet[oldAttachmentFilename] = struct{}{}
	}

	// New Ticket Attachments set:
	newAttachmentsSet := make(map[string]struct{}, len(rawTicketData.Attachments))

	for _, attachment := range rawTicketData.Attachments {
		filename, err := useCases.createFilename(user.ID, attachment)
		if err != nil {
			return err
		}

		newAttachmentsSet[filename] = struct{}{}
	}

	// Add new Attachments if it is not already exists:
	attachmentsToAdd := make([]*graphql.Upload, 0)

	for _, attachment := range rawTicketData.Attachments {
		filename, err := useCases.createFilename(user.ID, attachment)
		if err != nil {
			return err
		}

		if _, ok := oldAttachmentsSet[filename]; !ok {
			attachmentsToAdd = append(attachmentsToAdd, attachment)
		}
	}

	// Still used Attachments:
	stillUsedAttachments := make([]string, 0)

	for _, attachment := range ticket.Attachments {
		split := strings.Split(attachment.Link, "/")

		oldAttachmentFilename := split[len(split)-1]
		if _, ok := newAttachmentsSet[oldAttachmentFilename]; ok {
			stillUsedAttachments = append(stillUsedAttachments, attachment.Link)
		}
	}

	// Delete old Attachments if it is not used by Ticket now:
	attachmentsToDelete := make([]string, 0)

	for _, attachment := range ticket.Attachments {
		split := strings.Split(attachment.Link, "/")

		oldAttachmentFilename := split[len(split)-1]
		if _, ok := newAttachmentsSet[oldAttachmentFilename]; !ok {
			attachmentsToDelete = append(attachmentsToDelete, oldAttachmentFilename)
		}
	}

	if len(attachmentsToDelete) > 0 {
		deleteAttachmentErrors := useCases.fileStorageService.DeleteMany(ctx, attachmentsToDelete)
		if len(deleteAttachmentErrors) > 0 {
			concatenatedErrBuilder := strings.Builder{}
			concatenatedErrBuilder.WriteString("Failed to delete files:\n")

			for i, err := range deleteAttachmentErrors {
				// i + 1 due to index starts from zero
				concatenatedErrBuilder.WriteString(fmt.Sprintf("%d) %v\n", i+1, err))
			}

			// Logging errors for further investigate:
			logging.LogErrorContext(
				ctx,
				useCases.logger,
				concatenatedErrBuilder.String(),
				&customerrors.DeleteFileError{Message: concatenatedErrBuilder.String()},
			)
		}
	}

	var uploadedFiles []string
	if len(attachmentsToAdd) > 0 {
		uploadedFiles, err = useCases.UploadFiles(ctx, user.ID, attachmentsToAdd)
		if err != nil {
			return err
		}
	}

	updatedAttachments := append(stillUsedAttachments, uploadedFiles...)
	ticketData := entities.UpdateTicketDTO{
		ID:          rawTicketData.ID,
		CategoryID:  rawTicketData.CategoryID,
		Name:        rawTicketData.Name,
		Description: rawTicketData.Description,
		Price:       rawTicketData.Price,
		Quantity:    rawTicketData.Quantity,
		TagIDs:      tagIDs,
		Attachments: updatedAttachments,
	}

	return useCases.ticketsService.UpdateTicket(ctx, ticketData)
}

func (useCases *UseCases) DeleteTicket(ctx context.Context, accessToken string, id uint64) error {
	user, err := useCases.GetMe(ctx, accessToken)
	if err != nil {
		return err
	}

	ticket, err := useCases.GetTicketByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if Ticket belongs to User:
	if ticket.UserID != user.ID {
		return &customerrors.PermissionDeniedError{
			Message: fmt.Sprintf(
				"User with ID=%d is not owner of Ticket with ID=%d",
				user.ID,
				ticket.ID,
			),
		}
	}

	attachmentsToDelete := make([]string, 0, len(ticket.Attachments))

	for _, attachment := range ticket.Attachments {
		split := strings.Split(attachment.Link, "/")
		oldAttachmentFilename := split[len(split)-1]
		attachmentsToDelete = append(attachmentsToDelete, oldAttachmentFilename)
	}

	if len(attachmentsToDelete) > 0 {
		deleteAttachmentErrors := useCases.fileStorageService.DeleteMany(ctx, attachmentsToDelete)
		if len(deleteAttachmentErrors) > 0 {
			concatenatedErrBuilder := strings.Builder{}
			concatenatedErrBuilder.WriteString("Failed to delete files:\n")

			for i, err := range deleteAttachmentErrors {
				// i + 1 due to index starts from zero
				concatenatedErrBuilder.WriteString(fmt.Sprintf("%d) %v\n", i+1, err))
			}

			// Logging errors for further investigate:
			logging.LogErrorContext(
				ctx,
				useCases.logger,
				concatenatedErrBuilder.String(),
				&customerrors.DeleteFileError{Message: concatenatedErrBuilder.String()},
			)
		}
	}

	return useCases.ticketsService.DeleteTicket(ctx, ticket.ID)
}

func (useCases *UseCases) createFilename(userID uint64, file *graphql.Upload) (string, error) {
	fileExtension := path.Ext(file.Filename)
	if !validateFileExtension(fileExtension, useCases.validationConfig.FileAllowedExtensions) {
		return "", &customerrors.InvalidFileExtensionError{Message: fileExtension}
	}

	if !validateFileSize(file.Size, useCases.validationConfig.FileMaxSize) {
		return "", &customerrors.InvalidFileSizeError{Message: strconv.FormatInt(file.Size, 10)}
	}

	filename := security.RawEncode(
		[]byte(fmt.Sprintf("%d:%s", userID, file.Filename)),
	) + fileExtension

	return filename, nil
}

func (useCases *UseCases) processRawTicket(
	ticket entities.RawTicket,
	tags []entities.Tag,
) *entities.Ticket {
	processedTags := make([]entities.Tag, len(ticket.TagIDs))
	for tagIndex := range ticket.TagIDs {
		processedTags[tagIndex] = entities.Tag{ID: ticket.TagIDs[tagIndex]}
	}

	if tags != nil { // Soft processing if tags were received not to have distributed monolith antipattern.
		tagsMap := make(map[uint32]entities.Tag)
		for _, tag := range tags {
			tagsMap[tag.ID] = tag
		}

		for i, tag := range processedTags {
			if _, ok := tagsMap[tag.ID]; ok {
				processedTags[i].Name = tagsMap[tag.ID].Name
			}
		}
	}

	return &entities.Ticket{
		ID:          ticket.ID,
		UserID:      ticket.UserID,
		CategoryID:  ticket.CategoryID,
		Name:        ticket.Name,
		Description: ticket.Description,
		Price:       ticket.Price,
		Quantity:    ticket.Quantity,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
		Tags:        processedTags,
		Attachments: ticket.Attachments,
	}
}

func (useCases *UseCases) prepareTagsForCreation(tags []string) []entities.CreateTagDTO {
	tagsData := make([]entities.CreateTagDTO, len(tags))
	for i, tag := range tags {
		tagsData[i] = entities.CreateTagDTO{Name: tag}
	}

	return tagsData
}
