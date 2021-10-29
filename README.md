# voice-api-server

Golang voice api server 
 - Preprocessing audio converting(encoding and decoding) using FFmpeg 
 - Serving Speech(Speech To Text) and Synthesis(Text To Speech) service
 - Supported Rest API

##
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

### Voice Recognize
Audio Format
 - wav, mp3, m4a(aac), Raw pcm, ogg, flac

Request
```
$ curl -F metadata="{\"encoding\":\"wav\",\"targetLanguage\":\"ko\", \"sttMode\":1, \"encodingOpt\":{\"channel\":1, \"sampleRate\": 16000, \"sampleFmt\": \"S16LE\"}}" -F media=@test40.wav https://[SERVICE_URL]/v1/speech
```

Response
```bash
--2c941089d64863ff88d066d3a9ff37ad8cf468eca4fe81abe4223e2aca1b
Content-Disposition: form-data; name="voiceResult"

{"resultType":"start"}
--2c941089d64863ff88d066d3a9ff37ad8cf468eca4fe81abe4223e2aca1b
Content-Disposition: form-data; name="voiceResult"

{"resultType":"full","speechResult":{"text":"RECOGNIZE FROM YOUR STT CORE ENGIN SERVER","startTime":0.1,"endTime":40}}
--2c941089d64863ff88d066d3a9ff37ad8cf468eca4fe81abe4223e2aca1b
Content-Disposition: form-data; name="voiceResult"

{"resultType":"end","speechInfo":{"reqFileSize":1280078,"transCodec":"wav","convFileSize":1280078,"speechInputTime":40}}
--2c941089d64863ff88d066d3a9ff37ad8cf468eca4fe81abe4223e2aca1b--
```
