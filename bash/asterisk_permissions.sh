#!/bin/bash

# Asterisk Configuration Files
chown -R pbx-dummy:pbx /etc/asterisk;
chmod 050 /etc/asterisk;
chmod 040 /etc/asterisk/*;

# Asterisk Spool
chown -R pbx-dummy:pbx /var/spool/asterisk;
chmod 050 /var/spool/asterisk;
chmod 070 /var/spool/asterisk/{dictate,meetme,monitor,outgoing,recording,system,tmp,voicemail};

# Astersik Log
chown -R pbx-dummy:pbx /var/log/asterisk;
chmod 070 /var/log/asterisk;
chmod 070 /var/log/asterisk/{cdr-csv,cdr-custom,cel-custom};

# Asterisk Lib
chown -R pbx-dummy:pbx /var/lib/asterisk;
find /var/lib/asterisk -type f -exec chmod 040 {} +;
find /var/lib/asterisk -type d -exec chmod 070 {} +;
chmod 060 /var/lib/asterisk/astdb.sqlite3;

# Asterisk Cache
chown pbx-dummy:pbx /var/cache/asterisk;
chmod 070 /var/cache/asterisk;

# Asterisk Binaries
chown root:pbx /usr/sbin/{asterisk,astcanary,astdb2bdb,astdb2sqlite3,astgenkey,astversion,safe_asterisk};
chmod 550 /usr/sbin/{asterisk,astcanary,astdb2bdb,astdb2sqlite3,astgenkey,astversion,safe_asterisk};
