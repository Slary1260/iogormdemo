/*
 * @Author: tj
 * @Date: 2022-11-23 16:07:49
 * @LastEditors: tj
 * @LastEditTime: 2022-11-23 16:21:38
 * @FilePath: \demo\internal\dbgorm\mode.go
 */
package dbgorm

type DbMode string

const (
	Mastermode DbMode = "master" // 主库
	Slavemode  DbMode = "slave"  // 从库
)
