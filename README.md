# e2_go
A basic Eternity II solver in go. Done as a learning exercise to teach myself the go language (https://golang.org) which I enjoyed very much. It was never really a 'better' e2 solver, nothing new in the approach here.

It is designed to read puzzles in the format used here https://groups.yahoo.com/neo/groups/eternity_two/files/Brendan/Puzzles/

Discussion about the puzzle on https://groups.yahoo.com/neo/groups/eternity_two/info 

It handles upto  8x8 with no problems. Above this it takes a fair bit of time to solve.

One my system 
* 3x3 - 0m0.017s
* 4x4 - 0m0.019s
* 5x5 - 0m0.024s
* 7x7 - 0m0.182s
* 8x7 - 0m38.857s
* 8x8 - 122m59.957s

One of the  problems with this puzzle is that the more constraints you add, the slower it runs and the cost/benifit is not always clear! This version runs with the minumum of constraint checks.

I tried various attempts at using go routines for concurrancy but ended up removing most/all of this at it didn't make things any faster (most of my attempt made it slower infact). 

There is no restictions/license in using any of the code here. Use it as you wish. Its probably not very well written since it was a learning exercise. 

Build and install via (assuming you've installed go)

```
go install github.com/billtraill/e2_go
```
and run using pieces files in the above link

```
e2_go pieces_08x07.txt
```

