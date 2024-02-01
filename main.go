package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

// ファイルの内容をハイライトして返す関数
func highlightCode(fileName string, code string) (string, error) {
	lexer := lexers.Match(fileName)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return "", err
	}

	style := styles.Get("nord")
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	var b bytes.Buffer
	err = formatter.Format(&b, style, iterator)
	if err != nil {
		return "", err
	}
	return b.String(), nil

}

func cat(r *bufio.Reader, fileName string) {
	for {
		// 改行文字が見つかるまで読み込む
		buf, err := r.ReadBytes('\n')
		if err == io.EOF {
			break
		}

		highlighted, err := highlightCode(fileName, string(buf))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error highlighting code:", err)
			continue
		}
		// 標準出力に書き込む
		fmt.Fprintf(os.Stdout, "%s", highlighted)
	}
}

func main() {
	// コマンドライン引数を解析
	flag.Parse()

	if flag.NArg() == 0 {
		// 引数がない場合は標準入力を表示
		cat(bufio.NewReader(os.Stdin), "")
	}

	for i := 0; i < flag.NArg(); i++ {
		// ファイルを開く
		f, err := os.Open(flag.Arg(i))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: error reading from %s: %s\n", os.Args[0], flag.Arg(i), err)
			continue
		}
		// ファイルの内容を表示
		cat(bufio.NewReader(f), flag.Arg(i))
		f.Close()
	}
}
