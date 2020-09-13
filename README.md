# a9sharvest

a9sharvest is a CLI tool for [Harvest](https://www.getharvest.com/).

# Installation

```shell
$ brew tap anynines/a9sharvest
$ brew install a9sharvest
```

# Usage

```shell
$ export ACCOUNT_ID=12345
$ export TOKEN=12345.pt.Avfe-WEFWEF...D4z

$ export TAGS="[meeting_orga_lane],[support_lane]"
$ export FROM="20200901" # 14 days ago by default
$ export TO="20200914" # today by default

$ a9sharvest group
[unknown] = 68.17
[meeting_orga_lane] = 0.4
[support_lane] = 2.93
```

# Manual Release Building

```shell
git tag -a stable-0.1.0
GOOS=darwin GOARCH=amd64 go build -ldflags "-X github.com/anynines/a9sharvest/pkg/version.Version=stable-0.1.0"
```
