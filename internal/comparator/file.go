package comparator

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

const chunkSize = 64 * 1024 // 64 KB — optimal for disk I/O and CPU cache

// Result holds the outcome of a file comparison.
type Result struct {
	Identical  bool  // true if the files are byte-for-byte identical
	DiffOffset int64 // byte offset of first difference; -1 if identical or sizes differ
	SizeA      int64 // size of file A in bytes
	SizeB      int64 // size of file B in bytes
}

// CompareFiles compares two files byte-by-byte and returns a Result.
// Returns an error if either file cannot be stat'd, opened, or read.
//
// Algorithm:
//  1. os.Stat() both files — early exit if sizes differ (zero I/O)
//  2. Both empty → IDENTICAL (zero I/O)
//  3. Read in 64 KB chunks via bufio.Reader
//  4. bytes.Equal() per chunk — first mismatch triggers early exit with exact DiffOffset
func CompareFiles(pathA, pathB string) (Result, error) {
	infoA, err := os.Stat(pathA)
	if err != nil {
		return Result{}, err
	}
	infoB, err := os.Stat(pathB)
	if err != nil {
		return Result{}, err
	}

	sizeA, sizeB := infoA.Size(), infoB.Size()

	// Early exit: different sizes → DIFFERENT (zero I/O)
	if sizeA != sizeB {
		return Result{Identical: false, DiffOffset: -1, SizeA: sizeA, SizeB: sizeB}, nil
	}

	// Both empty → IDENTICAL
	if sizeA == 0 {
		return Result{Identical: true, DiffOffset: -1, SizeA: 0, SizeB: 0}, nil
	}

	fA, err := os.Open(pathA)
	if err != nil {
		return Result{}, err
	}
	defer fA.Close()

	fB, err := os.Open(pathB)
	if err != nil {
		return Result{}, err
	}
	defer fB.Close()

	readerA := bufio.NewReaderSize(fA, chunkSize)
	readerB := bufio.NewReaderSize(fB, chunkSize)

	bufA := make([]byte, chunkSize)
	bufB := make([]byte, chunkSize)
	var offset int64

	for {
		nA, errA := io.ReadFull(readerA, bufA)
		nB, errB := io.ReadFull(readerB, bufB)

		// Since sizes are equal nA should always equal nB; handle defensively.
		n := nA
		if nB < n {
			n = nB
		}

		if n > 0 && !bytes.Equal(bufA[:n], bufB[:n]) {
			// Locate the exact byte offset of the first difference.
			for i := 0; i < n; i++ {
				if bufA[i] != bufB[i] {
					return Result{
						Identical:  false,
						DiffOffset: offset + int64(i),
						SizeA:      sizeA,
						SizeB:      sizeB,
					}, nil
				}
			}
		}

		offset += int64(n)

		// io.ErrUnexpectedEOF: read n < chunkSize but n > 0 (last chunk).
		// io.EOF: nothing left to read.
		if errA == io.EOF || errA == io.ErrUnexpectedEOF {
			break
		}
		if errA != nil {
			return Result{}, errA
		}
		if errB != nil {
			return Result{}, errB
		}
	}

	return Result{Identical: true, DiffOffset: -1, SizeA: sizeA, SizeB: sizeB}, nil
}
