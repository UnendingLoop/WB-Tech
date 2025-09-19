package integration

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"sortClone/cmd"
	"sortClone/internal/model"
	"sortClone/internal/reader"
	"sortClone/internal/utils"
)

var (
	bigFilePath         string
	bigFileLinesCount   int64
	bigFileSize         int64
	smallFilePath       string
	smallFileLinesCount int64 = 1000 // кол-во строк для малого тестового файла
	smallFileSize       int64
	first10linesBigFile = []string{}
	first10linesSorted  = []string{}
	dstFileName         string
)

func TestStdInGroupedFlagsHumanSort(t *testing.T) {
	os.Args = []string{
		os.Args[0],
		"-ruHk", // группируем флаги
		"2",
	}

	input := "000\t1Tb\n000\t1Тб\n111\t1Tb\n222\t1Гб\n333\t1Mb\n444\t1Кб\n555\t1К\n666\t\n000\t1Тб\n777\t\n"
	expOutput := "000\t1Тб\n222\t1Гб\n333\t1Mb\n555\t1К\n777\t\n"

	// подменяем стдин
	oldStdin := os.Stdin                   // сохраняем оригинал
	defer func() { os.Stdin = oldStdin }() // восстанавливаем при окончании работы

	r, w, _ := os.Pipe()
	_, _ = w.Write([]byte(input))
	_ = w.Close()
	os.Stdin = r // подменяем stdin

	// перехватываем вывод
	buf := &bytes.Buffer{}
	cmd.OutputDST = buf
	defer func() { cmd.OutputDST = nil }()

	// выполняем саму сортировку с отсчетом времени
	start := time.Now()
	cmd.Execute()
	end := time.Since(start)

	// подводим итоги
	t.Logf("Running ExecuteSort took %vmsec.", end.Milliseconds())
	t.Logf("Generated testfile size: %vmb", bigFileSize/1024/1024)

	if buf.String() != expOutput {
		t.Logf("Expected output is:\n %v", expOutput)
		t.Logf("Actual output is:\n %v", buf.String())
		t.Fatalf("Output is not sorted!")
	}
}

// тут сделать вывод в файл
func TestSmallFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "TestTMPdir-*")
	if err != nil {
		t.Logf("Failed to pre-create TMP dir: %v", err)
	}

	dstFileName = path.Join(tmpDir, "testTMPresult.txt")

	os.Args = []string{
		os.Args[0],
		"-rn",
		"-k",
		"1",
		"-o",
		dstFileName,
		smallFilePath,
	}

	// выполняем саму сортировку с отсчетом времени
	start := time.Now()
	cmd.Execute()
	end := time.Since(start)

	// Cравниваем результаты
	sizeOK, linesOK := true, true

	destFile, errf := os.Open(dstFileName)
	defer func() {
		err := destFile.Close()
		if err != nil {
			t.Logf("Failed to close output file: %v", err)
		}
	}()

	if errf != nil {
		t.Fatalf("Failed to open DST file after sorting: %v", errf)
	}

	destInfo, erri := destFile.Stat()
	if erri != nil {
		t.Fatalf("Failed to fetch DST file info after opening: %v", erri)
	}

	if destInfo.Size() != smallFileSize {
		sizeOK = false
	}

	// считаем кол-во строк в выходном файле
	scanner := bufio.NewScanner(destFile)
	scanner.Buffer(make([]byte, 0, 1000), 1010)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if len(lines) != int(smallFileLinesCount) {
		linesOK = false
	}
	// подводим итоги
	t.Logf("Running ExecuteSort took %vmsec.", end.Milliseconds())
	t.Logf("Generated testfile size: %vb", smallFileSize)
	if !sizeOK {
		t.Logf("Size of output file: %vb", destInfo.Size())
		t.Fatalf("Discrepancy in input and output file sizes!")

	}
	if !linesOK {
		t.Logf("N of lines in input file:\n %v", smallFileLinesCount)
		t.Logf("N of lines in output file:\n %v", len(lines))
		t.Fatalf("Discrepancy in lines count!")
	}
}

func TestBigFileNumericReverseSortFirstColumn(t *testing.T) {
	// подсовываем аргументы для обработки в Cobra
	os.Args = []string{
		os.Args[0],
		"-n", // числовая сортировка
		"-r", // обратный порядок
		"-k",
		"1",         // по первой колонке - она числовая
		bigFilePath, // имя файла
	}

	// перехватываем вывод
	buf := &bytes.Buffer{}
	cmd.OutputDST = buf
	defer func() { cmd.OutputDST = nil }()

	// выполняем саму сортировку с отсчетом времени
	start := time.Now()
	cmd.Execute()
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
		log.Printf("Failed to generate mock-files: %v", err)
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
		err := os.Remove(dstFileName)
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

	maxFileSize := reader.MaxFileSize + reader.MaxFileSize*2/10 // делаем размер 120% от порога дробления файла на врем. файлы
	human := []string{}
	for k := range utils.Multipliers { // грузим актуальные суффиксы из карты основной программы
		human = append(human, strings.ToLower(k))
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

		if info, _ := tmpBigFile.Stat(); info.Size() > maxFileSize {
			bigFileSize = info.Size()
			break
		}
	}

	info, _ := tmpSmallFile.Stat()
	smallFileSize = info.Size()

	defer tmpBigFile.Close()
	defer tmpSmallFile.Close()

	bigFilePath = tmpBigFile.Name()
	smallFilePath = tmpSmallFile.Name()

	return nil
}
