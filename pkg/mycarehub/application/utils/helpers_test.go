package utils

import (
	"encoding/hex"
	"reflect"
	"testing"
	"time"

	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/tj/assert"
)

const (
	// DefaultSaltLen is the length of generated salt for the user is 256
	DefaultSaltLen = 256
	// DefaultKeyLen is the length of encoded key in PBKDF2 function is 512
	DefaultKeyLen = 512
)

func TestEncryptUID(t *testing.T) {
	type args struct {
		rawUID  string
		options *Options
	}

	customOptions := Options{
		// salt length should be greater than 0
		SaltLen:      0,
		Iterations:   2,
		KeyLen:       1,
		HashFunction: extension.DefaultHashFunction,
	}
	tests := []struct {
		name      string
		args      args
		want      string
		want1     string
		wantError bool
	}{
		{
			name: "success: correct default options have been used to encrypt uid",
			args: args{
				rawUID:  "1235",
				options: nil,
			},
			wantError: false,
		},
		{
			name: "failure: incorrect custom options have been used to encrypt uid",
			args: args{
				rawUID:  "1235",
				options: &customOptions,
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			salt, encoded := EncryptUID(tt.args.rawUID, tt.args.options)
			if tt.wantError {
				encodedBytes, err := hex.DecodeString(encoded)
				if err != nil {
					t.Error("Encrypted uid not hex encoded properly")
				}
				assert.Equal(t, len(encodedBytes), DefaultKeyLen)
			}
			if !tt.wantError {
				if !reflect.DeepEqual(len([]byte(salt)), DefaultSaltLen) {
					t.Error("Received length of salt:", len([]byte(salt)), "Expected length of salt:", DefaultSaltLen)
					return
				}
				encodedBytes, err := hex.DecodeString(encoded)
				if err != nil {
					t.Error("Encrypted uid not hex encoded properly")
				}
				assert.Equal(t, len(encodedBytes), DefaultKeyLen)
			}

		})
	}
}

func TestCompareUID(t *testing.T) {
	salt, encoded := EncryptUID("1234", nil)
	type args struct {
		rawUID     string
		salt       string
		encodedUID string
		options    *Options
	}
	tests := []struct {
		name      string
		args      args
		want      bool
		wantError bool
	}{
		{
			name: "success: correct uid supplied that has been encrypted correctly",
			args: args{
				rawUID:     "1234", // this is the same uid that was encrypted
				salt:       salt,
				encodedUID: encoded,
				options:    nil,
			},
			want:      true,
			wantError: false,
		},
		{
			name: "failure: incorrect uid supplied that has been encrypted correctly",
			args: args{
				rawUID:     "4567", // this uid was never encrypted
				salt:       salt,
				encodedUID: encoded,
				options:    nil,
			},
			want:      false,
			wantError: true,
		},
		{
			name: "failure: wrong custom options have been used to encrypt uid",
			args: args{
				rawUID:     "12345",
				salt:       "some random salt",
				encodedUID: "uncoded string",
				options:    nil,
			},
			want:      false,
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isEncypted := CompareUID(tt.args.rawUID, tt.args.salt, tt.args.encodedUID, tt.args.options)
			if !tt.wantError {
				assert.True(t, isEncypted)
				assert.Equal(t, tt.want, isEncypted)
			}
			if tt.wantError {
				assert.False(t, isEncypted)
				assert.Equal(t, tt.want, isEncypted)
			}
		})
	}
}

func TestGetHourMinuteSecond(t *testing.T) {
	type args struct {
		hour   time.Duration
		minute time.Duration
		second time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "Happy case",
			args: args{
				hour:   20,
				minute: 1,
				second: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetHourMinuteSecond(tt.args.hour, tt.args.minute, tt.args.second)
			assert.NotNil(t, got)
		})
	}
}

func TestNextAllowedLoginTime(t *testing.T) {
	type args struct {
		trials int
	}
	tests := []struct {
		name      string
		args      args
		wantError bool
	}{
		{
			name: "Happy case",
			args: args{
				trials: 3,
			},
			wantError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NextAllowedLoginTime(tt.args.trials)
			assert.NotNil(t, got)
		})
	}
}
