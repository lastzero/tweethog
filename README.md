TweetHog - Stream, filter and react to Twitter status updates
=============================================================

[![Build Status](https://travis-ci.org/lastzero/tweethog.png?branch=master)][ci]
[![Code Quality](https://goreportcard.com/badge/github.com/lastzero/tweethog)][code quality]
[![GitHub issues](https://img.shields.io/github/issues/lastzero/tweethog.svg)][issues]
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)][license]

[ci]: https://travis-ci.org/lastzero/tweethog
[code quality]: https://travis-ci.org/lastzero/tweethog
[issues]: https://github.com/lastzero/tweethog/issues
[license]: https://github.com/lastzero/tweethog/blob/master/LICENSE

This tool provides an easy way to stream, filter and optionally
react to tweets based on topic, language and user profile. It is similar
to commercial SaaS offerings such as TweetFull, RoboLike or Twitfox. However...

* TweetHog is free, fast and doesn't require setting up a server
* You stay in control of your data and don't need to give access to your Twitter account
* You can improve your Go skills if you choose to adapt it to your specific needs

![TweetHog](logo.png)

Installation
------------
Make sure you got the latest version of Go installed on your computer.
It can be downloaded for free at https://golang.org/dl/.
You might need to install Xcode and agree to its license first on OS X.

Then open a terminal and enter:

```
go get github.com/lastzero/tweethog/cmd/tweethog
```

To update it afterwards:

```
go get -u github.com/lastzero/tweethog/cmd/tweethog
```

Note: Your PATH, GOPATH and GOROOT environment variables
need to be set correctly for `go get` to work as expected.

Configuration
-------------
In order to stream status updates, TweetHog needs access to the Twitter API:

1. Copy [config.example.yml](config.example.yml) to config.yml in the directory you wish to execute the command
2. Create your own Twitter API keys & tokens on https://apps.twitter.com/app/new
3. Put them in config.yml by replacing the placeholders

Hint: You can use the `--config-file` flag to use a different config file.

Note: Due to an [issue](https://github.com/dghubble/go-twitter/issues/61)
with the go-twitter library, you won't see any error when using
invalid credentials. We're working on it.

Example Usage
-------------
```
# tweethog -t cat -l en --like
Started streaming Twitter status updates on 2017-10-26 11:51:58...
Topics       : cat
Languages    : en
URLs         : false
Retweets     : false
Replies      : false
Via          : false
Max mentions : 1
Max tags     : 2
Like tweets  : true

ID: 921269579980115969  User: @onthelooseliam  Following: 563  Followers: 825  Likes: 9248
>>> My cat is playing with my bun and I'm terrified
Liked ❤️
```

Flags
-----
Name                          | Description
------------------------------|------------------------------------------------------------------
--consumer-key value          | Twitter API consumer key
--consumer-secret value       | Twitter API consumer secret
--access-token value          | Twitter API access token
--access-secret value         | Twitter API access token secret
--topic value, -t value       | Stream filter topic e.g. cat, dog, fish
--lang value, -l value        | Stream filter language e.g. en, de, fr
--min-followers value         | User min followers (default: 0)
--max-followers value         | User max followers, 0 for unlimited (default: 0)
--min-following value         | User min following (default: 0)
--max-following value         | User max following, 0 for unlimited (default: 0)
--max-tags value              | Max number of hash #tags (default: 0)
--max-mentions value          | Max number of @mentions (default: 0)
--retweets                    | Include tweets starting with RT
--replies                     | Include tweets starting with @
--via                         | Include tweets containing via @
--urls                        | Include tweets containing URLs
--like                        | Like all matching tweets
--smart-like                  | Likes tweets with GetRandomInt delay and rate limit
--config-file value, -c value | YAML config filename (default: "config.yml")

All of the flags above (except `config-file` of course) can be
set in the config file as a default. Example:

```
max-mentions: 1
max-tags: 2
lang: en
urls: true
retweets: false
replies: false
via: false
```

About
-----
This tool was created by [Michael Mayer](https://blog.liquidbytes.net/about)
in an attempt to learn Go the agile way: Start with something small that works
and then refactor and improve. You can use it for free and at your own risk
under the terms of the MIT license.

Please feel free to send an e-mail if you have any questions or just want to say hello.
Contributions are most welcome, even if it's just a bug report or a tiny pull-request to fix a typo.
