package backend

import (
	"com/privatesquare/go/ldap-cli/ldap"
	m "com/privatesquare/go/ldap-cli/model"
	"fmt"
	"log"
	"os"
	"strings"
)

func GroupExists(l *ldap.Conn, ldapConn m.LDAPConn, cn, ou string) bool {
	if cn == "" || ou == "" {
		log.Fatal("cn and ou are required paramters for checking if a group exists")
	}
	var isExists bool
	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("cn=%s,ou=%s,%s", cn, ou, ldapConn.GroupBaseDN),
		ldap.ScopeSingleLevel, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=groupOfUniqueNames))",
		[]string{"cn"},
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

func GetGroups(l *ldap.Conn, ldapConn m.LDAPConn, ou string) []string {
	if ou == "" {
		log.Fatal("ou is a required paramter for getting a list of groups")
	}
	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("ou=%s,%s", ou, ldapConn.GroupBaseDN),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=groupOfUniqueNames))",
		[]string{"cn"},
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}
	var groupsList []string
	for _, entry := range sr.Entries {
		groupsList = append(groupsList, entry.GetAttributeValue("cn"))
	}
	return groupsList
}

func AddGroup(l *ldap.Conn, ldapConn m.LDAPConn, cn, ou, memberId string) {
	if cn == "" || ou == "" {
		log.Fatal("cn and ou are required paramters for adding a group")
	}
	if memberId == "" {
		memberId = "NO_SUCH_USER"
	} else {
		if !UserExists(l, ldapConn, memberId) {
			log.Printf("Group %s could not be created because the user %s does not exist\n", cn, memberId)
			os.Exit(1)
		}
	}
	if !GroupExists(l, ldapConn, cn, ou) {
		a := ldap.NewAddRequest(fmt.Sprintf("cn=%s,ou=%s,%s", cn, ou, ldapConn.GroupBaseDN))
		a.Attribute("objectClass", []string{"groupOfUniqueNames", "top"})
		a.Attribute("cn", []string{cn})
		a.Attribute("uniqueMember", []string{fmt.Sprintf("uid=%s,%s", memberId, ldapConn.UserBaseDN)})
		err := l.Add(a)
		if err != nil {
			log.Println("Group entry could not be added :", err)
		} else {
			log.Printf("Group %s is created\n", cn)
		}
	} else {
		log.Printf("Group %s already exists\n", cn)
	}
}

func DeleteGroup(l *ldap.Conn, ldapConn m.LDAPConn, cn, ou string) {
	if cn == "" || ou == "" {
		log.Fatal("cn and ou are required paramters for deleting a group")
	}
	if GroupExists(l, ldapConn, cn, ou) {
		d := ldap.NewDelRequest(fmt.Sprintf("cn=%s,ou=%s,%s", cn, ou, ldapConn.GroupBaseDN), nil)
		err := l.Del(d)
		if err != nil {
			log.Println("Group entry could not be deleted :", err)
		} else {
			log.Printf("Group %s is deleted\n", cn)
		}
	} else {
		log.Printf("Group %s does not exists, hence not deleting the group\n", cn)
	}
}

func GetGroupMembers(l *ldap.Conn, ldapConn m.LDAPConn, cn, ou string) []string {
	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("cn=%s,ou=%s,%s", cn, ou, ldapConn.GroupBaseDN),
		ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=groupOfUniqueNames))",
		[]string{"uniqueMember"},
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}
	var memberIdList []string
	memberDnList := sr.Entries[0].GetAttributeValues("uniqueMember")
	for _, memberDN := range memberDnList {
		memberDN = strings.Replace(memberDN, "uid=", "", -1)
		memberDN = strings.Replace(memberDN, fmt.Sprintf(",%s", ldapConn.UserBaseDN), "", -1)
		memberIdList = append(memberIdList, memberDN)
	}
	return memberIdList
}

func addMemberToGroup(l *ldap.Conn, ldapConn m.LDAPConn, cn, ou, memberId string) {
	if cn == "" || ou == "" || memberId == "" {
		log.Fatal("cn, ou and memberIds are required paramters for adding a member to a group")
	}
	if !GroupExists(l, ldapConn, cn, ou) {
		log.Printf("Group %s does not exist, cn or ou is incorrect\n", cn)
		os.Exit(1)
	}
	if !UserExists(l, ldapConn, memberId) {
		log.Printf("User %s does not exist, hence the user could not be added to the group %s\n", memberId, cn)
		os.Exit(1)
	}
	memberExists := false
	membersIdList := GetGroupMembers(l, ldapConn, cn, ou)
	for _, uniqueMemberId := range membersIdList {
		if uniqueMemberId == memberId {
			memberExists = true
		}
	}
	if !memberExists {
		modify := ldap.NewModifyRequest(fmt.Sprintf("cn=%s,ou=%s,%s", cn, ou, ldapConn.GroupBaseDN))
		modify.Add("uniqueMember", []string{fmt.Sprintf("uid=%s,%s", memberId, ldapConn.UserBaseDN)})
		err := l.Modify(modify)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("User %s is added to the group %s\n", memberId, cn)
	} else {
		log.Printf("User %s is already a member of the group %s\n", memberId, cn)
	}
}

