#include "trans_frames.h"

using namespace std;
using namespace cv;



inline bool isFlowCorrect(Point2f u)
{
	return !cvIsNaN(u.x) && !cvIsNaN(u.y) && fabs(u.x) < 1e9 && fabs(u.y) < 1e9;
}


static Vec3b computeColor(float fx, float fy)
{
	static bool first = true;

	// relative lengths of color transitions:
	// these are chosen based on perceptual similarity
	// (e.g. one can distinguish more shades between red and yellow
	//  than between yellow and green)
	const int RY = 15;
	const int YG = 6;
	const int GC = 4;
	const int CB = 11;
	const int BM = 13;
	const int MR = 6;
	const int NCOLS = RY + YG + GC + CB + BM + MR;
	static Vec3i colorWheel[NCOLS];

	if (first)
	{
		int k = 0;

		for (int i = 0; i < RY; ++i, ++k)
			colorWheel[k] = Vec3i(255, 255 * i / RY, 0);

		/*
		for (int i = 0; i < YG; ++i, ++k)
		colorWheel[k] = Vec3i(255 - 255 * i / YG, 255, 0);

		for (int i = 0; i < GC; ++i, ++k)
		colorWheel[k] = Vec3i(0, 255, 255 * i / GC);

		for (int i = 0; i < CB; ++i, ++k)
		colorWheel[k] = Vec3i(0, 255 - 255 * i / CB, 255);

		for (int i = 0; i < BM; ++i, ++k)
		colorWheel[k] = Vec3i(255 * i / BM, 0, 255);

		for (int i = 0; i < MR; ++i, ++k)
		colorWheel[k] = Vec3i(255, 0, 255 - 255 * i / MR);
		*/
		first = false;

	}

	const float rad = sqrt(fx * fx + fy * fy);
	const float a = atan2(-fy, -fx) / (float)CV_PI;

	const float fk = (a + 1.0f) / 2.0f * (NCOLS - 1);
	const int k0 = static_cast<int>(fk);
	const int k1 = (k0 + 1) % NCOLS;
	const float f = fk - k0;

	Vec3b pix;

	for (int b = 0; b < 3; b++)
	{
		const float col0 = colorWheel[k0][b] / 255.0f;
		const float col1 = colorWheel[k1][b] / 255.0f;

		float col = (1 - f) * col0 + f * col1;

		if (rad <= 1)
			col = 1 - rad * (1 - col); // increase saturation with radius
		else
			col *= .75; // out of range

		pix[2 - b] = static_cast<uchar>(255.0 * col);
	}

	return pix;
}

static void drawOpticalFlow(const Mat_<float>& flowx, const Mat_<float>& flowy, Mat& dst, float maxmotion = -1)
{
	dst.create(flowx.size(), CV_8UC3);
	dst.setTo(Scalar::all(0));

	// determine motion range:
	float maxrad = maxmotion;

	if (maxmotion <= 0)
	{
		maxrad = 1;
		for (int y = 0; y < flowx.rows; ++y)
		{
			for (int x = 0; x < flowx.cols; ++x)
			{
				Point2f u(flowx(y, x), flowy(y, x));

				if (!isFlowCorrect(u))
					continue;

				maxrad = max(maxrad, sqrt(u.x * u.x + u.y * u.y));
			}
		}
	}

	for (int y = 0; y < flowx.rows; ++y)
	{
		for (int x = 0; x < flowx.cols; ++x)
		{
			Point2f u(flowx(y, x), flowy(y, x));

			if (isFlowCorrect(u))
				dst.at<Vec3b>(y, x) = computeColor(u.x / maxrad, u.y / maxrad);
		}
	}
}
static void showFlow(const char* name, const Mat& c_flow)
{
	Mat planes[2];
	split(c_flow, planes);

	Mat flowx(planes[0]);
	Mat flowy(planes[1]);

	Mat out;
	drawOpticalFlow(flowx, flowy, out, 10);

	//imshow(name, out);
	cvtColor(out, out, CV_RGB2GRAY);
	normalize(out, out, 0, 255, NORM_MINMAX);
	imwrite(name, out);
}


double cflows(Mat frame0, Mat frame1, unsigned char*&  pFlowData, int&  iFlowDataSize, char*  pDumpName, int iLogFlag)
{
	Mat fflow;
	//frame0.copyTo(cframe0);
	//frame1.copyTo(cframe1);
	const int64 start = getTickCount();
	calcOpticalFlowFarneback(frame0, frame1, fflow, 0.702, 5, 10, 2, 7, 1.5, cv::OPTFLOW_FARNEBACK_GAUSSIAN);


	const double timeSec = (getTickCount() - start) / getTickFrequency();
	if (iLogFlag != 0)
	{
		cout << "cpu Farn : " << timeSec << " sec" << endl;
	}


	showFlow(pDumpName, fflow);

	iFlowDataSize = fflow.rows*fflow.cols * 2 * sizeof(float);

	pFlowData = (unsigned char*)malloc(iFlowDataSize);
	if (pFlowData != NULL)
	{
		memcpy(pFlowData, fflow.datastart, iFlowDataSize);
	}

	return timeSec;
}

