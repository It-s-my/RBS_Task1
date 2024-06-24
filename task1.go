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

func main() {
 fileName := flag.String("file", "", "имя файла")
 flag.Parse()

 if *fileName == "" {
  fmt.Println("Необходимо указать имя файла")
  return
 }

 file, read := ioutil.ReadFile(*fileName)
 if read != nil {
  fmt.Println("Ошибка чтения файла:", read)
  return
 }

 content := string(file)
 lines := strings.Split(content, "\n")

 for _, line := range lines {
  u, read := url.Parse(line)
  if read != nil || u.Scheme == "" || u.Host == "" {
   continue
  }

  resp, read := http.Get(line)
  if read != nil {
   fmt.Println(line, ":", read)
   continue
  }
  defer resp.Body.Close()

  body, read := ioutil.ReadAll(resp.Body)
  if read != nil {
   fmt.Println(line, ":", read)
   continue
  }

  outputFileName := strings.Replace(line, "://", "_", -1) + ".txt"
  outputFile, create := os.Create(outputFileName)
  if create != nil {
   fmt.Println(line, ":", create)
   continue
  }
  defer outputFile.Close()

  outputFile.Write(body)

  fmt.Println(line, ":", "Результат сохранен в файл", outputFileName)
 }
}
