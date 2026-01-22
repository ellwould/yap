CREATE TABLE `pbx.sql`
(
  `id` INT NOT NULL,
  `pbx_name` VARCHAR(75) NOT NULL,
  `pbx_site_address_id` INT NOT NULL,
  `pbx_invoice_address_id` INT NOT NULL,
  `group_id` INT NOT NULL,
  `date_added` DATETIME NOT NULL,
  `pbx_active` BOOLEAN NOT NULL,
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;
