package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rdsarjito/marketplace-backend/constants"
	"github.com/rdsarjito/marketplace-backend/domain/dto/response"
	"github.com/rdsarjito/marketplace-backend/repositories"
)

type ProvinceCityHandler struct {
	provinceCityRepo repositories.ProvinceCityRepository
}

func NewProvinceCityHandler(provinceCityRepo repositories.ProvinceCityRepository) *ProvinceCityHandler {
	return &ProvinceCityHandler{
		provinceCityRepo: provinceCityRepo,
	}
}

func (h *ProvinceCityHandler) GetListProvince(c *fiber.Ctx) error {
	provinces, err := h.provinceCityRepo.GetListProvince()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(constants.ErrExternalAPI, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(constants.MsgDataRetrieved, provinces))
}

func (h *ProvinceCityHandler) GetDetailProvince(c *fiber.Ctx) error {
	provID := c.Params("prov_id")
	if provID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Province ID is required", nil))
	}

	province, err := h.provinceCityRepo.GetDetailProvince(provID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse(constants.ErrProvinceNotFound, nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(constants.MsgDataRetrieved, province))
}

func (h *ProvinceCityHandler) GetListCity(c *fiber.Ctx) error {
	provID := c.Params("prov_id")
	if provID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Province ID is required", nil))
	}

	cities, err := h.provinceCityRepo.GetListCity(provID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(constants.ErrExternalAPI, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(constants.MsgDataRetrieved, cities))
}

func (h *ProvinceCityHandler) GetDetailCity(c *fiber.Ctx) error {
	cityID := c.Params("city_id")
	if cityID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("City ID is required", nil))
	}

	city, err := h.provinceCityRepo.GetDetailCity(cityID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse(constants.ErrCityNotFound, nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(constants.MsgDataRetrieved, city))
}
