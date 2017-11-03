# Orion Setup

### Server Setup

Orion is a hosted application that has several prerequisites. It expects a unix environment; other operating sytems are untested and not officially supported. To run Orion:

* install git on your server (ex. `yum install git`)
* install go on your server
    1. visit `https://golang.org/dl/`
    2. select an appropriate version (this application was built with `1.9.2`)
    3. untar the downloaded file (ex. `tar xvzf go1.9.2.src.tar.gz -C /usr/local/go`)
    4. add the needed env variables; you can do this quickly by adding the following to your resource file/profile:
        ```
        export GOPATH=/home/ec2-user/go
        export GOBIN=/home/ec2-user/go/bin
        PATH=$PATH:$GOBIN:/usr/local/go/bin
        ```
    5. verify your installation (ex. `go version`)
* install the following needed analysis tools:
    * gas (`github.com/GoAST/gas`)

Your server needs (at minimum) internet access to GitHub IP ranges. These ranges can be found [on this page](https://help.github.com/articles/github-s-ip-addresses/). If you're deploying Orion for an GitHub enterprise solution, it will need to be deployed on a server with network access and with access to your GitHub enterprise IP ranges. 

### GitHub Setup

Orion also needs a few things setup in GitHub to work. You'll need:

* a GitHub account (a service account is recommended)