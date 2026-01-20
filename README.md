![image](https://github.com/ellwould/yap/blob/main/yap_logo.jpeg)

![image](https://github.com/ellwould/yap/blob/main/yap_software.jpeg)

<br>

# YAP (Yet Another PBX) :telephone_receiver:
A GUI to administrate a multi-tenanted SIP Server, YAP is written in Go and uses the Asterisk framework.

<br>

(Tested with Ubuntu version 24.04.3 and Asterisk certified version 20.7-cert8)

<br>

## The user_account_type_reference :ledger: Table

| user_account_type | description |
|-------------------|-------------|
| 100 | A YAP admin account can create, read, update and delete all user accounts, groups and PBXs |
| 101 | A YAP regular account can read all user accounts, groups and PBXs |
| 200 | A group admin can read and update thier own PBX(s) and group |
| 201 | A group regular account can read thier own PBX(s) and group |
| 300 | A PBX admin account can read and update thier own PBX |
| 301 | A PBX regular account can read thier own PBX |

<br>

## An Example of the user_account :bust_in_silhouette: Table

| user_id | email             | first_name | last_name    | user_account_type | group_id | pbx_id | date_added | account_active |
|---------|-------------------|------------|--------------|-------------------|----------|--------|------------|----------------|
| 1       | yap@ell.today     | Elliot     | Keavney      | 100               | 1        | 1      | 20/10/2025 | yes            |
| 2       | john@example.com  | John       | not_provided | 200               | 294583   | 1      | 04/11/2025 | yes            |
| 3       | jane@example.net  | Jane       | not_provided | 301               | 294583   | 100    | 03/01/2026 | no             |
| 4       | frank@example.com | Frank      | not_provided | 201               | 294583   | 1      | 02/01/2026 | yes            |

<br>

## An Example of the group :busts_in_silhouette: Table

| group_id | group_name                    | group_site_address_id | group_invoice_address_id | date_added | group_active | note |
|----------|-------------------------------|-----------------------|--------------------------|------------|--------------|------|
| 1        | system                        | 1                     | 1                        | 20/10/2025 | yes          | created during YAP install |
| 294583   | (typically a company name of a VoIP reseller) | 2                     | 2                        | 04/11/2025 | yes          | no notes |

<br>

## An Example of the group_site_address :closed_book: Table

| group_site_address_id | street_road | city_town_village | postcode_zip_code | country | contact_email | contact_number |
|-----------------------|-------------|-------------------------------|-------------------|---------|---------------|----------------|
| 1                     | system      | system            | system            | system  | system        | system         |
| 2                     | (typically the VoIP resellers street/road) | (typically the VoIP resellers city/town/village) | (typically the VoIP resellers postcode/zip code) | (typically the VoIP resellers country) | support@example.com | +441614960000

<br>

## An Example of the group_invoice_address :green_book: Table

| group_invoice_address_id | street_road | city_town_village | postcode_zip_code | country | contact_email | contact_number |
|-----------------------|-------------|-------------------------------|-------------------|---------|---------------|----------------|
| 1                     | system      | system            | system            | system  | system        | system         |
| 2                     | (typically the VoIP resellers street/road) | (typically the VoIP resellers city/town/village) | (typically the VoIP resellers postcode/zip code) | (typically the VoIP resellers country) | accounts@example.com | +441514960000

<br>

## An Example of the pbx :desktop_computer: Table

| pbx_id | pbx_name                                | pbx_site_address_id | pbx_invoice_address_id | group_id | date_added | pbx_active |
|--------|-----------------------------------------------|---------------------|------------------------|----------|------------|------------|
| 1      | system                                        | 1                   | 1                      | 1        | 20/10/2025 | yes        |
| 100    | (typically a company that is a customer of a VoIP reseller) | 2                   | 2                      | 294583   | 02/01/2026 |  yes        |

<br>

## An Example of the pbx_site_address :blue_book: Table

| pbx_site_address_id | street_road | city_town_village | postcode_zip_code | country | contact_email | contact_number |
|-----------------------|-------------|-------------------------------|-------------------|---------|---------------|----------------|
| 1                     | system      | system            | system            | system  | system        | system         |
| 2                     | (customers street/road) | (customers city/town/village) | (customers postcode/zip code) | (customers country) | sales@example.net | +441414960000

<br>

## An Example of the pbx_invoice_address :orange_book: Table

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

## Download, Compile and Install :computer: Asterisk

<br>

>[!NOTE]
>For more detailed instructions on compiling and installing Asterisk please see my guide: [Compiling a Phone System (Asterisk 20.7-cert8) on Ubuntu 24.04.3](https://ellwould.medium.com/compiling-a-phone-system-asterisk-20-5-2-on-ubuntu-22-04-2-bf746b4d862c)

<br>

**1) Install Asterisk dependencies:**
```
apt install unixodbc odbc-mariadb wget build-essential libjansson-dev autoconf libxml2-dev libncurses-dev libedit-dev uuid-dev libsqlite3-dev libnewt-dev automake unixodbc-dev sqlite3 libsrtp2-dev libtool libssl-dev libcurl4-gnutls-dev
```

<br>

**2) Change to the /root directory (must be root):**
```
cd /root
```

<br>

**3) Download Asterisk source code:**
```
wget https://downloads.asterisk.org/pub/telephony/certified-asterisk/asterisk-certified-20.7-cert8.tar.gz
```

