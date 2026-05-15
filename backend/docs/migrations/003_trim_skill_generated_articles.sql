DROP PROCEDURE IF EXISTS `drop_skill_generated_article_index`;
DROP PROCEDURE IF EXISTS `drop_skill_generated_article_column`;

DELIMITER //

CREATE PROCEDURE `drop_skill_generated_article_index`(IN p_index_name VARCHAR(128))
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'skill_generated_articles'
      AND index_name = p_index_name
  ) THEN
    SET @stmt = CONCAT('ALTER TABLE `skill_generated_articles` DROP INDEX `', p_index_name, '`');
    PREPARE drop_index_stmt FROM @stmt;
    EXECUTE drop_index_stmt;
    DEALLOCATE PREPARE drop_index_stmt;
  END IF;
END//

CREATE PROCEDURE `drop_skill_generated_article_column`(IN p_column_name VARCHAR(128))
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'skill_generated_articles'
      AND column_name = p_column_name
  ) THEN
    SET @stmt = CONCAT('ALTER TABLE `skill_generated_articles` DROP COLUMN `', p_column_name, '`');
    PREPARE drop_column_stmt FROM @stmt;
    EXECUTE drop_column_stmt;
    DEALLOCATE PREPARE drop_column_stmt;
  END IF;
END//

DELIMITER ;

CALL `drop_skill_generated_article_index`('idx_skill_generated_articles_status');
CALL `drop_skill_generated_article_index`('idx_skill_generated_articles_publish_status');

CALL `drop_skill_generated_article_column`('keyword_segments');
CALL `drop_skill_generated_article_column`('html_preview');
CALL `drop_skill_generated_article_column`('tags');
CALL `drop_skill_generated_article_column`('cover_image_type');
CALL `drop_skill_generated_article_column`('images');
CALL `drop_skill_generated_article_column`('sources');
CALL `drop_skill_generated_article_column`('hot_topics');
CALL `drop_skill_generated_article_column`('style_profile');
CALL `drop_skill_generated_article_column`('word_count');
CALL `drop_skill_generated_article_column`('model_provider');
CALL `drop_skill_generated_article_column`('model_name');
CALL `drop_skill_generated_article_column`('prompt_version');
CALL `drop_skill_generated_article_column`('skill_version');
CALL `drop_skill_generated_article_column`('humanize_status');
CALL `drop_skill_generated_article_column`('status');
CALL `drop_skill_generated_article_column`('publish_status');
CALL `drop_skill_generated_article_column`('publish_payload');
CALL `drop_skill_generated_article_column`('error_message');

DROP PROCEDURE IF EXISTS `drop_skill_generated_article_index`;
DROP PROCEDURE IF EXISTS `drop_skill_generated_article_column`;
