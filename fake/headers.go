/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-03-07 10:06:06
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-03-20 16:21:48
 * @FilePath: \car_hub\fake\headers.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package fake

import (
	"math/rand"

	"github.com/gocolly/colly/v2"
)

var uas = [...]string{
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.111 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.1",
	"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/536.3",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/535.24",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_0) AppleWebKit/536.3",
}

func SetChe168Headers(r *colly.Request) {

	SetHeaders(r, "www.che168.com", "http://www.che168.com", "https://car.autohome.com.cn")
}

func GetUserAgent() string {
	n := rand.Intn(len(uas))
	return uas[n]
}

func SetHeaders(r *colly.Request, host string, origin string, referer string) {

	r.Headers.Set("Host", host)
	r.Headers.Set("Connection", "keep-alive")
	r.Headers.Set("Accept", "*/*")
	r.Headers.Set("Origin", origin)
	r.Headers.Set("Referer", referer)
	r.Headers.Set("Accept-Encoding", "gzip, deflate")
	r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
	r.Headers.Set("User-Agent", GetUserAgent())
}
