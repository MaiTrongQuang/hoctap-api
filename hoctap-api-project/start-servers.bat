@echo off
cd /d "%~dp0"
echo Starting HocTap API Dashboard...
echo.
echo Starting Go Server with integrated HTML dashboard on port 8080...
start "HocTap Server" cmd /k "cd /d "%~dp0" && go run main.go"

echo.
echo Waiting for server to start...
timeout /t 3 /nobreak >nul

echo.
echo ================================
echo   HocTap API Dashboard Ready!
echo ================================
echo.
echo Server:         http://localhost:8080
echo Web Dashboard:  http://localhost:8080
echo API Endpoints:  http://localhost:8080/api/users
echo Health Check:   http://localhost:8080/health
echo.
echo Press any key to open the dashboard in your browser...
pause >nul

start http://localhost:8080
