package cmd

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"sortClone/internal/model"
	"sortClone/internal/reader"
	"sortClone/internal/utils"
)

var (
	bigFilePath         string
	bigFileLinesCount   int64
	bigFileSize         int64
	smallFilePath       string
	smallFileLinesCount int64 = 1000
	smallFileSize       int64
	first10linesBigFile = []string{}
	first10linesSorted  = []string{}
)

func TestStdIn(t *testing.T) {
}

func TestBigFileNumericReverseSortFirstColumn(t *testing.T) {
	os.Args = []string{
		os.Args[0],
		"-n", // числовая сортировка
		"-r", // обратный порядок
		"-k",
		"1",         // по первой колонке - она числовая
		bigFilePath, // имя файла
	}
	buf := &bytes.Buffer{}
	OutputDST = buf

	// выполняем саму сортировку с отсчетом времени
	start := time.Now()
	Execute()
	end := time.Since(start)

	// проверка результата
	sortedLines := strings.Split(buf.String(), "\n")
	sortedLines = sortedLines[:len(sortedLines)-1]
	isSorted := true
	for i := range 5 {
		first10linesSorted = append(first10linesSorted, sortedLines[i]+"\n")
	}

	for i := 1; i < len(sortedLines); i++ {
		a, erra := strconv.Atoi(strings.Split(sortedLines[i-1], "\t")[0])
		b, errb := strconv.Atoi(strings.Split(sortedLines[i], "\t")[0])
		if erra != nil || errb != nil {
			t.Fatalf("Failed to convert lines %d or %d to int: %v,%v", i-1, i, erra, errb)
		}
		if a < b {
			isSorted = false
			break
		}
	}

	// подводим итоги
	t.Logf("Running ExecuteSort took %vmsec.", end.Milliseconds())
	t.Logf("Generated testfile size: %vmb", bigFileSize/1024/1024)
	if bigFileLinesCount != int64(len(sortedLines)) {
		t.Logf("Lines N in BIG file: %v", bigFileLinesCount)
		t.Logf("Lines N in sorted output: %v", len(sortedLines))
		t.Fatalf("Discrepancy in input and output lines count!")

	}
	if !isSorted {
		t.Logf("First 10 lines from BIG file:\n %v", first10linesBigFile)
		t.Logf("First 10 lines from sort result:\n %v", first10linesSorted)
		t.Fatalf("Output is not sorted!")
	}
}

func TestMain(m *testing.M) {
	// генерим тестовый файл с мокамии
	err := prepareTestFiles()
	if err != nil {
		log.Printf("Failed to generate mock-file: %v", err)
		os.Exit(1)
	}

	code := m.Run()

	// чистим хлам после тестов
	if err := os.Remove(bigFilePath); err != nil {
		log.Printf("Failed to remove original BIG mock-file: %v", err)
		os.Exit(1)
	}
	if err := os.Remove(smallFilePath); err != nil {
		log.Printf("Failed to remove original SMALL mock-file: %v", err)
		os.Exit(1)
	}
	if model.OptsContainer.WriteToFile != "" {
		err := os.Remove(model.OptsContainer.WriteToFile)
		if err != nil {
			log.Printf("Failed to remove DST mock-file: %v", err)
			os.Exit(1)
		}
	}

	os.Exit(code)
}

func prepareTestFiles() error {
	var lineSize int64 = 24
	rand.New(rand.NewSource(lineSize))

	// создаём временный файл
	tmpBigFile, err := os.CreateTemp("", "bigfile-*.txt")
	if err != nil {
		return err
	}
	tmpSmallFile, err := os.CreateTemp("", "smallfile-*.txt")
	if err != nil {
		return err
	}

	// генерируем строки
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZабвгдежзиклмнопрстуфхцчшщъыьэюяАБВГДЕЖЗИКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ")
	numbers := []rune("0123456789")

	fileSize := reader.MaxFileSize + reader.MaxFileSize*2/10 // делаем размер 120% от порога дробления файла на врем. файлы
	human := []string{}
	for k := range utils.Multipliers { // грузим актуальные суффиксы из карты основной программы
		human = append(human, k)
	}

	for {
		bigFileLinesCount++
		b := make([]rune, lineSize/3)
		c := make([]rune, lineSize/3)
		h := human[rand.Intn(len(human))]

		for l := range b {
			b[l] = letters[rand.Intn(len(letters))]
			c[l] = numbers[rand.Intn(len(numbers))]
		}

		generatedLine := fmt.Sprintf("%d\t%s\t%s\t%s%s\n", bigFileLinesCount, string(b), string(c), string(c), h)

		if _, err := fmt.Fprint(tmpBigFile, generatedLine); err != nil {
			return err
		}

		if bigFileLinesCount <= smallFileLinesCount {
			if _, err := fmt.Fprint(tmpSmallFile, generatedLine); err != nil {
				return err
			}
		}

		if bigFileLinesCount <= 5 {
			first10linesBigFile = append(first10linesBigFile, generatedLine)
		}

		if info, _ := tmpBigFile.Stat(); info.Size() > fileSize {
			bigFileSize = info.Size()
			break
		}
	}

	info, _ := tmpSmallFile.Stat()
	smallFileSize = info.Size()

	defer tmpBigFile.Close()
	defer tmpSmallFile.Close()

	buf := &bytes.Buffer{}
	OutputDST = buf
	bigFilePath = tmpBigFile.Name()
	smallFilePath = tmpSmallFile.Name()

	return nil
}
