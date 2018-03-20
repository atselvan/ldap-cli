#!/bin/bash

# Testing the validations
./ldap-cli -addMembers
./ldap-cli -addMembers -cn group1
./ldap-cli -addMembers -cn group1 -ou nexus

./ldap-cli -addMembers -cn group0 -ou nexus -memberIds C00000
./ldap-cli -addMembers -cn group1 -ou nex -memberIds C00001
./ldap-cli -addMembers -cn group1 -ou nexus -memberIds C00000

# Testing the functionality

./ldap-cli -addMembers -cn group1 -ou nexus -memberIds C00001
./ldap-cli -addMembers -cn group1 -ou nexus -memberIds C00002
./ldap-cli -addMembers -cn group1 -ou nexus -memberIds C00002,C00003
./ldap-cli -addMembers -cn group1 -ou nexus -memberIds C00004,C00005
./ldap-cli -addMembers -cn group1 -ou nexus -memberIds C00006,C00007,C00008,C00009,C00010