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

# Generate strong passwords using the OpenSSL cryptographic libary
mariadb_root_password=(`openssl rand -base64 40 | tr "/" a | tr "=" a | tr "+" a`);
mariadb_yap_password=(`openssl rand -base64 40 | tr "/" a | tr "=" a | tr "+" a`);
mariadb_pbx_password=(`openssl rand -base64 40 | tr "/" a | tr "=" a | tr "+" a`);
mariadb_temp_password=(`openssl rand -base64 40 | tr "/" a | tr "=" a | tr "+" a`);

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
  printf " ╔════════════════════╗ \n";
  printf " ║ Please run as root ║ \n";
  printf " ╚════════════════════╝ \n";
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

# Install wget
apt update;
apt install -y wget;

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
  wget -P /root https://go.dev/dl/$go_tar;
  hash_result="$(shasum -a 256 /root/$go_tar | cut -d " " -f 1)"
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

# MariaDB install and setup
apt update;
apt install mariadb-server -y;
systemctl daemon-reload;

# Secure the MariaDB server
mysql -u root -D mysql -e "SET PASSWORD FOR 'root'@'localhost' = PASSWORD('$mariadb_root_password')";
mysql -u root -D mysql -e "DELETE FROM mysql.user WHERE User='';";
mysql -u root -D mysql -e "DELETE FROM mysql.user WHERE User='root' AND Host NOT IN ('localhost', '127.0.0.1', '::1');";
mysql -u root -D mysql -e "DROP DATABASE IF EXISTS test";
mysql -u root -D mysql -e "DELETE FROM mysql.db WHERE Db='test' OR Db='test\\_%';"
mysql -u root -D mysql -e "DELETE FROM mysql.user WHERE User='';";
mysql -u root -e "FLUSH PRIVILEGES;";

# Create a .my.cnf file in /root and populate with account details for MariaDB root
my_cnf=/root/.my.cnf;
touch $my_cnf;
echo "[client]"  >> $my_cnf;
echo "user=root" >> $my_cnf;
echo "password=$mariadb_root_password" >> $my_cnf;
echo "socket=/run/mysqld/mysqld.sock"  >> $my_cnf;

# Create databases
mysql -u root -e "CREATE DATABASE asterisk;";
mysql -u root -e "CREATE DATABASE yap;";
mysql -u root -e "FLUSH PRIVILEGES;";

# Create YAP MaraiDB user
mysql -u root -e "CREATE USER 'yap'@'localhost' IDENTIFIED BY '$mariadb_yap_password';";
mysql -u root -e "FLUSH PRIVILEGES;";

# Add the YAP MariaDB user password to the YAP configuration file
string_update_file="/etc/yap/yap.env";
search_string="<REPLACE_YAP_PASSWORD>";
replace_string="$mariadb_yap_password";
string_update;

# Create PBX MaraiDB user
mysql -u root -e "CREATE USER 'pbx'@'localhost' IDENTIFIED BY '$mariadb_pbx_password';";
mysql -u root -e "FLUSH PRIVILEGES;";

# Start MariaDB on boot
systemctl enable mariadb;

#----------------------------------------------------------------------

# Create a system user named pbx with no shell for the Asterisk daemon
useradd -r -s /bin/false pbx;

# Lock the pbx user
usermod -L pbx;

# Install Astrisk dependencies
apt install -y unixodbc odbc-mariadb wget build-essential libjansson-dev autoconf libxml2-dev libncurses-dev libedit-dev uuid-dev libsqlite3-dev libnewt-dev automake unixodbc-dev sqlite3 libsrtp2-dev libtool libssl-dev libcurl4-gnutls-dev;

# Download Asterisk source code
wget https://downloads.asterisk.org/pub/telephony/certified-asterisk/$asterisk_version.tar.gz;

# Download the Asterisk teams PGP signature
wget https://downloads.asterisk.org/pub/telephony/certified-asterisk/$asterisk_version.tar.gz.asc;

# Import the Asterisk teams public key from the Ubuntu key server
gpg --keyserver keyserver.ubuntu.com --recv 0x$asterisk_key;

# Verify the compressed tar file against the Asterisk teams PGP signature using GPG
gpg --verify $asterisk_version.tar.gz.asc $asterisk_version.tar.gz;

# Conditional statment based on the return code of the GPG command used to verify Asterisk source code
if [ $? != 0 ]
then
  rm /root/asterisk_tar; 
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

# Alembic for Asterisk tables

# Install Alembic and python3-mysqldb
apt install -y python3-mysqldb alembic

# Create temp MariaDB user for Alembic
mysql -u root -e "CREATE USER 'temp'@'localhost' IDENTIFIED BY '$mariadb_temp_password';";
mysql -u root -e "FLUSH PRIVILEGES;";

# Grant permissions for temp MariaDB user
mysql -u root -e "GRANT INSERT, CREATE, ALTER, REFERENCES, INDEX ON asterisk.* TO 'temp'@'localhost';";
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

#----------------------------------------------------------------------

printf $bg_green;
printf $text_bold_white;
printf " ╔═══════════════════════════════════════╗ \n";
printf " ║ YAP has been installed in /usr/bin    ║ \n";
printf " ║ To run YAP type \"systemctl start yap\" ║ \n"; 
printf " ╚═══════════════════════════════════════╝ \n";
printf $reset_colour;
