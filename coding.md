
### Finding Allocations

Side note: There is a new testing approach: TBDD = Test n' Benchmark Driven
Development.

On the first run we got this result:

```
$ go test -run=😶 -bench=Benchmark_WithInitStoreByToken .
PASS
Benchmark_WithInitStoreByToken-4	  100000	     17297 ns/op	    9112 B/op	     203 allocs/op
ok  	github.com/corestoreio/csfw/store	2.569s

```

Quite shocking to use 203 allocs for just figuring out the current store view
within a request.

Now compile your tests into an executable binary:

```
$ go test -c
```

This compilation reduces the noise in the below output trace log.

```
$ GODEBUG=allocfreetrace=1 ./store -test.run=😶 -test.bench=Benchmark_WithInitStoreByToken -test.benchtime=10ms 2>trace.log
```

Now open the trace.log file (around 26MB) and investigate all the allocations
and refactor your code. Once finished you can achieve results like:

```
$ go test -run=NONE -bench=Benchmark_WithInitStoreByToken .
PASS
Benchmark_WithInitStoreByToken-4	 2000000	       826 ns/op	     128 B/op	       5 allocs/op
ok  	github.com/corestoreio/csfw/store	2.569s
```

### Profiling

```
$ go test -cpuprofile=cpu.out -benchmem -memprofile=mem.out -run=NONE -bench=NameOfBenchmark -v
$ go tool pprof packageName.test cpu.out
Entering interactive mode (type "help" for commands)
(pprof) top5
560ms of 1540ms total (36.36%)
Showing top 5 nodes out of 112 (cum >= 60ms)
      flat  flat%   sum%        cum   cum%
     180ms 11.69% 11.69%      400ms 25.97%  runtime.mallocgc
```

- `flat` is how much time is spent inside of a function.
- `cum` shows how much time is spent in a function, and also in any code called by a function.

For memory profile:

```
Sample value selection option (for heap profiles):
  -inuse_space      Display in-use memory size
  -inuse_objects    Display in-use object counts
  -alloc_space      Display allocated memory size
  -alloc_objects    Display allocated object counts

$ go tool pprof -alloc_objects packageName.test mem.out
```

### Bound Check Elimination

[http://klauspost-talks.appspot.com/2016/go17-compiler.slide](http://klauspost-talks.appspot.com/2016/go17-compiler.slide)

```
$ go build -gcflags="-d=ssa/check_bce/debug=1" bounds.go
or
$ go test -gcflags="-d=ssa/check_bce/debug=1" .
```

Success - Check bounds outside the loop. 

### Running Benchmark

Assuming we have already an existing file called `bm_baseline.txt`.

```
$ go test -v -run=🤐 -bench=. -count=10 . > bm_baseline_new.txt
```

After running above command to generate the second benchmark statistics file
we run:

```
$ benchstat bm_baseline.txt bm_baseline_new.txt
```

[https://godoc.org/rsc.io/benchstat](https://godoc.org/rsc.io/benchstat)

#### Other development helpers

- [go get github.com/maruel/panicparse/cmd/pp](https://github.com/maruel/panicparse)
- [go get github.com/alecthomas/gometalinter](https://github.com/alecthomas/gometalinter)

A preconfigured linter file `lint` has been included in this repoistory.
