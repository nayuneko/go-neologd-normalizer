package normalizer

/**
 * https://github.com/neologd/mecab-ipadic-neologd/wiki/Regexp.ja
 */

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/text/unicode/norm"
)

var (
	// ハイフンマイナスっぽい文字を置換
	repHypon = makeReplacerOldnews([]string{"˗", "֊", "‐", "‑", "‒", "–", "⁃", "⁻", "₋", "−"}, "-")
	// 長音記号っぽい文字を置換
	repChoon = makeReplacerOldnews([]string{"﹣", "－", "ｰ", "—", "―", "─", "━"}, "ー")
	// チルダっぽい文字は削除
	repChil = makeReplacerOldnews([]string{"~", "∼", "∾", "〜", "〰", "～"}, "")
	// 前処理で一旦記号を半角から全角に置換する
	repKigouBefore = makeTrReplacer(
		"!\"#$%&'()*+,-./:;<=>?@[¥]^_`{|}~｡､･｢｣",
		"！”＃＄％＆’（）＊＋，－．／：；＜＝＞？＠［￥］＾＿｀｛｜｝〜。、・「」",
	)
	// 全てのReplacerを結合
	repAll = makeReplacer(repHypon, repChoon, repChil, repKigouBefore)

	// 全角->半角に置換する記号
	repKigouAfter = makeTrReplacer(
		"！”＃＄％＆’（）＊＋，－．／：；＜＞？＠［￥］＾＿｀｛｜｝〜",
		"!\"#$%&'()*+,-./:;<>?@[¥]^_`{|}~",
	)
	repKigou = strings.NewReplacer(repKigouAfter...)

	// 1回以上連続する長音記号は1回に置換
	regChoon2 = regexp.MustCompile(`[ー]{2,}`)
	// 1つ以上の半角スペースは、1つの半角スペースに置換
	regSpace2 = regexp.MustCompile(`[ ]{2,}`)

	// ひらがな・全角カタカナ・半角カタカナ・漢字・全角記号
	blocks = `\x{4E00}-\x{9FFF}\x{3040}-\x{309F}\x{30A0}-\x{30FF}\x{3000}-\x{303F}\x{FF00}-\x{FFEF}`
	// 半角英数字
	latin = `\x{0000}-\x{007F}`
	// 「ひらがな・全角カタカナ・半角カタカナ・漢字・全角記号」間に含まれる半角スペースは削除
	regB2B = regexp.MustCompile(fmt.Sprintf("([%s]) ([%s])", blocks, blocks))
	// 「ひらがな・全角カタカナ・半角カタカナ・漢字・全角記号」と「半角英数字」の間に含まれる半角スペースは削除
	regB2L = regexp.MustCompile(fmt.Sprintf("([%s]) ([%s])", blocks, latin))
	regL2B = regexp.MustCompile(fmt.Sprintf("([%s]) ([%s])", latin, blocks))
)

func makeReplacerOldnews(chars []string, to string) []string {
	Oldnews := make([]string, len(chars)*2)
	for i, c := range chars {
		ii := i * 2
		Oldnews[ii] = c
		Oldnews[ii+1] = to
	}
	return Oldnews
}

func makeReplacer(oldnews ...[]string) *strings.Replacer {
	var l int
	for _, oldnew := range oldnews {
		l += len(oldnew)
	}
	params := make([]string, l)
	for _, oldnew := range oldnews {
		params = append(params, oldnew...)
	}
	return strings.NewReplacer(params...)
}

func makeTrReplacer(f, t string) []string {
	ff := []rune(f)
	tt := []rune(t)
	l1 := len(ff)
	l2 := len(tt)
	if l1 != l2 {
		panic(fmt.Sprintf("tr params failure: %d, %d", l1, l2))
	}
	oldNews := make([]string, l1+l2)
	for i, c := range ff {
		ii := i * 2
		oldNews[ii] = string(c) // string(ff[i])
		oldNews[ii+1] = string(tt[i])
	}
	return oldNews
}

func removeSpaceBetween(r *regexp.Regexp, s string) string {
	for r.FindStringIndex(s) != nil {
		s = r.ReplaceAllString(s, "$1$2")
	}
	return s
}

// NormalizeNeologd Neologd用にノーマライズ
func NormalizeNeologd(s string) string {
	// 半角->全角置換
	s = norm.NFKC.String(s)
	// もろもろ置換
	s = repAll.Replace(s)
	// 2文字以上の長音を1文字に
	s = regChoon2.ReplaceAllString(s, "ー")
	// 全角スペースを半角スペースに
	s = strings.Replace(s, "　", " ", -1)
	// 2文字以上のスペースを1文字に
	s = regSpace2.ReplaceAllString(s, " ")
	// 「全角 全角」、「全角 半角」、「半角 全角」の間のスペースを削除する
	s = removeSpaceBetween(regB2B, s)
	s = removeSpaceBetween(regB2L, s)
	s = removeSpaceBetween(regL2B, s)
	// 全角->半角 記号置換
	s = repKigou.Replace(s)
	// trim
	s = strings.Trim(s, " ")
	return s
}
