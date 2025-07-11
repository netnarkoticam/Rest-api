package server

import (
	"net/http"
	"resapi/internal/domain/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

// Repository интерфейс для работы с БД
type Repository interface {
	GetAllUser() ([]models.User, error)
	GetUser(int) (models.User, error)
	InsertUser(models.User) error
	DeleteUser(int) error
	PatchUser(id int, patch map[string]interface{}) (models.User, error)
	UpdateUser(user models.User) error
}

// Server — HTTP сервер
type Server struct {
	Db Repository
}

// GetUsersHandler возвращает всех пользователей
func (s *Server) GetUsersHandler(c *gin.Context) {
	users, err := s.Db.GetAllUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// SaveUserHandler сохраняет нового пользователя
func (s *Server) SaveUserHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.Db.InsertUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, "User %v was saved", user.ID)
}

// PatchUserHandler частично обновляет пользователя
func (s *Server) PatchUserHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var patch map[string]interface{}
	if err := c.ShouldBindJSON(&patch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.Db.PatchUser(id, patch)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUserHandler удаляет пользователя
func (s *Server) DeleteUserHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = s.Db.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateUserHandler обновляет пользователя целиком
func (s *Server) UpdateUserHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = id

	if err := s.Db.UpdateUser(user); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}