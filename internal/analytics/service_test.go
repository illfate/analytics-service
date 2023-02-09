package analytics_test

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

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
		ClientTime: time.Now(),
	}
	t.Run("happy path", func(t *testing.T) {
		repo := mock.NewMockRepository(ctrl)

		repo.EXPECT().InsertEvents(gomock.Any(), []analytics.Event{event}).Return(nil)

		service := analytics.NewService(repo)
		service.StartCreatingEvents(context.TODO(), func(err error) {
			require.NoError(t, err)
		})
		err := service.CreateEvents(context.Background(), []analytics.Event{event})
		assert.NoError(t, err)
		service.Close()
	})

	t.Run("error", func(t *testing.T) {
		repo := mock.NewMockRepository(ctrl)

		repo.EXPECT().InsertEvents(gomock.Any(), []analytics.Event{event}).Return(io.EOF)

		service := analytics.NewService(repo)
		service.StartCreatingEvents(context.TODO(), func(err error) {
			require.NotNil(t, err)
			assert.Equal(t, fmt.Errorf("failed to insert events: %w", io.EOF), err)
		})
		err := service.CreateEvents(context.Background(), []analytics.Event{event})
		require.NoError(t, err)
		service.Close()
	})
}
