@echo off
setlocal

set BINARY=mcp-inspect.exe
set CMD=.\cmd\inspect

if "%1"=="" goto build_all
if "%1"=="build" goto build_all
if "%1"=="build-main" goto build_main
if "%1"=="build-mocks" goto build_mocks
if "%1"=="demo" goto demo
if "%1"=="demo-json" goto demo_json
if "%1"=="clean" goto clean
echo Unknown target: %1
echo Usage: make.bat [build^|build-main^|build-mocks^|demo^|demo-json^|clean]
exit /b 1

:build_all
call %0 build-main
if errorlevel 1 exit /b 1
call %0 build-mocks
goto end

:build_main
echo Building mcp-inspect...
go build -o %BINARY% %CMD%
if errorlevel 1 (echo Build failed & exit /b 1)
echo Done: %BINARY%
goto end

:build_mocks
echo Building mock servers...
go build -o testdata\mock-servers\safe\mock-server.exe       .\testdata\mock-servers\safe
go build -o testdata\mock-servers\destructive\mock-server.exe  .\testdata\mock-servers\destructive
go build -o testdata\mock-servers\hidden-tools\mock-server.exe .\testdata\mock-servers\hidden-tools
if errorlevel 1 (echo Mock build failed & exit /b 1)
echo Done: mock servers
goto end

:demo
call %0 build
if errorlevel 1 exit /b 1
echo Running demo...
%BINARY% --config testdata\demo-config.json --output demo-report.html
goto end

:demo_json
call %0 build
if errorlevel 1 exit /b 1
%BINARY% --config testdata\demo-config.json --format json
goto end

:clean
echo Cleaning...
if exist %BINARY% del %BINARY%
if exist demo-report.html del demo-report.html
if exist mcp-report.html del mcp-report.html
if exist testdata\mock-servers\safe\mock-server.exe        del testdata\mock-servers\safe\mock-server.exe
if exist testdata\mock-servers\destructive\mock-server.exe   del testdata\mock-servers\destructive\mock-server.exe
if exist testdata\mock-servers\hidden-tools\mock-server.exe  del testdata\mock-servers\hidden-tools\mock-server.exe
echo Done.
goto end

:end
endlocal
