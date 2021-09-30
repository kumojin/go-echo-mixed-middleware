package mixed

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	echo "github.com/labstack/echo/v4"

	"github.com/stretchr/testify/assert"
)

func TestMixedAuthenticationSucceedsOnFirstAuth(t *testing.T) {
	auth1 := succeedingMiddleware("hello")
	auth2 := failingMiddleware(errors.New("boom"))

	ctx := newContext()
	md := Mixed()(auth1, auth2)

	handlerFunc := md(nil)
	err := handlerFunc(ctx)

	assert.NoError(t, err)
	assert.Equal(t, "hello", ctx.Get("KEY"))
}

func TestMixedAuthenticationSucceedsOnSecondAuth(t *testing.T) {
	auth1 := failingMiddleware(errors.New("boom"))
	auth2 := succeedingMiddleware("hello2")

	ctx := newContext()
	md := Mixed()(auth1, auth2)

	handlerFunc := md(nil)
	err := handlerFunc(ctx)

	assert.NoError(t, err)
	assert.Equal(t, "hello2", ctx.Get("KEY"))
}

func TestMixedAuthenticationFailsOnBothFailedAuths(t *testing.T) {
	auth1 := failingMiddleware(errors.New("boom"))
	auth2 := failingMiddleware(errors.New("boom2"))

	ctx := newContext()
	md := Mixed()(auth1, auth2)

	handlerFunc := md(nil)
	err := handlerFunc(ctx)

	assert.Equal(t, "boom2", err.Error())
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
