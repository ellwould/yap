# YAP
YAP (Yet Another PBX) - A GUI to administrate a multi-tenanted SIP Server, YAP is written in Go and uses the Asterisk framework.

---

# Asterisk:

### Underscores are used for contexts names because MariaDB does not like hyphens in table names.

### Transport naming format:
- IPv4-UDP
- IPv4-TCP
- IPv4-TLS
- IPv6-UDP
- IPv6-TCP
- IPv6-TLS

### SIP Trunks naming format:
- PBX-330-ST-1 (Endpoint)
- PBX-330-ST-1 (AOR)
- PBX-330-ST-1 (AUTH)
- PBX_330_IN (Context)
<br>

- PBX-330-ST-2 (Endpoint)
- PBX-330-ST-2 (AOR)
- PBX-330-ST-2 (AUTH)
- PBX_330_IN (Context)

### Extensions naming format:
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
