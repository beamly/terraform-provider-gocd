# GoCD SDK 0.6.5

[![GoDoc](https://godoc.org/github.com/drewsonne/go-gocd/gocd?status.svg)](https://godoc.org/github.com/drewsonne/go-gocd/gocd)
[![Build Status](https://travis-ci.org/drewsonne/go-gocd.svg?branch=master)](https://travis-ci.org/drewsonne/go-gocd)
[![codecov](https://codecov.io/gh/drewsonne/go-gocd/branch/master/graph/badge.svg)](https://codecov.io/gh/drewsonne/go-gocd)
[![codebeat badge](https://codebeat.co/badges/1ea74899-2337-4ea6-aaeb-2cc8037fe362)](https://codebeat.co/projects/github-com-drewsonne-go-gocd-master)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fdrewsonne%2Fgo-gocd.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fdrewsonne%2Fgo-gocd?ref=badge_shield)

## CLI

CLI tool to interface with GoCD Server.

### Usage

#### Installation

##### Homebrew

``` bash
brew tap drewsonne/tap
brew install go-gocd
```

##### Manual
Download the latest release from [https://github.com/drewsonne/go-gocd/releases](https://github.com/drewsonne/go-gocd/releases),
and place the binary in your `$PATH`.

#### Quickstart

```
$ gocd configure
? GoCD Server (should contain '/go/' suffix) https://my-go-server:8154/go/
? Client Username my_user
? Client Password *****
? Skip SSL certificate validation (y/N) N
$ gocd list-agents
```

#### Configuration
The library can either be configured using environment variables, cli flags, or a yaml config file.

| Name | CLI Flag | YAML | Environment Variable |
|------|----------|------|----------------------|
| GoCD Server (with `/go/` suffix) | `--server` | `server` | `$GOCD_SERVER` |
| Username | `--username` | `username` | `$GOCD_USERNAME` |
| Password | `--password` | `password` | `$GOCD_PASSWORD` |
| Skip HTTPS/SSL Certification Check | `--skip_ssl_check` | `skip_ssl_check` | `$GOCD_SKIP_SSL_CHECK` |

##### YAML Config File

Run `gocd configure` to launch a wizard which will create a file at `~/.gocd.conf`, or create the file manually:

```yaml
default:
  server: https://goserver:8154/go
  username: admin
  password: mypassword
  skip_ssl_check: true
```

##### Configuration Profiles
Authentication credentials for multiple gocd servers can be stored by using the `--profile` flag.
Configuration Profiles can be created using:
```bash
gocd --profile other-server configure
```
Which will create a new configuration block in `~/.beamly.conf`

Configuration profiles can be used by specifying `--profile` before your command
```bash
gocd --profile other-server list-agents
```

#### Help

    $ gocd -help

## Library

### Usage

Construct a new GoCD client and supply the URL to your GoCD server and if required, username and password. Then use the
various services on the client to access different parts of the GoCD API.
For example:

```go
package main
import (
    "github.com/drewsonne/go-gocd/gocd"
    "context"
    "fmt"
)

func main() {
    cfg := gocd.Configuration{
        Server: "https://my_gocd/go/",
        Username: "ApiUser",
        Password: "MySecretPassword",
    }

    c := cfg.Client()

    // list all agents in use by the GoCD Server
    var a []*gocd.Agent
    var err error
    var r *gocd.APIResponse
    if a, r, err = c.Agents.List(context.Background()); err != nil {
        if r.HTTP.StatusCode == 404 {
            fmt.Println("Couldn't find agent")
        } else {
        	panic(err)
        }
    }

    fmt.Println(a)
}
```

If you wish to use your own http client, you can use the following idiom

```go
package main

import (
    "github.com/drewsonne/go-gocd/gocd"
	"net/http"
    "context"
)

func main() {
    client := gocd.NewClient(
        &gocd.Configuration{},
        &http.Client{},
    )
    client.Login(context.Background())
}
```

## Roadmap ##
This library is still in pre-release. It was initially developed to be an interface for a [gocd terraform provider](https://github.com/drewsonne/terraform-provider-gocd),
which, at this stage, will heavily influence the direction of this library. A list of new features and the expected release
schedule for those features can be found in the project for this github repository.

## Background ##
This library's structure was initially inspired by [https://github.com/google/go-github](https://github.com/google/go-github).
There may still be some vestigial code and structures from this library which will be removed in future revisions and
before v1.0.0 of this library.

## License ##

This library is distributed under the Apache 2.0 license found in [LICENSE](./LICENSE) file.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fdrewsonne%2Fgo-gocd.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fdrewsonne%2Fgo-gocd?ref=badge_large)