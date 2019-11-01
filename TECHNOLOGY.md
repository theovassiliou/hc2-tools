# The technology behind the scenes

The tools are build around the fact, that we believe it is wise to manage the Fibaro scenes, which are in fact lua-scripts outside of the Fibaro HC2 system. For us the main reason was, that the editing capabilities of the web-interface where limited. No code completion, etc. was available. Also testing, was kind of impossible.

As software developers tend to use a version control for their developments the need for a file based development environment was getting even more obvious. Luckily the Fibaro HC2 exposes as (not so well) documented REST API. If you are interested to learn more about it take a lookg at `http://enterYourFibaroIP/docs`if you have a Fibaro HC2.

## The Idea

So we build the hc2-tools around the idea, that every scene is represented as a file on the disk. As the scenes are in fact lua-scrpits we are using for lua-scripts the suffix `.lua`. We found it helpfull that we follow the structuring of the Fibaro HC2 which associates each scene to a room, and each room to a section of your house.

While a room can have a name, it is referenced internally via a `roomID` the same way a scene has a name, but it is referenced internally via a `sceneID`. In order to relate a lua-script file on your file system to a particular scene in the Fibaro HC2 system, it is required to store the relevant information somewhere. We decided to store the information in the file itself, in a comment. We call it the `FIBARO_GIT_HOOK`-header.  

## FIBARO_GIT_HOOK-header

The concept of the `FIBARO_GIT_HOOK` is best introduced with an example.

```comment
--[[ FIBARO_GIT_HOOK - DO NOT CHANGE AS IT WILL BE DISCARDED
@sceneID=133
@name="FibaroDimmerExample"
@roomID=0
@autostart=false
@runConfig=TRIGGER_AND_MANUAL
@maxRunningInstance=2
@type="com.fibaro.luaScene"
@isLua=true
--]]
```

If you download your scenes via the `hc2Download` command the fields will be filled with the actual values that they have in your Fibaro HC2 system, so you do not really have to care, when you edit existing scenes.

If you write new scenes, or you would like to update the values with the next upload here are the options.

- `@sceneID: int` Integer containing the sceneID as used in the Fibaro HC2 system. -1 if a new scene should be created on upload
- `@name: string*` String containing the name of the scene. Value will be used/changed on upload
- `@roomID: int*`  Integer containing the roomID identifying the room. `0` indicated the "undefined" room.
- `@autostart: boolean*`  True or false, indicating whether the scene should be started on boot-time.
- `@runConfig: enum*` One of `TRIGGER_AND_MANUAL` `MANUAL_ONLY` or `DISABLED`
- `@maxRunningInstance: int*` The maximum number of instances running simultanously
- `@type: string*` Should be "com.fibaro.luaScene"
- `@isLua: boolean` indicates whether a scene is a lua-script. Should be true.

Fields marked with `*` can be updated. All other should be kept unchanged.

## How to identify the Fibaro HC2 internal IDs

Sometimes it is necessary to identify a particular sceneID or roomID within the Fibaro HC2 system. Go to you Fibaro HC2 system, select the scene or room and take a look at your URL in your browser. It should look like `http://192.163.174.22/fibaro/de/scenes/edit.html?id=148#bookmark-advanced` for a scene or like `http://192.163.174.22/fibaro/de/rooms/edit.html?id=12` for a room.

The sceneID for the above scene is `148` while the roomID in this example would be `12`.

## Libaries for Fibaro HC2 systems, mimicking the lua `require` statement

Unfortunately, the Fibaro HC2 system does *not* support the lua `require` statement. So writing reusable code is a "copy and paste" discpline.

For example: You have debug function, that you would like to use in each script.

```lua
Debug = function ( color, message )
  local print = false;
  if(dbgLvl ~= nil and assert(type(dbgLvl) == "number", "dbgLvl expects a number (1..4)")) then

    if(dbgLvl >= 1 and color == "red") then
      print = true;
    elseif (dbgLvl >= 2 and color == "green") then
      print = true;
    elseif (dbgLvl >= 3 and color == "blue") then
      print = true;
    elseif (dbgLvl >= 4 and color == "orange") then
      print = true;
    elseif (color ~= "red" and color ~= "green" and color ~= "blue" and color ~= "orange") then
      print =true ;
    end
  else  
    print = true;
  end

  if (print) then
    fibaro:debug(string.format('<%s style="color:%s;">%s', "span", color, message, "span"))
  end
end
```

