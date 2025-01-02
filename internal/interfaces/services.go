package interfaces

type SsoService interface {
	SsoRepository
}

type ToysService interface {
	ToysRepository
}

type FileStorageService interface {
	FileStorageRepository
}

type TicketsService interface {
	TicketsRepository
}
