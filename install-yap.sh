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
go_tar="go1.26.0.linux-amd64.tar.gz"
go_tar_hash="aac1b08a0fb0c4e0a7c1555beb7b59180b05dfc5a3d62e40e9de90cd42f88235"

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
apt install wget;

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
      printf $clear_screen;
      printf $bg_red;
      printf $text_bold_white;
      printf " ╔══════════════════════════════════════════════════════════╗ \n";
      printf " ║ The hash for $go_tar does not match! ║ \n";
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

cp /root/yap/env/yap.env /etc/yap/yap.env

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

#----------------------------------------------------------------------

printf $bg_green;
printf $text_bold_white;
printf " ╔═══════════════════════════════════════╗ \n";
printf " ║ YAP has been installed in /usr/bin    ║ \n";
printf " ║ To run YAP \"type systemctl start yap\" ║ \n"; 
printf " ╚═══════════════════════════════════════╝ \n";
printf $reset_colour;
