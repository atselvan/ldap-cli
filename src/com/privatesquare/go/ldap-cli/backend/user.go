package backend

import (
	"com/privatesquare/go/ldap-cli/ldap"
	m "com/privatesquare/go/ldap-cli/model"
	"fmt"
	"log"
	"strings"
)

func UserExists(l *ldap.Conn, ldapConn m.LDAPConn, uid string) bool{
	if uid == "" {
		log.Fatal("uid is a required paramter for checking if a user exists")
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

func AddUser(l *ldap.Conn, ldapConn m.LDAPConn, user m.UserDetails){
	if user.Uid == "" || user.Cn == "" || user.Sn == "" || user.Mail == "" || user.UserPassword == "" {
		log.Fatal("User details are incomplete for adding a user. Required parameters: uid, cn, sn, mail, password")
	}
	if !UserExists(l, ldapConn, user.Uid){
		a := ldap.NewAddRequest(fmt.Sprintf("uid=%s,%s", user.Uid, ldapConn.UserBaseDN))
		a.Attribute("objectClass" ,[]string{"person", "organizationalPerson", "inetOrgPerson", "top"})
		a.Attribute("uid" ,[]string{user.Uid})
		a.Attribute("cn" ,[]string{strings.Title(user.Cn)})
		a.Attribute("sn" ,[]string{strings.Title(user.Sn)})
		a.Attribute("displayName" ,[]string{fmt.Sprintf("%s %s", strings.Title(user.Cn), strings.Title(user.Sn))})
		a.Attribute("mail" ,[]string{user.Mail})
		a.Attribute("userPassword" ,[]string{user.UserPassword})
		err := l.Add(a)
		if err != nil {
			log.Println("User entry could not be added :",err)
		} else {
			log.Printf("User %s is added\n", user.Uid)
		}
	}else{
		log.Printf("User %s already exists\n", user.Uid)
	}
}

func DeleteUser(l *ldap.Conn, ldapConn m.LDAPConn, uid string){
	if uid == "" {
		log.Fatal("uid is a required paramter for deleting a user")
	}
	if UserExists(l, ldapConn, uid){
		d := ldap.NewDelRequest(fmt.Sprintf("uid=%s,%s", uid, ldapConn.UserBaseDN), nil)
		err := l.Del(d)
		if err != nil {
			fmt.Println("User could not be deleted :",err)
		} else {
			log.Printf("User %s is deleted\n", uid)
		}
	}else{
		log.Printf("User %s does not exist, hence could not be deleted", uid)
	}
}

func ModifyUser(){
	// TODO
}