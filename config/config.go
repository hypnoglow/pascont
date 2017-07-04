package config

import (
	"encoding/json"
	"io"
)

// Config is a struct holding configuration data.
type Config struct {
	Socket   configSocket   `json:"socket"`
	Database configDatabase `json:"database"`
	JWT      configJWT      `json:"jwt"`
	Session  configSession  `json:"session"`
}

type configSocket struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type configDatabase struct {
	DriverName       string              `json:"driver_name"`
	User             string              `json:"user"`
	Password         string              `json:"password"`
	Host             string              `json:"host"`
	Port             string              `json:"port"`
	DatabaseName     string              `json:"database_name"`
	ConnectionParams map[string][]string `json:"connection_params"`
}

type configJWT struct {
	Issuer    string `json:"issuer"`
	SecretKey string `json:"secret_key"`
}

type configSession struct {
	SecretKey string `json:"secret_key"`
}

// FromJSON returns a Config with data read from r as json.
func FromJSON(r io.Reader) Config {
	conf := Config{}

	d := json.NewDecoder(r)
	if err := d.Decode(&conf); err != nil {
		panic(err)
	}

	return conf
}
