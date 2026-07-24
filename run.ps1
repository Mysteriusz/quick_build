ri qb.exe 2> $null
echo $args

go build -o qb.exe
if(test-path qb.exe) {
	./qb.exe $args
}

