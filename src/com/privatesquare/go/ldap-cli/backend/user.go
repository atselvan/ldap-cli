package backend

import (
	"com/privatesquare/go/ldap-cli/ldap"
	m "com/privatesquare/go/ldap-cli/model"
	"com/privatesquare/go/ldap-cli/utils"
	"fmt"
	"log"
	"regexp"
	"strings"
)

func UserExists(l *ldap.Conn, ldapConn m.LDAPConn, uid string) bool {
	if uid == "" {
		log.Fatal("uid is a required parameter for checking if a user exists")
	}
	var isExists bool
	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("uid=%s,%s", uid, ldapConn.UserBaseDN),
		ldap.ScopeSingleLevel, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=organizationalPerson))",
		[]string{"uid"},
		nil,
	)
	_, err := l.Search(searchRequest)
	if err != nil {
		isExists = false
	} else {
		isExists = true
	}
	return isExists
}

func AddUser(l *ldap.Conn, ldapConn m.LDAPConn, user m.UserDetails) {
	if user.Uid == "" || user.Cn == "" || user.Sn == "" || user.Mail == "" || user.UserPassword == "" {
		log.Fatal("User details are incomplete for adding a user. Required parameters: uid, cn, sn, mail, password")
	}
	if !UserExists(l, ldapConn, user.Uid) {
		a := ldap.NewAddRequest(fmt.Sprintf("uid=%s,%s", user.Uid, ldapConn.UserBaseDN))
		a.Attribute("objectClass", []string{"person", "organizationalPerson", "inetOrgPerson", "top", "userExtras"})
		a.Attribute("uid", []string{user.Uid})
		a.Attribute("cn", []string{strings.Title(user.Cn)})
		a.Attribute("sn", []string{strings.Title(user.Sn)})
		a.Attribute("displayName", []string{fmt.Sprintf("%s %s", strings.Title(user.Cn), strings.Title(user.Sn))})
		a.Attribute("mail", []string{user.Mail})
		a.Attribute("userPassword", []string{user.UserPassword})
		a.Attribute("status", []string{"Active"})
		err := l.Add(a)
		if err != nil {
			log.Println("User entry could not be added :", err)
		} else {
			log.Printf("User %s is added\n", user.Uid)
		}
	} else {
		log.Printf("User %s already exists\n", user.Uid)
	}
}

func DeleteUser(l *ldap.Conn, ldapConn m.LDAPConn, uid string) {
	if uid == "" {
		log.Fatal("uid is a required paramter for deleting a user")
	}
	if UserExists(l, ldapConn, uid) {
		d := ldap.NewDelRequest(fmt.Sprintf("uid=%s,%s", uid, ldapConn.UserBaseDN), nil)
		err := l.Del(d)
		if err != nil {
			fmt.Println("User could not be deleted :", err)
		} else {
			log.Printf("User %s is deleted\n", uid)
		}
	} else {
		log.Printf("User %s does not exist, hence could not be deleted", uid)
	}
}

func getAllUserAccounts(l *ldap.Conn) []m.UserDetails {
	searchRequest := ldap.NewSearchRequest(
		"ou=users,o=sccm",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=inetOrgPerson))",
		[]string{"uid", "cn", "sn", "mail", "userPassword", "status"},
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}
	var usersList []m.UserDetails

	for _, user := range sr.Entries {
		userDetails := m.UserDetails{
			Uid:          user.GetAttributeValue("uid"),
			Cn:           user.GetAttributeValue("cn"),
			Sn:           user.GetAttributeValue("sn"),
			Mail:         user.GetAttributeValue("mail"),
			UserPassword: user.GetAttributeValue("userPassword"),
			Status:       user.GetAttributeValue("status"),
		}
		usersList = append(usersList, userDetails)
	}

	return usersList
}

// Filters personal account and other technical accounts
// TODO get only builder accounts also
func FilterUserAccounts(l *ldap.Conn) ([]m.UserDetails, []m.UserDetails) {
	usersList := getAllUserAccounts(l)
	var personalAccounts []m.UserDetails
	var technicalAccounts []m.UserDetails

	personalAccountRegExp := "^[A-Z]{4}|[(A-Z)|(a-z)]{1}[0-9]{5}|[(A-Z)|(a-z)]{2}[0-9]{4}|[A-Z]{2}[0-9]{2}"

	r, _ := regexp.Compile(personalAccountRegExp)

	for _, user := range usersList {
		if r.MatchString(user.Uid) && !strings.Contains(user.Uid, "_") {
			personalAccounts = append(personalAccounts, user)
		} else {
			technicalAccounts = append(technicalAccounts, user)
		}
	}
	return personalAccounts, technicalAccounts
}

func ListAllUsersAccounts(l *ldap.Conn) {
	userAccounts := getAllUserAccounts(l)
	for _, ua := range userAccounts {
		fmt.Printf("UID:    %s\nCN:     %s\nSN:     %s\nMail:   %s\nStatus: %s\n\n", ua.Uid, ua.Cn, ua.Sn, ua.Mail, ua.Status)
	}
}

// TODO Add validation - user details should not be empty
func setUserPassword(l *ldap.Conn, ldapConn m.LDAPConn, user m.UserDetails) {

	passwordModifyRequest := ldap.NewPasswordModifyRequest(fmt.Sprintf("uid=%s,%s", user.Uid, ldapConn.UserBaseDN), "", "")
	passwordModifyResponse, err := l.PasswordModify(passwordModifyRequest)
	if err != nil {
		log.Fatalf("Password could not be changed: %s", err.Error())
	}

	_ = passwordModifyResponse.GeneratedPassword

	if err != nil {
		log.Println("Password could not be changed :", err)
	} else {
		log.Printf("User %s password is set\n", user.Uid)
	}
}

func ResetAllPersonalUserAccountsPassword(l *ldap.Conn, ldapConn m.LDAPConn) {
	personalAccounts, _ := FilterUserAccounts(l)
	for _, pa := range personalAccounts {
		setUserPassword(l, ldapConn, pa)
	}
}

func ModifyUser() {
	// TODO
}
