package interaction

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func Message(
	message string,
) {
	println(fmt.Sprintf("🤖    %s", message))
}

func Error(
	message string,
) {
	println(fmt.Sprintf("💣    %s", message))
}

func Mark(skip int) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	_, file, line, _ := runtime.Caller(skip)
	println(fmt.Sprintf("📝   %s@%d", strings.TrimPrefix(file, wd), line))
}

func Choose(
	message string,
	count int,
	option func(index int) (string, string),
) string {
	shouldBreakLine := false

	Message(message)
	keys := []string{}
	values := []string{}
	for index := 0; index < count; index++ {
		key, value := option(index)
		keys = append(keys, key)
		values = append(values, value)
	}

	var in string
	for {
		if shouldBreakLine {
			println("👇    ")
		} else {
			print("👉    ")
		}
		for index, key := range keys {
			output := fmt.Sprintf("%s:%s    ", key, values[index])
			if shouldBreakLine {
				println(output)
			} else {
				print(output)
			}
		}
		println()

		fmt.Scan(&in)
		var found = false
		for _, key := range keys {
			if key == in {
				found = true
			}
		}
		if found {
			break
		}
		println("もう一度入力してください。")
	}

	return in
}
