/*
	Implementation written by Nybbit.
	If you have any questions, contact me either on Twitter (@Nybbit) or on the Pretendo discord server (@Nybbit#5412).
	Based off https://github.com/trapexit/wiiqt/blob/master/WiiQt/ash.cpp
*/

/*
	Notes and C++ equivalents of parts of the code

	- Go does have support for pointer arithmetic
	- Go is 64 bit, but this implementation works on both 32 bit and 64 bit
	- Most modern CPUs are little endian,
		but the original (https://github.com/giantpune/wii-system-menu-player/blob/master/source/utils/ash.cpp)
		implementation was written for big endian


	[C++ CODE]
	- *reinterpret_cast<uint32*>(in.data())
	OR (less safe)
	- *(uint32*)(&data[0])

	[GOLANG TRANSLATION]
	- *(*uint32)(unsafe.Pointer(&data[0]))
*/

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"unsafe"
)

// IsLittleEndian - Whether or not the CPU is little endian. Most modern CPUs are but if yours isn't for some reason, change this
const IsLittleEndian = true

// loadAshFileIntoSlice - Loads an ash file into a uint8 slice
func loadAshFileIntoSlice(path string) []uint8 {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	return b
}

// BYTESWAP32 - If it's little endian, it sets it to big endian, otherwise it doesn't do anything
func BYTESWAP32(in uint32) uint32 {
	if IsLittleEndian {
		value := make([]uint8, unsafe.Sizeof(uint32(0)))
		binary.LittleEndian.PutUint32(value, in)
		return binary.BigEndian.Uint32(value)
	}
	return in
}

// BYTESWAP16 - If it's little endian, it sets it to big endian, otherwise it doesn't do anything
func BYTESWAP16(in uint16) uint16 {
	if IsLittleEndian {
		value := make([]uint8, unsafe.Sizeof(uint16(0)))
		binary.LittleEndian.PutUint16(value, in)
		return binary.BigEndian.Uint16(value)
	}
	return in
}

// isAshCompressed - Checks if an ASH file is compressed
func isAshCompressed(data []uint8) bool {
	return uint32(len(data)) > 0x10 &&
		(BYTESWAP32(*(*uint32)(unsafe.Pointer(&data[0])))&0xFFFFFF00) == 0x41534800
}

// decompressAsh - Decompresses an ASH file
func decompressAsh(data []uint8) []uint8 {
	// Check to make sure that it is a compressed ASH file
	if !isAshCompressed(data) {
		log.Fatalln("Ash is not compressed")
	}

	var r [32]uint32
	var count uint32

	var memAddr = uint64(uintptr(unsafe.Pointer(&data[0]))) // in
	r[4] = 0x8000000
	var inDiff = int64(memAddr - uint64(r[4])) // hack to support higher memory addresses than crediar's version

	r[5] = 0x415348
	r[6] = 0x415348

	r[5] = BYTESWAP32(*(*uint32)(unsafe.Pointer(uintptr(int64(r[4])+inDiff) + 4))) // "Possible misuse of unsafe.Pointer" lol
	r[5] = r[5] & 0x00FFFFFF

	var size = r[5]

	fmt.Printf("Decompressed size: %d\n", size)

	crap2 := make([]uint8, size)
	if uint32(len(crap2)) != size {
		log.Fatalln("Out of memory 1")
	}
	for i := uint32(0); i < size; i++ {
		crap2[i] = '\000'
	}

	var memAddr2 = uint64(uintptr(unsafe.Pointer(&crap2[0])))
	r[3] = 0x9000000
	var outDiff = int64(memAddr2 - uint64(r[3]))

	var o = r[3]

	r[24] = 0x10
	r[28] = BYTESWAP32(*(*uint32)(unsafe.Pointer(uintptr(int64(r[4])+inDiff) + 8)))
	r[25] = 0
	r[29] = 0
	r[26] = BYTESWAP32(*(*uint32)(unsafe.Pointer(uintptr(int64(r[4])+inDiff) + 0xC)))
	r[30] = BYTESWAP32(*(*uint32)(unsafe.Pointer(uintptr(int64(r[4]) + int64(r[28]) + inDiff))))
	r[28] = r[28] + 4

	// HACK, pointer to RAM
	var largerSize uint32

	if size > 0x100000 {
		largerSize = size
	} else {
		largerSize = 0x100000
	}
	crap3 := make([]uint8, largerSize)
	if uint32(len(crap3)) != largerSize {
		log.Fatalln("Out of memory 2")
	}
	for i := uint32(0); i < largerSize; i++ {
		crap3[i] = '\000'
	}

	var memAddr3 = uint64(uintptr(unsafe.Pointer(&crap3[0])))
	r[8] = 0x84000000
	var outDiff2 = int64(memAddr3 - uint64(r[8]))

	// memset(reinterpret_cast<void*>(r[8] + outDiff2), 0, 0x100000);
	// Set the next 0x100000 bytes to 0 starting from the pointer at r[8] + outDiff2
	for i := int64(0); i < 0x100000; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(int64(r[8]) + outDiff2 + i*int64(unsafe.Sizeof(uint8(0)))))) = uint8(0)
	}

	r[9] = r[8] + 0x07FE
	r[10] = r[9] + 0x07FE
	r[11] = r[10] + 0x1FFE
	r[31] = r[11] + 0x1FFE
	r[23] = 0x200
	r[22] = 0x200
	r[27] = 0

