### 备忘

- 微信公众号链接： http://mp.weixin.qq.com/s?__biz=MjM5NzAyMDIwMA==&mid=2653672007&idx=1&sn=4274fa8d8f215475b8a50e947c637b78&chksm=bd3f25008a48ac165db889db05bc970d1987356a991f9189f7c723e07929230a629cf761a274#rd 这里面的query_string 都是什么含义？

在微信公众号的文章链接中，四个参数的含义如下：


__biz: 微信公众号的唯一标识符，每个公众号都有自己特定的__biz。

mid: 每篇文章的唯一标识符，可以通过这个标识符找到对应的文章。

idx: 如果一篇文章有多页内容，idx表示当前页面是第几页。

sn: 一篇文章的唯一标识符，与mid不同的是，sn是加密后的标识符。

这4个参数的作用，是在微信公众号的文章统计和推荐功能中起到关键作用的。公众号可以根据这些参数来获取文章的访问数据，比如阅读量、转发量、点赞量等数据，以此来了解文章的受欢迎程度和用户的行为特征，进而优化公众号的内容和营销策略。同时，这些参数也可以用于微信公众号的推荐功能，通过分析用户的访问历史、行为特征等信息，向用户推荐文章


- 相关解释：https://www.zhihu.com/question/29788207

![Alt text](1691403154209.png)

https://mp.weixin.qq.com/s?__biz=MzkyNDM0NjM3OA==&amp;mid=2247548522&amp;idx=1&amp;sn=13b6bebd02d1b510bf3f4d4a05bd0425&amp;chksm=c1d54f0ff6a2c619eeb684262ae7311f1e6033ebb8f762594f88b61a348b7ab7bad8166c9082#rd
https://mp.weixin.qq.com/s?__biz=MzkyNDM0NjM3OA==&amp;mid=2247548522&amp;idx=2&amp;sn=194240c6a96d04d9016c5d4c73680df8&amp;chksm=c1d54f0ff6a2c6190ac5a8a9c4c64c300709f75c388c8d76102a9057bb5c8be2303ce6f55ed6#rd
https://mp.weixin.qq.com/s?__biz=MzkyNDM0NjM3OA==&amp;mid=2247548522&amp;idx=3&amp;sn=b9bf1ef6b2530c473a76f9e5583c6b80&amp;chksm=c1d54f0ff6a2c6191e314fdad39dbf1494c0ca920c55cf6d9de95a16360c19f8c8bfdf2032f6#rd
http://mp.weixin.qq.com/s?__biz=MzkyNDM0NjM3OA==&amp;mid=2247548522&amp;idx=4&amp;sn=8fe063b6eb3f36de99b2baf62ff9944c&amp;chksm=c1d54f0ff6a2c619d5fc5af529221a6f906b7fef182c65df833ab8b3bd0bed1633d3c69747e1#rd


上面是截图内公众号一次推送里面的4篇文章的链接地址，其中:
__biz: 是相同的
mid: 是相同的
idx: 是顺序索引值
sn: 是不同的


