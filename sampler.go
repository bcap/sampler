package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"time"
)

func main() {
	ratio, seed := parse()
	panicOnErr(Sample(ratio, seed, os.Stdin, os.Stdout))
}

func parse() (ratio float64, seed int64) {
	var percent float64
	var sample int64

	flag.Float64Var(&percent, "percent", 0, "Define sampling by a percentage. Example: 20 means 20% of the samples will pass, which is the same as -ratio 0.2 or -sample 5")
	flag.Float64Var(&ratio, "ratio", 0, "Define sampling by a ratio. Example: 0.2 means 20% of the samples will pass, which is the same as -percent 20 or -sample 5")
	flag.Int64Var(&sample, "sample", 0, "Define sampling by a denominator. Example: 5 means \"one in five of the samples will pass\", which is the same as -percent 20 or -ratio 0.2")
	flag.Int64Var(&seed, "seed", -1, "Use this seed for randomness instead of the current time")
	flag.Parse()

	flagsUsed := 0
	if percent != 0 {
		flagsUsed++
	}
	if ratio != 0 {
		flagsUsed++
	}
	if sample != 0 {
		flagsUsed++
	}

	if flagsUsed == 0 {
		panicOnErr(fmt.Errorf("must define either -percent, -ratio or -sample"))
	} else if flagsUsed > 1 {
		panicOnErr(fmt.Errorf("only one of -percent, -ratio or -sample should be used"))
	}

	if percent < 0.0 || percent > 100.0 {
		panicOnErr(fmt.Errorf("-percent must be a float in the [0,100] range"))
	} else if ratio < 0 || ratio > 1.0 {
		panicOnErr(fmt.Errorf("-ratio must be a float in the [0,1] range"))
	} else if sample < 0 {
		panicOnErr(fmt.Errorf("-sample must be a positive integer"))
	}

	if percent != 0 {
		ratio = percent / 100.0
	} else if sample != 0 {
		ratio = 1 / float64(sample)
	}

	if seed == -1 {
		seed = time.Now().UnixMilli()
	}

	return
}

var ErrBadRatio = errors.New("ratio must be in the [0,1] range")

func Sample(ratio float64, seed int64, reader io.Reader, writer io.Writer) error {
	if ratio < 0 || ratio > 1.0 {
		return ErrBadRatio
	}

	var max int64 = math.MaxInt64 >> 2
	target := int64(float64(max) * ratio)
	random := rand.New(rand.NewSource(seed))

	bufReader := bufio.NewReader(reader)
	bufWriter := bufio.NewWriter(writer)
	defer bufWriter.Flush()

	for {
		line, err := bufReader.ReadString('\n')
		if err == io.EOF && line == "" {
			return nil
		}
		if err != nil && err != io.EOF {
			return err
		}

		if ratio == 1.0 || random.Int63n(max) < target {
			if _, err := bufWriter.WriteString(line); err != nil {
				return err
			}
		}
	}
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
