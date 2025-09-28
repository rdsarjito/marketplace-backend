package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
	"github.com/rdsarjito/marketplace-backend/constants"
	"github.com/rdsarjito/marketplace-backend/domain/dto/request"
	"github.com/rdsarjito/marketplace-backend/domain/dto/response"
	"github.com/rdsarjito/marketplace-backend/services"
)

type ShopHandler struct {
	shopService services.ShopService
	validator   *validator.Validate
}

func NewShopHandler(shopService services.ShopService) *ShopHandler {
	return &ShopHandler{
		shopService: shopService,
		validator:   validator.New(),
	}
}

func (h *ShopHandler) MyShop(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	shop, err := h.shopService.MyShop(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse(err.Error(), nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(constants.MsgDataRetrieved, shop))
}

func (h *ShopHandler) GetListShop(c *fiber.Ctx) error {
	shops, err := h.shopService.GetListShop()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(err.Error(), nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(constants.MsgDataRetrieved, shops))
}

func (h *ShopHandler) GetDetailShop(c *fiber.Ctx) error {
	shopID, err := strconv.Atoi(c.Params("id_toko"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Invalid shop ID", nil))
	}

	shop, err := h.shopService.GetDetailShop(shopID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse(err.Error(), nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(constants.MsgDataRetrieved, shop))
}

func (h *ShopHandler) UpdateProfileShop(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	shopID, err := strconv.Atoi(c.Params("id_toko"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Invalid shop ID", nil))
	}

	var req request.UpdateShopRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Invalid request body", err.Error()))
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Validation failed", err.Error()))
	}

	shop, err := h.shopService.UpdateProfileShop(userID, shopID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(err.Error(), nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(constants.MsgShopUpdated, shop))
}
