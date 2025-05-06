package credm

import (
	"context"
)

type CredClient interface {
	EncryptText(ctx context.Context, secret []byte) ([]byte, error)
	DecryptText(ctx context.Context, encryptedSecret []byte) ([]byte, error)
}
