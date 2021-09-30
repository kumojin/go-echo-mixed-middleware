package mixed

import (
	echo "github.com/labstack/echo/v4"
)

// Mixed returns a single middleware function composed of the two middlewares.
// Upon being called, it is internally going to call the first middleware, and if there
// is no error, return immediately with the answer.
//
// Otherwise, it is going to call the second middleware, and return their response.
func Mixed() func(handler1, handler2 echo.MiddlewareFunc) echo.MiddlewareFunc {
	return func(handler1, handler2 echo.MiddlewareFunc) echo.MiddlewareFunc {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			a1 := handler1(next)
			a2 := handler2(next)

			return func(c echo.Context) error {
				if a1(c) == nil {
					return nil
				}
				return a2(c)
			}
		}
	}
}
