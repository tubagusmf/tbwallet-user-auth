package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tubagusmf/tbwallet-user-auth/internal/model"
)

type UserHandler struct {
	userUsecase model.IUserUsecase
}

func NewUserHandler(e *echo.Echo, userUsecase model.IUserUsecase) {
	handlers := &UserHandler{
		userUsecase: userUsecase,
	}

	routeUser := e.Group("v1/user")
	routeUser.GET("", handlers.GetAll, AuthMiddleware)
	routeUser.GET("/:id", handlers.GetByID, AuthMiddleware)
	routeUser.POST("/register", handlers.Create)
	routeUser.PUT("/update/:id", handlers.Update, AuthMiddleware)
	routeUser.DELETE("/delete/:id", handlers.Delete, AuthMiddleware)
	routeUser.POST("/login", handlers.Login)
	routeUser.POST("/logout", handlers.Logout, AuthMiddleware)
}

func (h *UserHandler) GetAll(c echo.Context) error {
	var filter model.User
	filter.Name = c.QueryParam("name")
	filter.Email = c.QueryParam("email")

	users, err := h.userUsecase.GetAll(c.Request().Context(), filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch users")
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   users,
	})
}

func (h *UserHandler) GetByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	claim, ok := c.Request().Context().Value(model.BearerAuthKey).(model.CustomClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	if claim.UserID != id {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	user, err := h.userUsecase.GetByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   user,
	})
}

func (h *UserHandler) Create(c echo.Context) error {
	var body model.CreateUserInput
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if body.Name == "" || body.Email == "" || body.PasswordHash == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "All fields are required")
	}

	accessToken, err := h.userUsecase.Create(c.Request().Context(), body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	return c.JSON(http.StatusCreated, Response{
		Status:      http.StatusCreated,
		Message:     "User registered successfully",
		AccessToken: accessToken,
	})
}

func (h *UserHandler) Update(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	claim, ok := c.Request().Context().Value(model.BearerAuthKey).(model.CustomClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	if claim.UserID != id {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	var body model.UpdateUserInput
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	user, err := h.userUsecase.Update(c.Request().Context(), id, body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user")
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "User updated successfully",
		Data:    user,
	})
}

func (h *UserHandler) Delete(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	claim, ok := c.Request().Context().Value(model.BearerAuthKey).(model.CustomClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	if claim.UserID != id {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	err = h.userUsecase.Delete(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete user")
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "User deleted successfully",
	})
}

func (h *UserHandler) Login(c echo.Context) error {
	var body model.LoginInput
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if body.Email == "" || body.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "All fields are required")
	}

	accessToken, err := h.userUsecase.Login(c.Request().Context(), body)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid email or password")
	}

	return c.JSON(http.StatusOK, Response{
		Status:      http.StatusOK,
		Message:     "User logged in successfully",
		AccessToken: accessToken,
	})
}

func (h *UserHandler) Logout(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing token")
	}

	err := h.userUsecase.Logout(c.Request().Context(), model.UserSession{Token: token})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to logout")
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "Logout successful",
	})
}
