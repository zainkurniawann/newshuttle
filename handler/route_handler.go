package handler

import (
	"fmt"
	"log"
	"shuttle/models/dto"
	"shuttle/services"
	"shuttle/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RouteHandlerInterface interface {
	GetAllRoutesByAS(c *fiber.Ctx) error
	GetAllRouteAssignments(c *fiber.Ctx) error
	GetSpecRouteByAS(c *fiber.Ctx) error
	GetAllRoutesByDriver(c *fiber.Ctx) error
	AddRoute(c *fiber.Ctx) error
	UpdateStudentOrder(c *fiber.Ctx) error
	UpdateRoute(c *fiber.Ctx) error 
	DeleteRoute(c *fiber.Ctx) error

	GetDriverDistance(c *fiber.Ctx) error
}

type routeHandler struct {
	routeService services.RouteServiceInterface
}

func NewRouteHttpHandler(routeService services.RouteServiceInterface) RouteHandlerInterface {
	return &routeHandler{
		routeService: routeService,
	}
}

func (h *routeHandler) GetDriverDistance(c *fiber.Ctx) error {
	// Dummy data (seharusnya dari database)
	driverStart := [2]float64{-7.773161987268529, 110.3747415099855}
	students := [][2]float64{
		{-7.7115806244533305, 110.41349437928434}, // Siswa 3
		{-7.763653089789303, 110.42236650540428}, // Siswa 1
		{-7.703233845448127, 110.43105534860973}, // Siswa 2
	}
	school := [2]float64{-7.715987795086408, 110.40701270626889}

	// Hitung total jarak tempuh driver
	totalDistance := h.routeService.GetTotalDistance(driverStart, students, school)

	return c.JSON(fiber.Map{
		"message":       "Total jarak tempuh berhasil dihitung",
		"total_distance": totalDistance,
	})
}


func (handler *routeHandler) GetAllRoutesByAS(c *fiber.Ctx) error {
	// Ambil schoolUUID dari token
	schoolUUID, ok := c.Locals("schoolUUID").(string)
	if !ok {
		return utils.BadRequestResponse(c, "Invalid token or schoolUUID", nil)
	}

	// Ambil query parameter untuk pagination
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		return utils.BadRequestResponse(c, "Invalid page number", nil)
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		return utils.BadRequestResponse(c, "Invalid limit number", nil)
	}

	sortField := c.Query("sort_by", "route_name")
	sortDirection := c.Query("direction", "asc")

	if sortDirection != "asc" && sortDirection != "desc" {
		return utils.BadRequestResponse(c, "Invalid sort direction, use 'asc' or 'desc'", nil)
	}

	// Panggil service untuk mendapatkan data dan total items
	routes, totalItems, err := handler.routeService.GetAllRoutesByAS(page, limit, sortField, sortDirection, schoolUUID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to fetch routes", nil)
	}

	// Hitung total halaman
	totalPages := (totalItems + limit - 1) / limit
	if page > totalPages {
		if totalItems > 0 {
			return utils.BadRequestResponse(c, "Page number out of range", nil)
		}
		page = 1
	}

	// Hitung start dan end
	start := (page-1)*limit + 1
	if totalItems == 0 || start > totalItems {
		start = 0
	}

	end := start + len(routes) - 1
	if end > totalItems {
		end = totalItems
	}

	if len(routes) == 0 {
		start = 0
		end = 0
	}

	// Response dengan metadata pagination
	response := fiber.Map{
		"data": routes,
		"meta": fiber.Map{
			"current_page":   page,
			"total_pages":    totalPages,
			"per_page_items": limit,
			"total_items":    totalItems,
			"showing":        fmt.Sprintf("Showing %d-%d of %d", start, end, totalItems),
		},
	}

	return utils.SuccessResponse(c, "Routes fetched successfully", response)
}

