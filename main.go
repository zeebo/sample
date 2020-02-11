package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/zeebo/errs"
	"github.com/zeebo/wyhash"
)

func handle2(_ interface{}, err error) { handle(err) }

func handle(err error) {
	if err != nil {
		log.Fatalf("%+v", errs.Wrap(err))
	}
}

type posline struct {
	line string
	pos  int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <lines>\n", os.Args[0])
		os.Exit(2)
	}

	lines, err := strconv.ParseInt(os.Args[1], 10, 64)
	handle(err)

	rng := wyhash.RNG(time.Now().UnixNano())
	res := make([]posline, 0, lines)

	// fill our reservoir
	scanner := bufio.NewScanner(os.Stdin)
	for n := 0; scanner.Scan(); n++ {
		if len(res) < cap(res) {
			res = append(res, posline{pos: n, line: scanner.Text()})
		} else if i := rng.Intn(n); i < len(res) {
			res[i] = posline{pos: n, line: scanner.Text()}
		}
	}
	handle(scanner.Err())

	// sort by line number
	sort.Slice(res, func(i, j int) bool { return res[i].pos < res[j].pos })

	// flush our reservoir
	bw := bufio.NewWriter(os.Stdout)
	for _, pl := range res {
		handle2(bw.WriteString(pl.line))
		handle(bw.WriteByte('\n'))
	}
	handle(bw.Flush())
}
