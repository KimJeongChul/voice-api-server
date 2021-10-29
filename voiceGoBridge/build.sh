#!/bin/bash
g++ -fpic -std=c++11 -shared voiceGoBridge.cpp -L/usr/lib/ -I/usr/include/glib-2.0/ -I/usr/lib/x86_64-linux-gnu/glib-2.0/include  -L. -ldl -lswresample -o libVoiceGoBridge.so