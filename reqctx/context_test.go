package reqctx_test

import (
	"context"
	"testing"

	"github.com/uji/ness-api-function/nesserr"
	"github.com/uji/ness-api-function/reqctx"
)

func TestGetAuthenticationInfo(t *testing.T) {
	t.Run("AuthenticationInfo is registered", func(t *testing.T) {
		ainfo := reqctx.NewAuthenticationInfo("TeamID#0", "UserID#0")
		ctx := reqctx.NewRequestContext(context.Background(), ainfo)

		res, err := reqctx.GetAuthenticationInfo(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if res != ainfo {
			t.Fatal()
		}
	})

	t.Run("AuthenticationInfo is not registered", func(t *testing.T) {
		_, err := reqctx.GetAuthenticationInfo(context.Background())
		if err != nesserr.ErrUnauthorized {
			t.Fatal(err)
		}
	})
}
