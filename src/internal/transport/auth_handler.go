package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xtsank/mypills-super-service/internal/service/command"
	"github.com/xtsank/mypills-super-service/internal/transport/dto"

	"github.com/xtsank/mypills-super-service/internal/service"
)

type AuthHandler struct {
	authService service.IAuthService
}

func NewAuthHandler(authService service.IAuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input dto.CreateUserDto

	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := command.CreateUserCmd{
		Login:      input.Login,
		Password:   input.Password,
		Sex:        input.Sex,
		Weight:     input.Weight,
		Age:        input.Age,
		IsPregnant: input.IsPregnant,
		IsDriver:   input.IsDriver,
		Illnesses:  input.Illnesses,
		Allergies:  input.Allergies,
	}

	createdUser, err := h.authService.Register(c.Request.Context(), cmd)
	if err != nil {
		//if errors.Is(err, service.ErrUserAlreadyExists) {
		//	c.JSON(http.StatusConflict, gin.H{"error": "user with this login already exists"})
		//	return
		//}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": createdUser.ID})
}
