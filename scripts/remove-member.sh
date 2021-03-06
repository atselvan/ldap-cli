#!/bin/bash

# Testing the validations
./ldap-cli -removeMembers
./ldap-cli -removeMembers -cn group1
./ldap-cli -removeMembers -cn group1 -ou nexus

./ldap-cli -removeMembers -cn group0 -ou nexus -memberIds C00000
./ldap-cli -removeMembers -cn group1 -ou nex -memberIds C00001
./ldap-cli -removeMembers -cn group1 -ou nexus -memberIds C00000

# Testing the functionality

./ldap-cli -removeMembers -cn group1 -ou nexus -memberIds C00002
./ldap-cli -removeMembers -cn group1 -ou nexus -memberIds C00003,C00004,C00005
./ldap-cli -removeMembers -cn group1 -ou nexus -memberIds C00002,C00003,C00004,C00005
./ldap-cli -removeMembers -cn group1 -ou nexus -memberIds C00006,C00007,C00008,C00009,C00010

./ldap-cli -removeMembers -cn group1 -ou nexus -memberIds C00001