![image](https://github.com/ellwould/yap/blob/main/yap_logo.jpeg)

![image](https://github.com/ellwould/yap/blob/main/yap_software.jpeg)

<br>

# YAP (Yet Another PBX)
A GUI to administrate a multi-tenanted SIP Server, YAP is written in Go and uses the Asterisk framework.

<br>

---

<br>

# Asterisk

<br>

## Naming Format:

<br>

>[!TIP]
>Underscores are used for contexts names because MariaDB does not like hyphens in table names.

<br>

**Transport naming format:**
- IPv4-UDP
- IPv4-TCP
- IPv4-TLS
- IPv6-UDP
- IPv6-TCP
- IPv6-TLS

<br>

**SIP Trunks naming format:**
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

**Extensions naming format:**
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

## Compile & Install Asterisk:

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

## Post Asterisk Install:

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

**If Firewall Enabled:**
```
systemctl enable asterisk
```
```
systemctl disable asterisk
```

<br>

>[!TIP]
>**To find faults with Asterisk:**
>```
>/usr/sbin/asterisk -mqfv -C /etc/asterisk/asterisk.conf
>```

<br>

>[!TIP]
>**To enter Asterisk:**
>```
>sudo -u pbx asterisk -rvvvvv
>```

<br>

---

<br>

# MariaDB

<br>

## alembic:

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

## Generate a ECDSA Certificate Authority (CA) key and self-signed certificate and generate a ECDSA server key and certificate for MariaDB using the OpenSSL cryptographic library:

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

>[!NOTE]
>For a list of abbreviations and there meanings used throughout this repository please refer to this [README](https://github.com/Ellwould/information_technology_and_telecommunication_abbreviations)

<br>

>[!IMPORTANT]
>All third-party product and/or company names and logos are trademarks™ or registered® trademarks and remain the property of their respective holders/owners. Unless specifically identified as such, use of third party trademarks does not imply any affiliation with or endorsement between Elliot Michael Keavney and the owners of those trademarks.
