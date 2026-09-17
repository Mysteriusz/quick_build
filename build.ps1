if ($(test-path -path "./bin") -eq $false){
	new-item -path "." -name "bin" -itemtype "Directory"
}

go build -o ./bin/qb.exe

