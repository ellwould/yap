CREATE TABLE `user_account`
(
  `id` INT NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `first_name` VARCHAR(100) NOT NULL,
  `last_name` VARCHAR(100) NOT NULL,
  `user_account_type_id` INT NOT NULL,
  `group_id` INT NOT NULL,
  `pbx_id` INT NOT NULL,
  `date_added` DATETIME NOT NULL,
  `account_active` BOOLEAN NOT NULL,
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;
