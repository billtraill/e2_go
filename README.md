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
* 9x8 - 43m48.379s

One of the  problems with this puzzle is that the more constraints you add, the slower it runs and the cost/benifit is not always clear! This version runs with the minumum of constraint checks.

I tried various attempts at using go routines for concurrancy but ended up removing most/all of this at it didn't make things any faster (most of my attempt made it slower infact). 

There is no restictions/license in using any of the code here. Use it as you wish. Its probably not very well written since it was a learning exercise. 

Build and install via (assuming you've installed go)

```
go install github.com/billtraill/e2_go
```
and run using pieces files in the above link

```
e2_go pieces_07x07.txt
```

It outputs every time it progresses at least to where it got before or higher. Output looks like

```
Placed: 48 Thursday, 08-Oct-15 14:30:17 BST
Number of iterations: 3030415
 1/ 0  7/ 0  6/ 0  7/ 0  5/ 0  6/ 0  1/ 0                  1                1                1                1                1                1                1 
 7/ 1  3/ 0  2/ 0  3/ 1  2/ 0  3/ 1  1/ 0                  2                4                9               16               47               96               72 
 4/ 3  3/ 0  3/ 0  4/ 2  3/ 1  4/ 1  1/ 0                396              905             1467             2775             4997             8728             6647 
 2/ 1  1/ 0  2/ 1  2/ 0  2/ 0  2/ 0  1/ 0              25305            40736            58854            78915            97880           113303            67409 
 5/ 3  2/ 0  1/ 0  1/ 0  2/ 0  1/ 0  1/ 0             209139           230869           217324           185415           137608            85548            38441 
 4/ 2  2/ 0  1/ 0  1/ 0  1/ 0  1/ 0  1/ 0              98258            54030            22309             7179             1490              145               59 
 1/ 0  1/ 0  1/ 0  1/ 0  1/ 0  1/ 0  1/ 0                 55               21                6                4                1                1                1 
Current no of combinations: 58993488691200
    0        0        0        0        0        0        0    
 0  1  1  1  5  1  1 24  3  3  7  1  1 21  3  3 23  3  3  2  0 
    2        5        9        6        6        8        1    

    2        5        9        6        6        8        1    
 0 12  5  5 28  7  7 44  5  5 42  8  8 32  9  9 45  6  6  6  0 
    3        4        9        6        4        7        1    

    3        4        9        6        4        7        1    
 0 22  6  6 31  6  6 49  9  9 38  4  4 26  7  7 33  4  4 18  0 
    3        7        7        9        5        8        3    

    3        7        7        9        5        8        3    
 0 19  4  4 27  7  7 46  6  6 48  8  8 29  4  4 35  9  9 17  0 
    2        4        8        9        5        6        2    

    2        4        8        9        5        6        2    
 0 14  8  8 25  4  4 34  9  9 43  5  5 40  9  9 36  8  8  9  0 
    2        5        5        6        7        4        1    

    2        5        5        6        7        4        1    
 0 13  7  7 41  6  6 30  7  7 47  8  8 39  5  5 37  9  9 16  0 
    2        7        4        8        5        8        2    

    2        7        4        8        5        8        2    
 0  3  1  1  8  2  2 11  2  2 15  3  3 20  1  1 10  3  3  4  0 
    0        0        0        0        0        0        0    


finished solution 
Number of iterations to solution: 3030415

```
The numbers top left are the size of edge pair list in that position on the board and its current interation. On the right is the number of times it has placed a tile at that location. 

The current number of combinations is basically an estimate of what it would take to traverse the complete solution space (it is the lenght of the current edge pair lists multiplied together). Its not accurate in anyway it just gives a lose idea of the size of the search space to be traversed.

Number of iterations is the number of times we went round the search loop (aprox total no of times a location was visited).
