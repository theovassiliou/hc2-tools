# hc2DownloadScene

Download one or all scenes from a fibaroHC2 system. It can be choosen whether the file(s) should be located in directory of choice or locally.

It is a powerful tool intended to create a round-trip workflow. Plays together with `hc2UploadScene` and `hc2SceneInteract` commands.

## Usage

[NOTE: We assume that you have configured access to your Fibaro HC2 system as described in [CONFIGURATION](../../README.md#configuring-your-installation)]

`hc2DownloadScene --dir tmp --scene-id 55` downloads the script with the Fibaro HC2 sceneId `55` and saves it in directory `tmp`.

The file name is constructed from the scene name as defined in HC2 and suffixed with `.lua`

If stored in a directory, sub directories will be created in the form `./zoneName/roomName/`

If the directories do not exist they will be created.

```shell

 hc2DownloadScene -h

  Usage: hc2DownloadScene [options]

  Options:
  --log-level, -l  Log level, one of panic, fatal, error, warn or warning, info, debug, trace
                   (default info)
  --cfg-file, -c   The config file to use (default /Users/the/.hc2-tools/config.json)
  --init, -i       Create a default config file as defined by cfg-file, if set. If not set
                   ~/.hc2-tools/config.json will be created.
  --test, -t       Just print information about the contacted HC2 system
  --version, -v    display version
  --help, -h       display help

  HC2 options:
  --user, -u       Username for HC2 authentication
  --password, -p   Password for HC2 authentication
  --url            URL of the Fibaro HC2 system, in the form http://...

  Scene options:
  --create-header  If set create the FIBARO_GIT_HEADER if none present (default true)
  --scene-id, -s   The sceneId that shall be used. If none given, all scenes will be downloaded.
                   (default -1)
  --dir, -d        Where to search for the included libraries (default ./download)

  Version:
    hc2DownloadScene 1.0.0

  Read more:
    github.com/theovassiliou/hc2-tools

```

***

## Examples

`hc2DownloadScene -u admin -p theAdminPassword --url http://192.166.22.10 -i -t` tests the connection to your Fibaro HC2 system and creates a config-file in `~/.hc2-tool/*` so that you do not have to reenter the information in subsequent to the hc2-tools calls.

`hc2DownloadScene --dir tmp --scene-id 55` downloads the script with the Fibaro HC2 sceneId `55` and saves it in directory `tmp`.
If a scene with scenedID `55` does not exists program will return with exit code 1

`hc2DownloadScene -dir tmp` downloads all scripts from the Fibaro HC2 and saves it in directory `tmp`.

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
...
