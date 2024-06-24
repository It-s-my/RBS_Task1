package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "  %s -src <source_file> -dst <destination_directory>\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	fileName := flag.String("src", "", "Имя файла")                  //флаг для имени файла
	dstPtr := flag.String("dst", "", "Название конечной директории") //флаг для названия конечной директории
	flag.Parse()

	if *fileName == "" {
		fmt.Println("Необходимо указать имя файла") //проверка на наличие имени файла
		return
	}

	if *dstPtr == "" {
		fmt.Println("Необходимо указать название конечной директории") //проверка на наличие названия конечной директории
		return
	}

	file, err := os.ReadFile(*fileName)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err) //проверка чтения файла
		return
	}

	err = os.MkdirAll(*dstPtr, 0777)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	fmt.Println("Directory", *dstPtr, "created successfully")
	//создание директории

	content := string(file)               //получение содержимого файла
	lines := strings.Split(content, "\n") //разбиение содержимого файла на строки

	for _, line := range lines { //Проверка на корректность ссылок, неправильные сразу отбрасываем
		u, read := url.Parse(line)
		if read != nil && u.Scheme == "" && u.Host == "" {
			continue
		}

		resp, read := http.Get(line) //Отправляем гет запросы оставшимся ссылкам
		if read != nil {
			fmt.Println(line, ":", read)
			continue
		}
		defer resp.Body.Close() //закрытие потока

		body, read := ioutil.ReadAll(resp.Body) //чтение ответа
		if read != nil {
			fmt.Println(line, ":", read)
			continue
		}
		//Сохранение результатов в файлы
		outputFileName := *dstPtr + "/" + strings.Replace(line, "://", "_", -1) + ".txt"
		outputFile, create := os.Create(outputFileName)
		if create != nil {
			fmt.Println(line, ":", create)
			continue
		}
		defer outputFile.Close() //закрытие потока

		outputFile.Write(body)

		fmt.Println(line, ":", "Результат сохранен в файл", outputFileName) //Отчет о результате работы
	}
}
