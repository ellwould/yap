<p align="center">
<img src="https://github.com/yet-another-pbx/yap/blob/main/image/yap_logo.jpeg" alt="YAP (Yet Another PBX) Logo">
</p>

<p align="center">
<img src="https://github.com/yet-another-pbx/yap/blob/main/image/yap_software.jpeg" alt="YAP (Yet Another PBX) Software">
</p>

<br>

<h3 align="center">All YAP (Yet Another PBX) Go, BASH, SQL, HTML and CSS code, configuration files, guides and instructions are 100% designed, written and devolped by a human programmer - IN THE INTEREST OF SECURITY AND QUALITY AI WILL NEVER BE USED</h3>

<br>

<p align="center">
<a href="https://no-ai-icon.com/statement/?url=github.com/yet-another-pbx/yap">
<img src="https://github.com/yet-another-pbx/yap/blob/main/image/human_created_content.jpeg" alt="No AI used to create YAP (Yet Another PBX)" width="150" height="150"></a>
&nbsp &nbsp &nbsp &nbsp &nbsp &nbsp &nbsp &nbsp;
<a href="https://no-ai-icon.com/statement/?url=github.com/yet-another-pbx/yap">
<img src="https://github.com/yet-another-pbx/yap/blob/main/image/no_ai.jpeg" alt="YAP (Yet Another PBX) is 100% created by a human" width="150" height="150"></a>
</p>

<br>

# YAP (Yet Another PBX) :telephone_receiver:
YAP (Yet Another PBX) - A GUI to administrate a multi-tenanted SIP Server, automatically creates invoices and connects to 3rd party accountancy software to send invoices. YAP is written in Go and uses the Asterisk framework.

(Tested with Ubuntu version 24.04.3 and Asterisk certified version 20.7-cert8)

<br>

# YAP Website for Manual/Guides
https://yap.ell.today

<br>

# Demo Sites:
- Demo logged in as a YAP Admin (100) account - https://yap100.ell.today/yap

<br>

- Demo logged in as a YAP Customer Admin (200) account - https://yap200.ell.today/yap

<br>

- Demo logged in as a YAP PBX Admin (300) account - https://yap300.ell.today/yap

<br>

## Main Menu Page (Logged in as a YAP Admin):

