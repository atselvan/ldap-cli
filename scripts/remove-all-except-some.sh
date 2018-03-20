#!/bin/bash

# Testing the validations
./ldap-cli -removeAllExceptSome
./ldap-cli -removeAllExceptSome -cn group1
./ldap-cli -removeAllExceptSome -cn group1 -ou nexus

./ldap-cli -removeAllExceptSome -cn group0 -ou nexus -memberIds C00000
./ldap-cli -removeAllExceptSome -cn group1 -ou nex -memberIds C00001
./ldap-cli -removeAllExceptSome -cn group1 -ou nexus -memberIds C00000

# Testing the functionality

./ldap-cli -removeAllExceptSome -cn group1 -ou nexus -memberIds C00002,C00010
