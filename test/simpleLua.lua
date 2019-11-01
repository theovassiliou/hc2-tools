
require('test/beginningOfLine');
--[[
%% properties
%% events
311 CentralSceneEvent
%% globals
-]]
-- Some arbirty code, without further meaning
-- SOME LOCAL VARIABLES ---------------------------------------------------
-- In which room is this script runnig
-- require('commented/zb');

local myRoomId = 9; 
-- What should be the default action if script is started manually
local defaultKeyID = 1;
local defautlKeyAttribute = 'Pressed2' ;
dbgLvl = 4;

 require('include/withSomeSpace');

-- Some arbirty code, without further meaning
toggleLightAndSuspendAuto = function (theLight)
    if(fibaro:getGlobal("LightStateOffice") == 'Auto') then 
      fibaro:setGlobal("LightStateOffice", "Manual");
      Debug('green', 'Switching to Manual')
      fibaro:sleep(500); -- Wait a little bit so that lights settle
      switchOn(theLight);
    else
      if(isLightOff(theLight) and  isAllLightOff(myRoomId) ~= true ) then
        switchOn(theLight);
      else
        switchOff(theLight);
        if(isAllLightOff(myRoomId)) then
          fibaro:setGlobal("LightStateOffice", "Auto");
          Debug('green', 'Switching back to Auto')
          switchOn(jT.Buero.BZIndirektesLicht);
          fibaro:sleep(750);
          switchOff(jT.Buero.BZIndirektesLicht);
        end
      end
    end
  end
  
