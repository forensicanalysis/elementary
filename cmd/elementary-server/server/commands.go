package server

import (
	"github.com/forensicanalysis/elementary/daggy"
)

func Commands(cp daggy.CommandProvider) []*Command {
	return []*Command{
		ListTables(),
		SelectItems(),
		LoadFile(),
		ListTree(),
		ListTasks(cp),
		Files(),
		Logs(),
		ErrorsCommand(),
		Label(),
		Labels(),
		Query(),
		RunTask(cp),
	}
}
