package routes

import "github.com/gin-gonic/gin"

func Routes(server *gin.Engine) {
	server.POST("/create-backend", createBackEndNameSpace)
	server.Any("/:namespace", getBackendNameSpace)

}
