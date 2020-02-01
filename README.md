# cft
A codeforces toolkit

Put the configuration file in $HOME/.cft.yaml

A sample Configuration:
```
user: tourist
author: Gennady Korotkevich
workdir: /Users/reponroy/Desktop/code/codeforces
margin: 1e-6
templates:
  c++: /Users/reponroy/Desktop/code/templates/cpptemplate.cpp
  python: /Users/reponroy/Desktop/code/templates/pytemplate.py
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
```
