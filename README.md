## Concurrent file system crawler
Implementation of a recusive file crawler in Go  
See "Handout.pdf" for description   
### How to run 
Execute 'main.go' with the following command  
`go run main.go pattern [directory]`  
where 'pattern' is a Bash regular expression surrounded by quotes and 'directory' is the directory to search. If 'directory' is not specified, it searches the current directory.  
Example usage   
`go run main.go 'a*.[ch]' my_directory`  
### Version  
Developed on  
go version go1.7.5 darwin/amd64 
