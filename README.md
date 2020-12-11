# istio-crd-tool

## 简介

istio-crd-tool是istio crd yaml文件导入、导出的一个工具，基于go-gin实现rest api，基于cobra实现命令行

## 如何使用istio-crd-tool

### 配置文件

在etc目录下

### api启动程序

cmd/main.go

### 命令行

先编译
```
go build -o ice client/client.go
```

查看命令行帮助
```
./istio-crd-tool --help
Simple cllient to interact with ice api

Usage:
  istio-crd-tool [flags]
  istio-crd-tool [command]

Available Commands:
  export      export api call istio
  help        Help about any command
  import      import api call istio
  version     Print the version.

Flags:
  -h, --help   help for ice

Use "istio-crd-tool [command] --help" for more information about a command.

```

查看export命令行帮助
```
./istio-crd-tool export --help
Call the istio api to export crd config.

Usage:
  istio-crd-tool export [flags]

Flags:
  -h, --help               help for export
  -n, --namespace string   namespace to export
```
  
查看import命令行帮助  
```
$ ./istio-crd-tool import --help
Call the istio api to import crd config.

Usage:
  istio-crd-tool import [flags]

Flags:
  -f, --file string        file path to import
  -h, --help               help for import
  -n, --namespace string   namespace to import
```

### 生成swagger API Doc
```
swag init -g pkg/route/routes.go
```