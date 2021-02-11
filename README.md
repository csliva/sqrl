# Sqrl
Sqrl is a squaring CLI that will resize and crop or add whitespace to make square images.

## Build
First run `go get`
Then run `go build`

Add the new executable to your path 

GLOBAL OPTIONS:
   --size value, -s value  pixel size of square (default: 1000)
   --file value, -f value  specify a filename to prevent resizing the whole folder
   --expand, -e            square by adding whitespace (default: false)
   --help, -h              show help (default: false)
