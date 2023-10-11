package main

import (
	"fmt"
	"solid/srp"
)

func main() {
	title := "초전도치의 비밀"
	content := "초전도치는 비밀이다."

	report1 := srp.NewReport(title, content)

	saver := srp.ReportSaver{}

	err := saver.SaveToFile(report1, "report.txt")
	if err != nil {
		panic(err)
	}

	loader := srp.ReportLoader{}

	_, err = loader.LoadFromFile("report1.txt")
	if err != nil {
		panic(err)
	}

	fmt.Println(report1.Title == title)
	fmt.Println(report1.Content == content)
}
