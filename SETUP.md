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
        export GOPATH=/home/<your-user>/go
        export GOBIN=/home/<your-user>/go/bin
        PATH=$PATH:$GOBIN:/usr/local/go/bin
        ```
    5. verify your installation (ex. `go version`)
* install the following needed analysis tools:
    * gas (`github.com/GoAST/gas`)

Your server needs (at minimum) internet access to GitHub IP ranges. These ranges can be found [on this page](https://help.github.com/articles/github-s-ip-addresses/). If you're deploying Orion for an GitHub enterprise solution, it will need to be deployed on a server with network access and with access to your GitHub enterprise IP ranges. 

### GitHub Setup

Orion also needs a few things setup in GitHub to work. You'll need:

* a GitHub account (a service account is heavily recommended, but any account will work)
* a GitHub authentication token for the above account ([see here](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/) for more info), which needs to be added as `GH_AUTH_TOKEN` to the env (ex. `export GH_AUTH_TOKEN=<your token here>`)

 To onboard a repository to use Orion, we need to add a GitHub webhook to the repo. [View this page](https://developer.github.com/webhooks/creating/#setting-up-a-webhook) for more info on adding a webhook to your project.

 * the webhook's Payload URL should be of the form `http://<your-server-address>:8080/analyze`
 * the content type is `application/json`
 * for "which events...", select `Let me select individual events` and then check the box for `Pull request`
 * ensure that the `Active` box is checked

 That's it! Everything should now work and whenever a PR is submitted to your repo, orion will scan and post the results in a comment.