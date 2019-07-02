package util

import "github.com/kowshikRoy/cft/model"

var Mapping = map[string]model.Language{

	// Language specific configuration
	"c++": model.Language{
		Name: "c++",
		Compiler: map[string]string{
			"windows": "g++.exe",
			"darwin":  "g++",
			"linux":   "g++",
		},
		Extension:  ".cpp",
		BuildFlags: "-static -DONLINE_JUDGE -lm -s -x c++ -Wl,--stack=268435456 -O2 -std=c++11 -D__USE_MINGW_ANSI_STDIO=0 ",
		OutputFileExtension: map[string]string{
			"windows": ".exe",
			"darwin":  "",
			"linux":   "",
		},
	},

	// todo
}
