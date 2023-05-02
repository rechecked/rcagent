05/01/2023 - 1.0.3
==================
- Added cpuPercent and memPercent to processes output (#14)
- Added delta values to the network output when delta=1 parameter passed
- Fixed network check where delta was not applied causing check to not run unless delta was passed
- Fixed perfdata output for warn/crit values missing the ; when only critical is set (#16)

03/29/2023 - 1.0.2
==================
- Fixed issue with windows services not having proper status (#22)
- Fixed empty windows services JSON output to be [] instead of null

03/10/2023 - 1.0.1
==================
- Added Access-Control-Allow-Origin header for CORS requests
- Fixed Content-Type header not being set properly
- Fixed empty plugins JSON output to be [] instead of null

02/14/2023 - 1.0.0
==================
- Initial release