<br>

**4) Download the Asterisk teams PGP signature:**
```
wget https://downloads.asterisk.org/pub/telephony/certified-asterisk/asterisk-certified-20.7-cert8.tar.gz.asc
```

<br>

**5) Import the Asterisk teams public key from the Ubuntu key server:**
```
gpg --keyserver keyserver.ubuntu.com --recv 0xF2FC93DB7587BD1FB49E045A5D984BE337191CE7
```

<br>

**6) Verify the compressed tar file against the Asterisk teams PGP signature using GPG:**
```
gpg --verify asterisk-certified-20.7-cert8.tar.gz.asc asterisk-certified-20.7-cert8.tar.gz
```
**Output should show a good signature from the Asterisk team:**
```
gpg: Signature made Mon Jan 12 16:33:21 2026 UTC
gpg:                using RSA key F2FC93DB7587BD1FB49E045A5D984BE337191CE7
gpg: Good signature from "Asterisk Development Team <asteriskteam@digium.com>" [unknown]
gpg: WARNING: This key is not certified with a trusted signature!
gpg:          There is no indication that the signature belongs to the owner.
Primary key fingerprint: F2FC 93DB 7587 BD1F B49E  045A 5D98 4BE3 3719 1CE7
```
>[!WARNING]
>**IF THE OUTPUT SHOWS A BAD SIGNATURE LIKE THIS:**
>```
>gpg: Signature made Mon Jan 12 16:33:21 2026 UTC
>gpg:                using RSA key F2FC93DB7587BD1FB49E045A5D984BE337191CE7
>gpg: BAD signature from "Asterisk Development Team <asteriskteam@digium.com>" [unknown]
>```
>**DELETE THE asterisk-certified-20.7-cert8.tar.gz FILE!!!**

<br>

**7) If the signature was good decompress and untar the Asterisk source code**
```
tar -xvzf asterisk-certified-20.7-cert8.tar.gz
```

<br>

**8) Change the working directory to the asterisk-certified-20.7-cert8 directory:**
```
cd /root/asterisk-certified-20.7-cert8
```

<br>

**9) Run the configure script:**
```
./configure
```

<br>

**10) Run the menu selection system to decide which modules should be compiled:**
```
make menuselect
```

### In menuselect enable the following options:

**Applications:**
- [x] app_voicemail_odbc
- [x] app_attended_transfer
- [x] app_blind_transfer
- [x] app_statsd

**Call Detail Recording:**
- [x] cdr_csv
- [x] cdr_odbc

**PBX Modules:**
- [x] pbx_realtime

**Resource Modules:**
- [x] res_stasis_mailbox
- [x] res_endpoint_stats
- [x] res_pjsip_history
- [x] res_prometheus

**Core Sound Packages:**
- [x] CORE-SOUNDS-EN_GB-WAV
- [x] CORE-SOUNDS-EN_GB-ULAW
- [x] CORE-SOUNDS-EN_GB-ALAW
- [x] CORE-SOUNDS-EN_GB-GSM
- [x] CORE-SOUNDS-EN_GB-G729
- [x] CORE-SOUNDS-EN_GB-G722
- [x] CORE-SOUNDS-EN_GB-SLN16
- [x] CORE-SOUNDS-EN_GB-SIREN7
- [x] CORE-SOUNDS-EN_GB-SIREN14

**Music On Hold File Packages:**
- [x] MOH-OPSOUND-ULAW
- [x] MOH-OPSOUND-ALAW
- [x] MOH-OPSOUND-GSM
- [x] MOH-OPSOUND-G729
- [x] MOH-OPSOUND-G722
- [x] MOH-OPSOUND-SLN16
- [x] MOH-OPSOUND-SIREN7
- [x] MOH-OPSOUND-SIREN14

**Extras Sound Packages:**
- [x] EXTRA-SOUNDS-EN_GB-WAV
- [x] EXTRA-SOUNDS-EN_GB-ULAW
- [x] EXTRA-SOUNDS-EN_GB-ALAW
- [x] EXTRA-SOUNDS-EN_GB-GSM
- [x] EXTRA-SOUNDS-EN_GB-G729
- [x] EXTRA-SOUNDS-EN_GB-G722
- [x] EXTRA-SOUNDS-EN_GB-SLN16
- [x] EXTRA-SOUNDS-EN_GB-SIREN7
- [x] EXTRA-SOUNDS-EN_GB-SIREN14

<br>

**11) Run the configure script again:**
```
./configure
```

<br>

**12) Compile the source code:**
```
make
```

<br>

**13) Install Asterisk:**
```
make install
```

<br>

## Post :arrow_right: Asterisk Install

<br>

