# Orion API Specification

This document details the Orion API endpoints, inputs and responses.

### Index

**Route:** /

**Method(s):** Any

**Input(s):** None

**Response Codes:** 200

The index endpoint will take a user to the homepage, which will display helpful information about the Orion application. No private data should be displayed here. The index page will always return a `200` response as long as the app is running.

---

### Analysis 

**Route:** /analyze

**Method(s):** POST

**Input(s):** See [Github API spec for PR events](https://developer.github.com/v3/activity/events/types/#pullrequestevent).

**Response Code(s):** 200

The analysis route recieves a `PullRequestEvent` from Github whenever a configured repo receives a PR. The application then performs static code analysis on the code in the PR and reports findings (if any) as a comment in the PR itself. This route should never return anything other than a `200` response, because the message sent is only an initiation (Github does not expect a response on the other end).
