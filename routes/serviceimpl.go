package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
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

// ------------------------------
func getBackendNameSpace(context *gin.Context) {
	namespace := context.Param("namespace")
	values, err := getBackendNameSpaceFromRedis(context, namespace)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	var activeURLs []dto.ServiceUrlDTO
	for _, url := range values {
		if url.Active {
			activeURLs = append(activeURLs, url)
		}
	}
	// todo for now get first one
	urlToServeTheRequest := activeURLs[0]
	method := context.Request.Method // Get the HTTP method type
	if method == http.MethodGet {
		resp, _ := http.Get(urlToServeTheRequest.Url)
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		context.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
		return
	}
	// end todo

	context.JSON(http.StatusOK, gin.H{"urls": values, "active": activeURLs})
}

func setBackendNameSpaceInRedis(context *gin.Context, key string, value []dto.ServiceUrlDTO) error {
	jsonData, err := json.Marshal(value)
	err = db.RedisClient.Set(context, key, jsonData, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return nil
}

func getBackendNameSpaceFromRedis(context *gin.Context, key string) ([]dto.ServiceUrlDTO, error) {
	data, err := db.RedisClient.Get(context, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get  Data: %w", err)
	}
	var values []dto.ServiceUrlDTO
	err = json.Unmarshal([]byte(data), &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}
