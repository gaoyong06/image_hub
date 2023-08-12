### 内容服务（Image Hub）

|          |                                 |
| -------- | ------------------------------- |
| 作者     | [高勇](mailto:gaoyong06@qq.com) |
| 创建时间 | 2023-04-19 08:35:00             |

#### 介绍

imageHub是一个图片内容源服务，主要服务 content_service 的内容源采集 现在的目标是采集,主流的微信公众号,主流的图片网站内容


#### 工作原理
1. 读取directoryPath所有html的文件将各个文件中的img标签的data-src内的值取出来如果重复出现(出现次数大于1),则记录到变量params["dataSrcRepeat"]中

2. html文件名规则为："%Y%m%d_%H%M%S"_"序号.html", 例如: 20230109_111900_1.html

3. 新建一个colly queue队列

4. 读取html的内容, 判断该html内的图片类型(头像，壁纸，背景图，表情包)是头像,背景,套图,壁纸,表情 哪一种

5. 新建一个onePageSpider设置params["path"]和params["dataSrcRepeat"],添加到上述colly queue处理队列中

6. 队列的c.OnHTML中使用onePageSpider.Process 处理队列中的各个任务

7. 调用one_page.go中的ParseData方法将html字符串解析到Article结构体

8. 在func_map.go中定义了各个微信号的自定义处理函数，调用wechat_微信号.go(如：wechat_touxiangshe.go) 对特殊的微信公众号的的sections过解析处理

9. 通过onePageSpider.Process调用base_spider.go中的Process方法将上述解析到的Article和sections保存到db,支持重复覆盖方式写入


#### 新公众号内容导入步骤(新方法)
```
go run main.go -c ../../configs/config.yaml -d D:/work/wechat_download_data/html/Dump-0511-08-43-21/
```

1. 上面这个命令会自动在image_hub库中字段创建公众号对应的数据表，如果数据表不存在的话
2. 使用下面的shell命令,将image_hub中所有tbl_article_后缀名, 合并到数据表tbl_article
```
./merge_article_table.sh
```
3. 注意: tbl_article_后缀名 中,后缀名为微信号，但是如果微信号中如果包含中划线"-", 会将中划线"-"替换为下划线"_"


#### 新公众号内容导入步骤(旧方法)

1. 查看工作号"微信号", 假设微信号为:abc

2. 在image_hub/cmd/image_hub/main.go中修改需要被处理的html文档目录

```
 dir = "D:/work/wechat_download_data/html/Dump-0423-11-39-39"
```

