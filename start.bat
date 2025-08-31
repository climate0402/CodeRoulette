@echo off
echo Starting CodeRoulette with Docker Compose...
docker-compose up -d
echo.
echo Services started! Access the application at:
echo Frontend: http://localhost:3000
echo Backend API: http://localhost:8080
echo.
echo To stop services, run: docker-compose down
pause
