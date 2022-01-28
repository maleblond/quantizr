package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func getRandomPixelMatrix() [][]int {
	rand.Seed(time.Now().UTC().UnixNano())

	inputMatrix := make([][]int, 8)

	for i := range inputMatrix {
		inputMatrix[i] = make([]int, 8)

		for j := range inputMatrix[i] {
			inputMatrix[i][j] = rand.Intn(256)
		}
	}

	return inputMatrix
}

func getStaticPixelMatrix() [][]int {
	return [][]int{
		{140, 144, 147, 140, 140, 155, 179, 175},
		{144, 152, 140, 147, 140, 148, 167, 179},
		{152, 155, 136, 167, 163, 162, 152, 172},
		{168, 145, 156, 160, 152, 155, 136, 160},
		{162, 148, 156, 148, 140, 136, 147, 162},
		{147, 167, 140, 155, 155, 140, 136, 162},
		{136, 156, 123, 167, 162, 144, 140, 147},
		{148, 155, 136, 155, 152, 147, 147, 136},
	}
}
func getDCTMatrix(inputMatrix [][]int) [][]int {
	dctMatrix := make([][]int, 8)

	for i := range dctMatrix {
		dctMatrix[i] = make([]int, 8)

		for j := range dctMatrix {
			temp := 0.0

			for x := 0; x < 8; x++ {
				for y := 0; y < 8; y++ {
					temp += math.Cos(float64(2*x+1)*math.Pi*float64(i)/16.0) * math.Cos(float64(2*y+1)*math.Pi*float64(j)/16.0) * float64(inputMatrix[x][y]-128)
				}
			}

			coeff := 1.0

			if i == 0 {
				coeff /= math.Sqrt(2)
			}

			if j == 0 {
				coeff /= math.Sqrt(2)
			}

			temp *= coeff / 4

			dctMatrix[i][j] = int(math.Round(temp))
		}
	}

	return dctMatrix
}

func getQuantizedMatrix(dctMatrix [][]int) [][]int {
	quantizationMatrix := [][]int{
		{3, 5, 7, 9, 11, 13, 15, 17},
		{5, 7, 9, 11, 13, 15, 17, 19},
		{7, 9, 11, 13, 15, 17, 19, 21},
		{9, 11, 13, 15, 17, 19, 21, 23},
		{11, 13, 15, 17, 19, 21, 23, 25},
		{13, 15, 17, 19, 21, 23, 25, 27},
		{15, 17, 19, 21, 23, 25, 27, 29},
		{17, 19, 21, 23, 25, 27, 29, 31},
	}
	// quantizationMatrix := [][]int{
	// 	{255, 255, 255, 255, 255, 255, 255, 255},
	// 	{55, 60, 60, 70, 95, 130, 255, 255},
	// 	{70, 65, 80, 120, 200, 255, 255, 255},
	// 	{70, 85, 110, 145, 255, 255, 255, 255},
	// 	{90, 110, 185, 255, 255, 255, 255, 255},
	// 	{120, 175, 255, 255, 255, 255, 255, 255},
	// 	{245, 255, 255, 255, 255, 255, 255, 255},
	// 	{255, 255, 255, 255, 255, 255, 255, 255},
	// }

	quantizedMatrix := make([][]int, 8)

	for x := range dctMatrix {
		quantizedMatrix[x] = make([]int, 8)

		for y := range dctMatrix[x] {
			quantizedMatrix[x][y] = int(math.Round(float64(dctMatrix[x][y]) / float64(quantizationMatrix[x][y])))
		}
	}

	return quantizedMatrix
}

func printMatrix(matrix [][]int) {
	fmt.Println("----------------------")
	for _, line := range matrix {
		fmt.Println(line)
	}

	fmt.Println("----------------------")
}

func main() {
	// Steps to encode a jpeg
	// 1- Divide the pixel matrix into 8x8 blocks
	// 2- Convert to YCbCr color space. We want to separate luminosity information (Y) from color information (CbCr)
	//    since human vision is pretty bad at detecting color, but much more sensible to luminosity.
	// 3- Generate a "DCT matrix" for every 8x8 block, which allows to compress more the "high-frequency" information of an image (i.e. color / brightness changing very rapidly), than the "low-frequency information".
	//    Human vision is pretty bad at noticing high frequency changes (perceptually, you'll just see an "average").
	//    The DCT matrix is also a 8x8 block, but numbers don't represent a pixel value (or channel value), it's just a coefficient of a specific function.
	//    It's a reversible function. Forward DCT function, Inverse DCT function
	// 4- Quantize that "DCT matrix" to actually compress. This step is pretty simple, we just divide each number of the DCT matrix by a given number, and then round the result (this is where the lossy compression happens).
	//    This should result in a bunch of `0` values, which makes compression more efficient. 
	// 5- More compression (?)

	// DCT
	// If you have a set of N numbers (i.e. [8,10,15]), it can be represented as a set of N numbers representing coefficients of specific cosines functions
	// If [x,y,z] are our coefficient, x * f(2) + y * g(2) + z * h(2) will give us `3`
	//
	// In 2d, it's way more complicated: https://en.wikipedia.org/wiki/Discrete_cosine_transform

	// pixelMatrix := getRandomPixelMatrix()
	pixelMatrix := getStaticPixelMatrix()

	fmt.Println("Printing pixel matrix:")
	printMatrix(pixelMatrix)

	fmt.Println("Printing DCT matrix....")
	dctMatrix := getDCTMatrix(pixelMatrix)
	printMatrix(dctMatrix)

	fmt.Println("Printing quantized DCT matrix....")
	quantizedMatrix := getQuantizedMatrix(dctMatrix)
	printMatrix(quantizedMatrix)
}
