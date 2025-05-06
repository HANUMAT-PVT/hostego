package utilities

import (
	"encoding/json"
	"errors"
	"hash/fnv"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// func GetNonWebNewRelicTxn(txnName string) *newrelic.Transaction {
// 	app := newrelic_setup.GetNewRelicApp(viper.GetString(config_constants.VKEYS_NEWRELIC_APP_NAME))
// 	return app.StartTransaction(txnName)
// }

func GoRoutinePanicHandler() {
	if r := recover(); r != nil {
		log.Errorf("goroutine paniqued, recovering with panic handler: %+v , stack: %v", r, string(debug.Stack()))
	}
}

func ConvertStructToMapStringInterface(ipStruct interface{}) (map[string]interface{}, error) {
	b, _ := json.Marshal(&ipStruct)

	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	return m, err
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func IsContains(inputSlice []string, eleToFound string) bool {
	for _, ele := range inputSlice {
		if ele == eleToFound {
			return true
		}
	}
	return false
}

func GetStringSliceFromInterfaceSlice(input []interface{}) []string {
	var stringSlice []string
	for _, stringEle := range input {
		stringSlice = append(stringSlice, stringEle.(string))
	}
	return stringSlice
}

func ConvertInterfaceToMapStringInterface(value interface{}) (map[string]interface{}, error) {
	jsonStr, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	jsonValue := new(map[string]interface{})
	err = json.Unmarshal(jsonStr, &jsonValue)
	if err != nil {
		return nil, err
	}

	return *jsonValue, nil
}

func HashStringToInt64(text string) (number int64) {
	hash := fnv.New64a()
	hash.Write([]byte(text))
	value := hash.Sum64()
	return int64(value)
}

func ConvertEpochToDate(epochStr string) (date string, err error) {
	epoch, err := strconv.ParseInt(epochStr, 10, 64)
	if len(epochStr) == 13 {
		epoch = epoch / 1000
	}
	// Convert epoch to time.Time
	t := time.Unix(epoch, 0)

	// Format time.Time to desired format
	return t.Format("2006-01-02"), nil
}

func ConvertEpochToTime(epochStr string) (time.Time, error) {
	// Parse the epoch string to an integer
	epoch, err := strconv.ParseInt(epochStr, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	// Adjust for millisecond epoch if necessary
	if len(epochStr) == 13 {
		epoch = epoch / 1000
	}

	// Convert epoch to time.Time
	return time.Unix(epoch, 0), nil
}

func ConvertEpochIntToTime(epoch int64) (time.Time, error) {
	// Parse the epoch string to an integer

	// Convert epoch to time.Time
	return time.UnixMilli(epoch), nil
}

func ConvertTimeToEpoch(t *time.Time, inMilliseconds bool) int64 {
	if t == nil {
		return 0
	}
	// Parse the epoch string to an integer
	if inMilliseconds {
		return t.UnixMilli()
	}
	return t.Unix()
}

func GetBaseAndQuoteAsset(symbol string) (string, string, error) {
	assets := strings.Split(symbol, "/")
	if len(assets) != 2 {
		return "", "", errors.New("invalid symbol: " + symbol)
	}

	return strings.ToLower(strings.ReplaceAll(assets[0], " ", "")), strings.ToLower(strings.ReplaceAll(assets[1], " ", "")), nil
}

func IsNumeric(input string) bool {
	re := regexp.MustCompile(`^\d+$`)
	return re.MatchString(input)
}

func CheckTargetIsThereInKey(config map[string]interface{}, key string, target string) bool {
	/*
		example of config :
		{
				"EXCHANGES_ELIGIBLE":[
					"COINSWITCHX"
				]
		}
	*/

	keyConfig, ok := config[key].([]interface{})
	if !ok {
		return true
	}
	if len(keyConfig) == 0 {
		return true
	}

	for _, data := range keyConfig {
		str, ok := data.(string)
		if !ok {
			return false
		}
		if str == target {
			return true
		}
	}

	return false

}
