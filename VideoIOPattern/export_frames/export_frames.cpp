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

#endif


#include "decode_frames.h"
#include "trans_frames.h"


int ParseCmd(int iArgc, char*  pArgv[], char** ppInputURL, int* piTargetWidth, int*  piTargetHeight, float*  pfSeekPos, int* piFrameCount, int* piInterval, int* pLogFlag, char** ppDumpURL, char** ppPostFix);


int main(int argc, char*  argv[])
{
#ifdef _WIN32
	_CrtSetDbgFlag(_CRTDBG_ALLOC_MEM_DF | _CRTDBG_LEAK_CHECK_DF);
#endif
	S_Frames_Output*   pFrameOutArray[MAX_ARRAY_COUNT] = { 0 };
	int                iFrameOutCount = MAX_ARRAY_COUNT;
	char*   pInputMediaUrl = NULL;
	char*   pOutpuPicDir = NULL;
	char*   pPostFix = NULL;
	int     iDstWidth = 0;
	int     iDstHeight = 0;
	int     iFrameCount = 0;
	float   fSeekPoint = 0;
	int     iLogFlag = 0;
	int     iInterval = 0;
	int     iFrameArrayOutCount = 0;
	int iRet = 0;

	iRet = ParseCmd(argc, argv, &pInputMediaUrl, &iDstWidth, &iDstHeight, &fSeekPoint, &iFrameCount, &iInterval, &iLogFlag, &pOutpuPicDir, &pPostFix);
	if (iRet != 0)
	{
		return 1;
	}

	if (pInputMediaUrl == NULL || iDstHeight == 0 || iDstWidth == 0 || fSeekPoint < 0 || fSeekPoint >= 1)
	{
		printf("invalid parameters or not enough parameters!\n");
		return 1;
	}

	//Set default parameter if need
	if (pPostFix == NULL)
	{
		pPostFix = "jpg";
	}

	if (iFrameCount == 0)
	{
		iFrameCount = 10;
	}

	if (pOutpuPicDir == NULL)
	{
		pOutpuPicDir = ".";
	}

	iRet = DoFrameExport(pInputMediaUrl, fSeekPoint, iDstWidth, iDstHeight, AV_PIX_FMT_BGR24, iFrameCount, pFrameOutArray, 256, NULL, iInterval, &iFrameArrayOutCount);
	CalOptFlow(pFrameOutArray, iFrameArrayOutCount, pInputMediaUrl, pOutpuPicDir, pPostFix, iLogFlag);
	ReleaseFrameOutput(pFrameOutArray, iFrameArrayOutCount);
#ifdef _WIN32
	_CrtDumpMemoryLeaks();
#endif
	return 0;
}


int ParseCmd(int iArgc, char*  pArgv[], char** ppInputURL, int* piTargetWidth, int*  piTargetHeight, float*  pfSeekPos, int* piFrameCount, int* piInterval, int* pLogFlag, char** ppDumpURL, char** ppPostFix)
{
	int  iIndex = 0;
	float    fRate = 0;
	int     iRet = 0;
	char*   pCurOpt = NULL;
	char*   pCurValue = NULL;

	*piTargetWidth = 0;
	*piTargetHeight = 0;
	*pfSeekPos = 0;
	*piFrameCount = 0;
	*pLogFlag = 0;

	if (iArgc < 9)
	{
		printf("need more parameters!\n");
		printf("-i input_url     set the input url\n"
			"-ss time_off        set the start time offset pencent\n"
			"-s size             set frame size(WxH or abbreviation)\n"
			"-o dump url         set the dump url \n"
			"-interval interval  set the loop count \n"
			"-postfix flag       set dump postfix \n"
			"-c frame count      set frame count \n"
			"-log log_mode       set log output flag \n"
		);
		return ERR_INVALID_PARAMETERS;
	}

	while (iIndex < iArgc)
	{
		pCurOpt = pArgv[iIndex];
		if (strcmp("-i", pCurOpt) == 0)
		{
			*ppInputURL = pArgv[++iIndex];
		}

		if (strcmp("-ss", pCurOpt) == 0)
		{
			iRet = sscanf(pArgv[++iIndex], "%f", pfSeekPos);
		}

		if (strcmp("-s", pCurOpt) == 0)
		{
			iRet = sscanf(pArgv[++iIndex], "%dx%d", piTargetWidth, piTargetHeight);
		}

		if (strcmp("-o", pCurOpt) == 0)
		{
			*ppDumpURL = pArgv[++iIndex];
		}

		if (strcmp("-c", pCurOpt) == 0)
		{
			iRet = sscanf(pArgv[++iIndex], "%d", piFrameCount);
		}

		if (strcmp("-postfix", pCurOpt) == 0)
		{
			*ppPostFix = pArgv[++iIndex];
		}

		if (strcmp("-interval", pCurOpt) == 0)
		{
			iRet = sscanf(pArgv[++iIndex], "%d", piInterval);
		}

		if (strcmp("-log", pCurOpt) == 0)
		{
			iRet = sscanf(pArgv[++iIndex], "%d", pLogFlag);
		}

		if (iRet == -1)
		{
			printf("invalid parameter!\n");
			return ERR_INVALID_PARAMETERS;
		}
		iIndex++;
	}

	return 0;

}
