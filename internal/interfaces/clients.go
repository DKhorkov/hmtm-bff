package interfaces

import (
	"github.com/DKhorkov/hmtm-notifications/api/protobuf/generated/go/notifications"
	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
	"github.com/DKhorkov/hmtm-tickets/api/protobuf/generated/go/tickets"
	"github.com/DKhorkov/hmtm-toys/api/protobuf/generated/go/toys"
)

type SsoClient interface {
	sso.AuthServiceClient
	sso.UsersServiceClient
}

type ToysClient interface {
	toys.CategoriesServiceClient
	toys.ToysServiceClient
	toys.TagsServiceClient
	toys.MastersServiceClient
}

type TicketsClient interface {
	tickets.TicketsServiceClient
	tickets.RespondsServiceClient
}

type NotificationsClient interface {
	notifications.EmailsServiceClient
}
