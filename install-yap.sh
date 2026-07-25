#!/bin/bash

# Install Script for YAP (Yet Another PBX)

#----------------------------------------------------------------------

# Clear Screen
clear_screen="\033[H\033[2J";

# American National Standards Institute (ANSI) reset colour code
reset_colour="\033[0m";

# American National Standards Institute (ANSI) text colour code
text_bold_black="\033[1;30m";
text_bold_white="\033[1;37m";

# American National Standards Institute (ANSI) background colour codes
bg_red="\033[41m";
bg_green="\033[42m";
bg_yellow="\033[43m";
bg_purple="\033[45m";

# Golang version
go_tar="go1.26.0.linux-amd64.tar.gz";
go_tar_hash="aac1b08a0fb0c4e0a7c1555beb7b59180b05dfc5a3d62e40e9de90cd42f88235";

# Variable for Asterisk version
asterisk_version="asterisk-certified-20.7-cert8";

# Variable for Asterisk Public key
asterisk_key="F2FC93DB7587BD1FB49E045A5D984BE337191CE7";

# Function to search and replace text in files
function string_update {
  sed -i "s/$search_string/$replace_string/" $string_update_file;
};

#----------------------------------------------------------------------

# Check user is root otherwise exit script
if [ "$EUID" -ne 0 ]
then
  printf $clear_screen;
  printf $bg_yellow;
  printf $text_bold_white;
  printf " ╔═══════════════════════════════════════════╗ \n";
  printf " ║ Please run the YAP install script as root ║ \n";
  printf " ╚═══════════════════════════════════════════╝ \n";
  printf $reset_colour;
  exit;
fi;

#----------------------------------------------------------------------

# Check YAP has been cloned from GitHub
if [ ! -d "/root/yap" ]
then
  printf $clear_screen;
  printf $bg_red;
  printf $text_bold_white;
  printf " ╔══════════════════════════════════════════════════════════════════════════════╗ \n";
  printf " ║ Directory yap does not exist in /root.                                       ║ \n";
  printf " ║ Please run commands: \"cd /root && git clone https://github.com/ellwould/yap\" ║ \n";
  printf " ║ and run the install script again.                                            ║ \n";
  printf " ╚══════════════════════════════════════════════════════════════════════════════╝ \n";
  printf $reset_colour;
  exit;
fi;

#----------------------------------------------------------------------

# YAP Install Title
printf $clear_screen;
printf "\n";
printf "    ██╗   ██╗  █████╗  ██████╗     ██╗ ███╗   ██╗ ███████╗ ████████╗  █████╗  ██╗      ██╗         ███████╗  ██████╗ ██████╗  ██╗ ██████╗  ████████╗\n";
printf "    ╚██╗ ██╔╝ ██╔══██╗ ██╔══██╗    ██║ ████╗  ██║ ██╔════╝ ╚══██╔══╝ ██╔══██╗ ██║      ██║         ██╔════╝ ██╔════╝ ██╔══██╗ ██║ ██╔══██╗ ╚══██╔══╝\n";
printf "     ╚████╔╝  ███████║ ██████╔╝    ██║ ██╔██╗ ██║ ███████╗    ██║    ███████║ ██║      ██║         ███████╗ ██║      ██████╔╝ ██║ ██████╔╝    ██║   \n";
printf "      ╚██╔╝   ██╔══██║ ██╔═══╝     ██║ ██║╚██╗██║ ╚════██║    ██║    ██╔══██║ ██║      ██║         ╚════██║ ██║      ██╔══██╗ ██║ ██╔═══╝     ██║   \n";
printf "       ██║    ██║  ██║ ██║         ██║ ██║ ╚████║ ███████║    ██║    ██║  ██║ ███████╗ ███████╗    ███████║ ╚██████╗ ██║  ██║ ██║ ██║         ██║   \n";
printf "       ╚═╝    ╚═╝  ╚═╝ ╚═╝         ╚═╝ ╚═╝  ╚═══╝ ╚══════╝    ╚═╝    ╚═╝  ╚═╝ ╚══════╝ ╚══════╝    ╚══════╝  ╚═════╝ ╚═╝  ╚═╝ ╚═╝ ╚═╝         ╚═╝   \n";
printf "\n";

