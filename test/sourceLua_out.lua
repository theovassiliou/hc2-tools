
-- require('lib/ignoreExpand/someLib');
-- require('lib/ignoreExpand/CommentedSomeLib');

-- require('lib/includedLib');
-- LIBRARY BEGIN
-- DO NOT MODIFY THE CODE
-- file includedLib.lua
Here is some other code
bla bla
-- end of file
-- LIBRARY END

-- require('lib/subDir/includeAnotherLib')
-- LIBRARY BEGIN
-- DO NOT MODIFY THE CODE
-- includeAnotherLib in dir subDir
bla bla 
another bla 
-- End of file .
-- LIBRARY END

-- require('lib/subDir/includeAnotherLibCommented')

-- require('test/ignoreExpand/someLib');
