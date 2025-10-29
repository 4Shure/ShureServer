package routes

import (
	"4shure/cmd/internal/service"
	"4shure/cmd/internal/utils"
	"4shure/cmd/internal/utils/apierror"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type UserService interface {
	GetUsers() ([]*service.UserResponse, apierror.ErrorResponse)
	GetUser(rawId, subId string) (*service.UserResponse, apierror.ErrorResponse)
	CreateUser(req *service.CreateUserRequest) apierror.ErrorResponse
	Login(req *service.UserLoginRequest) (*service.UserLoginResponse, apierror.ErrorResponse)
	ConfirmSignup(req *service.ConfirmSignupRequest) apierror.ErrorResponse
}

type DefaultUserRoute struct {
	UserService UserService
}

func NewUserDefault(userService UserService) *DefaultUserRoute {
	return &DefaultUserRoute{UserService: userService}
}

func (u *DefaultUserRoute) GetUsers(c echo.Context) error {
	users, apierr := u.UserService.GetUsers()
	if apierr != nil {
		return c.JSON(apierr.Code(), apierr)
	}

	resp := echo.Map{"users": users}
	return c.JSON(http.StatusOK, &resp)
}

func (u *DefaultUserRoute) GetUser(c echo.Context) error {
	rawId := strings.TrimSpace(c.Param("id"))
	if rawId == "" {
		return c.JSON(http.StatusBadRequest, apierror.NewMissingParamError("id"))
	}

	data, err := utils.ParseTokenDataCtx(c)
	if err != nil {
		return c.JSON(401, apierror.InvalidAuthTokenError)
	}

	user, apierr := u.UserService.GetUser(rawId, data.Sub)
	if apierr != nil {
		return c.JSON(apierr.Code(), apierr)
	}
	return c.JSON(http.StatusOK, user)
}

func (u *DefaultUserRoute) CreateUser(c echo.Context) error {
	var req service.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, apierror.MalformedBodyError)
	}

	err := u.UserService.CreateUser(&req)
	if err != nil {
		return c.JSON(err.Code(), err)
	}
	return c.NoContent(http.StatusCreated)
}

func (u *DefaultUserRoute) CreateLogin(c echo.Context) error {
	var req service.UserLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, apierror.MalformedBodyError)
	}

	resp, apierr := u.UserService.Login(&req)
	if apierr != nil {
		return c.JSON(apierr.Code(), apierr)
	}
	return c.JSON(http.StatusOK, resp)
}

func (u *DefaultUserRoute) VerifySignup(c echo.Context) error {
	var req service.ConfirmSignupRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, apierror.MalformedBodyError)
	}

	apierr := u.UserService.ConfirmSignup(&req)
	if apierr != nil {
		return c.JSON(apierr.Code(), apierr)
	}
	return c.NoContent(http.StatusOK)
}
