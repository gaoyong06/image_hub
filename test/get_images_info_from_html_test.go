package test

import (
	"fmt"
	"image_hub/spiders"
	"regexp"
	"testing"
)

func TestGetImagesInfoFromHTML(t *testing.T) {

	htmlUrl := ""
	htmlStr := `
		<section
		donone="shifuMouseDownPic('shifu_p_028')"
		label="Copyright Reserved by PLAYHUDONG."
		style="
		background: #ffffff;
		text-align: center;
		border-style: none;
		clear: both;
		overflow: hidden;
		margin: 1em auto;
		"
	>
		<img
		class=""
		data-ratio="0.9525"
		data-src="http://mmbiz.qpic.cn/mmbiz_jpg/8GsU5rmbW2wicKkuOibVs7DYTkZUC8xsFibeXKickmqwzxiaPtpXnuj0fYjDialJk3gYh7MUFJuTnMg0FJiceCnUibUBzg/0?wx_fmt=jpeg"
		data-type="jpeg"
		data-w="1200"
		src="D:/work/wechat_download_data/images\2017/08/07/185349/20170807_185349_1_5_kz20.jpeg"
		style="margin-right: 2%; width: 49% !important"
		title="http://7xo6kd.com1.z0.glb.clouddn.com/upload-ueditor-image-20170807-1502101072587016207.jpg"
		/><img
		class=""
		data-ratio="0.9525"
		data-src="http://mmbiz.qpic.cn/mmbiz_jpg/8GsU5rmbW2wicKkuOibVs7DYTkZUC8xsFibN75quTLNef2Nw1WD5gCC5jC5lt9tEW44asH61GOL7egMnGJkezUKxg/0?wx_fmt=jpeg"
		data-type="jpeg"
		data-w="1200"
		src="D:/work/wechat_download_data/images\2017/08/07/185349/20170807_185349_1_6_5rc9.jpeg"
		style="width: 49% !important"
		title="http://7xo6kd.com1.z0.glb.clouddn.com/upload-ueditor-image-20170807-1502101062894038958.jpg"
		/>
	</section>
	`
	imgRegex, err := regexp.Compile(`<\s*img[^>]*src\s*=\s*["']?([^"']+)["']?[^>]*>`)
	if err != nil {
		panic("failed to compile imgRegex: " + err.Error())
	}

	imgTags := imgRegex.FindAllString(htmlStr, -1)
	fmt.Printf("\n========================= imgTags ==========================\n\n %#v\n", imgTags)

	imgs, filteredImgs, err := spiders.InferImageTypeFromHTML(htmlUrl, htmlStr)
	if err != nil {
		panic("spiders.GetImagesInfoFromHTML falied. err : " + err.Error())
	}
	fmt.Printf("\n========================= imgs ==========================\n\n %#v\n", imgs)
	fmt.Printf("\n========================= filteredImgs ==========================\n\n %#v\n", filteredImgs)

}

func TestGetImagesInfoFromHTML1(t *testing.T) {

	htmlStr := `<section donone="shifuMouseDownPic('shifu_p_028')" label="Copyright Reserved by PLAYHUDONG." style="background:#ffffff; text-align: center;border-style: none; clear: both; overflow: hidden;margin: 1em auto;"><img class="" data-ratio="0.9525" data-src="http://mmbiz.qpic.cn/mmbiz_jpg/8GsU5rmbW2wicKkuOibVs7DYTkZUC8xsFibeXKickmqwzxiaPtpXnuj0fYjDialJk3gYh7MUFJuTnMg0FJiceCnUibUBzg/0?wx_fmt=jpeg" data-type="jpeg" data-w="1200" src="D:/work/wechat_download_data/images\2017/08/07/185349/20170807_185349_1_5_kz20.jpeg" style="
	margin-right: 2%;
	width: 49% !important;" title="http://7xo6kd.com1.z0.glb.clouddn.com/upload-ueditor-image-20170807-1502101072587016207.jpg"><img class="" data-ratio="0.9525" data-src="http://mmbiz.qpic.cn/mmbiz_jpg/8GsU5rmbW2wicKkuOibVs7DYTkZUC8xsFibN75quTLNef2Nw1WD5gCC5jC5lt9tEW44asH61GOL7egMnGJkezUKxg/0?wx_fmt=jpeg" data-type="jpeg" data-w="1200" src="D:/work/wechat_download_data/images\2017/08/07/185349/20170807_185349_1_6_5rc9.jpeg" style="
	width: 49% !important;" title="http://7xo6kd.com1.z0.glb.clouddn.com/upload-ueditor-image-20170807-1502101062894038958.jpg"></section>>`

	imgRegex, err := regexp.Compile(`<\s*img[^>]*src\s*=\s*["']?([^"']+)["']?[^>]*>`)
	if err != nil {
		panic("failed to compile imgRegex: " + err.Error())
	}

	imgTags := imgRegex.FindAllString(htmlStr, -1)

	fmt.Printf("\n========================= imgTags ==========================\n\n %#v\n", imgTags)

}
