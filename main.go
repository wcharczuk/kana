package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"

	"github.com/blend/go-sdk/ansi/slant"
	"github.com/blend/go-sdk/sh"
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

const beginning = "\033[F"
const up = "\033[A"

func prompt(kana, roman string) bool {
	actual := sh.Promptf("\n%s? ", kana)
	if strings.ToLower(actual) == strings.ToLower(roman) {
		fmt.Printf("correct!")
		return true
	}
	fmt.Printf("incorrect! (%s)", roman)
	return false
}

func random(values map[string]string) (kana, roman string) {
	var keys []string
	for key := range values {
		keys = append(keys, key)
	}
	randomIndex := rand.Intn(len(keys))
	kana = keys[randomIndex]
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

func main() {
	includeKatakana := shortBoolP("katakana", "k", true, "If we should quiz katakana")
	includeHiragana := shortBoolP("hiragana", "h", true, "If we should quiz hiragana")
	flag.Parse()

	slant.Print(os.Stdout, "KANA")
	fmt.Printf("katakana: %v, hiragana: %v\n", *includeKatakana, *includeHiragana)

	var correct, total int
	go func() {
		var sets []map[string]string
		if *includeKatakana {
			sets = append(sets, katakana)
		}
		if *includeHiragana {
			sets = append(sets, hiragana)
		}
		final := merge(sets...)
		var last, kana, roman string
		for {
			kana, roman = random(final)
			if kana == last {
				continue
			}
			last = kana
			if prompt(kana, roman) {
				correct++
			}
			total++
		}
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
	if total > 0 {
		fmt.Printf("\nSession totals: %d/%d (%.2f%%)\n", correct, total, (float64(correct)/float64(total))*100)
	}
}
