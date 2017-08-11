#ifndef __DECODE_FRAMES__
#define __DECODE_FRAMES__

#ifdef _WIN32

#ifdef _DEBUG
#ifndef DBG_NEW
#define DBG_NEW new (_NORMAL_BLOCK , __FILE__ , __LINE__)
#define new DBG_NEW
#endif

#define _CRTDBG_MAP_ALLOC
#include <stdlib.h>
#include <crtdbg.h>
#endif



extern "C"
{
#define  __STDC_CONSTANT_MACROS
	//#define _snprintf snprintf

#include "libavutil/imgutils.h"
#include "libavutil/timestamp.h"
#include "libavcodec/avcodec.h"
#include "libavformat/avformat.h"
#include "libavutil/avutil.h"
#include "libavutil/imgutils.h"
#include "libavutil/mem.h"
#include "libavutil/fifo.h"
#include "libavutil/error.h"
#include "libswscale/swscale.h"
};
#else  
//Linux...
#ifdef __cplusplus  
extern "C"
{
#endif  
#include "libavutil/imgutils.h"
#include "libavutil/timestamp.h"
#include "libavcodec/avcodec.h"
#include "libavformat/avformat.h"
#include "libavutil/avutil.h"
#include "libavutil/mem.h"
#include "libavutil/fifo.h"
#include "libavutil/error.h"
#include "libswscale/swscale.h"
#ifdef __cplusplus  
};
#endif  
#endif  

#include "stdio.h"
#include "stdlib.h"

#define ERR_INVALID_PARAMETERS    0x80000001
#define ERR_CANNOT_FIND_VIDEO              0x80000002
#define ERR_CANNOT_FIND_VIDEO_DECODER              0x80000003
#define ERR_CANNOT_OPEN_VIDEO_DECODER              0x80000004
#define ERR_CANNOT_OPEN_VIDEO_SCALE                0x80000005
#define ERR_CANNOT_OPEN_MJPEG                      0x80000006
#define ERR_CANNOT_OPEN_MJPEG_ENCODER              0x80000007
#define ERR_CANNOT_ALLOC_MEMORY                    0x80000008
#define ERR_ENCODING_JPG                           0x80000009
#define MAX_FRAME_COUNT                            0x100


typedef struct
{
	//the pixel fmt is packed BRG24
	unsigned char*   pFrameData[MAX_FRAME_COUNT];
	unsigned char*   pFlowData[MAX_FRAME_COUNT];
	unsigned int     iWidth;
	unsigned int     iHeight;
	int              iFrameCount;
}S_Frames_Output;

int InitScaleCTX(int  iSrcWidth, int iSrcHeight, enum AVPixelFormat  eSrcPixelFmt, int  iDstWidth, int iDstHeight, enum AVPixelFormat  eDstPixelFmt, struct SwsContext** ppSwsCTX);
int DoFrameExport(char*  pInputURL, int iSeekPoint, int  iDstWidth, int iDstHeight, enum AVPixelFormat  eDstPixelFmt, int iFrameCount, S_Frames_Output** ppsFramesOutput, int iMaxOutputSize, char* pExtra, int iInterval, int* piActArrayCount);
int ConvVideoPacket(AVFrame*  pAVFrame, int iTargetWidth, int iTargetHeight, enum AVPixelFormat  eDstPixelFmt, struct SwsContext*  pSws_ctx, S_Frames_Output* psFramesOutput, int iFrameIndex);
void ReleaseFrameOutput(S_Frames_Output**  ppFrameOutput, int  iFrameOutputCount);
#endif
