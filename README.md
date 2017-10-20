TweetHog - Stream, filter and react to Twitter status updates
=============================================================

[![GitHub release](http://img.shields.io/github/release/lastzero/tweethog.svg?style=flat-square)][release]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]

[release]: https://github.com/lastzero/tweethog/releases
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
Make sure you got the latest version of Go installed on your computer. It can be downloaded for free at https://golang.org/dl/.

Then open a terminal and enter:

```
go get github.com/lastzero/tweethog
```

Configuration
-------------
In order to stream status updates, TweetHog needs access to the Twitter API:

1. Copy (or rename) [config.example.yml](config.example.yml) to config.yml
2. Create your own Twitter API keys & tokens on https://apps.twitter.com/app/new
3. Put them in config.yml by replacing the placeholders

Hint: You can use the `--config` flag to use a different config file.

Note: Due to an [issue](https://github.com/dghubble/go-twitter/issues/61)
with the go-twitter library, you won't see any error when using
invalid credentials. We're working on it.

Example Usage
-------------
```
# tweethog -t cat -l en --like
Started streaming Twitter status updates on Mon, 16 Oct 2017 15:17:02 CEST...
Topics      : cat
Languages   : en
Like tweets : true
...
ID: 919915098877190144  User: @lindsseybb  Following: 78  Followers: 52  Likes: 617
>>> Cat nip plants are resilient y'all. My pot has been tipped over numerous times,
I haven't watered it once and it just keeps coming back.
Liked ❤️
```

Flags
-----
Name                     | Description
-------------------------|------------------------------------------------------------------
--topic value, -t value  | Stream filter topic e.g. cat, dog, fish
--lang value, -l value   | Stream filter language e.g. en, de, fr
--max-followers value    | User max followers, 0 for unlimited (default: 5000)
--min-followers value    | User min followers (default: 5)
--max-following value    | User max following, 0 for unlimited (default: 5000)
--min-following value    | User min following (default: 5)
--max-tags value         | Max number of hash #tags (default: 2)
--max-mentions value     | Max number of @mentions (default: 1)
--retweets               | Include tweets starting with RT
--replies                | Include tweets starting with @
--via                    | Include tweets containing via @
--urls                   | Include tweets containing URLs
--like                   | Like all matching tweets
--smart-like             | Likes tweets with random delay and rate limit
--config value, -c value | Config file name (default: "config.yml")
--help, -h               | show help
--version, -v            | print the version

About
-----
This tool was created by [Michael Mayer](https://blog.liquidbytes.net/about). It is still under
development and not very mature. You can use it for free under the terms of the MIT license.

Please feel free to send an e-mail if you have any questions or just want to say hello.
Contributions are most welcome, even if it's just a bug report or a tiny pull-request to fix a typo.
