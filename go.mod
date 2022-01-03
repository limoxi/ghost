module github.com/limoxi/ghost

go 1.16

require (
	github.com/getsentry/sentry-go v0.6.1
	github.com/gin-gonic/gin v1.7.7
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
	gorm.io/driver/mysql v1.2.2
	gorm.io/gorm v1.22.4
)

//replace gorm.io/gorm => github.com/limoxi/gorm latest
