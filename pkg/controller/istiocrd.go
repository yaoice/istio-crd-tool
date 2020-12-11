package controller

import (
	"bytes"
	"fmt"
	"github.com/yaoice/istio-crd-tool/pkg/config"
	"github.com/yaoice/istio-crd-tool/pkg/log"
	"github.com/yaoice/istio-crd-tool/pkg/models"
	"github.com/yaoice/istio-crd-tool/pkg/util"
	"github.com/ghodss/yaml"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"istio.io/istio/pilot/pkg/config/kube/crd"
	istiomodel "istio.io/istio/pilot/pkg/model"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	once sync.Once
	icc  *istioCrdController
)

type IstioCrdInterface interface {
	Export(c *gin.Context)
	Import(c *gin.Context)
}

type istioCrdController struct {
	client map[string]*crd.Client
	rwm *sync.RWMutex
}

func NewIstioCrdController() *istioCrdController {
	once.Do(func() {
		client := make(map[string]*crd.Client)
		kcs := utils.GetKubeConfigsMap()
		for _, kc := range kcs.Clusters {
			kubePath := filepath.Join(utils.GetKubeConfigDir(), kc.Name + "_kubeconfig.yaml")
			_, err := os.Stat(kubePath)
			if err != nil {
				log.Errorf("kubePath %s doesn't exist. %v", kubePath, err.Error())
				continue
			}
			c, err := crd.NewClient(kubePath, "", istiomodel.IstioConfigTypes, "")
			if err != nil {
				log.Errorf("init istiocrd fail %v", err.Error())
				continue
			}
			client[kc.Name] = c
		}
		icc = &istioCrdController{
			client: client,
			rwm: new(sync.RWMutex),
		}
	})
	return icc
}

// @Summary Export
// @Description Istio crd导出
// @Param   cluster   path    string  true     "the k8s cluster of istio crd"
// @Param   namespace   path    string  true     "the namespace of istio crd"
// @Success 200 {array} byte OK
// @Failure 401 {object} models.ErrorResponse Unauthorized
// @Failure 404 {object} models.ErrorResponse Not Found
// @Failure 500 {object} models.ErrorResponse Internal Server Error
// @router /clusters/{cluster}/namespaces/{namespace}/export [get]
func (i *istioCrdController) Export(c *gin.Context) {
	cluster := c.Param("cluster")
	namespace := c.Param("namespace")
	crds := config.GetString(config.FLAG_KEY_ISTIO_CRD)

	var configList []istiomodel.Config

	for _, crd := range strings.Split(crds, ",") {
		crdList, err := i.listCrds(crd, cluster, namespace)
		if err != nil {
			message := fmt.Sprintf("istio list %s error: %v", crd, err.Error())
			log.Errorf(message)
			code := http.StatusInternalServerError
			c.JSON(
				code,
				models.ErrorResponse{
					Code:    code,
					Message: message,
				},
			)
			return
		}
		configList = append(configList, crdList...)
	}
	c.String(http.StatusOK, i.yamlOutput(configList, cluster))
}

