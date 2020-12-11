package models

type ClusterSlice []string

func (c ClusterSlice) Len() int { return len(c) }
func (c ClusterSlice) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c ClusterSlice) Less(i, j int) bool { return c[i] < c[j] }
