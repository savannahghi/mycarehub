package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
)

func TestIfCoverExistsInSlice(t *testing.T) {
	src := []base.Cover{
		{
			IdentifierHash: base.CreateCoverHash(base.Cover{
				PayerName:      "payer1",
				PayerSladeCode: 1,
				MemberNumber:   "mem1",
				MemberName:     "name1",
			}),
			PayerName:      "payer1",
			PayerSladeCode: 1,
			MemberNumber:   "mem1",
			MemberName:     "name1",
		},
		{
			IdentifierHash: base.CreateCoverHash(base.Cover{
				PayerName:      "payer2",
				PayerSladeCode: 2,
				MemberNumber:   "mem2",
				MemberName:     "name2",
			}),
			PayerName:      "payer2",
			PayerSladeCode: 2,
			MemberNumber:   "mem2",
			MemberName:     "name2",
		},
	}

	tests := []struct {
		name      string
		args      base.Cover
		srcCovers []base.Cover
		want      bool
	}{
		{
			name: "valid: exists_1",
			args: base.Cover{
				IdentifierHash: base.CreateCoverHash(base.Cover{
					PayerName:      "payer1",
					PayerSladeCode: 1,
					MemberNumber:   "mem1",
					MemberName:     "name1",
				}),
				PayerName:      "payer1",
				PayerSladeCode: 1,
				MemberNumber:   "mem1",
				MemberName:     "name1",
			},
			srcCovers: src,
			want:      true,
		},

		{
			name: "valid: exists_2",
			args: base.Cover{
				IdentifierHash: base.CreateCoverHash(base.Cover{
					PayerName:      "payer2",
					PayerSladeCode: 2,
					MemberNumber:   "mem2",
					MemberName:     "name2",
				}),
				PayerName:      "payer2",
				PayerSladeCode: 2,
				MemberNumber:   "mem2",
				MemberName:     "name2",
			},
			srcCovers: src,
			want:      true,
		},

		{
			name: "invalid: does not exist_1",
			args: base.Cover{
				IdentifierHash: base.CreateCoverHash(base.Cover{
					PayerName:      "payer1",
					PayerSladeCode: 1,
					MemberNumber:   "mem11",
					MemberName:     "name11",
				}),
				PayerName:      "payer1",
				PayerSladeCode: 1,
				MemberNumber:   "mem11",
				MemberName:     "name11",
			},
			srcCovers: src,
			want:      false,
		},

		{
			name: "invalid: does not exist_2",
			args: base.Cover{
				IdentifierHash: base.CreateCoverHash(base.Cover{
					PayerName:      "payer1",
					PayerSladeCode: 1,
					MemberNumber:   "mem1",
					MemberName:     "name11",
				}),
				PayerName:      "payer1",
				PayerSladeCode: 1,
				MemberNumber:   "mem1",
				MemberName:     "name11",
			},
			srcCovers: src,
			want:      false,
		},

		{
			name: "invalid: does not exist_3",
			args: base.Cover{
				IdentifierHash: base.CreateCoverHash(base.Cover{
					PayerName:      "payer2",
					PayerSladeCode: 2,
					MemberNumber:   "mem22",
					MemberName:     "name2",
				}),
				PayerName:      "payer2",
				PayerSladeCode: 2,
				MemberNumber:   "mem22",
				MemberName:     "name2",
			},
			srcCovers: src,
			want:      false,
		},

		{
			name: "invalid: does not exist_4",
			args: base.Cover{
				IdentifierHash: base.CreateCoverHash(base.Cover{
					PayerName:      "payer3",
					PayerSladeCode: 2,
					MemberNumber:   "mem2",
					MemberName:     "name2",
				}),
				PayerName:      "payer3",
				PayerSladeCode: 2,
				MemberNumber:   "mem2",
				MemberName:     "name2",
			},
			srcCovers: src,
			want:      false,
		},

		{
			name: "invalid: does not exist_5",
			args: base.Cover{
				IdentifierHash: base.CreateCoverHash(base.Cover{
					PayerName:      "payer3",
					PayerSladeCode: 2,
					MemberNumber:   "mem2",
					MemberName:     "name2",
				}),
				PayerName:      "payer3",
				PayerSladeCode: 2,
				MemberNumber:   "mem2",
				MemberName:     "name2",
			},
			srcCovers: []base.Cover{},
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := utils.IfCoverExistsInSlice(tt.srcCovers, tt.args)
			assert.Equal(t, tt.want, resp)
		})
	}
}

func TestAddHashToCovers(t *testing.T) {
	unhashedCovers := []base.Cover{
		{

			PayerName:      "payer1",
			PayerSladeCode: 1,
			MemberNumber:   "mem1",
			MemberName:     "name1",
		},
		{

			PayerName:      "payer2",
			PayerSladeCode: 2,
			MemberNumber:   "mem2",
			MemberName:     "name2",
		},
	}

	// check the covers don't have hash identifiers yet'
	for _, cover := range unhashedCovers {
		assert.Nil(t, cover.IdentifierHash)
	}

	// now hash the covers. This shold pass and return the covers
	hashedCovers1 := utils.AddHashToCovers(unhashedCovers)
	assert.Equal(t, len(unhashedCovers), len(hashedCovers1))

	hashedCovers2 := []base.Cover{
		{
			IdentifierHash: base.CreateCoverHash(base.Cover{
				PayerName:      "payer1",
				PayerSladeCode: 1,
				MemberNumber:   "mem1",
				MemberName:     "name1",
			}),
			PayerName:      "payer1",
			PayerSladeCode: 1,
			MemberNumber:   "mem1",
			MemberName:     "name1",
		},
		{
			IdentifierHash: base.CreateCoverHash(base.Cover{
				PayerName:      "payer2",
				PayerSladeCode: 2,
				MemberNumber:   "mem2",
				MemberName:     "name2",
			}),
			PayerName:      "payer2",
			PayerSladeCode: 2,
			MemberNumber:   "mem2",
			MemberName:     "name2",
		},
	}

	// check the covers do have hash identifiers yet'
	for _, cover := range hashedCovers2 {
		assert.NotNil(t, cover.IdentifierHash)
	}

	// now hash the covers. This shold fail and return an empty slice
	hashedCovers3 := utils.AddHashToCovers(hashedCovers2)
	assert.Equal(t, 0, len(hashedCovers3))
}

func TestMatchAndReturn(t *testing.T) {
	tests := []struct {
		old  bool
		new  bool
		want bool
	}{
		{old: false, new: true, want: true},
		{old: true, new: false, want: false},
		{old: true, new: true, want: true},
		{old: false, new: false, want: false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			resp := utils.MatchAndReturn(tt.old, tt.new)
			assert.Equal(t, tt.want, resp)
		})
	}
}
