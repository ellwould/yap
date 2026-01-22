CREATE TABLE `group_site_address.sql`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `address_line_1` VARCHAR(75) NOT NULL,
  `address_line_2` VARCHAR(75),
  `city_town_village` VARCHAR(75) NOT NULL,
  `postcode_zip_code` VARCHAR(7),
  `county_state_region` VARCHAR(75),
  `country` VARCHAR(75) NOT NULL,
  `contact_email` VARCHAR(255),
  `contact_number` VARCHAR(16),
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;
