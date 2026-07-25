#!/bin/bash

# Update email script for YAP (Yet Another PBX) admin with account ID 1

#----------------------------------------------------------------------

# Clear Screen
clear_screen="\033[H\033[2J";

# American National Standards Institute (ANSI) reset colour code
reset_colour="\033[0m";

# American National Standards Institute (ANSI) text colour code
text_bold_white="\033[1;37m";

# American National Standards Institute (ANSI) background colour codes
bg_red="\033[41m";
bg_green="\033[42m";
bg_yellow="\033[43m";

#----------------------------------------------------------------------

# Check user is root otherwise exit script
if [ "$EUID" -ne 0 ]
then
  printf $clear_screen;
  printf $bg_yellow;
  printf $text_bold_white;
  printf " ╔══════════════════════════════════════════════════════╗ \n";
  printf " ║ Please run the update YAP admin email script as root ║ \n";
  printf " ╚══════════════════════════════════════════════════════╝ \n";
  printf $reset_colour;
  exit;
fi;

#----------------------------------------------------------------------

# Enter a new email
printf "$clear_screen";
printf "\n";
read -p "    Enter a new email for the YAP admin account with account ID 1 (type exit to stop the update YAP admin email script): " email;
printf "\n";
if [[ $email = "" ]]
then
  printf $clear_screen;
  printf $bg_yellow;
  printf $text_bold_white;
  printf " ╔══════════════════════════╗ \n";
  printf " ║ Email cannot be empty    ║ \n";
  printf " ║ Press return to continue ║ \n";
  printf " ╚══════════════════════════╝ \n";
  printf $reset_colour;
  read -p "";
  source ./update-yap-admin-email.sh;
elif [[ $email = "exit" ]] || [[ $email = "Exit" ]]
then
  exit;
fi;

# Update email for YAP admin user with account ID 1
mysql -u root -e "UPDATE yap.user_account SET email = '$email' WHERE id = 1;";

#----------------------------------------------------------------------

printf $clear_screen;
printf $bg_green;
printf $text_bold_white;
printf " ╔═══════════════════════════════════════╗ \n";
printf " ║ Update YAP admin email script finshed ║ \n";
printf " ╚═══════════════════════════════════════╝ \n";
printf "$reset_colour";
exit;
