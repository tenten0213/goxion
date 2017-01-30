package main

import (
	"github.com/sclevine/agouti"
	"github.com/urfave/cli"
	"github.com/mitchellh/go-ps"
	"log"
	"io"
	"encoding/csv"
	"os"
	"strconv"
	"regexp"
	"strings"
	"time"
	"runtime"
)

func main() {

	app := cli.NewApp()
	app.Name = "goxion"
	app.Version = "0.1.0"
	app.Author = "tenten0213"
	app.Email = ""
	app.Usage = ""
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "username, u", Usage: "userId for login. ENV is GOXION_USER."},
		cli.StringFlag{Name: "password, p", Usage: "password for login. ENV is GOXION_PASSWORD."},
		cli.StringFlag{Name: "pjid, i", Usage: "projectid for login. ENV is GOXION_PJID."},
		cli.StringFlag{Name: "url, r", Usage: "login page's url. ENV is GOXION_URL."},
		cli.StringFlag{Name: "matrixpath, m", Usage: "path to matrix file. ENV is GOXION_MATRIX."},
	}
	app.Action = func(c *cli.Context) error {
		username := os.Getenv("GOXION_USER"); if username == "" {
			username = c.String("username")
		}

		password := os.Getenv("GOXION_PASSWORD"); if password == "" {
			password = c.String("password")
		}

		pjid := os.Getenv("GOXION_PJID"); if pjid == "" {
			pjid = c.String("pjid")
		}

		url := os.Getenv("GOXION_URL"); if url == "" {
			url = c.String("url")
		}

		matrixPath := os.Getenv("GOXION_MATRIX"); if matrixPath == "" {
			matrixPath = c.String("matrixpath")
		}

		if username == "" || password == "" ||
			pjid == "" || url == "" || matrixPath == "" {
			log.Printf("username = %v", username)
			log.Printf("password = %v", password)
			log.Printf("pjid = %v", pjid)
			log.Printf("url = %v", url)
			log.Printf("matrixpath = %v", matrixPath)
			log.Fatalln("Failed to load requre parameters")
			return nil
		}

		processes, err := ps.Processes()

		for i := range processes {
			if processes[i].Executable() == "safaridriver" {
				process, err := os.FindProcess(processes[i].Pid()); if err != nil {
					log.Fatalf("Failed to kill driver process:%v", err)
					return nil
				}
				process.Kill()
				break
			}
		}

		var driver *agouti.WebDriver
		if runtime.GOOS == "darwin" {
			driver = agouti.Selenium(agouti.Browser("safari"))
		} else if runtime.GOOS == "windows" {
			driver = agouti.Selenium(agouti.Browser("internet explorer"))
		} else {
			log.Fatalln("Don't support your os")
			return nil
		}

		if err := driver.Start(); err != nil {
			log.Fatalf("Failed to start driver:%v", err)
			defer driver.Stop()
			return nil
		}

		page, err := driver.NewPage()
		if err != nil {
			log.Fatalf("Failed to open page:%v", err)
			defer driver.Stop()
			return nil
		}

		if err := page.Navigate(url); err != nil {
			log.Fatalf("Failed to navigate:%v", err)
			defer driver.Stop()
			return nil
		}

		usernameField := page.FindByID("input_1")
		usernameField.Fill(username)

		passwordField := page.FindByID("input_2")
		passwordField.Fill(password)

		pjidField := page.FindByID("input_3")
		pjidField.Fill(pjid)
		if err := page.FindByClass("credentials_input_submit").Submit(); err != nil {
			log.Fatalf("Failed to login:%v", err)
			defer driver.Stop()
			return nil
		}
		time.Sleep(1 * time.Second)
		log.Println("User Login Success.")

		credentials, err := page.FindByID("credentials_table_header").Text()
		if err != nil {
			log.Fatalf("Failed to get credentials:%v", err)
			defer driver.Stop()
			return nil
		}
		r := regexp.MustCompile("\\[(.+?)\\]")
		keys := r.FindAllString(credentials, -1)
		fp, err := os.Open(matrixPath)
		if err != nil {
			log.Fatalf("Failed to open file%v", err)
			defer driver.Stop()
			return nil
		}
		reader := csv.NewReader(fp)

		reader.Comma = ','
		reader.LazyQuotes = true

		// create map for matrix
		matrix := map[string]string{}
		var elem = []string{"A","B","C","D","E","F","G","H","I","J"}
		var rowNum = 0
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
			for colNum := range record {
				matrix[elem[colNum] + strconv.Itoa(rowNum + 1)] = record[colNum]
			}
			rowNum++
		}
		var challenge = ""
		for i := range keys {
			challenge += matrix[strings.Replace(strings.Replace(keys[i],"[", "", -1), "]", "", -1)]
		}
		log.Printf("challenge = %v", challenge)

		f5challenge := page.FindByID("input_1")
		f5challenge.Fill(challenge)
		if err := page.FindByClass("credentials_input_password").Submit(); err != nil {
			log.Fatalf("Failed to Credential Challenge:%v", err)
		}
		return nil
	}
	app.Run(os.Args)
}
