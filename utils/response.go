package utils

import (
    "github.com/gofiber/fiber/v2"
)

type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Status  bool        `json:"status"`
    Data    interface{} `json:"data,omitempty"`
}

// Success Response (Code 200)
func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
    return c.Status(fiber.StatusOK).JSON(Response{
        Code:    fiber.StatusOK,
        Message: message,
        Status:  true,
        Data:    data,
    })
}

// Created Response (Code 201)
func CreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
    return c.Status(fiber.StatusCreated).JSON(Response{
        Code:    fiber.StatusCreated,
        Message: message,
        Status:  true,
        Data:    data,
    })
}

// Bad Request Response (Code 400)
func BadRequestResponse(c *fiber.Ctx, message string, data interface{}) error {
    return c.Status(fiber.StatusBadRequest).JSON(Response{
        Code:    fiber.StatusBadRequest,
        Message: message,
        Status:  false,
        Data:    data,
    })
}

// Unauthorized Response (Code 401)
func UnauthorizedResponse(c *fiber.Ctx, message string, data interface{}) error {
    return c.Status(fiber.StatusUnauthorized).JSON(Response{
        Code:    fiber.StatusUnauthorized,
        Message: message,
        Status:  false,
        Data:    data,
    })
}

// Forbidden Response (Code 403)
func ForbiddenResponse(c *fiber.Ctx, message string, data interface{}) error {
    return c.Status(fiber.StatusForbidden).JSON(Response{
        Code:    fiber.StatusForbidden,
        Message: message,
        Status:  false,
        Data:    data,
    })
}

// Not Found Response (Code 404)
func NotFoundResponse(c *fiber.Ctx, message string, data interface{}) error {
    return c.Status(fiber.StatusNotFound).JSON(Response{
        Code:    fiber.StatusNotFound,
        Message: message,
        Status:  false,
        Data:    data,
    })
}

// Internal Server Error Response (Code 500)
func InternalServerErrorResponse(c *fiber.Ctx, message string, data interface{}) error {
    return c.Status(fiber.StatusInternalServerError).JSON(Response{
        Code:    fiber.StatusInternalServerError,
        Message: message,
        Status:  false,
        Data:    data,
    })
}

// Custom Error Response
func ErrorResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
    return c.Status(statusCode).JSON(Response{
        Code:    statusCode,
        Message: message,
        Status:  false,
        Data:    data,
    })
}