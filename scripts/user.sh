#!/bin/bash

./ldap-cli -addUser -uid C00000
./ldap-cli -addUser -uid C00000 -cn User
./ldap-cli -addUser -uid C00000 -cn User -sn zero
./ldap-cli -addUser -uid C00000 -cn User -sn zero -mail something@privatesquare.in

./ldap-cli -addUser -uid C00001 -cn User -sn One -mail something@privatesquare.in -password welcome
./ldap-cli -addUser -uid C00002 -cn User -sn Two -mail something@privatesquare.in -password welcome
./ldap-cli -addUser -uid C00003 -cn User -sn Three -mail something@privatesquare.in -password welcome
./ldap-cli -addUser -uid C00004 -cn User -sn Four -mail something@privatesquare.in -password welcome
./ldap-cli -addUser -uid C00005 -cn User -sn Five -mail something@privatesquare.in -password welcome
./ldap-cli -addUser -uid C00006 -cn User -sn Six -mail something@privatesquare.in -password welcome
./ldap-cli -addUser -uid C00007 -cn User -sn Seven -mail something@privatesquare.in -password welcome
./ldap-cli -addUser -uid C00008 -cn User -sn Eight -mail something@privatesquare.in -password welcome
./ldap-cli -addUser -uid C00009 -cn User -sn Nine -mail something@privatesquare.in -password welcome
./ldap-cli -addUser -uid C00010 -cn User -sn Ten -mail something@privatesquare.in -password welcome
./ldap-cli -addUser -uid C00011 -cn User -sn Eleven -mail something@privatesquare.in -password welcome
./ldap-cli -addUser -uid C00012 -cn User -sn Twelve -mail something@privatesquare.in -password welcome


./ldap-cli -deleteUser -uid C00011
./ldap-cli -deleteUser -uid C00012
