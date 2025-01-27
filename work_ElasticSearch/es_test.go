package work_ElasticSearch

import "testing"

func Test_GetESClient(t *testing.T) {
	_, err := GetElasticSearchClient()
	if err != nil {
		return
	}
}
