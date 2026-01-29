package handlers

import (
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
)


const categoryBasePath = "/api/v1/categories/"


type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}


// HandleCategories - GET | POST /api/v1/categories
func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// HandleCategoryByID - GET | PUT | DELETE /api/v1/categories/{id}
func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetByID(w, r)
	case http.MethodPut:
		h.Update(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// GetAll - GET /api/v1/categories
func (h *CategoryHandler) GetAll(w http.ResponseWriter, _ *http.Request) {
	categories, err := h.service.GetAll()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, "Categories fetched successfully", categories)
}

// Create - POST /api/v1/categories
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := decodeJSON(r, &category); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.Create(&category); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, "Category created successfully", category)
}

// GetByID - GET /api/v1/categories/{id}
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := mustGetID(w, r, categoryBasePath)
	if !ok {
		return
	}

	category, err := h.service.GetByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, "Category fetched successfully", category)
}

// Update - PUT /api/v1/categories/{id}
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := mustGetID(w, r, categoryBasePath)
	if !ok {
		return
	}

	var category models.Category
	if err := decodeJSON(r, &category); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	category.ID = id
	if err := h.service.Update(&category); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, "Category updated successfully", category)
}

// Delete - DELETE /api/v1/categories/{id}
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := mustGetID(w, r, categoryBasePath)
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// writeJSON(w, http.StatusNoContent, "Category deleted successfully", nil)
	writeJSON(w, http.StatusOK, "Category deleted successfully", nil)
}
