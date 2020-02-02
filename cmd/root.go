package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	git "gopkg.in/src-d/go-git.v4"
)

var (
	path         string
	excludes     []string
	output       string
	includeBares bool
	verbose      bool
	rootCmd      = &cobra.Command{
		Use:   "find_git_repos",
		Short: "find_git_repos will recursively find all git repos under a folder",
		Long:  `find_git_repos could exclude certain folders, ignore bare repos, and export to csv or "git clone" command for easily import`,
		Run: func(cmd *cobra.Command, args []string) {
			walkPath()
		},
	}
)

// Execute executes the root command.s
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode")
	rootCmd.PersistentFlags().BoolVarP(&includeBares, "includeBares", "b", false, "include bare git folder")
	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", ".", "path to find git repos")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "git", "output format(git, csv)")
	rootCmd.PersistentFlags().StringArrayVarP(&excludes, "excludes", "e", []string{}, "excluded paths")
}

func walkPath() {
	if output == "csv" {
		fmt.Println("path,url")
	}
	excludes = append(excludes, "ddd")
	if verbose {
		fmt.Println(excludes)
		fmt.Println("||||||walking")
	}
	pathParam := path
	sep := string(os.PathSeparator)
	if sep == "\\" {
		pathParam = strings.ReplaceAll(pathParam, "/", "\\")
	}
	if sep == "/" {
		pathParam = strings.ReplaceAll(pathParam, "\\", "/")
	}
	if len(pathParam) > 1 && pathParam[0:2] == "."+sep {
		pathParam = pathParam[2:]
	}
	err := filepath.Walk(pathParam,
		func(path string, info os.FileInfo, err error) error {
			for _, exclude := range excludes {
				if len(path) >= len(exclude) && path[0:len(exclude)] == exclude {
					return nil
				}
			}
			if verbose {
				fmt.Println(excludes)
				fmt.Println(path)
			}
			if err != nil {
				fmt.Println(err)
				return err
			}
			if info.IsDir() {
				isGit, url, isBare := getGitRepoStatus(path)
				if isGit {
					excludes = append(excludes, path)
				}
				if isBare && includeBares == false {
					return nil
				}
				if isGit {
					if output == "git" {
						if url == "" {
							fmt.Printf("echo git clone %s %s\n", url, path)
						}
						fmt.Printf("git clone %s %s\n", url, path)
					} else if output == "csv" {
						w := csv.NewWriter(os.Stdout)
						w.Write([]string{path, url})
						w.Flush()
					}
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	// _, errClone := git.PlainClone("./selfhosted", false, &git.CloneOptions{
	// 	URL:      "https://github.com/awesome-selfhosted/awesome-selfhosted.git",
	// 	Progress: os.Stdout,
	// })
	// fmt.Println(errClone)
}

//check if folder is a git repo and if it is a bare repo, and return first remote url if has remote
func getGitRepoStatus(path string) (isGit bool, url string, isBare bool) {
	// fmt.Println(path)
	r, err := git.PlainOpen(path)
	if err != nil {
		isGit = false
		return
	}
	isGit = true
	_, e := r.Worktree()
	if e != nil {
		if e.Error() == "worktree not available in a bare repository" {
			isBare = true
		} else {
			fmt.Println(e)
			return
		}
	}
	remotes, errRemotes := r.Remotes()
	if errRemotes != nil {
		fmt.Println(errRemotes)
		return
	} else if len(remotes) == 0 || len(remotes[0].Config().URLs) == 0 {
		url = ""
	} else {
		url = remotes[0].Config().URLs[0]
	}
	// status, e := tree.Status()
	// if e != nil {
	// 	fmt.Println(e)
	// 	return
	// }
	return
}
