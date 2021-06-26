package extensions

import "sort"

// Checks if inside the array exists the item
func ArrayIncludes(arr []string, predicate string) bool {
	sort.Strings(arr)
	i := sort.SearchStrings(arr, predicate)
	return i < len(arr) && arr[i] == predicate
}

// Checks if inside the array exists the item
func ArrayIncludesT(arr []interface{}, predicate interface{}) bool {
	for i := range arr {
		if predicate == arr[i] {
			return true
		}
	}

	return false
}
