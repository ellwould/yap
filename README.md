![image](https://github.com/ellwould/yap/blob/main/yap_logo.jpeg)

![image](https://github.com/ellwould/yap/blob/main/yap_software.jpeg)

<br>

# YAP (Yet Another PBX) :telephone_receiver:
A GUI to administrate a multi-tenanted SIP Server, YAP is written in Go and uses the Asterisk framework.

<br>

## The user_account_type_reference Table

| user_account_type | description |
|-------------------|-------------|
| 100 | A YAP admin account can create, read, update and delete all user accounts, groups and PBXs |
| 101 | A YAP regular account can read all user accounts, groups and PBXs |
| 200 | A group admin can read and update thier own PBX(s) and group |
| 201 | A group regular account can read thier own PBX(s) and group |
| 300 | A PBX admin account can read and update thier own PBX |
| 301 | A PBX regular account can read thier own PBX |

<br>

## An Example of the user_account Table :bust_in_silhouette:

| user_id | email             | first_name | last_name    | user_account_type | group_id | pbx_id | date_added | account_active |
|---------|-------------------|------------|--------------|-------------------|----------|--------|------------|----------------|
| 1       | yap@ell.today     | Elliot     | Keavney      | 100               | 1        | 1      | 20/10/2025 | yes            |
| 2       | john@example.com  | John       | not_provided | 200               | 294583   | 1      | 04/11/2025 | yes            |
| 3       | jane@example.net  | Jane       | not_provided | 301               | 294583   | 100    | 03/01/2026 | no             |
| 4       | frank@example.com | Frank      | not_provided | 201               | 294583   | 1      | 02/01/2026 | yes            |

<br>

## An Example of the group Table

| group_id | group_name                    | group_site_address_id | group_invoice_address_id | date_added | group_active | note |
|----------|-------------------------------|-----------------------|--------------------------|------------|--------------|------|
| 1        | system                        | 1                     | 1                        | 20/10/2025 | yes          | created during YAP install |
| 294583   | (typically a company name of a VoIP reseller) | 2                     | 2                        | 04/11/2025 | yes          | no notes |

<br>

## An Example of the group_site_address Table

| group_site_address_id | street_road | city_town_village | postcode_zip_code | country | contact_email | contact_number |
|-----------------------|-------------|-------------------------------|-------------------|---------|---------------|----------------|
| 1                     | system      | system            | system            | system  | system        | system         |
| 2                     | (typically the VoIP resellers street/road) | (typically the VoIP resellers city/town/village) | (typically the VoIP resellers postcode/zip code) | (typically the VoIP resellers country) | support@example.com | +441614960000

<br>

## An Example of the group_invoice_address Table

| group_invoice_address_id | street_road | city_town_village | postcode_zip_code | country | contact_email | contact_number |
|-----------------------|-------------|-------------------------------|-------------------|---------|---------------|----------------|
| 1                     | system      | system            | system            | system  | system        | system         |
| 2                     | (typically the VoIP resellers street/road) | (typically the VoIP resellers city/town/village) | (typically the VoIP resellers postcode/zip code) | (typically the VoIP resellers country) | accounts@example.com | +441514960000

<br>

## An Example of the pbx Table

| pbx_id | pbx_name                                | pbx_site_address_id | pbx_invoice_address_id | group_id | date_added | pbx_active |
|--------|-----------------------------------------------|---------------------|------------------------|----------|------------|------------|
| 1      | system                                        | 1                   | 1                      | 1        | 20/10/2025 | yes        |
| 100    | (typically a company that is a customer of a VoIP reseller) | 2                   | 2                      | 294583   | 02/01/2026 |  yes        |

<br>

## An Example of the pbx_site_address Table

| pbx_site_address_id | street_road | city_town_village | postcode_zip_code | country | contact_email | contact_number |
|-----------------------|-------------|-------------------------------|-------------------|---------|---------------|----------------|
| 1                     | system      | system            | system            | system  | system        | system         |
| 2                     | (customers street/road) | (customers city/town/village) | (customers postcode/zip code) | (customers country) | sales@example.net | +441414960000

<br>

## An Example of the pbx_invoice_address Table

| pbx_invoice_address_id | street_road | city_town_village | postcode_zip_code | country | contact_email | contact_number |
|-----------------------|-------------|-------------------------------|-------------------|---------|---------------|----------------|
| 1                     | system      | system            | system            | system  | system        | system         |
| 2                     | (customers street/road) | (customers city/town/village) | (customers postcode/zip code) | (customers country) | accounts@example.net | +441414960000

<br>

---

<br>

# Asterisk :telephone:

<br>

## Naming Format :label:

<br>

>[!TIP]
>**Underscores are used for contexts names because MariaDB does not like hyphens in table names.**

<br>

**Transport :taxi: naming format:**
- IPv4-UDP
- IPv4-TCP
- IPv4-TLS
- IPv6-UDP
- IPv6-TCP
- IPv6-TLS

<br>

**SIP Trunks :left_right_arrow: naming format:**
- PBX-330-ST-1 (Endpoint)
- PBX-330-ST-1 (AOR)
- PBX-330-ST-1 (AUTH)
- PBX_330_IN (Context)

<br>

- PBX-330-ST-2 (Endpoint)
- PBX-330-ST-2 (AOR)
- PBX-330-ST-2 (AUTH)
- PBX_330_IN (Context)

<br>

**Extensions :calling: naming format:**
- PBX-330-EXT-200 (Endpoint)
- PBX-330-EXT-200 (AOR)
- PBX-330-EXT-200 (AUTH)
- PBX_330_OUT (Context)

<br>

- PBX-330-EXT-201 (Endpoint)
- PBX-330-EXT-201 (AOR)
- PBX-330-EXT-201 (AUTH)
- PBX_330_OUT (Context)

<br>

## Compile and Install :computer: Asterisk

<br>

```
apt install unixodbc odbc-mariadb wget build-essential libjansson-dev autoconf libxml2-dev libncurses5-dev libedit-dev uuid-dev libsqlite3-dev libnewt-dev automake unixodbc-dev sqlite libsrtp2-dev libtool libssl-dev libcurl4-gnutls-dev
```
```
tar -xvzf asterisk.tar.gz
```
```
cd asterisk
```
```
./configure
```
```
make menuselect
```
```
./configure
```
```
make
```
```
make install
```
```
make samples
```

<br>

## Post :arrow_right: Asterisk Install

<br>

```
cd /etc/asterisk
```
```
mkdir SAMPLES
```
```
mv * ./SAMPLES/
```
```
useradd -r -s /bin/false pbx
```
```
usermod -L pbx
```
```
useradd -r -s /bin/false pbx-dummy
```
```
useradd -L pbx-dummy
```

**In /etc/asterisk/asterisk.conf:**
```
    runuser = pbx    ; The user to run as.
    rungroup = pbx   ; The group to run as.
```
```
cp /root/multi-tenant-asterisk/systemd/asterisk.service /usr/lib/systemd/system/
```
```
systemctl daemon-reload
```

<br>
<br>

>[!TIP]
>**If a firewall :fire::bricks: is enabled, set Asterisk to automatically start :green_circle: at boot:**
>```
>systemctl enable asterisk
>```
>**To stop :red_circle: Asterisk starting at boot:**
>```
>systemctl disable asterisk
>```

<br>

>[!TIP]
>**To enter :door: Asterisk:**
>```
>sudo -u pbx asterisk -rvvvvv
>```

<br>

>[!TIP]
>**To find :mag: and fix :screwdriver: faults :warning: with Asterisk:**
>```
>/usr/sbin/asterisk -mqfv -C /etc/asterisk/asterisk.conf
>```

<br>

---

<br>

# MariaDB :file_cabinet: Setup

<br>

## Python :snake: Alembic :card_index_dividers:

```
apt install python3-mysqldb alembic
```
```
cp /root/asterisk/contrib/ast-db-manage/config.ini /root/asterisk/contrib/ast-db-manage/config.ini.sample
```
```
alembic -c /root/asterisk/contrib/ast-db-manage/config.ini upgrade head
```

<br>

## Generate a Self-Signed ECDSA Key :key: and Certificate :scroll: for MariaDB Using the OpenSSL Cryptographic Library :books:

<br>

**1) Generate the Certificate Authority (CA) key:**
```
openssl ecparam -genkey -name secp384r1 -out yap-ca.key;
```

<br>

**2) Generate a certificate Authority Certificate (CA) with expiry of 7300 days (20 years):**
```
openssl req -x509 -new -SHA384 -nodes -key yap-ca.key -days 7300 -out yap-ca.crt;
```