int  AddOptFlow(S_Frames_Output*  psFrameOutput, unsigned char*  pFlowData, int iFlowDataSize, int iIndex)
{
	psFrameOutput->pFlowData[iIndex] = pFlowData;
	return 0;
}


int  CalOptFlow(S_Frames_Output**  ppsFrameOutput, int iFrameOutputCount, char*  pURL, char*   pOutputDir, char*  pPostFix, int iLogFlag)
{
	char  strTmp[1024] = { 0 };
	int  iRet = 0;
	int  iIndex = 0;
	int  iLoop = 0;
	unsigned char*  pFlowData = NULL;
	int             iFlowDataSize = 0;
	S_Frames_Output*   psFrameOutput = NULL;
	char  strDir[1024] = { 0 };

	do
	{
		if (ppsFrameOutput == NULL)
		{
			break;
		}

		if (pOutputDir != NULL)
		{
			strcpy(strDir, pOutputDir);
			if (strlen(strDir) > 0 && strDir[strlen(strDir) - 1] != '/')
			{
				strDir[strlen(strDir)] = '/';
			}
		}

		for (iLoop = 0; iLoop < iFrameOutputCount; iLoop++)
		{
			psFrameOutput = ppsFrameOutput[iLoop];
			for (iIndex = 1; iIndex < psFrameOutput->iFrameCount; iIndex++)
			{
				memset(strTmp, 0, 1024);
				pFlowData = NULL;
				iFlowDataSize = 0;

				cv::Mat grey_prev, grey_cur;
				cv::Mat matPre(psFrameOutput->iHeight, psFrameOutput->iWidth, CV_8UC3, psFrameOutput->pFrameData[iIndex - 1]);
				cvtColor(matPre, grey_prev, CV_BGR2GRAY);

				memset(strTmp, 0, 1024);
				if (iIndex == 1)
				{
					sprintf(strTmp, "%s%s_frame_%04d_%04d.%s", strDir, pURL, iLoop, iIndex - 1, pPostFix);
					cv::imwrite(strTmp, matPre);
				}

				cv::Mat matCur(psFrameOutput->iHeight, psFrameOutput->iWidth, CV_8UC3, psFrameOutput->pFrameData[iIndex]);
				cvtColor(matCur, grey_cur, CV_BGR2GRAY);

				memset(strTmp, 0, 1024);
				sprintf(strTmp, "%s%s_frame_%04d_%04d.%s", strDir, pURL, iLoop, iIndex, pPostFix);
				cv::imwrite(strTmp, matCur);

				sprintf(strTmp, "%s%s_flow_%04d_%04d_%04d.%s", strDir, pURL, iLoop, iIndex - 1, iIndex, pPostFix);
				cflows(grey_prev, grey_cur, pFlowData, iFlowDataSize, strTmp, iLogFlag);
				AddOptFlow(psFrameOutput, pFlowData, iFlowDataSize, iIndex - 1);
			}
		}

	} while (0);

	return 0;
}

int  DumpFrameAndFlow(S_Frames_Output*  psFrameOutput)
{
	char strDumpName[265] = { 0 };
	int iIndex = 0;

	for (iIndex = 0; iIndex < psFrameOutput->iFrameCount; iIndex++)
	{
		if (psFrameOutput->pFrameData[iIndex] != NULL)
		{
			cv::Mat matFrame(psFrameOutput->iHeight, psFrameOutput->iWidth, CV_8UC3, psFrameOutput->pFrameData[iIndex]);
			memset(strDumpName, 0, 256);
			sprintf(strDumpName, "dump_frame_%d.png", iIndex);
			cv::imwrite(strDumpName, matFrame);
		}
	}

	for (iIndex = 0; iIndex < psFrameOutput->iFrameCount; iIndex++)
	{
		if (psFrameOutput->pFlowData[iIndex] != NULL)
		{
			Mat planes[2];
			cv::Mat matFlow(psFrameOutput->iHeight, psFrameOutput->iWidth, CV_32FC2, psFrameOutput->pFlowData[iIndex]);
			memset(strDumpName, 0, 256);
			sprintf(strDumpName, "dump_flow_%d.png", iIndex);
			split(matFlow, planes);

			Mat flowx(planes[0]);
			Mat flowy(planes[1]);

			Mat out;
			drawOpticalFlow(flowx, flowy, out, 10);
			imwrite(strDumpName, out);
		}
	}

	return 0;
}
