package api

//待优化
func EdgeSetScheduler(es *EdgeSet) error {
	for key, _ := range es.SetOnline {
		es.SetScore[key] = 1
	}
	return nil
}

//待优化
func EdgeNodeScheduler(en *EdgeNode) error {
	for key, _ := range en.NodeOnline {
		en.NodeScore[key] = 1
	}
	return nil
}
