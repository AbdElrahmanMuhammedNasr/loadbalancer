package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"loadbalancer/db"
	"loadbalancer/dto"
	"net/http"
)

func createBackEndNameSpace(context *gin.Context) {
	backendDTO := dto.BackendDTO{}
	if err := context.ShouldBindJSON(&backendDTO); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := setBackendNameSpaceInRedis(context, backendDTO.NameSpace, backendDTO.ServiceUrl); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"name": backendDTO.NameSpace})
}

func setBackendNameSpaceInRedis(context *gin.Context, key string, value []dto.ServiceUrlDTO) error {
	jsonData, err := json.Marshal(value)
	err = db.RedisClient.Set(context, key, jsonData, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return nil
}
