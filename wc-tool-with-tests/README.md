
# Word count exercise

Write a command line program that implements Unix wc like functionality.


## How to execute the program
The below command will create the build
``` 
make build
```
Usage
```
 ./wc -l -c -w ./test-data/file1.txt                            (all flags)
 ./wc  -w ./test-data/file1.txt                                 (any single flag)
 ./wc ./test-data/file1.txt                                     (no flag)
 ./wc -l -c -w ./test-data/file1.txt ./test-data/file2.txt      (multiple files)
 ./wc -l -w ./test-data/file1.txt ./test-data/file2.txt         (multiple files, few flgs)
 ./wc ./test-data/file1.txt ./test-data/file2.txt               (multiple files, no flg)
 ./wc -w
    abc
    def ghi jkl                                                 (stdin can also be provided, with combination of flags or no flag)
```


## About the program
* As per the expectations program will throw the appropritate error messages.
* Goroutine is implemented so program can spin up one goroutine per file.
* Program is tested with over 35 files, it can process the files simultaneously.
* Except for stdin, entire the file is not getting read, instead files are getting read line by line, word by word and char by char and then calculating the respective linse, words or counts.
* Progam has been load tested with these files :- https://github.com/ravexina/shakespeare-plays-dataset-scraper/tree/master/shakespeare-db