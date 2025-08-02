package wasmtools

import (
	"context"
	"io"
)

type instance interface {
	Close(ctx context.Context) error
	Run(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer, fsMap map[string]string, args ...string) error
}

var _ instance = &Instance{}
