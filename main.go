package main

import (
	"Nie-Mand/dfile/pkg"
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
		Run("npm ci").Copy(".", ".").
		Run("npm run build")

	
	dockerfile.NextStage()

	dockerfile.
		From("nginx").
		Copy("/dist", "/usr/share/nginx/html").
		Cmd("nginx -g daemon off;").
		Expose(80)

	dockerfile.Save()
}