#----------------------------------------------------------------------

# Enter email
read -p "    Enter the email of the YAP admin account with account ID 1 (type exit to stop the YAP install script): " email;
printf "\n";
if [[ $email = "" ]]
then
  printf $clear_screen;
  printf $bg_yellow;
  printf $text_bold_white;
  printf " ╔════════════════════════════╗ \n";
  printf " ║ Email name cannot be empty ║ \n";
  printf " ║ Press return to continue   ║ \n";
  printf " ╚════════════════════════════╝ \n";
  printf $reset_colour;
  printf "\n";
  read -p "";
  source ./install-yap.sh;
elif [[ $email = "exit" ]] || [[ $email = "Exit" ]]
then
  exit;
fi;

# UK VAT registered status 
read -p "    Enter UK VAT registered status, option can be yes/no (type exit to stop the YAP install script): " uk_vat_reg;
printf "\n";
if [[ $uk_vat_reg = "" ]]
then
  printf $clear_screen;
  printf $bg_yellow;
  printf $text_bold_white;
  printf " ╔══════════════════════════════════════════╗ \n";
  printf " ║ UK VAT registered status cannot be empty ║ \n";
  printf " ║ Press return to continue                 ║ \n";
  printf " ╚══════════════════════════════════════════╝ \n";
  printf $reset_colour;
  read -p "";
  source ./install-yap.sh;
elif [[ $uk_vat_reg = "exit" ]] || [[ $uk_vat_reg = "Exit" ]]
then
  exit;
elif [[ $uk_vat_reg = "yes" ]] || [[ $uk_vat_reg = "Yes" ]] || [[ $uk_vat_reg = "YES" ]] || [[ $uk_vat_reg = "y" ]] || [[ $uk_vat_reg = "Y" ]]
then  
  uk_vat_reg="yes";
elif [[ $uk_vat_reg = "no" ]] || [[ $uk_vat_reg = "No" ]] || [[ $uk_vat_reg = "NO" ]] || [[ $uk_vat_reg = "n" ]] || [[ $uk_vat_reg = "N" ]]
then
  uk_vat_reg="no";
else
  printf $clear_screen;
  printf $bg_yellow;
  printf $text_bold_white;
  printf " ╔═══════════════════════════════════════════════════════════════════╗ \n";
  printf " ║ Invalid option for UK VAT registered status, option can be yes/no ║ \n";
  printf " ║ Press return to continue                                          ║ \n";
  printf " ╚═══════════════════════════════════════════════════════════════════╝ \n";
  printf $reset_colour;
  read -p "";
  source ./install-yap.sh
fi;

#----------------------------------------------------------------------

# Generate strong passwords using the OpenSSL cryptographic libary
mariadb_root_password=(`openssl rand -base64 40 | tr "/" a | tr "=" a | tr "+" a`);
mariadb_yap_password=(`openssl rand -base64 40 | tr "/" a | tr "=" a | tr "+" a`);
mariadb_pbx_password=(`openssl rand -base64 40 | tr "/" a | tr "=" a | tr "+" a`);
mariadb_temp_password=(`openssl rand -base64 40 | tr "/" a | tr "=" a | tr "+" a`);

#----------------------------------------------------------------------

apt update;
apt install wget \
            python3-mysqldb \
            alembic \
            unixodbc \
            odbc-mariadb \
            build-essential \
            libjansson-dev \
            autoconf \
            libxml2-dev \
            libncurses-dev \
            libedit-dev \
            uuid-dev \
            libsqlite3-dev \
            libnewt-dev \
            automake \
            unixodbc-dev \
            sqlite3 \
            libsrtp2-dev \
            libtool \
            libssl-dev \
            libcurl4-gnutls-dev \
            mariadb-server -y;
