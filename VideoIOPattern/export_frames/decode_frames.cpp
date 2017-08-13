#include "decode_frames.h"

#ifdef _WIN32

#ifdef _DEBUG
#ifndef DBG_NEW
#define DBG_NEW new (_NORMAL_BLOCK , __FILE__ , __LINE__)
#define new DBG_NEW
#endif
#endif

#define _CRTDBG_MAP_ALLOC
#include <stdlib.h>
#include <crtdbg.h>

#endif

int DoFrameExport(char*  pInputURL, float fSeekPointPercent, int  iDstWidth, int iDstHeight, enum AVPixelFormat  eDstPixelFmt, int iFrameCount, S_Frames_Output** ppsFramesOutput,
	int iMaxFrameOutputSize, char* pExtra, int iInterval, int* piActArrayCount)
{
	AVFormatContext *pIfmt_ctx = NULL;
	AVStream *pInputStream = NULL;
	AVStream *pVideoStream = NULL;
	AVPacket pkt;
	AVFrame*   pAVFrame = NULL;
	AVFrame*   pAVFrameOutputArray[MAX_FRAME_COUNT] = { 0 };
	AVCodec*   pCodec = NULL;
	AVCodecContext*   pCodecCtx = NULL;
	struct SwsContext *pSws_ctx = NULL;
	int               iTimeStart = 0;
	long long         iFirstVideoTimeInMs = AV_NOPTS_VALUE;
	int iRet = 0;
	char   strErrorInfo[1024] = { 0 };
	int  iIndex = 0;
	int  iVideoIndex = -1;
	long long  illCurPTS = 0;
	int        iGotFrameCount = 0;
	int        iSeekConv = 0;
	int        iCount = 0;
	long long   illDuration = 0;

	av_register_all();
	avformat_network_init();
	do
	{
		iTimeStart = (int)time(NULL);
		if ((iRet = avformat_open_input(&pIfmt_ctx, pInputURL, 0, 0)) < 0)
		{
			memset(strErrorInfo, 0, 1024);
			av_strerror(iRet, strErrorInfo, 1024);
			printf("Could not open input file: %s , error info:%s\n", pInputURL, strErrorInfo);
			break;
		}


		if ((iRet = avformat_find_stream_info(pIfmt_ctx, 0)) < 0)
		{
			memset(strErrorInfo, 0, 1024);
			av_strerror(iRet, strErrorInfo, 1024);
			printf("Failed to retrieve input stream information for %s , error info:%s\n", pInputURL, strErrorInfo);
			break;
		}


		for (iIndex = 0; iIndex < pIfmt_ctx->nb_streams; iIndex++)
		{
			pInputStream = pIfmt_ctx->streams[iIndex];
			if (pInputStream->codecpar->codec_type == AVMEDIA_TYPE_VIDEO)
			{
				iVideoIndex = iIndex;
				break;
			}
		}

		if (iVideoIndex == -1)
		{
			iRet = ERR_CANNOT_FIND_VIDEO;
			break;
		}

		//Open Video Codec
		pCodec = avcodec_find_decoder(pIfmt_ctx->streams[iVideoIndex]->codecpar->codec_id);
		if (pCodec == NULL)
		{
			printf("codec not found\n");
			iRet = ERR_CANNOT_FIND_VIDEO_DECODER;
			break;
		}

		pCodecCtx = avcodec_alloc_context3(pCodec);
		avcodec_parameters_to_context(pCodecCtx, pIfmt_ctx->streams[iVideoIndex]->codecpar);
		iRet = avcodec_open2(pCodecCtx, pCodec, NULL);
		if (iRet < 0)
		{
			memset(strErrorInfo, 0, 1024);
			av_strerror(iRet, strErrorInfo, 1024);
			printf("Could not open codec error info:%s\n", strErrorInfo);
			break;
		}

		if (iDstHeight != pIfmt_ctx->streams[iVideoIndex]->codecpar->height ||
			iDstWidth != pIfmt_ctx->streams[iVideoIndex]->codecpar->width ||
			eDstPixelFmt != (enum AVPixelFormat)pIfmt_ctx->streams[iVideoIndex]->codecpar->format)
		{
			iRet = InitScaleCTX(pIfmt_ctx->streams[iVideoIndex]->codecpar->width, pIfmt_ctx->streams[iVideoIndex]->codecpar->height, (enum AVPixelFormat)pIfmt_ctx->streams[iVideoIndex]->codecpar->format,
				iDstWidth, iDstHeight, eDstPixelFmt, &pSws_ctx);
			if (iRet != 0)
			{
				break;
			}
		}

		iRet = av_seek_frame(pIfmt_ctx, -1, 0, AVSEEK_FLAG_BACKWARD);
		if (iRet != 0)
		{
			break;
		}

		pAVFrame = av_frame_alloc();
		av_init_packet(&pkt);
		while (iFirstVideoTimeInMs == AV_NOPTS_VALUE)
		{
			iRet = av_read_frame(pIfmt_ctx, &pkt);
			if (pkt.stream_index == iVideoIndex)
			{
				if (pkt.pts != AV_NOPTS_VALUE)
				{
					iFirstVideoTimeInMs = (long long)((pkt.pts * av_q2d(pIfmt_ctx->streams[pkt.stream_index]->time_base)) * 1000);
				}
				else
				{
					iFirstVideoTimeInMs = (long long)((pkt.dts * av_q2d(pIfmt_ctx->streams[pkt.stream_index]->time_base)) * 1000);
				}
				break;
			}
		}

		iSeekConv = (pIfmt_ctx->duration*fSeekPointPercent)/AV_TIME_BASE;

		if(iInterval > 0)
		{
			illDuration = pIfmt_ctx->duration*(1-fSeekPointPercent) / 1000;
			iCount = illDuration / (iInterval * 1000);
			if (iCount == 0)
			{
				iCount = 1;
			}
		}
		else
		{
			iCount = 1;
		}
			
		if(iCount >= MAX_ARRAY_COUNT)
		{
			iCount = MAX_ARRAY_COUNT -1;
		}

		for (iIndex = 0; iIndex < iCount; iIndex++)
		{
			iGotFrameCount = 0;
			ppsFramesOutput[iIndex] = (S_Frames_Output*)malloc(sizeof(S_Frames_Output));
			memset(ppsFramesOutput[iIndex], 0, sizeof(S_Frames_Output));

			ppsFramesOutput[iIndex]->iHeight = iDstHeight;
			ppsFramesOutput[iIndex]->iWidth = iDstWidth;
			if (ppsFramesOutput[iIndex] == NULL)
			{
				break;
			}

			iRet = av_seek_frame(pIfmt_ctx, -1, (iSeekConv+iIndex*iInterval)*AV_TIME_BASE, AVSEEK_FLAG_BACKWARD);
			if (iRet != 0)
			{
				break;
			}

			while (iGotFrameCount < iFrameCount && iGotFrameCount < MAX_FRAME_COUNT)
			{
				av_init_packet(&pkt);
				iRet = av_read_frame(pIfmt_ctx, &pkt);
				if(iRet != 0)
				{
					break;
				}
				if (pkt.stream_index == iVideoIndex && iRet == 0)
				{
					iRet = avcodec_send_packet(pCodecCtx, &pkt);
					iRet = avcodec_receive_frame(pCodecCtx, pAVFrame);
					if (iRet == 0)
					{
						if (pAVFrame->pts == AV_NOPTS_VALUE)
						{
							illCurPTS = (long long)((pAVFrame->pkt_dts * av_q2d(pIfmt_ctx->streams[pkt.stream_index]->time_base)) * 1000);
						}
						else
						{
							illCurPTS = (long long)((pAVFrame->pts * av_q2d(pIfmt_ctx->streams[pkt.stream_index]->time_base)) * 1000);
						}
						if ((illCurPTS - iFirstVideoTimeInMs) >= iIndex*iInterval * 1000)
						{
							iRet = ConvVideoPacket(pAVFrame, iDstWidth, iDstHeight, eDstPixelFmt, pSws_ctx, ppsFramesOutput[iIndex], iGotFrameCount);
							if (iRet == 0)
							{
								iGotFrameCount++;
							}
						}
					}
				}
			}
			ppsFramesOutput[iIndex]->iFrameCount = iGotFrameCount;
		}


		*piActArrayCount = iCount;

		//Get File Video Frame

	} while (0);


	if (pSws_ctx != NULL)
	{
		sws_freeContext(pSws_ctx);
	}

	if (pCodecCtx != NULL)
	{
		avcodec_close(pCodecCtx);
	}

	if (pIfmt_ctx != NULL)
	{
		avformat_free_context(pIfmt_ctx);
	}

	if (pAVFrame != NULL)
	{
		av_frame_free(&pAVFrame);
	}

	avformat_network_deinit();
	return iRet;
}

