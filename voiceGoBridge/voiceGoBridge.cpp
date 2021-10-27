#include <iostream>
#include <fstream>
#include <algorithm>
#include <string.h>
#include <mutex>
#include <string>
#include <thread>
#include <map>

#define CODEC_CAP_DELAY 0x0020

#ifdef __cplusplus
extern "C" {
#include "/mnt/dev/gopath/src/voice-api-server/_obj/_cgo_export.h"
#include <unistd.h>
#include <libavcodec/avcodec.h>
#include <libavformat/avformat.h>
#include <libavformat/avio.h>
#include <libswscale/swscale.h>
#include <libavutil/imgutils.h>
#include <libavutil/samplefmt.h>
#include <libavutil/timestamp.h>
#include <libswresample/swresample.h>
#include <libavutil/opt.h>
#endif

struct buffer_data {
    uint8_t *ptr;
    size_t size;
    int context;
};

// Log C Logging function
void Log(int inCtx, const char *inFunc, const char *inMsg) {
    // Call Go Logging function
    C2GoLog(inCtx,(char*)inFunc,(char*)inMsg);
}

// decodeToPcmBuffer Decoding pcm buffer (input audio format: mp3, flac, ogg)
int decodeToPcmBuffer(char *buffer, int bufferSize,char* targetBuffer,int targetBufferSize,char* outCodec) {
    int inCtx=-1;
    int rc = 0;
    int data_size;
    uint8_t *avio_ctx_buffer = NULL;
    size_t avio_ctx_buffer_size = 4096;

    AVIOContext *avio_ctx = NULL;
    struct SwrContext *swr_ctx;
    struct buffer_data bd = {0};
    bd.ptr = (uint8_t *) buffer;
    bd.size = bufferSize;

    av_register_all();
    AVFrame *frame = av_frame_alloc();
    if (!frame) {
        Log(inCtx, "decodeToPcmBuffer", "AVFrame Allocation Error");
        return -2;
    }

    AVFormatContext *formatContext = NULL;
    if (!(formatContext = avformat_alloc_context())) {
        Log(inCtx, "decodeToPcmBuffer", "AVFormatContext Allocation Error");
        av_free(frame);
        return -3;
    }

    avio_ctx_buffer = (uint8_t *) av_malloc(avio_ctx_buffer_size);
    if (!avio_ctx_buffer) {
        Log(inCtx, "decodeToPcmBuffer", "avio_ctx_buffer Allocation Error");
        av_free(frame);
        if (formatContext!=NULL) av_freep(formatContext);
        return -4;
    }

    avio_ctx = avio_alloc_context(
            avio_ctx_buffer,
            avio_ctx_buffer_size,
            0,
            &bd,
            [](void *opaque, uint8_t *buf, int buf_size) {
                struct buffer_data *bd = (struct buffer_data *) opaque;
                buf_size = FFMIN(buf_size, bd->size);
                memcpy(buf, bd->ptr, buf_size);
                bd->ptr += buf_size;
                bd->size -= buf_size;
                return buf_size;
            }, NULL, NULL
    );

    if (!avio_ctx) {
        Log(inCtx, "decodeToPcmBuffer", "avio_ctx Allocation Error");
        av_free(frame);
        if (formatContext!=NULL) av_freep(formatContext);
        av_freep(avio_ctx_buffer);
        return -5;
    }

    formatContext->pb = avio_ctx;

    //Open AVContext
    if (avformat_open_input(&formatContext, NULL, NULL, NULL) != 0) {
        Log(inCtx, "decodeToPcmBuffer", "avformat open input Error");
        av_free(frame);
        if (formatContext!=NULL) av_freep(formatContext);
        if (avio_ctx) {
            av_freep(avio_ctx->buffer);
            av_freep(avio_ctx);
        }
        return -6;
    }

    //Find Stream Info
    if (avformat_find_stream_info(formatContext, NULL) < 0) {
        Log(inCtx, "decodeToPcmBuffer", "avformat find stream info Error");
        av_free(frame);
        if (formatContext!=NULL) av_freep(formatContext);
        if (avio_ctx) {
            av_freep(avio_ctx->buffer);
            av_freep(avio_ctx);
        }
        return -7;
    }

    std::cout<<"Format context: "<<formatContext<<std::endl;

    //Find audio Stream
    AVCodec *cdc = nullptr;
    int streamIndex = av_find_best_stream(formatContext, AVMEDIA_TYPE_AUDIO, -1, -1, &cdc, 0);
    if (streamIndex < 0) {
        Log(inCtx, "decodeToPcmBuffer", "find audio stream info Error");
        avformat_close_input(&formatContext);
        av_free(frame);
        if (formatContext!=NULL) av_freep(formatContext);
        if (avio_ctx) {
            av_freep(avio_ctx->buffer);
            av_freep(avio_ctx);
        }
        return -8;
    }

    AVStream *audioStream = formatContext->streams[streamIndex];
    AVCodecContext *codecContext = audioStream->codec;
    codecContext->codec = cdc;

    //Open Codec
    if (avcodec_open2(codecContext, codecContext->codec, NULL) != 0) {
        Log(inCtx, "decodeToPcmBuffer", "avcodec open 2 Error");
        avformat_close_input(&formatContext);
        av_free(frame);
        if (formatContext!=NULL) av_freep(formatContext);
        if (avio_ctx) {
            av_freep(avio_ctx->buffer);
            av_freep(avio_ctx);
        }
        return -9;
    }

    // Record codec info 
    memcpy(outCodec,avcodec_get_name(cdc->id),3);

    std::string sr = std::to_string(codecContext->sample_rate);
    char const *s_rate = sr.c_str();

    Log(inCtx, "decodeToPcmBuffer", avcodec_get_name(cdc->id));
    Log(inCtx, "decodeToPcmBuffer", s_rate);

    //Prepare Resampler
    swr_ctx = swr_alloc_set_opts(NULL,
                                 AV_CH_LAYOUT_MONO,
                                 AV_SAMPLE_FMT_S16,
                                 16000,
                                 codecContext->channel_layout,
                                 codecContext->sample_fmt,
                                 codecContext->sample_rate,
                                 0,
                                 NULL
    );

    //Requested input sample format 31 is invalid

    if (!swr_ctx) {
        Log(inCtx, "decodeToPcmBuffer", "swr allocation Error");
        avformat_close_input(&formatContext);
        av_free(frame);
        if (formatContext!=NULL) av_freep(formatContext);
        if (avio_ctx) {
            av_freep(avio_ctx->buffer);
            av_freep(avio_ctx);
        }
        return -10;
    }

    if (swr_init(swr_ctx) < 0) {
        Log(inCtx, "decodeToPcmBuffer", "swr init Error");
        av_freep(swr_ctx);
        avformat_close_input(&formatContext);
        av_free(frame);
        if (formatContext!=NULL) av_freep(formatContext);
        if (avio_ctx) {
            av_freep(avio_ctx->buffer);
            av_freep(avio_ctx);
        }
        return -11;
    }

    AVPacket readingPacket;
    av_init_packet(&readingPacket);

    // Read Packet
    bool loopFlag=true;
    while (av_read_frame(formatContext, &readingPacket) == 0 && loopFlag) {
        if (readingPacket.stream_index == audioStream->index) {
            AVPacket decodingPacket = readingPacket;
            int dst_bufsize;
            uint8_t **dst_data = NULL;
            while (decodingPacket.size > 0) {
                int gotFrame = 0;
                int result = avcodec_decode_audio4(codecContext, frame, &gotFrame, &decodingPacket);
                if (result >= 0 && gotFrame) {
                    decodingPacket.size -= result;
                    decodingPacket.data += result;
                    data_size = av_get_bytes_per_sample(codecContext->sample_fmt);
                    
                    int ret;
                    int dst_linesize;

                    // Calculate expected out_num_sample
                    int out_num_samples = av_rescale_rnd(
                            swr_get_delay(swr_ctx, codecContext->sample_rate) + frame->nb_samples,
                            16000,
                            codecContext->sample_rate,
                            AV_ROUND_UP
                    );

                    ret = av_samples_alloc_array_and_samples(
                            &dst_data,
                            &dst_linesize,
                            1,
                            out_num_samples,
                            AV_SAMPLE_FMT_S16,
                            0
                    );

                    ret = swr_convert(
                            swr_ctx,
                            dst_data,
                            out_num_samples,
                            (const uint8_t **) &frame->data[0],
                            frame->nb_samples
                    );

                    dst_bufsize = av_samples_get_buffer_size(&dst_linesize, 1, ret, AV_SAMPLE_FMT_S16, 1);

                    if(dst_bufsize>0) {
                        if ((rc + dst_bufsize) <= targetBufferSize) {
                            memcpy(targetBuffer + rc, (char *) dst_data[0], dst_bufsize);
                            rc = rc + dst_bufsize;
                        } else {
                            loopFlag = false;
                        }
                    } else {
                        char errs[AV_ERROR_MAX_STRING_SIZE+20];
                        av_make_error_string(errs, AV_ERROR_MAX_STRING_SIZE, dst_bufsize);
                        Log(inCtx, "decodeToPcmBuffer", errs);
                        //loopFlag=false;
                    }
                } else {
                    decodingPacket.size = 0;
                    decodingPacket.data = nullptr;
                }
            }
        }
            (&readingPacket);
    }

    if (codecContext->codec->capabilities & CODEC_CAP_DELAY) {
        av_init_packet(&readingPacket);
        int gotFrame = 0;
    }

    cleanUp:
    Log(inCtx, "decodeToPcmBuffer", "CleanUp AV-Related Objs.");
    if (swr_ctx) {
        swr_free(&swr_ctx);
    }
    av_free(frame);
    if (avio_ctx) {
        av_freep(&avio_ctx->buffer);
        av_freep(&avio_ctx);
    }
    avcodec_close(codecContext);
    avformat_close_input(&formatContext);
    return rc;
}

// resampleToPcmBuffer Resample pcm
int resampleToPcmBuffer(int src_ch_layout,int src_rate,int src_sample_fmt,char *buffer, int bufferSize,char* targetBuffer,int targetBufferSize){
    int ret;
    int dst_linesize;
    std::string logStr="";

    int64_t t_src_ch_layout;
    switch(src_ch_layout){
        case 1:
            t_src_ch_layout=AV_CH_LAYOUT_MONO;
            break;
        case 2:
            t_src_ch_layout=AV_CH_LAYOUT_STEREO;
            break;
        default:
            return -101;
    }
    int divide_base;
    AVSampleFormat t_src_sample_fmt;
    switch(src_sample_fmt){
        case 1:
            divide_base=2;
            t_src_sample_fmt=AV_SAMPLE_FMT_S16;
            break;
        case 2:
            divide_base=4;
            t_src_sample_fmt=AV_SAMPLE_FMT_FLT;
            break;
        default:
            return -100;
    }

    // Prepare resampler
    struct SwrContext *swr_ctx;
    swr_ctx = swr_alloc_set_opts(
        NULL,
        AV_CH_LAYOUT_MONO,
        AV_SAMPLE_FMT_S16,
        16000,
        t_src_ch_layout,
        t_src_sample_fmt,
        src_rate,
        0,
        NULL
    );

    // Initialize resampler
    if(swr_init(swr_ctx)<0){
        Log(-1,"resampleToPcmBuffer","swr_init error");
        swr_free(&swr_ctx);
        return -500;
    }

    // Calculate expected out_num_sample
    int out_num_samples = av_rescale_rnd(
            swr_get_delay(swr_ctx, src_rate) + bufferSize / divide_base / src_ch_layout,
            16000,
            src_rate,
            AV_ROUND_UP
    );

    // Allocation buffer
    uint8_t **dst_data = NULL;
    ret = av_samples_alloc_array_and_samples(
        &dst_data,
        &dst_linesize,
        AV_CH_LAYOUT_MONO,
        out_num_samples,
        AV_SAMPLE_FMT_S16,
        0
    );

    // Do Resampling
    ret = swr_convert(
            swr_ctx,
            dst_data,
            out_num_samples,
            (const uint8_t **) &buffer,
            bufferSize / divide_base / src_ch_layout
    );

    // Calculate destination buffer size
    int dst_bufsize = av_samples_get_buffer_size(&dst_linesize, 1, ret, AV_SAMPLE_FMT_S16, 1);

    int copySize=0;
    if(dst_bufsize>=targetBufferSize) copySize=targetBufferSize;
    else copySize=dst_bufsize;
    logStr=logStr+"Resampled BufferSize:"+std::to_string(dst_bufsize)+" TargetBufferSize:"+std::to_string(targetBufferSize)+" CopySize:"+std::to_string(copySize);
    Log(-1,"resampleToPcmbuffer",logStr.c_str());
    memcpy(targetBuffer,dst_data[0],copySize);

    if(dst_data) {
        av_freep(&dst_data[0]);
        av_freep(&dst_data);
    }
    swr_free(&swr_ctx);
    return copySize;
}

#ifdef __cplusplus
}
#endif