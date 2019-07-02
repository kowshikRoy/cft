package cmd

import (
	"bytes"
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	var yamlExample = []byte(`
	user: tourist
	author: Gennady Korotkevich
	language: c++
	workdir: /Users/reproy/Desktop/code/codeforces
	margin: 1e-6
	templates:
		c++: /Users/reproy/Desktop/code/templates/cpptemplate.c++
		python: /Users/reproy/Desktop/code/templates/pytemplate.py
	buildConfig:
		c++:
			compiler:
				windows: g++.exe
				darwin: g++
				linux: g++
			extension: .cpp
			buildFlags: -static -DONLINE_JUDGE -lm -s -x c++ -Wl,--stack=268435456 -O2 -std=c++11 -D__USE_MINGW_ANSI_STDIO=0 
			outputFileExtension:
				windows: .exe
				darwin:
				linux:
		
	`)
	viper.ReadConfig(bytes.NewBuffer(yamlExample))
	fmt.Println(viper.AllSettings())

}
