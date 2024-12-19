package interfaces

import "context"

type UseCases interface {
	SsoService
	ToysService
	UploadFile(ctx context.Context, filename string, file []byte) (string, error)
}
