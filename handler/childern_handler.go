package handler

import (
	"net/http"

	"shuttle/models/dto"
	"shuttle/services"
	"shuttle/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ChildernHandlerInterface interface {
	GetAllChilderns(c *fiber.Ctx) error
	GetSpecChildern(c *fiber.Ctx) error
	UpdateChildern(c *fiber.Ctx) error
}

type ChildernHandler struct {
	ChildernService services.ChildernServiceInterface
	DB              *sqlx.DB
}

func NewChildernHandler(childernService services.ChildernServiceInterface) *ChildernHandler {
	return &ChildernHandler{
		ChildernService: childernService,
	}
}

func (handler *ChildernHandler) GetAllChilderns(c *fiber.Ctx) error {
	id, ok := c.Locals("userUUID").(string)
	if !ok || id == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "User UUID is missing or invalid",
		})
	}
	childernsDTO, total, err := handler.ChildernService.GetAllChilderns(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch students",
			"details": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"data":  childernsDTO,
		"total": total,
	})
}

func (handler *ChildernHandler) GetSpecChildern(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}
	studentDTO, err := handler.ChildernService.GetSpecChildern(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch data",
		})
	}
	return c.Status(http.StatusOK).JSON(studentDTO)
}

func (handler *ChildernHandler) UpdateChildern(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"code":    http.StatusBadRequest,
			"message": "Student ID is required",
			"status":  false,
		})
	}
	username := c.Locals("user_name")
	if username == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"code":    http.StatusUnauthorized,
			"message": "Unauthorized",
			"status":  false,
		})
	}
	var studentReqDTO dto.StudentRequestByParentDTO
	if err := c.BodyParser(&studentReqDTO); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"code":    http.StatusBadRequest,
			"message": "Invalid request data",
			"status":  false,
		})
	}
	if err := utils.ValidateStruct(c, &studentReqDTO); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"code":    http.StatusBadRequest,
			"message": "Validation error: " + err.Error(),
			"status":  false,
		})
	}
	err := handler.ChildernService.UpdateChildern(id, studentReqDTO, username.(string))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "Failed to update student data",
			"status":  false,
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "Student updated successfully",
		"status":  true,
	})
}

func (handler *ChildernHandler) UpdateChildernStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Student ID is required",
			"status":  false,
		})
	}
	username, ok := c.Locals("user_name").(string)
	if !ok || username == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"message": "Unauthorized",
			"status":  false,
		})
	}
	var studentReqDTO dto.StudentStatusRequestByParentDTO
	if err := c.BodyParser(&studentReqDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid request data",
			"status":  false,
		})
	}
	if err := utils.ValidateStruct(c, studentReqDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Validation error: " + err.Error(),
			"status":  false,
		})
	}
	if err := handler.ChildernService.UpdateChildernStatus(id, studentReqDTO, username); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Failed to update student data: " + err.Error(),
			"status":  false,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Student updated successfully",
		"status":  true,
	})
}