module github.com/limoxi/ghost

go 1.16

require (
	github.com/getsentry/sentry-go v0.6.1
	github.com/gin-gonic/gin v1.7.4
	github.com/jinzhu/gorm v1.9.12
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
)

replace github.com/jinzhu/gorm => github.com/limoxi/gorm v1.9.120
