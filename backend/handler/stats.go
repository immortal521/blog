package handler

import (
	"blog-server/response"
	"blog-server/service"

	"github.com/gofiber/fiber/v2"
)

type IStatsHandler interface {
	GetDashboardStats(c *fiber.Ctx) error
}

type statsHandler struct {
	svc service.IStatsService
}

func (s *statsHandler) GetDashboardStats(c *fiber.Ctx) error {
	res, err := s.svc.GetDashboardStats(c.UserContext())
	if err != nil {
		return err
	}

	result := response.Success(res)
	return c.JSON(result)
}

func NewStatsHandler(svc service.IStatsService) IStatsHandler {
	return &statsHandler{svc}
}