// @Summary Import
// @Description Istio crd导入
// @Param   cluster   path    string  true     "the k8s cluster of istio crd"
// @Param   namespace   path    string  true     "the namespace of istio crd"
// @Produce multipart/form-data
// @Success 200 {array} byte OK
// @Failure 401 {object} models.ErrorResponse Unauthorized
// @Failure 404 {object} models.ErrorResponse Not Found
// @Failure 500 {object} models.ErrorResponse Internal Server Error
// @router /clusters/{cluster}/namespaces/{namespace}/import [post]
func (i *istioCrdController) Import(c *gin.Context) {
	cluster := c.Param("cluster")
	namespace := c.Param("namespace")
	crdConfig, err := c.FormFile("crdconfig")
	if err != nil {
		message := fmt.Sprintf("import crdconfig error: %v", err.Error())
		log.Errorf(message)
		code := http.StatusBadRequest
		c.JSON(
			code,
			models.ErrorResponse{
				Code:    code,
				Message: message,
			},
		)
		return
	}

	f, err := crdConfig.Open()
	if err != nil {
		code := http.StatusInternalServerError
		message := fmt.Sprintf("Could not open file %s", crdConfig.Filename)
		c.JSON(
			code,
			models.ErrorResponse{
				Code:    code,
				Message: message,
			})
		return
	}

	crdBytes, err := ioutil.ReadAll(f)
	if err != nil {
		code := http.StatusInternalServerError
		message := fmt.Sprintf("Could not read file %s", crdConfig.Filename)
		c.JSON(
			code,
			models.ErrorResponse{
				Code:    code,
				Message: message,
			})
		return
	}

	configList, _, err := crd.ParseInputs(string(crdBytes))
	if err != nil {
		code := http.StatusInternalServerError
		message := fmt.Sprintf("crd ParseInputs error: %v", err.Error())
		c.JSON(
			code,
			models.ErrorResponse{
				Code:    code,
				Message: message,
			})
		return
	}

	count := 0
	errSlice := make([]error, 10)
	wg := sync.WaitGroup{}
	wg.Add(len(configList))
	for _, config := range configList {
		go func(config istiomodel.Config) {
			defer wg.Done()
			config.Namespace = namespace
			config.CreationTimestamp = time.Now().Local()
			rev, err := i.client[cluster].Create(config)
			if err != nil {
				current := i.client[cluster].Get(config.Type, config.Name, config.Namespace)
				config.ResourceVersion = current.ResourceVersion
				// clear resourceVersion for rollback
				current.ResourceVersion = ""
//				config.ResourceVersion = utils.CreateCaptcha()
				updateRev, err := i.client[cluster].Update(config)
				if err != nil {
					i.rwm.Lock()
					// if the config create fail, break loop and return error
					message := fmt.Sprintf("Created/Updated config error: %v", err.Error())
					log.Errorf(message)
					errSlice = append(errSlice, err)
					count++
					i.rwm.Unlock()
					return
				} else {
					log.Infoln("Updated config success", "key", config.Key(), "revision", updateRev, "config", config)
				}
			} else {
				log.Infoln("Created config success", "key", config.Key(), "revision", rev, "config", config)
			}
		}(config)
	}
	wg.Wait()

	if count > 0 {
		c.JSON(
			http.StatusInternalServerError,
			models.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("create config error in namespace %s, error: %v", namespace, errSlice),
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		models.ErrorResponse{
			Code:    http.StatusOK,
			Message: fmt.Sprintf("create config success in namespace %s", namespace),
		},
	)
}

func (i *istioCrdController) listCrds(tpe, cluster, namespace string) (configList []istiomodel.Config, err error) {
	client := i.client[cluster]
	switch tpe {
	case "virtualservice":
		{
			configList, err = client.List(istiomodel.VirtualService.Type, namespace)
		}
	case "serviceentry":
		{
			configList, err = client.List(istiomodel.ServiceEntry.Type, namespace)
		}
	case "destinationrule":
		{
			configList, err = client.List(istiomodel.DestinationRule.Type, namespace)
		}
	case "gateway":
		{
			configList, err = client.List(istiomodel.Gateway.Type, namespace)
		}
	}
	return configList, err
}

func (i *istioCrdController) GetClient() map[string]*crd.Client {
	return i.client
}


func (i *istioCrdController) yamlOutput(configList []istiomodel.Config, cluster string) string {
	buf := bytes.NewBuffer(make([]byte, 0))
	descriptor := i.client[cluster].ConfigDescriptor()
	for _, config := range configList {
		schema, exists := descriptor.GetByType(config.Type)
		if !exists {
			log.Errorf("Unknown kind %q for %v", crd.ResourceName(config.Type), config.Name)
			continue
		}
		obj, err := crd.ConvertConfig(schema, config)
		if err != nil {
			log.Errorf("Could not decode %v: %v", config.Name, err)
			continue
		}
		obj.SetObjectMeta(meta_v1.ObjectMeta{
			Name:        config.Name,
			Labels:      config.Labels,
			CreationTimestamp: meta_v1.Time{},
		})

		bytes, err := yaml.Marshal(obj)
		if err != nil {
			log.Errorf("Could not convert %v to YAML: %v", config, err)
			continue
		}

		buf.Write(bytes)
		buf.WriteString("---\n")
	}
	return buf.String()
}
