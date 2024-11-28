package interfaces

import "github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"

type SsoGrpcClient interface {
	sso.AuthServiceClient
	sso.UsersServiceClient
}
