# hc2UploadScene

hc2UploadScene uploads a referenced lua-script to the Fibaro HC2 systems.

It is a powerful tool intended to create a round-trip workflow. Plays together with hc2UploadScene and hc2DownloadScene commands.

## Usage

[NOTE: We assume that you have configured access to your Fibaro HC2 system as described in [CONFIGURATION](../../README.md#configuring-your-installation)]

`hc2UploadScene ExampleScene.lua` uploads the file `ExampleScene.lua`.

***

```txt
> hc2UploadScene -h

  Usage: hc2UploadScene [options] <lua-script>

  <lua-script> the file to be uploaded

  Options:
  --log-level, -l    Log level, one of panic, fatal, error, warn or warning, info, debug, trace
                     (default info)
  --cfg-file, -c     The config file to use (default /Users/the/.hc2-tools/config.json)
  --init, -i         Create a default config file as defined by cfg-file, if set. If not set
                     ~/.hc2-tools/config.json will be created.
  --test, -t         Just print information about the contacted HC2 system
  --create-header    Create the FIBARO_GIT_HEADER if set
  --dont-upload, -d  Don't upload the file but print only
  --version, -v      display version
  --help, -h         display help

  HC2 options:
  --user, -u         Username for HC2 authentication
  --password, -p     Password for HC2 authentication
  --url              URL of the Fibaro HC2 system, in the form http://...

  Scene options:
  --scene-id, -s     The sceneId that shall be used. If none given, create a new scene and implies
                     createHeader if header is missing (default -1)
  --room-id, -r      The roomId that shall be used. Implies createHeader if header is missing
                     (default -1)
  --scene-name       The scene name that shall be used. If none given and no header in file, than
                     take filename without file extenion and implies createHeader if header is
                     missing

  Require Expand options:
  --dont-expand      Don't expand the require statements
  --expand-path, -e  Where to search for the included libraries

  Version:
    hc2UploadScene 1.0.0

  Read more:
    github.com/theovassiliou/hc2-tools
```

***

## Examples

`hc2UploadScene -s 55 ExampleScene.lua` uploads the file `ExampleScene.lua` and uses the sceneID `55`. If the sceneID does not exists `hc2UploadScene` returns with error and exit code `1`.

## Library expansion

hc2Uploads enables the support of lua `require()` statements by expanding the library referenced in the `require()`statement inline in the file, before uploading

- `--expand-path` defines the search-path-root where to search for the included libraries.
  For example
    using `require('lib/Debug')` in the lua scene requires the source code inlined in the uploaded lua script.
    With `--expand-path ~/hc2/` the `hc2UploadScene`-tool will look for the required scene at `~/hc2/lib/Debug.lua`
- `--dont-expand` prohibits the expansion, and keeps the lua script as it is.

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
