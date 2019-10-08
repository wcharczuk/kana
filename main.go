package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/blend/go-sdk/ansi"
)

const (
	maxRepeatHistory     = 5
	weightDefault        = 1.0
	weightIncreaseFactor = 8.0
	weightDecreaseFactor = 2.0
	weightMax            = 512.0
	weightMin            = 0.0625
)

func main() {
	includeKatakana := flagBoolP("katakana", "k", true, "If we should quiz katakana")
	includeHiragana := flagBoolP("hiragana", "h", true, "If we should quiz hiragana")
	limit := flagIntP("limit", "l", 0, "A limit for the number of kana to test")
	flag.Parse()

	var totalAnswered, totalCorrect int
	var times []time.Duration

	var sets []map[string]string
	if *includeKatakana {
		sets = append(sets, katakana)
	}
	if *includeHiragana {
		sets = append(sets, hiragana)
	}
	values := mergeSets(sets...)

	if *limit > 0 {
		values = selectCount(values, *limit)
	}

	weights := createWeights(values)
	total := make(map[string]int)
	incorrect := make(map[string]int)
	kanaTimes := make(map[string][]time.Duration)

	finish := func() {
		fmt.Println()
		fmt.Println("Complete!")
		if totalAnswered > 0 {
			if totalCorrect > 0 {
				fmt.Printf("Total score: %d/%d (%.2f%%)\n", totalCorrect, totalAnswered, (float64(totalCorrect)/float64(totalAnswered))*100)
			} else {
				fmt.Printf("Total score: 0/%d 0.0%%\n", totalAnswered)
			}
			fmt.Printf("Total times: p95 %v, p50: %v\n", percentileOfDuration(times, 95.0).Round(time.Millisecond), percentileOfDuration(times, 50.0).Round(time.Millisecond))
			printResults(total, incorrect, values, weights, kanaTimes)
		}
		os.Exit(0)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fatal(fmt.Errorf("%v", r))
			}
		}()

		var history []string
		var kana, roman string
		var start time.Time
		var elapsed time.Duration
		var isCorrect bool
		var err error

		effectiveMaxRepeatHistory := maxRepeatHistory
		if maxRepeatHistory >= len(values) {
			effectiveMaxRepeatHistory = len(values) >> 1
		}

		for {
			kana, roman = selectWeighted(values, weights)

			if listHas(history, kana) {
				continue
			}
			history = listAddFixedLength(history, kana, effectiveMaxRepeatHistory)

			start = time.Now()

			if isCorrect, err = ask(kana, roman); err != nil {
				if err == errQuit {
					finish()
				}
			} else if isCorrect {
				totalAnswered++
				incrementCount(total, kana)
				decreaseWeight(weights, kana)
				totalCorrect++
				fmt.Printf("(%d/%d) correct!\n", totalCorrect, totalAnswered)
			} else {
				totalAnswered++
				incrementCount(total, kana)
				increaseWeight(weights, kana)
				incrementCount(incorrect, kana)
				fmt.Printf("(%d/%d) incorrect (%s)!\n", totalCorrect, totalAnswered, roman)
			}

			elapsed = time.Since(start)

			kanaTimes[kana] = append(kanaTimes[kana], elapsed)
			times = append(times, elapsed)
		}
	}()

	waitSigInt()
	finish()
}

var katakana = map[string]string{
	"ア": "a",
	"イ": "i",
	"ウ": "u",
	"エ": "e",
	"オ": "o",
	"カ": "ka",
	"キ": "ki",
	"ク": "ku",
	"ケ": "ke",
	"コ": "ko",
	"サ": "sa",
	"シ": "shi",
	"ス": "su",
	"セ": "se",
	"ソ": "so",
	"ナ": "na",
	"ニ": "ni",
	"ヌ": "nu",
	"ネ": "ne",
	"ノ": "no",
	"ハ": "ha",
	"ヒ": "hi",
	"フ": "fu",
	"ヘ": "he",
	"ホ": "ho",
	"マ": "ma",
	"ミ": "mi",
	"ム": "mu",
	"メ": "me",
	"モ": "mo",
	"ラ": "ra",
	"リ": "ri",
	"ル": "ru",
	"レ": "re",
	"ロ": "ro",
	"ワ": "wa",
	"ヲ": "wo",
	"ヤ": "ya",
	"ユ": "yu",
	"ヨ": "yo",
	"タ": "ta",
	"チ": "chi",
	"ツ": "tsu",
	"テ": "te",
	"ト": "to",
	"ン": "n",
	"ガ": "ga",
	"ギ": "gi",
	"グ": "gu",
	"ゲ": "ge",
	"ゴ": "go",
	"ザ": "za",
	"ジ": "zi",
	"ズ": "zu",
	"ゼ": "ze",
	"ゾ": "zo",
	"ヂ": "di",
	"ヅ": "du",
	"デ": "de",
	"ド": "do",
	"バ": "ba",
	"ビ": "bi",
	"ブ": "bu",
	"ベ": "be",
	"ボ": "bo",
	"パ": "pa",
	"ピ": "pi",
	"プ": "pu",
	"ペ": "pe",
	"ポ": "po",
}

