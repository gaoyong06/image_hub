### 内容服务（Image Hub）

|          |                                 |
| -------- | ------------------------------- |
| 作者     | [高勇](mailto:gaoyong06@qq.com) |
| 创建时间 | 2023-04-19 08:35:00             |

#### 介绍

imageHub,是一个图片内容源服务，主要服务 content_service 的内容源采集 现在的目标是采集,主流的微信公众号,主流的图片网站内容

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
