package utils

import (
	"reflect"
	"testing"
)

func TestAvailableIPs(t *testing.T) {
	type args struct {
		cidr string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		want1   string
		wantErr bool
	}{
		{"TestAvailableIp", args{cidr: "192.168.45.0/24"}, 253, "192.168.45.1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNIP, ipServer, gotIpList, err := AvailableIPs(tt.args.cidr)
			if (err != nil) != tt.wantErr {
				t.Errorf("AvailableIPs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNIP != tt.want {
				t.Errorf("AvailableIPs() got = %v, want %v", gotNIP, tt.want)
			}
			if !reflect.DeepEqual(len(gotIpList), tt.want) {
				t.Errorf("AvailableIPs() got1 = %v, want %v", gotIpList, tt.want)
			}
			if ipServer != tt.want1 {
				t.Errorf("AvailableIPs() got1 = %v, want %v", ipServer, tt.want1)
			}
		})
	}
}
