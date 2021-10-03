package mixed

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	echo "github.com/labstack/echo/v4"
)

func TestMixed(t *testing.T) {
	type args struct {
		mw1          echo.MiddlewareFunc
		mw2          echo.MiddlewareFunc
		preverveKeys []string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "TestMixedAuthenticationSucceedsOnFirstAuth",
			args: args{
				mw1:          succeedingMiddleware(204),
				mw2:          failingMiddleware(errors.New("boom")),
				preverveKeys: make([]string, 0),
			},
			want: "204",
		},
		{
			name: "TestMixedAuthenticationSucceedsOnSecondAuth",
			args: args{
				mw1:          failingMiddleware(errors.New("boom")),
				mw2:          succeedingMiddleware(204),
				preverveKeys: make([]string, 0),
			},
			want: "204",
		},
		{
			name: "TestMixedAuthenticationFailsOnBothFailedAuths",
			args: args{
				mw1:          failingMiddleware(errors.New("boom")),
				mw2:          failingMiddleware(errors.New("boom2")),
				preverveKeys: make([]string, 0),
			},
			want:    "boom2",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				got string
				ctx = newContext()
				md  = Mixed(tt.args.preverveKeys)(tt.args.mw1, tt.args.mw2)

				handlerFunc = md(nil)
				err         = handlerFunc(ctx)
			)

			switch tt.wantErr {
			case true:
				if err == nil {
					t.Errorf("got %v, want %v", err, tt.wantErr)
				}

				got = err.Error()
			default:
				if err != nil {
					t.Errorf("got %v, want %v", err, tt.wantErr)
				}

				got = fmt.Sprint(ctx.Response().Status)
			}

			if got != tt.want {
				t.Errorf("Mixed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMixedPreserveKeys(t *testing.T) {
	type args struct {
		mw1          echo.MiddlewareFunc
		mw2          echo.MiddlewareFunc
		preverveKeys []string
		ctx          echo.Context
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestMixedAuthenticationSucceedsOnFirstAuth",
			args: args{
				mw1:          succeedingMiddleware(204),
				mw2:          failingMiddleware(errors.New("boom")),
				preverveKeys: []string{"key1", "key2"},
				ctx:          newContext(),
			},
			want: "val1",
		},
		{
			name: "TestMixedAuthenticationSucceedsOnSecondAuth",
			args: args{
				mw1:          failingMiddleware(errors.New("boom")),
				mw2:          succeedingMiddleware(204),
				preverveKeys: []string{"key1", "key2"},
				ctx:          newContext(),
			},
			want: "val1",
		},
		{
			name: "TestMixedAuthenticationFailsOnBothFailedAuths",
			args: args{
				mw1:          failingMiddleware(errors.New("boom")),
				mw2:          failingMiddleware(errors.New("boom2")),
				preverveKeys: []string{"key1", "key2"},
				ctx:          newContext(),
			},
			want: "val2",
		},
		{
			name: "TestMixedAuthenticationPreserveInitalKey",
			args: args{
				mw1:          succeedingMiddleware(204),
				mw2:          failingMiddleware(errors.New("boom")),
				preverveKeys: []string{"key1", "key2", "key3"},
				ctx: func() echo.Context {
					ctx := newContext()
					ctx.Set("key3", "val3")
					return ctx
				}(),
			},
			want: "val1val3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				got = ""
				md  = Mixed(tt.args.preverveKeys)(tt.args.mw1, tt.args.mw2)
			)

			md(nil)(tt.args.ctx)
			for _, key := range tt.args.preverveKeys {
				if val := tt.args.ctx.Get(key); val != nil {
					got = fmt.Sprintf("%s%v", got, tt.args.ctx.Get(key))
				}
			}

			if got != tt.want {
				t.Errorf("Preserve keys = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMixedKeepPreserveKeyForSubMiddleware(t *testing.T) {
	var (
		mw1          = validateContextKeyMiddleware("key3", "val3", t)
		mw2          = failingMiddleware(errors.New("boom"))
		preverveKeys = []string{"key1", "key2", "key3"}
		ctx          = func() echo.Context {
			ctx := newContext()
			ctx.Set("key3", "val3")
			return ctx
		}()
		want = "val3"
		md   = Mixed(preverveKeys)(mw1, mw2)
		got  string
	)

	md(nil)(ctx)

	for _, key := range preverveKeys {
		if val := ctx.Get(key); val != nil {
			got = fmt.Sprintf("%s%v", got, ctx.Get(key))
		}
	}

	if got != want {
		t.Errorf("Preserve keys = %v, want %v", got, want)
	}
}

func TestMixedKeepPreviousContextStore(t *testing.T) {
	var (
		mw1          = succeedingMiddleware(204)
		mw2          = failingMiddleware(errors.New("boom"))
		preverveKeys = []string{}
		ctx          = func() echo.Context {
			ctx := newContext()
			ctx.Set("key", "val")
			return ctx
		}()
		want = "val"
		md   = Mixed(preverveKeys)(mw1, mw2)
		got  string
	)

	md(nil)(ctx)

	got = fmt.Sprint(ctx.Get("key"))
	if got != want {
		t.Errorf("Preserved key = %v, want %v", got, want)
	}
}

func succeedingMiddleware(status int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("key1", "val1")
			c.NoContent(status)
			return nil
		}
	}
}

func failingMiddleware(err error) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("key2", "val2")
			return err
		}
	}
}

func validateContextKeyMiddleware(key, val string, t *testing.T) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			got := fmt.Sprint(c.Get(key))
			if got != val {
				t.Errorf("ContextKey = %v, want %v", got, val)
			}
			return nil
		}
	}
}

func newContext() echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}
