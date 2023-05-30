package main

import (
	"github.com/Nie-Mand/dfile/pkg"
)

func main() {
	dockerfile := pkg.Dockerfile{}
	dockerfile.Init()

	dockerfile.
		From("node").
		ImageVersion("16").
		ImageAlias("build").
		WorkDir("/app").
		Copy("package.json", ".").
		Copy("package-lock.json", ".").
		Run("npm ci").
		Copy(".", ".").
		BuildEnvs("NODE_ENV", "production").
		Envs("PORT", "3000").
		Envs("API_URL", "").
		Run("npm run build")

	dockerfile.NextStage()

	dockerfile.
		From("nginx").
		Copy("/dist", "/usr/share/nginx/html").
		Cmd("nginx", "-g", "daemon off;").
		Expose(80)

	err := dockerfile.Save()
	if err != nil {
		panic(err)
	} else {
		println("Dockerfile generated successfully!")
	}
}
