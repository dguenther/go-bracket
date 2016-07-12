go-bracket [![Build Status](https://travis-ci.org/dguenther/go-bracket.svg?branch=master)](https://travis-ci.org/dguenther/go-bracket)
===========

An API client for calling Challonge and Smash.GG and receiving uniform
results.

Currently, the client only supports a limited subset of the fields available
from the APIs. Additionally, it only supports fetching data, not updating it.
Feel free to open a PR if you'd like to add support for more fields or
features.

How to use
==========
Instantiate the client with API keys:

`client := bracket.NewClient(challongeUser, challongeApiKey)`

Then pass a web URL into the fetch function:

`b := client.FetchBracket("http://challonge.com/xyfuz5c3")`

or

`b := client.FetchBracket("https://smash.gg/tournament/super-smash-sundays-48/brackets/14221/50133/165583")`