![YAP (Yet Another PBX) Main Menu Page](https://github.com/yet-another-pbx/yap/blob/main/image/yap_main_menu.jpeg)

<br>

## Customer Page (Logged in as a YAP Admin):

![YAP (Yet Another PBX) Customer Page 1](https://github.com/yet-another-pbx/yap/blob/main/image/customer_page_1.jpeg)


![YAP (Yet Another PBX) Customer Page 3](https://github.com/yet-another-pbx/yap/blob/main/image/customer_page_3.jpeg)


![YAP (Yet Another PBX) Customer Page 4](https://github.com/yet-another-pbx/yap/blob/main/image/customer_page_4.jpeg)

<br>

## PBX Page (Logged in as a YAP Admin):

![YAP (Yet Another PBX) PBX Page 1](https://github.com/yet-another-pbx/yap/blob/main/image/pbx_page_1.jpeg)

![YAP (Yet Another PBX) PBX Page 2](https://github.com/yet-another-pbx/yap/blob/main/image/pbx_page_2.jpeg)

<br>

## Extension Page (Logged in as a YAP Admin):

![YAP (Yet Another PBX) Ext Page 1](https://github.com/yet-another-pbx/yap/blob/main/image/ext_page_1.jpeg)

![YAP (Yet Another PBX) Ext Page 2](https://github.com/yet-another-pbx/yap/blob/main/image/ext_page_2.jpeg)

<br>

## Invoice Page (Logged in as a YAP Admin): 

![YAP (Yet Another PBX) Invoice Page 1](https://github.com/yet-another-pbx/yap/blob/main/image/invoice_page_1.jpeg)

![YAP (Yet Another PBX) Invoice Page 2](https://github.com/yet-another-pbx/yap/blob/main/image/invoice_page_2.jpeg)

<br>

## Service/Product Page (Logged in as a YAP Admin):

![YAP (Yet Another PBX) Service Product Page](https://github.com/yet-another-pbx/yap/blob/main/image/service_product_page.jpeg)

<br>

## Accounting Software Page (Logged in as a YAP Admin):

![YAP (Yet Another PBX) Accounting Software Page](https://github.com/yet-another-pbx/yap/blob/main/image/accounting_software_page.jpeg)

<br>

## YAP User Account Permissions :heavy_check_mark: ( ✅ = Allowed | ❌ = Prohibited | ⛔ = Not Applicable )

|  | YAP Admin (100) | Customer Admin (200) | Customer Regular (201) | PBX Admin (300) | PBX Regular (301) | PBX Read Only (302) |  Customer Invoice (400) |
|-------------------------------------------------|:--:|:--:|:--:|:--:|:--:|:--:|:--:|
| View Own User Account                           | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Update Own User Account                         | ✅\* | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Delete Own User Account                         | ✅\** | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
|  | YAP Admin (100) | Customer Admin (200) | Customer Regular (201) | PBX Admin (300) | PBX Regular (301) | PBX Read Only (302) | Customer Invoice (400) |
| Create a YAP Admin<br>(100) User Account        | ✅\* | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| View a YAP Admin<br>(100) User Account          | ✅\* | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Update a YAP Admin (100)                        | ✅\* | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Delete a YAP Admin (100)                        | ✅\** | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
|  | YAP Admin (100) | Customer Admin (200) | Customer Regular (201) | PBX Admin (300) | PBX Regular (301) | PBX Read Only (302) | Customer Invoice (400) |
| Create a Customer Admin (200)                   | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| View a Customer Admin (200)                     | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Update a Customer Admin (200)                   | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Delete a Customer Admin (200)                   | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
|  | YAP Admin (100) | Customer Admin (200) | Customer Regular (201) | PBX Admin (300) | PBX Regular (301) | PBX Read Only (302) | Customer Invoice (400) |
| Create a Customer Regular<br>(201) User Account | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| View a Customer Regular<br>(201) User Account   | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ |
| Update a Customer Regular<br>(201) User Account | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Delete a Customer Regular<br>(201) User Account | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
|  | YAP Admin (100) | Customer Admin (200) | Customer Regular (201) | PBX Admin (300) | PBX Regular (301) | PBX Read Only (302) | Customer Invoice (400) |
| Create a PBX Admin (300)                        | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| View a PBX Admin<br>(300) User Account          | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ |
| Update a PBX Admin<br>(300) User Account        | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Delete a PBX Admin<br>(300) User Account        | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
|  | YAP Admin (100) | Customer Admin (200) | Customer Regular (201) | PBX Admin (300) | PBX Regular (301) | PBX Read Only (302) | Customer Invoice (400) |
| Create a PBX Regular<br>(301) User Account      | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| View a PBX Regular<br>(301) User Account        | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ |
| Update a PBX Regular<br>(301) User Account      | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Delete a PBX Regular<br>(301) User Account      | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
|  | YAP Admin (100) | Customer Admin (200) | Customer Regular (201) | PBX Admin (300) | PBX Regular (301) | PBX Read Only (302) | Customer Invoice (400) |
| Create a PBX Read Only<br>(302) User Account    | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| View a PBX Read Only<br>(302) User Account      | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ |
| Update a PBX Read Only<br>(302) User Account    | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Delete a PBX Read Only<br>(302) User Account    | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
|  | YAP Admin (100) | Customer Admin (200) | Customer Regular (201) | PBX Admin (300) | PBX Regular (301) | PBX Read Only (302) | Customer Invoice (400) |
| Create a Customer Invoice (400) User Account    | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| View a Customer Invoice (400) User Account      | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Update a Customer Invoice (400) User Account    | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Delete a Customer Invoice (400) User Account    | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
|  | YAP Admin (100) | Customer Admin (200) | Customer Regular (201) | PBX Admin (300) | PBX Regular (301) | PBX Read Only (302) | Customer Invoice (400) |
| View Own Customer                               | ⛔ | ✅ | ✅ | ⛔ | ⛔ | ⛔ | ❌ |
| Update Own Customer                             | ⛔ | ❌ | ❌ | ⛔ | ⛔ | ⛔ | ❌ |
| Delete Own Customer                             | ⛔ | ❌ | ❌ | ⛔ | ⛔ | ⛔ | ❌ |
|  | YAP Admin (100) | Group Admin (200) | Customer Regular (201) | PBX Admin (300) | PBX Regular (301) | PBX Read Only (302) | Customer Invoice (400) |
| Create a Customer                               | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| View a Customer                                 | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Update a Customer                               | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Delete a Customer                               | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
|  | YAP Admin (100) | Customer Admin (200) | Customer Regular (201) | PBX Admin (300) | PBX Regular (301) | PBX Read Only (302) | Customer Invoice (400) |
| View Own PBX                                    | ⛔ | ⛔ | ⛔ | ✅ | ✅ | ✅ | ⛔ |
| Update Own PBX                                  | ⛔ | ⛔ | ⛔ | ❌ | ❌ | ❌ | ⛔ |
| Delete Own PBX                                  | ⛔ | ⛔ | ⛔ | ❌ | ❌ | ❌ | ⛔ |
|  | YAP Admin (100) | Customer Admin (200) | Customer Regular (201) | PBX Admin (300) | PBX Regular (301) | PBX Read Only (302) | Customer Invoice (400) |
| Create Another PBX                              | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ |
| View Another PBX                                | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ |
| Update Another PBX                              | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ |
| Delete Another PBX                              | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
|  | YAP Admin (100) | Customer Admin (200) | Customer Regular (201) | PBX Admin (300) | PBX Regular (301) | PBX Read Only (302) | Customer Invoice (400) |
| Create an Extension                             | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ |
| View an Extension                               | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ |
| Update an Extension                             | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ |
| Delete an Extension                             | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ |
|  | YAP Admin (100) | Customer Admin (200) | Customer Regular (201) | PBX Admin (300) | PBX Regular (301) | PBX Read Only (302) | Customer Invoice (400) |
| Create a Customer Invoice                       | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| View a Customer Invoice                         | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ✅ |
| Update a Customer Invoice                       | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Delete a Customer Invoice                       | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
|  | YAP Admin (100) | Customer Admin (200) | Customer Regular (201) | PBX Admin (300) | PBX Regular (301) | PBX Read Only (302) | Customer Invoice (400) |
| Create a Service/Product                        | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| View a Service/Product                          | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Update a Service/Product                        | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Delete a Service/Product                        | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
|  | YAP Admin (100) | Customer Admin (200) | Customer Regular (201) | PBX Admin (300) | PBX Regular (301) | PBX Read Only (302) | Customer Invoice (400) |
| Create a Supplier                               | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| View a Supplier                                 | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Update a Supplier                               | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Delete a Supplier                               | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
|  | YAP Admin (100) | Customer Admin (200) | Customer Regular (201) | PBX Admin (300) | PBX Regular (301) | PBX Read Only (302) | Customer Invoice (400) |
| Create a Sales Tax Rate                         | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| View a Sales Tax Rate                           | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Update a Sales Tax Rate                         | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Delete a Sales Tax Rate                         | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
|  | YAP Admin (100) | Customer Admin (200) | Customer Regular (201) | PBX Admin (300) | PBX Regular (301) | PBX Read Only (302) | Customer Invoice (400) |
| Connect to the Accounting Software              | ✅\*** | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Send Invoices                                   | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Update a Single Customer Details                | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Update All Customer Details                     | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |

<br>

\*Only the YAP Admin (100) account with account ID 1 can create and delete other YAP Admin (100) accounts
<br>
\**The YAP Admin (100) account with account ID 1 cannot be deleted
<br>
\***Only the YAP Admin (100) account with account ID 1 can connect to the accounting software

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

**Extensions :calling: naming format:**
- 200-13995787344150533 (Endpoint)
- 200-13995787344150533 (AOR)
- 200-13995787344150533 (AUTH)
- inbound_13995787344150533 (Context)

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

**9) Make a directory for the SQLite database**
```
mkdir /var/lib/asterisk/database
```

<br>

**10) Change the SQLite database directory permissions and group:**
```
chmod 570 /var/lib/asterisk/database && chgrp pbx /var/lib/asterisk/database
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

**4) Create a database named asterisk in MariaDB**
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
DNS.1 = localhost
DNS.2 = 127.0.0.1
DNS.3 = (FQDN)
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
<br>

## Adding the YAP self-signed CA certificate to the certificate store :convenience_store::

<br>

**1) Copy the YAP self-signed CA certificate to the certificate store:**
```
cp /etc/mysql/yap-ca.crt /usr/local/share/ca-certificates
```

<br>

**2) Update the CA certificates:**
```
update-ca-certificates
```

<br>

---

<br>
<br>

>[!NOTE]
>For a list of abbreviations and there meanings used throughout this repository please refer to this [README](https://github.com/Ellwould/information_technology_and_telecommunication_abbreviations)

<br>

>[!IMPORTANT]
>All third-party product and/or company names and logos are trademarks™ or registered® trademarks and remain the property of their respective holders/owners. Unless specifically identified as such, use of third party trademarks does not imply any affiliation with or endorsement between YAP (Yet Another PBX) or its contributor(s) by the owners of those trademarks.
