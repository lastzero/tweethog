TweetHog - Stream, filter and like Twitter status updates
=========================================================

Installation
------------
Make sure you got the latest version of Go installed on your computer.

It can be downloaded for free at https://golang.org/dl/.

Then open a terminal and enter:

```
go get github.com/lastzero/tweethog
```

Configuration
-------------
In order to stream status updates, TweetHog needs access to the Twitter API:

1. Copy (or rename) config.example.yml to config.yml
2. Create your own Twitter API keys & tokens on https://apps.twitter.com/app/new
3. Put them in config.yml by replacing the placeholders

Hint: You can use the `--config` flag to use a different config file.

Note: Due to an [issue](https://github.com/dghubble/go-twitter/issues/61)
with the go-twitter library, you won't see any error when using
invalid credentials. We're working on it.

Usage
-----
```
# tweethog -t cat -t dog -l en --no-retweets --like
Started streaming Twitter status updates on Mon, 16 Oct 2017 15:17:02 CEST...
Topics      : cat, dog
Languages   : en
URLs        : true
Retweets    : false
Like tweets : true
...
ID: 919915098877190144  Date: Mon Oct 16 13:17:05 +0000 2017  User: @lindsseybb  Following: 78  Followers: 52  Likes: 617
>>> Cat nip plants are resilient y'all. My pot has been tipped over numerous times, I haven't watered it once and it just keeps coming back.
Liked ❤️
```

Flags
-----
```
   --topic value, -t value   Stream filter topic (cat, dog, fish, ...)
   --lang value, -l value    Stream filter language (en, de, fr, ...)
   --max-followers value     User max follower count (0 for unlimited) (default: 5000)
   --min-followers value     User min follower count (default: 5)
   --max-following value     User max following count (0 for unlimited) (default: 5000)
   --min-following value     User min following count (default: 5)
   --no-retweets             Exclude tweets starting with RT or @
   --no-urls                 Exclude tweets containing URLs
   --like                    Like tweets
   --config value, -c value  Config file name (default: "config.yml")
   --help, -h                show help
   --version, -v             print the version
```

About
-----
This tool was created by [Michael Mayer](https://blog.liquidbytes.net/about).
You can use it for free under the terms of the MIT license.

Please feel free to send an e-mail if you have any questions or just want to say hello.
Contributions are most welcome, even if it's just a bug report or a tiny pull-request to fix a typo.
