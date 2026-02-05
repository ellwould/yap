-- Create YAP tables

CREATE TABLE `group`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `group_name` VARCHAR(100) NOT NULL,
  `date_added` DATETIME DEFAULT NOW() NOT NULL,
  `group_active` BOOLEAN NOT NULL,
  `note` VARCHAR(255),
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `group_invoice_address`
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

CREATE TABLE `group_site_address`
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

CREATE TABLE `pbx`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `pbx_name` VARCHAR(75) NOT NULL,
  `group_id` BIGINT UNSIGNED NOT NULL,
  `date_added` DATETIME DEFAULT NOW() NOT NULL,
  `pbx_active` BOOLEAN NOT NULL,
  `note` VARCHAR(255),
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `pbx_invoice_address`
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

CREATE TABLE `pbx_site_address`
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

CREATE TABLE `user_account`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `first_name` VARCHAR(100) NOT NULL,
  `last_name` VARCHAR(100) NOT NULL,
  `user_account_type_id` SMALLINT UNSIGNED NOT NULL,
  `group_id` BIGINT UNSIGNED NOT NULL,
  `pbx_id` BIGINT UNSIGNED NOT NULL,
  `date_added` DATETIME DEFAULT NOW() NOT NULL,
  `account_active` BOOLEAN NOT NULL,
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `user_account_type`
(
  `id` SMALLINT UNSIGNED NOT NULL,
  `type` VARCHAR(20) NOT NULL,
  `description` VARCHAR(100) NOT NULL,
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

----------------------------------------------------------------------------------------------------

-- Add pbx_id column to Asterisk tables

ALTER TABLE `ps_endpoints`
ADD COLUMN `pbx_id` BIGINT UNSIGNED NOT NULL;

ALTER TABLE `ps_aors`
ADD COLUMN `pbx_id` BIGINT UNSIGNED NOT NULL;

ALTER TABLE `ps_auths`
ADD COLUMN `pbx_id` BIGINT UNSIGNED NOT NULL;

----------------------------------------------------------------------------------------------------

-- Alter existing data types

ALTER TABLE `ps_endpoints`
MODIFY COLUMN `aors` varchar(255) NOT NULL;

ALTER TABLE `ps_endpoints`
MODIFY COLUMN `auth` varchar(255) NOT NULL;

ALTER TABLE `ps_aors`
MODIFY COLUMN `id` varchar(255) NOT NULL;

----------------------------------------------------------------------------------------------------

-- Add index to columns

ALTER TABLE `pbx`
ADD INDEX `index___pbx__group_id` (`group_id`);

ALTER TABLE `user_account`
ADD INDEX `index___user_account__user_account_type_id` (`user_account_type_id`);

ALTER TABLE `user_account`
ADD INDEX `index___user_account__group_id` (`group_id`);

ALTER TABLE `user_account`
ADD INDEX `index___user_account__pbx_id` (`pbx_id`);

ALTER TABLE `ps_endpoints`
ADD INDEX `index___ps_endpoints__aors` (`aors`);

ALTER TABLE `ps_endpoints`
ADD INDEX `index___ps_endpoints__auth` (`auth`);

ALTER TABLE `ps_endpoints`
ADD INDEX `index___ps_endpoints__pbx_id` (`pbx_id`);

ALTER TABLE `ps_aors`
ADD INDEX `index___ps_aors__pbx_id` (`pbx_id`);

ALTER TABLE `ps_auths`
ADD INDEX `index___ps_auths__pbx_id` (`pbx_id`);

----------------------------------------------------------------------------------------------------

-- Create foreign key constraints

ALTER TABLE `group_invoice_address`
ADD CONSTRAINT fk___group_invoice_address___group
FOREIGN KEY (`id`)
REFERENCES `group` (`id`)
ON DELETE CASCADE;

ALTER TABLE `group_site_address`
ADD CONSTRAINT fk___group_site_address___group
FOREIGN KEY (`id`)
REFERENCES `group` (`id`)
ON DELETE CASCADE;

ALTER TABLE `pbx_invoice_address`
ADD CONSTRAINT fk___pbx_invoice_address___pbx
FOREIGN KEY (`id`)
REFERENCES `pbx` (`id`)
ON DELETE CASCADE;

ALTER TABLE `pbx_site_address`
ADD CONSTRAINT fk___pbx_site_address___pbx
FOREIGN KEY (`id`)
REFERENCES `pbx` (`id`)
ON DELETE CASCADE;

ALTER TABLE `user_account`
ADD CONSTRAINT fk___user_account___user_account_type
FOREIGN KEY (`user_account_type_id`)
REFERENCES `user_account_type` (`id`);

ALTER TABLE `user_account`
ADD CONSTRAINT fk___user_account___group
FOREIGN KEY (`group_id`)
REFERENCES `group` (`id`);

ALTER TABLE `user_account`
ADD CONSTRAINT fk___user_account___pbx
FOREIGN KEY (`pbx_id`)
REFERENCES `pbx` (`id`);

ALTER TABLE `ps_endpoints`
ADD CONSTRAINT fk___ps_endpoints___pbx
FOREIGN KEY (`pbx_id`)
REFERENCES `pbx` (`id`)
ON DELETE CASCADE;

ALTER TABLE `ps_aors`
ADD CONSTRAINT fk___ps_aors___ps_endpoints
FOREIGN KEY (`id`)
REFERENCES `ps_endpoints` (`id`)
ON DELETE CASCADE;

ALTER TABLE `ps_aors`
ADD CONSTRAINT fk___ps_aors___pbx
FOREIGN KEY (`pbx_id`)
REFERENCES `pbx` (`id`)
ON DELETE CASCADE;

ALTER TABLE `ps_auths`
ADD CONSTRAINT fk___ps_auths___ps_endpoints
FOREIGN KEY (`id`)
REFERENCES `ps_endpoints` (`id`)
ON DELETE CASCADE;

ALTER TABLE `ps_auths`
ADD CONSTRAINT fk___ps_auths___pbx
FOREIGN KEY (`pbx_id`)
REFERENCES `pbx` (`id`)
ON DELETE CASCADE;

----------------------------------------------------------------------------------------------------

-- Create View(s)

CREATE VIEW `view___account_detail` AS
SELECT
  `user_account`.`user_account_type_id`,
  `user_account`.`first_name` AS 'user_account_first_name',
  `user_account`.`last_name` AS 'user_account_last_name',
  `user_account`.`email` AS 'user_account_email',
  `user_account_type`.`type` AS 'user_account_type',
  `user_account`.`date_added` AS 'user_account_date_added',
  `group`.`group_name`,
  `user_account`.`group_id`,
  `pbx`.`pbx_name`,
  `user_account`.`pbx_id`,
  `user_account_type`.`description` AS 'user_account_type_description',
  `group_site_address`.`address_line_1` AS 'group_site_address_line_1',
  `group_site_address`.`address_line_2` AS 'group_site_address_line_2',
  `group_site_address`.`city_town_village` AS 'group_site_city_town_village',
  `group_site_address`.`postcode_zip_code` AS 'group_site_postcode_zip_code',
  `group_site_address`.`county_state_region` AS 'group_site_county_state_region',
  `group_site_address`.`country` AS 'group_site_country',
  `group_site_address`.`contact_email` AS 'group_site_contact_email',
  `group_site_address`.`contact_number` AS 'group_site_contact_number',
  `group_invoice_address`.`address_line_1` AS 'group_invoice_address_line_1',
  `group_invoice_address`.`address_line_2` AS 'group_invoice_address_line_2',
  `group_invoice_address`.`city_town_village` AS 'group_invoice_city_town_village',
  `group_invoice_address`.`postcode_zip_code` AS 'group_invoice_postcode_zip_code',
  `group_invoice_address`.`county_state_region` AS 'group_invoice_county_state_region',
  `group_invoice_address`.`country` AS 'group_invoice_country',
  `group_invoice_address`.`contact_email` AS 'group_invoice_contact_email',
  `group_invoice_address`.`contact_number` AS 'group_invoice_contact_number',
  `pbx_site_address`.`address_line_1` AS 'pbx_site_address_line_1',
  `pbx_site_address`.`address_line_2` AS 'pbx_site_address_line_2',
  `pbx_site_address`.`city_town_village` AS 'pbx_site_city_town_village',
  `pbx_site_address`.`postcode_zip_code` AS 'pbx_site_postcode_zip_code',
  `pbx_site_address`.`county_state_region` AS 'pbx_site_county_state_region',
  `pbx_site_address`.`country` AS 'pbx_site_country',
  `pbx_site_address`.`contact_email` AS 'pbx_site_contact_email',
  `pbx_site_address`.`contact_number` AS 'pbx_site_contact_number',
  `pbx_invoice_address`.`address_line_1` AS 'pbx_invoice_address_line_1',
  `pbx_invoice_address`.`address_line_2` AS 'pbx_invoice_address_line_2',
  `pbx_invoice_address`.`city_town_village` AS 'pbx_invoice_city_town_village',
  `pbx_invoice_address`.`postcode_zip_code` AS 'pbx_invoice_postcode_zip_code',
  `pbx_invoice_address`.`county_state_region` AS 'pbx_invoice_county_state_region',
  `pbx_invoice_address`.`country` AS 'pbx_invoice_country',
  `pbx_invoice_address`.`contact_email` AS 'pbx_invoice_contact_email',
  `pbx_invoice_address`.`contact_number` AS 'pbx_invoice_contact_number'
FROM `user_account`
INNER JOIN `user_account_type`
ON `user_account`.`user_account_type_id` = `user_account_type`.`id`
INNER JOIN `group`
ON `user_account`.`group_id` = `group`.`id`
INNER JOIN `pbx`
ON `user_account`.`pbx_id` = `pbx`.`id`
INNER JOIN `group_site_address`
ON `user_account`.`group_id` = `group_site_address`.`id`
INNER JOIN `group_invoice_address`
ON `user_account`.`group_id` = `group_invoice_address`.`id`
INNER JOIN `pbx_site_address`
ON `user_account`.`pbx_id` = `pbx_site_address`.`id`
INNER JOIN `pbx_invoice_address`
ON `user_account`.`pbx_id` = `pbx_invoice_address`.`id`;

----------------------------------------------------------------------------------------------------

-- Insert data to YAP tables

INSERT INTO `group` (`id`, `group_name`, `group_active`, `note`)
VALUES (1, 'system', 1, 'created during YAP install');

INSERT INTO `group_invoice_address` (`id`,	`address_line_1`,	`address_line_2`,	`city_town_village`, `postcode_zip_code`,	`county_state_region`, `country`,	`contact_email`, `contact_number`)
VALUES (1, 'system', 'system', 'system', 'system', 'system', 'system', 'system', 'system');

INSERT INTO `group_site_address` (`id`,	`address_line_1`,	`address_line_2`,	`city_town_village`, `postcode_zip_code`,	`county_state_region`, `country`,	`contact_email`, `contact_number`)
VALUES (1, 'system', 'system', 'system', 'system', 'system', 'system', 'system', 'system');

INSERT INTO `pbx` (`id`, `pbx_name`, `group_id`, `pbx_active`, `note`)
VALUES (1, 'system', 1, 1, 'created during YAP install');

INSERT INTO `pbx_invoice_address` (`id`,	`address_line_1`,	`address_line_2`,	`city_town_village`, `postcode_zip_code`,	`county_state_region`, `country`,	`contact_email`, `contact_number`)
VALUES (1, 'system', 'system', 'system', 'system', 'system', 'system', 'system', 'system');

INSERT INTO `pbx_site_address` (`id`,	`address_line_1`,	`address_line_2`,	`city_town_village`, `postcode_zip_code`,	`county_state_region`, `country`,	`contact_email`, `contact_number`)
VALUES (1, 'system', 'system', 'system', 'system', 'system', 'system', 'system', 'system');

INSERT INTO `user_account_type` (`id`, `type`, `description`)
VALUES
(100, 'YAP Admin (100)', 'A YAP admin account can create, read, update and delete all user accounts, groups and PBXs'),
(101, 'YAP Read Only (101)', 'A YAP regular account can read all user accounts, groups and PBXs'),
(200, 'Group Admin (200)', 'A group admin can read and update thier own PBX(s) and group'),
(201, 'Group Read Only (201)', 'A group regular account can read thier own PBX(s) and group'),
(300, 'PBX Admin (300)', 'A PBX admin account can read and update thier own PBX'),
(301, 'PBX Read Only (301)', 'A PBX regular account can read thier own PBX');