In traditionall lua programming you would define a libary by saving this code for example in a file `Debug.lua`, preferably organized in a directory named for example `library/`.

A resulting Fibaro HC2 scene-script would then look like this

`CentralSceneExample.lua`

```lua
--[[
%% properties
%% events
589 CentralSceneEvent
%% globals
--]]

require('library/Debug')

Debug("green", "CentralSceneEvent triggered")

--[[ FIBARO_GIT_HOOK - DO NOT CHANGE AS IT WILL BE DISCARDED
@sceneID=22
@name="CentralSceneExample"
@roomID=0
@autostart=false
@runConfig=TRIGGER_AND_MANUAL
@maxRunningInstance=2
@type="com.fibaro.luaScene"
@isLua=true
--]]
```

The parts of the scene are:

1. the header that registeres a `CentralSceneEvent`. This is Fibaro HC2 proprietary.
2. a require statement that includes the library `Debug` in the sub-directory `library`
3. `Debug("green", ...)` that issues a Debug message
4. `FIBARO_GIT_HOOK` header

But as said, unfortunately the Fibaro HC2 system does not support the usage of the `require` statement.

So you typically copy and paste the "library code" into each scene.

While it is not only very cumbersome to copy this into every scene, it makes also maintenance hard. Imagine you change something in the debug function, and you have to manually replace the code in each script.

Using the hc2-tools we enable lua-script developers to use the require statement, by expanding each require statement before uploading a script to the Fibaro HC2 system.

Example:

`hc2Upload --expand-path ~/fibaroCoding/ CentralSceneExample.lua`

`--expand-path` defines a path to the directory containing the lua-libraries. `CentralSceneExample.lua` indicates the lua-script to be uploaded. The sceneID in this file is set to `22`. `hc2UploadScene` expands the `require('library/Debug')` statement by including the file content of `~/fibaroCoding/library/Debug.lua` for the `require`- statement and uploads the script to the Fibaro HC2 system.

The result looks as follows:

```lua
--[[
%% properties
%% events
589 CentralSceneEvent
%% globals
--]]

--^ require('library/Debug')
-- LIBRARY BEGIN -------------------------
-- DO NOT MODIFY THE CODE
Debug = function ( color, message )
  local print = false;
  if(dbgLvl ~= nil and assert(type(dbgLvl) == "number", "dbgLvl expects a number (1..4)")) then

    if(dbgLvl >= 1 and color == "red") then
      print = true;
    elseif (dbgLvl >= 2 and color == "green") then
      print = true;
    elseif (dbgLvl >= 3 and color == "blue") then
      print = true;
    elseif (dbgLvl >= 4 and color == "orange") then
      print = true;
    elseif (color ~= "red" and color ~= "green" and color ~= "blue" and color ~= "orange") then
      print =true ;
    end
  else  
    print = true;
  end

  if (print) then
    fibaro:debug(string.format('<%s style="color:%s;">%s', "span", color, message, "span"))
  end
end
-- LIBRARY END -------------------------

Debug("green", "CentralSceneEvent triggered")

--[[ FIBARO_GIT_HOOK - DO NOT CHANGE AS IT WILL BE DISCARDED
@sceneID=22
@name="CentralSceneExample"
@roomID=0
@autostart=false
@runConfig=TRIGGER_AND_MANUAL
@maxRunningInstance=2
@type="com.fibaro.luaScene"
@isLua=true
--]]
```

The `require('library/Debug')` has been commented as `--^ require('library/Debug')`, the file content has been inserted and surounded by the lines

```lua
-- LIBRARY BEGIN -------------------------
-- DO NOT MODIFY THE CODE
...
-- LIBRARY END -------------------------
```

Setting the option `--dont-expand` for the `hc2UploadScene` command bypasses the expansion mechanis and leaves the `require` statements untouched before uploading.

Evaluating the effect of the options *without* uploading a file can be achieved using the `-d` (`--dont-upload`)
