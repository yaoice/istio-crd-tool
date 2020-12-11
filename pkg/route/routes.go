package route

import (
	_ "github.com/yaoice/istio-crd-tool/docs"
	"github.com/yaoice/istio-crd-tool/pkg/controller"
	"github.com/yaoice/istio-crd-tool/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"os"
)

// @title Swagger ice
// @version 1.0
// @description This is a ice server.
// @contact.name iceyao
// @contact.url https://www.xxx.com
// @contact.email xiabingyao@tencent.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /apis/v1
func InstallRoutes(r *gin.Engine) {
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// a ping api test
	r.GET("/ping", controller.Ping)

	// config reload
	r.Any("/-/reload", func(c *gin.Context) {
		log.Info("===== Server Stop! Cause: Config Reload. =====")
		os.Exit(1)
	})

	rootGroup := r.Group("/apis/v1")
	{
		pathPrefix := "/clusters/:cluster/namespaces/:namespace"
		// for test
		icc := controller.NewIstioCrdController()
		rootGroup.GET("/clusters", icc.GetClusters)
		rootGroup.GET(pathPrefix + "/export", icc.Export)
		rootGroup.POST(pathPrefix + "/import", icc.Import)
	}
}
