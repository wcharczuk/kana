package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/ansi/slant"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/mathutil"
	"github.com/blend/go-sdk/sh"
)

const (
	showIncorrect = 10
	maxHistory    = 5
	weightFactor  = 2.0
	weightMin     = 0.125
)

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
	"フ": "hu",
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

func prompt(kana, roman string) bool {
	actual := sh.Promptf("%s? ", kana)
	if strings.ToLower(actual) == strings.ToLower(roman) {
		fmt.Println("correct!")
		return true
	}
	fmt.Printf("incorrect! (%s)\n", roman)
	return false
}

// createWeights creates a weight map for a given set of values.
func createWeights(values map[string]string) map[string]float64 {
	output := make(map[string]float64)
	for key := range values {
		output[key] = 1.0
	}
	return output
}

func increaseWeight(weights map[string]float64, value string) {
	if weight, ok := weights[value]; ok {
		weights[value] = weight * weightFactor
	}
}

func decreaseWeight(weights map[string]float64, value string) {
	if weight, ok := weights[value]; ok {
		if weight <= weightMin {
			return
		}
		weights[value] = weight / weightFactor
	}
}

func random(weights map[string]float64, values map[string]string) (kana, roman string) {
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

func merge(sets ...map[string]string) map[string]string {
	output := make(map[string]string)
	for _, set := range sets {
		for key, value := range set {
			output[key] = value
		}
	}
	return output
}

func shortBoolP(long, short string, defaultValue bool, usage string) *bool {
	var value bool
	flag.BoolVar(&value, long, defaultValue, usage)
	flag.BoolVar(&value, short, defaultValue, usage+" (shorthand)")
	return &value
}

func shortIntP(long, short string, defaultValue int, usage string) *int {
	var value int
	flag.IntVar(&value, long, defaultValue, usage)
	flag.IntVar(&value, short, defaultValue, usage+" (shorthand)")
	return &value
}

func incrementWrong(wrong map[string]int, key string) {
	if count, ok := wrong[key]; !ok {
		wrong[key] = 1
	} else {
		wrong[key] = count + 1
	}
}

func printWrong(wrong map[string]int, values map[string]string, weights map[string]float64) {
	if len(wrong) == 0 {
		return
	}
	columns := []string{
		"Kana (Roman)",
		"Count",
		"Selection Weight",
	}
	var rows [][]string
	for kana, count := range wrong {
		rows = append(rows, []string{
			fmt.Sprintf("%s (%s)", kana, values[kana]),
			strconv.Itoa(count),
			fmt.Sprintf("%.2f", weights[kana]),
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][1] > rows[j][1]
	})

	effectiveLimit := min(showIncorrect, len(rows))

	fmt.Printf("Incorrect Answers (Top %d):\n", effectiveLimit)
	if len(rows) > effectiveLimit {
		fatal(ansi.Table(os.Stdout, columns, rows[:effectiveLimit]))
	} else {
		fatal(ansi.Table(os.Stdout, columns, rows))
	}
}

func waitSigInt() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
}

func applyLimit(values map[string]string, limit int) map[string]string {
	if len(values) <= limit {
		return values
	}
	output := make(map[string]string)
	for key, value := range values {
		output[key] = value
		if len(output) == limit {
			break
		}
	}
	return output
}

// inList returns if a value is present in a list
func inList(list []string, value string) bool {
	for _, listValue := range list {
		if listValue == value {
			return true
		}
	}
	return false
}

// addFixedList adds a value to a given list
func addFixedList(list []string, value string, max int) []string {
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

func main() {
	includeKatakana := shortBoolP("katakana", "k", true, "If we should quiz katakana")
	includeHiragana := shortBoolP("hiragana", "h", true, "If we should quiz hiragana")
	limit := shortIntP("limit", "l", 0, "A limit for the number of kana to test")
	flag.Parse()

	slant.Print(os.Stdout, "KANA")
	fmt.Printf("katakana: %v, hiragana: %v, limit: %v\n", formatBoolP(includeKatakana), formatBoolP(includeHiragana), formatIntP(limit))

	var correct, total int
	var times []time.Duration

	var sets []map[string]string
	if *includeKatakana {
		sets = append(sets, katakana)
	}
	if *includeHiragana {
		sets = append(sets, hiragana)
	}
	values := merge(sets...)

	if *limit > 0 {
		values = applyLimit(values, *limit)
	}

	weights := createWeights(values)
	wrong := make(map[string]int)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fatal(ex.New(r))
			}
		}()

		var history []string
		var kana, roman string
		var start time.Time

		effectiveMaxHistory := maxHistory
		if *limit > 0 && maxHistory >= *limit {
			effectiveMaxHistory = (*limit) >> 1
		}

		for {
			kana, roman = random(weights, values)

			if inList(history, kana) {
				continue
			}
			history = addFixedList(history, kana, effectiveMaxHistory)

			start = time.Now()
			if prompt(kana, roman) {
				correct++
				decreaseWeight(weights, kana)
			} else {
				increaseWeight(weights, kana)
				incrementWrong(wrong, kana)
			}
			times = append(times, time.Since(start))
			total++
		}
	}()

	waitSigInt()
	if total > 0 {
		fmt.Println()
		if correct > 0 {
			fmt.Printf("Session score: %d/%d (%.2f%%)\n", correct, total, (float64(correct)/float64(total))*100)
		} else {
			fmt.Printf("Session score: 0/%d 0.0%%\n", total)
		}
		fmt.Printf("Session times: p95 %v, p50: %v\n", mathutil.PercentileOfDuration(times, 95.0).Round(time.Millisecond), mathutil.PercentileOfDuration(times, 50.0).Round(time.Millisecond))
		printWrong(wrong, values, weights)
	}
}
