package main

import (
	"fmt"
	"github.com/yaoice/istio-crd-tool/pkg/config"
	"github.com/yaoice/istio-crd-tool/pkg/log"
	"github.com/yaoice/istio-crd-tool/pkg/route"
	"github.com/yaoice/istio-crd-tool/pkg/util"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.SetupDumpStacksTrap()
	r := gin.Default()
	m := config.GetString(config.FLAG_KEY_GIN_MODE)
	gin.SetMode(m)

	route.InstallRoutes(r)
	serverBindAddr := fmt.Sprintf("%s:%d", config.GetString(config.FLAG_KEY_SERVER_HOST), config.GetInt(config.FLAG_KEY_SERVER_PORT))
	log.Infof("Run server at %s", serverBindAddr)
	r.Run(serverBindAddr) // listen and serve
}
