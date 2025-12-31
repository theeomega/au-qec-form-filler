package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"golang.org/x/term"
)

// --- CONFIGURATION ---
const (
	LoginURL                  = "https://portals.au.edu.pk/qec/login.aspx"
	BaseURL                   = "https://portals.au.edu.pk/qec/"
	FirstPerformaURL          = BaseURL + "p1.aspx"
	TeacherEvaluationURL      = BaseURL + "p10.aspx"
	OnlineLearningFeedbackURL = BaseURL + "p10a_learning_online_form.aspx"
	UserAgent                 = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

// --- STYLING HELPERS ---

func RGB(hex string) func(string) string {
	var r, g, b uint8
	fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	return func(s string) string {
		return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", r, g, b, s)
	}
}

var (
	ColorBlue    = RGB("#0c7ec9")
	ColorPurple  = RGB("#bf06bf")
	ColorGold    = RGB("#f5c816")
	ColorCyan    = RGB("#00ffff")
	ColorGreen   = RGB("#00ff00")
	ColorRed     = RGB("#ff0000")
	ColorGrey    = RGB("#8a8a8a")
	ColorMagenta = RGB("#d700d7")
)

// --- HTTP CLIENT SETUP ---
type HiddenFields struct {
	ViewState          string
	EventValidation    string
	ViewStateGenerator string
}

var (
	client *http.Client
	reader *bufio.Reader
)

func init() {
	jar, _ := cookiejar.New(nil)
	client = &http.Client{Jar: jar}
	reader = bufio.NewReader(os.Stdin)
}

func makeRequest(method, urlStr string, data url.Values) (*http.Response, error) {
	var req *http.Request
	var err error
	if method == "POST" && data != nil {
		req, err = http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req, err = http.NewRequest("GET", urlStr, nil)
		if err != nil {
			return nil, err
		}
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Referer", LoginURL)
	return client.Do(req)
}

func getHiddenFields(doc *goquery.Document) HiddenFields {
	vs, _ := doc.Find("input[name='__VIEWSTATE']").Attr("value")
	ev, _ := doc.Find("input[name='__EVENTVALIDATION']").Attr("value")
	vsg, _ := doc.Find("input[name='__VIEWSTATEGENERATOR']").Attr("value")
	return HiddenFields{vs, ev, vsg}
}

// --- UI FUNCTIONS ---

func printBanner() {
	fmt.Print("\033[H\033[2J")
	topArt := `
 ██████╗ ███████╗ ██████╗      ██████╗ ██████╗ ███╗   ███╗ █████╗ ████████╗██╗ ██████╗ ███╗   ██╗
██╔═══██╗██╔════╝██╔════╝     ██╔═══██╗╚════██╗████╗ ████║██╔══██╗╚══██╔══╝██║██╔═══██╗████╗  ██║
██║   ██║█████╗  ██║          ██║   ██║ █████╔╝██╔████╔██║███████║   ██║   ██║██║   ██║██╔██╗ ██║`
	
	bottomArt := `██║▄▄ ██║██╔══╝  ██║          ██║   ██║██╔═══╝ ██║╚██╔╝██║██╔══██║   ██║   ██║██║   ██║██║╚██╗██║
╚██████╔╝███████╗╚██████╗     ╚██████╔╝███████╗██║ ╚═╝ ██║██║  ██║   ██║   ██║╚██████╔╝██║ ╚████║
 ╚══▀▀═╝ ╚══════╝ ╚═════╝      ╚═════╝ ╚══════╝╚═╝     ╚═╝╚═╝  ╚═╝   ╚═╝   ╚═╝ ╚═════╝ ╚═╝  ╚═══╝`

	fmt.Println(ColorBlue(topArt))
	fmt.Println(ColorPurple(bottomArt))
	fmt.Println()
	fmt.Println(ColorCyan("Developer GitHub: https://github.com/theeomega"))
	fmt.Println("\n\n")
}

func prompt(label string, isPassword bool) string {
	fmt.Print(ColorGold(label))
	if isPassword {
		// Raw mode for asterisk masking
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			bytePass, _ := term.ReadPassword(int(os.Stdin.Fd()))
			return string(bytePass)
		}
		defer term.Restore(int(os.Stdin.Fd()), oldState)

		var password []byte
		for {
			b := make([]byte, 1)
			_, err := os.Stdin.Read(b)
			if err != nil { break }
			char := b[0]

			if char == 3 { // Ctrl+C
				term.Restore(int(os.Stdin.Fd()), oldState)
				os.Exit(1)
			} else if char == 13 { // Enter
				fmt.Print("\r\n")
				break
			} else if char == 127 || char == 8 { // Backspace
				if len(password) > 0 {
					password = password[:len(password)-1]
					fmt.Print("\b \b")
				}
			} else {
				password = append(password, char)
				fmt.Print("*")
			}
		}
		return string(password)
	}

	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func printPanel(title, subtitle string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)
	t.Style().Options.SeparateRows = false
	t.Style().Box.PaddingLeft = "  "
	t.Style().Box.PaddingRight = "  "
	
	header := text.FormatUpper.Apply(title)
	t.AppendRow(table.Row{text.AlignCenter.Apply(header, 60)})
	if subtitle != "" {
		t.AppendRow(table.Row{text.AlignCenter.Apply(subtitle, 60)})
	}
	t.Render()
}

// --- MAIN LOGIC ---

func main() {
	printBanner()

	loginSuccess := false
	var username string
	
	// FIX: Declare hf here so it is available everywhere in main()
	var hf HiddenFields 

	// --- LOGIN ---
	for !loginSuccess {
		username = prompt("Enter your ID: ", false)
		password := prompt("Enter your password: ", true)

		success := func() bool {
			resp, err := makeRequest("GET", LoginURL, nil)
			if err != nil { return false }
			doc, _ := goquery.NewDocumentFromReader(resp.Body)
			resp.Body.Close()

			hf = getHiddenFields(doc) // Updates the outer hf variable
			if hf.ViewState == "" { return false }

			data := url.Values{}
			data.Set("ctl00$ContentPlaceHolder2$ddlcampus", "Islamabad")
			data.Set("ctl00$ContentPlaceHolder2$ddlUserType", "Student/Alumni")
			data.Set("ctl00$ContentPlaceHolder2$txt_regid", username)
			data.Set("ctl00$ContentPlaceHolder2$txt_password", password)
			data.Set("ctl00$ContentPlaceHolder2$btnAccountlogin", "Login")
			data.Set("__VIEWSTATE", hf.ViewState)
			data.Set("__EVENTVALIDATION", hf.EventValidation)
			data.Set("__VIEWSTATEGENERATOR", hf.ViewStateGenerator)

			loginResp, err := makeRequest("POST", LoginURL, data)
			if err != nil { return false }
			defer loginResp.Body.Close()

			bodyBytes, _ := io.ReadAll(loginResp.Body)
			if strings.Contains(strings.ToLower(string(bodyBytes)), "logout") {
				return true
			}
			return false
		}()

		if success {
			fmt.Println(ColorGreen("\n[ - ] Login successful\n"))
			loginSuccess = true
		} else {
			fmt.Println(ColorRed("\n[ X ] Login failed. Please try again.\n"))
		}
	}

	// --- TEACHER EVALUATION SETUP ---
	resp, _ := makeRequest("GET", TeacherEvaluationURL, nil)
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	hf = getHiddenFields(doc) // Now valid because hf is declared at top of main

	// Enhanced Struct to hold Course and Grade state
	type Teacher struct {
		Value  string
		Name   string
		Course string
		Grade  string
	}
	var teacherList []*Teacher

	doc.Find("select[name='ctl00$ContentPlaceHolder2$ddlTeacher'] option").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			val, _ := s.Attr("value")
			teacherList = append(teacherList, &Teacher{
				Value:  val, 
				Name:   strings.TrimSpace(s.Text()),
				Course: "...", // Placeholder
				Grade:  "-",
			})
		}
	})

	choice := prompt("\nDo you want to give custom grades? (default 'n') [y/n]: ", false)
	
	// Map to store for later proformas
	customGrades := make(map[string]string)
	teacherCourseMap := make(map[string]string)

	if strings.ToLower(choice) == "y" {
		
		for i, teacher := range teacherList {
			// 1. Fetch Course Info for current teacher
			data := url.Values{}
			data.Set("__EVENTTARGET", "ctl00$ContentPlaceHolder2$ddlTeacher")
			data.Set("__VIEWSTATE", hf.ViewState)
			data.Set("__VIEWSTATEGENERATOR", hf.ViewStateGenerator)
			data.Set("__EVENTVALIDATION", hf.EventValidation)
			data.Set("ctl00$ContentPlaceHolder2$ddlTeacher", teacher.Value)

			r, _ := makeRequest("POST", TeacherEvaluationURL, data)
			tDoc, _ := goquery.NewDocumentFromReader(r.Body)
			r.Body.Close()
			hf = getHiddenFields(tDoc)

			// Scrape Course
			courseName := "Unknown"
			tDoc.Find("[id*='lblCourse']").Each(func(i int, s *goquery.Selection) {
				courseName = strings.TrimSpace(s.Text())
			})
			
			// Update the struct
			teacher.Course = courseName
			
			// 2. Render the Table (Refreshes every iteration)
			fmt.Print("\033[H\033[2J") // Clear Screen
			printPanel("Custom Grade Assignment", "Assign per-teacher grades")

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.SetStyle(table.StyleRounded)
			t.AppendHeader(table.Row{"No.", "Teacher", "Course", "Grade"})
			t.Style().Color.Header = text.Colors{text.FgMagenta, text.Bold}
			
			for idx, tObj := range teacherList {
				// Highlight current row
				teacherName := tObj.Name
				if idx == i {
					teacherName = ColorCyan(tObj.Name) + " ◄"
				} else if idx < i {
					teacherName = ColorGrey(tObj.Name)
				}
				
				gradeDisplay := tObj.Grade
				if idx < i {
					gradeDisplay = ColorGreen(gradeDisplay)
				} else if idx == i {
					gradeDisplay = ColorGold("Pending")
				}

				t.AppendRow(table.Row{idx + 1, teacherName, tObj.Course, gradeDisplay})
			}
			t.Render()

			// 3. Prompt
			fmt.Printf("\n%s %s %s\n\n", 
				ColorGold("Grade for"), 
				ColorCyan(teacher.Name), 
				ColorGold("("+courseName+")"))
			
			fmt.Println("A = Strongly Agree | B = Agree | C = Disagree | D = Strongly Disagree\n")
			
			valid := false
			var grade string
			for !valid {
				grade = strings.ToUpper(prompt("Enter the Grade: ", false))
				if grade == "A" || grade == "B" || grade == "C" || grade == "D" {
					valid = true
				}
			}
			
			// Save State
			teacher.Grade = grade
			customGrades[teacher.Value] = grade
			teacherCourseMap[teacher.Value] = courseName
		}

		// --- FINAL RENDER (Updates last row from Pending -> Grade) ---
		fmt.Print("\033[H\033[2J")
		printPanel("Custom Grade Assignment", "Assign per-teacher grades")
		tFinal := table.NewWriter()
		tFinal.SetOutputMirror(os.Stdout)
		tFinal.SetStyle(table.StyleRounded)
		tFinal.AppendHeader(table.Row{"No.", "Teacher", "Course", "Grade"})
		tFinal.Style().Color.Header = text.Colors{text.FgMagenta, text.Bold}
		for idx, tObj := range teacherList {
			tFinal.AppendRow(table.Row{idx + 1, tObj.Name, tObj.Course, ColorGreen(tObj.Grade)})
		}
		tFinal.Render()

		// Success message
		fmt.Println(ColorGreen("\n[ - ] Custom grades recorded."))
	}

	// --- 1. FILLING PROFORMA 1 (Subjects) ---
	resp, _ = makeRequest("GET", FirstPerformaURL, nil)
	doc, _ = goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	hf = getHiddenFields(doc) // Now works fine

	// Reuse struct logic for parsing subject options
	var subjects []*Teacher
	doc.Find("select[name='ctl00$ContentPlaceHolder2$cmb_courses'] option").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			val, _ := s.Attr("value")
			subjects = append(subjects, &Teacher{Value: val, Name: strings.TrimSpace(s.Text())})
		}
	})

	subjectGradeMap := make(map[string]string)
	for tVal, grade := range customGrades {
		if cName, ok := teacherCourseMap[tVal]; ok {
			for _, sub := range subjects {
				if strings.Contains(strings.ToLower(sub.Name), strings.ToLower(cName)) || 
				   strings.Contains(strings.ToLower(cName), strings.ToLower(sub.Name)) {
					subjectGradeMap[sub.Value] = grade
				}
			}
		}
	}

	if len(subjects) == 0 {
		fmt.Println(ColorGold("[ - ] Subjects evaluation proforma already filled"))
	} else {
		for _, sub := range subjects {
			grade := "A"
			if g, ok := subjectGradeMap[sub.Value]; ok { grade = g }

			data := url.Values{}
			data.Set("__EVENTTARGET", "ctl00$ContentPlaceHolder2$cmb_courses")
			data.Set("__VIEWSTATE", hf.ViewState)
			data.Set("__VIEWSTATEGENERATOR", hf.ViewStateGenerator)
			data.Set("__EVENTVALIDATION", hf.EventValidation)
			data.Set("ctl00$ContentPlaceHolder2$cmb_courses", sub.Value)
			for i := 1; i <= 12; i++ {
				data.Set(fmt.Sprintf("ctl00$ContentPlaceHolder2$q%d", i), grade)
			}
			data.Set("ctl00$ContentPlaceHolder2$btnSave", "Submit Proforma")

			r, _ := makeRequest("POST", FirstPerformaURL, data)
			sDoc, _ := goquery.NewDocumentFromReader(r.Body)
			r.Body.Close()
			hf = getHiddenFields(sDoc)
			
			fmt.Printf("      [-] Submitted for %s\n", ColorGrey(sub.Name))
		}
		fmt.Println(ColorGreen("[ - ] Subjects evaluation proforma filled successfully."))
	}

	// --- 2. FILLING PROFORMA 10 (Teachers) ---
	resp, _ = makeRequest("GET", TeacherEvaluationURL, nil)
	doc, _ = goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	hf = getHiddenFields(doc)
	
	var submissionTeachers []*Teacher
	doc.Find("select[name='ctl00$ContentPlaceHolder2$ddlTeacher'] option").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			val, _ := s.Attr("value")
			submissionTeachers = append(submissionTeachers, &Teacher{Value: val, Name: strings.TrimSpace(s.Text())})
		}
	})

	if len(submissionTeachers) == 0 {
		fmt.Println(ColorGold("[ - ] Teacher evaluation forms already filled"))
	} else {
		for _, t := range submissionTeachers {
			data := url.Values{}
			data.Set("__EVENTTARGET", "ctl00$ContentPlaceHolder2$ddlTeacher")
			data.Set("__VIEWSTATE", hf.ViewState)
			data.Set("__VIEWSTATEGENERATOR", hf.ViewStateGenerator)
			data.Set("__EVENTVALIDATION", hf.EventValidation)
			data.Set("ctl00$ContentPlaceHolder2$ddlTeacher", t.Value)
			
			r, _ := makeRequest("POST", TeacherEvaluationURL, data)
			tDoc, _ := goquery.NewDocumentFromReader(r.Body)
			r.Body.Close()
			hf = getHiddenFields(tDoc)

			grade := "A"
			if g, ok := customGrades[t.Value]; ok { grade = g }

			formData := url.Values{}
			formData.Set("__VIEWSTATE", hf.ViewState)
			formData.Set("__VIEWSTATEGENERATOR", hf.ViewStateGenerator)
			formData.Set("__EVENTVALIDATION", hf.EventValidation)
			formData.Set("ctl00$ContentPlaceHolder2$ddlTeacher", t.Value)
			for i := 1; i <= 16; i++ {
				formData.Set(fmt.Sprintf("ctl00$ContentPlaceHolder2$q%d", i), grade)
			}
			formData.Set("ctl00$ContentPlaceHolder2$q20", "Good instructor")
			formData.Set("ctl00$ContentPlaceHolder2$q21", "Good course")
			formData.Set("ctl00$ContentPlaceHolder2$btnSave", "Save Proforma Proforma")

			rFinal, _ := makeRequest("POST", TeacherEvaluationURL, formData)
			rFinal.Body.Close()
			fmt.Printf("      [-] Submitted for %s (Grade: %s)\n", t.Name, grade)
		}
		fmt.Println(ColorGreen("[ - ] Teacher evaluations proforma filled successfully."))
	}

	// --- 3. FILLING PROFORMA 10a (Online Learning) ---
	resp, _ = makeRequest("GET", OnlineLearningFeedbackURL, nil)
	doc, _ = goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	hf = getHiddenFields(doc)

	var onlineSubjects []*Teacher
	doc.Find("select[name='ctl00$ContentPlaceHolder1$cmb_courses'] option").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			val, _ := s.Attr("value")
			onlineSubjects = append(onlineSubjects, &Teacher{Value: val, Name: strings.TrimSpace(s.Text())})
		}
	})

	if len(onlineSubjects) == 0 {
		fmt.Println(ColorGold("[ - ] Online Learning Feedback Proformas already filled"))
	} else {
		for _, sub := range onlineSubjects {
			data := url.Values{}
			data.Set("__EVENTTARGET", "ctl00$ContentPlaceHolder1$cmb_courses")
			data.Set("__VIEWSTATE", hf.ViewState)
			data.Set("__VIEWSTATEGENERATOR", hf.ViewStateGenerator)
			data.Set("__EVENTVALIDATION", hf.EventValidation)
			data.Set("ctl00$ContentPlaceHolder1$cmb_courses", sub.Value)

			r, _ := makeRequest("POST", OnlineLearningFeedbackURL, data)
			oDoc, _ := goquery.NewDocumentFromReader(r.Body)
			r.Body.Close()
			hf = getHiddenFields(oDoc)

			grade := "A"
			if g, ok := subjectGradeMap[sub.Value]; ok { grade = g }

			formData := url.Values{}
			formData.Set("__VIEWSTATE", hf.ViewState)
			formData.Set("__VIEWSTATEGENERATOR", hf.ViewStateGenerator)
			formData.Set("__EVENTVALIDATION", hf.EventValidation)
			formData.Set("ctl00$ContentPlaceHolder1$cmb_courses", sub.Value)
			for i := 1; i <= 16; i++ {
				formData.Set(fmt.Sprintf("ctl00$ContentPlaceHolder1$q%d", i), grade)
			}
			formData.Set("ctl00$ContentPlaceHolder1$q20", "Good online learning experience")
			formData.Set("ctl00$ContentPlaceHolder1$btnSave", "Submit Proforma")

			rFinal, _ := makeRequest("POST", OnlineLearningFeedbackURL, formData)
			rFinal.Body.Close()
			fmt.Printf("      [-] Submitted for %s\n", ColorGrey(sub.Name))
		}
		fmt.Println(ColorGreen("[-] Online learning feedback proformas filled successfully.\n\n"))
	}

	fmt.Println("\nPress Enter to close the program...")
	reader.ReadString('\n')
}
