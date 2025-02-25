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

func (useCases *UseCases) RegisterUser(ctx context.Context, userData entities.RegisterUserDTO) (uint64, error) {
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

func (useCases *UseCases) ChangePassword(
	ctx context.Context,
	accessToken string,
	oldPassword string,
	newPassword string,
) error {
	return useCases.ssoService.ChangePassword(ctx, accessToken, oldPassword, newPassword)
}

func (useCases *UseCases) ForgetPassword(ctx context.Context, accessToken string) error {
	return useCases.ssoService.ForgetPassword(ctx, accessToken)
}

func (useCases *UseCases) GetMe(ctx context.Context, accessToken string) (*entities.User, error) {
	return useCases.ssoService.GetMe(ctx, accessToken)
}

func (useCases *UseCases) RefreshTokens(ctx context.Context, refreshToken string) (*entities.TokensDTO, error) {
	return useCases.ssoService.RefreshTokens(ctx, refreshToken)
}

func (useCases *UseCases) GetUserByID(ctx context.Context, id uint64) (*entities.User, error) {
	return useCases.ssoService.GetUserByID(ctx, id)
}

func (useCases *UseCases) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	return useCases.ssoService.GetUserByEmail(ctx, email)
}

func (useCases *UseCases) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	return useCases.ssoService.GetAllUsers(ctx)
}

