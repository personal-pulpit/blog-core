package handlers

import (
	"blog/api/helpers"
	"blog/internal/repository"
	"blog/internal/service/category"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService category.CategoryService
}

type createCategoryInput struct {
	Name string `json:"name" binding:"required" example:"Technology"`
}

func NewCategoryHandler(categoryService category.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// @Summary Get all categories
// @Description Retrieve all categories
// @Tags categories
// @Produce json
// @Success 200 {object} helpers.HttpResponse{data=map[string]interface{}} "Successfully retrieved categories"
// @Failure 500 {object} helpers.HttpResponse{data=map[string]interface{}} "Failed to fetch categories"
// @Router /api/v1/categories [get]
func (c *CategoryHandler) GetAll(ctx *gin.Context) {
	categories, err := c.categoryService.GetAllCategories(ctx)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusInternalServerError,
			"Failed to fetch categories",
			map[string]interface{}{
				"error": err.Error(),
				"count": 0,
			})
		return
	}

	helpers.RespondWithSuccess(ctx, http.StatusOK,
		"Successfully retrieved categories",
		map[string]interface{}{
			"categories": categories,
			"count":      len(categories),
		})
}

// @Summary Get a category by ID
// @Description Retrieve a category by its ID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} helpers.HttpResponse{data=map[string]interface{}} "Category retrieved successfully"
// @Failure 404 {object} helpers.HttpResponse{data=map[string]interface{}} "Category not found"
// @Failure 500 {object} helpers.HttpResponse{data=map[string]interface{}} "Failed fetch category"
// @Router /api/v1/categories/{id} [get]
func (c *CategoryHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	ID := helpers.StringToInt(id)

	category, err := c.categoryService.GetCategoryByID(ctx, uint(ID))
	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			helpers.RespondWithError(ctx, http.StatusNotFound,
				"Category not found", map[string]interface{}{
					"category_id": ID,
					"error":       err.Error(),
				})
			return
		}

		helpers.RespondWithError(ctx, http.StatusInternalServerError, "Failed fetch category", map[string]interface{}{
			"category_id": id,
			"error":       err.Error(),
		})
		return
	}

	helpers.RespondWithSuccess(ctx, http.StatusOK,
		"Category retrieved successfully", map[string]interface{}{
			"category": category,
		})
}

// @Summary Create a new category
// @Description Create a new category with the provided name
// @Tags categories
// @Accept json
// @Produce json
// @Param category body createCategoryInput true "Category input"
// @Success 201 {object} helpers.HttpResponse{data=map[string]interface{}} "Category created successfully"
// @Failure 400 {object} helpers.HttpResponse{data=map[string]interface{}} "Invalid category data"
// @Failure 500 {object} helpers.HttpResponse{data=map[string]interface{}} "Failed to create category"
// @Router /api/v1/categories [post]
func (c *CategoryHandler) Create(ctx *gin.Context) {
	var ci = new(createCategoryInput)

	err := ctx.ShouldBindJSON(ci)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusBadRequest,
			"Invalid category data", map[string]interface{}{
				"validation_errors": err,
				"provided_data":     ci,
			})
		return
	}

	category, err := c.categoryService.CreateCategory(ctx, ci.Name)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusInternalServerError,
			"Failed to create category", map[string]interface{}{
				"error": err.Error(),
			})
		return
	}

	helpers.RespondWithSuccess(ctx, http.StatusCreated,
		"Category created successfully", map[string]interface{}{
			"category": category,
		})
}

// @Summary Delete a category by ID
// @Description Delete a category by its ID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} helpers.HttpResponse{data=map[string]interface{}} "Category deleted successfully"
// @Failure 500 {object} helpers.HttpResponse{data=map[string]interface{}} "Failed to delete category"
// @Failure 404 {object} helpers.HttpResponse{data=map[string]interface{}} "Category not found"
// @Router /api/v1/categories/{id} [delete]
func (c *CategoryHandler) DeleteByID(ctx *gin.Context) {
	id := ctx.Param("id")

	ID := helpers.StringToInt(id)

	err := c.categoryService.DeleteCategoryByID(ctx, uint(ID))
	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			helpers.RespondWithError(ctx, http.StatusNotFound,
				"Category not found", map[string]interface{}{
					"category_id": ID,
					"error":       err.Error(),
				})
			return
		}

		helpers.RespondWithError(ctx, http.StatusInternalServerError,
			"Failed to delete category", map[string]interface{}{
				"error": err.Error(),
			})
		return
	}

	helpers.RespondWithSuccess(ctx, http.StatusOK,
		"Category deleted successfully", map[string]interface{}{
			"category_id": id,
		})
}
