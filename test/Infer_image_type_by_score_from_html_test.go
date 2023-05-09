package test

import (
	"fmt"
	"image_hub/spiders"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
)

// Test spiders.InferImageTypeByScoreFromHTML function
// TODO: 还有bug
func TestInferImageTypeByScoreFromHTML(t *testing.T) {

	fmt.Println("======================= RUN TestInferImageTypeByScoreFromHTML ======================")

	// 后面一定要加反斜杠
	directoryPath := "D:/work/wechat_download_data/html/test5/"

	// 新建测试文件的前缀
	updatedFilePrefix := "updated_"

	// Read all HTML files in the directory
	fileList, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	// Loop through each file
	for _, file := range fileList {
		// Check if the file is an HTML file
		if filepath.Ext(file.Name()) == ".html" || !strings.HasPrefix(file.Name(), updatedFilePrefix) {
			// Add to the wait group
			wg.Add(1)
			go func(file os.FileInfo) {
				defer wg.Done()
				// Open the file

				htmlUrl := directoryPath + file.Name()

				htmlFile, err := os.Open(htmlUrl)
				if err != nil {
					panic(err)
				}
				defer htmlFile.Close()

				// Read the file into a string
				htmlBytes, err := ioutil.ReadAll(htmlFile)
				if err != nil {
					panic(err)
				}
				htmlStr := string(htmlBytes)

				// Get image info from HTML using the GetImageInfoFromHTML function
				imgsInfo, filteredImgs, err := spiders.InferImageTypeFromHTML(htmlUrl, htmlStr)
				if err != nil {
					panic(err)
				}

				fmt.Printf("=================fileName: %s,  len(imgsInfo): %d, filteredImgs: %#v \n", file.Name(), len(imgsInfo), filteredImgs)

				// Lock the mutex
				mu.Lock()

				// Loop through each image to add the custom class and border
				for _, img := range imgsInfo {
					// fmt.Printf("=================fileName: %s,  imgInfo: %#v\n", file.Name(), img)

					// Escape any regex special characters in img["src"]
					imgSrc := regexp.QuoteMeta(img["src"].(string))

					// Generate regex pattern to match the img tag with any attributes in any order
					imgRegex := regexp.MustCompile(fmt.Sprintf(`(?i)<img.*?src=['"]%s['"].*?>`, imgSrc))

					// Convert the img size from bytes to kb
					sizeKB := fmt.Sprintf("%.2f", img["size"].(float64)/1024)

					// Generate overlay text displaying image information in the top right corner
					overlayText := fmt.Sprintf(`<div style="position:relative">

								<img src="%s" class="custom-class" style="max-width: 100%%; height: auto; display: block;">
									<div style="position:absolute; top: 0; right: 0; padding: 10px; background-color: rgba(0,255,0,0.5); color: white; font-size: 12px; font-weight: bold;">
										Ratio: %f<br>Width: %f<br>Height: %f<br>Format: %s<br>Type: %s<br>Shape: %s<br>Size: %s kb
									</div>
							</div>`, img["src"], img["ratio"], img["width"], img["height"], img["format"], img["type"], img["shape"], sizeKB)

					// The ReplaceAllStringFunc() method replaces each matched image tag with the above overlay text
					htmlStr = imgRegex.ReplaceAllStringFunc(htmlStr, func(match string) string {
						return overlayText
					})

				}
				// Unlock the mutex
				mu.Unlock()

				// Create the updated file with a path relative to the original file
				updatedFilePath := filepath.Dir(directoryPath+file.Name()) + "/" + updatedFilePrefix + file.Name()
				updatedFile, err := os.Create(updatedFilePath)
				if err != nil {
					panic(err)
				}
				// Write the updated HTML string to the file
				_, err = updatedFile.WriteString(htmlStr)
				if err != nil {
					panic(err)
				}

				updatedFile.Close()
				fmt.Printf("===== success. updatedFile: %s\n", updatedFilePath)

			}(file)

		}
	}
	// Wait for all goroutines to finish
	wg.Wait()
}
