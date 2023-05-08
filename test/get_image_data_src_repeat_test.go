package test

import (
	"fmt"
	"image_hub/spiders"
	"testing"
)

// 读取directoryPath所有html的文件将各个文件中的img标签的data-src内的值取出来，如果值重复出现的(出现次数>1)。则打印出来
func TestGetImageDataSrcRepeat(t *testing.T) {

	// 后面一定要加反斜杠
	// directoryPath := "D:/work/wechat_download_data/html/test6/"
	directoryPath := "D:/work/wechat_download_data/html/test7/"

	dataSrcRepeat, err := spiders.GetImageDataSrcRepeat(directoryPath)

	if err != nil {
		panic(err)
	}
	// Print the Repeat of data-src values
	fmt.Printf("\n\n ============= dataSrcRepeat======================= \n\n%#v\n\n", dataSrcRepeat)
}
