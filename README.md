# voice-api-server

Golang voice api server 
 - Preprocessing audio converting(encoding and decoding) using FFmpeg 
 - Serving Speech(Speech To Text) and Synthesis(Text To Speech) service
 - Supported Rest API

![image](https://user-images.githubusercontent.com/10591350/139374435-d192b956-9b97-4314-acc9-fbb9e1a0319a.png)

| API                             | PATH          | Method          |
|---------------------------------|---------------|-----------------|
| Voice Recognize(Speech To Text) | /v1/speech    | handleSpeech    |
| Voice Synthesis(Text To Speech) | /v1/synthesis | handleSynthesis |

## Prequisition
Install Library
```bash
# apt-get install libresample-dev libavcodec-dev libavformat-dev libglib2.0-dev libswscale-dev libmp3lame-dev nasm yasm -y
```

Build FFmpeg
- https://github.com/FFmpeg/FFmpeg

## CGO
voiceExtension/voiceExtension.go
```golang
//#cgo CFLAGS: -I[YOUR_PATH]/voice-api-server/voiceGoBridge
//#include "[YOUR_PATH]/voice-api-server/voiceGoBridge/voiceGoBridge.h"
```
``` bash
$ go tool cgo voiceExtension/voiceExtension.go 
```

## VoiceGoBridge
Build
```bash
$ cd voiceGoBridge
$ ./build.sh
# cp libVoiceGoBridge.so /usr/lib/
```

## Server Config
```json
{
    "listenPort": [SERVER_LISTEN_PORT],
    "ssl":[SSL_OPTION],
    "certPemPath":[CERT_PATH],
    "keyPemPath": [KEY_PATH],
    "rcvAudioSavePath": [RECEIVE_AUDIO_SAVE_PATH],
    "pcmSavePath": [PCM_SAVE_PATH],
    "speechResultPath": [SPEECH_RESULT_PATH],
    "logPath": [LOG_PATH],
    "logPeriod": [LOG_ROTATE_PERIOD_DAY]
}

example:
{
    "listenPort": "9096",
    "ssl":1,
    "certPemPath":"/data/voiceApiServer/cert.pem",
    "keyPemPath":"/data/voiceApiServer/cert/key.pem",
    "rcvAudioSavePath": "/data/voiceApiServer/audioReceived",
    "pcmSavePath": "/data/voiceApiServer/pcmSaved",
    "speechResultPath": "/data/voiceApiServer/speechResult",
    "logPath": "/data/voiceApiServer/log",
    "logLevel": "debug",
    "logPeriod": 60
}   
```

## Run
```bash
$ go build
$ ./voice-api-server
```

