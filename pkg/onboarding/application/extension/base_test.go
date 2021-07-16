package extension_test

// func TestNewBaseExtensionImpl(t *testing.T) {
// 	type args struct {
// 		fc firebasetools.IFirebaseClient
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want extension.BaseExtension
// 	}{
// 		{
// 			name: "Happy test",
// 			args: args{
// 				fc: &firebasetools.MockFirebaseClient{},
// 			},
// 			want: &extension.BaseExtensionImpl{
// 				fc:firebasetools.IFirebaseClient{
// 					firebasetools.InitFirebase()
// 				},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := extension.NewBaseExtensionImpl(tt.args.fc); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NewBaseExtensionImpl() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
