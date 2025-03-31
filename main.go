package main

import (
	"go-notify/dto"
	"go-notify/jobs"
	"go-notify/mail"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	_ "github.com/joho/godotenv"
)

func main() {
	var err error
	// err = godotenv.Load(".env")
	// if err != nil {
	// 	panic(err)
	// }

	mail.M = mail.NewMail()
	mail.M.SetAuth()

	var wg = &sync.WaitGroup{}
	wp := jobs.NewWorkerPool(5, 25)

	wg.Add(int(wp.Worker))

	wp.Do(wg)

	r := gin.Default()
	r.POST("/api/mail", func(ctx *gin.Context) {
		var req dto.Body
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, dto.Response{
				Message: "error while binding body into struct",
				Data: map[string]interface{}{
					"error": err.Error(),
				},
			})
			return
		}
		wp.Add(req)
		ctx.JSON(200, gin.H{
			"message": "ok",
		})
	})

	r.POST("/api/mail/v2", func(ctx *gin.Context) {
		var req dto.Body
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, dto.Response{
				Message: "error while binding body into struct",
				Data: map[string]interface{}{
					"error": err.Error(),
				},
			})
			return
		}
		err = mail.M.SendMail(req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dto.Response{
				Message: "failed to send email",
				Data: map[string]interface{}{
					"error": err.Error(),
				},
			})
			return
		}
		ctx.JSON(200, gin.H{
			"message": "ok",
		})
	})

	r.Run(":8080")
	wg.Wait()
}
