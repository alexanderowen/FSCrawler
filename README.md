## Concurrent file system crawler
Implementation of a recusive file crawler in Go  
See "Handout.pdf" for description   
### How to run 
Execute 'main.go' with the following command  
`go run main.go pattern`  
where 'pattern' is a Bash regular expression surrounded by quotes. For example  
`go run main.go 'a*.[ch]'`  
### Other notes
Does not currently follow the handout to a T.  
..* Writes results in non-alphabetic order (does not use heap to store results, writes to stdout when match found)
