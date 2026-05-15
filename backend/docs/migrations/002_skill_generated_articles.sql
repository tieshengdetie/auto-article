CREATE TABLE IF NOT EXISTS `skill_generated_articles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `task_id` varchar(64) NOT NULL DEFAULT '',
  `platform` varchar(32) NOT NULL DEFAULT '',
  `keyword` varchar(128) NOT NULL DEFAULT '',
  `category` varchar(64) NOT NULL DEFAULT '',
  `title` varchar(512) NOT NULL DEFAULT '',
  `title_options` text,
  `summary` text,
  `markdown_content` longtext,
  `cover_image_url` varchar(1024) NOT NULL DEFAULT '',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_skill_generated_articles_task_id` (`task_id`),
  KEY `idx_skill_generated_articles_platform` (`platform`),
  KEY `idx_skill_generated_articles_keyword` (`keyword`),
  KEY `idx_skill_generated_articles_category` (`category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
