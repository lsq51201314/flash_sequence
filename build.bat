@echo off
cd %cd%
echo 正在生成图标 。。。
echo IDI_ICON1 ICON "main.ico" > main.rc
windres -o main.syso main.rc
echo 正在编译Windows平台 。。。
go build -ldflags "-s -w" -o flash_sequence.exe
echo 编译完成 。。。
pause