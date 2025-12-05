package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
	"github.com/rdsarjito/marketplace-backend/constants"
	"github.com/rdsarjito/marketplace-backend/domain/dto/request"
	"github.com/rdsarjito/marketplace-backend/domain/dto/response"
	"github.com/rdsarjito/marketplace-backend/services"
)

type AuthHandler struct {
	authService services.AuthService
	validator   *validator.Validate
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator.New(),
	}
}

func (h *AuthHandler) RegisterUser(c *fiber.Ctx) error {
	var req request.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Invalid request body", err.Error()))
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Validation failed", err.Error()))
	}

	authResponse, err := h.authService.RegisterUser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(err.Error(), nil))
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(constants.MsgUserRegistered, authResponse))
}

func (h *AuthHandler) LoginUser(c *fiber.Ctx) error {
	var req request.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Invalid request body", err.Error()))
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Validation failed", err.Error()))
	}

	authResponse, err := h.authService.LoginUser(&req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse(err.Error(), nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(constants.MsgUserLoggedIn, authResponse))
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req request.ForgotPasswordRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Invalid request body", err.Error()))
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Validation failed", err.Error()))
	}

	forgotResponse, err := h.authService.ForgotPassword(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(err.Error(), nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(forgotResponse.Message, nil))
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req request.ResetPasswordRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Invalid request body", err.Error()))
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Validation failed", err.Error()))
	}

	resetResponse, err := h.authService.ResetPassword(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(err.Error(), nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(resetResponse.Message, nil))
}

// Google OAuth
func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
    // Build Google OAuth URL and redirect
    clientID := os.Getenv("GOOGLE_CLIENT_ID")
    redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
    if clientID == "" || redirectURL == "" {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse("Google OAuth not configured", nil))
    }
    scope := url.QueryEscape("openid email profile")
    oauthURL := fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&prompt=consent", clientID, url.QueryEscape(redirectURL), scope)
    return c.Redirect(oauthURL, fiber.StatusTemporaryRedirect)
}

func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
    code := c.Query("code")
    if code == "" {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Missing code", nil))
    }

    clientID := os.Getenv("GOOGLE_CLIENT_ID")
    clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
    redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
    if clientID == "" || clientSecret == "" || redirectURL == "" {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse("Google OAuth not configured", nil))
    }

    // Exchange code for tokens
    tokenResp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
        "code":          {code},
        "client_id":     {clientID},
        "client_secret": {clientSecret},
        "redirect_uri":  {redirectURL},
        "grant_type":    {"authorization_code"},
    })
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Failed to exchange code", err.Error()))
    }
    defer tokenResp.Body.Close()

    var tokenData struct {
        AccessToken string `json:"access_token"`
        IdToken     string `json:"id_token"`
        ExpiresIn   int    `json:"expires_in"`
        TokenType   string `json:"token_type"`
        Scope       string `json:"scope"`
    }
    if err := json.NewDecoder(tokenResp.Body).Decode(&tokenData); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Invalid token response", err.Error()))
    }

    // Fetch userinfo
    req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
    req.Header.Set("Authorization", "Bearer "+tokenData.AccessToken)
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Failed to fetch userinfo", err.Error()))
    }
    defer resp.Body.Close()

    var userinfo struct {
        Email         string `json:"email"`
        VerifiedEmail bool   `json:"email_verified"`
        Name          string `json:"name"`
        Picture       string `json:"picture"`
        Sub           string `json:"sub"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&userinfo); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Invalid userinfo response", err.Error()))
    }
    if userinfo.Email == "" || !userinfo.VerifiedEmail {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse("Email not verified", nil))
    }

    // Login or create user
    authResp, err := h.authService.LoginWithGoogle(userinfo.Email, userinfo.Name, userinfo.Picture)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(err.Error(), nil))
    }

    // Redirect to frontend with token
    frontendURL := os.Getenv("FRONTEND_URL")
    if frontendURL == "" {
        frontendURL = "http://localhost:5173"
    }
    redirect := fmt.Sprintf("%s/login?token=%s", strings.TrimRight(frontendURL, "/"), authResp.Token)
    return c.Redirect(redirect, fiber.StatusTemporaryRedirect)
}
