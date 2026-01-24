ALTER TABLE `group_invoice_address`
ADD CONSTRAINT fk_____group_invoice_address__id_____group__id
FOREIGN KEY (`id`)
REFERENCES `group` (`id`)
ON DELETE CASCADE;

ALTER TABLE `group_site_address`
ADD CONSTRAINT fk_____group_site_address__id_____group__id
FOREIGN KEY (`id`)
REFERENCES `group` (`id`)
ON DELETE CASCADE;

ALTER TABLE `pbx_invoice_address`
ADD CONSTRAINT fk_____pbx_invoice_address__id_____group__id
FOREIGN KEY (`id`)
REFERENCES `group` (`id`)
ON DELETE CASCADE;

ALTER TABLE `pbx_site_address`
ADD CONSTRAINT fk_____pbx_site_address__id_____group__id
FOREIGN KEY (`id`)
REFERENCES `group` (`id`)
ON DELETE CASCADE;

ALTER TABLE `user_account`
ADD CONSTRAINT fk_____user_account__user_account_type_id_____user_account_type__id
FOREIGN KEY (`id`)
REFERENCES `user_account_type` (`id`)
