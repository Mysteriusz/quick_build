ri qb.exe 2> $null

go build -o qb.exe
if(test-path qb.exe) {
	./qb.exe
}

