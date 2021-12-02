@echo off

rem 获取命令行参数，用于将文件拖入bat上使用
set "files=%*"

rem 如果直接运行
if "%files%x"=="x" (
	cfg.exe -c ./conf.json
	goto finish
) 
	
rem 用','分割多个命令行传入文件
set files=%files: =,%

rem 将命令行拖入文件以files参数传入
%~dp0cfg.exe -c conf.json -files %files%
goto finish

:finish
pause