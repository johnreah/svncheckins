@echo off

for %%r in (
	https://<svn-cloud-host>/projectA
	https://<svn-cloud-host>/projectB
	https://<svn-cloud-host>/projectC
	
) do ^
echo %%r & ^
echo. & ^
svn info %%r/trunk & ^
svn info %%r/branches & ^
echo ----------------------------------------------------------------------------- & ^
echo.

