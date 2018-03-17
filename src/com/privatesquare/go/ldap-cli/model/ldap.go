package model

type LDAPConn struct {
	Hostname     string `json:"hostname"`
	Port         string `json:"port"`
	BindUser     string `json:"bindUser"`
	BaseDN       string `json:"baseDN"`
	UserBaseDN   string `json:"userBaseDN"`
	GroupBaseDN  string `json:"groupBaseDN"`
	BindPassword string
}

type UserDetails struct {
	Uid          string
	Cn           string
	Sn           string
	Mail         string
	UserPassword string
}
