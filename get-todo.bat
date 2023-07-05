@echo off

set URL=http://localhost:8080/todo
set COUNT=100

for /l %%i in (1,1,%COUNT%) do (
    curl %URL%
)
