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

type TRXHandler struct {
	trxService services.TRXService
	validator  *validator.Validate
}

func NewTRXHandler(trxService services.TRXService) *TRXHandler {
	return &TRXHandler{
		trxService: trxService,
		validator:  validator.New(),
	}
}

func (h *TRXHandler) GetListTRX(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	trxs, err := h.trxService.GetListTRX(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(err.Error(), nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(constants.MsgDataRetrieved, trxs))
}

func (h *TRXHandler) GetDetailTRX(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	trxID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Invalid transaction ID", nil))
	}

	trx, err := h.trxService.GetDetailTRX(userID, trxID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse(err.Error(), nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(constants.MsgDataRetrieved, trx))
}

func (h *TRXHandler) CreateTRX(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	var req request.CreateTRXRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Invalid request body", err.Error()))
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Validation failed", err.Error()))
	}

	trx, err := h.trxService.CreateTRX(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(err.Error(), nil))
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(constants.MsgTransactionCreated, trx))
}

// CheckPayment checks payment status for a transaction
func (h *TRXHandler) CheckPayment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	trxID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Invalid transaction ID", nil))
	}

	trx, err := h.trxService.CheckPaymentStatus(userID, trxID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(err.Error(), nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse("Payment status checked successfully", trx))
}
