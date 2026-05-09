package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/req"
	"github.com/xtsank/mypills-super-service/src/internal/transport/middleware"

	"github.com/xtsank/mypills-super-service/src/internal/service"
	_ "github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
)

type AuthHandler struct {
	authService service.IAuthService
}

func NewAuthHandler(i do.Injector) (*AuthHandler, error) {
	authService := do.MustInvoke[service.IAuthService](i)

	return &AuthHandler{authService: authService}, nil
}

func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
}

// Register godoc
// @Summary      Регистрация нового пользователя
// @Description  Создает аккаунт и сохраняет профиль пользователя
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input  body      req.CreateUserDto  true  "Данные пользователя"
// @Success      201    {object}  res.AuthResDto     "Пользователь зарегистрирован"
// @Failure      400    {object}  errors.AppError    "Невалидные входные данные"
// @Failure      409    {object}  errors.AppError    "Пользователь уже существует"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input req.CreateUserDto

	err := c.BindJSON(&input)
	if err != nil {
		_ = c.Error(errors.ErrInvalidInput.WithError(err))
		return
	}

	cmd := command.NewCreateUserCmd(
		input.Login,
		input.Password,
		input.Sex,
		input.Weight,
		input.Age,
		input.IsPregnant,
		input.IsDriver,
		input.Illnesses,
		input.Allergies,
	)

	result, err := h.authService.Register(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Set(middleware.ResponsePayloadKey, result)
	c.Set(middleware.ResponseStatusKey, http.StatusCreated)
}

// Login godoc
// @Summary      Авторизация пользователя
// @Description  Проверяет логин и пароль, возвращает токен
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input  body      req.LoginUserDto  true  "Данные для входа"
// @Success      200    {object}  res.AuthResDto    "Успешный вход"
// @Failure      400    {object}  errors.AppError   "Невалидные входные данные"
// @Failure      401    {object}  errors.AppError   "Неверные учетные данные"
// @Failure      404    {object}  errors.AppError   "Пользователь не найден"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input req.LoginUserDto

	err := c.BindJSON(&input)
	if err != nil {
		_ = c.Error(errors.ErrInvalidInput.WithError(err))
		return
	}

	cmd, err := command.NewLoginUserCmd(
		input.Login,
		input.Password,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	result, err := h.authService.Login(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Set(middleware.ResponsePayloadKey, result)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}
