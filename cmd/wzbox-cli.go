package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/infra-whizz/wzbox"
)

func app(c *cli.Context) error {
	files := strings.Split(c.String("files"), ",")

	packageName := c.String("package")
	if packageName == "" {
		packageName = "main"
	}

	outName := c.String("out")
	if outName == "" {
		outName = strings.Split(files[0], ".")[0] + ".go"
	}

	structName := c.String("struct")
	if structName == "" {
		structName = strings.Title(strings.Split(files[0], ".")[0])
	}

	gen := wzbox.NewWzBox().
		SetOutputFilename(outName).
		SetCompression(c.Bool("compress")).
		SetStructName(structName).
		SetPackageName(packageName)
	for _, fname := range files {
		gen.AddFile(fname)
	}
	out, err := gen.Generate()
	if err != nil {
		return err
	}

	fmt.Printf(out)

	return nil
}

func main() {
	appname := "wzbox"
	app := &cli.App{
		Version: "0.1",
		Name:    appname,
		Usage:   "Utility to generate Go sources with the embedded files and use them in Go projects",
		Action:  app,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Aliases:  []string{"f"},
				Name:     "files",
				Usage:    "Comma-separated files",
				Required: true,
			},
			&cli.StringFlag{
				Aliases: []string{"p"},
				Name:    "package",
				Usage:   "Name of the package (default: main)",
			},
			&cli.StringFlag{
				Aliases: []string{"s"},
				Name:    "struct",
				Usage:   "Name of the struct (default: first static filename, title-case)",
			},
			&cli.BoolFlag{
				Aliases: []string{"c"},
				Name:    "compress",
				Usage:   "Compress content (ZIP). NOTE: This might degrade performance of your project.",
			},
			&cli.StringFlag{
				Aliases: []string{"o"},
				Name:    "out",
				Usage:   "Output filename (default: first static filename lower-case with .go extension)",
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}

}
