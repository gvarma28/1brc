# 1brc

Credits: https://github.com/gunnarmorling/1brc/

To create the input `measurements.txt` file, run `python3 create_measurements.py 1_000_000_000`

The task is to write a program which reads the file, calculates the min, mean, and max temperature value per weather station, and emits the results on stdout like this (i.e. sorted alphabetically by station name, and the result values per station in the format <min>/<mean>/<max>, rounded to one fractional digit):
// output format: <weather-station>=<min>/<mean>/<max>


## Observations:
### Solution 1: ~2m10s
- Naive approach
- Single threaded
### Solution 2: ~24s
- Parallel workers, processing records in individual maps and consolidating each map at the end.
