package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rdsarjito/marketplace-backend/constants"
	"github.com/rdsarjito/marketplace-backend/domain/dto/request"
	"github.com/rdsarjito/marketplace-backend/domain/dto/response"
	"github.com/rdsarjito/marketplace-backend/services"
	"github.com/rdsarjito/marketplace-backend/storage"
)

type ProductHandler struct {
	productService services.ProductService
	storage        storage.MediaStorage
	validator      *validator.Validate
}

func NewProductHandler(productService services.ProductService, mediaStorage storage.MediaStorage) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		storage:        mediaStorage,
		validator:      validator.New(),
	}
}

func (h *ProductHandler) GetListProduct(c *fiber.Ctx) error {
	products, err := h.productService.GetListProduct()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(err.Error(), nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(constants.MsgDataRetrieved, products))
}

func (h *ProductHandler) GetDetailProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Invalid product ID", nil))
	}

	product, err := h.productService.GetDetailProduct(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse(err.Error(), nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(constants.MsgDataRetrieved, product))
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	var req request.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Invalid request body", err.Error()))
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Validation failed", err.Error()))
	}

	product, err := h.productService.CreateProduct(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(err.Error(), nil))
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(constants.MsgProductCreated, product))
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Invalid product ID", nil))
	}

	var req request.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Invalid request body", err.Error()))
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Validation failed", err.Error()))
	}

	product, err := h.productService.UpdateProduct(userID, id, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(err.Error(), nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(constants.MsgProductUpdated, product))
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Invalid product ID", nil))
	}

	err = h.productService.DeleteProduct(userID, id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(err.Error(), nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(constants.MsgProductDeleted, nil))
}

func (h *ProductHandler) UploadProductPhoto(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Invalid product ID", nil))
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("File not found", nil))
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse("Unable to read file", err.Error()))
	}
	defer src.Close()

	objectName := fmt.Sprintf("products/%d/%d_%s", id, time.Now().UnixNano(), sanitizeFilename(file.Filename))
	url, err := h.storage.Upload(c.UserContext(), objectName, src, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse("Upload failed", err.Error()))
	}

	prod, err := h.productService.AddProductPhoto(userID, id, url)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(err.Error(), nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse("Photo uploaded", prod))
}

func sanitizeFilename(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")

	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') ||
			(r >= '0' && r <= '9') ||
			r == '.' || r == '_' || r == '-' {
			return r
		}
		return '-'
	}, name)
}
