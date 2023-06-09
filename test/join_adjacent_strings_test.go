package test

import (
	"fmt"
	"image_hub/pkg/utils"
	"testing"
)

func TestJoinAdjacentStrings(t *testing.T) {
	newTexts := []string{
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"儿童节可以向我要红包",
		"我可以给你，但是过几天父亲节",
		"你要是不送我礼物",
		"可别怪爸爸翻脸了",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"今天记得去幼儿园接我哦",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"虽然年龄过了",
		"但智商上我觉得我还是适合过六一的",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"  好看的皮囊千篇一律",
		"可爱的灵魂要过六一",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"我是小朋友这件事你们都知道 吧",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"🍬𝚌𝚑𝚒𝚕𝚍𝚛𝚎𝚗'𝚜𝚍𝚊𝚢🎡✧ ☆ ✩ 🌈 ✩ ☆ ✧要做这个六月最快乐的崽",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"❥ꃅꍏᖘᖘꌩ ꉓꃅꀤ꒒ꀸꋪꍟꈤ'ꌗ\u00a0ꀸ♡ꌩ𝟞𝟙限定超甜小可爱在线卖萌",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"凭我每个月赚的那点钱",
		"我难道不该过六一吗",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		" 会的姿势比你爸妈还多",
		"还要过61呢",
		"这边建议你直接过69",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"情侣头像",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"👇🏻👇🏻👇🏻",
		"你们要的",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"有时候不去追问别人的原因，也是一种体贴。",
		"",
	}

	newTexts1 := utils.JoinAdjacentStrings(newTexts)

	fmt.Printf("len(newTexts): %+v\n", len(newTexts))
	fmt.Printf("newTexts1: %#v\n", newTexts1)
	fmt.Printf("len(newTexts1): %+v\n", len(newTexts1))
}
