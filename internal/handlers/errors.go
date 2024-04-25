package handlers

import "github.com/gin-gonic/gin"

var responseUserNotFound = gin.H{"error": "user not found. please login again"}

var responseInvalidReqBody = gin.H{"error": "invalid request payload"}

var responseAccountNotFound = gin.H{"error": "acco not found"}

var responseAlreadyLoggedIn = gin.H{"error": "you are already logged in"}

func notFoundResponse(resource string) gin.H {
	return gin.H{"error": resource + " not found"}
}
