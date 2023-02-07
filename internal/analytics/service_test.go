package analytics_test

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/illfate/analytics-service/internal/analytics"
	"github.com/illfate/analytics-service/internal/mock"
)

func TestServiceCreateEvents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	event := analytics.Event{
		ClientTime: "time",
	}
	t.Run("happy path", func(t *testing.T) {
		repo := mock.NewMockRepository(ctrl)

		repo.EXPECT().InsertEvents(gomock.Any(), event).Return(nil)

		service := analytics.NewService(repo)
		err := service.CreateEvents(context.Background(), event)
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		repo := mock.NewMockRepository(ctrl)

		repo.EXPECT().InsertEvents(gomock.Any(), event).Return(io.EOF)

		service := analytics.NewService(repo)
		err := service.CreateEvents(context.Background(), event)
		require.NotNil(t, err)
		assert.Equal(t, fmt.Errorf("failed to insert events: %w", io.EOF), err)
	})
}
