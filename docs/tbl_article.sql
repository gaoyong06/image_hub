CREATE TABLE `tbl_article` (
	`sn` VARCHAR(255) CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci' NOT NULL COMMENT '一篇文章的唯一标识符',
	`mid` BIGINT UNSIGNED NOT NULL DEFAULT 0  COMMENT '每次推送文章的唯一标识符',
	`idx` INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '如果一次推送有多篇文章，idx表示当前页面是第几个',
	`biz` VARCHAR(255) CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci' NOT NULL COMMENT '微信公众号的唯一标识符',
	`author` VARCHAR(255) CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci' NOT NULL DEFAULT '' COMMENT '公众号作者名称',
	`wechat_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '公众号微信号',
	`title` VARCHAR(255) CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci' NOT NULL DEFAULT '' COMMENT '文章标题',
	`tags` MEDIUMTEXT CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci' DEFAULT NULL COMMENT '合集标签',
	`sections` MEDIUMTEXT CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci' DEFAULT NULL COMMENT '文章分段，一篇文章(article)由多个分段(section)组成',
	`local_path` VARCHAR(255) CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci' NOT NULL DEFAULT '' COMMENT '文章本地保存路径',
	`publish_time` TIMESTAMP NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '文章发布时间',
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
	`deleted_at` TIMESTAMP NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '删除时间',
	PRIMARY KEY (`sn`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  COMMENT = '文章表';