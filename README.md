# a9sharvest

a9sharvest is a CLI tool for [Harvest](https://www.getharvest.com/).

# Installation

```shell
$ brew tap anynines/tap
$ brew install a9sharvest
```

# Requirements

Before you can use `a9sharvest` you need a token and the account ID.
You can get both by accessing the [developer tools](https://id.getharvest.com/developers)
and creating a new token.

# Usage

```shell
$ export ACCOUNT_ID=12345
$ export TOKEN=12345.pt.Avfe-WEFWEF...D4z

$ export TAGS="[meeting_orga_lane],[support_lane]"
$ export SKIP_PROJECT_IDS="12345,6789" # do not skip any project by default
$ export FROM="20200901" # 14 days ago by default
$ export TO="20200914" # today by default

$ a9sharvest group
[unknown] = 68.17
[meeting_orga_lane] = 0.4
[support_lane] = 2.93
```

# Manual Release Building

```shell
git tag -a v1.1.0
GOOS=darwin GOARCH=amd64 go build -ldflags "-X github.com/anynines/a9sharvest/pkg/version.Version=v1.1.0"
```

# Links

 - [Developer tools](https://id.getharvest.com/developers)
 - [API documentation](http://help.getharvest.com/api-v2/)
