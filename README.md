# OBS_Deezer
I guess it doesn't have to be used for OBS but yeah..

---
## What is this?
This is a small go program to get the currently playing track from Deezer. I'm using if for OBS but you can use it for whatever.

---
## Why aren't you using the Deezer API?
Caching on the API by the looks of things. In order to get the song playing *RIGHT THIS SECOND*, we use the profile page.
Also this way we don't need to oauth.

---
## Usage
You need to compile it, it's go so it's not that hard..
Literially ```go build src/github.com/SilverCory/OBS_Deezer```

There's some arguments, use the -help argument to see them.