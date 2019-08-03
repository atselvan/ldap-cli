package main

import (
	b "com/privatesquare/go/ldap-cli/backend"
	"com/privatesquare/go/ldap-cli/ldap"
	m "com/privatesquare/go/ldap-cli/model"
	"flag"
	"fmt"
	"log"
)

func main() {

	//action
	DeleteUserfromAllGroups := flag.Bool("DeleteUserfromAllGroups", false, "Deletes a user from all groups. Required parameter: uid")
	addUser := flag.Bool("addUser", false, "Add a user. Required parameters: uid,cn,sn mail, password")
	deleteUser := flag.Bool("deleteUser", false, "Deletes a user. Required parameter: uid")
	addGroup := flag.Bool("addGroup", false, "Add a Group. Required parameter: cn, ou. Optional parameter: memberIds")
	deleteGroup := flag.Bool("deleteGroup", false, "Delete a Group. Required Parameter: cn, ou")
	addMembers := flag.Bool("addMembers", false, "Add a member to a group. Required Parameter: cn, ou, memberIds")
	removeMembers := flag.Bool("removeMembers", false, "Remove a member to a group. Required Parameter: cn, ou, memberIds")
	removeAllExceptSome := flag.Bool("removeAllExceptSome", false, "Remove all members except the memberIds input from the group. Required Parameter: cn, ou, memberIds")
	listAllUsers := flag.Bool("listAllUsers", false, "Lists all user accounts")
	resetAllPersonalAccountsPassword := flag.Bool("resetAllPersonalAccountPasswords", false, "This task will reset the passwords of all personal user accounts.")

	//parameters
	connConfigFile := flag.String("connConfigFile", "./conf/ldap-connection.json", "Connection details of the LDAP server.")
	bindPassword := flag.String("bindPassword", "", "LDAP bind password.")
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
	} else if *addGroup == true {
		b.AddGroup(l, ldapConn, *cn, *ou, *memberIds)
	} else if *deleteGroup == true {
		b.DeleteGroup(l, ldapConn, *cn, *ou)
	} else if *addMembers == true {
		b.AddMembersToGroup(l, ldapConn, *cn, *ou, *memberIds)
	} else if *removeMembers == true {
		b.RemoveMembersFromGroup(l, ldapConn, *cn, *ou, *memberIds)
	} else if *removeAllExceptSome == true {
		b.RemoveAllMembersExceptSome(l, ldapConn, *cn, *ou, *memberIds)
	} else if *DeleteUserfromAllGroups == true {
		b.DeleteUserfromAllGroups(l, ldapConn, *uid)
	} else if *listAllUsers == true {
		b.ListAllUsersAccounts(l)
	} else if *resetAllPersonalAccountsPassword == true {
		b.ResetAllPersonalUserAccountsPassword(l, ldapConn)
	} else {
		log.Println("Select a valid action flag.")
	}
}
