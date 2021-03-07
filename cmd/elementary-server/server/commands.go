package server

import "github.com/forensicanalysis/elementary/plugin"

func Commands(cp plugin.Provider) []*Command {
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
