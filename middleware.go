package mixed

import (
	echo "github.com/labstack/echo/v4"
)

// Mixed returns a single middleware function composed of the middlewares.
// Upon being called, it is internally going to call the first middleware, and if there
// is no error, return immediately with the answer.
//
// Otherwise, it is going to call the next middlewares, and return their response.
func Mixed(preserveKeys []string) func(handlers ...echo.MiddlewareFunc) echo.MiddlewareFunc {
	return func(handlers ...echo.MiddlewareFunc) echo.MiddlewareFunc {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			var tempContext echo.Context

			return func(c echo.Context) error {
				var err error

				for _, handler := range handlers {
					tempContext = copyContext(c, preserveKeys)
					defer tempContext.Echo().ReleaseContext(tempContext)

					err = handler(next)(tempContext)
					if err == nil {
						copyResponse(c, tempContext, preserveKeys)

						return nil
					}
				}

				if err != nil {
					copyResponse(c, tempContext, preserveKeys)
					c.Response().WriteHeader(tempContext.Response().Status)

					return err
				}

				return nil
			}
		}
	}
}
