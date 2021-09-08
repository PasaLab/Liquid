package main

type JobSorter []Job

func (s JobSorter) Len() int {
	return len(s)
}
func (s JobSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s JobSorter) Less(i, j int) bool {
	return s[i].CreatedAt < s[j].CreatedAt
}
