package backend

import (
	"com/privatesquare/go/ldap-cli/ldap"
	m "com/privatesquare/go/ldap-cli/model"
	"log"
)

func ModifyUserDetails(l *ldap.Conn, ldapConn m.LDAPConn, user m.UserDetails){
	// Add a description, and replace the mail attributes
	modify := ldap.NewModifyRequest("uid=ndrake,ou=users,dc=privatesquare,dc=in")
	//modify.Add("description", []string{"An example user"})
	modify.Replace("mail", []string{"user@example.org"})

	err := l.Modify(modify)
	if err != nil {
		log.Fatal(err)
	}
}
