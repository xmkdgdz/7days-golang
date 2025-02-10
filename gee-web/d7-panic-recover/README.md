# 错误恢复(Panic Recover)

## defer

panic 会导致程序被中止，但是在退出前，会先处理完当前协程上已经 defer 的任务，执行完成后再退出。

defer 的任务执行完成后，panic 会继续被抛出，导致程序非正常结束。

## recover

recover 函数可以避免因为 panic 发生而导致整个程序终止，recover 函数只在 defer 中生效。

```go
func test_recover() {
	defer func() {
		fmt.Println("defer func")
		if err := recover(); err != nil {
			fmt.Println("recover success")
		}
	}()

	arr := []int{1, 2, 3}
	fmt.Println(arr[4])
	fmt.Println("after panic")
}

func main() {
	test_recover()
	fmt.Println("after recover")
}
```

结果：

```console
defer func
recover success
after recover
```