loc_81332124:

	if r[25] != 0x1F {
		goto loc_81332140
	}

	r[0] = r[26] >> 31
	r[26] = BYTESWAP32(*(*uint32)(unsafe.Pointer(uintptr(int64(r[4]+r[24]) + inDiff))))
	r[25] = 0
	r[24] = r[24] + 4
	goto loc_8133214C

loc_81332140:

	r[0] = r[26] >> 31
	r[25] = r[25] + 1
	r[26] = r[26] << 1

loc_8133214C:

	if r[0] == 0 {
		goto loc_81332174
	}

	r[0] = r[23] | 0x8000
	*(*uint16)(unsafe.Pointer(uintptr(int64(r[31]) + outDiff2))) = uint16(BYTESWAP16(uint16(r[0])))
	r[0] = r[23] | 0x4000
	*(*uint16)(unsafe.Pointer(uintptr(int64(r[31]) + 2 + outDiff2))) = uint16(BYTESWAP16(uint16(r[0])))

	r[31] = r[31] + 4
	r[27] = r[27] + 2
	r[23] = r[23] + 1
	r[22] = r[22] + 1

	goto loc_81332124

loc_81332174:

	r[12] = 9
	r[21] = r[25] + r[12]
	var t = r[21]
	if r[21] > 0x20 {
		goto loc_813321AC
	}

	r[21] = (^(r[12] - 0x20)) + 1
	r[6] = r[26] >> r[21]
	if t == 0x20 {
		goto loc_8133219C
	}

	r[26] = r[26] << r[12]
	r[25] = r[25] + r[12]
	goto loc_813321D0

loc_8133219C:

	r[26] = BYTESWAP32(*(*uint32)(unsafe.Pointer(uintptr(int64(r[4]+r[24]) + inDiff))))
	r[25] = 0
	r[24] = r[24] + 4
	goto loc_813321D0

loc_813321AC:

	r[0] = (^(r[12] - 0x20)) + 1
	r[6] = r[26] >> r[0]
	r[26] = BYTESWAP32(*(*uint32)(unsafe.Pointer(uintptr(int64(r[4]+r[24]) + inDiff))))
	r[0] = (^(r[21] - 0x40)) + 1
	r[24] = r[24] + 4
	r[0] = r[26] >> r[0]
	r[6] = r[6] | r[0]
	r[25] = r[21] - 0x20
	r[26] = r[26] << r[25]