func (h *routeHandler) GetAllRouteAssignments(c *fiber.Ctx) error {
    log.Println("Starting GetAllRouteAssignments handler")

    // Ambil query parameter untuk pagination
    page, err := strconv.Atoi(c.Query("page", "1"))
    if err != nil || page < 1 {
        log.Printf("Invalid page number: %s\n", c.Query("page"))
        return utils.BadRequestResponse(c, "Invalid page number", nil)
    }
    log.Printf("Page number: %d\n", page)

    limit, err := strconv.Atoi(c.Query("limit", "10"))
    if err != nil || limit < 1 {
        log.Printf("Invalid limit number: %s\n", c.Query("limit"))
        return utils.BadRequestResponse(c, "Invalid limit number", nil)
    }
    log.Printf("Limit: %d\n", limit)

    // Ambil query parameter untuk sorting
    sortField := c.Query("sort_by", "route_name") // Default sort field
    sortDirection := c.Query("direction", "asc") // Default sort direction
    log.Printf("Sort field: %s, Sort direction: %s\n", sortField, sortDirection)

    // Validasi sort direction
    if sortDirection != "asc" && sortDirection != "desc" {
        log.Printf("Invalid sort direction: %s\n", sortDirection)
        return utils.BadRequestResponse(c, "Invalid sort direction, use 'asc' or 'desc'", nil)
    }

    // Panggil service dengan pagination dan sorting
    log.Println("Fetching route assignments from service")
    routes, totalItems, err := h.routeService.GetAllRouteAssignments(page, limit, sortField, sortDirection)
    if err != nil {
        log.Printf("Failed to fetch route assignments: %v\n", err)
        return utils.InternalServerErrorResponse(c, "Failed to fetch route assignments", nil)
    }
    log.Printf("Fetched %d route assignments. Total items: %d\n", len(routes), totalItems)

    // Hitung total halaman
    totalPages := (totalItems + limit - 1) / limit
    log.Printf("Total pages: %d\n", totalPages)

    if page > totalPages {
        log.Printf("Page number %d out of range (total pages: %d)\n", page, totalPages)
        if totalItems > 0 {
            return utils.BadRequestResponse(c, "Page number out of range", nil)
        }
        page = 1
    }

    // Hitung start dan end item untuk metadata
    start := (page-1)*limit + 1
    end := start + len(routes) - 1
    if totalItems == 0 || start > totalItems {
        start = 0
        end = 0
    }
    if end > totalItems {
        end = totalItems
    }
    log.Printf("Pagination metadata - Start: %d, End: %d\n", start, end)

    // Buat respons JSON dengan metadata pagination
    response := fiber.Map{
        "data": routes,
        "meta": fiber.Map{
            "current_page":   page,
            "total_pages":    totalPages,
            "per_page_items": limit,
            "total_items":    totalItems,
            "showing":        fmt.Sprintf("Showing %d-%d of %d", start, end, totalItems),
        },
    }
    log.Println("Successfully prepared response")
    return utils.SuccessResponse(c, "Route assignments fetched successfully", response)
}

func (handler *routeHandler) GetSpecRouteByAS(c *fiber.Ctx) error {
	routeNameUUID := c.Params("id")
	driverUUID, err := handler.routeService.GetDriverUUIDByRouteName(routeNameUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get driver UUID"})
	}
	routeResponse, err := handler.routeService.GetSpecRouteByAS(routeNameUUID, driverUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(routeResponse)
}

func (handler *routeHandler) GetAllRoutesByDriver(c *fiber.Ctx) error {
	driverUUID, ok := c.Locals("userUUID").(string)
	if !ok {
		log.Println("Token does not contain driver UUID")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Token does not contain driver UUID"})
	}
	if _, err := uuid.Parse(driverUUID); err != nil {
		log.Println("Invalid UUID format:", driverUUID)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID format"})
	}
	log.Println("Fetching routes for driver UUID:", driverUUID)
	routes, err := handler.routeService.GetAllRoutesByDriver(driverUUID)
	if err != nil {
		log.Println("Error fetching routes:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch routes"})
	}
	log.Println("Routes fetched successfully for driver UUID:", driverUUID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"routes": routes})
}

