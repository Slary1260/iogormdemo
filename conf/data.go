/*
 * @Author: tj
 * @Date: 2022-11-23 15:35:27
 * @LastEditors: tj
 * @LastEditTime: 2022-11-23 16:43:43
 * @FilePath: \demo\conf\data.go
 */
package conf

type ServerConfig struct {
	DatabaseConfigs []*DatabaseCfg
}

type DatabaseCfg struct {
	Port    int
	Mode    string
	Host    string
	User    string
	PassWd  string
	DBName  string
	Charset string
}
