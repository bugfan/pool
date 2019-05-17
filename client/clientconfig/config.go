package clientconfig

import (
	"fmt"
	"io/ioutil"
	"runtime"

	"github.com/olebedev/config"
)

var cfg *config.Config

var default_config = `
server:
  server: ""
  client_id: ""
`

var clientConf = "/etc/client.yml"

func init() {
	var err error
	if runtime.GOOS == "darwin" {
		clientConf = "config.yml"
	}
	cfg, err = config.ParseYamlFile(clientConf)
	if err != nil {
		fmt.Print("config parse error or not exists, use default")
		cfg, _ = config.ParseYaml(default_config)
	}
}

func String() string {
	yml, _ := config.RenderYaml(cfg.Root)
	return yml
}

func Save() {
	yml, _ := config.RenderYaml(cfg.Root)
	d := []byte(yml)
	ioutil.WriteFile(clientConf, d, 0644)
}

func Get(path string) (string, error) {
	return cfg.String(path)
}

func Set(path, val string) error {
	return cfg.Set(path, val)
}
