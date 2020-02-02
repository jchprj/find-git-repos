# Install
```
go get -u github.com/jchprj/find_git_repos
```

# Help
find_git_repos -h
```
find_git_repos could exclude certain folders, ignore bare repos, and export to csv or "git clone" command for easily import

Usage:
  find_git_repos [flags]

Flags:
  -e, --excludes stringArray   excluded paths
  -h, --help                   help for find_git_repos
  -b, --includeBares           include bare git folder
  -o, --output string          output format(git, csv) (default "git")
  -p, --path string            path to find git repos (default ".")
  -v, --verbose                verbose mode
```
# Usage
* Exclude folders by prefix  
  ```
  find_git_repos -p . -e aa -e bb
  ```

* Include bare repos
  Default is excluded, to include:  
  ```
  find_git_repos -b
  ```

#  Features

* Omitting subfolders of a git folder
* Will convert input path seperator to os path seperator

# Develop
* Golang environment

* Golang dependencies

```
go get -u gopkg.in/src-d/go-git.v4/...
go get -u github.com/spf13/cobra/cobra
```

* VS Code  
  Debug configuration  
  Go install task