# monzo-crawler
This is a simple web crawler that crawls only the web links available from the domain monzo.com
## How to the program

1) Clone the repo [link here](git@github.com:Mutusva/monzo-crawler.git) or download the zip file
   from greenhouse [link](https://app.greenhouse.io/tests/a0905d764a69b392bc88c55a82f62501?utm_medium=email&utm_source=TakeHomeTest&utm_source=Automated) and extract the files

2) run `go mod tidy`

3) To run tests run `make test`

4) cd into the monzo-crawler folder and run the commang `go run ./cmd/main.go` to run the program with default configurations



## Notes
- I added a worker package to be able to execute the crawling in parallel but had a challenge putting it all together in the time I had.