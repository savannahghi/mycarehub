package utils_test

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/tj/assert"
)

func TestEncryptPIN(t *testing.T) {
	type args struct {
		rawPwd  string
		options *utils.Options
	}

	customOptions := utils.Options{
		// salt length should be greater than 0
		SaltLen:      0,
		Iterations:   2,
		KeyLen:       1,
		HashFunction: utils.DefaultHashFunction,
	}
	tests := []struct {
		name      string
		args      args
		want      string
		want1     string
		wantError bool
	}{
		{
			name: "success: correct default options have been used to encrypt pin",
			args: args{
				rawPwd:  "1235",
				options: nil,
			},
			wantError: false,
		},
		{
			name: "failure: incorrect custom options have been used to encrypt pin",
			args: args{
				rawPwd:  "1235",
				options: &customOptions,
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			salt, encoded := utils.EncryptPIN(tt.args.rawPwd, tt.args.options)
			if tt.wantError {
				if reflect.DeepEqual(len([]byte(salt)), utils.DefaultSaltLen) {
					t.Error("Received length of salt:", len([]byte(salt)), "Expected length of salt:", utils.DefaultSaltLen)
					return
				}
				encodedBytes, err := hex.DecodeString(encoded)
				if err != nil {
					t.Error("Encrypted Password not hex encoded properly")
				}
				if reflect.DeepEqual(len(encodedBytes), utils.DefaultKeyLen) {
					t.Error("Received length of password:", len(encodedBytes), "Expected length of password:", utils.DefaultKeyLen)
					return
				}
			}
			if !tt.wantError {
				if !reflect.DeepEqual(len([]byte(salt)), utils.DefaultSaltLen) {
					t.Error("Received length of salt:", len([]byte(salt)), "Expected length of salt:", utils.DefaultSaltLen)
					return
				}
				encodedBytes, err := hex.DecodeString(encoded)
				if err != nil {
					t.Error("Encrypted Password not hex encoded properly")
				}
				if !reflect.DeepEqual(len(encodedBytes), utils.DefaultKeyLen) {
					t.Error("Received length of password:", len(encodedBytes), "Expected length of password:", utils.DefaultKeyLen)
					return
				}
			}

		})
	}
}

func TestComparePIN(t *testing.T) {
	salt, encoded := utils.EncryptPIN("1234", nil)
	type args struct {
		rawPwd     string
		salt       string
		encodedPwd string
		options    *utils.Options
	}
	tests := []struct {
		name      string
		args      args
		want      bool
		wantError bool
	}{
		{
			name: "success: correct pin supplied that has been encrypted correctly",
			args: args{
				rawPwd:     "1234", // this is the same password that was encrypted
				salt:       salt,
				encodedPwd: encoded,
				options:    nil,
			},
			want:      true,
			wantError: false,
		},
		{
			name: "failure: incorrect pin supplied that has been encrypted correctly",
			args: args{
				rawPwd:     "4567", // this password was never encrypted
				salt:       salt,
				encodedPwd: encoded,
				options:    nil,
			},
			want:      false,
			wantError: true,
		},
		{
			name: "failure: wrong custom options have been used to encrypt pin",
			args: args{
				rawPwd:     "12345",
				salt:       "some random salt",
				encodedPwd: "uncoded string",
				options:    nil,
			},
			want:      false,
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isEncypted := utils.ComparePIN(tt.args.rawPwd, tt.args.salt, tt.args.encodedPwd, tt.args.options)
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
