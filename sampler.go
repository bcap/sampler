package main

import (
	"bufio"
	"errors"
	"io"
	"math"
	"math/rand"
)

var ErrBadRatio = errors.New("ratio must be in the [0,1] range")

func Sample(reader io.Reader, writer io.Writer, ratio float64, seed int64) error {
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