var hiragana = map[string]string{
	"あ": "a",
	"い": "i",
	"う": "u",
	"え": "e",
	"お": "o",
	"か": "ka",
	"き": "ki",
	"く": "ku",
	"け": "ke",
	"こ": "ko",
	"さ": "sa",
	"し": "shi",
	"す": "su",
	"せ": "se",
	"そ": "so",
	"た": "ta",
	"ち": "chi",
	"つ": "tsu",
	"て": "te",
	"と": "to",
	"な": "na",
	"に": "ni",
	"ぬ": "nu",
	"ね": "ne",
	"の": "no",
	"は": "ha",
	"ひ": "hi",
	"ふ": "fu",
	"へ": "he",
	"ほ": "ho",
	"ま": "ma",
	"み": "mi",
	"む": "mu",
	"め": "me",
	"も": "mo",
	"や": "ya",
	"ゆ": "yu",
	"よ": "yo",
	"ら": "ra",
	"り": "ri",
	"る": "ru",
	"れ": "re",
	"ろ": "ro",
	"わ": "wa",
	"を": "wo",
	"ん": "n",
	"が": "ga",
	"ぎ": "gi",
	"ぐ": "gu",
	"げ": "ge",
	"ご": "go",
	"ざ": "za",
	"じ": "ji",
	"ず": "zu",
	"ぜ": "ze",
	"ぞ": "zo",
	"だ": "da",
	"ぢ": "di",
	"づ": "du",
	"で": "de",
	"ど": "do",
	"ば": "ba",
	"び": "bi",
	"ぶ": "bu",
	"べ": "be",
	"ぼ": "bo",
	"ぱ": "pa",
	"ぴ": "pi",
	"ぷ": "pu",
	"ぺ": "pe",
	"ぽ": "po",
}

var errQuit = errors.New("should quit")

func promptf(format string, args ...interface{}) string {
	fmt.Fprintf(os.Stdout, format, args...)
	scanner := bufio.NewScanner(os.Stdin)
	var output string
	if scanner.Scan() {
		output = scanner.Text()
	}
	return output
}

func ask(question, expected string) (bool, error) {
	actual := promptf("%s? ", question)
	switch strings.ToLower(strings.TrimSpace(actual)) {
	case "quit", "q":
		return false, errQuit

	}
	if strings.ToLower(actual) == strings.ToLower(expected) {
		return true, nil
	}
	return false, nil
}

func createWeights(values map[string]string) map[string]float64 {
	output := make(map[string]float64)
	for key := range values {
		output[key] = weightDefault
	}
	return output
}

func increaseWeight(weights map[string]float64, value string) {
	if weight, ok := weights[value]; ok {
		if weight < weightMax {
			weights[value] = weight * weightIncreaseFactor
		}
	}
}

func decreaseWeight(weights map[string]float64, value string) {
	if weight, ok := weights[value]; ok {
		if weight <= weightMin {
			return
		}
		weights[value] = weight / weightDecreaseFactor
	}
}

func selectWeighted(values map[string]string, weights map[string]float64) (kana, roman string) {
	// collect "weighted" choices
	type weightedChoice struct {
		Key    string
		Weight float64
	}
	var keys []weightedChoice
	for key := range values {
		keys = append(keys, weightedChoice{
			Key:    key,
			Weight: weights[key],
		})
	}

	// sort by weight ascending
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Weight < keys[j].Weight
	})

	// sum all the weights, assign to indexes
	totals := make([]float64, len(keys))
	var runningTotal float64
	for index, wc := range keys {
		runningTotal += wc.Weight
		totals[index] = runningTotal
	}
	randomValue := rand.Float64() * runningTotal
	randomIndex := sort.SearchFloat64s(totals, randomValue)

	kana = keys[randomIndex].Key
	roman = values[kana]
	return
}

func mergeSets(sets ...map[string]string) map[string]string {
	output := make(map[string]string)
	for _, set := range sets {
		for key, value := range set {
			output[key] = value
		}
	}
	return output
}

func flagBoolP(long, short string, defaultValue bool, usage string) *bool {
	var value bool
	flag.BoolVar(&value, long, defaultValue, usage)
	flag.BoolVar(&value, short, defaultValue, usage+" (shorthand)")
	return &value
}

func flagIntP(long, short string, defaultValue int, usage string) *int {
	var value int
	flag.IntVar(&value, long, defaultValue, usage)
	flag.IntVar(&value, short, defaultValue, usage+" (shorthand)")
	return &value
}

func incrementCount(values map[string]int, key string) {
	if count, ok := values[key]; !ok {
		values[key] = 1
	} else {
		values[key] = count + 1
	}
}

