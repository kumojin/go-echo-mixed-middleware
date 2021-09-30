package mixed

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	echo "github.com/labstack/echo/v4"
)

func TestMixed(t *testing.T) {
	type args struct {
		mw1 echo.MiddlewareFunc
		mw2 echo.MiddlewareFunc
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
				mw1: succeedingMiddleware("hello"),
				mw2: failingMiddleware(errors.New("boom")),
			},
			want: "hello",
		},
		{
			name: "TestMixedAuthenticationSucceedsOnSecondAuth",
			args: args{
				mw1: failingMiddleware(errors.New("boom")),
				mw2: succeedingMiddleware("hello2"),
			},
			want: "hello2",
		},
		{
			name: "TestMixedAuthenticationFailsOnBothFailedAuths",
			args: args{
				mw1: failingMiddleware(errors.New("boom")),
				mw2: failingMiddleware(errors.New("boom2")),
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
				md  = Mixed()(tt.args.mw1, tt.args.mw2)

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

				got = ctx.Get("KEY").(string)
			}

			if got != tt.want {
				t.Errorf("Mixed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func succeedingMiddleware(succeedingValue string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("KEY", succeedingValue)
			return nil
		}
	}
}

func failingMiddleware(err error) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return err
		}
	}
}

func newContext() echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}