loc_813321D0:

	r[12] = uint32(uint16(BYTESWAP16(uint16(*(*uint16)(unsafe.Pointer(uintptr(int64(r[31])+outDiff2) - 2))))))
	r[31] -= 2
	r[27] = r[27] - 1
	r[0] = r[12] & 0x8000
	r[12] = (r[12] & 0x1FFF) << 1
	if r[0] == 0 {
		goto loc_813321F8
	}

	*(*uint16)(unsafe.Pointer(uintptr(int64(r[9]+r[12]) + outDiff2))) = uint16(BYTESWAP16(uint16(r[6])))
	r[6] = (r[12] & 0x3FFF) >> 1 // extrwi %r6, %r12, 14,17
	if r[27] != 0 {
		goto loc_813321D0
	}

	goto loc_81332204

loc_813321F8:

	*(*uint16)(unsafe.Pointer(uintptr(int64(r[8]+r[12]) + outDiff2))) = uint16(BYTESWAP16(uint16(r[6])))
	r[23] = r[22]
	goto loc_81332124

loc_81332204:

	r[23] = 0x800
	r[22] = 0x800

loc_8133220C:

	if r[29] != 0x1F {
		goto loc_81332228
	}

	r[0] = r[30] >> 31
	r[30] = BYTESWAP32(*(*uint32)(unsafe.Pointer(uintptr(int64(r[4]+r[28]) + inDiff))))
	r[29] = 0
	r[28] = r[28] + 4
	goto loc_81332234

loc_81332228:

	r[0] = r[30] >> 31
	r[29] = r[29] + 1
	r[30] = r[30] << 1

loc_81332234:

	if r[0] == 0 {
		goto loc_8133225C
	}

	r[0] = r[23] | 0x8000
	*(*uint16)(unsafe.Pointer(uintptr(int64(r[31]) + outDiff2))) = uint16(BYTESWAP16(uint16(r[0])))
	r[0] = r[23] | 0x4000
	*(*uint16)(unsafe.Pointer(uintptr(int64(r[31]) + 2 + outDiff2))) = uint16(BYTESWAP16(uint16(r[0])))

	r[31] = r[31] + 4
	r[27] = r[27] + 2
	r[23] = r[23] + 1
	r[22] = r[22] + 1

	goto loc_8133220C

loc_8133225C:

	r[12] = 0xB
	r[21] = r[29] + r[12]
	t = r[21]
	if r[21] > 0x20 {
		goto loc_81332294
	}

	r[21] = (^(r[12] - 0x20)) + 1
	r[7] = r[30] >> r[21]
	if t == 0x20 {
		goto loc_81332284
	}

	r[30] = r[30] << r[12]
	r[29] = r[29] + r[12]
	goto loc_813322B8

loc_81332284:

	r[30] = BYTESWAP32(*(*uint32)(unsafe.Pointer(uintptr(int64(r[4]+r[28]) + inDiff))))
	r[29] = 0
	r[28] = r[28] + 4
	goto loc_813322B8

loc_81332294:

	r[0] = (^(r[12] - 0x20)) + 1
	r[7] = r[30] >> r[0]
	r[30] = BYTESWAP32(*(*uint32)(unsafe.Pointer(uintptr(int64(r[4]+r[28]) + inDiff))))
	r[0] = (^(r[21] - 0x40)) + 1
	r[28] = r[28] + 4
	r[0] = r[30] >> r[0]
	r[7] = r[7] | r[0]
	r[29] = r[21] - 0x20
	r[30] = r[30] << r[29]

