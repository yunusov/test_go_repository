package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	//conv()
	gs()
}

func gs() {
	/*
		Задание 2. Graceful shutdown
	*/
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 1
		for {
			time.Sleep(time.Second)
			i++
			select {
			case <-ctx.Done():
				fmt.Println("Выхожу из программы... ")
				return
			default:
				fmt.Println(i * 2)
			}
		}
	}()
	
	go initGracefulShutdown(stop, &wg)
	wg.Wait()
	fmt.Println("End")
}

func initGracefulShutdown(cancelFunc context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs
	fmt.Println("Graceful Shutdown")
	cancelFunc()
}

func action(cancelCtx context.Context) error {
	time.Sleep(time.Second)
	return errors.New("failed")
}

/* ----- */

func conv() {
	/*
		Задание 2. Конвейер
	*/
	for {
		n, isStop := readInt()
		if isStop {
			break
		}
		var wg sync.WaitGroup
		wg.Add(2)
		nChan := make(chan int)
		go func() {
			nChan <- n
			close(nChan)
		}()
		sqrChan := calcSqr(nChan, &wg)
		multCh := multipleTwo(sqrChan, &wg)
		wg.Wait()
		printResult(<-multCh)
	}
}

func readInt() (n int, isStop bool) {
	for {
		fmt.Print("Ввод: ")
		reader := bufio.NewReader(os.Stdin)
		result, err := reader.ReadString('\r')
		if err != nil {
			panic(err)
		}
		if result == "стоп\r" {
			n = 0
			isStop = true
			return
		}
		result = strings.Replace(result, "\r", "", 1)
		n, err = strconv.Atoi(result)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		return
	}
}

func calcSqr(ch chan int, wg *sync.WaitGroup) chan int {
	defer wg.Done()
	resultChan := make(chan int)
	go func() {
		defer close(resultChan)
		n := <-ch
		n *= n
		fmt.Println("Квадрат:", n)
		resultChan <- n
	}()
	return resultChan
}

func multipleTwo(ch chan int, wg *sync.WaitGroup) chan int {
	defer wg.Done()
	resultChan := make(chan int)
	go func() {
		defer close(resultChan)
		n := <-ch
		n *= 2
		resultChan <- n
	}()
	return resultChan
}

func printResult(n int) {
	fmt.Println("Произведение:", n)
}
