TweetHog - Stream, filter and react to Twitter status updates
=============================================================

[![Build Status](https://travis-ci.org/lastzero/tweethog.png?branch=master)][ci]
[![Code Quality](https://goreportcard.com/badge/github.com/lastzero/tweethog)][code quality]
[![Test Coverage](https://gocover.io/_badge/github.com/lastzero/tweethog)][test coverage]
[![GitHub issues](https://img.shields.io/github/issues/lastzero/tweethog.svg)][issues]
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)][license]

[ci]: https://travis-ci.org/lastzero/tweethog
[code quality]: https://travis-ci.org/lastzero/tweethog
[test coverage]: http://gocover.io/github.com/lastzero/tweethog
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
# ./tweethog -t trump,cat
Topics       : trump,cat
Languages    : en
URLs         : false
Retweets     : false
Replies      : false
Via          : false
Max mentions : 1
Max tags     : 2
Like tweets  : false
2017/10/26 14:37:56 Starting Twitter stream...

ID: 923529145740558337  Date: 2017-10-26 12:38:01  User: @SonjaFazzio  Following: 128  Followers: 70  Likes: 6431
>>> Trump's fiscal year begins soon. Give it 90 days...time will reveal all this so called winning.
I'm not holding my breath.

ID: 923529152241848325  Date: 2017-10-26 12:38:02  User: @terrypruitt  Following: 407  Followers: 179  Likes: 151
>>> “Cat stretches first thing when he wakes up because he never knows when he will need to pounce.” #postrunwisdom #baltimore1/2

ID: 923529161184108544  Date: 2017-10-26 12:38:05  User: @Wombleloyalist  Following: 1134  Followers: 770  Likes: 652
>>> And i have Donald trump as his assistant
```

Commands
--------

Command   | Description
----------|-------------------------------------------------------------
config    | Displays all configuration values
filter    | Shows all matching tweets without performing any action
like      | Automatically likes all matching tweets
smartlike | Likes tweets with random delay and rate limit

Global Flags
------------

Name                          | Description
------------------------------|------------------------------------------------------------------
--config-file value, -c value | YAML config filename (default: "config.yml")
--consumer-key value          | Twitter API consumer key
--consumer-secret value       | Twitter API consumer secret
--access-token value          | Twitter API access token
--access-secret value         | Twitter API access token secret


Filter Flags
------------
Name                          | Description
------------------------------|------------------------------------------------------------------
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

All of the flags above can be set in the config file as a default. Example:

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
