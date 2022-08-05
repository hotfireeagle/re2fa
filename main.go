package main

func Parse(re string) {
	for _, c := range re {
		if isValueLetter(c) {

		}
	}
}

func isValueLetter(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func main() {

}
