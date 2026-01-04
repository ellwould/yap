![image](https://github.com/ellwould/yap/blob/main/yap_logo.jpeg)

![image](https://github.com/ellwould/yap/blob/main/yap_software.jpeg)

# YAP
YAP (Yet Another PBX) - A GUI to administrate a multi-tenanted SIP Server, YAP is written in Go and uses the Asterisk framework.

---

# Asterisk:

### Underscores are used for contexts names because MariaDB does not like hyphens in table names.

**Transport naming format:**
- IPv4-UDP
- IPv4-TCP
- IPv4-TLS
- IPv6-UDP
- IPv6-TCP
- IPv6-TLS

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

---

## Compile & Install Asterisk:

- apt install unixodbc odbc-mariadb wget build-essential libjansson-dev autoconf libxml2-dev libncurses5-dev libedit-dev uuid-dev libsqlite3-dev libnewt-dev automake unixodbc-dev sqlite libsrtp2-dev libtool libssl-dev libcurl4-gnutls-dev
- tar -xvzf asterisk.tar.gz
- cd asterisk
- ./configure
- make menuselect
- ./configure
- make
- make install
- make samples

---

## Post Install
- cd /etc/asterisk
- mkdir SAMPLES
- mv * ./SAMPLES/
- useradd -r -s /bin/false pbx
- usermod -L pbx
- useradd -r -s /bin/false pbx-dummy
- useradd -L pbx-dummy
- In /etc/asterisk/asterisk.conf:<br>
    runuser = pbx    ; The user to run as.<br>
    rungroup = pbx   ; The group to run as.<br>
- cp /root/multi-tenant-asterisk/systemd/asterisk.service /usr/lib/systemd/system/
- systemctl daemon-reload

<br>

- systemctl enable asterisk   [IF FIREWALL ENABLED]
- systemctl disable asterisk  [IF FIREWALL ENABLED]

---

## To Find Faults:

/usr/sbin/asterisk -mqfv -C /etc/asterisk/asterisk.conf

---

## Enter Asterisk:

sudo -u pbx asterisk -rvvvvv

---

## alembic:

1) apt install python3-mysqldb alembic<br>
2) cp /root/asterisk/contrib/ast-db-manage/config.ini /root/asterisk/contrib/ast-db-manage/config.ini.sample<br>
3) alembic -c /root/asterisk/contrib/ast-db-manage/config.ini upgrade head<br>

---

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

---

>[!NOTE]
>For a list of abbreviations and there meanings used throughout this repository please refer to this [README](https://github.com/Ellwould/information_technology_and_telecommunication_abbreviations)

<br>

>[!IMPORTANT]
>All third-party product and/or company names and logos are trademarks™ or registered® trademarks and remain the property of their respective holders/owners. Unless specifically identified as such, use of third party trademarks does not imply any affiliation with or endorsement between Elliot Michael Keavney and the owners of those trademarks.
