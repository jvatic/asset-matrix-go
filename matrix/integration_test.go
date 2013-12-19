package matrix

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	. "launchpad.net/gocheck"
)

// Hook gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type IntegrationSuite struct{}

var _ = Suite(&IntegrationSuite{})

type beginsChecker struct {
	*CheckerInfo
}

var BeginsWith Checker = &beginsChecker{
	&CheckerInfo{Name: "BeginsWith", Params: []string{"obtained", "expected prefix"}},
}

func (checker *beginsChecker) Check(params []interface{}, names []string) (result bool, err string) {
	defer func() {
		if v := recover(); v != nil {
			result = false
			err = fmt.Sprint(v)
		}
	}()

	if a, ok := params[0].(string); ok {
		if b, ok := params[1].(string); ok {
			return strings.HasPrefix(a, b), ""
		} else {
			return false, fmt.Sprintf("%v (%T) is not a string", params[1], params[1])
		}
	} else {
		return false, fmt.Sprintf("%v (%T) is not a string", params[0], params[0])
	}
}

func (s *IntegrationSuite) TestRequire(c *C) {
	manifest := NewManifest([]string{"./support/test_1/input"}, "./support/test_1/output", os.Stdout)
	if err := manifest.ScanInputDirs(); err != nil {
		c.Error(err)
	}

	if err := manifest.EvaluateDirectives(); err != nil {
		c.Error(err)
	}

	if err := manifest.ConfigureHandlers(); err != nil {
		c.Error(err)
	}

	if err := manifest.WriteOutput(); err != nil {
		c.Error(err)
	}

	file, err := os.Open("./support/test_1/output/file_1.js")
	if err != nil {
		c.Error(err)
	}
	defer file.Close()
	data := new(bytes.Buffer)
	if _, err := io.Copy(data, file); err != nil {
		c.Error(err)
	}

	expectedOutput := `console.log("file 2");
console.log("file 1");
console.log("file 3");`

	c.Assert(data.String(), BeginsWith, expectedOutput)
}
