package main

import (
	// 標準パッケージ

	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync"
	"time"

	// ここは独自パッケージ

	// _をつけた場合は パッケージ内のinit関数のみ実行される
	"go-sample/wiredrawing"
	_ "go-sample/wiredrawing"

	"go-sample/wiredrawing/inputter"
	"go-sample/wiredrawing/parallel"

	"golang.org/x/sys/windows"
)

// var command *cobra.Command = new(cobra.Command)

// 割り込み監視用
var signal_chan chan os.Signal = make(chan os.Signal)

// ガベージコレクションを任意の時間間隔で実行
func regularsGabageCollection() {

	for {
		time.Sleep(5 * time.Second)
		runtime.GC()
	}
}

// 入力内容を保持するための履歴ファイル
var fileNameToSaveInputHistory string = ".goshell.log"

func main() {

	// var err error = exec.Command("php", "-i").Run()
	// output, err := exec.Command("php", "-i").Output()
	// fmt.Println(err)
	// fmt.Println(output)
	// var outputString string = string(output)
	// fmt.Println(outputString)

	// コマンドの実行結果をPipeで受け取る
	command := exec.Command("php", "-i")
	buffer, err := command.StdoutPipe()

	if err != nil {
		panic(err)
	}
	command.Start()

	// // StdoutPipeの実行結果のbufferを読み取る
	// loadBuffer := func(buffer io.ReadCloser) bool {
	// 	fmt.Println("bufferの読み取り完了------>")
	// 	scanner := bufio.NewScanner(buffer)
	// 	for {

	// 		if scanner.Scan() == true {
	// 			fmt.Println(scanner.Text())
	// 			continue
	// 		}
	// 		// 読み取り失敗
	// 		// fmt.Println("読み取り失敗")
	// 		break
	// 	}
	// 	defer fmt.Println("bufferの読み取り完了------>")
	// 	return true
	// }
	// loadBuffer(buffer)
	wiredrawing.LoadBuffer(buffer)

	// コンソールの監視
	signal.Notify(
		signal_chan,
		os.Interrupt,
		os.Kill,
		windows.SIGKILL,
		windows.SIGHUP,
		windows.SIGINT,
		windows.SIGTERM,
		windows.SIGQUIT,
		windows.Signal(0x13),
		windows.Signal(0x14), // Windowsの場合 SIGTSTPを認識しないためリテラルで指定する
	)

	// GCを実行
	go regularsGabageCollection()

	var exit chan int = make(chan int)
	// 割り込み対処を実行するGoルーチン
	go parallel.InterruptProcess(exit, signal_chan)

	// go (func(exit chan int) {
	// 	echo := fmt.Print
	// 	var s os.Signal
	// 	for {
	// 		s, _ = <-signal_chan
	// 		if s == syscall.SIGHUP {
	// 			echo("[syscall.SIGHUP].\r\n")
	// 			// 割り込みを無視
	// 			exit <- 1
	// 		} else if s == syscall.SIGTERM {
	// 			echo("[syscall.SIGTERM].\r\n")
	// 			exit <- 2
	// 		} else if s == os.Kill {
	// 			echo("[os.Kill].\r\n")
	// 			// 割り込みを無視
	// 			exit <- 3
	// 		} else if s == os.Interrupt {
	// 			if runtime.GOOS != "darwin" {
	// 				echo("[os.Interrupt].\r\n")
	// 			}
	// 			// 割り込みを無視
	// 			exit <- 4
	// 		} else if s == syscall.Signal(0x14) {
	// 			if runtime.GOOS != "darwin" {
	// 				echo("[syscall.SIGTSTP].\r\n")
	// 			}
	// 			// 割り込みを無視
	// 			exit <- 5
	// 		} else if s == syscall.SIGQUIT {
	// 			echo("[syscall.SIGQUIT].\r\n")
	// 			exit <- 6
	// 		}
	// 	}
	// })(exit)

	go func(exit chan int) {
		var echo = fmt.Print
		var code int = 0
		for {
			code, _ = <-exit

			if code == 1 {
				os.Exit(code)
			} else if code == 4 {
				echo("[Ignored interrupt].\r\n")
			} else {
				if runtime.GOOS != "darwin" {
					echo("[Ignored interrupt].\r\n")
				}
			}
		}
	}(exit)

	// ----------------------------------------------
	// 標準入力を可能にする
	// 標準入力の開始
	// ----------------------------------------------

	var waiter *sync.WaitGroup = new(sync.WaitGroup)
	waiter.Add(1)
	go inputter.StandByInput(waiter)
	waiter.Wait()
	// scanner := bufio.NewScanner(os.Stdin)

	// var inputText string = ""
	// for {
	// 	fmt.Print("  >> ")
	// 	var isOk bool = scanner.Scan()
	// 	if isOk != true {
	// 		fmt.Println("scanner.Scan()が失敗")
	// 		// scannerを初期化
	// 		scanner = nil
	// 		scanner = bufio.NewScanner(os.Stdin)
	// 		continue
	// 	}
	// 	inputText = scanner.Text()
	// 	// 入力内容が exit ならアプリケーションを終了
	// 	if len(inputText) > 0 {
	// 		if inputText == "exit" {
	// 			os.Exit(1)
	// 		}
	// 		fmt.Print(" ==> ")
	// 		fmt.Println(inputText)
	// 	}
	// }
	// // 標準入力の終了

	// // cobraコマンドの初期化
	// command.Use = "使い方"
	// command.Short = "some descritpion"
	// command.Long = "some long description"
	// command.Run = func(cmd *cobra.Command, arguments []string) {
	// 	fmt.Println(arguments)
	// }

	// // cobraの実行
	// var err error = command.Execute()
	// if err != nil {
	// 	fmt.Print("Some Error Happend")
	// 	panic(err)
	// 	os.Exit(-1)
	// }
	// // 外部パッケージの構造体のポインタ変数を作成する
	// var article *wiredrawing2.Article = new(wiredrawing2.Article)
	// article.SetTitle("set the title")
	// article.SetDescription("set the description")
	// fmt.Println(article)
	// wiredrawing.Print()
	// fmt.Println("test")
	// var wg sync.WaitGroup
	// fmt.Println(wg)
	// fmt.Println(quote.Hello())

	// // execute concurrency
	// var wg2 *sync.WaitGroup = new(sync.WaitGroup)
	// var result string = functionForConcurrency(wg2)
	// fmt.Println(result)
	// wiredrawing.Print()
	// samplepackage.CallableFunctionFromOtherPackage()

	// // goroutineのテスト
	// var _waiter *sync.WaitGroup = new(sync.WaitGroup)
	// echo := fmt.Println
	// _waiter.Add(1)
	// go (func(waiter *sync.WaitGroup) {
	// 	// waiter.Add(1)
	// 	echo("これはGoroutineの実行中です")
	// 	defer waiter.Done()
	// })(_waiter)
	// // time.Sleep(10)
	// _waiter.Wait()
}

// --------------------------------------
// 並行処理で実行するための関数
// --------------------------------------
// func functionForConcurrency(waiter *sync.WaitGroup) string {
// 	waiter.Add(1)
// 	// waitGroupをカウントダウンさせる
// 	defer waiter.Done()
// 	return "Return the some data you want to back"
// }
