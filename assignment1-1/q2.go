package cos418_hw1_1

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
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
}

// Read integers from the file `fileName` and return sum of all values.
// This function must launch `num` go routines running
// `sumWorker` to find the sum of the values concurrently.
// You should use `checkError` to handle potential errors.
// Do NOT modify function signature.
func sum(num int, fileName string) int {
	// TODO: implement me
	// HINT: use `readInts` and `sumWorkers`
	// HINT: used buffered channels for splitting numbers between workers
	send := make(chan int, num)
	recv := make(chan int, num)
	for i := 0; i < num; i++ {
		go sumWorker(send, recv)
	}
	// 分块读取文件
	const BufferSize = 100
	file, err := os.Open(fileName)
	checkError(err)
	defer file.Close()
	lastNum := ""
	currNums := ""
	buffer := make([]byte, BufferSize)
	// 创建"readInts"协同进程
	var wg sync.WaitGroup
	for {
		bytesread, err := file.Read(buffer)
		// err value can be io.EOF, which means that we reached the end of
		// file, and we have to terminate the loop. Note the fmt.Println lines
		// will get executed for the last chunk because the io.EOF gets
		// returned from the Read function only on the *next* iteration, and
		// the bytes returned will be 0 on that read.
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		if i := bytes.LastIndexByte(buffer, '\n'); i < bytesread-1 {
			currNums = lastNum + string(buffer[:i])
			lastNum = string(buffer[i:bytesread])
		} else {
			currNums = lastNum + string(buffer[:bytesread])
			lastNum = ""
		}
		wg.Add(1)
		go func(source string) {
			rd := strings.NewReader(source)
			nums, _ := readInts(rd)
			for _, tosend := range nums {
				send <- tosend
			}
			wg.Done()
		}(currNums)
	}
	// 等所有"readInts"进程退出之后，关闭channel，防止阻塞
	wg.Wait()
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
