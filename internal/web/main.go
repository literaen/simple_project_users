package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/literaen/simple_project/users/internal/app"

	"github.com/gin-gonic/gin"
)

type Web struct {
	app    *app.App
	Server *http.Server
}

func NewWeb() *Web {
	app, err := app.InitApp()
	if err != nil {
		log.Fatal(err)
	}

	return &Web{
		app: app,
	}
}

func (w *Web) Init() *gin.Engine {
	r := gin.Default()

	w.Server = &http.Server{
		Addr:    fmt.Sprintf(":%s", w.app.Config.PORT),
		Handler: r,
	}

	return r
}

func (w *Web) Run() {
	go func() {
		if err := w.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Канал для получения сигнала завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := w.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}
}

func (w *Web) Shutdown(ctx context.Context) error {
	log.Println("Shutting down Gin server...")
	return w.Server.Shutdown(ctx)
}