func AddMembersToGroup(l *ldap.Conn, ldapConn m.LDAPConn, cn, ou, memberIds string) {
	if cn == "" || ou == "" || memberIds == "" {
		log.Fatal("cn, ou and memberIds are required paramters for adding a member to a group")
	}
	if !GroupExists(l, ldapConn, cn, ou) {
		log.Printf("Group %s does not exist, cn or ou is incorrect\n", cn)
		os.Exit(1)
	}
	membersList := strings.Split(memberIds, ",")
	for _, member := range membersList {
		addMemberToGroup(l, ldapConn, cn, ou, member)
	}
}

func removeMemberFromGroup(l *ldap.Conn, ldapConn m.LDAPConn, cn, ou, memberId string) {
	if cn == "" || ou == "" || memberId == "" {
		log.Fatal("cn, ou and memberIds are required paramters for removing a member from a group")
	}
	if !GroupExists(l, ldapConn, cn, ou) {
		log.Printf("Group %s does not exist, cn or ou is incorrect\n", cn)
		os.Exit(1)
	}
	memberExists := false
	membersIdList := GetGroupMembers(l, ldapConn, cn, ou)
	for _, uniqueMemberId := range membersIdList {
		if uniqueMemberId == memberId {
			memberExists = true
		}
	}
	if len(membersIdList) == 1 {
		log.Printf("A group must contain atleast one member, hence cannot remove the last member of the group\n")
		os.Exit(1)
	}
	if memberExists {
		modify := ldap.NewModifyRequest(fmt.Sprintf("cn=%s,ou=%s,%s", cn, ou, ldapConn.GroupBaseDN))
		modify.Delete("uniqueMember", []string{fmt.Sprintf("uid=%s,%s", memberId, ldapConn.UserBaseDN)})
		err := l.Modify(modify)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("User %s is removed from the group %s\n", memberId, cn)
	} else {
		log.Printf("User %s is not a member of the group %s\n", memberId, cn)
	}
}

func RemoveMembersFromGroup(l *ldap.Conn, ldapConn m.LDAPConn, cn, ou, memberIds string) {
	if cn == "" || ou == "" || memberIds == "" {
		log.Fatal("cn, ou and memberIds are required paramters for adding a member to a group")
	}
	if !GroupExists(l, ldapConn, cn, ou) {
		log.Printf("Group %s does not exist, cn or ou is incorrect\n", cn)
		os.Exit(1)
	}
	membersList := strings.Split(memberIds, ",")
	for _, member := range membersList {
		removeMemberFromGroup(l, ldapConn, cn, ou, member)
	}
}

func RemoveAllMembersExceptSome(l *ldap.Conn, ldapConn m.LDAPConn, cn, ou, memberIds string) {
	if cn == "" || ou == "" || memberIds == "" {
		log.Fatal("cn, ou and memberIds are required paramters for removing members from a group")
	}
	if !GroupExists(l, ldapConn, cn, ou) {
		log.Printf("Group %s does not exist, cn or ou is incorrect\n", cn)
		os.Exit(1)
	}
	var removeMembersList []string
	existingMembersList := GetGroupMembers(l, ldapConn, cn, ou)
	retainMembersList := strings.Split(memberIds, ",")
	for _, member := range retainMembersList {
		var loopList []string
		for _, existingMember := range existingMembersList {
			if existingMember != member {
				loopList = append(loopList, existingMember)
			}
		}
		existingMembersList = loopList
	}
	removeMembersList = existingMembersList
	for _, member := range removeMembersList {
		removeMemberFromGroup(l, ldapConn, cn, ou, member)
	}
}
func DeleteUserfromAllGroups(l *ldap.Conn, ldapConn m.LDAPConn, uid string) {
	ToolsList := [...]string{"bitbucket", "hudson", "nexus", "nexusIQ", "sccm", "sonar", "Fortify", "subversion"}
	//ToolsList := [...]string{"Fortify","subversion"}
	for _, tools := range ToolsList {
		list := GetGroups(l, ldapConn, tools)
		for _, Group := range list {
			membersIdList := GetGroupMembers(l, ldapConn, Group, tools)
			for _, uniqueMemberId := range membersIdList {
				if uniqueMemberId == uid {
					removeMemberFromGroup(l, ldapConn, Group, tools, uid)
				}
			}
		}
	}
}
