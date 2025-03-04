package main

import (
	"bytes"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBadRatio(t *testing.T) {
	testBadRatio := func(ratio float64) {
		in := bytes.Buffer{}
		out := bytes.Buffer{}
		err := Sample(ratio, 0, &in, &out)
		require.Equal(t, ErrBadRatio, err)
	}

	testBadRatio(-0.1)
	testBadRatio(1.1)
}

func TestSampling(t *testing.T) {
	numLines := 1000000
	lines := genRandomLines(numLines, 5, 10, 2, 10)

	// empty strings -> empty results
	test(t, []string{}, []string{}, 0)
	test(t, []string{}, []string{}, 0.5)
	test(t, []string{}, []string{}, 1)

	// all filtered out
	test(t, lines, []string{}, 0)

	// no sampling
	test(t, lines, lines, 1)

	// 50% sampling
	result := test(t, lines, nil, 0.5)
	expectedLength := int(float64(numLines) * 0.5)
	assertWithPrecision(t, expectedLength, len(result), 0.01) // 1% precision

	// 10% sampling
	result = test(t, lines, nil, 0.1)
	expectedLength = int(float64(numLines) * 0.1)
	assertWithPrecision(t, expectedLength, len(result), 0.01) // 1% precision

	// 1% sampling
	result = test(t, lines, nil, 0.01)
	expectedLength = int(float64(numLines) * 0.01)
	assertWithPrecision(t, expectedLength, len(result), 0.05) // 5% precision

	// 0.1% sampling
	result = test(t, lines, nil, 0.001)
	expectedLength = int(float64(numLines) * 0.001)
	assertWithPrecision(t, expectedLength, len(result), 0.05) // 5% precision

	// 0.01% sampling
	result = test(t, lines, nil, 0.0001)
	expectedLength = int(float64(numLines) * 0.0001)
	assertWithPrecision(t, expectedLength, len(result), 0.05) // 5% precision
}

func TestSeed(t *testing.T) {
	numLines := 1000
	lines := genRandomLines(numLines, 5, 10, 2, 10)

	result1Seed0 := testWithSeed(t, lines, nil, 0.5, 0)
	result2Seed0 := testWithSeed(t, lines, nil, 0.5, 0)
	result3Seed1 := testWithSeed(t, lines, nil, 0.5, 1)
	result4Seed1 := testWithSeed(t, lines, nil, 0.5, 1)
	result5Seed2 := testWithSeed(t, lines, nil, 0.5, 2)

	assert.Equal(t, result1Seed0, result2Seed0)
	assert.Equal(t, result3Seed1, result4Seed1)
	assert.NotEqual(t, result1Seed0, result3Seed1)
	assert.NotEqual(t, result1Seed0, result5Seed2)
	assert.NotEqual(t, result3Seed1, result5Seed2)
}

func test(t *testing.T, inputLines []string, outputLines []string, ratio float64) []string {
	return testWithSeed(t, inputLines, outputLines, ratio, 0)
}

func testWithSeed(t *testing.T, inputLines []string, outputLines []string, ratio float64, seed int64) []string {
	in := bytes.Buffer{}
	for _, line := range inputLines {
		in.WriteString(line)
		in.Write([]byte{'\n'})
	}

	out := bytes.Buffer{}

	err := Sample(ratio, seed, &in, &out)
	require.NoError(t, err)

	resultLines := []string{}
	str := out.String()
	if str != "" {
		resultLines = strings.Split(str, "\n")
		for len(resultLines) > 0 && resultLines[len(resultLines)-1] == "" {
			resultLines = resultLines[0 : len(resultLines)-1]
		}
	}
	if outputLines != nil {
		assert.Equal(t, len(outputLines), len(resultLines), "The amount of resulting lines is different than expected")
		require.Equal(t, outputLines, resultLines)
	}
	return resultLines
}

func genRandomLines(lines int, wordsMin int, wordsMax int, wordSizeMin int, wordSizeMax int) []string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	random := rand.New(rand.NewSource(1))
	result := make([]string, lines)
	for idx, _ := range result {
		numWords := wordsMin + random.Intn(wordsMax-wordsMin)
		line := strings.Builder{}
		for wordIdx := 0; wordIdx < numWords; wordIdx++ {
			wordSize := wordSizeMin + random.Intn(wordSizeMax-wordSizeMin)
			for charIdx := 0; charIdx < wordSize; charIdx++ {
				randomRune := letters[random.Intn(len(letters))]
				line.WriteByte(randomRune)
			}
			if wordIdx < numWords-1 {
				line.WriteByte(' ')
			}
		}
		result[idx] = line.String()
	}
	return result
}

func assertWithPrecision(t *testing.T, expected int, value int, precision float64) {
	min := int(float64(expected) * (1.0 - precision))
	max := int(float64(expected) * (1.0 + precision))
	assert.Greater(t, value, min)
	assert.Less(t, value, max)
}
