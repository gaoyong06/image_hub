package utils

import (
	"net/url"

	log "github.com/sirupsen/logrus"
)

// 获取从文章内容页面meta内获取到的文章url的query string参数
// <meta content="http://mp.weixin.qq.com/s?__biz=MjM5NzAyMDIwMA==&amp;mid=2653562471&amp;idx=1&amp;sn=5a209eca9a0c9d92d484dadfa516a807&amp;chksm=bd3ed1208a49583679dddb80f504983511b6bc9d63c89242dd3df68daebd587a78b8fea1afa0#rd"/>
func GetArticleUrlQueryParams(ogUrl string) (url.Values, error) {

	decodedUrl, err := url.QueryUnescape(ogUrl)
	if err != nil {
		log.Errorf("url.QueryUnescape failed. err: %+v\n", err)
		return nil, err
	}
	u, err := url.Parse(decodedUrl)
	if err != nil {
		log.Errorf("url.Parse failed. err: %+v\n", err)
		return nil, err
	}

	queryParams := u.Query()
	return queryParams, nil
}
