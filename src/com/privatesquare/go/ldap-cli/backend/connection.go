package backend

import (
	m "com/privatesquare/go/ldap-cli/model"
	"encoding/json"
	"io/ioutil"
	"log"
)

func GetConnectionDetails(filename, bindPassword string) m.LDAPConn {
	if bindPassword == "" {
		log.Fatal("bindPassword is a required parameter for making a lDAP connection")
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	var ldapConn m.LDAPConn
	json.Unmarshal(data, &ldapConn)

	ldapConn.BindPassword = bindPassword

	return ldapConn
}