<br>

**3) Generate a key for the MariaDB server:**
```
openssl ecparam -genkey -name secp384r1 -out mariadb.key;
```

<br>

**4) Generate a CSR (Certificate Signing Request):**
```
openssl req -new -SHA384 -key mariadb.key -nodes -out mariadb.csr;
```

<br>

**5) Generate an extensions file:**
```
touch extensions.ext
```

<br>

**6) The contents of the extensions.ext file:**
```
authorityKeyIdentifier = keyid, issuer
basicConstraints = critical, CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = (FQDN)
```

<br>

**7) Generate and sign the mariadb.crt. It is vaild for 1825 days (5 years):**
```
openssl x509 -req -SHA384 -extfile extensions.ext -days 1825 -in mariadb.csr -CA yap-ca.crt -CAkey yap-ca.key -CAcreateserial -out mariadb.crt;
```

<br>

---

<br>
<br>

>[!NOTE]
>For a list of abbreviations and there meanings used throughout this repository please refer to this [README](https://github.com/Ellwould/information_technology_and_telecommunication_abbreviations)

<br>

>[!IMPORTANT]
>All third-party product and/or company names and logos are trademarks™ or registered® trademarks and remain the property of their respective holders/owners. Unless specifically identified as such, use of third party trademarks does not imply any affiliation with or endorsement between Elliot Michael Keavney and the owners of those trademarks.
