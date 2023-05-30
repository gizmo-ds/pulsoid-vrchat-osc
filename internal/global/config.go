package global

import "github.com/gizmo-ds/pulsoid-vrchat-osc/internal/config"

var Config *config.Config

func LoadConfig(filename string) error {
	conf, err := config.LoadFormFile(filename)
	if err != nil {
		return err
	}
	Config = conf
	return nil
}
