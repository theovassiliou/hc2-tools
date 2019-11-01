--[[
%% properties
%% globals
--]]

--[[
-----------------------------------------------------------------------------
-- SECURE TURN OFF LIGHTS-- Schlaffzimmer
-----------------------------------------------------------------------------

-- LIST OF GLOBAL VARIABLES NEEDED ------------------------------------------

NONE

--]]

-- SOME LOCAL VARIABLES -----------------------------------------------------
local roomId = 5;


-- USING HOMETABLE IS RECOMMENDED -------------------------------------------
--local jT = json.decode(fibaro:getGlobalValue("HomeTable"));


-- SOME LOCAL VARIABLES -----------------------------------------------------

if (fibaro:countScenes() > 1) then
  fibaro:abort();
end

--------------------------------------------------------------------------------
-- The classical debug function. Prints a *color*ed message.
-- @function Debug
-- The classical debug function. Prints a *color*ed message.
-- @param color color One of, *red*, *blue*, *green*, *yellow*, *orange*, *white*
-- @param message message the message to print.
-- @return nothing
-------------------------------------------------------------------------------
Debug = function ( color, message )
  fibaro:debug(string.format('<%s style="color:%s;">%s', "span", color, message, "span"))
end

--- WANT TO SHOW DEBUG INFO, SET BELOW VALUE TO TRUE
debug = true;


local startSource = fibaro:getSourceTrigger();
-- Do not want to do anything if manually triggered
-- we are reacting on power or value change
if(startSource == nil) then
  return
else
  Debug('red', "startSource: " .. startSource["type"]);
--  if (( startSource["type"] == "other")) then return end
end


-- usefull functions for wall plugs

function main()
  fibaro:startScene(105,{roomId});
end

main()
