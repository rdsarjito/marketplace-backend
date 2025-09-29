package handlers

import (
    "fmt"
    "strconv"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/go-playground/validator/v10"
    "github.com/rdsarjito/marketplace-backend/constants"
    "github.com/rdsarjito/marketplace-backend/domain/dto/request"
    "github.com/rdsarjito/marketplace-backend/domain/dto/response"
    "github.com/rdsarjito/marketplace-backend/services"
)

type ProductHandler struct {
	productService services.ProductService
	validator      *validator.Validate
}

func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
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
    // Save file to uploads folder
    path := fmt.Sprintf("uploads/%d_%s", time.Now().UnixNano(), file.Filename)
    if err := c.SaveFile(file, path); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse("Upload failed", err.Error()))
    }
    // Build public URL (served via /uploads static) + cache-busting
    url := fmt.Sprintf("/%s?v=%d", path, time.Now().Unix())
    prod, err := h.productService.AddProductPhoto(userID, id, url)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(err.Error(), nil))
    }
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse("Photo uploaded", prod))
}
