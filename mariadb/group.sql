CREATE TABLE `group.sql`
(
  `id` INT NOT NULL,
  `group_name` VARCHAR(100) NOT NULL,
  `group_site_address_id` INT NOT NULL,
  `group_invoice_address_id` INT NOT NULL,
  `date_added` DATETIME NOT NULL,
  `group_active` BOOLEAN NOT NULL,
  `note` VARCHAR(255),
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;
