#ifndef VOICE_BRIDGE_H
#define VOICE_BRIDGE_H

#ifdef __cplusplus
extern "C" {
#include <libavcodec/avcodec.h>
#endif

int decodeToPcmBuffer(char *buffer, int bufferSize,char* targetBuffer,int targetBufferSize,char* outCodec);
int decodeToPcmM4aFile(char *saveFilePath, char* targetBuffer,int targetBufferSize,char* outCodec);

int resampleToPcmBuffer(int src_ch_layout,int src_rate,int src_samplle_fmt,char *buffer, int bufferSize,char* targetBuffer,int targetBufferSize);


#ifdef __cplusplus
}
#endif

#endif