/*
 * @Author: tj
 * @Date: 2022-10-10 09:52:23
 * @LastEditors: tj
 * @LastEditTime: 2022-10-25 14:09:02
 * @FilePath: \backend\internal\common\common.go
 */
package common

import (
	"github.com/shopspring/decimal"
)

// DecimalPlacesHandle 小数位数处理（四舍五入）
func DecimalPlacesHandle(data float64, places int) float64 {
	ret, _ := decimal.NewFromFloat(data).Round(int32(places)).Float64()

	return ret
}
