package pkg

import (
	"os"
	"strings"
)

type Dockerfile struct {
	name string
	stages []Stage
	currentStage int
}

func (df *Dockerfile) Init() {
	df.stages = append(df.stages, Stage{})
	df.currentStage = 0
}

func (df *Dockerfile) NextStage() {
	df.stages = append(df.stages, Stage{})
	df.currentStage++
}

func (df *Dockerfile) GetFilename() string {
	if df.name != "" {
		return df.name + ".Dockerfile"
	}

	return "Dockerfile"
}

func (df *Dockerfile) From(image string) (*Dockerfile) {
	df.stages[df.currentStage].From(image)
	return df
}

func (df *Dockerfile) ImageAlias(as string) (*Dockerfile) {
	df.stages[df.currentStage].ImageAlias(as)
	return df
}

func (df *Dockerfile) ImageVersion(version string) (*Dockerfile) {
	df.stages[df.currentStage].ImageVersion(version)
	return df
}

func (df *Dockerfile) WorkDir(workDir string) (*Dockerfile) {
	df.stages[df.currentStage].WorkDir(workDir)
	return df
}

func (df *Dockerfile) Cmd(cmd string) (*Dockerfile) {
	df.stages[df.currentStage].Cmd(cmd)
	return df
}

func (df *Dockerfile) Expose(port int) (*Dockerfile) {
	df.stages[df.currentStage].Expose(port)
	return df
}

func (df *Dockerfile) Run(cmd string) (*Dockerfile) {
	df.stages[df.currentStage].Run(cmd)
	return df
}

func (df *Dockerfile) Copy(from string, to string) (*Dockerfile) {
	df.stages[df.currentStage].Copy(from, to)
	return df
}

func (df *Dockerfile) BuildEnvs(k string, v string) (*Dockerfile) {
	df.stages[df.currentStage].SetBuildEnv(k, v)
	return df
}

func (df *Dockerfile) Envs(k string, v string) (*Dockerfile) {
	df.stages[df.currentStage].SetEnv(k, v)
	return df
}

func (df *Dockerfile) GetCode() (string, error) {
	stages := make([]string, 0)
	for _, stage := range df.stages {
		code, err := stage.GetCode()
		if err != nil {
			return "", err
		}
		stages = append(stages, code)
	}

	return strings.Join(stages, "\n\n"), nil
}

func (df *Dockerfile) Save() error {
	f, err := os.Create(df.GetFilename())
	if err != nil {
		return err
	}

	defer f.Close()

	content, err := df.GetCode()
	if err != nil {
		return err
	}
	
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	
	return nil
}


