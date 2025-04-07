package interfaces

import (
	"github.com/DKhorkov/hmtm-notifications/api/protobuf/generated/go/notifications"
	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
	"github.com/DKhorkov/hmtm-tickets/api/protobuf/generated/go/tickets"
	"github.com/DKhorkov/hmtm-toys/api/protobuf/generated/go/toys"
)

//go:generate mockgen -source=clients.go -destination=../../mocks/clients/sso_client.go -package=mockclients -exclude_interfaces=ToysClient,TicketsClient,NotificationsClient
type SsoClient interface {
	sso.AuthServiceClient
	sso.UsersServiceClient
}

//go:generate mockgen -source=clients.go -destination=../../mocks/clients/toys_client.go -package=mockclients -exclude_interfaces=SsoClient,TicketsClient,NotificationsClient
type ToysClient interface {
	toys.CategoriesServiceClient
	toys.ToysServiceClient
	toys.TagsServiceClient
	toys.MastersServiceClient
}

//go:generate mockgen -source=clients.go -destination=../../mocks/clients/tickets_client.go -package=mockclients -exclude_interfaces=ToysClient,SsoClient,NotificationsClient
type TicketsClient interface {
	tickets.TicketsServiceClient
	tickets.RespondsServiceClient
}

//go:generate mockgen -source=clients.go -destination=../../mocks/clients/notifications_client.go -package=mockclients -exclude_interfaces=ToysClient,TicketsClient,SsoClient
type NotificationsClient interface {
	notifications.EmailsServiceClient
}
