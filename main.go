package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    router.GET("/", func(context *gin.Context) {
        context.JSON(http.StatusOK, gin.H {
            "message": "base route",
        })
    })

    router.GET("/kaboom", func(context *gin.Context) {
        context.JSON(http.StatusOK, gin.H {
            "message": "MountVesuvius",
        })
    })

    router.GET("/:val", func(context *gin.Context) {
        val := context.Param("val")

        context.JSON(http.StatusOK, gin.H {
            "message": val,
        })
    })

    router.Run()
}
