package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

type FileCompare struct {
	file1, file2 []byte
}

const dataStart int64 = 9624
const dataEnd int64 = 13601
const checksumLocation int64 = 13603
const fileName string = "Files/Pokemon - Versione Gialla - Speciale Edizione Pikachu (Italy) (GBC,SGB Enhanced).sav"

func NewFileCompare(nFile1 string, nFile2 string) FileCompare {
	file1, err := os.ReadFile(nFile1)
	check(err)
	file2, err := os.ReadFile(nFile2)
	check(err)
	return FileCompare{
		file1: file1,
		file2: file2,
	}
}

type EndLessThanStartError struct {
	start, end int64
}

func (e EndLessThanStartError) Error() string {
	return fmt.Sprintf("end (%d) is less than start (%d)", e.end, e.start)
}

type EndOutOfBoundError struct {
	bound, end int64
}

func (e EndOutOfBoundError) Error() string {
	return fmt.Sprintf("end (%d) is out of bound(%d)", e.end, e.bound)
}

type CompareResult struct {
	Index int64
	A, B  byte
}

func (f *FileCompare) Compare(start, end int64) (cr []CompareResult, err error) {
	s1, s2 := f.normalize()
	cr = make([]CompareResult, 0)

	if start >= end {
		err = EndLessThanStartError{
			start: start,
			end:   end,
		}
		return
	}
	if end > int64(len(s1)) {
		err = EndOutOfBoundError{
			bound: int64(len(s1)),
			end:   end,
		}
		return
	}

	for i := start; i < end; i++ {
		if s1[i] != s2[i] {
			cr = append(cr, CompareResult{
				Index: i,
				A:     s1[i],
				B:     s2[i],
			})
		}

	}
	return
}

func (f *FileCompare) normalize() (s1, s2 []byte) {
	l := 0
	if len(f.file1) >= len(f.file2) {
		l = len(f.file1)
	} else {
		l = len(f.file2)
	}
	s1 = f.file1[:l]
	s2 = f.file2[:l]
	return
}

func main() {

	//nFile contiene la path relativa a ./Files/"tuofile" inserita dall'utente
	nFile1, nFile2 := "Files/", "Files/"
	nFile1 += insFile(nFile1, 1)
	nFile2 += insFile(nFile2, 2)

	files := NewFileCompare(nFile1, nFile2)
	var str string
	for {
		fmt.Print("insrisci l'inizio del file ")
		_, err := fmt.Scan(&str)
		check(err)
		start, _ := strconv.ParseInt(str, 16, 64)
		fmt.Print("insrisci la fine del file ")
		_, err = fmt.Scan(&str)
		check(err)
		end, _ := strconv.ParseInt(str, 16, 64)
		compare, err := files.Compare(start, end)
		if err != nil {
			fmt.Println(err)
			continue
		}
		scelta(compare, files)

		break
	}

}

func check(e error) {
	if e != nil && e != io.EOF {
		panic(e)
	}
}

func insFile(nFile string, n int) string {

	for {
		fmt.Printf("Inserire il nome del file %v  (deve essere contenuto nella cartella Files)\t", n)
		_, err := fmt.Scan(&nFile)
		check(err)
		_, err = os.Stat("Files/" + nFile)
		if os.IsNotExist(err) {
			fmt.Print("il file non esiste, reinseriscilo\n")
		} else {
			break
		}
	}

	return nFile
}

func PrintAllCompare(result []CompareResult) {
	if len(result) == 0 {
		fmt.Printf("File are the same\n")
		return
	}

	for i, compareResult := range result {

		fmt.Printf("%d) Difference at byte %X: %X %X\n", i+1, compareResult.Index, compareResult.A, compareResult.B)
	}
}

func scelta(compare []CompareResult, files FileCompare) {
	var s int
	exit := 0
	for exit != 1 {
		fmt.Print("\n\n\nInserisci l'opzione desiderata per:\n")
		fmt.Print("1: stampa di tutti i valori diversi\n")
		fmt.Print("2: stampa di tutti i valori diversi consecutivi\n")
		fmt.Print("3: calcolo e controllo checksum\n")
		fmt.Print("4: cambiare dei byte\n")
		fmt.Print("5: per uscire dal programma\t")
		for {
			_, err := fmt.Scan(&s)
			check(err)
			if s <= 5 && s >= 1 {
				break
			}
			fmt.Print("reinseriscila scimmia\n")
		}

		switch s {
		case 1:
			PrintAllCompare(compare)
		case 2:
			PrintNeighboors(compare)
		case 3:
			checksum(files)
		case 4:
			WriteSwtich(files, compare)
		case 5:
			exit = 1
		}
	}

}

