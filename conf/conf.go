package conf

import (
	"sync"

	"github.com/srrmendez/private-api-offers/model"
	"github.com/srrmendez/services-interface-tools/pkg/config"
)

var properties *Properties
var syncProperties sync.Once

func GetProps() *Properties {
	syncProperties.Do(func() {
		//dir := "./config/conf.yaml"
		dir := "/var/www/api-offers/config/conf.yaml"

		if err := config.LoadEnvFromYamlFile(dir, &properties); err != nil {
			panic(err)
		}

	})

	return properties
}

type Category struct {
	Category model.CategoryType `yaml:"category"`
	Type     model.OfferType    `yaml:"type"`
}

type Properties struct {
	App struct {
		Path       string `yaml:"appPath"`
		Port       int    `yaml:"port"`
		LogAddress string `yaml:"logAddress"`
	} `yaml:"app"`
	Database struct {
		Host               string `yaml:"host"`
		Port               int    `yaml:"port"`
		Database           string `yaml:"database"`
		Table              string `yaml:"table"`
		SupplementaryTable string `yaml:"supplementaryTable"`
	} `yaml:"database"`
	Categories         map[string]Category `yaml:"categories"`
	PrivateApiTracking struct {
		Host string `yaml:"host"`
	} `yaml:"privateApiTracking"`
}
