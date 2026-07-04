-- Create YAP tables

CREATE TABLE `customer`
(
  `id` VARCHAR(255) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `uk_based` ENUM('yes', 'no', 'n/a') NOT NULL,
  `reselling_minutes` ENUM('no', 'yes', 'n/a') NOT NULL,
  `consumer_type` VARCHAR(255) NOT NULL,
  `uk_vat_registered` ENUM('yes', 'no', 'n/a') NOT NULL,
  `uk_vat_number` VARCHAR(20),
  `pbx_limit` SMALLINT UNSIGNED NOT NULL,
  `pbx_setup_price` DECIMAL(8,2) NOT NULL,
  `pbx_rental_price` DECIMAL(8,2) NOT NULL,
  `pbx_cease_price` DECIMAL(8,2) NOT NULL,
  `pbx_contract_length` VARCHAR(255),
  `sip_ext_setup_price` DECIMAL(8,2) NOT NULL,
  `sip_ext_rental_price` DECIMAL(8,2) NOT NULL,
  `sip_ext_cease_price` DECIMAL(8,2) NOT NULL,
  `sip_ext_contract_length` VARCHAR(255),
  `date_time_added` DATETIME DEFAULT NOW() NOT NULL,
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `consumer_type_lookup` (
  `consumer_type` VARCHAR(255),
  PRIMARY KEY (`consumer_type`)
)
ENGINE = InnoDB;

CREATE TABLE `customer_invoice_address`
(
  `id` VARCHAR(255) NOT NULL,
  `address_line_1` VARCHAR(255),
  `address_line_2` VARCHAR(255),
  `city_town_village` VARCHAR(255),
  `county_state_region` VARCHAR(255),
  `postcode_zip_code` VARCHAR(7),
  `country` VARCHAR(255),
  `contact_email` VARCHAR(255),
  `contact_number` VARCHAR(20),
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `customer_site_address`
(
  `id` VARCHAR(255) NOT NULL,
  `address_line_1` VARCHAR(255),
  `address_line_2` VARCHAR(255),
  `city_town_village` VARCHAR(255),
  `county_state_region` VARCHAR(255),
  `postcode_zip_code` VARCHAR(7),
  `country` VARCHAR(255),
  `contact_email` VARCHAR(255),
  `contact_number` VARCHAR(20),
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `pbx`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `customer_id` VARCHAR(255) NOT NULL,
  `sip_extension_limit` SMALLINT UNSIGNED NOT NULL,
  `date_time_added` DATETIME DEFAULT NOW() NOT NULL,
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `pbx_invoice_address`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `address_line_1` VARCHAR(255),
  `address_line_2` VARCHAR(255),
  `city_town_village` VARCHAR(255),
  `county_state_region` VARCHAR(255),
  `postcode_zip_code` VARCHAR(7),
  `country` VARCHAR(255),
  `contact_email` VARCHAR(255),
  `contact_number` VARCHAR(20),
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `pbx_site_address`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `address_line_1` VARCHAR(255),
  `address_line_2` VARCHAR(255),
  `city_town_village` VARCHAR(255),
  `county_state_region` VARCHAR(255),
  `postcode_zip_code` VARCHAR(7),
  `country` VARCHAR(255),
  `contact_email` VARCHAR(255),
  `contact_number` VARCHAR(20),
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
  `date_time_added` DATETIME DEFAULT NOW() NOT NULL,
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

CREATE TABLE `invoice_item` (
  `id` BIGINT UNSIGNED AUTO_INCREMENT NOT NULL,
  `customer_id` VARCHAR(255) NOT NULL,
  `tag` VARCHAR(255),
  `good_service_name` VARCHAR(255) NOT NULL,
  `sell_price` DECIMAL(8,2) NOT NULL,
  `sales_tax_rate` DECIMAL(5,2) NOT NULL,
  `sales_tax_status` VARCHAR(255) NOT NULL,
  `bill_item_once` ENUM('yes', 'no') NOT NULL,
  `item_on_hold` ENUM('yes', 'no') NOT NULL,
  `contract_length` VARCHAR(255),
  `contract_start_date` DATE,
  `date_time_added` DATETIME DEFAULT NOW() NOT NULL,
  PRIMARY KEY(`id`)
)
ENGINE = InnoDB;
  
CREATE TABLE `sales_tax_rate_lookup` (
  `sales_tax_rate` DECIMAL(5,2),
  PRIMARY KEY(`sales_tax_rate`)
)
ENGINE = InnoDB;
  
CREATE TABLE `sales_tax_status_lookup` (
  `sales_tax_status` VARCHAR(255),
  PRIMARY KEY(`sales_tax_status`)
)
ENGINE = InnoDB;

CREATE TABLE `good_service` (
  `name` VARCHAR(255) NOT NULL,
  `good_service_type` VARCHAR(255) NOT NULL,
  `supplier_name` VARCHAR(255) NOT NULL,
  `contract_length` VARCHAR(255),
  `date_time_added` DATETIME DEFAULT NOW() NOT NULL,
  PRIMARY KEY(`name`)
)
ENGINE = InnoDB;

CREATE TABLE `good_service_type_lookup` (
  `good_service_type` VARCHAR(255),
  PRIMARY KEY(`good_service_type`)
)
ENGINE = InnoDB;

CREATE TABLE `supplier` (
  `name` VARCHAR(255) NOT NULL,
  `date_time_added` DATETIME DEFAULT NOW() NOT NULL,
  PRIMARY KEY(`name`)
)
ENGINE = InnoDB;

CREATE TABLE `contract_length_lookup` (
  `contract_length` VARCHAR(255),
  PRIMARY KEY(`contract_length`)
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
ADD INDEX `index___customer__pbx_contract_length` (`pbx_contract_length`);

ALTER TABLE `customer`
ADD INDEX `index___customer__sip_ext_contract_length` (`sip_ext_contract_length`);

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

ALTER TABLE `invoice_item`
ADD INDEX `index___invoice_item__customer_id` (`customer_id`);

ALTER TABLE `invoice_item`
ADD INDEX `index___invoice_item__sales_tax_rate` (`sales_tax_rate`);

ALTER TABLE `invoice_item`
ADD INDEX `index___invoice_item__sales_tax_status` (`sales_tax_status`);

ALTER TABLE `invoice_item`
ADD INDEX `index___invoice_item__good_service_name` (`good_service_name`);

ALTER TABLE `invoice_item`
ADD INDEX `index___invoice_item__contract_length` (`contract_length`);

ALTER TABLE `good_service`
ADD INDEX `index___good_service__good_service_type` (`good_service_type`);

ALTER TABLE `good_service`
ADD INDEX `index___good_service__supplier_name` (`supplier_name`);

ALTER TABLE `good_service`
ADD INDEX `index___good_service__contract_length` (`contract_length`);

----------------------------------------------------------------------------------------------------

-- Create foreign key constraints

ALTER TABLE `customer`
ADD CONSTRAINT fk___customer___consumer_type_lookup
FOREIGN KEY (`consumer_type`)
REFERENCES `consumer_type_lookup` (`consumer_type`);

ALTER TABLE `customer`
ADD CONSTRAINT fk___customer__pbx_contract_length___contract_length_lookup
FOREIGN KEY (`pbx_contract_length`)
REFERENCES `contract_length_lookup` (`contract_length`);

ALTER TABLE `customer`
ADD CONSTRAINT fk___customer__sip_ext_contract_length___contract_length_lookup
FOREIGN KEY (`sip_ext_contract_length`)
REFERENCES `contract_length_lookup` (`contract_length`);

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

ALTER TABLE `invoice_item`
ADD CONSTRAINT fk___invoice_item___customer_id
FOREIGN KEY (`customer_id`)
REFERENCES `customer` (`id`)
ON DELETE CASCADE;

ALTER TABLE `invoice_item`
ADD CONSTRAINT fk___invoice_item___sales_tax_rate_lookup
FOREIGN KEY (`sales_tax_rate`)
REFERENCES `sales_tax_rate_lookup` (`sales_tax_rate`);

ALTER TABLE `invoice_item`
ADD CONSTRAINT fk___invoice_item___sales_tax_status_lookup
FOREIGN KEY (`sales_tax_status`)
REFERENCES `sales_tax_status_lookup` (`sales_tax_status`);

ALTER TABLE `invoice_item`
ADD CONSTRAINT fk___invoice_item___good_service
FOREIGN KEY (`good_service_name`)
REFERENCES `good_service` (`name`);

ALTER TABLE `invoice_item`
ADD CONSTRAINT fk___invoice_item___contract_length_lookup
FOREIGN KEY (`contract_length`)
REFERENCES `contract_length_lookup` (`contract_length`);

ALTER TABLE `good_service`
ADD CONSTRAINT fk___good_service___good_service_type_lookup
FOREIGN KEY (`good_service_type`)
REFERENCES `good_service_type_lookup` (`good_service_type`);

ALTER TABLE `good_service`
ADD CONSTRAINT fk___good_service___supplier
FOREIGN KEY (`supplier_name`)
REFERENCES `supplier` (`name`);

ALTER TABLE `good_service`
ADD CONSTRAINT fk___good_service___contract_length_lookup
FOREIGN KEY (`contract_length`)
REFERENCES `contract_length_lookup` (`contract_length`);

----------------------------------------------------------------------------------------------------

-- Create View(s)

CREATE VIEW `view___account_detail` AS
SELECT
  `user_account`.`id` AS 'user_account_id',
  `user_account`.`user_account_type_id`,
  `user_account`.`first_name` AS 'user_account_first_name',
  `user_account`.`last_name` AS 'user_account_last_name',
  `user_account`.`email` AS 'user_account_email',
  `user_account_type`.`type` AS 'user_account_type',
  `user_account`.`date_time_added` AS 'user_account_date_time_added',
  `user_account`.`customer_id`,
  `customer`.`name` AS 'customer_name',
  `user_account`.`pbx_id`,
   `pbx`.`name` AS 'pbx_name',
  `user_account_type`.`permission` AS 'user_account_type_permission',
  IFNULL(`customer_site_address`.`address_line_1`, '') AS 'customer_site_address_line_1',
  IFNULL(`customer_site_address`.`address_line_2`, '') AS 'customer_site_address_line_2',
  IFNULL(`customer_site_address`.`city_town_village`, '') AS 'customer_site_city_town_village',
  IFNULL(`customer_site_address`.`county_state_region`, '') AS 'customer_site_county_state_region',
  IFNULL(`customer_site_address`.`postcode_zip_code`, '') AS 'customer_site_postcode_zip_code',
  IFNULL(`customer_site_address`.`country`, '') AS 'customer_site_country',
  IFNULL(`customer_site_address`.`contact_email`, '') AS 'customer_site_contact_email',
  IFNULL(`customer_site_address`.`contact_number`, '') AS 'customer_site_contact_number',
  IFNULL(`customer_invoice_address`.`address_line_1`, '') AS 'customer_invoice_address_line_1',
  IFNULL(`customer_invoice_address`.`address_line_2`, '') AS 'customer_invoice_address_line_2',
  IFNULL(`customer_invoice_address`.`city_town_village`, '') AS 'customer_invoice_city_town_village',
  IFNULL(`customer_invoice_address`.`county_state_region`, '') AS 'customer_invoice_county_state_region',
  IFNULL(`customer_invoice_address`.`postcode_zip_code`, '') AS 'customer_invoice_postcode_zip_code',
  IFNULL(`customer_invoice_address`.`country`, '') AS 'customer_invoice_country',
  IFNULL(`customer_invoice_address`.`contact_email`, '') AS 'customer_invoice_contact_email',
  IFNULL(`customer_invoice_address`.`contact_number`, '') AS 'customer_invoice_contact_number',
  IFNULL(`pbx_site_address`.`address_line_1`, '') AS 'pbx_site_address_line_1',
  IFNULL(`pbx_site_address`.`address_line_2`, '') AS 'pbx_site_address_line_2',
  IFNULL(`pbx_site_address`.`city_town_village`, '') AS 'pbx_site_city_town_village',
  IFNULL(`pbx_site_address`.`county_state_region`, '') AS 'pbx_site_county_state_region',
  IFNULL(`pbx_site_address`.`postcode_zip_code`, '') AS 'pbx_site_postcode_zip_code',
  IFNULL(`pbx_site_address`.`country`, '') AS 'pbx_site_country',
  IFNULL(`pbx_site_address`.`contact_email`, '') AS 'pbx_site_contact_email',
  IFNULL(`pbx_site_address`.`contact_number`, '') AS 'pbx_site_contact_number',
  IFNULL(`pbx_invoice_address`.`address_line_1`, '') AS 'pbx_invoice_address_line_1',
  IFNULL(`pbx_invoice_address`.`address_line_2`, '') AS 'pbx_invoice_address_line_2',
  IFNULL(`pbx_invoice_address`.`city_town_village`, '') AS 'pbx_invoice_city_town_village',
  IFNULL(`pbx_invoice_address`.`county_state_region`, '') AS 'pbx_invoice_county_state_region',
  IFNULL(`pbx_invoice_address`.`postcode_zip_code`, '') AS 'pbx_invoice_postcode_zip_code',
  IFNULL(`pbx_invoice_address`.`country`, '') AS 'pbx_invoice_country',
  IFNULL(`pbx_invoice_address`.`contact_email`, '') AS 'pbx_invoice_contact_email',
  IFNULL(`pbx_invoice_address`.`contact_number`, '') AS 'pbx_invoice_contact_number'
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
  `customer`.`date_time_added` AS 'customer_date_time_added',
  `customer`.`uk_based` AS 'customer_uk_based',
  `customer`.`consumer_type` AS 'customer_consumer_type',
  `customer`.`uk_vat_registered` AS 'customer_uk_vat_registered',
  IFNULL(`customer`.`uk_vat_number`, '') AS 'customer_uk_vat_number',
  `customer`.`reselling_minutes` AS 'customer_reselling_minutes',
  `customer`.`pbx_limit` AS 'customer_pbx_limit',
  `customer`.`pbx_setup_price` AS 'customer_pbx_setup_price',
  `customer`.`pbx_rental_price` AS 'customer_pbx_rental_price',
  `customer`.`pbx_cease_price` AS 'customer_pbx_cease_price',
  IFNULL(`customer`.`pbx_contract_length`, '') AS 'customer_pbx_contract_length',
  `customer`.`sip_ext_setup_price` AS 'customer_sip_ext_setup_price',
  `customer`.`sip_ext_rental_price` AS 'customer_sip_ext_rental_price',
  `customer`.`sip_ext_cease_price` AS 'customer_sip_ext_cease_price',
  IFNULL(`customer`.`sip_ext_contract_length`, '') AS 'customer_sip_ext_contract_length',
  IFNULL(`customer_site_address`.`address_line_1`, '') AS 'customer_site_address_line_1',
  IFNULL(`customer_site_address`.`address_line_2`, '') AS 'customer_site_address_line_2',
  IFNULL(`customer_site_address`.`city_town_village`, '') AS 'customer_site_city_town_village',
  IFNULL(`customer_site_address`.`county_state_region`, '') AS 'customer_site_county_state_region',
  IFNULL(`customer_site_address`.`postcode_zip_code`, '') AS 'customer_site_postcode_zip_code',
  IFNULL(`customer_site_address`.`country`, '') AS 'customer_site_country',
  IFNULL(`customer_site_address`.`contact_email`, '') AS 'customer_site_contact_email',
  IFNULL(`customer_site_address`.`contact_number`, '') AS 'customer_site_contact_number',
  IFNULL(`customer_invoice_address`.`address_line_1`, '') AS 'customer_invoice_address_line_1',
  IFNULL(`customer_invoice_address`.`address_line_2`, '') AS 'customer_invoice_address_line_2',
  IFNULL(`customer_invoice_address`.`city_town_village`, '') AS 'customer_invoice_city_town_village',
  IFNULL(`customer_invoice_address`.`county_state_region`, '') AS 'customer_invoice_county_state_region',
  IFNULL(`customer_invoice_address`.`postcode_zip_code`, '') AS 'customer_invoice_postcode_zip_code',
  IFNULL(`customer_invoice_address`.`country`, '') AS 'customer_invoice_country',
  IFNULL(`customer_invoice_address`.`contact_email`, '') AS 'customer_invoice_contact_email',
  IFNULL(`customer_invoice_address`.`contact_number`, '') AS 'customer_invoice_contact_number'
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
  `pbx`.`date_time_added` AS 'pbx_date_time_added',
  `pbx`.`sip_extension_limit` AS 'pbx_sip_extension_limit',
  IFNULL(`pbx_site_address`.`address_line_1`, '') AS 'pbx_site_address_line_1',
  IFNULL(`pbx_site_address`.`address_line_2`, '') AS 'pbx_site_address_line_2',
  IFNULL(`pbx_site_address`.`city_town_village`, '') AS 'pbx_site_city_town_village',
  IFNULL(`pbx_site_address`.`county_state_region`, '') AS 'pbx_site_county_state_region',
  IFNULL(`pbx_site_address`.`postcode_zip_code`, '') AS 'pbx_site_postcode_zip_code',
  IFNULL(`pbx_site_address`.`country`, '') AS 'pbx_site_country',
  IFNULL(`pbx_site_address`.`contact_email`, '') AS 'pbx_site_contact_email',
  IFNULL(`pbx_site_address`.`contact_number`, '') AS 'pbx_site_contact_number',
  IFNULL(`pbx_invoice_address`.`address_line_1`, '') AS 'pbx_invoice_address_line_1',
  IFNULL(`pbx_invoice_address`.`address_line_2`, '') AS 'pbx_invoice_address_line_2',
  IFNULL(`pbx_invoice_address`.`city_town_village`, '') AS 'pbx_invoice_city_town_village',
  IFNULL(`pbx_invoice_address`.`county_state_region`, '') AS 'pbx_invoice_county_state_region',
  IFNULL(`pbx_invoice_address`.`postcode_zip_code`, '') AS 'pbx_invoice_postcode_zip_code',
  IFNULL(`pbx_invoice_address`.`country`, '') AS 'pbx_invoice_country',
  IFNULL(`pbx_invoice_address`.`contact_email`, '') AS 'pbx_invoice_contact_email',
  IFNULL(`pbx_invoice_address`.`contact_number`, '') AS 'pbx_invoice_contact_number'
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
  `pbx`.`id` AS 'pbx_id',
  `pbx`.`name` AS 'pbx_name',
  `customer`.`id` AS 'customer_id',
  `customer`.`name` AS 'customer_name'
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
  `pbx`.`id` AS 'pbx_id',
  `pbx`.`name` AS 'pbx_name',
  `customer`.`id` AS 'customer_id',
  `customer`.`name` AS 'customer_name'
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

CREATE VIEW `view___invoice_item` AS
SELECT DISTINCT
  `customer`.`id` AS 'customer_id',
  `customer`.`name` AS 'customer_name',
  `customer`.`uk_based` AS 'customer_uk_based',
  `customer`.`reselling_minutes` AS 'customer_reselling_minutes',
  `customer`.`uk_vat_registered` AS 'customer_uk_vat_registered',
  IFNULL(`customer`.`uk_vat_number`, '') AS 'customer_uk_vat_number',
  `invoice_item`.`id` AS 'invoice_item_id',
  IFNULL(`invoice_item`.`tag`, '') AS 'invoice_item_tag',
  `invoice_item`.`sell_price` AS 'invoice_item_sell_price',
  `invoice_item`.`date_time_added` AS 'invoice_item_date_time_added',
  `invoice_item`.`sales_tax_rate` AS 'invoice_item_sales_tax_rate',
  `invoice_item`.`sales_tax_status` AS 'invoice_item_sales_tax_status',
  `invoice_item`.`bill_item_once` AS 'invoice_bill_item_once',
  `invoice_item`.`item_on_hold` AS 'invoice_item_on_hold',
  IFNULL(`invoice_item`.`contract_length`, '') AS 'invoice_item_contract_length',
  IFNULL(`invoice_item`.`contract_start_date`, '') AS 'invoice_item_contract_start_date',
  `good_service`.`name` AS 'good_service_name',
  `good_service`.`good_service_type`,
  `good_service`.`supplier_name` AS 'good_service_supplier_name',
  `good_service`.`contract_length` AS 'good_service_contract_length'
FROM `customer`
INNER JOIN `invoice_item`
ON `invoice_item`.`customer_id` = `customer`.`id`
INNER JOIN `good_service`
ON `good_service`.`name` = `invoice_item`.`good_service_name`;

----------------------------------------------------------------------------------------------------

-- Insert data to YAP tables

INSERT INTO `sales_tax_rate_lookup` (`sales_tax_rate`)
VALUES
  (20),
  (5),
  (0);

INSERT INTO `sales_tax_status_lookup` (`sales_tax_status`)
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

INSERT INTO `good_service_type_lookup` (`good_service_type`)
VALUES
  ('Services'),
  ('Products');

INSERT INTO `contract_length_lookup` (`contract_length`)
VALUES
  ('1 Day'),
  ('1 Week'),
  ('1 Month'),
  ('3 Months'),
  ('6 Months'),
  ('12 Months'),
  ('18 Months'),
  ('24 Months'),
  ('36 Months'),
  ('48 Months'),
  ('60 Months');

INSERT INTO `customer` (`id`, `name`, `uk_based`, `consumer_type`, `uk_vat_registered`, `uk_vat_number`, `reselling_minutes`, `pbx_limit`, `pbx_setup_price`, `pbx_rental_price`, `pbx_cease_price`, `pbx_contract_length`, `sip_ext_setup_price`, `sip_ext_rental_price`, `sip_ext_cease_price`, `sip_ext_contract_length`)
VALUES (1, 'system', 'n/a', 'n/a', 'n/a', 'n/a', NULL, 0, 0, 0, 0, NULL, 0, 0, 0, NULL);

INSERT INTO `customer_invoice_address` (`id`,	`address_line_1`,	`address_line_2`,	`city_town_village`, `postcode_zip_code`,	`county_state_region`, `country`,	`contact_email`, `contact_number`)
VALUES (1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);

INSERT INTO `customer_site_address` (`id`,	`address_line_1`,	`address_line_2`,	`city_town_village`, `postcode_zip_code`,	`county_state_region`, `country`,	`contact_email`, `contact_number`)
VALUES (1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);

INSERT INTO `pbx` (`id`, `name`, `customer_id`, `sip_extension_limit`)
VALUES (1, 'system', 1, 0);
  
INSERT INTO `pbx_invoice_address` (`id`,	`address_line_1`,	`address_line_2`,	`city_town_village`, `postcode_zip_code`,	`county_state_region`, `country`,	`contact_email`, `contact_number`)
VALUES (1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);

INSERT INTO `pbx_site_address` (`id`,	`address_line_1`,	`address_line_2`,	`city_town_village`, `postcode_zip_code`,	`county_state_region`, `country`,	`contact_email`, `contact_number`)
VALUES (1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);

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
      &#9989* Update Own User Account<br>
      &#9989* Delete Own User Account<br>
      <div class="main-menu-space"></div>
    </td>
  </tr>
  <tr>
    <td>
      &#9989* Create a YAP Admin (100) User Account<br>
      &#9989* View a YAP Admin (100) User Account<br>
      &#9989* Update a YAP Admin (100) User Account<br>
      &#9989* Delete a YAP Admin (100) User Account<br>
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
      *Only the YAP Admin (100) account with account ID 1<br>
       can create and delete other YAP Admin (100) accounts<br>
      *The YAP Admin (100) account with account ID 1<br>
       cannot be deleted or edited<br>
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
      &#9989 Create a SIP Extension<br>
      &#9989 View a SIP Extension<br>
      &#9989 Update a SIP Extension<br>
      &#9989 Delete a SIP Extension<br>
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
      &#9989 Create a SIP Extension<br>
      &#9989 View a SIP Extension<br>
      &#9989 Update a SIP Extension<br>
      &#9989 Delete a SIP Extension<br>
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
      &#9989 Create a SIP Extension<br>
      &#9989 View a SIP Extension<br>
      &#9989 Update a SIP Extension<br>
      &#9989 Delete a SIP Extension<br>
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
      &#9989 Create a SIP Extension<br>
      &#9989 View a SIP Extension<br>
      &#9989 Update a SIP Extension<br>
      &#9989 Delete a SIP Extension<br>
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
      &#9989 Create a SIP Extension<br>
      &#9989 View a SIP Extension<br>
      &#9989 Update a SIP Extension<br>
      &#9989 Delete a SIP Extension<br>
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
      &#10060 Create a SIP Extension<br>
      &#9989 View a SIP Extension<br>
      &#10060 Update a SIP Extension<br>
      &#10060 Delete a SIP Extension<br>
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
      &#10060 Create a SIP Extension<br>
      &#10060 View a SIP Extension<br>
      &#10060 Update a SIP Extension<br>
      &#10060 Delete a SIP Extension<br>
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
