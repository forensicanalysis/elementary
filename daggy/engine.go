package daggy

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/hashicorp/logutils"
	"github.com/hashicorp/terraform/dag"
	"github.com/hashicorp/terraform/tfdiags"
	"github.com/spf13/cobra"
)

type Engine struct {
	commands map[string]*cobra.Command
	mux      sync.Mutex
}

func New(commands []*cobra.Command) *Engine {
	setupLogging()
	engine := Engine{commands: map[string]*cobra.Command{}}
	for _, command := range commands {
		engine.commands[command.Name()] = command
	}
	return &engine
}

// Sort creates a direct acyclic graph of tasks.
func (e *Engine) Run(workflow *Workflow, storeDir string) error {
	// Create the dag
	graph := dag.AcyclicGraph{}
	tasks := map[uuid.UUID]Task{}
	outputToTaskIDs := map[string][]uuid.UUID{}
	var unavailableCommands []string

	for _, task := range workflow.Tasks {
		task.ID = uuid.New()
		graph.Add(task.ID)
		tasks[task.ID] = task

		cmd, ok := e.commands[task.Command]
		if !ok {
			unavailableCommands = append(unavailableCommands, task.Command)
			continue
		}

		if outputs, ok := cmd.Annotations["output"]; ok {
			for _, output := range strings.Split(outputs, ",") {
				outputToTaskIDs[output] = append(outputToTaskIDs[output], task.ID)
			}
		}
	}

	// Add edges / requirements
	for _, task := range workflow.Tasks {
		cmd, ok := e.commands[task.Command]
		if !ok {
			continue
		}

		if requires, ok := cmd.Annotations["requires"]; ok {
			for _, requirement := range strings.Split(requires, ",") {
				for _, output := range outputToTaskIDs[requirement] {
					graph.Connect(dag.BasicEdge(task.ID, output))
				}
			}
		}
	}

	w := &dag.Walker{Callback: func(v dag.Vertex) tfdiags.Diagnostics {
		task := tasks[v.(uuid.UUID)]

		if err := e.RunTask(task, storeDir); err != nil {
			return tfdiags.Diagnostics{tfdiags.Sourceless(tfdiags.Error, fmt.Sprint(v.(uuid.UUID)), err.Error())}
		}
		return nil
	}}
	w.Update(&graph)

	dagErr := w.Wait().Err()
	switch {
	case dagErr != nil && unavailableCommands != nil:
		return fmt.Errorf("unavailable commands: %v, run error: %w", unavailableCommands, dagErr)
	case dagErr != nil:
		return dagErr
	case unavailableCommands != nil:
		return fmt.Errorf("unavailable commands: %v", unavailableCommands)
	}
	return nil
}

func (e *Engine) RunTask(task Task, storeDir string) error {
	e.mux.Lock() // serialize tasks
	defer e.mux.Unlock()
	command, ok := e.commands[task.Command]
	if !ok {
		return errors.New("command not found")
	}

	var args []string
	for flag, value := range task.Arguments {
		args = append(args, toCmdline(flag, value)...)
	}
	args = append(args, storeDir)

	err := command.ParseFlags(args)
	if err != nil {
		return err
	}

	// command.SetArgs(args)
	if command.RunE == nil {
		return fmt.Errorf("command %s cannot be run", command.Name())
	}
	return command.RunE(command, command.Flags().Args())
}

func toCmdline(name string, i interface{}) []string {
	switch reflect.TypeOf(i).Kind() {
	case reflect.Slice:
		var s []string
		v := reflect.ValueOf(i)
		for i := 0; i < v.Len(); i++ {
			s = append(s, "--"+name, toCmdline2(v.Index(i)))
		}
		return s
	default:
		return []string{"--" + name, fmt.Sprint(i)}
	}
}

func toCmdline2(v reflect.Value) string {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Slice:
		var parts []string
		for i := 0; i < v.Len(); i++ {
			parts = append(parts, toCmdline2(v.Index(i)))
		}
		sort.Strings(parts)
		return strings.Join(parts, ",")
	case reflect.Map:
		var parts []string
		for _, k := range v.MapKeys() {
			i := v.MapIndex(k)
			parts = append(parts, fmt.Sprintf("%s=%s", k, i))
		}
		sort.Strings(parts)
		return strings.Join(parts, ",")
	default:
		return fmt.Sprint(v.Interface())
	}
}

func setupLogging() {
	// disable logging in github.com/hashicorp/terraform/dag
	log.SetOutput(&logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"TRACE", "OTHER"},
		MinLevel: "OTHER",
		Writer:   log.Writer(),
	})
}
