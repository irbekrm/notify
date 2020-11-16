package watch

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/irbekrm/notify/internal/github"
	"github.com/irbekrm/notify/internal/receiver"
	"github.com/irbekrm/notify/internal/store"
	"github.com/irbekrm/notify/mocks"
)

func TestClient_PollRepoFunc(t *testing.T) {
	tests := []struct {
		name             string
		repoName         string
		issueDescription string
		startTime        StartTime
		interval         time.Duration
		setup            func(ctx context.Context, ctrl *gomock.Controller, startTime StartTime, repoName string, issueDescription string) (github.Finder, receiver.Notifier, store.DB)
	}{
		{
			name:      "no issues found",
			startTime: StartTime{t: time.Now()},
			setup: func(ctx context.Context, ctrl *gomock.Controller, startTime StartTime, repoName string, issueDescription string) (github.Finder, receiver.Notifier, store.DB) {
				rp := mocks.NewMockFinder(ctrl)
				rp.
					EXPECT().
					RepoName().
					Return(repoName)
				rp.
					EXPECT().
					Find(ctx, startTime.t).
					Return([]github.Issue{}, nil)
				rec := mocks.NewMockNotifier(ctrl)
				db := mocks.NewMockDB(ctrl)
				return rp, rec, db
			},
		},
		{
			name:      "a seen issue found",
			startTime: StartTime{t: time.Now()},
			setup: func(ctx context.Context, ctrl *gomock.Controller, startTime StartTime, repoName string, issueDescription string) (github.Finder, receiver.Notifier, store.DB) {
				issue := mocks.NewMockIssue(ctrl)
				rp := mocks.NewMockFinder(ctrl)
				rp.
					EXPECT().
					RepoName().
					Return(repoName)
				rp.
					EXPECT().
					Find(ctx, startTime.t).
					Return([]github.Issue{issue}, nil)
				rec := mocks.NewMockNotifier(ctrl)
				db := mocks.NewMockDB(ctrl)
				db.
					EXPECT().
					FindIssue(ctx, issue, repoName).
					Return(true, nil)
				return rp, rec, db
			},
		},
		{
			name:      "an unseen issue found",
			startTime: StartTime{t: time.Now()},
			setup: func(ctx context.Context, ctrl *gomock.Controller, startTime StartTime, repoName string, issueDescription string) (github.Finder, receiver.Notifier, store.DB) {
				issue := mocks.NewMockIssue(ctrl)
				issue.
					EXPECT().
					Description().
					Return(issueDescription)
				rp := mocks.NewMockFinder(ctrl)
				rp.
					EXPECT().
					RepoName().
					Return(repoName)
				rp.
					EXPECT().
					Find(ctx, startTime.t).
					Return([]github.Issue{issue}, nil)
				rec := mocks.NewMockNotifier(ctrl)
				rec.
					EXPECT().
					Notify(issueNotification(issueDescription))

				db := mocks.NewMockDB(ctrl)
				db.
					EXPECT().
					FindIssue(ctx, issue, repoName).
					Return(false, nil)
				db.
					EXPECT().
					WriteIssue(ctx, issue, repoName).
					Return(nil)
				return rp, rec, db
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ctx := context.TODO()
			rp, rec, db := tt.setup(ctx, ctrl, tt.startTime, tt.repoName, tt.issueDescription)
			c := &Client{
				startTime: tt.startTime,
				interval:  tt.interval,
				rp:        rp,
				reciever:  rec,
				db:        db,
			}
			f := c.PollRepoFunc(ctx)
			f()
		})
	}
}
