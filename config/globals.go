package config // COMMANDS

var Commands = struct {
	Login         string
	Logout        string
	SysDisconnect string
	Subscribe     string
}{
	Login:         "/login",
	Logout:        "/logout",
	SysDisconnect: "/sys-disconnect",
	Subscribe:     "/subscribe",
}
