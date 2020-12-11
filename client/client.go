package main

import (
	"bytes"
	"fmt"
	log "github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const (
	version                  = "1.0"
	SERVICE_NAME             = "ice"
	FLAG_KEY_SERVER_HOST     = "server.host"
	FLAG_KEY_SERVER_PORT     = "server.port"
	BASE_PATH 				 = "apis/v1"
)

var (
	exportCS string
	importCS string
	exportNS string
	importNS string
	importFile string
	prefixURL string
	prefixApiURL string
)

var rootCmd = &cobra.Command{
	Use:   "ice",
	Short: "ice api client",
	Long:  "Simple cllient to interact with ice api",
	Run: func(cmd *cobra.Command, args []string) {
		runCmd()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version.",
	Long:  "The version of the dispatch service.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "get cluster list.",
	Long:  "Get the k8s cluster list.",
	Run: func(cmd *cobra.Command, args []string) {
		goCluster()
	},
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export api call istio",
	Long:  "Call the istio api to export crd config.",
	Run: func(cmd *cobra.Command, args []string) {
		cs := exportCS
		ns := exportNS
		goExport(cs, ns)
	},
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import api call istio",
	Long:  "Call the istio api to import crd config.",
	Run: func(cmd *cobra.Command, args []string) {
		cs := importCS
		ns := importNS
		goImport(cs, ns, importFile)
	},
}

// This function will hopefully display a welcome message
// based on the authentication token provided in login

func goPing() {
	url := fmt.Sprintf("%s/ping", prefixURL)
	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func goCluster() {
	url := fmt.Sprintf("%s/clusters", prefixApiURL)
	req, _ := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func goExport(cluster, namespace string) {
	pathPrefix := prefixApiURL + "/clusters/%s/namespaces/%s"
	url := fmt.Sprintf(pathPrefix + "/export", cluster, namespace)
	req, _ := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func goImport(cluster, namespace, path string) {
	currentPath, _ := os.Getwd()
	path = filepath.Join(currentPath, path)

	if _, err := os.Open(path); err != nil {
		panic(err)
	}

	pathPrefix := prefixApiURL + "/clusters/%s/namespaces/%s"
	url := fmt.Sprintf(pathPrefix + "/import", cluster, namespace)

	extraParams := map[string]string{
		"title":       "Istio Crd",
		"author":      "ice",
		"description": "Istio crd k8s yaml config in specific namespace.",
	}
	req, err := newfileUploadRequest(url, extraParams, "crdconfig", path)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

// This function will log you in via Json payload and return an auth token
// if successfull

func runCmd() {
	goPing()
}

func init() {
	viper.SetConfigType("toml")
	viper.SetConfigName(SERVICE_NAME)
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", SERVICE_NAME))   // path to look for the config file in
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s/", SERVICE_NAME)) // call multiple times to add many search paths
	viper.AddConfigPath("./etc/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		fmt.Fprintf(os.Stderr, "Fatal error config file: %s \n", err)
	}

	host := viper.GetString(FLAG_KEY_SERVER_HOST)
	port := viper.GetString(FLAG_KEY_SERVER_PORT)
	if port == "80" {
		prefixURL = fmt.Sprintf("http://%s", host)
	} else {
		prefixURL = fmt.Sprintf("http://%s:%s", host, port)
	}
	prefixApiURL = fmt.Sprintf("%s/%s", prefixURL, BASE_PATH)

	// Adding commands into the client
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(clusterCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(importCmd)

	exportFlags := exportCmd.Flags()
	importFlags := importCmd.Flags()

	exportFlags.StringVarP(&exportCS, "cluster", "c", "","cluster to export")
	exportFlags.StringVarP(&exportNS, "namespace", "n", "","namespace to export")
	viper.BindPFlag("cluster", exportFlags.Lookup("cluster"))
	viper.BindPFlag("namespace", exportFlags.Lookup("namespace"))

	importFlags.StringVarP(&importCS, "cluster", "c", "","cluster to import")
	importFlags.StringVarP(&importNS, "namespace", "n", "","namespace to import")
	importFlags.StringVarP(&importFile, "file", "f", "", "file path to import")
	viper.BindPFlag("cluster", importFlags.Lookup("cluster"))
	viper.BindPFlag("namespace", importFlags.Lookup("namespace"))
	viper.BindPFlag("file", importFlags.Lookup("file"))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
