/*
 * @Author: tj
 * @Date: 2022-09-20 10:04:42
 * @LastEditors: tj
 * @LastEditTime: 2022-11-23 15:36:11
 * @FilePath: \demo\conf\config.go
 */
package conf

import (
	"github.com/BurntSushi/toml"
)

var (
	cfg = &ServerConfig{}
)

func GetConfigInfo(path string) (*ServerConfig, error) {
	_, err := toml.DecodeFile(path, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
