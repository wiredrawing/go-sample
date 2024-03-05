package wiredrawing

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
)

//var scanner *bufio.Scanner
//
//var loadedBuffer string

// LoadBuffer ---------------------------------------------------------------------
// 引数に渡された io.ReadCloser 変数の中身を読み取り出力する
// 2060
// ---------------------------------------------------------------------
func LoadBuffer(buffer io.ReadCloser, previousLine *int, showBuffer bool, whenError bool, colorCode string) bool {
	var currentLine int

	const ensureLength int = 512

	currentLine = 0

	for {
		readData := make([]byte, ensureLength)
		n, err := buffer.Read(readData)
		readData = readData[:n]
		if err != nil {
			break
		}

		if n == 0 {
			break
		}

		from := currentLine
		to := currentLine + n
		if (currentLine + n) >= *previousLine {
			fmt.Fprintf(os.Stdout, "\033["+colorCode+"m")
			if from < *previousLine && *previousLine <= to {
				diff := *previousLine - currentLine
				tempSlice := readData[diff:]
				// 出力内容の表示フラグがtrueの場合のみ
				if showBuffer == true {
					fmt.Fprint(os.Stdout, string(tempSlice))
				}
			} else {
				// 出力内容の表示フラグがtrueの場合のみ
				if showBuffer == true {
					fmt.Fprint(os.Stdout, string(readData))
				}
			}
			// コンソールのカラーをもとにもどす
			fmt.Fprint(os.Stdout, "\033[0m")
		}
		currentLine += n
		readData = nil
	}
	// エラーチェック以外の場合
	if whenError != true {
		*previousLine = currentLine
	}
	// 使用したメモリを開放してみる
	runtime.GC()
	debug.FreeOSMemory()
	return true
}
