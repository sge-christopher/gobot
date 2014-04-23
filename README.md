gobot
=====
To install:  
```
go get github.com/sg-christopher/gobot
```

Usage:  
- Set env variable SPITFIRE_WORKSPACE to spitfire directory  

```
gobot h
```

##### NAME:  
gobot - handle the robot things

##### USAGE:
   gobot [global options] command [command options] [arguments...]

##### VERSION:
   0.0.1

##### COMMANDS:
   bundle, b	bundles each directory in the workspace if it contains a Gemfile
   pull, p	Pulls the branch that is checked out in each directory. Will ignore repos that have modified files
   heads, ch	Returns the current head + branch of each repo in the workspace
   help, h	Shows a list of commands or help for one command

##### GLOBAL OPTIONS:
   --version, -v	print the version
   --help, -h		show help
