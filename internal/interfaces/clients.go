package interfaces

import "github.com/DKhorkov/hmtm-sso/protobuf/generated/go/sso"

type SsoGrpcClient interface {
	sso.AuthServiceClient
	sso.UsersServiceClient
}