func (useCases *UseCases) AddToy(ctx context.Context, rawToyData entities.RawAddToyDTO) (uint64, error) {
	user, err := useCases.GetMe(ctx, rawToyData.AccessToken)
	if err != nil {
		return 0, err
	}

	uploadedFiles, err := useCases.UploadFiles(ctx, user.ID, rawToyData.Attachments)
	if err != nil {
		return 0, err
	}

	tagsData := make([]entities.CreateTagDTO, len(rawToyData.Tags))
	for i, tag := range rawToyData.Tags {
		tagsData[i] = entities.CreateTagDTO{
			Name: tag,
		}
	}

	tagIDs, err := useCases.toysService.CreateTags(ctx, tagsData)
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

func (useCases *UseCases) GetAllToys(ctx context.Context) ([]entities.Toy, error) {
	return useCases.toysService.GetAllToys(ctx)
}

func (useCases *UseCases) GetMasterToys(ctx context.Context, masterID uint64) ([]entities.Toy, error) {
	return useCases.toysService.GetMasterToys(ctx, masterID)
}

func (useCases *UseCases) GetMyToys(ctx context.Context, accessToken string) ([]entities.Toy, error) {
	user, err := useCases.GetMe(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return useCases.toysService.GetUserToys(ctx, user.ID)
}

func (useCases *UseCases) GetToyByID(ctx context.Context, id uint64) (*entities.Toy, error) {
	return useCases.toysService.GetToyByID(ctx, id)
}

func (useCases *UseCases) GetAllMasters(ctx context.Context) ([]entities.Master, error) {
	return useCases.toysService.GetAllMasters(ctx)
}

func (useCases *UseCases) GetMasterByID(ctx context.Context, id uint64) (*entities.Master, error) {
	return useCases.toysService.GetMasterByID(ctx, id)
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

func (useCases *UseCases) UploadFile(ctx context.Context, userID uint64, file *graphql.Upload) (string, error) {
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

func (useCases *UseCases) UploadFiles(ctx context.Context, userID uint64, files []*graphql.Upload) ([]string, error) {
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
		concatenatedErrBuilder.WriteString(fmt.Sprintf("%d) %s\n", i+1, err.Error()))
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
		_, span := useCases.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
		defer span.End()

		span.SetStatus(tracing.StatusError, concatenatedErrBuilder.String())
		return nil, &customerrors.UploadFileError{Message: concatenatedErrBuilder.String()}
	}

	// Return no err and any amount of uploaded files, if exists:
	return uploadedFiles, nil
}

func (useCases *UseCases) CreateTicket(
	ctx context.Context,
	rawTicketData entities.RawCreateTicketDTO,
) (uint64, error) {
	user, err := useCases.GetMe(ctx, rawTicketData.AccessToken)
	if err != nil {
		return 0, err
	}

	uploadedFiles, err := useCases.UploadFiles(ctx, user.ID, rawTicketData.Attachments)
	if err != nil {
		return 0, err
	}

	tagsData := make([]entities.CreateTagDTO, len(rawTicketData.Tags))
	for i, tag := range rawTicketData.Tags {
		tagsData[i] = entities.CreateTagDTO{
			Name: tag,
		}
	}

	tagIDs, err := useCases.toysService.CreateTags(ctx, tagsData)
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

	return useCases.processRawTicket(ctx, *rawTicket), nil
}

func (useCases *UseCases) processRawTicket(ctx context.Context, ticket entities.RawTicket) *entities.Ticket {
	processedTags := make([]entities.Tag, len(ticket.TagIDs))
	for tagIndex := range ticket.TagIDs {
		processedTags[tagIndex] = entities.Tag{ID: ticket.TagIDs[tagIndex]}
	}

	tags, err := useCases.toysService.GetAllTags(ctx)
	if err == nil { // Soft processing if tags were received not to have distributed monolith antipattern.
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
	}
}

func (useCases *UseCases) GetAllTickets(ctx context.Context) ([]entities.Ticket, error) {
	rawTickets, err := useCases.ticketsService.GetAllTickets(ctx)
	if err != nil {
		return nil, err
	}

	tickets := make([]entities.Ticket, len(rawTickets))
	for i, rawTicket := range rawTickets {
		tickets[i] = *useCases.processRawTicket(ctx, rawTicket)
	}

	return tickets, err
}

func (useCases *UseCases) GetUserTickets(ctx context.Context, userID uint64) ([]entities.Ticket, error) {
	rawTickets, err := useCases.ticketsService.GetUserTickets(ctx, userID)
	if err != nil {
		return nil, err
	}

	tickets := make([]entities.Ticket, len(rawTickets))
	for i, rawTicket := range rawTickets {
		tickets[i] = *useCases.processRawTicket(ctx, rawTicket)
	}

	return tickets, err
}

func (useCases *UseCases) GetMyTickets(
	ctx context.Context,
	accessToken string,
) ([]entities.Ticket, error) {
	user, err := useCases.GetMe(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return useCases.GetUserTickets(ctx, user.ID)
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

	ticket, err := useCases.ticketsService.GetTicketByID(ctx, respond.TicketID)
	if err != nil {
		return nil, err
	}

	master, err := useCases.toysService.GetMasterByID(ctx, respond.MasterID)
	if err != nil {
		return nil, err
	}

	// Check if Respond belongs to Ticket owner or to Master, which responded to Ticket.
	if ticket.UserID != user.ID && master.UserID != user.ID {
		_, span := useCases.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
		defer span.End()

		errorMessage := fmt.Sprintf(
			"User with ID=%d is not rather owner of Respond with ID=%d, or owner of Ticket with ID=%d",
			user.ID,
			id,
			ticket.ID,
		)

		span.SetStatus(tracing.StatusError, errorMessage)

		return nil, &customerrors.PermissionDeniedError{
			Message: errorMessage,
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

	ticket, err := useCases.ticketsService.GetTicketByID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	// Check if Ticket belongs to current User.
	if ticket.UserID != user.ID {
		_, span := useCases.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
		defer span.End()

		errorMessage := fmt.Sprintf(
			"Ticket with ID=%d does not belong to current User with ID=%d",
			ticketID,
			user.ID,
		)

		span.SetStatus(tracing.StatusError, errorMessage)

		return nil, &customerrors.PermissionDeniedError{
			Message: errorMessage,
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
) ([]entities.Email, error) {
	user, err := useCases.GetMe(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return useCases.notificationsService.GetUserEmailCommunications(ctx, user.ID)
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
	var avatar string
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
			if avatar, err = useCases.UploadFile(ctx, user.ID, rawUserProfileData.Avatar); err != nil {
				return err
			}
		}
	}

	userProfileData := entities.UpdateUserProfileDTO{
		AccessToken: rawUserProfileData.AccessToken,
		DisplayName: rawUserProfileData.DisplayName,
		Phone:       rawUserProfileData.Phone,
		Telegram:    rawUserProfileData.Telegram,
		Avatar:      &avatar,
	}

	return useCases.ssoService.UpdateUserProfile(ctx, userProfileData)
}

func (useCases *UseCases) createFilename(userID uint64, file *graphql.Upload) (string, error) {
	fileExtension := path.Ext(file.Filename)
	if !validateFileExtension(fileExtension, useCases.validationConfig.FileAllowedExtensions) {
		return "", &customerrors.InvalidFileExtensionError{Message: fileExtension}
	}

	if !validateFileSize(file.Size, useCases.validationConfig.FileMaxSize) {
		return "", &customerrors.InvalidFileSizeError{Message: strconv.FormatInt(file.Size, 10)}
	}

	filename := security.RawEncode([]byte(fmt.Sprintf("%d:%s", userID, file.Filename))) + fileExtension
	return filename, nil
}