func printResults(total, incorrect map[string]int, values map[string]string, weights map[string]float64, kanaTimes map[string][]time.Duration) {
	if len(values) == 0 {
		return
	}
	columns := []string{
		"Kana (Roman)",
		"Total",
		"Incorrect",
		"Selection Weight",
		"P95",
		"P50",
	}
	var rows [][]string
	for kana, roman := range values {
		totalCount, hasTotal := total[kana]
		incorrectCount := incorrect[kana]
		if hasTotal {
			rows = append(rows, []string{
				fmt.Sprintf("%s (%s)", kana, roman),
				strconv.Itoa(totalCount),
				strconv.Itoa(incorrectCount),
				fmt.Sprintf("%.2f", weights[kana]),
				fmt.Sprint(percentileOfDuration(kanaTimes[kana], 95.0).Round(time.Millisecond)),
				fmt.Sprint(percentileOfDuration(kanaTimes[kana], 50.0).Round(time.Millisecond)),
			})
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		// sort by p95
		iv, _ := time.ParseDuration(rows[i][4])
		jv, _ := time.ParseDuration(rows[j][4])
		return iv > jv
	})

	fmt.Println("Results:")
	fatal(ansi.Table(os.Stdout, columns, rows))
}

func waitSigInt() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
}

func selectCount(values map[string]string, count int) map[string]string {
	if len(values) <= count {
		return values
	}
	output := make(map[string]string)
	for key, value := range values {
		output[key] = value
		if len(output) == count {
			break
		}
	}
	return output
}

// listHas returns if a value is present in a list
func listHas(list []string, value string) bool {
	for _, listValue := range list {
		if listValue == value {
			return true
		}
	}
	return false
}

// addFixedList adds a value to a given list
func listAddFixedLength(list []string, value string, max int) []string {
	list = append(list, value)
	if len(list) < max {
		return list
	}
	return list[1:]
}

func min(values ...int) int {
	if len(values) == 0 {
		return 0
	}
	working := values[0]
	for _, value := range values[1:] {
		if value < working {
			working = value
		}
	}
	return working
}

// percentileOfDuration finds the relative standing in a slice of durations
func percentileOfDuration(input []time.Duration, percentile float64) time.Duration {
	if len(input) == 0 {
		return 0
	}
	return percentileSortedDurations(copySortDurations(input), percentile)
}

// percentileSortedDurations finds the relative standing in a sorted slice of durations
func percentileSortedDurations(sortedInput []time.Duration, percentile float64) time.Duration {
	index := (percentile / 100.0) * float64(len(sortedInput))
	if index == float64(int64(index)) {
		i := int(roundPlaces(index, 0))

		if i < 1 {
			return 0
		}
		return meanDurations([]time.Duration{sortedInput[i-1], sortedInput[i]})
	}

	i := int(roundPlaces(index, 0))
	if i < 1 {
		return time.Duration(0)
	}
	return sortedInput[i-1]
}

// copySortDurations copies and sorts an array of floats.
func copySortDurations(input []time.Duration) []time.Duration {
	inputCopy := copyDurations(input)
	sort.Sort(durations(inputCopy))
	return inputCopy
}

// copyDurations copies an array of time.Duration.
func copyDurations(input []time.Duration) []time.Duration {
	output := make([]time.Duration, len(input))
	copy(output, input)
	return output
}

// durations is an array of durations.
type durations []time.Duration

// Len implements sort.Sorter
func (d durations) Len() int {
	return len(d)
}

// Swap implements sort.Sorter
func (d durations) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Less implements sort.Sorter
func (d durations) Less(i, j int) bool {
	return d[i] < d[j]
}

// meanDurations gets the average of a slice of numbers
func meanDurations(input []time.Duration) time.Duration {
	if len(input) == 0 {
		return 0
	}

	sum := sumDurations(input)
	mean := uint64(sum) / uint64(len(input))
	return time.Duration(mean)
}

// sumDurations adds all the numbers of a slice together
func sumDurations(values []time.Duration) time.Duration {
	var total time.Duration
	for x := 0; x < len(values); x++ {
		total += values[x]
	}

	return total
}

// roundPlaces a float to a specific decimal place or precision
func roundPlaces(input float64, places int) float64 {
	if math.IsNaN(input) {
		return 0.0
	}

	sign := 1.0
	if input < 0 {
		sign = -1
		input *= -1
	}

	rounded := float64(0)
	precision := math.Pow(10, float64(places))
	digit := input * precision
	_, decimal := math.Modf(digit)

	if decimal >= 0.5 {
		rounded = math.Ceil(digit)
	} else {
		rounded = math.Floor(digit)
	}

	return rounded / precision * sign
}

func fatal(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func formatBoolP(value *bool) string {
	if value == nil {
		return "n/a"
	}
	return ansi.ColorLightWhite.Apply(strconv.FormatBool(*value))
}

func formatIntP(value *int) string {
	if value == nil || *value == 0 {
		return "n/a"
	}

	return ansi.ColorLightWhite.Apply(strconv.Itoa(*value))
}
