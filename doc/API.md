# Orion API Specification

This document details the Orion API endpoints, inputs and responses.

### Index

**Route:** `/`

**Method(s):** Any

**Input(s):** None

**Response Code(s):** 
* Success (200)

The index endpoint will take a user to the homepage, which will display helpful information about the Orion application. No private/sensitive data should be displayed here. The index page will always return a `200` response as long as the app is running.

---

### Analysis 

**Route:** `/analyze`

**Method(s):** `POST`

**Input(s):** See [Github API spec for PR events](https://developer.github.com/v3/activity/events/types/#pullrequestevent).

**Response Code(s):** 
* Success (`200`); this indicates the request was received and either:
  * the analysis has been successfully initiated
  * the event was not of type `opened`, and was therefore ignored
* Bad Request (`400`)
* Internal Error (`500`)

The analysis route recieves a `PullRequestEvent` from Github whenever a configured repo receives a PR. The application then performs static code analysis on the code in the PR and reports findings (if any) as a comment in the PR itself. If a repo is configured correctly and comments are not left, either A) there are no defects to report, or B) an error occurred within the app itself (in this case, consult the logs for more information). GitHub does not require a response, however one is sent anyway for debugging purposes. For more information on the result of a request, you can look at the message returned to GitHub under the `webhooks` section for the established repository.
