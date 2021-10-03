package mixed

import (
	echo "github.com/labstack/echo/v4"
)

// Mixed returns a single middleware function composed of the two middlewares.
// Upon being called, it is internally going to call the first middleware, and if there
// is no error, return immediately with the answer.
//
// Otherwise, it is going to call the second middleware, and return their response.
func Mixed(preserveKeys []string) func(handler1, handler2 echo.MiddlewareFunc) echo.MiddlewareFunc {
	return func(handler1, handler2 echo.MiddlewareFunc) echo.MiddlewareFunc {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			a1 := handler1(next)
			a2 := handler2(next)

			return func(c echo.Context) error {
				tempContext := copyContext(c, preserveKeys)
				defer tempContext.Echo().ReleaseContext(tempContext)
				if a1(tempContext) == nil {
					copyResponse(c, tempContext, preserveKeys)

					return nil
				}

				// Try the second middleware
				err := a2(c)
				if err != nil {
					// Return the first middleware error
					copyResponse(c, tempContext, preserveKeys)
					c.Response().WriteHeader(tempContext.Response().Status)
				}
				return err
			}
		}
	}
}
