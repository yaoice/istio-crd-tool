package utils

type Cluster struct {
	Name          string `yaml:"name"`
	APIServerHost string `yaml:"apiServerHost"`
}

type KubeConfigs struct {
	Clusters []Cluster `yaml:"clusters"`
}
