package server

import "github.com/forensicanalysis/elementary/pluginlib"

func Commands(cp pluginlib.Provider) []*Command {
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
