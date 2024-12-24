package main

import "fmt"

func main() {
	// 示例字符串，包含多个地址
	// s := "gddsa广西壮族自治区南宁市清秀区民族大道99号gddsa江苏省南京市雨花台区云密城123号1楼2单元"
	// s := "江苏省南京市雨花台区云密城123号1号楼2单元3层，在法国普罗旺斯的薰衣草田里，紫色的花海在阳光下闪烁着迷人的光芒>。江苏省南京市江宁区正方中路888号4号楼5单元6室，美国大峡谷的壮丽景色令人震撼，红色岩层在阳光下显得格外壮观。>广东省深圳市南山区科技园南区科苑路15号7号楼8单元9户，加拿大班夫国家公园的露易丝湖像一面镜子，倒映着周围的雪山。北京市海淀区中关村大街1号10号楼11单元12层，意大利的五渔村，彩色的房子沿着海岸线排列，美得如同一幅油画。上海市浦东新区陆家嘴环路1000号13号楼14单元15室，瑞士的阿尔卑斯山，雪山覆盖的山峰和绿色的草地构成了一幅美丽的画卷>。浙江省杭州市西湖区文三路555号16号楼17单元18户，澳大利亚的大堡礁，清澈的海水和五彩斑斓的珊瑚构成了一个梦幻般的海底世界。四川省成都市武侯区武侯大道222号19号楼20单元21层，西班牙的托莱多古城，古老的建筑和狭窄的街道让人仿佛穿越回中世纪。山东省济南市历下区经十路167号22号楼23单元24室，巴西的亚马逊雨林，茂密的热带雨林和丰富的生物多样性令人惊叹。湖北省武汉市洪山区珞瑜路188号25号楼26单元27户，美国的黄石国家公园，间歇泉和热泉池散发出独特的地热景观。重庆市渝中区解放碑步行街8号28号楼29单元30层，希腊的圣托里尼岛，白色的房子和蓝色的圆顶教堂点缀在蔚蓝的爱琴海边。广东省深圳市南山区科技园南区科苑路15号31号楼32单元33室，日本的富士山，山顶覆盖着皑皑白雪，周围环绕>着宁静的湖泊。江苏省南京市江宁区正方中路888号34号楼35单元36户，荷兰的风车村，古老的风车在田野中缓缓转动，构成一幅田园诗般的画面。北京市海淀区中关村大街1号37号楼38单元39层，挪威的峡湾，陡峭的山壁和清澈的湖水相互映衬，美得令人窒息。上海市浦东新区陆家嘴环路1000号40号楼41单元42室，新西兰的米尔福德峡湾，碧绿的湖水和高耸的山峰构成>了一幅绝美的自然风光。浙江省杭州市西湖区文三路555号43号楼44单元45户，印度尼西亚的巴厘岛，金色的沙滩和碧蓝的海"
	// s := "美国大峡谷的壮丽景色令人震撼，红色江苏省南京市雨花台区云密城123号1号楼2单元3层，衣草田里，紫>。江苏省南京市江宁区正方中路888号4号楼5单元6室，美国大峡谷的壮丽景色令人震撼，红色岩层在阳光下显得格外壮观。>广东省深圳市南山区科技园南区科苑路15号7号楼8单元9户，加"
	s := "江苏省南京市江宁区正方中路888号4号楼5单元6室，美国大峡谷的壮丽景色令人震撼，红色岩层在阳光下显得格外壮观"

	// 词典，包含可能的地址关键词
	dict := map[string]int{
		"省":   1,
		"自治区": 1, // 添加自治区
		"市":   4,
		"自治州": 4,
		"区":   7,
		"县":   7,
		"镇":   7,
		"乡":   7,
		"城":   9,
		"路":   9,
		"街":   9,
		"里":   9,
		"村":   9,
		"屯":   9,
		"组":   9,
		"道":   9,
		"栋":   12,
		"幢":   12,
		"大厦":  12,
		"小区":  12,
		"广场":  12,
		// "号":   12,
		"号楼": 15,
		"单元": 18,
		"层":  21,
		"室":  21,

		"户": 21,
		"房": 21,

		"州": 98, //追溯3个字
		"元": 99, //追溯2个字
		"厦": 99,
		"场": 99,
		"楼": 99,
	}

	runes := []rune(s) // 转换字符串为 []rune，按字符处理
	left := 0          // 记录左边的关键字
	// mid := 0
	right := 0              //记录右边的关键字
	sl := 0                 //记录左边界
	sr := 0                 //记录右边界
	addresses := []string{} // 存储匹配的地址
	// flag := true
	k_array := []int{}

	//------------------------------------------------------------------------------
	for i := sr; i < len(runes); i++ { //遍历

		if k, v := dict[string(runes[i])]; v { //判断关键字存在
			//判断关键字是否正常，找到真正的关键字
			if k == 99 { //往前追溯两个字的路径
				fmt.Println(string(runes[i-1 : i+1]))
				key, v := dict[string(runes[i-1:i+1])]
				if v {
					k = key

				}
			} else if k == 98 { //往前追溯三个字的路径
				fmt.Println(string(runes[i-2 : i+1]))
				key, v := dict[string(runes[i-2:i+1])]
				if v {
					k = key

				}
			} else if k == 7 { //存在该词且往前追溯两个字的
				fmt.Println(string(runes[i-1 : i+1]))
				key, v := dict[string(runes[i-1:i+1])]
				if v {
					k = key

				}
			}
			k_array = append(k_array, k) //记录k
			//得到真正的关键字大小后，处理字符串
			if k < right { //一直扫描到权重小于前面的字符
				//判断前面捕获的字符串是否合理
				if right-left > 13 { //合理将字符串append到结果集
					addresses = append(addresses, string(runes[sl:sr+1]))
				}

				//合不合理都将left和right重置到该处
				left = 0
				right = 0
				sl = i + 1
				sr = i + 1
				continue

			} else { //扫描正常情况下，持续更新right结果
				if left == right {
					sl = i - 2
				}
				//需要更新sl，之前sl未更新
				right = k
				sr = i - 2
			}

		}
	}
	fmt.Println(addresses)
	fmt.Println(k_array)
	//---------------------------------------------------------------------------------------
	/*for flag {
		// 遍历字符串，逐个字符判断是否匹配词典中的关键字
		for i := sr; i < len(runes); i++ {
			if k, v := dict[string(runes[i])]; v {
				// if k==0{

				// }
				if k >= right {
					if left == right {
						left = k
						sl = i - 2
						if k == 7 { //处理自治区的路径
							fmt.Println(string(runes[i-2 : i+1]))
							k, v := dict[string(runes[i-2:i+1])]
							if v {
								left = k
								sl = i - 6
								k_array = append(k_array, k) //记录k
							}
						}
					} else {
						if k == 99 { //往前追溯两个字的路径
							fmt.Println(string(runes[i-1 : i+1]))
							key, v := dict[string(runes[i-1:i+1])]
							if v {
								right = key
								sr = i
								k_array = append(k_array, key) //记录k
							}
						} else if k == 98 { //往前追溯三个字的路径
							fmt.Println(string(runes[i-2 : i+1]))
							key, v := dict[string(runes[i-2:i+1])]
							if v {
								right = key
								sr = i
								k_array = append(k_array, k) //记录k
							}
						} else if k == 7 { //存在该词且往前追溯两个字的
							fmt.Println(string(runes[i-1 : i+1]))
							key, v := dict[string(runes[i-1:i+1])]
							if v {
								right = key
								sr = i
								k_array = append(k_array, key) //记录k
							} else {
								right = 7
								sr = i
								k_array = append(k_array, k) //记录k
							}
						} else { //正常路径
							// if right < left {
							// 	right = k
							// 	sr = i
							// } else {
							// 	left = k
							// 	sl = i
							// }
							right = k
							sr = i
							k_array = append(k_array, k) //记录k
						}

					}
				} else {
					break
				}
			}

		}
		//除字符超出外的退出条件和append条件
		if sl == sr {
			flag = false
		} else {
			if right-left > 12 { //用以筛选过短的地址
				addresses = append(addresses, string(runes[sl:sr+1]))
			}

		}

		//处理下一次要扫描的字符串
		sl = sr + 2
		sr = sr + 2
		left = 0
		right = 0
		//超出则退出循环
		if sr > len(runes) {
			flag = false
		}
	}
	for _, val := range addresses {
		fmt.Println(val)
	}
	fmt.Println(k_array)
	// fmt.Println(addresses)
	*/
}
