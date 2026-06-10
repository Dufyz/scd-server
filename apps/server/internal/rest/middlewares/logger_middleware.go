package middlewares

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func LoggerMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now().Local()

			err := next(c)

			req := c.Request()
			res := c.Response()

			fields := []zap.Field{
				zap.String("latency", time.Since(start).String()),
				zap.String("host", req.Host),
				zap.String("request", req.Method+" "+req.RequestURI),
				zap.Int("status", res.Status),
				zap.Int64("size", res.Size),
				zap.String("user_agent", req.UserAgent()),
			}

			n := res.Status
			switch {
			case n >= 500:
				logger.Error("Server error", fields...)
			case n >= 400:
				logger.Warn("Client error", fields...)
			case n >= 300:
				logger.Info("Redirection", fields...)
			default:
				logger.Info("Success", fields...)
			}

			return err
		}
	}
}
