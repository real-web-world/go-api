package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/real-web-world/go-api/bootstrap"
	"github.com/real-web-world/go-api/global"
	"github.com/real-web-world/go-api/routes"
)

// @title buffge blog API
// @version 1.0
// @description blog api

// @contact.name buffge
// @contact.url https://github.com/buffge/
// @contact.email buffge.com@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @tag.name user
// @tag.description 用户

// @tag.name aliPay
// @tag.description 支付宝

// @host blog-api.buffge.com
// @BasePath /
func main() {
	engine := bootstrap.InitApp()
	routes.RouterSetup(engine)
	httpAddr := fmt.Sprintf("%s:%d", global.Conf.HTTPHost,
		global.Conf.HTTPPort)
	srv := &http.Server{
		Addr:    httpAddr,
		Handler: engine,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("server shutdown err:", err)
		return
	}
	log.Println("server shutdown success")
}
