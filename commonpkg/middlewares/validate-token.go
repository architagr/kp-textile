package middlewares

import (
	"commonpkg/models"
	"commonpkg/token"
	"net/http"
	"strings"

	gin "github.com/gin-gonic/gin"
)

func ReturnUnauthorized(context *gin.Context) {
	context.AbortWithStatusJSON(http.StatusUnauthorized, models.CommonResponse{
		Errors: []models.ErrorDetail{
			{
				ErrorCode:    models.ErrorUnauthorized,
				ErrorMessage: "You are not authorized to access this path",
			},
		},
		StatusCode:   http.StatusUnauthorized,
		ErrorMessage: "Unauthorized access",
	})
}

func ValidateTokenMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {

		tokenString := context.Request.Header.Get(models.AuthHeaderName)

		if len(tokenString) > 0 {

			tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
			valid, claims := token.VerifyToken(tokenString)
			if !valid {
				ReturnUnauthorized(context)
			}
			if len(context.Keys) == 0 {
				context.Keys = make(map[string]interface{})
			}
			context.Keys[models.ContextKey_Username] = claims.Username
			context.Keys[models.ContextKey_Roles] = claims.Roles
			context.Next()
		} else {
			ReturnUnauthorized(context)
		}
	}
}
