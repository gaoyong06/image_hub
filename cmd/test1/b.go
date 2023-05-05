package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type ImageType int

const (
	AvatarImage    ImageType = iota
	WallpaperImage           // Use implicit assignment here
	WeChatBackgroundImage
	WeChatEmoji
)

// 对磁盘上的图片文件，进行判断，如果图片的尺寸，大小，不能做社交网络的头像，手机壁纸，微信朋友的背景图，微信的表情包，则返回false, 反正返回true
func isImageValid(imgPath string, imgType ImageType) (bool, string, error) {
	// Open the image file
	f, err := os.Open(imgPath)
	if err != nil {
		log.Printf("Failed to open image file: %v\n", err)
		return false, "failed to open file", err
	}
	defer f.Close()

	// Decode the image to get its dimensions
	img, _, err := image.DecodeConfig(f)
	if err != nil {
		log.Printf("Failed to decode image: %v\n", err)
		return false, "failed to decode image", err
	}

	// Get the file size
	fi, err := f.Stat()
	if err != nil {
		log.Printf("Failed to get file info: %v\n", err)
		return false, "failed to get file info", err
	}
	fileSize := fi.Size()

	return checkImageValid(img, fileSize, imgPath, imgType)
}

func checkImageValid(img image.Config, fileSize int64, imgPath string, imgType ImageType) (bool, string, error) {
	var (
		valid      bool
		reason     string
		dimensions string
		sizeLimit  int64
	)

	// Check the dimensions
	switch imgType {
	case AvatarImage:
		dimensions = "200x200"
		sizeLimit = 1 * 1024 * 1024
		if img.Width < 200 || img.Height < 200 {
			valid = false
			reason = fmt.Sprintf("dimension < %s", dimensions)
		}
	case WallpaperImage:
		dimensions = "1080x1920"
		sizeLimit = 5 * 1024 * 1024
		if img.Width < 1080 || img.Height < 1920 {
			valid = false
			reason = fmt.Sprintf("dimension < %s", dimensions)
		}
	case WeChatBackgroundImage:
		dimensions = "1080x1920"
		sizeLimit = 5 * 1024 * 1024
		if img.Width < 1080 || img.Height < 1920 {
			valid = false
			reason = fmt.Sprintf("dimension < %s", dimensions)
		}
	case WeChatEmoji:
		dimensions = "96x96"
		sizeLimit = 64 * 1024
		if img.Width != 96 || img.Height != 96 {
			valid = false
			reason = fmt.Sprintf("dimension != %s", dimensions)
		}
	default:
		return false, "", fmt.Errorf("unknown image type")
	}

	// Check the file size
	if fileSize > sizeLimit {
		valid = false
		reason = fmt.Sprintf("file size > %d", sizeLimit)
	}

	return valid, reason, nil
}

func processImagePaths(imagePath interface{}, imgType ImageType) (chan bool, chan bool) {
	wait := make(chan bool)
	done := make(chan bool)

	// Handle processing of images concurrently
	go func() {
		var wg sync.WaitGroup
		errCh := make(chan error)

		errLog := filepath.Join(filepath.Dir(fmt.Sprintf("%v", imagePath)), fmt.Sprintf("error_%s.log", time.Now().Format("2006-01-02_15-04-05")))
		f, err := os.Create(errLog)
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		switch imagePath.(type) {
		case string: // Single file
			if isValid, reason, err := isImageValid(fmt.Sprintf("%v", imagePath), imgType); err != nil {
				log.Printf("Failed to validate image %s: %s", imagePath, err)
			} else if !isValid {
				msg := fmt.Sprintf("%s: %s\n", imagePath, reason)
				_, err := f.WriteString(msg)
				if err != nil {
					log.Printf("error writing to log file: %v\n", err)
				}
			} else {
				fmt.Printf("%s is valid\n", imagePath) // Handle image processing here...
			}
		case []interface{}: // Directory or array of files
			for _, p := range imagePath.([]interface{}) {
				processImagePaths(p, imgType) // Recursively call function for processing subdirectories and individual files
			}
		default:
			log.Fatal("Invalid input type")
		}

		wg.Wait()
		close(errCh)
		done <- true

		// Check for errors
		if err := <-errCh; err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for processing to start before continuing
	go func() {
		for {
			select {
			case <-wait:
				return
			default:
				// Sleep to avoid spinning
				// Could be changed with a WaitGroup, but I'll just use this for now.
				time.Sleep(time.Millisecond * 20)
			}
		}
	}()

	// Return the channels for use
	return wait, done
}

func main() {
	imagePath := []interface{}{"images", "subdirectory", "C:/Users/Username/Pictures/image.png", "C:/Users/Username/Pictures/image.jpg"}

	wait, done := processImagePaths(imagePath, AvatarImage)
	wait <- true
	<-done
}
