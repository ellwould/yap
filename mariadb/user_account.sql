CREATE TABLE `user_account`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `first_name` VARCHAR(100) NOT NULL,
  `last_name` VARCHAR(100) NOT NULL,
  `user_account_type_id` SMALLINT UNSIGNED NOT NULL,
  `group_id` BIGINT UNSIGNED NOT NULL,
  `pbx_id` BIGINT UNSIGNED NOT NULL,
  `date_added` DATETIME NOT NULL,
  `account_active` BOOLEAN NOT NULL,
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

ALTER TABLE `user_account`
ADD INDEX `index___user_account__group_id` (`group_id`);

ALTER TABLE `user_account`
ADD INDEX `index___user_account__pbx_id` (`pbx_id`);
