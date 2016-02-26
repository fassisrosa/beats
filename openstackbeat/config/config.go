// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

type Config struct {
	Openstackbeat OpenstackbeatConfig
}

type OpenstackbeatConfig struct {
	Period string `yaml:"period"`
	OpenStackVersion string `yaml:"openStackVersion"`
	OpenStackUrl string `yaml:"openStackUrl"`
}
