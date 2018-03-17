# Go LDAP CLI

A easy to configure and use go application which can be used as a command line interface to execute CRUD operations on a LDAP server.

## Required Libraries

[go-ldap](https://github.com/go-ldap/ldap)

The go-ldap package is copied over to this repository so that the CLI can be built and executed from a close environment.

[asn1-ber](https://github.com/go-asn1-ber/asn1-ber/tree/v1.2)

This package is a dependency of the go-ladp package. This package is copied over to this repository so that the CLI can be built and executed from a close environment.

## Features

* Add a user
* Delete a Users
* Add a group
* Delete a group
* Add a members to a group
* Remove a members from a group
* User and Group validations while performing the above operations
* Easy configuration

## Configiration File

Basic configuration details for the LDAP connection and based DN can be set in the [ldap-connection.json](./src/com/privatesquare/go/ldap-cli/conf/ldap-connection.json) file. This file is placed in the [conf](./src/com/privatesquare/go/ldap-cli/conf) directory.

connections file format:

```json
{
  "hostname": "192.168.178.2",
  "port": "10389",
  "bindUser": "cn=root,dc=privatesquare,dc=in",
  "baseDN": "dc=privatesquare,dc=in",
  "userBaseDN": "ou=users,dc=privatesquare,dc=in",
  "groupBaseDN": "ou=groups,dc=privatesquare,dc=in"
}
```

## Usage

```bash
Usage of ./ldap-cli:
  -addGroup
    	Add a Group. Required parameter: cn, ou. Optional parameter: memberId
  -addMembers
    	Add a member to a group. Required Parameter: cn, ou, memberId
  -addUser
    	Add a user. Required parameters: uid,cn,sn mail, password
  -bindPassword string
    	LDAP bind password.
  -cn string
    	Name
  -connConfigFile string
    	Connection details of the LDAP server. (default "./conf/ldap-connection.json")
  -deleteGroup
    	Delete a Group. Required Parameter: cn, ou
  -deleteUser
    	Deletes a user. Required parameter: uid
  -mail string
    	Email ID
  -memberIds string
    	Id of a user to be added to a group. Pass comma separated values if you wish to add more than one member
  -ou string
    	Organizational Unit
  -password string
    	User Password
  -removeMembers
    	Remove a member to a group. Required Parameter: cn, ou, memberId
  -sn string
    	Last Name
  -uid string
    	User ID
```

## Examples

Examples for running the cli can be found in the [scripts](./scripts) directory.