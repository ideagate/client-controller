package kubernetes

import (
	"reflect"
	"testing"

	"k8s.io/client-go/rest"
)

func Test_getKubeConfig(t *testing.T) {
	t.Skip()
	tests := []struct {
		name     string
		mockFunc func()
		want     *rest.Config
		wantErr  bool
	}{
		{
			name:     "positive",
			mockFunc: func() {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			got, err := getKubeConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("getKubeConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getKubeConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
