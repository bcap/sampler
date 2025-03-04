package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/alexflint/go-arg"
)

type args struct {
	Input   string   `arg:"-i,--input" help:"Input file path. Use - to read from stdin" default:"-"`
	Percent *float64 `arg:"-p,--percent" help:"Define sampling by a percentage. Example: 20 means 20% of the samples will pass, which is the same as -r 0.2 or -s 5"`
	Ratio   *float64 `arg:"-r,--ratio" help:"Define sampling by a ratio. Example: 0.2 means 20% of the samples will pass, which is the same as -p 20 or -s 5"`
	Sample  *int64   `arg:"-s,--sample" help:"Define sampling by a denominator. Example: 5 means \"one in five of the samples will pass\", which is the same as -p 20 or -r 0.2"`
	Seed    *int64   `arg:"-S,--seed" help:"Use this seed for randomness instead of the current time"`
}

type parsedArgs struct {
	args

	input io.Reader
	seed  int64
	ratio float64
}

func main() {
	args := parse()
	if err := Sample(args.input, os.Stdout, args.ratio, args.seed); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func parse() parsedArgs {
	args := args{}
	p := arg.MustParse(&args)

	var defined = 0
	if args.Percent != nil {
		defined++
	}
	if args.Ratio != nil {
		defined++
	}
	if args.Sample != nil {
		defined++
	}
	if defined == 0 {
		p.Fail("must define either -percent, -ratio or -sample")
	} else if defined > 1 {
		p.Fail("only one of -percent, -ratio or -sample should be used")
	}

	parsed := parsedArgs{
		args:  args,
		input: os.Stdin,
	}

	if args.Input != "-" {
		file, err := os.Open(args.Input)
		if err != nil {
			p.Fail(fmt.Sprintf("failed to open input file: %s", err))
		}
		parsed.input = file
	}

	if args.Seed == nil {
		parsed.seed = time.Now().UnixNano()
	} else {
		parsed.seed = *args.Seed
	}

	if args.Percent != nil {
		pct := *args.Percent
		if pct < 0 || pct > 100 {
			p.Fail("--percent must be a float in the [0,100] range")
		}
		parsed.ratio = *args.Percent / 100.0
		return parsed
	}

	if args.Ratio != nil {
		ratio := *args.Ratio
		if ratio < 0 || ratio > 1 {
			p.Fail("--ratio must be a float in the [0,1] range")
		}
		parsed.ratio = ratio
		return parsed
	}

	sample := *args.Sample
	if sample <= 0 {
		p.Fail("--sample must be a positive integer above zero")
	}
	parsed.ratio = 1.0 / float64(sample)
	return parsed
}
