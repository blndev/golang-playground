package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	textPtr := flag.String("text", "", "Text to parse.")
	metricPtr := flag.String("metric", "chars", "Metric {chars|words|lines};.")
	debugPtr := flag.Bool("debug", false, "Measure unique values of a metric.")
	helpPtr := flag.Bool("help", false, "Display this message.")
	helpSHortPtr := flag.Bool("?", false, "Display this message.")

	flag.Parse()

	if *helpPtr || *helpSHortPtr {
		flag.PrintDefaults()
		os.Exit(0)
	}

	fmt.Printf("textPtr: %s, metricPtr: %s, debug: %t\n", *textPtr, *metricPtr, *debugPtr)
}
