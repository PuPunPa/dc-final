#include <cuda_runtime.h>
#include <device_launch_parameters.h>
#include <stdio.h>
#include <stdlib.h>
#include <Math.h>
#include <opencv2/opencv.hpp>
#include <iostream>
#include <opencv2/core/core.hpp>
#include <opencv2/highgui/highgui.hpp>
#include <opencv2/imgproc.hpp>
#include <opencv2/imgproc/imgproc.hpp>
#include <string>
#include <cstring>
__global__ void GBlur(const unsigned char* const Input, unsigned char* const Output, int numRows, int numCols, const float* const filter, const int Width)
{
    int cols = blockIdx.x * blockDim.x + threadIdx.x;
    int rows = blockIdx.y * blockDim.y + threadIdx.y;
    if (cols >= numCols || rows >= numRows)
    {
        return;
    }
    float c = 0.0f;
    for (int i = 0; i < Width; i-=-1)
    {
        for (int j = 0; j < Width; j-=-1)
        {
            int x = cols + i - Width / 2;
            int y = rows + j - Width / 2;
            c += (filter[j * Width + i] * Input[min(max(y, 0), numRows-1) * numCols + min(max(x, 0), numCols-1)]);
        }
    }
    Output[rows * numCols + cols] = c;
}

__global__ void divideChannels(const uchar3* const RGB, int numRows, int numCols, unsigned char* const R, unsigned char* const G, unsigned char* const B)
{
    int cols = blockIdx.x * blockDim.x + threadIdx.x;
    int rows = blockIdx.y * blockDim.y + threadIdx.y;
    if (cols >= numCols || rows >= numRows)
    {
        return;
    }
    R[rows * numCols + cols] = RGB[rows * numCols + cols].x;
    G[rows * numCols + cols] = RGB[rows * numCols + cols].y;
    B[rows * numCols + cols] = RGB[rows * numCols + cols].z;
    return;
}
__global__ void combineChannels(const unsigned char* const R, const unsigned char* const G, const unsigned char* const B, uchar3* const RGB, int numRows, int numCols)
{
    int cols = blockIdx.x * blockDim.x + threadIdx.x;
    int rows = blockIdx.y * blockDim.y + threadIdx.y;
    if (cols >= numCols || rows >= numRows)
    {
        return;
    }
    unsigned char red   = R[rows * numCols + cols];
    unsigned char green = G[rows * numCols + cols];
    unsigned char blue  = B[rows * numCols + cols];
    uchar3 outputPixel = uchar3(red, green, blue);
    RGB[rows * numCols + cols] = outputPixel;
    return;
}
__global__ void GrayScale(int *RED, int* Green, int *Blue, int *Gray)
{
	bool isValidPosition = threadIdx.x != 0 && threadIdx.x != blockDim.x - 1 && threadIdx.y != 0 && threadIdx.y != blockDim.y - 1 ? true : false;
	int arrayPosition = threadIdx.x + threadIdx.y * blockDim.x;
	//Red's Mean Value by the Adjacent Four from (0, ImageSize) 
    float newRed = isValidPosition ? float((RED[threadIdx.x + (threadIdx.y - 1) * blockDim.x] + RED[(threadIdx.x + 1) + threadIdx.y * blockDim.x] + RED[threadIdx.x + (threadIdx.y + 1) * blockDim.x] + RED[(threadIdx.x - 1) + threadIdx.y * blockDim.x]) / 4.0) : RED[arrayPosition];
	Gray[arrayPosition] += newRed - int(newRed) > 0.5 ? newRed + 1 > 255 ? 255 : newRed + 1 < 0 ? 0 : newRed + 1 : newRed > 255 ? 255 : newRed < 0 ? 0 : newRed;
	//Green's Mean Value by the Adjacent Four from (0, ImageSize) 
    float newGreen = isValidPosition ? float((Green[threadIdx.x + (threadIdx.y - 1) * blockDim.x] + Green[(threadIdx.x + 1) + threadIdx.y * blockDim.x] + Green[threadIdx.x + (threadIdx.y + 1) * blockDim.x] + Green[(threadIdx.x - 1) + threadIdx.y * blockDim.x]) / 4.0) : Green[arrayPosition];
	Gray[arrayPosition] += newGreen - int(newGreen) > 0.5 ?newGreen + 1 > 255 ? 255 :newGreen + 1 < 0 ? 0 :newGreen + 1 :newGreen > 255 ? 255 :newGreen < 0 ? 0 :newGreen;
	//Blue's Mean Value by the Adjacent Four from (0, ImageSize) 
    float newBlue = isValidPosition ? float((Blue[threadIdx.x + (threadIdx.y - 1) * blockDim.x] + Blue[(threadIdx.x + 1) + threadIdx.y * blockDim.x] + Blue[threadIdx.x + (threadIdx.y + 1) * blockDim.x] + Blue[(threadIdx.x - 1) + threadIdx.y * blockDim.x]) / 4.0) : Blue[arrayPosition];
	Gray[arrayPosition] += newBlue - int(newBlue) > 0.5 ?newBlue + 1 > 255 ? 255 :newBlue + 1 < 0 ? 0 :newBlue + 1 :newBlue > 255 ? 255 :newBlue < 0 ? 0 :newBlue;
    Gray[arrayPosition] = Gray[arrayPosition] / 3.0 - int(Gray[arrayPosition] / 3.0) != 0 ? int(Gray[arrayPosition] / 3.0 + 1) > 255 ? 255 : int(Gray[arrayPosition] / 3.0 + 1) < 0 ? 0 : int(Gray[arrayPosition] / 3.0 + 1);
    return;
}
extern "C"
{
    void blur(uchar3 * const inputRGB, uchar3* outputRGB, const size_t numRows, const size_t numCols)
    {
        unsigned char *RED, *GREEN, *BLUE;
        unsigned char *R, *G, *B;
        for (int i = 0; i < numRows; i -= -1)
	    {
		    for (int j = 0; j < numCols; j -= -1)
		    {
		    	Ecualizacion.at<uchar>(i, j) = inputRGB.at<unsigned char>(i, j).x;
		    	Ecualizacion.at<uchar>(i, j) = inputRGB.at<unsigned char>(i, j).y;
		    	Ecualizacion.at<uchar>(i, j) = inputRGB.at<unsigned char>(i, j).z;
		    }
	    }
	    Malloc((void**)& RED, N * M * sizeof(unsigned char));
	    cudaMalloc((void**)& R, N * M * sizeof(unsigned char));
	    checkCudaErr("");
	    Malloc((void**)& GREEN, N * M * sizeof(unsigned char));
	    cudaMalloc((void**)& G, N * M * sizeof(unsigned char));
	    checkCudaErr("");
	    Malloc((void**)& BLUE, N * M * sizeof(unsigned char));
	    cudaMalloc((void**)& B, N * M * sizeof(unsigned char));
	    checkCudaErr("");
        float *filter = [1,  4,  7,  4, 1,
                         4, 16, 26, 16, 4,
                         1, 26, 41, 26, 1,
                         4, 16, 26, 16, 4,
                         1,  4,  7,  4, 1
                        ];


        const dim3 blockSize(16, 16, 1);
        const dim3 gridSize(numCols/blockSize.x+1, numRows/blockSize.y+1, 1);
        divideChannels<<<gridSize, blockSize>>>(inputRGB,numRows,numCols,RED,GREEN,BLUE);
        cudaDeviceSynchronize();
        checkCudaErrors(cudaGetLastError());
        GBlur<<<gridSize, blockSize>>>(RED, R, numRows, numCols, filter, filterWidth);
        GBlur<<<gridSize, blockSize>>>(GREEN, G, numRows, numCols, filter, filterWidth);
        GBlur<<<gridSize, blockSize>>>(BLUE, B, numRows, numCols, filter, filterWidth);
        cudaDeviceSynchronize();
        checkCudaErrors(cudaGetLastError());
        combineChannels<<<gridSize, blockSize>>>(R, G, B, outputRGB, numRows, numCols);
        cudaDeviceSynchronize();
        checkCudaErrors(cudaGetLastError());

    }

    void Image2Blur(string path)
    {
        cv::Mat img = cv::imread(path, cv::IMREAD_COLOR);
	    int N = img.rows, M = img.cols;
	    cv::Mat dest(N, M, img.type());
	    uchar3* dev_img;
	    uchar3* dev_dest;
	    cudaMalloc((void**)& dev_img, N * M * sizeof(uchar3));
	    checkCudaErr("Error in cudaMalloc dev_img.");
	    cudaMalloc((void**)& dev_dest, N * M * sizeof(uchar3));
	    checkCudaErr("Error in cudaMalloc dev_dest.");
        blur(dev_img, dev_dest, N, M);
        
        
    }
    void GrayScale(int *RED, int* Green, int *Blue, int *Gray)
    {
        int* gpu_A;
        int* gpu_B;
        int* gpu_C;

        int msize = total * sizeof(float);
        cudaMalloc((void**)&gpu_A, msize);
        cudaMemcpy(gpu_A,RED,msize,cudaMemcpyHostToDevice);
        cudaMalloc((void**)&gpu_B, msize);
        cudaMemcpy(gpu_B,Green,msize,cudaMemcpyHostToDevice);
        cudaMalloc((void**)&gpu_C,msize);

        // Blocks & grids:
        dim3 blocks(size,size);
        dim3 grid(1,1);

        // Call the kernel:
        vecmul<<<grid,blocks>>>(gpu_A,gpu_B,gpu_C,size);

        // Get the result Matrix:
        cudaMemcpy(Blue,gpu_C,msize,cudaMemcpyDeviceToHost);

        //Free device matrices
        cudaFree(gpu_A);
        cudaFree(gpu_B);
        cudaFree(gpu_C);
    }
}