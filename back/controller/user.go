package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/wonderfulsuccess/go-web-app/back/model"
)

// UserController contains CRUD handlers for the User model.
type UserController struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{db: db}
}

func (uc *UserController) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("", uc.List)
	group.POST("", uc.Create)
	group.GET(":id", uc.Get)
	group.PUT(":id", uc.Update)
	group.DELETE(":id", uc.Delete)
}

func (uc *UserController) List(c *gin.Context) {
	var users []model.User
	if err := uc.db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (uc *UserController) Get(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	if err := uc.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uc *UserController) Create(c *gin.Context) {
	var payload model.User
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.db.Create(&payload).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payload)
}

func (uc *UserController) Update(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var payload model.User
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payload.ID = id

	updates := map[string]interface{}{
		"name":  payload.Name,
		"email": payload.Email,
		"role":  payload.Role,
	}

	if err := uc.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	payload.ID = id
	c.JSON(http.StatusOK, payload)
}

func (uc *UserController) Delete(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.db.Delete(&model.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func parseID(idParam string) (uint, error) {
	v, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}
