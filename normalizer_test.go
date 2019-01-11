package normalizer

import "testing"

func Test_NormalizeNeologd(t *testing.T) {
	assert := func(expected, actual string) {
		if actual != expected {
			t.Errorf("got %v\nwant %v", actual, expected)
		}
	}
	// test
	assert("0123456789", NormalizeNeologd("０１２３４５６７８９"))
	assert("ABCDEFGHIJKLMNOPQRSTUVWXYZ", NormalizeNeologd("ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺ"))
	assert("abcdefghijklmnopqrstuvwxyz", NormalizeNeologd("ａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ"))
	assert("!\"#$%&'()*+,-./:;<>?@[¥]^_`{|}", NormalizeNeologd("！”＃＄％＆’（）＊＋，−．／：；＜＞？＠［￥］＾＿｀｛｜｝"))
	assert("＝。、・「」", NormalizeNeologd("＝。、・「」"))
	assert("ハンカク", NormalizeNeologd("ﾊﾝｶｸ"))
	assert("o-o", NormalizeNeologd("o₋o"))
	assert("majikaー", NormalizeNeologd("majika━"))
	assert("わい", NormalizeNeologd("わ〰い"))
	assert("スーパー", NormalizeNeologd("スーパーーーー"))
	assert("!#", NormalizeNeologd("!#"))
	assert("ゼンカクスペース", NormalizeNeologd("ゼンカク　スペース"))
	assert("おお", NormalizeNeologd("お             お"))
	assert("おお", NormalizeNeologd("      おお"))
	assert("おお", NormalizeNeologd("おお      "))
	assert("検索エンジン自作入門を買いました!!!", NormalizeNeologd("検索 エンジン 自作 入門 を 買い ました!!!"))
	assert("アルゴリズムC", NormalizeNeologd("アルゴリズム C"))
	assert("PRML副読本", NormalizeNeologd("　　　ＰＲＭＬ　　副　読　本　　　"))
	assert("Coding the Matrix", NormalizeNeologd("Coding the Matrix"))
	assert("南アルプスの天然水Sparking Lemonレモン一絞り", NormalizeNeologd("南アルプスの　天然水　Ｓｐａｒｋｉｎｇ　Ｌｅｍｏｎ　レモン一絞り"))
	assert("南アルプスの天然水-Sparking*Lemon+レモン一絞り", NormalizeNeologd("南アルプスの　天然水-　Ｓｐａｒｋｉｎｇ*　Ｌｅｍｏｎ+　レモン一絞り"))
}
