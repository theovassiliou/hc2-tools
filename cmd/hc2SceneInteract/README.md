# hc2SceneInteract

Start, top, enable or disable scenes in the Fibaro HC2 system. In addition you can retrieve the debug messages for an executing scene of the Fibaro HC2 system.

It is a powerful tool intended to create a round-trip workflow. Plays together with hc2UploadScene and hc2DownloadScene commands.

## Usage

[NOTE: We assume that you have configured access to your Fibaro HC2 system as described in [CONFIGURATION](../../README.md#configuring-your-installation)]

`hc2SceneInteract --action start -sceneId 55` triggers execution of scene 55.

***

```shell

hc2SceneInteract -h

  Usage: hc2SceneInteract [options]

  Options:
  --action, -a     Triggers a scene action for sceneID. One of start, stop, enable, disable.
  --get-debug, -g  Retrieve after starting the action the debug messages, while respecting the
                   tail-flag. Ignored with action enable or disable
  --init, -i       Create a default config file as defined by cfg-file, if set. If not set
                   ~/.hc2-tools/config.json will be created.
  --test, -t       Just print information about the contacted HC2 system
  --cfg-file, -c   The config file to use (default /Users/the/.hc2-tools/config.json)
  --log-level, -l  Log level, one of panic, fatal, error, warn or warning, info, debug, trace
                   (default info)
  --version, -v    display version
  --help, -h       display help

  HC2 options:
  --user, -u       Username for HC2 authentication
  --password, -p   Password for HC2 authentication
  --url            URL of the Fibaro HC2 system, in the form http://...

  Generic command options:
  --scene-id, -s   The sceneId that shall be used.
  --file, -f       sceneID is taken from <lua-script-file> with fibaro header. sceneID flag is
                   ignored.
  --tail           The -t option causes get-debug to not stop when all debug messages are read, but
                   rather to wait for additional data to be appended to the input.

  Version:
    hc2SceneInteract 1.0.0

  Read more:
    github.com/theovassiliou/hc2-tools

```

***

## Examples

`hc2SceneInteract -a start -s 55 -g --tail` starts scene 55, retrieves all debug messages so far, and waits for more. Not terminating. Hit Control-C to terminate execution.

`hc2SceneInteract -s 55 -g` retrieves all debug messages from scene 55.

`hc2SceneInteract -f myLuaScene.lua -a start` reads the file  myLuaScene.lua, parses the `FIBARO_GIT_HEADER`-headers, uses the scene-id found there, and start the referenced scene in the Fibaro HC2 system. Note that the lua-file is *not* uploaded to the system.

```shell
hc2UploadScene myLuaScene.lua
hc2SceneInteract -f myLuaScene.lua -a start -g
```

first uploads the `myLuaScene.lua` file. The second command then starts the scene and then get the debug messages.

## config-file

The file has the following structure

```json
{
    "url":"http://hc.2.ip.address",
    "username":"admin",
    "password":"theAdminPassword"
}
```

Please note that you do not have to use the admin user, but any user who has access to desired scenes.
