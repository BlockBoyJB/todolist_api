package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "todolist_api/docs"
	"todolist_api/internal/service"
)

func NewRouter(h *echo.Echo, services *service.Services) {
	h.Use(middleware.Recover())
	h.GET("/ping", ping)
	h.GET("/swagger/*", echoSwagger.WrapHandler)

	newAuthRouter(h.Group("/auth"), services.Auth, services.User)
	auth := &authMiddleware{auth: services.Auth}

	v1 := h.Group("/api/v1", auth.authHandler)
	newTaskRouter(v1.Group("/tasks"), services.Task)
}

func ping(c echo.Context) error {
	return c.NoContent(200)
}
