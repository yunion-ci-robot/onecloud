package ansible

import (
	"fmt"
	"strings"
)

// TODO
//
// quote parser

// ParseHostLine parse string representation of an ansible host
//
// Input should be in format "name key=value"
func ParseHostLine(s string) (host Host, err error) {
	s = strings.TrimSpace(s)
	parts := strings.Split(s, " ")
	if len(parts) == 0 {
		err = fmt.Errorf("Host name must not be empty")
		return
	}
	host.Name = parts[0]
	vars := map[string]string{}
	for _, kv := range parts[1:] {
		s := strings.SplitN(kv, "=", 2)
		if len(s) != 2 {
			err = fmt.Errorf("host var not in the form name=val: %s", kv)
			return
		}
		vars[s[0]] = s[1]
	}
	host.Vars = vars
	return
}

// ParseModuleLine parse string representation of an ansible module task
//
// The input should be in format "name key=value arg0 arg1".  The argN form is
// For module "command" and "shell"
func ParseModuleLine(s string) (mod Module, err error) {
	s = strings.TrimSpace(s)
	parts := strings.Split(s, " ")
	if len(parts) == 0 {
		err = fmt.Errorf("Module name must not be empty")
		return
	}
	mod.Name = parts[0]

	args := []string{}
	command := ""
	freeForm := false // command and shell module take free form arguments
	for _, part := range parts[1:] {
		if freeForm {
			command = command + " " + part
			continue
		}
		if strings.Contains(part, "=") {
			args = append(args, part)
		} else {
			freeForm = true
			command = command + " " + part
		}
	}
	if command != "" {
		args = append(args, command)
	}
	mod.Args = args
	return
}
