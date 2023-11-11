## CSV Line Filter

A simple prefilter for reading csv files in golang. You can use this to reduce the number 
of allocations made by the golang `encoding/csv` package when reading files containing 
lines that you do not need.