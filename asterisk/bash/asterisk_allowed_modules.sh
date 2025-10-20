#!/bin/bash

# Variables
group=pbx
path=/usr/lib/asterisk/modules

chown root:pbx /usr/lib/asterisk;
chmod 050 /usr/lib/asterisk;
chgrp $group $path;
chown root $path;
chmod 050 $path;
chgrp root $path/*;
chown root $path/*;
chmod 040 $path/*;

###############
# App modules #
###############

chgrp $group $path/app_cdr.so;      # Writes ad hoc record to CDR.
chgrp $group $path/app_dial.so;     # Used to connect channels together (i.e., make phone calls).
chgrp $group $path/app_playback.so; # Plays pairs of tones of specified frequencies (DTMF mostly).
#chgrp $group $path/app_read.so;    # Requests input of digits from callers and assigns
                                      # input to a variable.

##################
# Bridge modules #
##################

chgrp $group $path/bridge_simple.so; # Performs simple channel-to-channel bridging.

###############
# CDR modules #
###############

chgrp $group $path/cdr_adaptive_odbc.so; # Allows writing of CDRs through ODBC 
                                           # framework with ability to add custom fields.

###############
# CEL modules #
###############

###################
# Channel drivers #
###################

chgrp $group $path/chan_pjsip.so; # Session Initiation Protocol (SIP) channel driver.

#####################
# Codec Translators #
#####################

chgrp $group $path/codec_alaw.so; # A-law PCM codec used all over the world on 
                                    # the PSTN (except Canada/USA).
chgrp $group $path/codec_gsm.so;  # Global System for Mobile Communications (GSM) codec. 
                                    # Very poor sound quality.

######################
# Dialplan Functions #
######################

chgrp $group $path/func_callerid.so; # Gets/sets caller ID.
chgrp $group $path/func_cdr.so;      # Gets/sets CDR variable.
chgrp $group $path/func_odbc.so;     # ODBC lookups.
#chgrp $group $path/func_rand.so;    # Random number dialplan function
#chgrp $group $path/func_sha1.so;    # SHA-1 computation dialplan function
chgrp $group $path/func_strings.so;  # Includes over a dozen string manipulation functions.
#chgrp $group $path/func_logic.so;   # Includes ISNULL(), SET(), EXISTS(), IF(), 
                                       # IFTIME(), and IMPORT(); performs various logical functions.

#######################
# Format Interpreters #
#######################

chgrp $group $path/format_gsm.so;     # RPE-LTP (original GSM codec): .gsm
chgrp $group $path/format_wav.so;     # .wav
chgrp $group $path/format_wav_gsm.so; # GSM audio in a WAV container: .wav, .wav49

###############
# PBX Modules #
###############

chgrp $group $path/pbx_config.so;    # Provides the traditional dialplan 
                                       # language for Asterisk. Without this 
                                       # module, Asterisk cannot read extensions.conf.
# chgrp $group $path/pbx_loopback.so; # Realtime Switch
chgrp $group $path/pbx_realtime.so; # Realtime Switch - OBSOLETE! USE func_odbc.so

####################
# Resource Modules #
####################

########################
# res Calender modules #
########################

####################
# res odbc modules #
####################

chgrp $group $path/res_config_odbc.so;      # Pulls configuration information using ODBC.
chgrp $group $path/res_odbc.so;             # Provides common subroutines for other ODBC modules.
chgrp $group $path/res_odbc_transaction.so; # ODBC transaction resource (func_odbc.so dependancy).

#####################
# res other modules #
#####################

chgrp $group $path/res_clialiases.so; # This module provides the capability to create aliases to other CLI commands.

#####################
# res pjsip modules #
#####################

chgrp $group $path/res_pjproject.so; # PJPROJECT Log and Utility Support. 

chgrp $group $path/res_pjsip.so;                               # Basic SIP resource.
chgrp $group $path/res_pjsip_authenticator_digest.so;          # PJSIP authentication resource. 
chgrp $group $path/res_pjsip_caller_id.so;                     # PJSIP Caller ID Support.
chgrp $group $path/res_pjsip_endpoint_identifier_anonymous.so; # PJSIP Anonymous endpoint identifier.
chgrp $group $path/res_pjsip_endpoint_identifier_ip.so;        # PJSIP IP endpoint identifier.
chgrp $group $path/res_pjsip_endpoint_identifier_user.so;      # PJSIP username endpoint identifier.
chgrp $group $path/res_pjsip_outbound_authenticator_digest.so; # PJSIP authentication resource.
chgrp $group $path/res_pjsip_outbound_publish.so;              # PJSIP Outbound Publish Support.
chgrp $group $path/res_pjsip_outbound_registration.so;         # PJSIP Outbound Registration Support.
chgrp $group $path/res_pjsip_pubsub.so;                        # PJSIP event resource.
chgrp $group $path/res_pjsip_registrar.so;  # PJSIP Registrar Support.            
chgrp $group $path/res_pjsip_sdp_rtp.so;    # PJSIP SDP RTP/AVP stream handler.
chgrp $group $path/res_pjsip_session.so;    # PJSIP Session resource. 

chgrp $group $path/res_rtp_asterisk.so;     # Asterisk RTP Stack.

chgrp $group $path/res_sorcery_astdb.so;    # Sorcery Astdb Object Wizard. 
chgrp $group $path/res_sorcery_config.so;   # Sorcery Configuration File Object Wizard.
chgrp $group $path/res_sorcery_memory.so;   # Sorcery In-Memory Object Wizard.
chgrp $group $path/res_sorcery_realtime.so; # Sorcery Realtime Object Wizard.
chgrp $group $path/res_timing_timerfd.so;   # Timerfd Timing Interface.
