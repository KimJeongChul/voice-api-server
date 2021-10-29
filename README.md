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