package cmd

import (
	"log"
	"os/exec"
	"reflect"
	"strings"

	"github.com/PolarGeospatialCenter/warewulf-sync/pkg/warewulf"
	"github.com/spf13/cobra"
)

type SyncableObject interface {
	IdString() string
	DeleteCmd() [][]string
	UpdateCmd() [][]string
	NewCmd() [][]string
	Equals(interface{}) (bool, string)
}

func MakeSyncableMap(m interface{}) map[string]SyncableObject {
	result := make(map[string]SyncableObject)
	if reflect.ValueOf(m).Kind() == reflect.Map {
		for _, key := range reflect.ValueOf(m).MapKeys() {
			v := reflect.ValueOf(m).MapIndex(key)
			if v.Type().AssignableTo(reflect.TypeOf((*SyncableObject)(nil)).Elem()) {
				result[key.Interface().(string)] = v.Interface().(SyncableObject)
			}
		}
	}
	return result
}

func BuildSyncCommands(existing, desired map[string]SyncableObject) [][]string {
	commands := make([][]string, 0)
	for _, existingObject := range existing {
		if desiredObject, ok := desired[existingObject.IdString()]; !ok {
			commands = append(commands, existingObject.DeleteCmd()...)
		} else {
			if equal, reason := existingObject.Equals(desiredObject); !equal {
				log.Printf("Object %s doesn't match desired state: %s", desiredObject.IdString(), reason)
				commands = append(commands, desiredObject.UpdateCmd()...)
			}
		}
	}

	for _, desiredObject := range desired {
		if _, ok := existing[desiredObject.IdString()]; !ok {
			commands = append(commands, desiredObject.NewCmd()...)
			commands = append(commands, desiredObject.UpdateCmd()...)
		}
	}
	return commands
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		// load desired state from yaml
		db, err := warewulf.LoadYaml(cfg.GetString("config_path"))
		if err != nil {
			log.Fatalf("Unable to load files data: %v", err)
		}

		log.Print(db)

		// load warewulf
		wwdb, err := warewulf.LoadWwshDB()
		if err != nil {
			log.Fatalf("Unable to load warewulf database: %v", err)
		}
		log.Print(wwdb)

		syncCommands := make([][]string, 0)
		syncCommands = append(syncCommands, BuildSyncCommands(MakeSyncableMap(wwdb.Files), MakeSyncableMap(db.Files))...)
		syncCommands = append(syncCommands, BuildSyncCommands(MakeSyncableMap(wwdb.Bootstraps), MakeSyncableMap(db.Bootstraps))...)
		syncCommands = append(syncCommands, BuildSyncCommands(MakeSyncableMap(wwdb.Vnfs), MakeSyncableMap(db.Vnfs))...)
		syncCommands = append(syncCommands, BuildSyncCommands(MakeSyncableMap(wwdb.Nodes), MakeSyncableMap(db.Nodes))...)
		syncCommands = append(syncCommands, []string{"wwsh", "pxe", "-v", "--nodhcp"})

		for _, cmd := range syncCommands {
			log.Print(cmd)
			c := exec.Command(cmd[0], cmd[1:]...)

			stdErrOut, err := c.CombinedOutput()
			if err != nil {
				log.Fatalf("Error executing '%s': %v", strings.Join(cmd, " "), err)
			}
			log.Printf("Result: %s", stdErrOut)
		}
		// load existing state from warewulf
		// files, err := warewulf.ListFiles()
		// compare
		// fix
	},
}
