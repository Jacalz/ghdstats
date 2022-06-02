# Ghdstats

Ghdstats is a simple, but very fast, tool to fetch download counts from GitHub releases.
You just give it a user, or a user and repository, and it fetches all the necessary details for you.

## Installation

The cli command can be downloaded by running the following command (using Go 1.16 or above):
```
$ go install github.com/Jacalz/ghdstats@latest
```

## Usage

The general usage of the program works like this:
```
$ ghdstats [user] [repository, optional]
```

As an example, you can get all the data for [jacalz/rymdport](https://github.com/jacalz/rymdport):
```
$ ghdstats jacalz rymdport
```

The same command can also be written as this:
```
$ ghdstats jacalz/rymdport
```

The tool can also fetch all downloads for a given user or organization:
```
$ ghdstats jacalz
```

## Performance

The goal here is to be as fast as possible but still keeping very readable code.
All repositories for a user are downloaded and processed concurrently for better performance on multicore systems.

## Inspiration and resoning

A lot of the inspiration for this came from https://github.com/mmilidoni/github-downloads-count.
THe idea was to create a faster tool without any need for Python. The plan is to have binaries for as
many 64-bit platforms as possible.
