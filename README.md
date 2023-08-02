### 内容服务（Image Hub）

|          |                                 |
| -------- | ------------------------------- |
| 作者     | [高勇](mailto:gaoyong06@qq.com) |
| 创建时间 | 2023-04-19 08:35:00             |

#### 介绍

imageHub,是一个图片内容源服务，主要服务 content_service 的内容源采集 现在的目标是采集,主流的微信公众号,主流的图片网站内容


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


CREATE TABLE `tbl_article_gh_8c96baecf453`  (
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
