package rewards

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"

	"github.com/brave-intl/bat-go/services/rewards/model"
)

func TestService_GetCards(t *testing.T) {
	type tcGiven struct {
		cfg   *Config
		s3Svc s3Service
	}

	type tcExpected struct {
		cards CardBytes
		err   error
	}

	type testCase struct {
		name  string
		given tcGiven
		exp   tcExpected
	}

	tests := []testCase{
		{
			name: "error_get_object",
			given: tcGiven{
				cfg: &Config{
					Cards: &CardsConfig{},
				},
				s3Svc: &mockS3Service{
					fnGetObject: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
						return nil, model.Error("error")
					},
				},
			},
			exp: tcExpected{
				err: model.Error("error"),
			},
		},

		{
			name: "success",
			given: tcGiven{
				cfg: &Config{
					Cards: &CardsConfig{},
				},
				s3Svc: &mockS3Service{
					fnGetObject: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
						cards := CardBytes(`{ "card": [{"title": "<string>", "description": "<string>", "url": "<string>", "thumbnail": "<string>"}] }`)

						out := &s3.GetObjectOutput{
							Body: io.NopCloser(bytes.NewReader(cards)),
						}

						return out, nil
					},
				},
			},
			exp: tcExpected{
				cards: CardBytes(`{ "card": [{"title": "<string>", "description": "<string>", "url": "<string>", "thumbnail": "<string>"}] }`),
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			s := &Service{
				cfg:   tc.given.cfg,
				s3Svc: tc.given.s3Svc,
			}

			ctx := context.Background()

			actual, err := s.GetCardsAsBytes(ctx)

			assert.ErrorIs(t, err, tc.exp.err)
			assert.Equal(t, tc.exp.cards, actual)
		})
	}
}

type mockS3Service struct {
	fnGetObject func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

func (m *mockS3Service) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	if m.fnGetObject == nil {
		return &s3.GetObjectOutput{}, nil
	}

	return m.fnGetObject(ctx, params, optFns...)
}