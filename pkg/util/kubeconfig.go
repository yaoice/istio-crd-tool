// Copyright 2018 Naftis Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"github.com/yaoice/istio-crd-tool/pkg/config"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"path/filepath"
)

// Kubeconfig returns kube config list path.
func getKubeConfigsPath() string {
	if kubeConfig := config.GetString(config.FLAG_KEY_ISTIO_KUBECONFIGS); kubeConfig != "" {
		return kubeConfig
	}
	return filepath.Join("/etc", "multiple_k8s_clusters.yaml")
}

func GetKubeConfigDir() string {
	return filepath.Dir(getKubeConfigsPath())
}

func GetKubeConfigsMap() *KubeConfigs {
	var kcs KubeConfigs
	path := getKubeConfigsPath()
	kubeConfigsBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(kubeConfigsBytes, &kcs)
	if err != nil {
		panic(err)
	}
	return &kcs
}

