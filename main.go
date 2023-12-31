package main

import (
	"errors"
	"fmt"
	"math"
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
    for i, sample := range songData {
        songSamples[i] = float64(sample)
    }

    originalSamples := make([]float64, len(originalData))
    for i, sample := range originalData {
        originalSamples[i] = float64(sample)
    }

    // 時系列データを比較して、採点基準となる指標を算出する。
    // 例として、音程のズレを比較する指標を算出する。
    errorSamples := make([]float64, len(songSamples))
    for i, songSample := range songSamples {
        errorSamples[i] = songSample - originalSamples[i]
    }

    // 指標を元に、100点満点で採点する。
    // 例として、音程のズレが小さければ大きいほど、採点結果が高くなるようにする。
	me, _ := mean(errorSamples)
    score := 100 - math.Abs(me)

    // 採点結果を出力する。
    fmt.Println("採点結果:", score)
	fmt.Println("平均誤差:", me)
	fmt.Println("平均誤差の絶対値:", math.Abs(me))
	fmt.Println(originalSamples[100000:100100])
	fmt.Println(errorSamples[100000:100100])
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
