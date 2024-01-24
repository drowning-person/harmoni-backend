-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `like` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` bigint UNSIGNED NOT NULL,
  `liking_id` bigint NOT NULL COMMENT '点赞ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `target_user_id` bigint NOT NULL COMMENT '点赞对象所属用户ID',
  `object_id` bigint NOT NULL COMMENT '对象ID',
  `object_type` tinyint NOT NULL COMMENT '对象类型',
  PRIMARY KEY (`id`)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `like`;
-- +goose StatementEnd
