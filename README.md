# sampler

Small tool to sample input lines

```
% sampler --help
Usage: sampler [--input INPUT] [--percent PERCENT] [--ratio RATIO] [--sample SAMPLE] [--seed SEED]

Options:
  --input INPUT, -i INPUT
                         Input file path. Use - to read from stdin [default: -]
  --percent PERCENT, -p PERCENT
                         Define sampling by a percentage. Example: 20 means 20% of the samples will pass, which is the same as -r 0.2 or -s 5
  --ratio RATIO, -r RATIO
                         Define sampling by a ratio. Example: 0.2 means 20% of the samples will pass, which is the same as -p 20 or -s 5
  --sample SAMPLE, -s SAMPLE
                         Define sampling by a denominator. Example: 5 means "one in five of the samples will pass", which is the same as -p 20 or -r 0.2
  --seed SEED, -S SEED   Use this seed for randomness instead of the current time
  --help, -h             display this help and exit
```

## Example:

```
% seq 1000000 | sampler --percent .001 --seed 42
528617
661048
741589
841481
879430
894870
936412
```

## Install

`go install github.com/bcap/sampler@latest`