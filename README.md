# github-activity

See recent activity of a github user

Designed for https://roadmap.sh/projects/github-user-activity

## Building

**You must have Go installed**

I recommend using the build script `build.sh` to create a binary file with the proper name.

## Usage

Simply run

```
github-activity <username>
```
To get a short description of the user's recent public activity

## Example

Running 
```
$ github-activity ursaru-tudor
```

May output

```
Recent activity of ursaru-tudor
 - Pushed 1 commits to ursaru-tudor/roadmap-sh-github-activity
 - Made the repo ursaru-tudor/task-cli public
 - Pushed 4 commits to ursaru-tudor/roadmap-sh-github-activity
 - Pushed 1 commits to ursaru-tudor/task-cli
 - Pushed 1 commits to ursaru-tudor/task-cli
 - Pushed 2 commits to ursaru-tudor/task-cli
 - Is now watching golang/go
 - Made the repo ursaru-tudor/task-cli public
```