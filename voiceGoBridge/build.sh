#!/bin/bash
g++ -fpic -std=c++11 -shared voiceGoBridge.cpp -L/usr/lib/ -I/usr/include/x86_64-linux-gnu/libavcodec -L. -ldl -lavformat -lavdevice -lavutil -lswresample -o libVoiceGoBridge.so
