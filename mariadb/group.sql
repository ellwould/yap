CREATE TABLE `group.sql`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `group_name` VARCHAR(100) NOT NULL,
  `date_added` DATETIME NOT NULL,
  `group_active` BOOLEAN NOT NULL,
  `note` VARCHAR(255),
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;
