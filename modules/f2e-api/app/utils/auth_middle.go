package utils

import (
	"net/http"

	h "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func AuthSessionMidd(c *gin.Context) {
	auth, isServiceToken, err := h.SessionChecking(c)
	if !viper.GetBool("skip_auth") {
		if err != nil || auth != true {
			log.Debugf("error: %v, auth: %v", err.Error(), auth)
			c.Set("auth", auth)
			h.JSONR(c, http.StatusUnauthorized, err)
			c.Abort()
			return
		}
	}
	c.Set("auth", auth)
	c.Set("is_service_token", isServiceToken)
	c.Next()
}

func CORS() gin.HandlerFunc {
	return func(context *gin.Context) {
		if context.Request.Method == "OPTIONS" {
			context.AbortWithStatus(200)
		} else {
			context.Next()
		}
	}
}
