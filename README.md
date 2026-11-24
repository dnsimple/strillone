# Strillone

_Strillone_ is a service to publish the events generated from a DNSimple account to a messaging service, using the DNSimple webhooks.

[![ci](https://github.com/dnsimple/strillone/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/dnsimple/strillone/actions/workflows/ci.yml)

![Example](img-example.png)

## Usage

### Deploy the application

Strillone runs on Fly.io. See [Deployment](#deployment).

### Check the deploy

Make sure the app is properly deployed. If you access the homepage, you should see a JSON response like the following one:

```json
{"ping":"1458412047","what":"dnsimple-strillone"}
```

### Configure the Publishers

Configure the target of the messages. We currently support the following publishers:

- [Slack](#slack-configuration)

See below for the specific configurations.

### Create the webhook

Once you configured the publisher and generated the webhook URL, use the URL to create a new webhook in your DNSimple account.

## Slack configuration

Strillone integrates with Slack using the **Slack Incoming Webhook** feature.

### Step 1: Set up a Slack Incoming Webhook

1. Go to your Slack workspace and [create a new Incoming Webhook](https://docs.slack.dev/messaging/sending-messages-using-incoming-webhooks/)
2. Choose the channel where you want the DNSimple notifications to appear
3. Slack will provide you with a webhook URL that looks like this:

```bash
https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX
```

### Step 2: Create your Strillone webhook URL

To connect Strillone to your Slack webhook, you need to construct a special URL:

1. Take your Strillone application URL (e.g., `https://your-strillone-domain.com`)
2. Add `/slack` to it
3. Extract the unique path from your Slack webhook URL (the part after `services/`)
4. Combine them together

For example:

- Your Strillone app is deployed at: `https://strillone.dnsimple.tools`
- Your Slack webhook URL is: `https://hooks.slack.com/services/T12345/B67890/ABCDEFGHIJKLMNO`
- Your Strillone webhook URL will be: `https://strillone.dnsimple.tools/slack/T12345/B67890/ABCDEFGHIJKLMNO`

### Step 3: Add the webhook to DNSimple

Use this newly created URL when setting up your webhook in DNSimple:

- Either through the [DNSimple dashboard](https://support.dnsimple.com/articles/webhooks/)
- Or via the [DNSimple API](https://developer.dnsimple.com/v2/webhooks/webhooks/)

## Configuration

Strillone reads its configuration from environment variables. `.env.example` lists them with their defaults.

| Name         | Type   | Default                              | Description                                                         |
|--------------|--------|--------------------------------------|---------------------------------------------------------------------|
| BIND_ADDRESS | String | `""`                                 | The address the HTTP server listens on. Empty listens on all interfaces. |
| PORT         | Int    | `4000`                               | The port the HTTP server listens on.                                |
| DNSIMPLE_URL | String | `"https://app.dnsimple.com"`         | Dashboard link host; can be overridden per deployment.              |
| SLACK_URL    | String | `"https://hooks.slack.com/services"` | Slack incoming webhook base URL.                                    |
| LOG_LEVEL    | String | `"info"`                             | Log verbosity: `debug`, `info`, `warn`, or `error`.                 |

## Deployment

Strillone runs on Fly.io in the `dnsimple-ops` organization, one machine per
app:

| Environment | Fly app                      | URL                                        | Webhooks from                                                      |
|-------------|------------------------------|--------------------------------------------|--------------------------------------------------------------------|
| staging     | `dnsimple-strillone-staging` | `https://strillone-staging.dnsimple.tools` | The DNSimple sandbox account                                       |
| production  | `dnsimple-strillone`         | `https://strillone.dnsimple.tools`         | The DNSimple Core, DNSimple Domains and DNSimple Monitors accounts |

The hostnames are CNAMEs to each app's `fly.dev` name, managed in
[infra-dnsimple](https://github.com/dnsimple/infra-dnsimple), and Fly holds a
certificate for each (`fly certs list -a <app>`). Webhooks use these hostnames
rather than `fly.dev`, so moving the app only changes DNS.

The [Fly Deploy](.github/workflows/fly-deploy.yml) workflow deploys them.

1. Merge a PR to main. Once CI passes, the workflow deploys the merged commit
   to staging. Wait for that deploy before promoting: production only takes a
   commit on `main` that deployed to staging successfully.
2. Run the workflow with the `production` environment, from the Actions tab or
   with `gh workflow run fly-deploy.yml -f environment=production`. It deploys
   the commit staging runs, prints it in the run summary, and refuses if that
   commit is not on `main` or its staging deploy did not succeed.

Start every manual run from `main`. Both GitHub environments only accept runs
from `main`.

- **Staging** deploys main's head when run by hand, or any commit passed as
  `commit`, so a branch can be tried on staging. It stays there until the next
  green CI run on `main` replaces it.
- **Production** takes a `commit` only as a rollback, to a commit that
  production ran before: it must be a completed release among the latest 25
  releases in the production app's Fly release history, and have a green
  `fly/production` status. To release a new commit, deploy it to staging and promote staging.
- **Automatic deploys never move staging backwards.** Only the newest commit
  on `main` with a green CI run deploys. A re-run of an older CI run, or one
  that finishes after a newer one, deploys nothing.

A deploy replaces the app's one machine in place, or creates it if the app has
none, and refuses an app with more than one. It succeeds once that machine is
started on the new image and answers Fly's health check on `/`. The commit's
`fly/<environment>` status records the outcome. Production never moves
backwards on its own either: a promotion of a commit older than the one
production runs is refused, even one that waited behind another deploy.

A deploy that fails its health check leaves the machine on the failed image.
Check `fly logs`, then:

- **Staging:** run the workflow again for the same commit, or deploy a fix.
- **Production:** roll back by running the workflow with the commit
  production ran before. To retry a failed promotion instead, re-run its
  failed jobs from the Actions tab, which keeps the commit it chose, unless
  production has moved to a newer commit since.

Strillone has no secrets: each webhook URL carries its Slack token. Each GitHub
environment (`staging`, `production`) holds a `FLY_API_TOKEN` deploy token for
its app. A production promotion reads the staging machine with the staging
token.

```shell
fly logs -a dnsimple-strillone
```

## About the name

The word [strillone](https://en.wiktionary.org/wiki/strillone) (literally _someone who shouts a lot_, in practice the equivalent of _newspaper boy_) comes from Italian and it refers to the newspaper sellers in the street, who were used to yell the titles in the front page to catch the attention and sell more newspapers.

![Strillone](strillone.jpg)

Photo: [New York Media](http://nymag.com/daily/intelligencer/2013/06/fed-is-having-a-1936-moment.html)

## License

Copyright (c) 2016-2024 DNSimple Corporation. This is Free Software distributed under the MIT license.