**1) Create a system user named pbx with no shell for the Asterisk daemon:**
```
useradd -r -s /bin/false pbx
```

<br>

**2) Lock the pbx user**
```
usermod -L pbx
```

<br>

**3) Edit /etc/asterisk/asterisk.conf so Asterisk runs as the user pbx:**
```
runuser = pbx    ; The user to run as.
rungroup = pbx   ; The group to run as.
```

<br>

**4) Recursively change the group and permissions for Asterisk configuration files:**
```
chown -R root:pbx /etc/asterisk && chmod 550 /etc/asterisk && chmod 440 /etc/asterisk/*
```

<br>

**5) Recursively change the group and permissions for the Asterisk voicemail and call recording files:**
```
chown -R pbx:pbx /var/spool/asterisk && chmod -R 770 /var/spool/asterisk
```

<br>

**6) Recursively change the group and permissions for Asterisk log files**
```
chown -R pbx:pbx /var/log/asterisk && chmod -R 770 /var/log/asterisk
```

<br>

**7) Copy the Asterisk systemd service file:**
```
cp /root/yap/systemd/asterisk.service /usr/lib/systemd/system/
```

<br>

**8) Reload the systemd manager configuration:**
```
systemctl daemon-reload
```

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

<br>

**1) Install the Python interface for MySQL and Alembic**
```
apt install python3-mysqldb alembic
```

<br>

**2) Change the working directory**
```
cd /root/asterisk-certified-20.7-cert8/contrib/ast-db-manage
```

<br>

**3) Copy the sample configuration file and rename it (must add username, password and other details inside the config.ini)**

```
cp config.ini.sample config.ini
```

<br>

**4) Create a database name asterisk in MariaDB**
```
create database asterisk;
```

<br>

**5) Run Alembic**

```
alembic -c config.ini upgrade head
```

<br>

## Generate a Self-Signed ECDSA Key :key: and Certificate :scroll: for MariaDB Using the OpenSSL Cryptographic Library :books:

<br>

>[!CAUTION]
>**A more secure method than the steps listed below is to generate all the files on a seprate computer that is isolated away from the YAP server and then transfer the mariadb.crt, mariadb.key and yap-ca.crt to the YAP server using the SCP (Secure Copy Protocol). For better security the yap-ca.key should not be on the YAP server because it can be used to sign certificate signing requests.**
>

<br>

**1) Make a directory for the files and change into the directory (must be root):**
```
mkdir /root/mariadb-openssl && cd /root/mariadb-openssl
```

<br>

**2) Generate the Certificate Authority (CA) key:**
```
openssl ecparam -genkey -name secp384r1 -out yap-ca.key
```

<br>

**3) Generate a certificate Authority Certificate (CA) with expiry of 7300 days (20 years):**
```
openssl req -x509 -new -SHA384 -nodes -key yap-ca.key -days 7300 -out yap-ca.crt
```

<br>

**4) Generate a key for the MariaDB server:**
```
openssl ecparam -genkey -name secp384r1 -out mariadb.key
```

<br>

**5) Generate a CSR (Certificate Signing Request):**
```
openssl req -new -SHA384 -key mariadb.key -nodes -out mariadb.csr
```

<br>

**6) Generate an extensions file:**
```
touch extensions.ext
```

<br>

**7) The contents of the extensions.ext file:**
```
authorityKeyIdentifier = keyid, issuer
basicConstraints = critical, CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = (FQDN)
```

<br>

**8) Generate and sign the mariadb.crt. It is vaild for 1825 days (5 years):**
```
openssl x509 -req -SHA384 -extfile extensions.ext -days 1825 -in mariadb.csr -CA yap-ca.crt -CAkey yap-ca.key -CAcreateserial -out mariadb.crt
```

<br>

**9) Copy mariadb.crt, mariadb.key and yap-ca.crt to the /etc/mysql directory:**
```
cp /root/mariadb-openssl/{mariadb.crt,mariadb.key,yap-ca.crt} /etc/mysql/
```

<br>

**10) Change the mariadb.key permissions and group:**
```
chmod 440 /etc/mysql/mariadb.key && chgrp mysql /etc/mysql/mariadb.key
```

<br>

>[!TIP]
>**To view the CA certificate in the CLI:**
>```
>openssl x509 -noout -text -in yap-ca.crt
>```
>**To view the MariaDB server certificate in the CLI:**
>```
>openssl x509 -noout -text -in mariadb.crt
>```
>**To view the MariaDB server certificate signing request in the CLI:**
>```
>openssl req -text -noout -verify -in mariadb.csr
>```

<br>

---

<br>
<br>

>[!NOTE]
>For a list of abbreviations and there meanings used throughout this repository please refer to this [README](https://github.com/Ellwould/information_technology_and_telecommunication_abbreviations)

<br>

>[!IMPORTANT]
>All third-party product and/or company names and logos are trademarks™ or registered® trademarks and remain the property of their respective holders/owners. Unless specifically identified as such, use of third party trademarks does not imply any affiliation with or endorsement between Elliot Michael Keavney and the owners of those trademarks.
