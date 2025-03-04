# sampler

Small tool to sample input lines

```
Usage of ./sampler:
  -percent float
    	Define sampling by a percentage. Example: 20 means 20% of the samples will pass, which is the same as -ratio 0.2 or -sample 5
  -ratio float
    	Define sampling by a ratio. Example: 0.2 means 20% of the samples will pass, which is the same as -percent 20 or -sample 5
  -sample int
    	Define sampling by a denominator. Example: 5 means "one in five of the samples will pass", which is the same as -percent 20 or -ratio 0.2
  -seed int
    	Use this seed for randomness instead of the current time (default -1)
```

Example:

```
% seq 1000000 | ./sampler --percent .001 --seed 42
528617
661048
741589
841481
879430
894870
936412
```