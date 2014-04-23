package main

import (
	"bytes"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/mgutz/ansi"
	"os"
	"os/exec"
	"strings"
)

var (
	success string = ansi.ColorCode("green")
	fail    string = ansi.ColorCode("red")
	reset   string = ansi.ColorCode("reset")
	info    string = ansi.ColorCode("cyan")
	callout string = ansi.ColorCode("white+b")
	yellow  string = ansi.ColorCode("yellow")
)

// Function to get the directories in the given workspace ws
func getDirs(ws string) ([]string, error) {
	var dirs []string
	f, err := os.Open(ws)
	defer f.Close()

	if err != nil {
		return nil, err
	}

	files, err := f.Readdir(0)

	for _, v := range files {
		if v.IsDir() {
			dirs = append(dirs, v.Name())
		}
	}

	return dirs, nil
}

func gitStatus() (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("git", "status", "-s")

	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func gitPull() error {
	var stderr bytes.Buffer
	cmd := exec.Command("git", "pull", "--rebase")
	err := cmd.Run()
	cmd.Stderr = &stderr

	if err != nil {
		fmt.Println(stderr.String())
		return err
	}

	return nil
}

func gitCurrentHead() (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("git", "log", "HEAD~1..", "--oneline")
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		return "", err
	}

	return strings.Split(out.String(), "\n")[0], nil
}

func gitCurrentBranch() (string, error) {
	var out bytes.Buffer
	// Outputs only the current branch
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

func gitCommitsSince(commit string) (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("git", "log", commit+"..", "--oneline")
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func main() {
	errHandler := func(err error, shouldBreak bool) {
		if err != nil {
			fmt.Println(fail + err.Error() + reset)
			if shouldBreak {
				return
			}
		}
	}

	ws := os.Getenv("SPITFIRE_WORKSPACE")
	if len(ws) < 1 {
		panic("SPITFIRE_WORKSPACE is not set")
	}

	app := cli.NewApp()
	app.Name = "gobot"
	app.Usage = "handle the robot things"

	app.Commands = []cli.Command{
		{
			Name:      "bundle",
			ShortName: "b",
			Usage:     "bundles each directory in the workspace if it contains a Gemfile",
			Action: func(c *cli.Context) {
				// Gets a []string of directory names from our workspace
				dirs, err := getDirs(ws)
				errHandler(err, true)

				// Iterate over the directory names and bundle if there is a Gemfile
				for _, d := range dirs {
					curr_dir := ws + "/" + d

					if _, err := os.Stat(curr_dir + "/Gemfile"); err == nil {
						if err := os.Chdir(curr_dir); err == nil {
							out, err := exec.Command("bundle").Output()
							if err != nil {
								fmt.Printf("%s%s failed to bundle:\n%s%s\n", fail, d, string(out), reset)
							} else {
								println(success + "Bundled in " + d + reset)
							}
						}
					}
				}
			},
		},
		{
			Name:      "pull",
			ShortName: "p",
			Usage:     "Pulls the branch that is checked out in each directory. Will ignore repos that have modified files",
			Action: func(c *cli.Context) {
				dirs, err := getDirs(ws)
				errHandler(err, true)

				for _, d := range dirs {
					curr_dir := ws + "/" + d

					// Change to the working directory
					os.Chdir(curr_dir)
					m, err := gitStatus()

					if err != nil {
						fmt.Println(err)
					} else if len(m) > 0 {
						fmt.Printf("There are are modified files in %s.\n%s\n", d, m)
					} else {
						fmt.Printf("%sRebase...%s\n%s", callout, d, reset)
						currentHead, err := gitCurrentHead()
						errHandler(err, true)

						err = gitPull()
						errHandler(err, false)

						changes, err := gitCommitsSince(strings.Split(currentHead, " ")[0])
						errHandler(err, false)

						fmt.Printf(changes)
						fmt.Println(success + "OK" + reset)
					}
				}
				return
			},
		},
		{
			Name:      "head",
			ShortName: "ch",
			Action: func(*cli.Context) {
				dirs, err := getDirs(ws)
				errHandler(err, true)

				for _, d := range dirs {
					curr_dir := ws + "/" + d
					os.Chdir(curr_dir)

					msg, err := gitCurrentHead()
					errHandler(err, false)

					branch, err := gitCurrentBranch()
					errHandler(err, false)

					repo := callout + d + reset
					msgArr := strings.Split(msg, " ")
					commit := yellow + msgArr[0] + reset
					message := info + strings.Join(msgArr[1:], " ") + reset
					fmt.Println(repo + "(" + branch + ")\n" + commit + " " + message)
				}
				return
			},
		},
	}

	app.Run(os.Args)
}