int InitScaleCTX(int  iSrcWidth, int iSrcHeight, enum AVPixelFormat  eSrcPixelFmt, int  iDstWidth, int iDstHeight, enum  AVPixelFormat  eDstPixelFmt, struct SwsContext** ppSwsCTX)
{
	int iRet = 0;
	struct SwsContext *pSws_ctx = NULL;
	pSws_ctx = sws_getContext(iSrcWidth, iSrcHeight, eSrcPixelFmt,
		iDstWidth, iDstHeight, eDstPixelFmt,
		SWS_SINC, NULL, NULL, NULL);

	if (pSws_ctx == NULL)
	{
		printf("Impossible to create scale context for the conversion "
			"fmt:%s s:%dx%d -> fmt:%s s:%dx%d\n",
			av_get_pix_fmt_name(eSrcPixelFmt), iSrcWidth, iSrcHeight,
			av_get_pix_fmt_name(eDstPixelFmt), iDstWidth, iDstHeight);
		iRet = ERR_CANNOT_OPEN_VIDEO_SCALE;
	}
	else
	{
		printf("Init SwsContext done!\n");
		*ppSwsCTX = pSws_ctx;
	}

	return iRet;
}

int ConvVideoPacket(AVFrame*  pAVFrame, int iTargetWidth, int iTargetHeight, enum AVPixelFormat  eDstPixelFmt, struct SwsContext*  pSws_ctx, S_Frames_Output* psFramesOutput, int iFrameIndex)
{
	unsigned char* picture_buf = NULL;
	AVFrame* picture = NULL;
	int  iRet = 0;
	unsigned char *src_data[4] = { 0 };
	unsigned char *dst_data[4] = { 0 };
	int src_linesize[4] = { 0 };
	int dst_linesize[4] = { 0 };

	do
	{
		if (iTargetHeight != pAVFrame->height || iTargetWidth != pAVFrame->width || eDstPixelFmt != (enum AVPixelFormat)pAVFrame->format)
		{

			iRet = av_image_alloc(dst_data, dst_linesize, iTargetWidth, iTargetHeight, eDstPixelFmt, 1);
			if (iRet < 0)
			{
				printf("Could not allocate destination image\n");
				iRet = ERR_CANNOT_OPEN_VIDEO_SCALE;
				break;
			}

			/* convert to destination format */
			src_data[0] = pAVFrame->data[0];
			src_data[1] = pAVFrame->data[1];
			src_data[2] = pAVFrame->data[2];

			src_linesize[0] = pAVFrame->linesize[0];
			src_linesize[1] = pAVFrame->linesize[1];
			src_linesize[2] = pAVFrame->linesize[2];

			sws_scale(pSws_ctx, (const uint8_t * const*)src_data,
				src_linesize, 0, pAVFrame->height, dst_data, dst_linesize);
		}
		else
		{
			dst_data[0] = pAVFrame->data[0];
			dst_data[1] = pAVFrame->data[1];
			dst_data[2] = pAVFrame->data[2];

			dst_linesize[0] = pAVFrame->linesize[0];
			dst_linesize[1] = pAVFrame->linesize[1];
			dst_linesize[2] = pAVFrame->linesize[2];
		}
	} while (0);

	//Use BGR24 only
	psFramesOutput->pFrameData[iFrameIndex] = dst_data[0];
	psFramesOutput->iFrameCount++;
	return 0;
}

void ReleaseFrameOutput(S_Frames_Output**  ppFrameOutput, int  iFrameOutputCount)
{
	S_Frames_Output*    pFrameOutput = NULL;
	int iIndex = 0;
	int iLoop = 0;

	for (iLoop = 0; iLoop < iFrameOutputCount; iLoop++)
	{
		pFrameOutput = ppFrameOutput[iLoop];
		for (iIndex = 0; iIndex < pFrameOutput->iFrameCount; iIndex++)
		{
			if (pFrameOutput->pFrameData[iIndex] != NULL)
			{
				//av_freep(&(pFrameOutput->pFrameData[iIndex])).
				av_freep(&(pFrameOutput->pFrameData[iIndex]));
				pFrameOutput->pFrameData[iIndex] = NULL;
			}
		}

		for (iIndex = 0; iIndex < pFrameOutput->iFrameCount; iIndex++)
		{
			if (pFrameOutput->pFlowData[iIndex] != NULL)
			{
				free(pFrameOutput->pFlowData[iIndex]);
				pFrameOutput->pFlowData[iIndex] = NULL;
			}
		}

		free(ppFrameOutput[iLoop]);
	}
}
