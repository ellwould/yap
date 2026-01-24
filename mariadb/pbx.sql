CREATE TABLE `pbx`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `pbx_name` VARCHAR(75) NOT NULL,
  `group_id` BIGINT UNSIGNED NOT NULL,
  `date_added` DATETIME NOT NULL,
  `pbx_active` BOOLEAN NOT NULL,
  `note` VARCHAR(255),
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;
