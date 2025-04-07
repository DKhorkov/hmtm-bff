package interfaces

//go:generate mockgen -source=services.go -destination=../../mocks/services/sso_service.go -package=mockservices -exclude_interfaces=ToysService,FileStorageService,TicketsService,NotificationsService
type SsoService interface {
	SsoRepository
}

//go:generate mockgen -source=services.go -destination=../../mocks/services/toys_service.go -package=mockservices -exclude_interfaces=SsoService,FileStorageService,TicketsService,NotificationsService
type ToysService interface {
	ToysRepository
}

//go:generate mockgen -source=services.go -destination=../../mocks/services/file_storage_service.go -package=mockservices -exclude_interfaces=ToysService,SsoService,TicketsService,NotificationsService
type FileStorageService interface {
	FileStorageRepository
}

//go:generate mockgen -source=services.go -destination=../../mocks/services/tickets_service.go -package=mockservices -exclude_interfaces=ToysService,FileStorageService,SsoService,NotificationsService
type TicketsService interface {
	TicketsRepository
}

//go:generate mockgen -source=services.go -destination=../../mocks/services/notifications_service.go -package=mockservices -exclude_interfaces=ToysService,FileStorageService,TicketsService,SsoService
type NotificationsService interface {
	NotificationsRepository
}