func (handler *routeHandler) AddRoute(c *fiber.Ctx) error {
    // Ambil schoolUUID dari token
    schoolUUID, ok := c.Locals("schoolUUID").(string)
    if !ok || schoolUUID == "" {
        return utils.InternalServerErrorResponse(c, "Token does not contain schoolUUID", nil)
    }

    // Ambil username dari token
    username, ok := c.Locals("user_name").(string)
    if !ok || username == "" {
        return utils.InternalServerErrorResponse(c, "Token does not contain username", nil)
    }

    // Parse body request
    route := new(dto.RoutesRequestDTO)
    if err := c.BodyParser(route); err != nil {
        return utils.BadRequestResponse(c, "Invalid request body", nil)
    }

    // Validasi payload
    if err := utils.ValidateStruct(c, route); err != nil {
        return utils.BadRequestResponse(c, err.Error(), nil)
    }

    // Panggil service untuk menambahkan rute
    err := handler.routeService.AddRoute(*route, schoolUUID, username)
    if err != nil {
        // Tangani error spesifik
        switch {
		case strings.Contains(err.Error(), "Maximum seats exceeded"):
			return utils.BadRequestResponse(c, err.Error(), nil)
		case err.Error() == "same student not permitted":
			return utils.BadRequestResponse(c, "Same student not permitted", nil)
		case err.Error() == "student not found":
			return utils.BadRequestResponse(c, "Student not found", nil)
		case err.Error() == "driver not found":
			return utils.BadRequestResponse(c, "Driver not found", nil)
		case err.Error() == "driver already assigned to another route":
			return utils.BadRequestResponse(c, "Driver already assigned to another route", nil)
		default:
			return utils.InternalServerErrorResponse(c, "Failed to add route", err)
	}	
    }

    return utils.SuccessResponse(c, "Route added successfully", nil)
}

func (h *routeHandler) UpdateRoute(c *fiber.Ctx) error {
    var requestDTO dto.UpdateRouteRequest

    // Parsing body menjadi DTO
    if err := c.BodyParser(&requestDTO); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
    }

    // Mengambil parameter tambahan
    routeNameUUID := c.Params("id")
    schoolUUID, _ := c.Locals("schoolUUID").(string)
    username, _ := c.Locals("user_name").(string)

    // Memanggil service dengan DTO
    err := h.routeService.UpdateRoute(requestDTO, routeNameUUID, schoolUUID, username)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"message": "Route updated successfully"})
}

func (handler *routeHandler) UpdateStudentOrder(c *fiber.Ctx) error {
	// Ambil student_uuid dari URL
	studentUUID := c.Params("id")
	if _, err := uuid.Parse(studentUUID); err != nil {
		log.Println("Invalid UUID format:", studentUUID)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID format"})
	}

	// Parse new_order dari body request
	var updateStudentOrderDTO dto.UpdateStudentOrderDTO
	if err := c.BodyParser(&updateStudentOrderDTO); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Panggil service untuk update student order
	err := handler.routeService.UpdateStudentOrderByDriver(studentUUID, updateStudentOrderDTO.NewOrder)
	if err != nil {
		log.Println("Error updating student order:", err)
		if err.Error() == "student_uuid not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Student not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update student order"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Student order updated successfully"})
}

func (handler *routeHandler) DeleteRoute(c *fiber.Ctx) error {
	routenameUUID := c.Params("id")
	schoolUUID, ok := c.Locals("schoolUUID").(string)
	if !ok {
		return utils.InternalServerErrorResponse(c, "Token does not contain schoolUUID", nil)
	}
	username, ok := c.Locals("user_name").(string)
	if !ok {
		return utils.InternalServerErrorResponse(c, "Token does not contain username", nil)
	}
	if err := handler.routeService.DeleteRoute(routenameUUID, schoolUUID, username); err != nil {
		if err.Error() == "route not found" {
			return utils.NotFoundResponse(c, "Route not found", nil)
		}
		return utils.InternalServerErrorResponse(c, err.Error(), nil)
	}
	return utils.SuccessResponse(c, "Route deleted successfully", nil)
}