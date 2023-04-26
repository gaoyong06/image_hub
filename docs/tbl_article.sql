CREATE TABLE `tbl_article` (
	`mid` BIGINT UNSIGNED NOT NULL COMMENT '文章id 每篇文章的唯一标识符',
	`biz` VARCHAR(255) CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci' NOT NULL COMMENT '微信公众号的唯一标识符',
	`idx` INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '如果一篇文章有多页内容，idx表示当前页面是第几页',
	`sn` VARCHAR(255) CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci' NOT NULL DEFAULT '' COMMENT '一篇文章的唯一标识符，与mid不同的是，sn是加密后的标识符',
	`title` VARCHAR(255) CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci' NOT NULL DEFAULT '' COMMENT '标题',
	`author` VARCHAR(255) CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci' NOT NULL DEFAULT '' COMMENT '作者',
	`tags` MEDIUMTEXT CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci' DEFAULT NULL COMMENT '合集标签',
	`sections` MEDIUMTEXT CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci' DEFAULT NULL COMMENT '文章分段，一篇文章(article)由多个分段(section)组成',
	`local_path` VARCHAR(255) CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci' NOT NULL DEFAULT '' COMMENT '文章保存路径',
	`publish_time` TIMESTAMP NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '发布时间',
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
	`deleted_at` TIMESTAMP NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '删除时间',
	PRIMARY KEY (`mid`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  COMMENT = '文章表';
