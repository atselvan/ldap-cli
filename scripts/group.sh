#!/bin/bash

./ldap-cli -addGroup -cn group0

./ldap-cli -deleteGroup -cn group1 -ou nexus
./ldap-cli -deleteGroup -cn group2 -ou nexus
./ldap-cli -deleteGroup -cn group3 -ou nexus

./ldap-cli -addGroup -cn group1 -ou nexus -memberId C00001
./ldap-cli -addGroup -cn group2 -ou nexus
./ldap-cli -addGroup -cn group3 -ou nexus -memberId C0000