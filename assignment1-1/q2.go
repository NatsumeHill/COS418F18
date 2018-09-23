package cos418_hw1_1

import (
	"bufio"
	"io"
	"os"
	"strconv"
)

// Sum numbers from channel `nums` and output sum to `out`.
// You should only output to `out` once.
// Do NOT modify function signature.
func sumWorker(nums chan int, out chan int) {
	// TODO: implement me
	// HINT: use for loop over `nums`
	total := 0
	for num := range nums {
		total += num
	}
	out <- total
	println(">>>>>>>>>>>>>>>sum worker")
	println(total)
}

// Read integers from the file `fileName` and return sum of all values.
// This function must launch `num` go routines running
// `sumWorker` to find the sum of the values concurrently.
// You should use `checkError` to handle potential errors.
// Do NOT modify function signature.
type chunk struct {
	bufsize int
	offset  int64
}

func sum(num int, fileName string) int {
	// TODO: implement me
	// HINT: use `readInts` and `sumWorkers`
	// HINT: used buffered channels for splitting numbers between workers
	send := make(chan int, num)
	recv := make(chan int)
	for i := 0; i < num; i++ {
		go sumWorker(send, recv)
	}
	// 并发读取文件
	const BufferSize = 100
	f, err := os.Open(fileName)
	checkError(err)
	defer f.Close()
	// finfo, err := f.Stat()
	// checkError(err)
	// fileSize := int(finfo.Size())
	// // Number of go routines we need to spawn.
	// concurrency := fileSize / BufferSize
	// // buffer sizes that each of the go routine below should use. ReadAt
	// // returns an error if the buffer size is larger than the bytes returned
	// // from the file.
	// chunksizes := make([]chunk, concurrency)
	// // All buffer sizes are the same in the normal case. Offsets depend on the
	// // index. Second go routine should start at 100, for example, given our
	// // buffer size of 100.
	// for i := 0; i < concurrency; i++ {
	// 	chunksizes[i].bufsize = BufferSize
	// 	chunksizes[i].offset = int64(BufferSize * i)
	// }
	// // check for any left over bytes. Add the residual number of bytes as the
	// // the last chunk size.
	// if remainder := fileSize % BufferSize; remainder != 0 {
	// 	c := chunk{bufsize: remainder, offset: int64(concurrency * BufferSize)}
	// 	concurrency++
	// 	chunksizes = append(chunksizes, c)
	// }
	// var wg sync.WaitGroup
	// wg.Add(concurrency)
	// stringFromFile := ""
	// for i := 0; i < concurrency; i++ {
	// 	go func(chunksizes []chunk, i int) {
	// 		defer wg.Done()
	// 		chunk := chunksizes[i]
	// 		buffer := make([]byte, chunk.bufsize)
	// 		bytesread, err := f.ReadAt(buffer, chunk.offset)
	// 		checkError(err)
	// 		// rd := bufio.NewReader(strings.NewReader(string(buffer[:bytesread])))
	// 		// checkNum, err := readInts(rd)
	// 		// if err != nil && err != io.EOF {
	// 		// 	fmt.Println(err)
	// 		// 	return
	// 		// }
	// 		// for _, tosend := range checkNum {
	// 		// 	send <- tosend
	// 		// }
	// 		// fmt.Printf("%v\n", checkNum)
	// 		stringFromFile += string(buffer[:bytesread])
	// 	}(chunksizes, i)
	// }
	// wg.Wait()
	rd := bufio.NewReader(f)
	checkNum, err := readInts(rd)
	// fmt.Printf("%v\n", checkNum)
	checkError(err)
	for _, tosend := range checkNum {
		send <- tosend
	}
	close(send)
	total := 0
	for i := 0; i < num; i++ {
		total += <-recv
	}
	return total
}

// Read a list of integers separated by whitespace from `r`.
// Return the integers successfully read with no error, or
// an empty slice of integers and the error that occurred.
// Do NOT modify this function.
func readInts(r io.Reader) ([]int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var elems []int
	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return elems, err
		}
		elems = append(elems, val)
	}
	return elems, nil
}
