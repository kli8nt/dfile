package pkg

import (
	"fmt"
	"strings"
)

//  COPY, ARG, ENV, RUN[]

type Stage struct {
	from struct {
		as        string
		imageName string
		version   string
	}
	run       []string
	copy      []Pair
	expose    []int
	workDir   string
	cmd       []string
	commands  []TrackedCommand
	buildEnvs []Pair
	envs      []Pair
}

func (stage *Stage) From(image string) {
	stage.from.as = image
	stage.from.imageName = image
	stage.from.version = "latest"
}

func (stage *Stage) ImageAlias(as string) {
	stage.from.as = as
}

func (stage *Stage) ImageVersion(version string) {
	stage.from.version = version
}

func (stage *Stage) WorkDir(workDir string) {
	stage.workDir = workDir
}

func (stage *Stage) Cmd(cmd []string) {
	stage.cmd = cmd
}

func (stage *Stage) Expose(port int) {
	stage.expose = append(stage.expose, port)
}

func (stage *Stage) Run(cmd string) {
	stage.run = append(stage.run, cmd)
	stage.commands = append(stage.commands, TrackedCommand{command: "run", idx: len(stage.run) - 1})
}

func (stage *Stage) Copy(from string, to string) {
	stage.copy = append(stage.copy, Pair{first: from, second: to})
	stage.commands = append(stage.commands, TrackedCommand{command: "copy", idx: len(stage.copy) - 1})
}

func (stage *Stage) SetBuildEnv(key string, value string) {
	stage.buildEnvs = append(stage.buildEnvs, Pair{first: key, second: value})
}

func (stage *Stage) SetEnv(key string, value string) {
	stage.envs = append(stage.envs, Pair{first: key, second: value})
}

func (stage *Stage) GetCode() (string, error) {
	if stage.from.imageName == "" {
		return "", fmt.Errorf("You did not specify a Base Image")
	}

	statements := Statements{}
	statements.AddStatement("FROM", stage.from.imageName+":"+stage.from.version, "AS", stage.from.as)

	if stage.workDir != "" {
		statements.AddStatement("WORKDIR", stage.workDir)
	}

	for _, env := range stage.buildEnvs {
		if env.second == "" {
			statements.AddStatement("ARG", env.first)
		} else {
			statements.AddStatement("ARG", env.first+"="+env.second)
		}
	}

	for _, env := range stage.envs {
		if env.second == "" {
			statements.AddStatement("ARG", env.first)
			
		} else {
			statements.AddStatement("ARG", env.first+"="+env.second)
		}
		statements.AddStatement("ENV", env.first+"=$"+env.first)
	}

	for _, command := range stage.commands {
		if command.command == "run" && len(stage.run) > command.idx {
			instruction := stage.run[command.idx]
			statements.AddStatement("RUN", instruction)
		} else if command.command == "copy" && len(stage.copy) > command.idx {
			from := stage.copy[command.idx].first
			to := stage.copy[command.idx].second
			statements.AddStatement("COPY", from, to)
		} else {
			return "", fmt.Errorf("Unrecognized Command")
		}
	}

	for _, port := range stage.expose {
		statements.AddStatement("EXPOSE", fmt.Sprintf("%d", port))
	}

	if len(stage.cmd) != 1 {
		splittedCmd := map2(stage.cmd, wrapInQuotes)
		cmd := strings.Join(splittedCmd, ", ")
		statements.AddStatement("CMD", "["+cmd+"]")
	}

	return strings.Join(statements.instructions, "\n\n"), nil
}
