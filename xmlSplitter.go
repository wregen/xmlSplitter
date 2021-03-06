package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"strings"
	"time"
)

const (
	directoryName string = "xmlFiles"
)

var (
	noElements = flag.Int("n", 5, "Numri i elementeve per file")
)

type subDocument struct {
	Name    string `xml:"name"`
	Surname string `xml:"surname"`
	Age     int    `xml:"age"`
}

type document struct {
	XMLName      xml.Name      `xml:"document"`
	SubDocuments []subDocument `xml:"subDocument"`
}

func main() {
	flag.Parse()
	fileName, err := getFirstXMLFileFound()
	if err != nil {
		writeToLogFile(err.Error())
		return
	}
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		writeToLogFile(err.Error())
		return
	}
	var xmlData document
	err = xml.Unmarshal(b, &xmlData)
	if err != nil {
		writeToLogFile(err.Error())
		return
	}
	err = createSplittedXMLFiles(&xmlData, *noElements)
	if err != nil {
		writeToLogFile(err.Error())
		return
	}
	writeToLogFile("Ndarja e XML u krye me sukses")
}

func writeToLogFile(log string) {
	content := fmt.Sprintf("%s: %s\n", time.Now(), log)
	f, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
}

func getFirstXMLFileFound() (string, error) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return "", err
	}
	for _, file := range files {

		if strings.Contains(strings.ToUpper(file.Name()), ".XML") {
			return file.Name(), nil
		}
	}
	return "", errors.New("Nuk u gjet asnje file xml ne direktorine ku ekzekutohet programi")
}

func createSplittedXMLFiles(doc *document, xmlNumerOfElementsPerFile int) error {
	if xmlNumerOfElementsPerFile <= 0 {
		return errors.New("Numri i elementëve për file nuk mund të jetë më i vogël ose i barabartë me zero")
	}
	exists := func(name string) bool {
		if _, err := os.Stat(name); err != nil {
			if os.IsNotExist(err) {
				return false
			}
		}
		return true
	}
	if exists(directoryName) {
		err := os.RemoveAll(directoryName)
		if err != nil {
			return err
		}
	}
	err := os.Mkdir(directoryName, 0777)
	if err != nil {
		return err
	}
	div := float64(len(doc.SubDocuments)) / float64(xmlNumerOfElementsPerFile)
	nrFiles := int(math.Ceil(div))
	for i := 0; i < nrFiles; i++ {
		var fileContent document
		if i == nrFiles-1 {
			fileContent = document{
				XMLName:      doc.XMLName,
				SubDocuments: doc.SubDocuments[i*xmlNumerOfElementsPerFile:],
			}
		} else {
			fileContent = document{
				XMLName:      doc.XMLName,
				SubDocuments: doc.SubDocuments[i*xmlNumerOfElementsPerFile : (i+1)*xmlNumerOfElementsPerFile],
			}
		}
		cont, err := xml.MarshalIndent(fileContent, "", "   ")
		if err != nil {
			return err
		}
		splitFileName := path.Join(directoryName, fmt.Sprintf("splitFile%d.xml", i))
		content := append([]byte(xml.Header), cont...)
		ioutil.WriteFile(splitFileName, content, 0644)
	}
	return nil
}
