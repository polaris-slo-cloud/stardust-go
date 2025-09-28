package helper

import "strings"

func ExtendFilename(filename string, add string) string {
	lastDotIndex := strings.LastIndex(filename, ".")

	if lastDotIndex == -1 {
		return filename + add
	}

	base := filename[:lastDotIndex]
	extension := filename[lastDotIndex:]

	return base + add + extension
}
