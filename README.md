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

$ export LOG_LEVEL=debug # `debug`, `trace` or empty

$ export TAGS="[meeting_orga_lane],[support_lane]"

$ export SKIP_PROJECT_IDS="12345,6789" # do not skip any project by default
$ export ALLOWED_USER_IDS="3344,4444" # empty by default

$ export FROM="20200901" # 14 days ago by default
$ export TO="20200914" # today by default

$ a9sharvest group
  [customer_lane]     |  3.00
  [side_project_lane] |  3.21
  [standard_lane]     | 33.56
  [meeting_orga_lane] | 13.66
  [misc_lane]         |  8.75
  [support_lane]      | 23.69
```

# Manual Release Building

```shell
git tag -a v1.1.0
GOOS=darwin GOARCH=amd64 go build -ldflags "-X github.com/anynines/a9sharvest/pkg/version.Version=v1.1.0"
```

# Links

 - [Developer tools](https://id.getharvest.com/developers)
 - [API documentation](http://help.getharvest.com/api-v2/)
