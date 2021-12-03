@echo off

rem 设置文件参数
set param=

rem 延迟环境变量扩展, 获取命令行参数，用于将文件拖入bat上使用
setlocal EnableDelayedExpansion 
set "files=%*"
if not "%files%x"=="x" (
	set "param=-files %files: =,%"
)
endlocal & set param=%param%

rem 传参生成配置
%~dp0cfg.exe -c conf.json %param%