loc_813322B8:

	r[12] = uint32(uint16(BYTESWAP16(uint16(*(*uint16)(unsafe.Pointer(uintptr(int64(r[31])+outDiff2) - 2))))))
	r[31] -= 2
	r[27] = r[27] - 1
	r[0] = r[12] & 0x8000
	r[12] = (r[12] & 0x1FFF) << 1
	if r[0] == 0 {
		goto loc_813322E0
	}

	*(*uint16)(unsafe.Pointer(uintptr(int64(r[11]+r[12]) + outDiff2))) = uint16(BYTESWAP16(uint16(r[7])))
	r[7] = (r[12] & 0x3FFF) >> 1 // extrwi %r7, %r12, 14,17
	if r[27] != 0 {
		goto loc_813322B8
	}

	goto loc_813322EC

loc_813322E0:

	*(*uint16)(unsafe.Pointer(uintptr(int64(r[10]+r[12]) + outDiff2))) = uint16(BYTESWAP16(uint16(r[7])))
	r[23] = r[22]
	goto loc_8133220C

loc_813322EC:

	r[0] = r[5]

loc_813322F0:

	r[12] = r[6]

loc_813322F4:

	if r[12] < 0x200 {
		goto loc_8133233C
	}

	if r[25] != 0x1F {
		goto loc_81332318
	}

	r[31] = r[26] >> 31
	r[26] = BYTESWAP32(*(*uint32)(unsafe.Pointer(uintptr(int64(r[4]+r[24]) + inDiff))))
	r[24] = r[24] + 4
	r[25] = 0
	goto loc_81332324

loc_81332318:

	r[31] = r[26] >> 31
	r[25] = r[25] + 1
	r[26] = r[26] << 1

loc_81332324:

	r[27] = r[12] << 1
	if r[31] != 0 {
		goto loc_81332334
	}

	r[12] = uint32(uint16(BYTESWAP16(uint16(*(*uint16)(unsafe.Pointer(uintptr(int64(r[8]+r[27]) + outDiff2)))))))
	goto loc_813322F4

loc_81332334:

	r[12] = uint32(uint16(BYTESWAP16(uint16(*(*uint16)(unsafe.Pointer(uintptr(int64(r[9]+r[27]) + outDiff2)))))))
	goto loc_813322F4

loc_8133233C:

	if r[12] >= 0x100 {
		goto loc_8133235C
	}

	*(*uint8)(unsafe.Pointer(uintptr(int64(r[3]) + outDiff))) = uint8(r[12])
	r[3] = r[3] + 1
	r[5] = r[5] - 1
	if r[5] != 0 {
		goto loc_813322F0
	}

	goto loc_81332434

loc_8133235C:

	r[23] = r[7]

loc_81332360:

	if r[23] < 0x800 {
		goto loc_813323A8
	}

	if r[29] != 0x1F {
		goto loc_81332384
	}

	r[31] = r[30] >> 31
	r[30] = BYTESWAP32(*(*uint32)(unsafe.Pointer(uintptr(int64(r[4]+r[28]) + inDiff))))
	r[28] = r[28] + 4
	r[29] = 0
	goto loc_81332390

loc_81332384:

	r[31] = r[30] >> 31
	r[29] = r[29] + 1
	r[30] = r[30] << 1

loc_81332390:

	r[27] = r[23] << 1
	if r[31] != 0 {
		goto loc_813323A0
	}

	r[23] = uint32(uint16(BYTESWAP16(uint16(*(*uint16)(unsafe.Pointer(uintptr(int64(r[10]+r[27]) + outDiff2)))))))
	goto loc_81332360

loc_813323A0:

	r[23] = uint32(uint16(BYTESWAP16(uint16(*(*uint16)(unsafe.Pointer(uintptr(int64(r[11]+r[27]) + outDiff2)))))))
	goto loc_81332360

loc_813323A8:

	r[12] = r[12] - 0xFD
	r[23] = ^r[23] + r[3] + 1
	r[5] = ^r[12] + r[5] + 1
	r[31] = r[12] >> 3

	if r[31] == 0 {
		goto loc_81332414
	}

	count = r[31]

