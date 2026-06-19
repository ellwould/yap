-- Create YAP tables

CREATE TABLE `customer`
(
  `id` VARCHAR(255) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `date_added` DATETIME DEFAULT NOW() NOT NULL,
  `active` BOOLEAN NOT NULL,
  `uk_based` ENUM('yes', 'no', 'n/a') NOT NULL,
  `consumer_type` VARCHAR(255) NOT NULL,
  `uk_vat_status` VARCHAR(255) NOT NULL,
  `reselling_miniutes` ENUM('no', 'yes', 'n/a') NOT NULL,
  `pbx_limit` SMALLINT UNSIGNED DEFAULT 20 NOT NULL,
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `consumer_type_lookup` (
  `consumer_type` VARCHAR(255),
  PRIMARY KEY (`consumer_type`)
)
ENGINE = InnoDB;
  
CREATE TABLE `uk_vat_status_lookup` (
  `uk_vat_status` VARCHAR(255),
  PRIMARY KEY (`uk_vat_status`)
)
ENGINE = InnoDB;

CREATE TABLE `customer_invoice_address`
(
  `id` VARCHAR(255) NOT NULL,
  `address_line_1` VARCHAR(255) NOT NULL,
  `address_line_2` VARCHAR(255) NOT NULL,
  `city_town_village` VARCHAR(255) NOT NULL,
  `county_state_region` VARCHAR(255) NOT NULL,
  `postcode_zip_code` VARCHAR(7) NOT NULL,
  `country` VARCHAR(255) NOT NULL,
  `contact_email` VARCHAR(255) NOT NULL,
  `contact_number` VARCHAR(20) NOT NULL,
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `customer_site_address`
(
  `id` VARCHAR(255) NOT NULL,
  `address_line_1` VARCHAR(255) NOT NULL,
  `address_line_2` VARCHAR(255) NOT NULL,
  `city_town_village` VARCHAR(255) NOT NULL,
  `county_state_region` VARCHAR(255) NOT NULL,
  `postcode_zip_code` VARCHAR(7) NOT NULL,
  `country` VARCHAR(255) NOT NULL,
  `contact_email` VARCHAR(255) NOT NULL,
  `contact_number` VARCHAR(20) NOT NULL,
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `pbx`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `customer_id` VARCHAR(255) NOT NULL,
  `date_added` DATETIME DEFAULT NOW() NOT NULL,
  `active` BOOLEAN NOT NULL,
  `sip_extension_limit` SMALLINT UNSIGNED DEFAULT 100 NOT NULL,
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `pbx_invoice_address`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `address_line_1` VARCHAR(255) NOT NULL,
  `address_line_2` VARCHAR(255) NOT NULL,
  `city_town_village` VARCHAR(255) NOT NULL,
  `county_state_region` VARCHAR(255) NOT NULL,
  `postcode_zip_code` VARCHAR(7) NOT NULL,
  `country` VARCHAR(255) NOT NULL,
  `contact_email` VARCHAR(255) NOT NULL,
  `contact_number` VARCHAR(20) NOT NULL,
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `pbx_site_address`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `address_line_1` VARCHAR(255) NOT NULL,
  `address_line_2` VARCHAR(255) NOT NULL,
  `city_town_village` VARCHAR(255) NOT NULL,
  `county_state_region` VARCHAR(255) NOT NULL,
  `postcode_zip_code` VARCHAR(7) NOT NULL,
  `country` VARCHAR(255) NOT NULL,
  `contact_email` VARCHAR(255) NOT NULL,
  `contact_number` VARCHAR(20) NOT NULL,
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `user_account`
(
  `id` BIGINT UNSIGNED AUTO_INCREMENT NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `first_name` VARCHAR(255) NOT NULL,
  `last_name` VARCHAR(255) NOT NULL,
  `user_account_type_id` SMALLINT UNSIGNED NOT NULL,
  `customer_id` VARCHAR(255) NOT NULL,
  `pbx_id` BIGINT UNSIGNED NOT NULL,
  `date_added` DATETIME DEFAULT NOW() NOT NULL,
  `account_active` BOOLEAN NOT NULL,
UNIQUE (`email`),
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `user_account_type`
(
  `id` SMALLINT UNSIGNED NOT NULL,
  `type` VARCHAR(255) NOT NULL,
  `permission` VARCHAR(4000) NOT NULL,
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `invoice` (
  `id` INT UNSIGNED,
  `customer_id` VARCHAR(255) NOT NULL,
  `phone_number` VARCHAR(20) NOT NULL,
  `good_service_id` VARCHAR(255) NOT NULL,
  `sell_price` DECIMAL(8,2) NOT NULL,
  `uk_sales_tax_rate` DECIMAL(5,2) NOT NULL,
  `uk_sales_tax_status` VARCHAR(255) NOT NULL,
  `invoice_customer` ENUM('yes', 'no') NOT NULL,
  `one_off_charge` ENUM('yes', 'no') NOT NULL,
  `date_added` DATETIME DEFAULT NOW() NOT NULL,
  `comment` VARCHAR(255) NOT NULL,
  PRIMARY KEY(`id`)
)
ENGINE = InnoDB;
  
CREATE TABLE `uk_sales_tax_rate_lookup` (
  `uk_sales_tax_rate` DECIMAL(5,2),
  PRIMARY KEY(`uk_sales_tax_rate`)
)
ENGINE = InnoDB;
  
CREATE TABLE `uk_sales_tax_status_lookup` (
  `uk_sales_tax_status` VARCHAR(255),
  PRIMARY KEY(`uk_sales_tax_status`)
)
ENGINE = InnoDB;

CREATE TABLE `good_service` (
  `id` VARCHAR(255) NOT NULL,
  `good_service_type` VARCHAR(255) NOT NULL,
  `supplier_name` VARCHAR(255) NOT NULL,
  `buy_price` DECIMAL(8,2) NOT NULL,
  `date_added` DATETIME NOT NULL,
  PRIMARY KEY(`good_service_name`)
)
ENGINE = InnoDB;

CREATE TABLE `good_service_type_lookup` (
  `good_service_type` VARCHAR(255),
  PRIMARY KEY(`good_service_type`)
)
ENGINE = InnoDB;

----------------------------------------------------------------------------------------------------

-- Add pbx_id and and endpoint_type column to Asterisk tables

ALTER TABLE `ps_endpoints`
ADD COLUMN `endpoint_type` ENUM ('sip_extension', 'webrtc_extension') NOT NULL;

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
ALTER TABLE `customer`
ADD INDEX `index___customer__consumer_type` (`consumer_type`);

ALTER TABLE `customer`
ADD INDEX `index___customer__uk_vat_status` (`uk_vat_status`);

ALTER TABLE `pbx`
ADD INDEX `index___pbx__customer_id` (`customer_id`);

ALTER TABLE `user_account`
ADD INDEX `index___user_account__user_account_type_id` (`user_account_type_id`);

ALTER TABLE `user_account`
ADD INDEX `index___user_account__customer_id` (`customer_id`);

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

ALTER TABLE `invoice`
ADD INDEX `index___invoice__customer_id` (`customer_id`);

ALTER TABLE `invoice`
ADD INDEX `index___invoice__uk_sales_tax_rate` (`uk_sales_tax_rate`);

ALTER TABLE `invoice`
ADD INDEX `index___invoice__uk_sales_tax_status` (`uk_sales_tax_status`);

ALTER TABLE `invoice`
ADD INDEX `index___invoice__good_service_id` (`good_service_id`);

ALTER TABLE `invoice`
ADD INDEX `index___invoice__good_service_type` (`good_service_type`);

----------------------------------------------------------------------------------------------------

-- Create foreign key constraints

ALTER TABLE `customer`
ADD CONSTRAINT fk___customer___consumer_type_lookup
FOREIGN KEY (`consumer_type`)
REFERENCES `consumer_type_lookup` (`consumer_type`);

ALTER TABLE `customer`
ADD CONSTRAINT fk___customer___uk_vat_status_lookup
FOREIGN KEY (`uk_vat_status`)
REFERENCES `uk_vat_status_lookup` (`uk_vat_status`);

ALTER TABLE `pbx`
ADD CONSTRAINT fk___pbx___customer
FOREIGN KEY (`customer_id`)
REFERENCES `customer` (`id`)
ON DELETE CASCADE;

ALTER TABLE `customer_invoice_address`
ADD CONSTRAINT fk___customer_invoice_address___customer
FOREIGN KEY (`id`)
REFERENCES `customer` (`id`)
ON DELETE CASCADE;

ALTER TABLE `customer_site_address`
ADD CONSTRAINT fk___customer_site_address___customer
FOREIGN KEY (`id`)
REFERENCES `customer` (`id`)
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
ADD CONSTRAINT fk___user_account___customer
FOREIGN KEY (`customer_id`)
REFERENCES `customer` (`id`);

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

ALTER TABLE `invoice`
ADD CONSTRAINT fk___invoice___customer_id
FOREIGN KEY (`customer_id`)
REFERENCES `customer` (`id`)
ON DELETE CASCADE;

ALTER TABLE `invoice`
ADD CONSTRAINT fk___invoice___uk_sales_tax_rate_lookup
FOREIGN KEY (`uk_sales_tax_rate`)
REFERENCES `uk_sales_tax_rate_lookup` (`uk_sales_tax_rate`);

ALTER TABLE `invoice`
ADD CONSTRAINT fk___invoice___uk_sales_tax_status_lookup
FOREIGN KEY (`uk_sales_tax_status`)
REFERENCES `uk_sales_tax_status_lookup` (`uk_sales_tax_status`);

ALTER TABLE `invoice`
ADD CONSTRAINT fk___invoice___good_service
FOREIGN KEY (`good_service_id`)
REFERENCES `good_service` (`id`);

ALTER TABLE `good_service`
ADD CONSTRAINT fk___good_service___good_service_type_lookup
FOREIGN KEY (`good_service_type`)
REFERENCES `good_service_type_lookup` (`good_service_type`);

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
  `customer`.`name` AS 'customer_name',
  `user_account`.`customer_id`,
  `pbx`.`name` AS 'pbx_name',
  `user_account`.`pbx_id`,
  `user_account_type`.`permission` AS 'user_account_type_permission',
  `customer_site_address`.`address_line_1` AS 'customer_site_address_line_1',
  `customer_site_address`.`address_line_2` AS 'customer_site_address_line_2',
  `customer_site_address`.`city_town_village` AS 'customer_site_city_town_village',
  `customer_site_address`.`county_state_region` AS 'customer_site_county_state_region',
  `customer_site_address`.`postcode_zip_code` AS 'customer_site_postcode_zip_code',
  `customer_site_address`.`country` AS 'customer_site_country',
  `customer_site_address`.`contact_email` AS 'customer_site_contact_email',
  `customer_site_address`.`contact_number` AS 'customer_site_contact_number',
  `customer_invoice_address`.`address_line_1` AS 'customer_invoice_address_line_1',
  `customer_invoice_address`.`address_line_2` AS 'customer_invoice_address_line_2',
  `customer_invoice_address`.`city_town_village` AS 'customer_invoice_city_town_village',
  `customer_invoice_address`.`county_state_region` AS 'customer_invoice_county_state_region',
  `customer_invoice_address`.`postcode_zip_code` AS 'customer_invoice_postcode_zip_code',
  `customer_invoice_address`.`country` AS 'customer_invoice_country',
  `customer_invoice_address`.`contact_email` AS 'customer_invoice_contact_email',
  `customer_invoice_address`.`contact_number` AS 'customer_invoice_contact_number',
  `pbx_site_address`.`address_line_1` AS 'pbx_site_address_line_1',
  `pbx_site_address`.`address_line_2` AS 'pbx_site_address_line_2',
  `pbx_site_address`.`city_town_village` AS 'pbx_site_city_town_village',
  `pbx_site_address`.`county_state_region` AS 'pbx_site_county_state_region',
  `pbx_site_address`.`postcode_zip_code` AS 'pbx_site_postcode_zip_code',
  `pbx_site_address`.`country` AS 'pbx_site_country',
  `pbx_site_address`.`contact_email` AS 'pbx_site_contact_email',
  `pbx_site_address`.`contact_number` AS 'pbx_site_contact_number',
  `pbx_invoice_address`.`address_line_1` AS 'pbx_invoice_address_line_1',
  `pbx_invoice_address`.`address_line_2` AS 'pbx_invoice_address_line_2',
  `pbx_invoice_address`.`city_town_village` AS 'pbx_invoice_city_town_village',
  `pbx_invoice_address`.`county_state_region` AS 'pbx_invoice_county_state_region',
  `pbx_invoice_address`.`postcode_zip_code` AS 'pbx_invoice_postcode_zip_code',
  `pbx_invoice_address`.`country` AS 'pbx_invoice_country',
  `pbx_invoice_address`.`contact_email` AS 'pbx_invoice_contact_email',
  `pbx_invoice_address`.`contact_number` AS 'pbx_invoice_contact_number'
FROM `user_account`
INNER JOIN `user_account_type`
ON `user_account`.`user_account_type_id` = `user_account_type`.`id`
INNER JOIN `customer`
ON `user_account`.`customer_id` = `customer`.`id`
INNER JOIN `pbx`
ON `user_account`.`pbx_id` = `pbx`.`id`
INNER JOIN `customer_site_address`
ON `user_account`.`customer_id` = `customer_site_address`.`id`
INNER JOIN `customer_invoice_address`
ON `user_account`.`customer_id` = `customer_invoice_address`.`id`
INNER JOIN `pbx_site_address`
ON `user_account`.`pbx_id` = `pbx_site_address`.`id`
INNER JOIN `pbx_invoice_address`
ON `user_account`.`pbx_id` = `pbx_invoice_address`.`id`;

CREATE VIEW `view___customer_detail` AS
SELECT
  `customer`.`id` AS 'customer_id',
  `customer`.`name` AS 'customer_name',
  `customer`.`date_added` AS 'customer_date_added',
  `customer`.`active` AS 'customer_active',
  `customer`.`uk_based` AS 'customer_uk_based',
  `customer`.`consumer_type` AS 'customer_consumer_type',
  `customer`.`uk_vat_status` AS 'customer_uk_vat_status',
  `customer`.`reselling_miniutes` AS 'customer_reselling_miniutes',
  `customer`.`pbx_limit` AS 'customer_pbx_limit',
  `customer_site_address`.`address_line_1` AS 'customer_site_address_line_1',
  `customer_site_address`.`address_line_2` AS 'customer_site_address_line_2',
  `customer_site_address`.`city_town_village` AS 'customer_site_city_town_village',
  `customer_site_address`.`county_state_region` AS 'customer_site_county_state_region',
  `customer_site_address`.`postcode_zip_code` AS 'customer_site_postcode_zip_code',
  `customer_site_address`.`country` AS 'customer_site_country',
  `customer_site_address`.`contact_email` AS 'customer_site_contact_email',
  `customer_site_address`.`contact_number` AS 'customer_site_contact_number',
  `customer_invoice_address`.`address_line_1` AS 'customer_invoice_address_line_1',
  `customer_invoice_address`.`address_line_2` AS 'customer_invoice_address_line_2',
  `customer_invoice_address`.`city_town_village` AS 'customer_invoice_city_town_village',
  `customer_invoice_address`.`county_state_region` AS 'customer_invoice_county_state_region',
  `customer_invoice_address`.`postcode_zip_code` AS 'customer_invoice_postcode_zip_code',
  `customer_invoice_address`.`country` AS 'customer_invoice_country',
  `customer_invoice_address`.`contact_email` AS 'customer_invoice_contact_email',
  `customer_invoice_address`.`contact_number` AS 'customer_invoice_contact_number'
FROM `customer`
INNER JOIN `customer_site_address`
ON `customer`.`id` = `customer_site_address`.`id`
INNER JOIN `customer_invoice_address`
ON `customer`.`id` = `customer_invoice_address`.`id`;

CREATE VIEW `view___pbx_detail` AS
SELECT
  `pbx`.`id` AS 'pbx_id',
  `pbx`.`name` AS 'pbx_name',
  `pbx`.`customer_id`,
  `customer`.`name` AS 'customer_name',
  `pbx`.`date_added` AS 'pbx_date_added',
  `pbx`.`active` AS 'pbx_active',
  `pbx`.`sip_extension_limit` AS 'pbx_sip_extension_limit',
  `pbx_site_address`.`address_line_1` AS 'pbx_site_address_line_1',
  `pbx_site_address`.`address_line_2` AS 'pbx_site_address_line_2',
  `pbx_site_address`.`city_town_village` AS 'pbx_site_city_town_village',
  `pbx_site_address`.`county_state_region` AS 'pbx_site_county_state_region',
  `pbx_site_address`.`postcode_zip_code` AS 'pbx_site_postcode_zip_code',
  `pbx_site_address`.`country` AS 'pbx_site_country',
  `pbx_site_address`.`contact_email` AS 'pbx_site_contact_email',
  `pbx_site_address`.`contact_number` AS 'pbx_site_contact_number',
  `pbx_invoice_address`.`address_line_1` AS 'pbx_invoice_address_line_1',
  `pbx_invoice_address`.`address_line_2` AS 'pbx_invoice_address_line_2',
  `pbx_invoice_address`.`city_town_village` AS 'pbx_invoice_city_town_village',
  `pbx_invoice_address`.`county_state_region` AS 'pbx_invoice_county_state_region',
  `pbx_invoice_address`.`postcode_zip_code` AS 'pbx_invoice_postcode_zip_code',
  `pbx_invoice_address`.`country` AS 'pbx_invoice_country',
  `pbx_invoice_address`.`contact_email` AS 'pbx_invoice_contact_email',
  `pbx_invoice_address`.`contact_number` AS 'pbx_invoice_contact_number'
FROM `pbx`
INNER JOIN `customer`
ON `pbx`.`customer_id` = `customer`.`id`
INNER JOIN `pbx_site_address`
ON `pbx`.`id` = `pbx_site_address`.`id`
INNER JOIN `pbx_invoice_address`
ON `pbx`.`id` = `pbx_invoice_address`.`id`;

CREATE VIEW `view___sip_extension_detail` AS
SELECT DISTINCT
  `ps_auths`.`username` AS 'sip_username',
  `ps_auths`.`password` AS 'sip_password',
  IFNULL(`ps_endpoints`.`callerid`, '(NOT SET)') AS 'caller_id',
  IFNULL(`ps_endpoints`.`callerid_privacy`, 'allowed_not_screened (DEFAULT)') AS 'caller_id_privacy',
  IFNULL(`ps_endpoints`.`named_call_group`, '(NOT SET)') AS 'call_group',
  `ps_endpoints`.`allow` AS 'codec_allowed',
  IFNULL(`ps_endpoints`.`direct_media`, 'yes (DEFAULT)') AS 'direct_media',
  IFNULL(`ps_endpoints`.`direct_media_method`, 'invite (DEFAULT)') AS 'direct_media_method',
  IFNULL(`ps_endpoints`.`dtmf_mode`, 'rfc4733 (DEFAULT)') AS 'dtmf_mode',
  IFNULL(`ps_endpoints`.`force_rport`, 'yes (DEFAULT)') AS 'force_rport',  
  IFNULL(`ps_endpoints`.`from_user`, '(NOT SET)') AS 'from_sip_header_user',
  IFNULL(`ps_endpoints`.`from_domain`, '(NOT SET)') AS 'from_sip_header_domain',
  IFNULL(`ps_endpoints`.`permit`, '(NOT SET)') AS 'ip_address_allowed',
  IFNULL(`ps_endpoints`.`named_pickup_group`, '(NOT SET)') AS 'pickup_group',
  IFNULL(`ps_endpoints`.`media_encryption`, 'no (RECOMMENDED TO ENABLE TLS OR SETUP A VPN SERVER)') AS 'media_encryption_enabled',
  IFNULL(`ps_endpoints`.`stir_shaken`, 'no (DEFAULT)') AS 'stir_shaken_enabled',
  IFNULL(`ps_endpoints`.`stir_shaken_profile`, '(NOT SET)') AS 'stir_shaken_profile',
  `ps_contacts`.`endpoint` IS NOT NULL AS 'registered',
  `pbx`.`name` AS 'pbx_name',
  `pbx`.`id` AS 'pbx_id',
  `customer`.`name` AS 'customer_name',
  `customer`.`id` AS 'customer_id'
FROM `ps_endpoints`
INNER JOIN `ps_auths`
ON `ps_endpoints`.`id` = `ps_auths`.`id`
INNER JOIN `pbx`
ON `ps_endpoints`.`pbx_id` = `pbx`.`id`
LEFT JOIN `ps_contacts`
on `ps_endpoints`.`id` = `ps_contacts`.`endpoint`
INNER JOIN `customer`
ON `pbx`.`customer_id` = `customer`.`id`
WHERE `ps_endpoints`.`endpoint_type` = 'sip_extension';

CREATE VIEW `view___sip_extension_registered` AS
SELECT
  `ps_auths`.`username` AS 'sip_username',
  `ps_contacts`.`uri`,
  `ps_contacts`.`user_agent`,
  `pbx`.`name` AS 'pbx_name',
  `pbx`.`id` AS 'pbx_id',
  `customer`.`name` AS 'customer_name',
  `customer`.`id` AS 'customer_id'
FROM `ps_endpoints`
INNER JOIN `ps_auths`
ON `ps_endpoints`.`id` = `ps_auths`.`id`
INNER JOIN `pbx`
ON `ps_endpoints`.`pbx_id` = `pbx`.`id`
INNER JOIN `ps_contacts`
on `ps_endpoints`.`id` = `ps_contacts`.`endpoint`
INNER JOIN `customer`
ON `pbx`.`customer_id` = `customer`.`id`
WHERE `ps_endpoints`.`endpoint_type` = 'sip_extension';

----------------------------------------------------------------------------------------------------

-- Insert data to YAP tables

INSERT INTO `uk_sales_tax_rate_lookup` (`uk_sales_tax_rate`)
VALUES
  (20),
  (5),
  (0);

INSERT INTO `uk_sales_tax_status_lookup` (`uk_sales_tax_status`)
VALUES
  ('TAXABLE'),
  ('EXEMPT');

INSERT INTO `consumer_type_lookup` (`consumer_type`)
VALUES
  ('Residentail'),
  ('Sole Trader'),
  ('Partnership'),
  ('Limited Liability Partnership (LLP)'),
  ('Private Limited Company (LTD)'),
  ('Public Limited Company (PLC)'),
  ('Community Interest Company (CIC)'),
  ('n/a');

INSERT INTO `uk_vat_status_lookup` (`uk_vat_status`)
VALUES
  ('Registered'),
  ('Not Registered'),
  ('n/a');

INSERT INTO `good_service_type_lookup` (`good_service_type`)
VALUES
  ('Services'),
  ('Products');

INSERT INTO `customer` (`id`, `name`, `active`, `uk_based`, `consumer_type`, `uk_vat_status`, `reselling_miniutes`, `pbx_limit`)
VALUES (1, 'system', 0, 'n/a', 'n/a', 'n/a', 'n/a', 0);

INSERT INTO `customer_invoice_address` (`id`,	`address_line_1`,	`address_line_2`,	`city_town_village`, `postcode_zip_code`,	`county_state_region`, `country`,	`contact_email`, `contact_number`)
VALUES (1, 'system', 'system', 'system', 'system', 'system', 'system', 'system', 'system');

INSERT INTO `customer_site_address` (`id`,	`address_line_1`,	`address_line_2`,	`city_town_village`, `postcode_zip_code`,	`county_state_region`, `country`,	`contact_email`, `contact_number`)
VALUES (1, 'system', 'system', 'system', 'system', 'system', 'system', 'system', 'system');

INSERT INTO `pbx` (`id`, `name`, `customer_id`, `active`, `sip_extension_limit`)
VALUES (1, 'system', 1, 0, 0);
  
INSERT INTO `pbx_invoice_address` (`id`,	`address_line_1`,	`address_line_2`,	`city_town_village`, `postcode_zip_code`,	`county_state_region`, `country`,	`contact_email`, `contact_number`)
VALUES (1, 'system', 'system', 'system', 'system', 'system', 'system', 'system', 'system');

INSERT INTO `pbx_site_address` (`id`,	`address_line_1`,	`address_line_2`,	`city_town_village`, `postcode_zip_code`,	`county_state_region`, `country`,	`contact_email`, `contact_number`)
VALUES (1, 'system', 'system', 'system', 'system', 'system', 'system', 'system', 'system');

INSERT INTO `user_account_type` (`id`, `type`, `permission`)
VALUES
(100, 'YAP Admin (100)',
'<table>
  <tr>
    <td>
      <b>Key:</b><br>
      &#9989 = Allowed<br>
      &#10060 = Prohibited<br>
      &#9940 = Not Applicable<br>
    </td>
    <td>
      &#9989 View Own User Account<br>
      &#9989 Update Own User Account<br>
      &#9989 Delete Own User Account<br>
      <div class="main-menu-space"></div>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 Create a YAP Admin (100) User Account<br>
      &#9989 View a YAP Admin (100) User Account<br>
      &#9989 Update a YAP Admin (100) User Account<br>
      &#10060 Delete a YAP Admin (100) User Account<br>
    </td>
    <td>
      &#9989 Create a Customer Admin (200) User Account<br>
      &#9989 View a Customer Admin (200) User Account<br>
      &#9989 Update a Customer Admin (200) User Account<br>
      &#9989 Delete a Customer Admin (200) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 Create a Customer Regular (201) User Account<br>
      &#9989 View a Customer Regular (201) User Account<br>
      &#9989 Update a Customer Regular (201) User Account<br>
      &#9989 Delete a Customer Regular (201) User Account<br>
    </td>
    <td>
      &#9989 Create a PBX Admin (300) User Account<br>
      &#9989 View a PBX Admin (300) User Account<br>
      &#9989 Update a PBX Admin (300) User Account<br>
      &#9989 Delete a PBX Admin (300) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 Create a PBX Regular (301) User Account<br>
      &#9989 View a PBX Regular (301) User Account<br>
      &#9989 Update a PBX Regular (301) User Account<br>
      &#9989 Delete a PBX Regular (301) User Account<br>
    </td>
    <td>
      &#9989 Create a PBX Read Only (302) User Account<br>
      &#9989 View a PBX Read Only (302) User Account<br>
      &#9989 Update a PBX Read Only (302) User Account<br>
      &#9989 Delete a PBX Read Only (302) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 Create a Customer Invoice (400) User Account<br>
      &#9989 View a Customer Invoice (400) User Account<br>
      &#9989 Update a Customer Invoice (400) User Account<br>
      &#9989 Delete a Customer Invoice (400) User Account<br>
    </td>
    <td>
      Note: Customer Invoice Accounts Are Read Only<br>
      Accounts for Viewing Services and Goods<br>
      Billed to a Customer.
      <div class="main-menu-space"></div>
    </td>
  </tr>
  <tr>
    <td>
      &#9940 View Own Customer<br>
      &#9940 Update Own Customer<br>
      &#9940 Delete Own Customer<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#9989 Create a Customer<br>
      &#9989 View a Customer<br>
      &#9989 Update a Customer<br>
      &#9989 Delete a Customer<br>
    </td>
  <tr>
    <td>
      &#9940 View Own PBX<br>
      &#9940 Update Own PBX<br>
      &#9940 Delete Own PBX<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#9989 Create a PBX<br>
      &#9989 View a PBX<br>
      &#9989 Update a PBX<br>
      &#9989 Delete a PBX<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 Create a SIP Endpoint<br>
      &#9989 View a SIP Endpoint<br>
      &#9989 Update a SIP Endpoint<br>
      &#9989 Delete a SIP Endpoint<br>
    </td>
    <td>
      &#9989 Create a Customer Invoice<br>
      &#9989 View a Customer Invoice<br>
      &#9989 Update a Customer Invoice<br>
      &#9989 Delete a Customer Invoice<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 View YAP User Account Logs<br>
      &#9989 View Customer Logs<br>
      &#9989 View PBX Logs<br>
      &#9989 Download Logs<br>
    </td>
    <td>
      &#9989 View Server Information<br>
      &#9989 Download Server Information<br>
      &#9989 Set Resource Limits<br>
      &#9989 View Resource Limits<br>
    </td>
  </tr>
</table>'
),
(200, 'Customer Admin (200)',
'<table>
  <tr>
    <td>
      <b>Key:</b><br>
      &#9989 = Allowed<br>
      &#10060 = Prohibited<br>
      &#9940 = Not Applicable<br>
    </td>
    <td>
      &#9989 View Own User Account<br>
      &#9989 Update Own User Account<br>
      &#10060 Delete Own User Account<br>
      <div class="main-menu-space"></div>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a YAP Admin (100) User Account<br>
      &#10060 View a YAP Admin (100) User Account<br>
      &#10060 Update a YAP Admin (100) User Account<br>
      &#10060 Delete a YAP Admin (100) User Account<br>
    </td>
    <td>
      &#10060 Create a Customer Admin (200) User Account<br>
      &#9989 View a Customer Admin (200) User Account<br>
      &#10060 Update a Customer Admin (200) User Account<br>
      &#10060 Delete a Customer Admin (200) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a Customer Regular (201) User Account<br>
      &#9989 View a Customer Regular (201) User Account<br>
      &#9989 Update a Customer Regular (201) User Account<br>
      &#10060 Delete a Customer Regular (201) User Account<br>
    </td>
    <td>
      &#10060 Create a PBX Admin (300) User Account<br>
      &#9989 View a PBX Admin (300) User Account<br>
      &#9989 Update a PBX Admin (300) User Account<br>
      &#10060 Delete a PBX Admin (300) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a PBX Regular (301) User Account<br>
      &#9989 View a PBX Regular (301) User Account<br>
      &#9989 Update a PBX Regular (301) User Account<br>
      &#10060 Delete a PBX Regular (301) User Account<br>
    </td>
    <td>
      &#10060 Create a PBX Read Only (302) User Account<br>
      &#9989 View a PBX Read Only (302) User Account<br>
      &#9989 Update a PBX Read Only (302) User Account<br>
      &#10060 Delete a PBX Read Only (302) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a Customer Invoice (400) User Account<br>
      &#9989 View a Customer Invoice (400) User Account<br>
      &#10060 Update a Customer Invoice (400) User Account<br>
      &#10060 Delete a Customer Invoice (400) User Account<br>
    </td>
    <td>
      Note: Customer Invoice Accounts Are Read Only<br>
      Accounts for Viewing Services and Goods<br>
      Billed to a Customer.
      <div class="main-menu-space"></div>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 View Own Customer<br>
      &#9989 Update Own Customer<br>
      &#10060 Delete Own Customer<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#10060 Create a Customer<br>
      &#10060 View a Customer<br>
      &#10060 Update a Customer<br>
      &#10060 Delete a Customer<br>
    </td>
  <tr>
    <td>
      &#9940 View Own PBX<br>
      &#9940 Update Own PBX<br>
      &#9940 Delete Own PBX<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#9989 Create a PBX<br>
      &#9989 View a PBX<br>
      &#9989 Update a PBX<br>
      &#9989 Delete a PBX<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 Create a SIP Endpoint<br>
      &#9989 View a SIP Endpoint<br>
      &#9989 Update a SIP Endpoint<br>
      &#9989 Delete a SIP Endpoint<br>
    </td>
    <td>
      &#10060 Create a Customer Invoice<br>
      &#9989 View a Customer Invoice<br>
      &#10060 Update a Customer Invoice<br>
      &#10060 Delete a Customer Invoice<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 View YAP User Account Logs<br>
      &#9989 View Customer Logs<br>
      &#9989 View PBX Logs<br>
      &#9989 Download Logs<br>
    </td>
    <td>
      &#10060 View Server Information<br>
      &#10060 Download Server Information<br>
      &#10060 Set Resource Limits<br>
      &#9989 View Resource Limits<br>
    </td>
  </tr>
</table>'
),
(201, 'Customer Regular (201)',
'<table>
  <tr>
    <td>
      <b>Key:</b><br>
      &#9989 = Allowed<br>
      &#10060 = Prohibited<br>
      &#9940 = Not Applicable<br>
    </td>
    <td>
      &#9989 View Own User Account<br>
      &#10060 Update Own User Account<br>
      &#10060 Delete Own User Account<br>
      <div class="main-menu-space"></div>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a YAP Admin (100) User Account<br>
      &#10060 View a YAP Admin (100) User Account<br>
      &#10060 Update a YAP Admin (100) User Account<br>
      &#10060 Delete a YAP Admin (100) User Account<br>
    </td>
    <td>
      &#10060 Create a Customer Admin (200) User Account<br>
      &#10060 View a Customer Admin (200) User Account<br>
      &#10060 Update a Customer Admin (200) User Account<br>
      &#10060 Delete a Customer Admin (200) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a Customer Regular (201) User Account<br>
      &#10060 View a Customer Regular (201) User Account<br>
      &#10060 Update a Customer Regular (201) User Account<br>
      &#10060 Delete a Customer Regular (201) User Account<br>
    </td>
    <td>
      &#10060 Create a PBX Admin (300) User Account<br>
      &#9989 View a PBX Admin (300) User Account<br>
      &#9989 Update a PBX Admin (300) User Account<br>
      &#10060 Delete a PBX Admin (300) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a PBX Regular (301) User Account<br>
      &#9989 View a PBX Regular (301) User Account<br>
      &#9989 Update a PBX Regular (301) User Account<br>
      &#10060 Delete a PBX Regular (301) User Account<br>
    </td>
    <td>
      &#10060 Create a PBX Read Only (302) User Account<br>
      &#9989 View a PBX Read Only (302) User Account<br>
      &#9989 Update a PBX Read Only (302) User Account<br>
      &#10060 Delete a PBX Read Only (302) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a Customer Invoice (400) User Account<br>
      &#10060 View a Customer Invoice (400) User Account<br>
      &#10060 Update a Customer Invoice (400) User Account<br>
      &#10060 Delete a Customer Invoice (400) User Account<br>
    </td>
    <td>
      Note: Customer Invoice Accounts Are Read Only<br>
      Accounts for Viewing Services and Goods<br>
      Billed to a Customer.
      <div class="main-menu-space"></div>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 View Own Customer<br>
      &#10060 Update Own Customer<br>
      &#10060 Delete Own Customer<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#10060 Create a Customer<br>
      &#10060 View a Customer<br>
      &#10060 Update a Customer<br>
      &#10060 Delete a Customer<br>
    </td>
  <tr>
    <td>
      &#9940 View Own PBX<br>
      &#9940 Update Own PBX<br>
      &#9940 Delete Own PBX<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#9989 Create a PBX<br>
      &#9989 View a PBX<br>
      &#9989 Update a PBX<br>
      &#9989 Delete a PBX<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 Create a SIP Endpoint<br>
      &#9989 View a SIP Endpoint<br>
      &#9989 Update a SIP Endpoint<br>
      &#9989 Delete a SIP Endpoint<br>
    </td>
    <td>
      &#10060 Create a Customer Invoice<br>
      &#9989 View a Customer Invoice<br>
      &#10060 Update a Customer Invoice<br>
      &#10060 Delete a Customer Invoice<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 View YAP User Account Logs<br>
      &#10060 View Customer Logs<br>
      &#9989 View PBX Logs<br>
      &#9989 Download Logs<br>
    </td>
    <td>
      &#10060 View Server Information<br>
      &#10060 Download Server Information<br>
      &#10060 Set Resource Limits<br>
      &#9989 View Resource Limits<br>
    </td>
  </tr>
</table>'
),
(300, 'PBX Admin (300)',
'<table>
  <tr>
    <td>
      <b>Key:</b><br>
      &#9989 = Allowed<br>
      &#10060 = Prohibited<br>
      &#9940 = Not Applicable<br>
    </td>
    <td>
      &#9989 View Own User Account<br>
      &#9989 Update Own User Account<br>
      &#10060 Delete Own User Account<br>
      <div class="main-menu-space"></div>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a YAP Admin (100) User Account<br>
      &#10060 View a YAP Admin (100) User Account<br>
      &#10060 Update a YAP Admin (100) User Account<br>
      &#10060 Delete a YAP Admin (100) User Account<br>
    </td>
    <td>
      &#10060 Create a Customer Admin (200) User Account<br>
      &#10060 View a Customer Admin (200) User Account<br>
      &#10060 Update a Customer Admin (200) User Account<br>
      &#10060 Delete a Customer Admin (200) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a Customer Regular (201) User Account<br>
      &#10060 View a Customer Regular (201) User Account<br>
      &#10060 Update a Customer Regular (201) User Account<br>
      &#10060 Delete a Customer Regular (201) User Account<br>
    </td>
    <td>
      &#10060 Create a PBX Admin (300) User Account<br>
      &#9989 View a PBX Admin (300) User Account<br>
      &#10060 Update a PBX Admin (300) User Account<br>
      &#10060 Delete a PBX Admin (300) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a PBX Regular (301) User Account<br>
      &#9989 View a PBX Regular (301) User Account<br>
      &#9989 Update a PBX Regular (301) User Account<br>
      &#10060 Delete a PBX Regular (301) User Account<br>
    </td>
    <td>
      &#10060 Create a PBX Read Only (302) User Account<br>
      &#9989 View a PBX Read Only (302) User Account<br>
      &#9989 Update a PBX Read Only (302) User Account<br>
      &#10060 Delete a PBX Read Only (302) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a Customer Invoice (400) User Account<br>
      &#10060 View a Customer Invoice (400) User Account<br>
      &#10060 Update a Customer Invoice (400) User Account<br>
      &#10060 Delete a Customer Invoice (400) User Account<br>
    </td>
    <td>
      Note: Customer Invoice Accounts Are Read Only<br>
      Accounts for Viewing Services and Goods<br>
      Billed to a Customer.
      <div class="main-menu-space"></div>
    </td>
  </tr>
  <tr>
    <td>
      &#9940 View Own Customer<br>
      &#9940 Update Own Customer<br>
      &#9940 Delete Own Customer<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#10060 Create a Customer<br>
      &#10060 View a Customer<br>
      &#10060 Update a Customer<br>
      &#10060 Delete a Customer<br>
    </td>
  <tr>
    <td>
      &#9989 View Own PBX<br>
      &#9989 Update Own PBX<br>
      &#10060 Delete Own PBX<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#10060 Create a PBX<br>
      &#10060 View a PBX<br>
      &#10060 Update a PBX<br>
      &#10060 Delete a PBX<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 Create a SIP Endpoint<br>
      &#9989 View a SIP Endpoint<br>
      &#9989 Update a SIP Endpoint<br>
      &#9989 Delete a SIP Endpoint<br>
    </td>
    <td>
      &#10060 Create a Customer Invoice<br>
      &#10060 View a Customer Invoice<br>
      &#10060 Update a Customer Invoice<br>
      &#10060 Delete a Customer Invoice<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 View YAP User Account Logs<br>
      &#10060 View Customer Logs<br>
      &#9989 View PBX Logs<br>
      &#9989 Download Logs<br>
    </td>
    <td>
      &#10060 View Server Information<br>
      &#10060 Download Server Information<br>
      &#10060 Set Resource Limits<br>
      &#9989 View Resource Limits<br>
    </td>
  </tr>
</table>'
),
(301, 'PBX Regular (301)',
'<table>
  <tr>
    <td>
      <b>Key:</b><br>
      &#9989 = Allowed<br>
      &#10060 = Prohibited<br>
      &#9940 = Not Applicable<br>
    </td>
    <td>
      &#9989 View Own User Account<br>
      &#10060 Update Own User Account<br>
      &#10060 Delete Own User Account<br>
      <div class="main-menu-space"></div>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a YAP Admin (100) User Account<br>
      &#10060 View a YAP Admin (100) User Account<br>
      &#10060 Update a YAP Admin (100) User Account<br>
      &#10060 Delete a YAP Admin (100) User Account<br>
    </td>
    <td>
      &#10060 Create a Customer Admin (200) User Account<br>
      &#10060 View a Customer Admin (200) User Account<br>
      &#10060 Update a Customer Admin (200) User Account<br>
      &#10060 Delete a Customer Admin (200) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a Customer Regular (201) User Account<br>
      &#10060 View a Customer Regular (201) User Account<br>
      &#10060 Update a Customer Regular (201) User Account<br>
      &#10060 Delete a Customer Regular (201) User Account<br>
    </td>
    <td>
      &#10060 Create a PBX Admin (300) User Account<br>
      &#10060 View a PBX Admin (300) User Account<br>
      &#10060 Update a PBX Admin (300) User Account<br>
      &#10060 Delete a PBX Admin (300) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a PBX Regular (301) User Account<br>
      &#10060 View a PBX Regular (301) User Account<br>
      &#10060 Update a PBX Regular (301) User Account<br>
      &#10060 Delete a PBX Regular (301) User Account<br>
    </td>
    <td>
      &#10060 Create a PBX Read Only (302) User Account<br>
      &#9989 View a PBX Read Only (302) User Account<br>
      &#9989 Update a PBX Read Only (302) User Account<br>
      &#10060 Delete a PBX Read Only (302) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a Customer Invoice (400) User Account<br>
      &#10060 View a Customer Invoice (400) User Account<br>
      &#10060 Update a Customer Invoice (400) User Account<br>
      &#10060 Delete a Customer Invoice (400) User Account<br>
    </td>
    <td>
      Note: Customer Invoice Accounts Are Read Only<br>
      Accounts for Viewing Services and Goods<br>
      Billed to a Customer.
      <div class="main-menu-space"></div>
    </td>
  </tr>
  <tr>
    <td>
      &#9940 View Own Customer<br>
      &#9940 Update Own Customer<br>
      &#9940 Delete Own Customer<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#10060 Create a Customer<br>
      &#10060 View a Customer<br>
      &#10060 Update a Customer<br>
      &#10060 Delete a Customer<br>
    </td>
  <tr>
    <td>
      &#9989 View Own PBX<br>
      &#10060 Update Own PBX<br>
      &#10060 Delete Own PBX<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#10060 Create a PBX<br>
      &#10060 View a PBX<br>
      &#10060 Update a PBX<br>
      &#10060 Delete a PBX<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 Create a SIP Endpoint<br>
      &#9989 View a SIP Endpoint<br>
      &#9989 Update a SIP Endpoint<br>
      &#9989 Delete a SIP Endpoint<br>
    </td>
    <td>
      &#10060 Create a Customer Invoice<br>
      &#10060 View a Customer Invoice<br>
      &#10060 Update a Customer Invoice<br>
      &#10060 Delete a Customer Invoice<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 View YAP User Account Logs<br>
      &#10060 View Customer Logs<br>
      &#10060 View PBX Logs<br>
      &#10060 Download Logs<br>
    </td>
    <td>
      &#10060 View Server Information<br>
      &#10060 Download Server Information<br>
      &#10060 Set Resource Limits<br>
      &#9989 View Resource Limits<br>
    </td>
  </tr>
</table>'
),
(302, 'PBX Read Only (302)',
'<table>
  <tr>
    <td>
      <b>Key:</b><br>
      &#9989 = Allowed<br>
      &#10060 = Prohibited<br>
      &#9940 = Not Applicable<br>
    </td>
    <td>
      &#9989 View Own User Account<br>
      &#10060 Update Own User Account<br>
      &#10060 Delete Own User Account<br>
      <div class="main-menu-space"></div>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a YAP Admin (100) User Account<br>
      &#10060 View a YAP Admin (100) User Account<br>
      &#10060 Update a YAP Admin (100) User Account<br>
      &#10060 Delete a YAP Admin (100) User Account<br>
    </td>
    <td>
      &#10060 Create a Customer Admin (200) User Account<br>
      &#10060 View a Customer Admin (200) User Account<br>
      &#10060 Update a Customer Admin (200) User Account<br>
      &#10060 Delete a Customer Admin (200) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a Customer Regular (201) User Account<br>
      &#10060 View a Customer Regular (201) User Account<br>
      &#10060 Update a Customer Regular (201) User Account<br>
      &#10060 Delete a Customer Regular (201) User Account<br>
    </td>
    <td>
      &#10060 Create a PBX Admin (300) User Account<br>
      &#10060 View a PBX Admin (300) User Account<br>
      &#10060 Update a PBX Admin (300) User Account<br>
      &#10060 Delete a PBX Admin (300) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a PBX Regular (301) User Account<br>
      &#10060 View a PBX Regular (301) User Account<br>
      &#10060 Update a PBX Regular (301) User Account<br>
      &#10060 Delete a PBX Regular (301) User Account<br>
    </td>
    <td>
      &#10060 Create a PBX Read Only (302) User Account<br>
      &#10060 View a PBX Read Only (302) User Account<br>
      &#10060 Update a PBX Read Only (302) User Account<br>
      &#10060 Delete a PBX Read Only (302) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a Customer Invoice (400) User Account<br>
      &#10060 View a Customer Invoice (400) User Account<br>
      &#10060 Update a Customer Invoice (400) User Account<br>
      &#10060 Delete a Customer Invoice (400) User Account<br>
    </td>
    <td>
      Note: Customer Invoice Accounts Are Read Only<br>
      Accounts for Viewing Services and Goods<br>
      Billed to a Customer.
      <div class="main-menu-space"></div>
    </td>
  </tr>
  <tr>
    <td>
      &#9940 View Own Customer<br>
      &#9940 Update Own Customer<br>
      &#9940 Delete Own Customer<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#10060 Create a Customer<br>
      &#10060 View a Customer<br>
      &#10060 Update a Customer<br>
      &#10060 Delete a Customer<br>
    </td>
  <tr>
    <td>
      &#9989 View Own PBX<br>
      &#10060 Update Own PBX<br>
      &#10060 Delete Own PBX<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#10060 Create a PBX<br>
      &#10060 View a PBX<br>
      &#10060 Update a PBX<br>
      &#10060 Delete a PBX<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a SIP Endpoint<br>
      &#9989 View a SIP Endpoint<br>
      &#10060 Update a SIP Endpoint<br>
      &#10060 Delete a SIP Endpoint<br>
    </td>
    <td>
      &#10060 Create a Customer Invoice<br>
      &#10060 View a Customer Invoice<br>
      &#10060 Update a Customer Invoice<br>
      &#10060 Delete a Customer Invoice<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 View YAP User Account Logs<br>
      &#10060 View Customer Logs<br>
      &#10060 View PBX Logs<br>
      &#10060 Download Logs<br>
    </td>
    <td>
      &#10060 View Server Information<br>
      &#10060 Download Server Information<br>
      &#10060 Set Resource Limits<br>
      &#9989 View Resource Limits<br>
    </td>
  </tr>
</table>'
),
(400, 'Group Invoice (400)',
'<table>
  <tr>
    <td>
      <b>Key:</b><br>
      &#9989 = Allowed<br>
      &#10060 = Prohibited<br>
      &#9940 = Not Applicable<br>
    </td>
    <td>
      &#9989 View Own User Account<br>
      &#10060 Update Own User Account<br>
      &#10060 Delete Own User Account<br>
      <div class="main-menu-space"></div>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a YAP Admin (100) User Account<br>
      &#10060 View a YAP Admin (100) User Account<br>
      &#10060 Update a YAP Admin (100) User Account<br>
      &#10060 Delete a YAP Admin (100) User Account<br>
    </td>
    <td>
      &#10060 Create a Customer Admin (200) User Account<br>
      &#10060 View a Customer Admin (200) User Account<br>
      &#10060 Update a Customer Admin (200) User Account<br>
      &#10060 Delete a Customer Admin (200) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a Customer Regular (201) User Account<br>
      &#10060 View a Customer Regular (201) User Account<br>
      &#10060 Update a Customer Regular (201) User Account<br>
      &#10060 Delete a Customer Regular (201) User Account<br>
    </td>
    <td>
      &#10060 Create a PBX Admin (300) User Account<br>
      &#10060 View a PBX Admin (300) User Account<br>
      &#10060 Update a PBX Admin (300) User Account<br>
      &#10060 Delete a PBX Admin (300) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a PBX Regular (301) User Account<br>
      &#10060 View a PBX Regular (301) User Account<br>
      &#10060 Update a PBX Regular (301) User Account<br>
      &#10060 Delete a PBX Regular (301) User Account<br>
    </td>
    <td>
      &#10060 Create a PBX Read Only (302) User Account<br>
      &#10060 View a PBX Read Only (302) User Account<br>
      &#10060 Update a PBX Read Only (302) User Account<br>
      &#10060 Delete a PBX Read Only (302) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a Customer Invoice (400) User Account<br>
      &#10060 View a Customer Invoice (400) User Account<br>
      &#10060 Update a Customer Invoice (400) User Account<br>
      &#10060 Delete a Customer Invoice (400) User Account<br>
    </td>
    <td>
      Note: Customer Invoice Accounts Are Read Only<br>
      Accounts for Viewing Services and Goods<br>
      Billed to a Customer.
      <div class="main-menu-space"></div>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 View Own Customer<br>
      &#10060 Update Own Customer<br>
      &#10060 Delete Own Customer<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#10060 Create a Customer<br>
      &#10060 View a Customer<br>
      &#10060 Update a Customer<br>
      &#10060 Delete a Customer<br>
    </td>
  <tr>
    <td>
      &#9940 View Own PBX<br>
      &#9940 Update Own PBX<br>
      &#9940 Delete Own PBX<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#10060 Create a PBX<br>
      &#10060 View a PBX<br>
      &#10060 Update a PBX<br>
      &#10060 Delete a PBX<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a SIP Endpoint<br>
      &#10060 View a SIP Endpoint<br>
      &#10060 Update a SIP Endpoint<br>
      &#10060 Delete a SIP Endpoint<br>
    </td>
    <td>
      &#10060 Create a Customer Invoice<br>
      &#9989 View a Customer Invoice<br>
      &#10060 Update a Customer Invoice<br>
      &#10060 Delete a Customer Invoice<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 View YAP User Account Logs<br>
      &#10060 View Customer Logs<br>
      &#10060 View PBX Logs<br>
      &#10060 Download Logs<br>
    </td>
    <td>
      &#10060 View Server Information<br>
      &#10060 Download Server Information<br>
      &#10060 Set Resource Limits<br>
      &#10060 View Resource Limits<br>
    </td>
  </tr>
</table>'
);
