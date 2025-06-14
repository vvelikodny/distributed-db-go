package parser

import (
	"context"
	"testing"

	"concurrency-go/internal"

	"go.uber.org/zap"
)

func TestParser_ProcessCommand(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	tests := []struct {
		name    string
		query   string
		want    *internal.Command
		wantErr error
	}{
		{
			name:    "empty query",
			query:   "",
			want:    nil,
			wantErr: internal.ErrInvalidCommand,
		},
		{
			name:    "invalid command",
			query:   "INVALID",
			want:    nil,
			wantErr: internal.ErrInvalidCommand,
		},
		{
			name:    "SET with missing value",
			query:   "SET key",
			want:    nil,
			wantErr: internal.ErrInvalidSyntax,
		},
		{
			name:    "GET with extra value",
			query:   "GET key value",
			want:    nil,
			wantErr: internal.ErrInvalidSyntax,
		},
		{
			name:    "DEL with extra value",
			query:   "DEL key value",
			want:    nil,
			wantErr: internal.ErrInvalidSyntax,
		},
		{
			name:  "valid SET command",
			query: "SET key value",
			want: &internal.Command{
				Type:  "SET",
				Key:   "key",
				Value: "value",
			},
			wantErr: nil,
		},
		{
			name:  "valid GET command",
			query: "GET key",
			want: &internal.Command{
				Type: "GET",
				Key:  "key",
			},
			wantErr: nil,
		},
		{
			name:  "valid DEL command",
			query: "DEL key",
			want: &internal.Command{
				Type: "DEL",
				Key:  "key",
			},
			wantErr: nil,
		},
		{
			name:    "case sensitive command - lowercase",
			query:   "set key value",
			want:    nil,
			wantErr: internal.ErrInvalidCommand,
		},
		{
			name:    "case sensitive command - mixed case",
			query:   "SeT key value",
			want:    nil,
			wantErr: internal.ErrInvalidCommand,
		},
		{
			name:    "invalid key characters",
			query:   "SET key@ value",
			want:    nil,
			wantErr: internal.ErrInvalidSyntax,
		},
		{
			name:    "invalid value characters",
			query:   "SET key value@",
			want:    nil,
			wantErr: internal.ErrInvalidSyntax,
		},
		{
			name:  "valid special characters in key",
			query: "SET key_123 value",
			want: &internal.Command{
				Type:  "SET",
				Key:   "key_123",
				Value: "value",
			},
			wantErr: nil,
		},
		{
			name:  "valid special characters in value",
			query: "SET key value-123",
			want: &internal.Command{
				Type:  "SET",
				Key:   "key",
				Value: "value-123",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt // Create new variable for parallel tests
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Enable parallel execution
			p := New(logger)
			got, err := p.ProcessCommand(context.Background(), tt.query)

			if err != tt.wantErr {
				t.Errorf("ProcessCommand() error = %v, want %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if got.Type != tt.want.Type {
					t.Errorf("ProcessCommand() Type = %v, want %v", got.Type, tt.want.Type)
				}
				if got.Key != tt.want.Key {
					t.Errorf("ProcessCommand() Key = %v, want %v", got.Key, tt.want.Key)
				}
				if got.Value != tt.want.Value {
					t.Errorf("ProcessCommand() Value = %v, want %v", got.Value, tt.want.Value)
				}
			}
		})
	}
}
