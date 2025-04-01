package main

import (
	"go-notify/dto"
	"go-notify/jobs"
	"go-notify/mail"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
	// _ "github.com/joho/godotenv"
)

func main() {
	var err error
	err = godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	mail.M = mail.NewMail()
	mail.M.SetAuth()

	var wg = &sync.WaitGroup{}
	// You can use other settings for the number of workers and queues.
	// Experiment and adjust to your device
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
		log.Printf("binding sender:%s receiver:%s\n", req.Sender, req.Receiver)
		wp.Add(req, 0)
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
