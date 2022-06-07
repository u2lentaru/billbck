package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func XMLSplitter() {

	var fw *os.File
	const line_count = 10000000 //Размер файла (строк)
	const StartTag = "<InformationRegisterRecordSet.СреднееПотребление>"
	const StopTag = "</InformationRegisterRecordSet.СреднееПотребление>"
	const sourcename = "SP.xml"
	const filename = "sp_out"

	t := time.Now()
	fmt.Println("Started: " + t.Format("2006-01-02 15:04:05"))

	f, err := os.Open(sourcename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if checkFileIsExist(filename + ".xml") { // Если файл существует
		// fw, _ = os.OpenFile(filename, os.O_APPEND, 0666) // Открыть файл
		os.Remove(filename)
		fw, _ = os.Create(filename + ".xml") // Создать файл
		fmt.Println("Файл " + filename + ".xml" + " пересоздан")
	} else {
		fw, _ = os.Create(filename + ".xml") // Создать файл
		fmt.Println("Файл " + filename + ".xml" + " создан")
	}

	fl := false
	lc := 3
	fc := 1 //fle count

	// Чтение файла с ридером
	sc := bufio.NewScanner(f)

	w := bufio.NewWriter(fw) // Создаем новый объект Writer

	_, _ = w.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	_, _ = w.WriteString(" <V8Exch:_1CV8DtUD xmlns:V8Exch=\"http://www.1c.ru/V8/1CV8DtUD/\" xmlns:v8=\"http://v8.1c.ru/data\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">\n")
	_, _ = w.WriteString(" 	<V8Exch:Data>\n")

	for sc.Scan() {
		if strings.Contains(sc.Text(), StartTag) {
			// if strings.Contains(sc.Text(), "<CatalogObject.ТочкиУчета>") {
			fl = true
		}

		if fl {
			// fmt.Println(sc.Text())

			_, _ = w.WriteString(sc.Text() + "\n")

			lc++
		}

		if strings.Contains(sc.Text(), StopTag) {
			// if strings.Contains(sc.Text(), "</CatalogObject.ТочкиУчета>") {
			fl = false

			if lc > line_count {
				_, _ = w.WriteString(" 	</V8Exch:Data>\n")
				_, _ = w.WriteString(" </V8Exch:_1CV8DtUD>\n")

				w.Flush()
				fw.Close()

				if checkFileIsExist(filename + strconv.Itoa(fc) + ".xml") {
					os.Remove(filename + strconv.Itoa(fc) + ".xml")
					fw, _ = os.Create(filename + strconv.Itoa(fc) + ".xml")
					fmt.Println("Файл " + filename + strconv.Itoa(fc) + ".xml" + " пересоздан")
				} else {
					fw, _ = os.Create(filename + strconv.Itoa(fc) + ".xml")
					fmt.Println("Файл " + filename + strconv.Itoa(fc) + ".xml" + " создан")
				}

				fc++

				w = bufio.NewWriter(fw)

				_, _ = w.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
				_, _ = w.WriteString(" <V8Exch:_1CV8DtUD xmlns:V8Exch=\"http://www.1c.ru/V8/1CV8DtUD/\" xmlns:v8=\"http://v8.1c.ru/data\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">\n")
				_, _ = w.WriteString(" 	<V8Exch:Data>\n")

				lc = 3

			}
		}

	}

	_, _ = w.WriteString(" 	</V8Exch:Data>\n")
	_, _ = w.WriteString(" </V8Exch:_1CV8DtUD>\n")

	w.Flush()
	fw.Close()

	t = time.Now()
	fmt.Println("Stopped: " + t.Format("2006-01-02 15:04:05"))

}

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
