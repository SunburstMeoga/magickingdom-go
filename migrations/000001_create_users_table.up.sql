-- 创建用户表
CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `open_id` varchar(100) NOT NULL COMMENT '微信 OpenID',
  `union_id` varchar(100) DEFAULT NULL COMMENT '微信 UnionID',
  `session_key` varchar(100) DEFAULT NULL COMMENT '微信 SessionKey',
  `nickname` varchar(100) DEFAULT NULL COMMENT '昵称',
  `avatar_url` varchar(500) DEFAULT NULL COMMENT '头像 URL',
  `gender` int DEFAULT 0 COMMENT '性别：0-未知 1-男 2-女',
  `country` varchar(50) DEFAULT NULL COMMENT '国家',
  `province` varchar(50) DEFAULT NULL COMMENT '省份',
  `city` varchar(50) DEFAULT NULL COMMENT '城市',
  `language` varchar(20) DEFAULT NULL COMMENT '语言',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `status` int DEFAULT 1 COMMENT '状态：1-正常 0-禁用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_open_id` (`open_id`),
  KEY `idx_users_deleted_at` (`deleted_at`),
  KEY `idx_users_union_id` (`union_id`),
  KEY `idx_users_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

