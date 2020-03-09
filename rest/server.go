package rest

import (
	"github.com/gofiber/fiber"
	"github.com/minskylab/palapi"
)

type Service struct {
	port int64
	app  *fiber.App
	parentManager *palapi.Manager
}

func NewService(manager *palapi.Manager,port int64) (*Service, error) {
	return &Service{
		port: port,
		parentManager: manager,
		app:  fiber.New(),
	},nil
}

func (s *Service) Run() error {
	s.describe()
	return s.app.Listen(int(s.port))
}