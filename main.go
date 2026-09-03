package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"github.com/xuri/excelize/v2"
)

type Student struct {
	No   int
	Nama string
	NIM  string
}

type RenameItem struct {
	OldPath string
	NewPath string
	OldName string
	NewName string
}

func main() {
	fmt.Println("--- Tool Rename File HIMATIF ---")
	scanner := bufio.NewScanner(os.Stdin)

	excelFile, targetDir := resolvePaths()
	if excelFile == "" {
		fmt.Println("[ERROR] File Excel (.xlsx) tidak ditemukan di folder ini.")
		return
	}

	students, err := loadStudents(excelFile)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return
	}

	fmt.Printf("File Excel    : %s (%d mahasiswa)\n", filepath.Base(excelFile), len(students))
	fmt.Printf("Folder Target : %s\n", targetDir)

	files, err := listTargetFiles(targetDir, excelFile)
	if err != nil || len(files) == 0 {
		fmt.Printf("\n[INFO] Belum ada file sertifikat (PDF/JPG/PNG/JPEG) di folder \"%s\".\n", targetDir)
		fmt.Println("Silakan taruh file sertifikat di folder tersebut, lalu jalankan kembali.")
		return
	}

	withNIM := chooseFormat(scanner)
	items := matchStudentsWithFiles(students, files, targetDir, withNIM)
	if len(items) == 0 {
		fmt.Println("\n[PERINGATAN] Tidak ada file yang cocok dengan Nama, Nomor, atau NIM mahasiswa.")
		return
	}

	fmt.Printf("\nFile Cocok    : %d dari %d file\n\nPreview perubahan:\n", len(items), len(files))
	for i := 0; i < min(len(items), 5); i++ {
		fmt.Printf("  [%d] \"%s\" -> \"%s\"\n", i+1, items[i].OldName, items[i].NewName)
	}
	if len(items) > 5 {
		fmt.Printf("  ... dan %d file lainnya.\n", len(items)-5)
	}

	if !confirm(scanner, fmt.Sprintf("\nLanjutkan rename %d file? [Y/n]: ", len(items))) {
		fmt.Println("Proses dibatalkan.")
		return
	}

	executeRename(items)
}

func chooseFormat(scanner *bufio.Scanner) bool {
	for _, arg := range os.Args[1:] {
		lower := strings.ToLower(arg)
		if lower == "-nama-saja" || lower == "--nama-saja" || lower == "nama" {
			return false
		}
	}

	fmt.Println("\nPilihan format nama file:")
	fmt.Println("  [1] NIM dan Nama (contoh: 263307069_Diny Agata Rahmawati.ext)")
	fmt.Println("  [2] Nama saja    (contoh: Diny Agata Rahmawati.ext)")
	fmt.Print("Pilih format [1/2] (default 1): ")

	if scanner.Scan() {
		ans := strings.TrimSpace(strings.ToLower(scanner.Text()))
		return ans != "2" && ans != "nama" && ans != "nama saja" && ans != "nama aja"
	}
	return true
}

func resolvePaths() (string, string) {
	excelPath := ""
	targetDir := "."

	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			continue
		}
		if info, err := os.Stat(arg); err == nil {
			if info.IsDir() {
				targetDir = arg
			} else if strings.HasSuffix(strings.ToLower(arg), ".xlsx") {
				excelPath = arg
			}
		}
	}

	if excelPath == "" {
		entries, _ := os.ReadDir(".")
		for _, e := range entries {
			name := e.Name()
			if !e.IsDir() && strings.HasSuffix(strings.ToLower(name), ".xlsx") && !strings.HasPrefix(name, "~$") {
				excelPath = name
				break
			}
		}
	}

	return excelPath, targetDir
}