if [ $? != 0 ]
then
  printf $clear_screen;
  printf $bg_red;
  printf $text_bold_white;
  printf " ╔══════════════════════════════════════════════════╗ \n";
  printf " ║ Some or all software required failed to install, ║ \n";
  printf " ║ please re-run the install script again           ║ \n";
  printf " ╚══════════════════════════════════════════════════╝ \n";
  printf $reset_colour;
  exit;
fi;

#----------------------------------------------------------------------

# Remove any previous version of Go, download and install Go 1.26.0
printf $clear_screen;
printf $bg_purple;
printf $text_bold_white;
printf " ╔════════════════════════════════════════════════════════════╗ \n";
printf " ║ The installer removes any previous version of Go installed ║ \n";
printf " ║ Go version 1.26.0 will be downlaoded and installed         ║ \n";
printf " ╚════════════════════════════════════════════════════════════╝ \n";
printf $reset_colour;
printf "\n";
printf $text_bold_black;
read -p "   Would you like to continue? [Yes/No]: " response;
if [ $response == "Yes" ] || [ $response == "yes" ] || [ $response == "YES" ] || [ $response == "Y" ] || [ $response == "y" ]
then
  printf $reset_colour;
  # If the Go source code has not already been download in /root then download it
  if [ ! -f "/root/$go_tar" ]; then
    wget -P /root https://go.dev/dl/$go_tar;
  fi;
  hash_result="$(shasum -a 256 /root/$go_tar | cut -d " " -f 1)";
  if [ $hash_result != $go_tar_hash ]
    then
      rm /root/$go_tar;
      printf $clear_screen;
      printf $bg_red;
      printf $text_bold_white;
      printf " ╔══════════════════════════════════════════════════════════╗ \n";
      printf " ║ The hash for $go_tar does not match! ║ \n";
      printf " ║ The Go source code has been removed ║ \n";
      printf " ╚══════════════════════════════════════════════════════════╝ \n";
      printf $reset_colour;
      exit;
  fi;
  rm -rf /usr/local/go && tar -C /usr/local -xzf /root/$go_tar;
  export PATH=$PATH:/usr/local/go/bin;
  version_installed="$(go version)";
  printf $clear_screen;
  printf $bg_purple;
  printf $text_bold_white;
  printf " ╔═════════════════════════════════╗ \n";
  printf " ║ $version_installed ║ \n";
  printf " ╚═════════════════════════════════╝ \n";
  printf $reset_colour;
else
  printf $reset_colour;
  printf $clear_screen;
  exit;
fi;

#----------------------------------------------------------------------

# Create a system user named yap with no shell, no home directory and lock the account
useradd -r -s /bin/false yap;
usermod -L yap;

# Create Go directories in root home directory for compiling the source code
mkdir -p /root/go/{bin,pkg,src/yap};

# Copy YAP source code
cp /root/yap/go/yap.go /root/go/src/yap/yap.go;

# Export Go
export PATH=$PATH:/usr/local/go/bin;

# Remove old go.mod and create a Go mod file for YAP
rm /root/go/src/yap/go.mod;
cd /root/go/src/yap;
go mod init root/go/src/yap;
go mod tidy;

# Compile yap.go
cd /root/go/src/yap;
go build yap.go;

# Create directores used for configuration and HTML/CSS file
mkdir -p /etc/yap/html-css;

# Copy YAP configuration file
cp /root/yap/env/yap.env /etc/yap/yap.env;

