# ASH0 decompression written in native go

### How to install and use

1. Install the package with the command `go get https://github.com/PretendoNetwork/ASH0`

2. Import the `ash0` package into your project
```go
import (
	"github/PretendoNetwork/ASH0"
)

```
3. Load your ash0 file into a uint8 slice
```go
// Note: myFuncThatLoadsAFileIntoAByteSlice() is not included with the ash0 package
data := myFuncThatLoadsAFileIntoAByteSlice("myfile.ash")
```

4. (Optional) Check if the file is a compressed ash0 file
```go
if !ash0.IsAshCompressed(data) {
	// Oh no! It's not a compressed ash0 file!
}
```

5. Decompress the file
```go
out := ash0.Decompress(data)
```

### Notes

- If you're running the code on a big endian machine, be sure to set the `IsLittleEndian` boolean to false
- There's a lot of pointer arithmetic in the source code, so if VSCode warns you about `Possible misuse of unsafe.Pointer`, just laugh it off
- If you have any questions feel free to ask me in the Pretendo discord (make sure to ping me @Nybbit#5412)
