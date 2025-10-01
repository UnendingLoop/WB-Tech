// Package model contains a structure for storing settings to cut
package model

type CutConfig struct {
	Fields  string   //-f — указание номеров столбцов для вывода; номера через запятую, можно диапазоны: «-f 1,3-5»
	Delim   string   //-d — использовать другой разделитель (по умолчанию '\t')
	SepOnly bool     //-s – только строки, содержащие разделитель - строки без разделителя игнорируются (не выводятся)
	Source  []string // Имя/имена файлов для чтения данных
}