# Copy HTML/CSS start and end files
cp /root/yap/html-css/* /etc/yap/html-css/;

# Change executables file permissions, owner, group and move executables
chown root:yap /root/go/src/yap/yap;
chmod 050 /root/go/src/yap/yap;
mv /root/go/src/yap/yap /usr/bin/yap;

# Change YAP file permissions, owner and group
chown -R root:yap /etc/yap;
chmod 050 /etc/yap;
chmod 050 /etc/yap/html-css;
chmod 040 /etc/yap/yap.env;
chmod 040 /etc/yap/html-css/*;

# Change directroy to /root
cd /root;

# Copy Systemd service file and reload the systemd deamon
cp /root/yap/systemd/yap.service /usr/lib/systemd/system/;
systemctl daemon-reload;

#----------------------------------------------------------------------

# Download, verify and untar the Asterisk source code

# Change to root directory
cd /root;

# If the Asterisk source code has not already been download in /root then download it
if [ ! -f "/root/$asterisk_version.tar.gz" ]; then
  wget https://downloads.asterisk.org/pub/telephony/certified-asterisk/$asterisk_version.tar.gz;
fi;

# If the Asterisk teams PGP signature has not already been download in /root then download it
if [ ! -f "/root/$asterisk_version.tar.gz.asc" ]; then
  wget https://downloads.asterisk.org/pub/telephony/certified-asterisk/$asterisk_version.tar.gz.asc;
fi;

# Import the Asterisk teams public key from the Ubuntu key server
gpg --keyserver keyserver.ubuntu.com --recv 0x$asterisk_key;

# Verify the compressed tar file against the Asterisk teams PGP signature using GPG
gpg --verify $asterisk_version.tar.gz.asc $asterisk_version.tar.gz;

# Conditional statment based on the return code of the GPG command used to verify Asterisk source code
if [ $? != 0 ]
then
  rm /root/$asterisk_tar; 
  printf $clear_screen;
  printf $bg_red;
  printf $text_bold_white;
  printf " ╔═══════════════════════════════════════════════╗ \n";
  printf " ║ Verification failed for Asterisk source code, ║ \n";
  printf " ║ the Asterisk source code has been removed     ║ \n";
  printf " ╚═══════════════════════════════════════════════╝ \n";
  printf $reset_colour;
  exit;
fi;

# Untar
tar -xvzf /root/$asterisk_version.tar.gz;

#----------------------------------------------------------------------

# MariaDB install and setup

# Secure the MariaDB server - equivalent to the mysql_secure_installation shell script 
mysql -u root -D mysql -e "SET PASSWORD FOR 'root'@'localhost' = PASSWORD('$mariadb_root_password')";
mysql -u root -D mysql -e "DELETE FROM mysql.user WHERE User='';";
mysql -u root -D mysql -e "DELETE FROM mysql.user WHERE User='root' AND Host NOT IN ('localhost', '127.0.0.1', '::1');";
mysql -u root -D mysql -e "DROP DATABASE IF EXISTS test";
mysql -u root -D mysql -e "DELETE FROM mysql.db WHERE Db='test' OR Db='test\\_%';"
mysql -u root -D mysql -e "DELETE FROM mysql.user WHERE User='';";
mysql -u root -e "FLUSH PRIVILEGES;";

# Create a .my.cnf file in /root and populate with account details for MariaDB root
my_cnf=/root/.my.cnf;
# If any previous .my.cnf file exists remove it from /root
if [ -f "$my_cnf" ]; then
  rm -r $my_cnf;
fi;
touch $my_cnf;
echo "[client]"  >> $my_cnf;
echo "user=root" >> $my_cnf;
echo "password=$mariadb_root_password" >> $my_cnf;
echo "socket=/run/mysqld/mysqld.sock"  >> $my_cnf;

# Create database - asterisk and YAP tables will be stored in the same database named yap
mysql -u root -e "CREATE DATABASE yap;";
mysql -u root -e "FLUSH PRIVILEGES;";

# Drop any previous temp user and create a temp MariaDB user for Alembic
mysql -u root -e "DROP USER IF EXISTS 'temp'@'localhost';";
mysql -u root -e "CREATE USER 'temp'@'localhost' IDENTIFIED BY '$mariadb_temp_password';";
mysql -u root -e "FLUSH PRIVILEGES;";

# Grant privileges for temp MariaDB user
mysql -u root -e "GRANT SELECT, INSERT, UPDATE, CREATE, ALTER, REFERENCES, INDEX ON yap.* TO 'temp'@'localhost';";
mysql -u root -e "FLUSH PRIVILEGES;";

# Copy Alembic configuration script
cp /root/yap/alembic/config.ini /root/$asterisk_version/contrib/ast-db-manage/;

# Add the temp MariaDB user password to the Alembic configuration file
string_update_file="/root/$asterisk_version/contrib/ast-db-manage/config.ini";
search_string="<REPLACE_TEMP_PASSWORD>";
replace_string="$mariadb_temp_password";
string_update;

cd /root/$asterisk_version/contrib/ast-db-manage;
alembic -c config.ini upgrade head;

# Remove config.ini file
rm /root/$asterisk_version/contrib/ast-db-manage/config.ini;

# Remove temp MariaDB user
mysql -u root -e "DROP USER 'temp'@'localhost';";
mysql -u root -e "FLUSH PRIVILEGES;";

# Create YAP tables and insert data
mysql -u root -D yap -e "SOURCE /root/yap/mariadb/yap.sql;";

# Drop any previous YAP MaraiDB user and create a YAP MaraiDB user for YAP
mysql -u root -e "DROP USER IF EXISTS 'yap'@'localhost';";
mysql -u root -e "CREATE USER 'yap'@'localhost' IDENTIFIED BY '$mariadb_yap_password';";
mysql -u root -e "FLUSH PRIVILEGES;";

# Grant privileges for YAP MariaDB user
mysql -u root -e "GRANT SELECT ON yap.* TO 'yap'@'localhost';";
mysql -u root -e "GRANT INSERT, UPDATE, DELETE ON yap.user_account TO 'yap'@'localhost';";
mysql -u root -e "GRANT INSERT, UPDATE, DELETE ON yap.customer TO 'yap'@'localhost';";
mysql -u root -e "GRANT INSERT, UPDATE, DELETE ON yap.customer_invoice_address TO 'yap'@'localhost';";
mysql -u root -e "GRANT INSERT, UPDATE, DELETE ON yap.customer_site_address TO 'yap'@'localhost';";
mysql -u root -e "GRANT INSERT, UPDATE, DELETE ON yap.pbx TO 'yap'@'localhost';";
mysql -u root -e "GRANT INSERT, UPDATE, DELETE ON yap.pbx_site_address TO 'yap'@'localhost';";
mysql -u root -e "GRANT INSERT, UPDATE, DELETE ON yap.ps_aors TO 'yap'@'localhost';";
mysql -u root -e "GRANT INSERT, UPDATE, DELETE ON yap.ps_auths TO 'yap'@'localhost';";
mysql -u root -e "GRANT INSERT, UPDATE, DELETE ON yap.ps_endpoints TO 'yap'@'localhost';";
mysql -u root -e "GRANT INSERT, DELETE ON yap.invoice_item TO 'yap'@'localhost';";
mysql -u root -e "GRANT INSERT, UPDATE, DELETE ON yap.service_product TO 'yap'@'localhost';";
mysql -u root -e "GRANT INSERT, UPDATE, DELETE ON yap.supplier TO 'yap'@'localhost';";
mysql -u root -e "GRANT INSERT, UPDATE, DELETE ON yap.sales_tax_rate_lookup TO 'yap'@'localhost';";
mysql -u root -e "FLUSH PRIVILEGES;";

# Add the YAP MariaDB user password to the YAP configuration file
string_update_file="/etc/yap/yap.env";
search_string="<REPLACE_YAP_PASSWORD>";
replace_string="$mariadb_yap_password";
string_update;

# Add the UK VAT registered status to the YAP configuration file
string_update_file="/etc/yap/yap.env";
search_string="<REPLACE_VAT_REGISTERED_STATUS>";
replace_string="$uk_vat_reg";
string_update;

# Drop any previous PBX MaraiDB user and create a PBX MaraiDB user for Asterisk
mysql -u root -e "DROP USER IF EXISTS 'pbx'@'localhost';";
mysql -u root -e "CREATE USER 'pbx'@'localhost' IDENTIFIED BY '$mariadb_pbx_password';";
mysql -u root -e "FLUSH PRIVILEGES;";

# Grant privileges for PBX MariaDB user
mysql -u root -e "GRANT SELECT ON yap.iaxfriends TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.meetme TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.musiconhold TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.musiconhold_entry TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.ps_aors TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.ps_asterisk_publications TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.ps_auths TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT, INSERT, UPDATE, DELETE ON yap.ps_contacts TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.ps_domain_aliases TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.ps_endpoint_id_ips TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.ps_endpoints TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.ps_globals TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.ps_inbound_publications TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.ps_outbound_publishes TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.ps_registrations TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.ps_resource_list TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.ps_subscription_persistence TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.ps_systems TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.ps_transports TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.queue_members TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.queue_rules TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.queues TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.sippeers TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.stir_tn TO 'pbx'@'localhost';";
mysql -u root -e "GRANT SELECT ON yap.voicemail TO 'pbx'@'localhost';";
mysql -u root -e "FLUSH PRIVILEGES;";

# Start MariaDB on boot
systemctl enable mariadb;

#----------------------------------------------------------------------

# Create a system user named pbx with no shell for the Asterisk daemon
useradd -r -s /bin/false pbx;

# Lock the pbx user
usermod -L pbx;

cd /root/$asterisk_version;
./configure;

make menuselect.makeopts;
menuselect/menuselect --enable app_voicemail_odbc --enable pbx_realtime menuselect.makeopts;
menuselect/menuselect --enable-category MENUSELECT_CORE_SOUNDS --enable-category MENUSELECT_MOH --enable-category MENUSELECT_EXTRA_SOUNDS menuselect.makeopts;

./configure;

make;
make install;
cd /root;

# Copy Asterisk Configuration files and change file attributes
mkdir /etc/asterisk;
cp /root/yap/asterisk/* /etc/asterisk/;
chown root:root /etc/asterisk;
chmod 555 /etc/asterisk;
chown root:pbx /etc/asterisk/*;
chmod 440 /etc/asterisk/*;

# Make realtime-switch directory and change the attriburtes
mkdir /etc/asterisk/realtime-switch;
chown root:yap /etc/asterisk/realtime-switch;
chmod 575 /etc/asterisk/realtime-switch;

# Recursively change the group and permissions for the Asterisk voicemail and call recording files
chown -R pbx:pbx /var/spool/asterisk;
chmod -R 770 /var/spool/asterisk;

# Recursively change the group and permissions for Asterisk log file
chown -R pbx:pbx /var/log/asterisk;
chmod -R 770 /var/log/asterisk;

# Copy the Asterisk systemd service file
cp /root/yap/systemd/asterisk.service /usr/lib/systemd/system/;

# Reload the systemd manager configuration:
systemctl daemon-reload;

# Make a directory for the SQLite database
mkdir /var/lib/asterisk/database;

# Change the SQLite database directory permissions and group
chmod 570 /var/lib/asterisk/database;
chgrp pbx /var/lib/asterisk/database;

# Set Asterisk to automatically start at boot
systemctl enable asterisk;

#----------------------------------------------------------------------

# Setup ODBC

# Copy ODBC Configuration files
cp /root/yap/odbc/* /etc/;

# Add the PBX MariaDB user password to the ODBC configuration file
string_update_file="/etc/odbc.ini";
search_string="<REPLACE_PBX_PASSWORD>";
replace_string="$mariadb_pbx_password";
string_update;

#----------------------------------------------------------------------

# Create YAP admin user with account ID 1
mysql -u root -e "INSERT INTO yap.user_account (id, email, first_name, last_name, user_account_type_id, customer_id, pbx_id) VALUES (1, '$email', 'YAP', 'Admin', 100, 1, 1);";
mysql -u root -e "FLUSH PRIVILEGES;";

#----------------------------------------------------------------------

printf $bg_green;
printf $text_bold_white;
printf " ╔═══════════════════════════════════════╗ \n";
printf " ║ YAP has been installed in /usr/bin    ║ \n";
printf " ║ To run YAP type \"systemctl start yap\" ║ \n"; 
printf " ╚═══════════════════════════════════════╝ \n";
printf $reset_colour;
exit;
