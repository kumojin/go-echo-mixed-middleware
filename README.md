# go-echo-mixed-middleware

[![Test](https://github.com/kumojin/go-echo-mixed-middleware/actions/workflows/go.yaml/badge.svg?branch=main)](https://github.com/kumojin/go-echo-mixed-middleware/actions/workflows/go.yaml)

A composition middleware for Echo that will succeed if either of the middleware succeeds.

## Guide
### Installation

```
go get github.com/kumojin/go-echo-mixed-middleware
```

### Example
```go
package main

import (
  "errors"
  "net/http"
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
  mixed "github.com/kumojin/go-echo-mixed-middleware"
)

func main() {
  // Echo instance
  e := echo.New()

  // Middleware
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())

  // Routes
  e.GET("/", hello, mixed.Mixed(failingMiddleware, succeedingMiddleware))

  // Start server
  e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
  return c.String(http.StatusOK, "Hello, World!")
}

func succeedingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return nil
		}
	}
}

func failingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
            ctx.String(http.StatusUnauthorized, "Unauthorized")
			return errors.New("Unauthorized")
		}
	}
}
```

## Contribute

**Use issues for everything**

- For a small change, just send a PR.
- For bigger changes, please open an issue for discussion before sending a PR.
- PR should have:
  - Test case
  - Documentation
  - Example (If it makes sense)
- You can also contribute by:
  - Reporting issues
  - Suggesting new features or enhancements
  - Improve/fix documentation

## License

[MIT](https://github.com/kumojin/go-echo-mixed-middleware/blob/main/LICENSE)