func loadStudents(path string) ([]Student, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka file Excel: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("file Excel tidak memiliki sheet")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil || len(rows) < 2 {
		return nil, fmt.Errorf("data di file Excel tidak mencukupi")
	}

	idxNama, idxNim, idxNo := -1, -1, -1
	for i, col := range rows[0] {
		c := strings.ToLower(strings.TrimSpace(col))
		if strings.Contains(c, "nama") && idxNama == -1 {
			idxNama = i
		} else if (strings.Contains(c, "nim") || strings.Contains(c, "npm")) && idxNim == -1 {
			idxNim = i
		} else if (c == "no" || c == "no." || c == "nomor") && idxNo == -1 {
			idxNo = i
		}
	}

	if idxNama == -1 || idxNim == -1 {
		return nil, fmt.Errorf("kolom Nama atau NIM tidak ditemukan di header: %v", rows[0])
	}

	var students []Student
	for i, r := range rows[1:] {
		if idxNama >= len(r) || idxNim >= len(r) {
			continue
		}
		nama := strings.TrimSpace(r[idxNama])
		nim := strings.TrimSpace(r[idxNim])
		if nama == "" || nim == "" {
			continue
		}

		no := i + 1
		if idxNo != -1 && idxNo < len(r) {
			if n, err := strconv.Atoi(strings.TrimSuffix(strings.TrimSpace(r[idxNo]), ".0")); err == nil {
				no = n
			}
		}

		students = append(students, Student{No: no, Nama: nama, NIM: nim})
	}

	return students, nil
}

func listTargetFiles(dir string, excelPath string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	excelBase := filepath.Base(excelPath)
	certExts := map[string]bool{
		".pdf": true, ".jpg": true, ".jpeg": true, ".png": true,
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if strings.HasPrefix(name, ".") || name == excelBase || !certExts[ext] {
			continue
		}
		files = append(files, name)
	}

	return files, nil
}

func matchStudentsWithFiles(students []Student, files []string, dir string, withNIM bool) []RenameItem {
	var items []RenameItem
	used := make(map[string]bool)

	for _, s := range students {
		normNama := cleanStr(s.Nama)
		noStr := strconv.Itoa(s.No)

		for _, f := range files {
			if used[f] {
				continue
			}

			base := strings.TrimSuffix(f, filepath.Ext(f))
			normBase := cleanStr(base)

			if normBase == normNama || strings.Contains(normBase, normNama) || normBase == noStr || strings.Contains(base, s.NIM) {
				used[f] = true
				ext := filepath.Ext(f)

				newName := fmt.Sprintf("%s_%s%s", s.NIM, s.Nama, ext)
				if !withNIM {
					newName = fmt.Sprintf("%s%s", s.Nama, ext)
				}

				items = append(items, RenameItem{
					OldPath: filepath.Join(dir, f),
					NewPath: filepath.Join(dir, newName),
					OldName: f,
					NewName: newName,
				})
				break
			}
		}
	}

	return items
}

func executeRename(items []RenameItem) {
	fmt.Println("\nMemproses rename...")
	sukses, lewat, gagal := 0, 0, 0

	for i, item := range items {
		prefix := fmt.Sprintf("  [%d/%d]", i+1, len(items))
		if item.OldPath == item.NewPath {
			fmt.Printf("%s [LEWAT] Nama sudah sesuai: %s\n", prefix, item.OldName)
			lewat++
			continue
		}
		if _, err := os.Stat(item.NewPath); err == nil {
			fmt.Printf("%s [LEWAT] File tujuan sudah ada: %s\n", prefix, item.NewName)
			lewat++
			continue
		}
		if err := os.Rename(item.OldPath, item.NewPath); err != nil {
			fmt.Printf("%s [GAGAL] %s -> %s: %v\n", prefix, item.OldName, item.NewName, err)
			gagal++
		} else {
			fmt.Printf("%s [OK] %s -> %s\n", prefix, item.OldName, item.NewName)
			sukses++
		}
	}

	fmt.Printf("\nSelesai.\nBerhasil : %d file\n", sukses)
	if lewat > 0 {
		fmt.Printf("Dilewati : %d file\n", lewat)
	}
	if gagal > 0 {
		fmt.Printf("Gagal    : %d file\n", gagal)
	}
}

func cleanStr(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")

	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			b.WriteRune(r)
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

func confirm(scanner *bufio.Scanner, prompt string) bool {
	fmt.Print(prompt)
	if scanner.Scan() {
		ans := strings.TrimSpace(strings.ToLower(scanner.Text()))
		return ans == "" || ans == "y" || ans == "ya" || ans == "yes"
	}
	return false
}