func PrintNeighboors(result []CompareResult) {
	if len(result) == 0 {
		fmt.Printf("File are the same\n")
		return
	}
	for i, compareResult := range result {
		ricPrint(i, compareResult, result)
	}

}

// todo sistemare errori del print ricorsivo
func ricPrint(i int, compare CompareResult, result []CompareResult) {
	defer func() {
		if r := recover(); r != nil {

		}
	}()
	if compare.Index == result[i+1].Index {
		compare = result[i+1]
		ricPrint(i+1, compare, result)
		fmt.Printf("%d) Difference at byte %X: %X %X\n", i+1, compare.Index, compare.A, compare.B)
	}
}

/*
	func kMain(in []CompareResult) (out []CompareResult) {
		out = make([]CompareResult, 0)
		for i := 0; i < len(in); {
			outPart := keepNeighboursBetter(in[i+1:], in[i].Index)
			if len(outPart) != 0 {
				out = append(out, in[i])
				out = append(out, outPart...)
				i += len(out)
			}
		}
		return
	}

	func keepNeighboursBetter(in []CompareResult, indexToBeEqual int64) (out []CompareResult) {
		out = make([]CompareResult, 0)
		if len(in) < 2 {
			return
		}
		if in[0].Index-1 == indexToBeEqual {
			out = append(out, in[0])
			out = append(out, keepNeighboursBetter(in[1:], in[0].Index)...)
		}
		return
	}
*/
func checksum(files FileCompare) byte {
	var cs byte
	cs = 255
	for i := dataStart; i <= dataEnd; i++ {
		cs = cs - files.file1[i]
	}

	fmt.Printf("%X \n", files.file1[checksumLocation])
	fmt.Printf("%X", cs)
	return cs
}

func WriteSwtich(files FileCompare, compare []CompareResult) {
	var s int
	exit := 0
	for exit != 1 {
		fmt.Print("\n\nScegliere:\n")
		fmt.Print("1: Cambiare un nuovo byte\n")
		fmt.Print("2: Cambiare un byte con indirizzo personalizzato")
		fmt.Print("3: Non cambiare\n")
		for {
			_, err := fmt.Scan(&s)
			check(err)
			if s <= 3 && s >= 1 {
				break
			}
			fmt.Print("reinseriscila scimmia\n")
		}
		switch s {
		case 1:
			Write(files, compare)
		case 2:
			WritePers(files)
		case 3:
			exit = 1

		}

	}
}

func WritePers(files FileCompare) {
	var str string
	fmt.Print("inserisci l'indirizzo da sostituire")
	_, err := fmt.Scan(&str)
	check(err)
	ind, _ := strconv.ParseInt(str, 16, 64)
	var by byte
	fmt.Printf("\n Inserire il numero da sostituire ")
	for {
		_, err := fmt.Scanf("%X", &by)
		if err == nil {
			break
		}
	}

	files.file1[ind] = by
	files.file1[checksumLocation] = checksum(files)

	err = os.WriteFile(fileName, files.file1, 0666)
	check(err)
}

func Write(files FileCompare, compare []CompareResult) {

	var ind int
	fmt.Print("scegliere quale byte modificare inserendo il numero associato all'indirizzo desiderato" + "\n")
	for i := 0; i < len(compare); i++ {
		fmt.Printf("%v : %X\n", i, compare[i].Index)
	}
	for {
		_, err := fmt.Scan(&ind)
		check(err)
		if (ind < len(compare)) && ind >= 0 {
			break
		}
		fmt.Print("numero inesistente, reinserirlo\n")
	}
	var by byte
	fmt.Printf("\n Inserire il numero da sostituire ")
	for {
		_, err := fmt.Scanf("%X", &by)
		if err == nil {
			break
		}
	}

	files.file1[compare[ind].Index] = by
	files.file1[checksumLocation] = checksum(files)

	err := os.WriteFile(fileName, files.file1, 0666)
	check(err)
}

//todo miglioni di controlli errori
