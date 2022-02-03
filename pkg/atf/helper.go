package atf

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tidwall/gjson"
)

// validateResource validates the resource exists in state file
func validateResource(name string, validations []validation, getAPI GetAPIFunc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("[Validate Resource] resource %s not found", name)
		}

		id := rs.Primary.Attributes["id"]
		if id == "" {
			return fmt.Errorf("resource %s ID is not set", name)
		}

		resp, err := getAPI(rs.Primary.Attributes)
		if err != nil {
			return err
		}

		jsonBody, err := json.Marshal(resp)
		if err != nil {
			return err
		}
		jsonStr := string(jsonBody)
		for _, v := range validations {
			var result string
			if v.isJSON {
				result = gjson.Get(jsonStr, v.key).String()
			} else {
				result = rs.Primary.Attributes[v.key]
			}
			if result != fmt.Sprint(v.value) {
				return fmt.Errorf("validation failed for %s. On API response, expected %s = %s, but got %v",
					name, v.key, result, v.value)
			}
		}

		return nil
	}
}

// getLocalName truncates hpegl_vmaas_ and returns back remaining.
func getLocalName(res string) string {
	return res[len("hpegl_vmaas_"):]
}

func getTag(isResource bool) string {
	if isResource {
		return "resources"
	}

	return "data-sources"
}

func getType(isResource bool) string {
	if isResource {
		return "resource"
	}

	return "data"
}

func toInt(str string) int {
	i, _ := strconv.Atoi(str)

	return i
}

// newRand will create different random number if it is called from different
// go routine. This will ensure there will be no collision in random number and
// Parallel testing is possible
func newRand() *rand.Rand {
	s := myCaller()
	m := md5.New()
	_, err := m.Write([]byte(s))
	if err != nil {
		panic(err)
	}
	sourceStr := m.Sum(nil)
	var sourceInt int64
	for _, i := range sourceStr {
		sourceInt += int64(i)
	}

	return rand.New(rand.NewSource(sourceInt + time.Now().Unix()))
}

func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}

// MyCaller returns the caller of the function that called it :)
func myCaller() string {
	// Skip GetCallerFunctionName and the function to get the caller of
	return getFrame(6).Function
}

func path(v ...interface{}) string {
	if len(v) == 0 {
		return ""
	}

	var p string
	for _, val := range v {
		p += "." + fmt.Sprint(val)
	}

	return p[1:]
}
