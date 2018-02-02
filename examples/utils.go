package utils

import (
	twizo "github.com/twizoapi/lib-api-go"
	"os"
	"flag"
	"fmt"
	"strings"
	"bufio"
	"reflect"
	"bytes"
	"errors"
	"regexp"
	"encoding/json"
)

// exported settings
var SuppliedApiKey 		string
var SuppliedRegion              string

func getRegionStrings() ([]string) {
	regions := twizo.GetRegions();
	keys := make([]string, 0, len(regions))
	for k := range regions {
		keys = append(keys, string(k))
	}
	return keys
}

func Main() {
	key       := flag.String("key", "", "the api key (required)")
	region    := flag.String("region", string(twizo.APIRegionDefault), fmt.Sprintf("the region to use [%s]", strings.Join(getRegionStrings(), ",")))
	isVerbose := flag.Bool("verbose", false, "show interaction with api")
	isHelp    := flag.Bool("help", false, "Show the help")
	flag.Parse()

	if *isHelp == true {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *key == "" {
		flag.PrintDefaults()

		fmt.Fprintln(os.Stderr, "Error: No api key supplied")
		os.Exit(1)
	}

	// hack to allow extra environments to be passed in for testing
	if len(flag.Args()) == 1 {
		regionData := strings.Split(flag.Arg(0), ":")
		if len(regionData) == 2 {
			twizo.AddHostForRegion(twizo.APIRegion(regionData[0]), regionData[1])
			// switch region if default to the added one
			if *region == string(twizo.APIRegionDefault) {
				*region = regionData[0]
			}
		}
	}

	// check to see if we have a valid region
	regions := twizo.GetRegions()
	if _, ok := regions[twizo.APIRegion(*region)]; !ok {
		flag.PrintDefaults()

		fmt.Fprintf(os.Stderr, "Error: No such region [%s]\n", *region)
		os.Exit(1)
	}

	SuppliedApiKey = *key
	SuppliedRegion = *region

	if *isVerbose {
		twizo.DebugLogger.SetOutput(os.Stdout)
	}
}

func AskForInput(messageStr string, defaultStr string) (string, error) {
	fmt.Print(messageStr)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return defaultStr, err
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return defaultStr, nil
	}
	return text, nil
}

func Call(m map[string]interface{}, name string, params ... interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}

func AsString(s interface{}, params ...string) (string) {
	switch sType := s.(type) {
	case string:
		return sType
	case *string:
		if (sType == nil && len(params) > 0) {
			return params[0]
		} else if (sType == nil) {
			return "nil"
		} else {
			return *sType
		}
	}

	return "unknown"
}

func GetJsonFor(s interface{}) (string) {
	bytes, _ := json.MarshalIndent(s, "", "\t")

	return string(bytes)
}

func DumpStruct(s interface{}) (string) {
	var buffer bytes.Buffer
	sType := reflect.TypeOf(s)
	sValue := reflect.ValueOf(s)

	buffer.WriteString(fmt.Sprintf("%s\n", sType))


	re := regexp.MustCompile("^Get(?P<property>.+)$")
	for i := 0; i < sType.NumMethod(); i++ {
		method := sType.Method(i)
		// call only the Get methods
		result := re.FindStringSubmatch(method.Name)
		if (len(result) != 2) {
			// it is not a get method
			continue
		}

		value := sValue.MethodByName(method.Name).Call([]reflect.Value{})
		buffer.WriteString(fmt.Sprintf("[%s] -> [%s]\n", result[1], AsString(value[0].Interface())))

	}

	return buffer.String()
}


