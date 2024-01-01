package main

import (
	"errors"
	"fmt"
	"math"
	"math/cmplx"
	"os"
	//"time"
)

func main() {
    // 歌声と原曲の音声データを取得する。
    wavFileName, _ := os.Open("two.wav")
    songData := readWav(wavFileName)

    wavFileName, _ = os.Open("one.wav")
    originalData := readWav(wavFileName)

    // 音声データを時系列データに変換する。
    songSamples := make([]float64, len(songData))
	songSamplesComplex := make([]complex128, len(songData))
    for i, sample := range songData {
        songSamples[i] = float64(sample)
		songSamplesComplex[i] = complex(float64(sample), 0.0)
    }

    originalSamples := make([]float64, len(originalData))
	originalSamplesComplex := make([]complex128, len(originalData))
    for i, sample := range originalData {
        originalSamples[i] = float64(sample)
		originalSamplesComplex[i] = complex(float64(sample), 0.0)
    }

    // 時系列データを比較して、採点基準となる指標を算出する。
    // 例として、音程のズレを比較する指標を算出する。
    errorSamples := make([]float64, len(songSamples))
    for i, songSample := range songSamples {
        errorSamples[i] = songSample - originalSamples[i]
    }

	// FFTを行う。
	songSamplesComplex = FFT(songSamplesComplex, len(songSamplesComplex))
	originalSamplesComplex = FFT(originalSamplesComplex, len(originalSamplesComplex))

    // 指標を元に、100点満点で採点する。
    // 例として、音程のズレが小さければ大きいほど、採点結果が高くなるようにする。
	me, _ := mean(errorSamples)
    score := 100 - math.Abs(me)

    // 採点結果を出力する。
    fmt.Println("採点結果:", score)
	fmt.Println("平均誤差:", me)
	fmt.Println("平均誤差の絶対値:", math.Abs(me))
	//fmt.Println(songSamples[100000:100100])
	fmt.Println(songSamplesComplex[100000:100100])
	//fmt.Println(originalSamples[100000:100100])
	//fmt.Println(errorSamples[100000:100100])
}

func readWav(file *os.File) []byte {
	// wavファイルのヘッダーを読み飛ばす。
	header := make([]byte, 44)
	file.Read(header)

	// wavファイルのボディを読み取る。
	fileInfo, _ := file.Stat()
	body := make([]byte, fileInfo.Size()-44)
	file.Read(body)

	return body
}

func mean(input []float64) (float64, error) {

	if len(input) == 0 {
		return math.NaN(), errors.New("Empty slice")
	}

	sum := 0.0
	for _, x := range input {
		sum += x
	}

	return sum / float64(len(input)), nil
}

func FFT(x []complex128, n int) []complex128 {
	y := fft(x, n)
	complex_n := complex(float64(n), 0.0)
	for i := 0; i < n; i++ {
		y[i] = y[i] / complex_n
	}
	return y
}

func fft(a []complex128, n int) []complex128 {
	x := make([]complex128, n)
	copy(x, a)

	j := 0
	for i := 0; i < n; i++ {
		if i < j {
			x[i], x[j] = x[j], x[i]
		}
		m := n / 2
		for {
			if j < m {
				break
			}
			j = j - m
			m = m / 2
			if m < 2 {
				break
			}
		}
		j = j + m
	}
	kmax := 1
	for {
		if kmax >= n {
			return x
		}
		istep := kmax * 2
		for k := 0; k < kmax; k++ {
			theta := complex(0.0, -1.0*math.Pi*float64(k)/float64(kmax))
			for i := k; i < n; i += istep {
				j := i + kmax
				if j >= n {
					//fmt.Println("j >= n")
					break
				}
				temp := x[j] * cmplx.Exp(theta)
				x[j] = x[i] - temp
				x[i] = x[i] + temp
			}
		}
		kmax = istep
	}
}