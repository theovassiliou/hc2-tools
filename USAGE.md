# Usage Scenarios

This document contains examples configuration on how the hc2-tools tools can be combined to efficiently create, modify and upload lua scripts.

- [Usage Scenarios](#usage-scenarios)
  - [Prerequisite](#prerequisite)
  - [Download all lua scenes](#download-all-lua-scenes)
    - [Goal](#goal)
    - [Why should I do this](#why-should-i-do-this)
    - [Steps](#steps)
    - [Results](#results)
  - [Uploading from Visual Studio Code](#uploading-from-visual-studio-code)
    - [VSC Upload Integration Goal](#vsc-upload-integration-goal)
    - [Why should I do an upload integration](#why-should-i-do-an-upload-integration)
    - [Steps to the upload integrttion](#steps-to-the-upload-integrttion)
  - [Run a scene on the Fibaro HC2 from the terminal](#run-a-scene-on-the-fibaro-hc2-from-the-terminal)
    - [Why should I run start a scene from the command line](#why-should-i-run-start-a-scene-from-the-command-line)
    - [Steps to run a scene](#steps-to-run-a-scene)
  - [Run a scene on the Fibaro HC2 from Visual Studio Code](#run-a-scene-on-the-fibaro-hc2-from-visual-studio-code)
    - [Why should I run a scene from VSC](#why-should-i-run-a-scene-from-vsc)
    - [Steps to run a scene from VSC](#steps-to-run-a-scene-from-vsc)
  - [Upload a script on git commit](#upload-a-script-on-git-commit)
    - [Why should upload a script on git commit](#why-should-upload-a-script-on-git-commit)
    - [Prerequisites](#prerequisites)
    - [Steps to integrate into git commit](#steps-to-integrate-into-git-commit)
  - [FIBARO_GIT_HOOK Header](#fibarogithook-header)
    - [Code snipped for FIBARO_GIT_HOOK header](#code-snipped-for-fibarogithook-header)

## Prerequisite

All examples assume that you have [installed](README.md#Installation) all hc2-tools command line tools and that they there are included in your search path.

hc2-tools installs different tools that solve different task when developing lua scripts for the Fibaro HC2 system. Mainly

- [hc2DownloadScene](https://github.com/theovassiliou/hc2-tools/cmd/hc2DownloadScene/README.md)
- [hc2UploadScene](https://github.com/theovassiliou/hc2-tools/cmd/hc2UploadScene/README.md)
- [hc2SceneInteract](https://github.com/theovassiliou/hc2-tools/cmd/hc2SceneInteract/README.md)

Information on how a single file relates to a particular scene is captured in what is called [FIBARO_GIT_HOOK](#FIBARO_GIT_HOOK-Header)-Header.

## Download all lua scenes

This is the typical scenario when you are using the hc2-tools toolset, and you have already some scenes active in the FIBARO.

### Goal

You have a set of scenes created in the Fibaro system, and you would like to download all of them.

### Why should I do this

There are numerous reasons why you should do this. One is to be able to version control the scripts, what you haven't done until know. Another reason can be, to analyse the scripts, and perhaps act on this. For example delete all scenes that are disabled.

### Steps

```shell
hc2Download -d ./download -s -1 --create-header

INFO[0000] No LuaSpec in file: download/House/Livingroom/VSLLivingroom.lua
INFO[0000] Adding ...
INFO[0000] No LuaSpec in file: download/House/Floor/VSLFloor.lua
[...]
INFO[0031] retrieved 87 scenes
INFO[0031] created 87 files
INFO[0031] wrote 1555448 bytes  
```

`hc2DownloadScene` informs that no [FIBARO_GIT_HOOK](#hc2-tools-header) is present in the downloaded file and that it has been before the file has been saved (option `--create-header`).

### Results

As a result you will get in the download folder (`./download`) in this case 87 files.

```shell
tree ./download
./download
├── FibaroDimmerExample.lua
├── Storage
│   ├── SensorPool
│   │   ├── CentralSceneTest.lua
│   │   ├── GoodNight.lua
│   │   ├── KU\ PhiliioButtonv2WI.lua
│   │   ├── LeavingHome.lua
...
├── House
│   ├── Livingroom
│   │   ├── ObserveGlobals.lua
│   │   ├── RecordLuxReadings.lua
│   ├── Floor
│   │   ├── FLMainLightOnOff.lua
│   │   ├── CentralSceneHandler.lua
│   │   ├── RecordLuxReadings.lua
│   │   └── VSLFloor.lua
...
```

The first level of directories are the sections as defined in the FIBARO system. The 2nd level are the rooms in each section. On each level lua files can be located. In each lua-file the FIBARO_GIT_HOOK header has been added if it hasn' been present. It ensures that an updated file is replacing the exisitng one.

## Uploading from Visual Studio Code

When developing lua scripts it is not unusual that you are using an IDE (Integrated Development Environment) like for example Visual Studio Code. `hc2UploadScene` has been created to enable a seamlesss integration. The usual "copy from text editor" and "paste to webbrowser" approach is not really suitable for more serious development.

### VSC Upload Integration Goal

Configuring VCS in a way to be able with a hit of a button to upload the script currently being edited.

### Why should I do an upload integration

Such an integration speeds up the turn-around time for script development.

### Steps to the upload integrttion

In the example directory you will find [`tasks-lua.json`](examples/tasks-lua.json) that can be used as a starting point. In particular the task labeled `Upload` is of interest here:

```json
{
    "label": "Upload",
    "type": "shell",
    "command": "hc2UploadScene",
    "args": [
        "${file}"
    ],
    "problemMatcher": [
        "$vsls"
    ],
    "presentation": {
      "reveal": "never"
    }
},
```

add this to your `.vscode/tasks.json` file in your workspace.  Open a lua file, you have downloaded earlier via `hc2DownloadScene` and verify that a FIBARO_GIT_HOOK header has been included.

Run the task by selecting `Terminal->Run Task ...` `Upload`.

If everything works out well the file will be silently uploaded and live now on the FIBARO system. Go to the FIBARO system, select your scene and check whether your changes have arrived.

## Run a scene on the Fibaro HC2 from the terminal

Running a scene from your favorite IDE or the command line.

Triggering the execution of an uploaded scene and retrieve the debug messages for this scene.

### Why should I run start a scene from the command line

When navigating through the downloaded scenes on the terminal you would like to run a scene and check the current status of a specific, uploaded file.

### Steps to run a scene

You would like to execute the scene, with the sceneID 205 and then look at the debug messages.

```shell
hc2SceneInteract -s 205 -a start -g --tail
[DEBUG] 14:37:46: Motion detected
[DEBUG] 14:40:22: Motion detected
```

Passing `-s 205` to `hc2SceneInteract` triggers the execution of sceneID 205 with the action `-a start`. `-g` triggers the display of the existing debug messages, while `--tail` instructs `hc2SceneInteract` to wait for new debug messages.

```shell
hc2SceneInteract -f ObserveGlobals.lua -g
[DEBUG] 20:54:00: Global variable SleepState changed.
```

Running `hc2SceneInteract` with `-f ObserveGlobals.lua` extracts the sceneID from the `FIBARO_GIT_HOOK`-header and `-g` retrieves the debug messages that have been generated so far. No action is triggered.  

## Run a scene on the Fibaro HC2 from Visual Studio Code

### Why should I run a scene from VSC

After editing and uploading a scene you would like to run the scene and see the effect by looking into the debug messages to asses the success from within an IDE

### Steps to run a scene from VSC

In the example directory you will find [`tasks-lua.json`](examples/tasks-lua.json) that can be used as a starting point. In particular the tasks labeled `Start` and `Tail` are of interest here:

```json
{
    "label": "Start",
    "type": "shell",
    "command": "hc2SceneInteract",
    "args": [
        "-a",
        "start",
        "-f",
        "${file}"
    ],
    "problemMatcher": [
        "$vsls"
    ],
    "presentation": {
      "reveal": "never"
    }
},
{
    "label": "Tail",
    "type": "shell",
    "command": "hc2SceneInteract",
    "args": [
        "-g",
        "--tail",
        "-f",
        "${file}"
    ],
    "problemMatcher": [
        "$vsls"
    ],
    "presentation": {
      "reveal": "always"
    }
},

```

add this to your `.vscode/tasks.json` file in your workspace.  Open a lua file, that has FIBARO_GIT_HOOK header included.

Run the task by selecting `Terminal->Run Task ...` `Start`.

If everything works out well the file will be silently started on the FIBARO system.

In order to retrieve the debug messages for the same file, just run the task by selecting `Terminal->Run Task ...` `Tail`.

You will see in output area

```shell
[DEBUG] 14:37:46: Motion detected
[DEBUG] 14:40:22: Motion detected
```

## Upload a script on git commit

### Why should upload a script on git commit

Version control is an important practice in software development and git has become popular alternative. Automatically uploading a scene on a `git commit` command ensures consistency between your production system and the versioned scene.

### Prerequisites

You have created a git version controlled project, where your lua scenes live. You are ok to go, if you have a `.git` directory in the project root. If we call the project root `$(PRJ)`. Then you should have a `$(PRJ)/.git`.

### Steps to integrate into git commit

- Go to `$(PRJ)/.git/hooks`

- Create a file named `pre-commit` with the following contents

```shell
#!/bin/sh
echo "Performing pre-commit-hook ..."
hc2UploadScene -e `pwd` "`git diff --cached --name-only --diff-filter=ACM`"
if [ "$?" -ne "0" ]; then
  echo "uploading scenes failed. Aborting."
  exit 1
fi
exit 0
```

- modify a lua file suitable for upload e.g. `testFile.lua`

- `git commit -m "a commit message" testFile.lua`
  The file `testFile.lua` is being committed and via the `pre-commit` hook uploaded.
  Note: In case you get an error message that `hc2Upload` can not be found, use an absolute path to the executable.

- If you would like to commit a file *without* using the pre-commit-hook use `git commit --no-verify`
  
## FIBARO_GIT_HOOK Header

The FIBARO_GIT_HOOK header has the following structure

```lua
--[[ FIBARO_GIT_HOOK - DO NOT CHANGE AS IT WILL BE DISCARDED
@sceneID=105
@name="SwitchAllLightsOff"
@roomID=0
@autostart=false
@runConfig=MANUAL_ONLY
@maxRunningInstance=2
@type="com.fibaro.luaScene"
@isLua=true
--]]
```

- `@sceneID` contains the ID of the scene as maintained in the FIBARO system. If set to -1 on upload a new scene will be created.
- `@name` contains the scene name as used in the FIBARO system
- `@roomID` contains the ID of the room as maintained in the FIBARO system.
- `@autostart` true or false, indicating whether the script should start at start time
- `@runConfig` TRIGGER_AND_MANUAL, MANUAL_ONLY, DISABLED
- `@maxRunningInstance` indicating the maximum number of instances
- `@type="com.fibaro.luaScene"` should always be "com.fibaro.luaScene", script will not be uploaded if other.
- `@isLua=true` should be true

### Code snipped for FIBARO_GIT_HOOK header

In `examples/fibaro.code-snippets` you can find example on how you could install a code-snipped for the manual insertion of the FIBARO_GIT_HEADER to your workspace.

Copy the file `fibaro.code-snippets` into your `$(workspaceRoot)/.vscode` directory. This install the code snipped, which is enabled when typing `git_hook`.
