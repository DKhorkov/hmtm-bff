package interfaces

import (
	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
	"github.com/DKhorkov/hmtm-toys/api/protobuf/generated/go/toys"
)

type SsoGrpcClient interface {
	sso.AuthServiceClient
	sso.UsersServiceClient
}

type ToysGrpcClient interface {
	toys.CategoriesServiceClient
	toys.ToysServiceClient
	toys.TagsServiceClient
	toys.MastersServiceClient
}
