TweetHog - Stream, filter and like Twitter status updates
=========================================================

Build
-----
```
git clone https://github.com/lastzero/tweethog.git
cd tweethog
go get -u github.com/golang/dep/cmd/dep
dep ensure
go build
```

Configuration
-------------
1. Copy config.example.yml to config.yml
2. Create your own Twitter API keys & tokens on https://apps.twitter.com/app/new
3. Put them in config.yml

Usage
-----
```
# ./tweethog -t cat -l en --no-retweets --like
Started streaming Twitter status updates on Mon, 16 Oct 2017 15:17:02 CEST...
Topics      : cat
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