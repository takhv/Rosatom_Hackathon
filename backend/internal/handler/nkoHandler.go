package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"rosatom.ru/nko/internal/models"
	"rosatom.ru/nko/internal/repository"
)

type NKOHandler struct {
	repo repository.NKORepo
}

func NewNKOHandler(repo repository.NKORepo) *NKOHandler {
	return &NKOHandler{repo: repo}
}

func (h *NKOHandler) GetAllNKO(context echo.Context) error {
	nkos, err := h.repo.GetAllNKO()
	if err != nil {
		return context.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return context.JSON(http.StatusOK, nkos)
}

func (h *NKOHandler) GetByID(context echo.Context) error {
	idStr := context.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return context.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	nko, err := h.repo.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return context.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		return context.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Internal server error",
		})
	}
	return context.JSON(http.StatusOK, nko)
}

func (h *NKOHandler) GetNKOName(context echo.Context) error {
	idStr := context.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return context.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	nko, err := h.repo.GetNKOName(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return context.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		return context.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Internal server error",
		})
	}
	return context.JSON(http.StatusOK, nko)
}

func (h *NKOHandler) SearchNKO(c echo.Context) error {
	name := c.QueryParam("name")
	category := c.QueryParam("category")

	if name == "" && category == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Provide at least one search parameter (name or category)",
		})
	}

	nko, err := h.repo.SearchNKO(name, category)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "NKO not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Internal server error",
		})
	}

	return c.JSON(200, nko)
}

func (h *NKOHandler) CreateNKO(c echo.Context) error {
	var req models.CreateNKO

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid JSON format",
		})
	}

	if req.Name == "" {
		return c.JSON(400, map[string]string{
			"error": "Name is required",
		})
	}
	if req.Category == "" {
		return c.JSON(400, map[string]string{
			"error": "Category is required",
		})
	}
	if req.CityID == 0 {
		return c.JSON(400, map[string]string{
			"error": "City ID is required",
		})
	}

	nko := models.NKO{
		Name:                  req.Name,
		Category:              req.Category,
		Description:           req.Description,
		Volunteer_description: req.VolunteerDescription,
		Phone:                 req.Phone,
		Address:               req.Address,
		Logo_url:              req.LogoURL,
		Website_url:           req.WebsiteURL,
		Social_links:          req.SocialLinks,
		City_id:               req.CityID,
	}

	err := h.repo.CreateNKO(&nko)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create NKO: " + err.Error(),
		})
	}

	return c.JSON(201, map[string]interface{}{
		"message": "NKO created successfully",
		"data":    nko,
	})
}
