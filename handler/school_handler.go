package handler

import (
	"fmt"
	"strconv"
	"strings"

	"shuttle/logger"
	"shuttle/models/dto"
	"shuttle/services"
	"shuttle/utils"

	"github.com/gofiber/fiber/v2"
)

type SchoolHandlerInterface interface {
	GetAllSchools(c *fiber.Ctx) error
	GetSpecSchool(c *fiber.Ctx) error
	AddSchool(c *fiber.Ctx) error
	UpdateSchool(c *fiber.Ctx) error
	DeleteSchool(c *fiber.Ctx) error
}

type schoolHandler struct {
	schoolService services.SchoolService
}

func NewSchoolHttpHandler(schoolService services.SchoolService) SchoolHandlerInterface {
	return &schoolHandler{
		schoolService: schoolService,
	}
}

func (handler *schoolHandler) GetAllSchools(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		return utils.BadRequestResponse(c, "Invalid page number", nil)
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		return utils.BadRequestResponse(c, "Invalid limit number", nil)
	}

	sortField := c.Query("sort_by", "school_id")
	sortDirection := c.Query("direction", "desc")

	if sortDirection != "asc" && sortDirection != "desc" {
		return utils.BadRequestResponse(c, "Invalid sort direction, use 'asc' or 'desc'", nil)
	}

	if !isValidSortFieldForSchools(sortField) {
		return utils.BadRequestResponse(c, "Invalid sort field", nil)
	}

	schools, totalItems, err := handler.schoolService.GetAllSchools(page, limit, sortField, sortDirection)
	if err != nil {
		logger.LogError(err, "Failed to fetch paginated schools", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	totalPages := (totalItems + limit - 1) / limit

	if page > totalPages {
		if totalItems > 0 {
			return utils.BadRequestResponse(c, "Page number out of range", nil)
		} else {
			page = 1
		}
	}

	start := (page-1)*limit + 1
	if totalItems == 0 || start > totalItems {
		start = 0
	}

	end := start + len(schools) - 1
	if end > totalItems {
		end = totalItems
	}

	if len(schools) == 0 {
		start = 0
		end = 0
	}

	response := fiber.Map{
		"data": schools,
		"meta": fiber.Map{
			"current_page":   page,
			"total_pages":    totalPages,
			"per_page_items": limit,
			"total_items":    totalItems,
			"showing":        fmt.Sprintf("Showing %d-%d of %d", start, end, totalItems),
		},
	}

	return utils.SuccessResponse(c, "Schools fetched successfully", response)
}

func (handler *schoolHandler) GetSpecSchool(c *fiber.Ctx) error {
	id := c.Params("id")

	school, err := handler.schoolService.GetSpecSchool(id)
	if err != nil {
		logger.LogError(err, "Failed to fetch specific school", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "School fetched successfully", school)
}

func (handler *schoolHandler) AddSchool(c *fiber.Ctx) error {
	username, ok := c.Locals("user_name").(string)
	if !ok || username == "" {
		// Handle missing or invalid username
		return utils.UnauthorizedResponse(c, "Unauthorized access: user_name not found", nil)
	}

	// Parse request body
	school := new(dto.SchoolRequestDTO)
	if err := c.BodyParser(school); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	// Validate the parsed data
	if err := utils.ValidateStruct(c, school); err != nil {
		return utils.BadRequestResponse(c, strings.ToUpper(err.Error()[0:1])+err.Error()[1:], nil)
	}

	// Call the service to add the school
	if err := handler.schoolService.AddSchool(*school, username); err != nil {
		// Log the error and return an internal server error response
		logger.LogError(err, "Failed to create school", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	// Return success response if everything goes well
	return utils.SuccessResponse(c, "School created successfully", nil)
}


func (handler *schoolHandler) UpdateSchool(c *fiber.Ctx) error {
	id := c.Params("id")
	username := c.Locals("user_name").(string)

	school := new(dto.SchoolRequestDTO)
	if err := c.BodyParser(school); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	if err := utils.ValidateStruct(c, school); err != nil {
		return utils.BadRequestResponse(c, strings.ToUpper(err.Error()[0:1])+err.Error()[1:], nil)
	}

	if err := handler.schoolService.UpdateSchool(id, *school, username); err != nil {
		logger.LogError(err, "Failed to update school", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "School updated successfully", nil)
}

func (handler *schoolHandler) DeleteSchool(c *fiber.Ctx) error {
	id := c.Params("id")
	username := c.Locals("user_name").(string)

	force_delete := c.Query("force_delete")

	existingSchool, err := handler.schoolService.GetSpecSchool(id)
	if err != nil {
		logger.LogError(err, "Failed to fetch specific school", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	if (existingSchool.AdminUUID != "N/A" && existingSchool.AdminUUID != "") && force_delete != "true" {
		return utils.BadRequestResponse(c, "Warning: By deleting this school, the school admin will also be deleted, continue?", nil)
	}

	if err := handler.schoolService.DeleteSchool(id, username, existingSchool.AdminUUID); err != nil {
		logger.LogError(err, "Failed to delete school", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "School deleted successfully", nil)
}

func isValidSortFieldForSchools(field string) bool {
	allowedFields := map[string]bool{
		"school_name": true,
		"school_id":   true,
		"school_point": true,
	}
	return allowedFields[field]
}
