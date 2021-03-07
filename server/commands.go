package server

func Commands() []*Command {
	return []*Command{
		ListTables(),
		SelectItems(),
		LoadFile(),
		ListTree(),
		ListTasks(),
		Files(),
		Logs(),
		ErrorsCommand(),
		Label(),
		Labels(),
		Query(),
		RunTask(),
	}
}
