@echo off
set VERSION=1.0.0
set TARGET=%cd:~% 

for /f "delims=" %%i in ('git describe --always') do set DESCRIBE=%%i
for /f "delims=" %%i in ('git rev-list --count --first-parent HEAD') do set REVERSION=%%i
set VERSION=%VERSION% %REVERSION%.%DESCRIBE%
for /f "tokens=1-3 delims= " %%a in ('date /t') do set BUILDDT=%%a %time%

go build -ldflags "-X 'main.VERSION=%VERSION%' -X 'main.BUILDDT=%BUILDDT%'" -o %TARGET% .

