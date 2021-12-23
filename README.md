# Ghdstats

Ghdstats is a simple, but very fast, tool to fetch download counts from GitHub releases.
You just give it a user, or a user and repository, and it gets all the necessary details.

## Usage

The general usage of the program works like this:
```
$ ghdstats [user] [repository, optional]
```

As an example, you can get all the data for [fyne-io/calculator](https://github.com/fyne-io/calculator):
```
$ ghdstats fyne-io calculator
```

Or, you can alternatively get data for all repositories in [fyne-io](https://github.com/fyne-io):
```
$ ghdstats fyne-io
```

## Performance

The goal here is to be as fast as possible but still keeping very readable code.
All repositories for a user are downloaded and processed concurrently for better performance on multicore systems.

## Inspiration and resoning

A lot of the inspiration for this came from https://github.com/mmilidoni/github-downloads-count.
I wanted to create something faster that could run without having any Python-stuff installed
and easily be compiled for multiple platforms. The plan is to have binaries for as
many 64-bit platforms as possible.
