package params

// 需要过滤的图片img data-src属性值
// 一般不常用,通过spiders/webchat_微信号.go 里面写规则来过滤, 但是一些简单或临时的,可以通过在下面添加来过滤不需要的图
var SectionDirtyImgDataSrc = []string{
	"https://mmbiz.qlogo.cn/mmbiz_gif/L0owf7fTI9PcicvDCqOv1gkk8zfQSibB4JwmPxIG7l1NF91P4XvAwQO27jbyDVuctn5Wzk8qzkn0T5OEPjAuicoqA/0?wx_fmt=gif",
	"https://mmbiz.qpic.cn/mmbiz_jpg/Hw8QGloibSVHiaIBeibODxzpYEX0zSXvDrVmNlYMpPdJKA4d0LyEnnwaKicQ6TRoRsiagdjhIWPLOCNgnQ9vibXOLFkA/640?wx_fmt=jpeg",
	"https://mmbiz.qpic.cn/mmbiz_jpg/Hw8QGloibSVHiaIBeibODxzpYEX0zSXvDrVWoRZtuxwLCEkAia8jG9ib1Fvj1HlNwvMaLpH6UNViahyJicmFrJL3qaWIQ/640?wx_fmt=jpeg",
	"https://mmbiz.qpic.cn/mmbiz_gif/myr7atqOsbJOOvrgzialOvWwYVNdhgYay9tSRwMvnvk3I2PQRv8qr4vW7c1smFicaEBuaicIwNgTecbAnTWRm7trQ/640",

	//TODO: 这个不对,为什么没有通过重复图片把这个检测出来?
	// "https://mmbiz.qpic.cn/mmbiz_png/e6u7rIactyalTzFHnvuLtktqfQMUtD0ibZOSmQs7k7IpM1AhTjCGScXHJUuomR1MiaYlwrux7fmzxNmcx6t2uAZQ/640?wx_fmt=png&wx_lazy=1&wx_co=1",
}
