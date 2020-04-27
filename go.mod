module github.com/limoxi/ghost

go 1.13

require (
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gin-gonic/gin v1.5.0
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/jinzhu/gorm v1.9.11
	github.com/json-iterator/go v1.1.8 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.11 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0
	golang.org/x/sys v0.0.0-20191210023423-ac6580df4449 // indirect
	gopkg.in/go-playground/validator.v9 v9.30.2 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
)

replace github.com/jinzhu/gorm => github.com/limoxi/gorm v1.9.120
