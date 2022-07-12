# Ghdstats

Ghdstats is a simple, but very fast, tool to fetch download counts from GitHub releases.
You just give it a user, or a user and repository, and it fetches all the necessary details for you.

## Installation

TODO: This software is pending a rewrite in Rust (from a learning point-of-view).

## Performance

The goal here is to be as fast as possible but still keeping very readable code.
All repositories for a user are downloaded and processed concurrently for better performance on multicore systems.

## Inspiration and resoning

A lot of the inspiration for this came from https://github.com/mmilidoni/github-downloads-count.
THe idea was to create a faster tool without any need for Python. The plan is to have binaries for as
many 64-bit platforms as possible.
