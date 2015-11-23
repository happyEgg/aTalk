/*
敏感词判断
*/
package controller

//敏感词判断
func SensitiveWords(name interface{}) bool {
	words := []string{"admin", "毛泽东", "周恩来"}
	for _, word := range words {
		if name == word {
			return false
		}
	}
	return true
}
