package dto

import (
	"reflect"
	"testing"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/scalarutils"
)

func TestUsernameInput_ParseUsername(t *testing.T) {
	type fields struct {
		Username string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Happy case: valid username",
			fields: fields{
				Username: "username",
			},
			wantErr: false,
			want:    "username",
		},
		{
			name: "Sad case: invalid username",
			fields: fields{
				Username: "user@name",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Sad case: missing username",
			fields: fields{
				Username: "",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CMDUsernameInput{
				Username: tt.fields.Username,
			}
			got, err := u.ParseUsername()
			if (err != nil) != tt.wantErr {
				t.Errorf("UsernameInput.ParseUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UsernameInput.ParseUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNameInput_ParseName(t *testing.T) {
	type fields struct {
		FirstName string
		LastName  string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Happy case: valid input",
			fields: fields{
				FirstName: "user",
				LastName:  "name",
			},
			want:    "user name",
			wantErr: false,
		},
		{
			name: "Sad case: missing first name",
			fields: fields{
				LastName: "name",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Sad case: missing last name",
			fields: fields{
				FirstName: "user",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := CMDNameInput{
				FirstName: tt.fields.FirstName,
				LastName:  tt.fields.LastName,
			}
			got, err := s.ParseName()
			if (err != nil) != tt.wantErr {
				t.Errorf("NameInput.ParseName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NameInput.ParseName() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestDateInput_ParseDate(t *testing.T) {
	type fields struct {
		Year  string
		Month string
		Day   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    scalarutils.Date
		wantErr bool
	}{
		{
			name: "Happy case: valid date",
			fields: fields{
				Year:  "2000",
				Month: "1",
				Day:   "1",
			},
			want: scalarutils.Date{
				Year:  2000,
				Month: 1,
				Day:   1,
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid year",
			fields: fields{
				Year:  "invalid",
				Month: "1",
				Day:   "1",
			},
			want:    scalarutils.Date{},
			wantErr: true,
		},
		{
			name: "Sad case: invalid month",
			fields: fields{
				Year:  "2000",
				Month: "invalid",
				Day:   "1",
			},
			want:    scalarutils.Date{},
			wantErr: true,
		},
		{
			name: "Sad case: invalid day",
			fields: fields{
				Year:  "2000",
				Month: "1",
				Day:   "invalid",
			},
			want:    scalarutils.Date{},
			wantErr: true,
		},
		{
			name:    "Sad case: missing date",
			fields:  fields{},
			want:    scalarutils.Date{},
			wantErr: true,
		},
		{
			name: "Sad case: invalid date",
			fields: fields{
				Year:  "2000000",
				Month: "1",
				Day:   "1",
			},
			want:    scalarutils.Date{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &CMDDateInput{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			got, err := d.ParseDate()
			if (err != nil) != tt.wantErr {
				t.Errorf("DateInput.ParseDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DateInput.ParseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenderInput_ParseGender(t *testing.T) {
	type fields struct {
		Gender string
	}
	tests := []struct {
		name    string
		fields  fields
		want    enumutils.Gender
		wantErr bool
	}{
		{
			name: "Happy Case: valid gender",
			fields: fields{
				Gender: "Male",
			},
			want:    enumutils.GenderMale,
			wantErr: false,
		},
		{
			name: "Sad Case: invalid gender",
			fields: fields{
				Gender: "invalid",
			},
			want:    enumutils.Gender(""),
			wantErr: true,
		},
		{
			name:    "Sad Case: missing gender",
			fields:  fields{},
			want:    enumutils.Gender(""),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &CMDGenderInput{
				Gender: tt.fields.Gender,
			}
			got, err := g.ParseGender()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenderInput.ParseGender() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenderInput.ParseGender() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPhoneInput_ParsePhone(t *testing.T) {
	type fields struct {
		Phone string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Happy case: valid phone number",
			fields: fields{
				Phone: "0999999999",
			},
			want:    "0999999999",
			wantErr: false,
		},
		{
			name: "Sad case: invalid phone number",
			fields: fields{
				Phone: "invalid",
			},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Sad case: missing phone number",
			fields:  fields{},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &CMDPhoneInput{
				Phone: tt.fields.Phone,
			}
			got, err := p.ParsePhone()
			if (err != nil) != tt.wantErr {
				t.Errorf("PhoneInput.ParsePhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PhoneInput.ParsePhone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSendInviteInput_ParseSendInvite(t *testing.T) {
	type fields struct {
		SendInvite string
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: valid input, yes",
			fields: fields{
				SendInvite: "Yes",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Happy case: valid input, no",
			fields: fields{
				SendInvite: "No",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Sad case: invalid input",
			fields: fields{
				SendInvite: "true",
			},
			want:    false,
			wantErr: true,
		},
		{
			name:    "Sad case: missing input",
			fields:  fields{},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CMDSendInviteInput{
				SendInvite: tt.fields.SendInvite,
			}
			got, err := s.ParseSendInvite()
			if (err != nil) != tt.wantErr {
				t.Errorf("SendInviteInput.ParseSendInvite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SendInviteInput.ParseSendInvite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIDNumberInput_ParseIDNumber(t *testing.T) {
	type fields struct {
		IDNumber string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Happy case: valid IDNumber number",
			fields: fields{
				IDNumber: "0999999999",
			},
			want:    "0999999999",
			wantErr: false,
		},
		{
			name:    "Sad case: missing IDNumber number",
			fields:  fields{},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &CMDIDNumberInput{
				IDNumber: tt.fields.IDNumber,
			}
			got, err := i.ParseIDNumber()
			if (err != nil) != tt.wantErr {
				t.Errorf("IDNumberInput.ParseIDNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IDNumberInput.ParseIDNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStaffInput_ParseStaffNumber(t *testing.T) {
	type fields struct {
		StaffNumber string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Happy case: valid StaffNumber number",
			fields: fields{
				StaffNumber: "0999999999",
			},
			want:    "0999999999",
			wantErr: false,
		},
		{
			name:    "Sad case: missing StaffNumber number",
			fields:  fields{},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CMDStaffInput{
				StaffNumber: tt.fields.StaffNumber,
			}
			got, err := s.ParseStaffNumber()
			if (err != nil) != tt.wantErr {
				t.Errorf("StaffInput.ParseStaffNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StaffInput.ParseStaffNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
