package main

import (
	"path/filepath"
	"os"
	"fmt"
	"io/ioutil"
	"strings"
	//"golang.org/x/text/encoding/unicode"
	"regexp"
	"github.com/spf13/viper"
	"bufio"
	//"bytes"
)

func panicWhenErr(err error) {
	if err != nil {
		panic(err)
	}
}

func findVictimSql(input string) (result string) {
	r := strings.NewReplacer("\n", "", "\r", "")
	content := r.Replace(input)
	idRule := `deadlock victim="([a-z0-9]+)"`
	re := regexp.MustCompile(idRule)
	processId := re.FindStringSubmatch(content)
	sqlRule := `<process id="` + processId[1] + `".+<inputbuf>(.+)<\/inputbuf>`

	re = regexp.MustCompile(sqlRule)
	sql := re.FindStringSubmatch(content)
	if len(sql) >= 1 {
		result = string(sql[1])
	}
	return result
}

func config() {
	// todo: Default value
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func main() {
	var files []string
	config()
	rule := `objectname\=\"Survey.dbo\.(\w+)\"`
	directory := viper.GetString("deadlock.directory.path")
	reportPath := viper.GetString("deadlock.report.filepath")
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	panicWhenErr(err)
	f, err := os.Create(reportPath)
	panicWhenErr(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	fmt.Fprint(w, "Date,Datetime,Deadlock_table,XDL_filename,victim_SQL\n")
	for i, file := range files {
		if !strings.Contains(file, ".xdl") {
			continue
		}
		dat, err := ioutil.ReadFile(file)
		panicWhenErr(err)
		//dat, err = unicode.UTF8.NewDecoder().Bytes(dat)
		content := string(dat)
		re := regexp.MustCompile(rule)
		tableName := re.FindStringSubmatch(content)

		info, err := os.Stat(file)
		panicWhenErr(err)
		date := info.ModTime().Format("2006-01-02")
		dateTime := info.ModTime().Format("2006-01-02 15:04")
		sql := findVictimSql(content)
		if sql == "" {
			fmt.Print("**")
		}

		_, err = fmt.Fprintf(
			w, "%q,%q,%q,%q,%q\n", date, dateTime, tableName, info.Name(), strings.Replace(sql, "&apos;","'", -1))
		progress := float32(i + 1) / float32(len(files))
		fmt.Printf("\r Processing: %d/%d [%s%s]", i+1, len(files), strings.Repeat(">", int(progress*50)), strings.Repeat("*", 50-int(progress*50)))
	}

	panicWhenErr(err)
	w.Flush()
}
