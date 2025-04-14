@echo off
REM Script to run Ollama integration tests for Local AI Agents

echo Local AI Agents - Ollama Integration Tests
echo ========================================

REM Check if Ollama is installed
where ollama >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo Error: Ollama is not installed or not in PATH
    echo Please install Ollama from https://ollama.com/
    exit /b 1
)

REM Check if Ollama is running
echo Checking if Ollama is running...
curl -s http://localhost:11434 >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo Warning: Ollama is not running
    echo Starting Ollama...
    
    REM Try to start Ollama
    start /b ollama serve
    echo Waiting for Ollama to start...
    timeout /t 5 /nobreak >nul
    
    REM Check again if Ollama is running
    curl -s http://localhost:11434 >nul 2>&1
    if %ERRORLEVEL% NEQ 0 (
        echo Error: Could not start Ollama
        echo Please start Ollama manually and try again
        exit /b 1
    )
) else (
    echo Ollama is running
)

REM Check if the test model is available
echo Checking if test model is available...
set MODEL=llama3.3:latest
curl -s "http://localhost:11434/api/tags" | findstr "%MODEL%" >nul
if %ERRORLEVEL% NEQ 0 (
    echo Warning: Test model %MODEL% is not available
    echo Pulling model %MODEL%...
    ollama pull %MODEL%
    
    if %ERRORLEVEL% NEQ 0 (
        echo Error: Failed to pull model %MODEL%
        echo Tests will use whatever model is available
    ) else (
        echo Model %MODEL% pulled successfully
    )
) else (
    echo Model %MODEL% is available
)

REM Run the tests
echo Running tests...
cd ..
go test -v ./tests/...

REM Check test results
if %ERRORLEVEL% EQU 0 (
    echo All tests passed!
) else (
    echo Some tests failed
)

echo Tests completed
