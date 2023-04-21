package spiders

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

// 测试"1.5万"转为15000
func TestGetMileage(t *testing.T) {

	mileageStrSlice := []string{"0万", "1万", "1.5万", "33.4万"}

	for _, mileageStr := range mileageStrSlice {

		newMileageStr := strings.Replace(mileageStr, "万", "", 1)
		fmt.Println("newMileageStr: ", newMileageStr)

		mileage, err := strconv.ParseFloat(newMileageStr, 64)
		fmt.Println("mileage: ", mileage)

		if err != nil {
			fmt.Println("err: ", err)
		}

		if err == nil {
			fmt.Println(mileage * 10000)
		}
	}
}
