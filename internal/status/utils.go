package status

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Checkable interface {
	CheckValue() float64
	String() string
	PerfData(string, string) string
}

type CheckableAgainst interface {
	CheckValue() string
	String() string
}

type CheckableExtra interface {
	LongOutput() string
}

/*
type apiError struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}
*/

type CheckResult struct {
	Exitcode   int    `json:"exitcode"`
	Output     string `json:"output"`
	Perfdata   string `json:"perfdata"`
	LongOutput string `json:"longoutput"`
}

func (c *CheckResult) String() string {
	return fmt.Sprintf("CheckResult: output: %s | exitcode: %d", c.Output, c.Exitcode)
}

func GetCheckResult(chk Checkable, w, c string) CheckResult {
	exitcode := 0
	output := fmt.Sprintf("OK - %s", chk.String())
	cv := chk.CheckValue()
	if isInRange(w, cv) {
		exitcode = 1
		output = fmt.Sprintf("WARNING - %s", chk.String())
	}
	if isInRange(c, cv) {
		exitcode = 2
		output = fmt.Sprintf("CRITICAL - %s", chk.String())
	}

	perfdata := chk.PerfData(w, c)

	longoutput := ""
	ex, ok := chk.(CheckableExtra)
	if ok {
		longoutput = ex.LongOutput()
	}

	return CheckResult{
		exitcode, output, perfdata, longoutput,
	}
}

func GetCheckAgainstResult(chk CheckableAgainst, e string) CheckResult {
	exitcode := 0
	output := fmt.Sprintf("OK: %s", chk.String())
	cv := chk.CheckValue()
	if cv != e {
		exitcode = 2
		output = fmt.Sprintf("CRITICAL - %s", chk.String())
	}

	return CheckResult{
		exitcode, output, "", "",
	}
}

func createPerfData(pre, w, c string) string {
	var perf []string
	perf = append(perf, pre)
	if w != "" {
		perf = append(perf, w)
	} else if c != "" {
		// Still add the ; if we are going to be adding a critical value
		perf = append(perf, "")
	}
	if c != "" {
		perf = append(perf, c)
	}
	return strings.Join(perf, ";")
}

func isInRange(r string, value float64) bool {

	if r == "" {
		return false
	}

	// Get values that match regex (1 or 2 are valid)
	regexExp := `(-?[0-9]+(\.[0-9]+)?)`
	re := regexp.MustCompile(regexExp)
	matches := re.FindAllString(r, -1)
	if len(matches) == 0 {
		return false
	}
	p1, _ := strconv.ParseFloat(string(matches[0]), 64)

	if len(matches) == 2 {
		p2, _ := strconv.ParseFloat(string(matches[1]), 64)
		// Check for @10:20
		re = regexp.MustCompile(fmt.Sprintf(`^@%s:%s$`, matches[0], matches[1]))
		if re.MatchString(r) && (value < p1 || value > p2) {
			return true
		}
		// Check for 10:20
		re = regexp.MustCompile(fmt.Sprintf(`^%s:%s$`, matches[0], matches[1]))
		if re.MatchString(r) && value >= p1 && value <= p2 {
			return true
		}
	} else if len(matches) == 1 {
		// Check for 10
		re = regexp.MustCompile(fmt.Sprintf(`^%s$`, matches[0]))
		if re.MatchString(r) && (value > p1 || value < 0) {
			return true
		}
		// Check for ~:10
		re = regexp.MustCompile(fmt.Sprintf(`^~:%s$`, matches[0]))
		if re.MatchString(r) && value > p1 {
			return true
		}
		// Check for :10
		re = regexp.MustCompile(fmt.Sprintf(`^:%s$`, matches[0]))
		if re.MatchString(r) && (value > p1 || value < 0) {
			return true
		}
		// Check for 10:
		re = regexp.MustCompile(fmt.Sprintf(`^%s:$`, matches[0]))
		if re.MatchString(r) && value < p1 {
			return true
		}
	}

	// Assume false if bad value, we can validate
	// them properly in the future if we need to
	return false
}

/*
func errorCheckResult(err error) CheckResult {
	return CheckResult{
		Exitcode: 3,
		Output:   fmt.Sprintf("UNKOWN - %s", err),
	}
}
*/
