package utils

import (
	"net/http"

	h "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/gin-gonic/gin"
)

func AuthSessionMidd(c *gin.Context) {
	auth, err := h.SessionChecking(c)
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
