package voiceExtension

//#cgo CFLAGS: -I/home/dev/go/src/voice-api-server/voiceGoBridge
//#cgo LDFLAGS: -L/usr/lib/x86_64-linux-gnu -lswresample -lavformat -lavcodec -lavutil -lVoiceGoBridge
//#include <stdlib.h>
//#include "/home/dev/go/src/voice-api-server/voiceGoBridge/voiceGoBridge.h"
import "C"

import (
	"bytes"
	"log"
	"unsafe"
	"voice-api-server/common"

	"github.com/takuyaohashi/go-wav"
)

var voiceChannelMap map[int]*VoiceObj

type VoiceObj struct {
	Session	*common.SessionObj
}

func init() {
	voiceChannelMap = make(map[int]*VoiceObj)
}

//export C2GoLog
func C2GoLog(context C.int, funcName *C.char, logText *C.char) {
	gCtx := int(context)
	gFuncName := C.GoString(funcName)
	gLogText := C.GoString(logText)
	gSttCtx := voiceChannelMap[gCtx]
	if gSttCtx != nil {
		log.Println("[I] (", gSttCtx.Session.TransactionId, ") < C-Function::"+gFuncName, "> Msg:", gLogText)
	} else {
		log.Println("[I] ( UNDEFINED ) < C-Function::"+gFuncName, "> Msg:", gLogText)
	}
}

// GetPcmFromEncoded Decoding audio(mp3, flac, ogg) to pcm
func GetPcmFromEncoded(inBuffer []byte, inOutBuffer []byte) (string, int) {
	gReturnCodec := "pcm"
	cReturnCodec := C.CString(gReturnCodec)
	cInBufferLength := C.int(len(inBuffer))
	cOutBufferLength := C.int(len(inOutBuffer))
	cInBuffer := (*C.char)(unsafe.Pointer(&inBuffer[0]))
	cOutBuffer := (*C.char)(unsafe.Pointer(&inOutBuffer[0]))
	gRet := int(C.decodeToPcmBuffer(cInBuffer, cInBufferLength, cOutBuffer, cOutBufferLength, cReturnCodec))
	gReturnCodec = C.GoString(cReturnCodec)
	C.free(unsafe.Pointer(cReturnCodec))
	return gReturnCodec, gRet
}

// Decoding wav audio to pcm
func GetWavToPcm(inBuffer []byte, inOutBuffer []byte) (rc int) {
	cInBuffer := (*C.char)(unsafe.Pointer(&inBuffer[0]))
	cInBufferLength := C.int(len(inBuffer))
	cOutBuffer := (*C.char)(unsafe.Pointer(&inOutBuffer[0]))
	cOutBufferLength := C.int(len(inOutBuffer))
	
	// Parse WAV Header 
	var b bytes.Buffer
	b.Write(inBuffer[0:50])
	wavFile := bytes.NewReader(b.Bytes())
	wavReader := wav.NewReader(wavFile)
	wavReader.Parse()
	
	// Get Source SampleRate
	srcSampleRate := int(wavReader.GetSampleRate())
	srcChannels := int(wavReader.GetNumChannels())
	srcBitPerSample := int(wavReader.GetBitsPerSample())
	gSampleFmt := 1
	switch srcBitPerSample {
	case 16:
	case 32:
		gSampleFmt = 2
	}
	cSampleFmt := C.int(gSampleFmt)
	rc = int(C.resampleToPcmBuffer(C.int(srcChannels), C.int(srcSampleRate), C.int(cSampleFmt), cInBuffer, cInBufferLength, cOutBuffer, cOutBufferLength))
	return
}
