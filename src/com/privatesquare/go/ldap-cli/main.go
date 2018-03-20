package main

import (
	"com/privatesquare/go/ldap-cli/ldap"
	b "com/privatesquare/go/ldap-cli/backend"
	m "com/privatesquare/go/ldap-cli/model"
	"fmt"
	"log"
	"flag"
)

func main() {

	//action
	addUser := flag.Bool("addUser", false, "Add a user. Required parameters: uid,cn,sn mail, password")
	deleteUser := flag.Bool("deleteUser", false, "Deletes a user. Required parameter: uid")
	addGroup := flag.Bool("addGroup", false, "Add a Group. Required parameter: cn, ou. Optional parameter: memberId")
	deleteGroup := flag.Bool("deleteGroup", false, "Delete a Group. Required Parameter: cn, ou")
	addMembers := flag.Bool("addMembers", false, "Add a member to a group. Required Parameter: cn, ou, memberId")
	removeMembers := flag.Bool("removeMembers", false, "Remove a member to a group. Required Parameter: cn, ou, memberId")
	//parameters
	connConfigFile := flag.String("connConfigFile", "./conf/ldap-connection.json", "Connection details of the LDAP server.")
	bindPassword := flag.String("bindPassword", "welkom", "LDAP bind password.")
	uid := flag.String("uid", "", "User ID")
	cn := flag.String("cn", "", "Name")
	sn := flag.String("sn", "", "Last Name")
	mail := flag.String("mail", "", "Email ID")
	password := flag.String("password", "", "User Password")
	ou := flag.String("ou", "", "Organizational Unit")
	memberIds := flag.String("memberIds", "", "Id of a user to be added to a group. Pass comma separated values if you wish to add more than one member")
	flag.Parse()

	ldapConn := b.GetConnectionDetails(*connConfigFile, *bindPassword)

	// Make LDAP connection
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", ldapConn.Hostname, ldapConn.Port))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	err = l.Bind(ldapConn.BindUser, ldapConn.BindPassword)
	if err != nil {
		log.Fatal(err)
	}

	if *addUser == true {
		userDetails := m.UserDetails{Uid: *uid, Cn: *cn, Sn: *sn, Mail: *mail, UserPassword: *password}
		b.AddUser(l, ldapConn, userDetails)
	} else if *deleteUser == true {
		b.DeleteUser(l, ldapConn, *uid)
	} else if *addGroup == true{
		b.AddGroup(l, ldapConn, *cn, *ou, *memberIds)
	} else if *deleteGroup == true{
		b.DeleteGroup(l, ldapConn, *cn,*ou)
	}else if *addMembers == true{
		b.AddMembersToGroup(l, ldapConn, *cn, *ou, *memberIds)
	}else if *removeMembers == true{
		b.RemoveMembersFromGroup(l, ldapConn, *cn, *ou, *memberIds)
	}else {
		log.Println("Select a valid action flag.")
	}
}