3. 在image_hub库中创建文章表,表名为:tbl_article_微信号, 即: tbl_article_abc
```
CREATE TABLE `tbl_article_abc`  (
  `sn` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '一篇文章的唯一标识符',
  `mid` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '每次推送文章的唯一标识符',
  `idx` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '如果一次推送有多篇文章，idx表示当前页面是第几个',
  `biz` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '微信公众号的唯一标识符',
  `author` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '公众号作者名称',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '文章标题',
  `tags` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '合集标签',
  `sections` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '文章分段，一篇文章(article)由多个分段(section)组成',
  `local_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '文章本地保存路径',
  `publish_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '文章发布时间',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `deleted_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '删除时间',
  PRIMARY KEY (`sn`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '文章表' ROW_FORMAT = DYNAMIC;
```

4. 修改image_hub/model/article.go 中表名为tbl_article_abc
```
func (t *TblArticle) TableName() string {
	return "tbl_article_abc"
}
```
5. 在image_hub/spiders/func_map.go中按文章索引号增加自定义数据处理方法的map

```
  // 示例
	addFunc("abc_1", abc)
```

6. 在image_hub/spiders 下增加文字自定义处理方法,即4中map中的方法,文件名格式为: wechat_微信号.go, 即: image_hub/spiders/wechat_adc.go

7. 在image_hub/cmd/image_hub目录下，执行命令后,在 tbl_article_abc表中查看分析结果
```
 go run main.go -c ../../configs/config.yaml
```

#### 工具

1. 查看公众号内重复的图片
```
  // image_hub/cmd/image_hub/main.go 中
  // helper.GenerateHTML 会在D:/work/image_hub/test/目录下生成一个index-年月日时分秒.html的文件,用于观测重复图的内容
  helper.GenerateHTML("D:/work/image_hub/test/", dataSrcRepeat)
```

2. 查看公众号各个html的图片内容类型(这个工具目前不是特别的准确, 红框绿框可视化的部分没问题,但是检测内容类型的部分,还需优化)
```
  cd image_hub/cmd/image_border
  go run main.go
```

3. 在所有html中查找某个关键字在那个html中(类似grep)
```
  python search_html.py D:/work/wechat_download_data "丨人间值得"
```


#### 备忘
1. 增加了"动图","GIF动图"标签, 数据未导入 file:///D:/work/wechat_download_data/html/Dump-0425-09-02-16/20200319_120000_4.html


#### 目录处理记录
| | | | | |
| -------- | ------------------------------- | ------------------------------- | ------------------------------- | ------------------------------- |
| 序号 | 目录名 | 公众号 | 微信号 | 状态 |
| 1 | Dump-0421-11-15-39 | 头像社 | touxiangshe | done |
| 2 | Dump-0422-20-12-37 | 情侣头像原创榜 | seevanlove | done |
| 3 | Dump-0422-20-45-37 | 头像有点好看 | gh_8c96baecf453 | done |
| 4 | Dump-0422-20-54-12 | 头像即新欢 | gh_22c17e1db325 | done |
| 5 | Dump-0422-21-05-01 | 元気头像 | NiceWallpaper | done |
| 6 | Dump-0423-11-39-39 | 头像库 | touxiangcool | done |
| 7 | Dump-0423-12-06-34 | 头像文案 | fashionshijue | done |
| 8 | Dump-0423-19-16-40 | 你的小众头像 | h13031h | done |
| 9 | Dump-0423-19-23-35 | 换头像bo | htxb888 | done |
| 10 | Dump-0423-19-29-57 | 每日新头像 | gh_75640868571b | done |
| 11 | Dump-0423-19-36-06 | 头像备忘录 | DNTX9527 | done |
| 12 | Dump-0423-19-49-05 | 小鹿头像酱 | fairy_goods_thing | done |
| 13 | Dump-0423-19-56-42 | 梅头像 | MXLtou | done |
| 14 | Dump-0423-20-18-59 | 可爱cp头像 | tcgonglue | done |
| 15 | Dump-0423-20-51-56 | 要啥头像 | gh_cdb453299489 | done |
| 16 | Dump-0425-08-23-57 | 情侣头像大全 | qltxdq | done |
| 17 | Dump-0425-08-29-21 | 暮昭昭头像馆 | MzzTxg | done |
| 18 | Dump-0425-08-34-41 | 琉柒头像 | lik0894 | done  |
| 19 | Dump-0425-08-48-50 | 头像娣 | Txd777i | done  |
| 20 | Dump-0425-09-02-16 | 精选女生头像 | touxiang_520 | done  |
| 21 | Dump-0425-09-08-48 | 女生头像壁纸控 | touxiangdiss1 | done |
| 22 | Dump-0425-09-28-30 | 头像先生 | J79938 | done  |
| 23 | Dump-0425-10-08-29 | 头像味 | gh_bc125df08550 | done |
| 24 | Dump-0425-10-14-07 | 小怪兽头像 | gh_97a6f9e34972 | done |
| 25 | Dump-0425-10-20-23 | 二次元头像集 | cpdd52199 | done |
| 26 | Dump-0425-10-23-36 | 头像壁纸大全 | txbz001 | done |
| 27 | Dump-0425-10-30-17 | 搞怪沙雕头像 | youtiaotaolu | done |
| 28 | Dump-0425-10-35-23 | 头像记 | laixieee | done |
| 29 | Dump-0425-10-38-12 | 头像酱呀 | bizhi1994 |  |
| 30 | Dump-0425-10-42-59 | 头像文案 | fashionshijue |  |
| 31 | Dump-0428-22-35-17 | 头像辑 | touxiangh |  |
| 32 | Dump-0428-22-40-43 | 头像博主 | txbozhu |  |
| 33 | Dump-0428-22-46-16 | 头像号 | remenyt |  |
| 34 | Dump-0428-22-55-02 | 百合头像 | baihetouxiang |  |
| 35 | Dump-0428-23-03-29 | 背景头像 | meaijiepai |  |
| 36 | Dump-0429-18-36-56 | 头像录 | liaoshangbiji |  |
| 37 | Dump-0429-18-40-36 | 头像哒 | gh_367e376abfe0 |  |
| 38 | Dump-0429-18-44-47 | 女生头像宝藏集 | gh_a600aed1c30d |  |
| 39 | Dump-0503-21-16-12 | 超火情侣头像 | chanxuehuiyu |  |
| 40 | Dump-0504-10-02-15 | 头像社 | touxiangshe |  |
| 41 | Dump-0505-23-31-52 | 头像壁纸每日推荐 | touxbizhimeiriTJ |  |
| 42 | Dump-0507-10-40-20 | 搞怪头像大全 | gh_089775ff1457 |  |
| 43 | Dump-0507-12-39-56 | 头像微甜 | txwt-sweet |  |
| 44 | Dump-0508-00-15-59 | 女生头像壁纸 | nvshengtouxiang1 |  |
| 45 | Dump-0508-08-18-06 | ULzzang头像 | Ins-face |  |
| 46 | Dump-0508-08-21-43 | 头像书 | Txs5665 |  |
| 47 | Dump-0508-08-24-32 | 头像大全丫 | wulai969 |  |
| 48 | Dump-0509-09-13-18 | 可爱萌娃头像大全 | mwtx66695 |  |
| 49 | Dump-0509-09-18-51 | 萌娃头像库 | gh_bb13ee258433 |  |
| 50 | Dump-0509-09-20-43 | 可爱萌娃头像 | bqv8897 |  |
| 51 | Dump-0509-09-23-59 | 丸子妹头像 | bq6691 |  |
| 52 | Dump-0509-09-26-15 | 萌娃表情包可爱 | bqb598 |  |
| 53 | Dump-0510-10-43-36 | 搞怪头像大全 | gh_089775ff1457 |  |
| 54 | Dump-0510-10-47-14 | 搞怪沙雕头像 | youtiaotaolu |  |
| 55 | Dump-0510-10-49-15 | 搞怪头像合集 | gaoguaitx |  |
| 56 | Dump-0510-10-55-14 | 沙雕头像君 | bao_mihuaqi |  |
| 57 | Dump-0511-07-36-14 | 古风头像控 | gh_aca0610cf585 |  |
| 58 | Dump-0511-08-15-18 | 古风头像馆 | gftxg123 |  |
| 59 | Dump-0511-08-22-58 | 古风壁纸馆 | gfbzg-007 |  |
| 60 | Dump-0511-08-32-49 | 九栀头像 | afx1990 |  |
| 61 | Dump-0511-08-38-30 | 胖橘子呀 | Pang-Juziya |  |
| 62 | Dump-0511-08-43-21 | 玫竹斋 | meizhuzhai |  |
| 63 | Dump-0512-23-46-10 | 草莓头像 | touxiangforever |  |

 





 








内容包括：

1. 头像
2. 壁纸
3. 背景图
4. 表情包

5. 潮流扭蛋
6. 优惠券群
7. 恋爱话术

需求，根据颜色找图


抓取规则：
1. 头像社


微信公众号，通过聊天，找另一半，搜图，发支付二维码

ins 头像
复古风
氛围感
可爱头像
闺蜜头像
姐妹头像
遮脸男头
遮脸女头
红色系

根据图片的尺寸，判断图片是壁纸，头像，或是背景图

使用公众号合集的标签,做内容标签

目前支持的网站
| | | |
| -------- | ------------------------------- | ------------------------------- |
| 序号 | 公众号 | 网址 |
| 1 | 我要头像 | 2023-04-19 08:35:00 |

#### 线上环境

- 主机:
- 域名:
- 传输协议: http80 & https443
- 路由协议:
- 请求示例：
- 鉴权协议: 使用 appId+appSecret 鉴权

#### 测试环境

- 主机：
- 网关域名:
- 服务端口：80
- Nginx 配置:
- 部署路径：
- 程序端口: 8081
- 示例：
  - 设置 Host:
  - 调用方式一：
  - 调用方式二：

#### 测试环境部署步骤

- ip: 43.140.216.104
- 1、cd /home/gy/work/content_service
- 2、拉取最新代码
- 3、cd ./scripts 执行 sh build.sh content_service release 1.0.0
  - 1.0.0 为版本号，需改变
  - 将新编译的目录拷贝至 /home/gy/
  - cd ../out
  - cp -R content_service.release.1.0.0 /home/gy/
  - 修改 content_service 软连接指向 ../out/content_service.release.1.0.0 目录
  - ln -snf /home/gy/content_service.release.1.0.0 /home/gy/content_service
  - 重启服务 sudo supervisorctl restart content_service
  - 查询 supervisorctl 控制进程 supervisorctl status

#### 编译打包

- 测试包：./build/build.sh content_service debug 1.0.0
- 线上包：./build/build.sh content_service release 1.0.0

#### 注意事项



