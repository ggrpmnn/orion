# Orion Setup

### Server Setup

Orion is a hosted application that has several prerequisites. It expects a unix environment; other operating sytems are untested and not officially supported. To run Orion:

* install git on your server (e.g. by running `yum install git`)
* install `go` on your server
    1. visit `https://golang.org/dl/`
    2. select an appropriate version (this application was built with version `1.9.2`)
    3. untar the downloaded file (e.g. by running `tar xvzf go1.9.2.src.tar.gz -C /usr/local/go`)
    4. add the needed env variables; you can do this quickly by adding the following to your resource file/profile:
        ```
        export GOPATH=/home/<your-user>/go
        export GOBIN=/home/<your-user>/go/bin
        PATH=$PATH:$GOBIN:/usr/local/go/bin
        ```
    5. verify your installation (e.g. by running `go version`)
* install the following needed analysis tools:
    * gas (`github.com/GoAST/gas`)
    * more to come soon

Your server needs (at minimum) internet access to GitHub IP ranges, which can be found [on this page](https://help.github.com/articles/github-s-ip-addresses/). If you're deploying Orion for a GitHub enterprise solution, it will need to be deployed on a server with network access and with access to your GitHub enterprise IP ranges.

Once your server is configured, you'll need to download and build Orion. To do that, `git clone` this repo and run the command `go build -o orion <path-to-orion-project>/source/*.go` in the cloned directory. Then simply run Orion (`./orion`), and the application will be awaiting an incoming webhook message (see below)!

To scan a repo, Orion will `git clone` that repo in order to operate on the files. It will, by default, clone those repos into its own working folder (`<path-to-orion-project>/source`). If you need to change where Orion downloads source code for any reason, you can do that by setting a path into your environment under `ORION-WORKSPACE`. The user/group for Orion will need to be able to access this folder as well, so double check the permissions on any workspace folders before running.

### GitHub Setup

Orion also needs a few things setup in GitHub to work. You'll need:

* a GitHub account (a service account is heavily recommended, but any account will work); the username of the account should be added to the env as `GH_USERNAME` (e.g. `export GH_USERNAME=<your username here>`)
* a GitHub authentication token for the above account ([see here](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/) for more info), which needs to be added to the env as `GH_AUTH_TOKEN` (e.g. `export GH_AUTH_TOKEN=<your token here>`)

 To onboard a repository to use Orion, we need to add a GitHub webhook to the repo. [View this page](https://developer.github.com/webhooks/creating/#setting-up-a-webhook) for more info on adding a webhook to your project.

 * the webhook's Payload URL should be of the form `http://<your-server-address>:8080/analyze`
 * the content type is `application/json`
 * for "which events...", select `Let me select individual events` and then check the box for `Pull request`
 * ensure that the `Active` box is checked

 That's it! Everything should now work and whenever a PR is submitted to your repo, orion will scan and post the results in a comment.