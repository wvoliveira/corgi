package log

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Add(*gin.Context, model.User, any) error

	NewHTTP(*gin.RouterGroup)
	HTTPAdd(*gin.Context)
}

type service struct{}

// NewService creates a new user management service.
func NewService() Service {
	return service{}
}

// Add change specific link by ID.
func (s service) Add(c *gin.Context, user model.User, payload any) (err error) {
	// l := logger.Logger(c.Request.Context())

	fmt.Println("USER")
	fmt.Println(user)

	fmt.Println("Payload")
	fmt.Println(payload)
	return
}