loc_813323C0:

	r[31] = uint32(*(*uint8)(unsafe.Pointer(uintptr((int64(r[23]) + outDiff) - 1))))
	*(*uint8)(unsafe.Pointer(uintptr(int64(r[3]) + outDiff))) = uint8(r[31])

	r[31] = uint32(*(*uint8)(unsafe.Pointer(uintptr(int64(r[23]) + outDiff))))
	*(*uint8)(unsafe.Pointer(uintptr(int64(r[3]) + 1 + outDiff))) = uint8(r[31])

	r[31] = uint32(*(*uint8)(unsafe.Pointer(uintptr((int64(r[23]) + outDiff) + 1))))
	*(*uint8)(unsafe.Pointer(uintptr(int64(r[3]) + 2 + outDiff))) = uint8(r[31])

	r[31] = uint32(*(*uint8)(unsafe.Pointer(uintptr((int64(r[23]) + outDiff) + 2))))
	*(*uint8)(unsafe.Pointer(uintptr(int64(r[3]) + 3 + outDiff))) = uint8(r[31])

	r[31] = uint32(*(*uint8)(unsafe.Pointer(uintptr((int64(r[23]) + outDiff) + 3))))
	*(*uint8)(unsafe.Pointer(uintptr(int64(r[3]) + 4 + outDiff))) = uint8(r[31])

	r[31] = uint32(*(*uint8)(unsafe.Pointer(uintptr((int64(r[23]) + outDiff) + 4))))
	*(*uint8)(unsafe.Pointer(uintptr(int64(r[3]) + 5 + outDiff))) = uint8(r[31])

	r[31] = uint32(*(*uint8)(unsafe.Pointer(uintptr((int64(r[23]) + outDiff) + 5))))
	*(*uint8)(unsafe.Pointer(uintptr(int64(r[3]) + 6 + outDiff))) = uint8(r[31])

	r[31] = uint32(*(*uint8)(unsafe.Pointer(uintptr((int64(r[23]) + outDiff) + 6))))
	*(*uint8)(unsafe.Pointer(uintptr(int64(r[3]) + 7 + outDiff))) = uint8(r[31])

	r[23] = r[23] + 8
	r[3] = r[3] + 8

	count--

	if count != 0 {
		goto loc_813323C0
	}

	r[12] = r[12] & 7
	if r[12] == 0 {
		goto loc_8133242C
	}
loc_81332414:

	count = r[12]

loc_81332418:

	r[31] = uint32(*(*uint8)(unsafe.Pointer(uintptr(int64(r[23]) + outDiff - 1))))
	r[23] = r[23] + 1
	*(*uint8)(unsafe.Pointer(uintptr(int64(r[3]) + outDiff))) = uint8(r[31])
	r[3] = r[3] + 1

	count--

	if count != 0 {
		goto loc_81332418
	}

loc_8133242C:

	if r[5] != 0 {
		goto loc_813322F0
	}
loc_81332434:

	r[3] = r[0]

	var ret []uint8

	for i := uint32(0); i < r[3]; i += uint32(unsafe.Sizeof(uint8(0))) {
		ret = append(ret, *(*uint8)(unsafe.Pointer(uintptr(int64(o) + outDiff + int64(i)))))
	}

	return ret
}

func main() {
	fmt.Printf("Language: Go\n")
	var filepath = "test.ash"
	fmt.Printf("Checking file \"%s\"\n", filepath)

	var data = loadAshFileIntoSlice(filepath)

	fmt.Printf("isAshCompressed: %t\n", isAshCompressed(data))

	ioutil.WriteFile("out.ash", decompressAsh(data), 0644)

	if deepCompare("out.ash", "out_real.ash") {
		fmt.Println("Decompression successful")
	} else {
		fmt.Println("Decompression unsuccessful")
	}
}

// Stolen from stackoverflow to compare files
const chunkSize = 64000

func deepCompare(file1, file2 string) bool {
	// Check file size ...

	f1, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}

	f2, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				log.Fatal(err1, err2)
			}
		}

		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}
