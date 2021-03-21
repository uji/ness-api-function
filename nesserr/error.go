package nesserr

import "golang.org/x/xerrors"

var (
	ErrUnauthorized = xerrors.New("Unauthorized")
)
