# Strillone

_Strillone_ is a service to publish the events generated from a DNSimple account to a messaging service, using the DNSimple webhooks.

[![ci](https://github.com/dnsimple/strillone/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/dnsimple/strillone/actions/workflows/ci.yml)

![Example](img-example.png)

## Usage

### Deploy the application

You can use the following button to deploy the service to Heroku.

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/dnsimple/strillone)

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

- Your Strillone app is deployed at: `https://my-strillone-app.herokuapp.com`
- Your Slack webhook URL is: `https://hooks.slack.com/services/T12345/B67890/ABCDEFGHIJKLMNO`
- Your Strillone webhook URL will be: `https://my-strillone-app.herokuapp.com/slack/T12345/B67890/ABCDEFGHIJKLMNO`

### Step 3: Add the webhook to DNSimple

Use this newly created URL when setting up your webhook in DNSimple:
- Either through the [DNSimple dashboard](https://support.dnsimple.com/articles/webhooks/)
- Or via the [DNSimple API](https://developer.dnsimple.com/v2/webhooks/webhooks/)

## About the name

The word [strillone](https://en.wiktionary.org/wiki/strillone) (literally _someone who shouts a lot_, in practice the equivalent of _newspaper boy_) comes from Italian and it refers to the newspaper sellers in the street, who were used to yell the titles in the front page to catch the attention and sell more newspapers.

![Strillone](strillone.jpg)

Photo: [New York Media](http://nymag.com/daily/intelligencer/2013/06/fed-is-having-a-1936-moment.html)

## License

Copyright (c) 2016-2024 DNSimple Corporation. This is Free Software distributed under the MIT license.
