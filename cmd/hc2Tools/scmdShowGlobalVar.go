package main

import (
	"encoding/json"
	"fmt"
)

type showGlobalVar struct {
	PrettyPrint bool     `type:"flag" name:"prettyprint" help:"PrettyPrint the JSON values"`
	VarNames    []string `type:"arg" name:"varName" help:"Print value of a global variable. All if none given"`
}

const showGlobalVarUsage = "Print global variable value, all if none given"

func (cmd *showGlobalVar) Run() {
	var selectedVars = getGlobalVariables(cmd.VarNames)

	for _, variable := range selectedVars {
		var s []byte
		if cmd.PrettyPrint {
			s, _ = json.MarshalIndent(variable, "", "    ")
		} else {
			s, _ = json.Marshal(variable)
		}
		fmt.Printf("%s\n", s)
	}

}
