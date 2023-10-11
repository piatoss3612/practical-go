package srp

import (
	"fmt"
	"os"
	"strings"
)

var ReportFormat = "title: %s\ncontent: %s"

// Report는 보고서의 내용을 담는 구조체
type Report struct {
	Title   string
	Content string
}

// CreateReport는 새로운 보고서를 생성함.
func NewReport(title, content string) Report {
	return Report{
		Title:   title,
		Content: content,
	}
}

// FormatReport는 보고서를 형식에 맞게 출력함.
func (r *Report) FormatReport() string {
	return fmt.Sprintf(ReportFormat, r.Title, r.Content)
}

type ReportSaver struct{} // ReportSaver는 보고서를 파일에 저장하는 구조체

// SaveToFile는 보고서를 파일에 저장함.
func (rs *ReportSaver) SaveToFile(r Report, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fileContent := fmt.Sprintf("%s\n%s", r.Title, r.Content)

	_, err = file.WriteString(fileContent)
	if err != nil {
		return err
	}

	return nil
}

type ReportLoader struct{} // ReportLoader는 보고서를 파일에서 읽어오는 구조체

// LoadFromFile는 파일에서 보고서를 읽어옴.
func (rl *ReportLoader) LoadFromFile(filename string) (Report, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Report{}, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return Report{}, err
	}

	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		return Report{}, err
	}

	fileContent := string(bs)

	lines := strings.Split(fileContent, "\n")

	if len(lines) != 2 {
		return Report{}, fmt.Errorf("invalid format")
	}

	title := lines[0]
	content := lines[1]

	return NewReport(title, content), nil
}

// package main

// import (
// 	"fmt"
// 	"solid/srp"
// )

// func main() {
// 	title := "초전도치의 비밀"
// 	content := "초전도치는 비밀이다."

// 	report1 := srp.NewReport(title, content)

// 	saver := srp.ReportSaver{}

// 	err := saver.SaveToFile(report1, "report.txt")
// 	if err != nil {
// 		panic(err)
// 	}

// 	loader := srp.ReportLoader{}

// 	_, err = loader.LoadFromFile("report1.txt")
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(report1.Title == title)
// 	fmt.Println(report1.Content == content)
// }
