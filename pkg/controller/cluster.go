package controller

import (
	"github.com/yaoice/istio-crd-tool/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"sort"
	"strings"
)


// @Summary Get Clusters
// @Description 获取cluster列表
// @Success 200 {array} byte OK
// @Failure 401 {object} models.ErrorResponse Unauthorized
// @Failure 404 {object} models.ErrorResponse Not Found
// @Failure 500 {object} models.ErrorResponse Internal Server Error
// @router /clusters [get]
func (i *istioCrdController) GetClusters(c *gin.Context) {
	keys := reflect.ValueOf(i.client).MapKeys()
	strKeys := make(models.ClusterSlice, len(keys))
	for i := 0; i < len(keys); i++ {
		strKeys[i] = keys[i].String()
	}
	sort.Stable(strKeys)
	c.String(http.StatusOK, strings.Join(strKeys, " "))
}