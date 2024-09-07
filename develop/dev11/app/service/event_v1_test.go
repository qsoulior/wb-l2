package service

import (
	"context"
	"dev11/app/entity"
	"dev11/app/repo"
	"fmt"
	"reflect"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

func TestNewEventV1(t *testing.T) {
	t.Run("NilRepo", func(t *testing.T) {
		var want Event = nil

		if got := NewEventV1(nil); got != want {
			t.Errorf("NewEventV1() = %v, want %v", got, want)
		}
	})

	t.Run("ValidRepo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repo.NewMockEvent(ctrl)
		want := eventV1{repo: repo}

		if got := NewEventV1(repo); !reflect.DeepEqual(got, want) {
			t.Errorf("NewEventV1() = %v, want %v", got, want)
		}
	})
}

func Test_eventV1_GetForRange(t *testing.T) {
	dateStart, _ := time.Parse(time.DateOnly, "2010-05-20")
	dateEnd, _ := time.Parse(time.DateOnly, "2010-05-25")

	type args struct {
		userID    string
		dateStart time.Time
		dateEnd   time.Time
	}
	tests := []struct {
		name    string
		prepare func(repo *repo.MockEvent)
		args    args
		want    []entity.Event
		wantErr bool
	}{
		{"ValidRange", func(repo *repo.MockEvent) {
			repo.EXPECT().GetForRange(gomock.Any(), gomock.Eq(""), gomock.Eq(dateStart), gomock.Eq(dateEnd)).Return([]entity.Event{}, nil)
		}, args{"", dateStart, dateEnd}, []entity.Event{}, false},
		{"InvalidRange", func(repo *repo.MockEvent) {}, args{"", dateEnd, dateStart}, nil, true},
		{"RepoError", func(repo *repo.MockEvent) {
			repo.EXPECT().GetForRange(gomock.Any(), gomock.Eq(""), gomock.Eq(dateStart), gomock.Eq(dateEnd)).Return([]entity.Event{}, fmt.Errorf(""))
		}, args{"", dateStart, dateEnd}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := repo.NewMockEvent(ctrl)
			tt.prepare(repo)
			e := eventV1{repo: repo}

			got, err := e.GetForRange(ctx, tt.args.userID, tt.args.dateStart, tt.args.dateEnd)
			if (err != nil) != tt.wantErr {
				t.Errorf("eventV1.GetForRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("eventV1.GetForRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_eventV1_GetForDay(t *testing.T) {
	day, _ := time.Parse(time.DateTime, "2010-05-20 16:00:00")
	dateStart, _ := time.Parse(time.DateTime, "2010-05-20 00:00:00")
	dateEnd := dateStart.AddDate(0, 0, 1).Add(-time.Nanosecond)

	type args struct {
		userID string
		day    time.Time
	}
	tests := []struct {
		name    string
		prepare func(repo *repo.MockEvent)
		args    args
		want    []entity.Event
		wantErr bool
	}{
		{"ValidDay", func(repo *repo.MockEvent) {
			repo.EXPECT().GetForRange(gomock.Any(), gomock.Eq(""), gomock.Eq(dateStart), gomock.Eq(dateEnd)).Return([]entity.Event{}, nil)
		}, args{"", day}, []entity.Event{}, false},
		{"RepoError", func(repo *repo.MockEvent) {
			repo.EXPECT().GetForRange(gomock.Any(), gomock.Eq(""), gomock.Eq(dateStart), gomock.Eq(dateEnd)).Return([]entity.Event{}, fmt.Errorf(""))
		}, args{"", day}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := repo.NewMockEvent(ctrl)
			tt.prepare(repo)
			e := eventV1{repo: repo}

			got, err := e.GetForDay(ctx, tt.args.userID, tt.args.day)
			if (err != nil) != tt.wantErr {
				t.Errorf("eventV1.GetForDay() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("eventV1.GetForDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_eventV1_GetForWeek(t *testing.T) {
	week, _ := time.Parse(time.DateOnly, "2010-05-20")
	dateStart, _ := time.Parse(time.DateOnly, "2010-05-17")
	dateEnd := dateStart.Add(7*24*time.Hour - time.Nanosecond)

	type args struct {
		userID string
		week   time.Time
	}
	tests := []struct {
		name    string
		prepare func(repo *repo.MockEvent)
		args    args
		want    []entity.Event
		wantErr bool
	}{
		{"ValidWeek", func(repo *repo.MockEvent) {
			repo.EXPECT().GetForRange(gomock.Any(), gomock.Eq(""), gomock.Eq(dateStart), gomock.Eq(dateEnd)).Return([]entity.Event{}, nil)
		}, args{"", week}, []entity.Event{}, false},
		{"RepoError", func(repo *repo.MockEvent) {
			repo.EXPECT().GetForRange(gomock.Any(), gomock.Eq(""), gomock.Eq(dateStart), gomock.Eq(dateEnd)).Return([]entity.Event{}, fmt.Errorf(""))
		}, args{"", week}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := repo.NewMockEvent(ctrl)
			tt.prepare(repo)
			e := eventV1{repo: repo}

			got, err := e.GetForWeek(ctx, tt.args.userID, tt.args.week)
			if (err != nil) != tt.wantErr {
				t.Errorf("eventV1.GetForWeek() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("eventV1.GetForWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_eventV1_GetForMonth(t *testing.T) {
	month, _ := time.Parse(time.DateOnly, "2010-05-20")
	dateStart, _ := time.Parse(time.DateOnly, "2010-05-01")
	dateEnd := dateStart.AddDate(0, 1, 0).Add(-time.Nanosecond)

	type args struct {
		userID string
		month  time.Time
	}
	tests := []struct {
		name    string
		prepare func(repo *repo.MockEvent)
		args    args
		want    []entity.Event
		wantErr bool
	}{
		{"ValidDay", func(repo *repo.MockEvent) {
			repo.EXPECT().GetForRange(gomock.Any(), gomock.Eq(""), gomock.Eq(dateStart), gomock.Eq(dateEnd)).Return([]entity.Event{}, nil)
		}, args{"", month}, []entity.Event{}, false},
		{"RepoError", func(repo *repo.MockEvent) {
			repo.EXPECT().GetForRange(gomock.Any(), gomock.Eq(""), gomock.Eq(dateStart), gomock.Eq(dateEnd)).Return([]entity.Event{}, fmt.Errorf(""))
		}, args{"", month}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := repo.NewMockEvent(ctrl)
			tt.prepare(repo)
			e := eventV1{repo: repo}

			got, err := e.GetForMonth(ctx, tt.args.userID, tt.args.month)
			if (err != nil) != tt.wantErr {
				t.Errorf("eventV1.GetForMonth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("eventV1.GetForMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_eventV1_Create(t *testing.T) {
	validUUID := "18310e71-4df6-42c0-adf4-1a280013dd08"
	validEvent := entity.Event{Title: "event", UserID: validUUID}

	type args struct {
		event entity.Event
	}
	tests := []struct {
		name    string
		prepare func(repo *repo.MockEvent)
		args    args
		want    entity.Event
		wantErr bool
	}{
		{"ValidEvent", func(repo *repo.MockEvent) {
			repo.EXPECT().Create(gomock.Any(), gomock.Eq(validEvent)).Return(validEvent, nil)
		}, args{validEvent}, validEvent, false},
		{"InvalidEvent", func(repo *repo.MockEvent) {}, args{entity.EmptyEvent}, entity.EmptyEvent, true},
		{"RepoError", func(repo *repo.MockEvent) {
			repo.EXPECT().Create(gomock.Any(), gomock.Eq(validEvent)).Return(entity.EmptyEvent, fmt.Errorf(""))
		}, args{validEvent}, entity.EmptyEvent, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := repo.NewMockEvent(ctrl)
			tt.prepare(repo)
			e := eventV1{repo: repo}

			got, err := e.Create(ctx, tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("eventV1.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("eventV1.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_eventV1_Update(t *testing.T) {
	validUUID := "18310e71-4df6-42c0-adf4-1a280013dd08"
	validEvent := entity.Event{ID: validUUID, Title: "event", UserID: validUUID}

	type args struct {
		event entity.Event
	}
	tests := []struct {
		name    string
		prepare func(repo *repo.MockEvent)
		args    args
		want    entity.Event
		wantErr bool
	}{
		{"ValidEvent", func(repo *repo.MockEvent) {
			repo.EXPECT().Update(gomock.Any(), gomock.Eq(validEvent)).Return(validEvent, nil)
		}, args{validEvent}, validEvent, false},
		{"InvalidEvent", func(repo *repo.MockEvent) {}, args{entity.EmptyEvent}, entity.EmptyEvent, true},
		{"RepoError", func(repo *repo.MockEvent) {
			repo.EXPECT().Update(gomock.Any(), gomock.Eq(validEvent)).Return(entity.EmptyEvent, fmt.Errorf(""))
		}, args{validEvent}, entity.EmptyEvent, true},
		{"RepoErrorNotExist", func(r *repo.MockEvent) {
			r.EXPECT().Update(gomock.Any(), gomock.Eq(validEvent)).Return(entity.EmptyEvent, repo.ErrNotExist)
		}, args{validEvent}, entity.EmptyEvent, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := repo.NewMockEvent(ctrl)
			tt.prepare(repo)
			e := eventV1{repo: repo}

			got, err := e.Update(ctx, tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("eventV1.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("eventV1.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_eventV1_Delete(t *testing.T) {
	type args struct {
		userID string
		id     string
	}
	tests := []struct {
		name    string
		prepare func(repo *repo.MockEvent)
		args    args
		wantErr bool
	}{
		{"ValidEvent", func(repo *repo.MockEvent) {
			repo.EXPECT().Delete(gomock.Any(), gomock.Eq("1"), gomock.Eq("2")).Return(nil)
		}, args{"1", "2"}, false},
		{"RepoError", func(repo *repo.MockEvent) {
			repo.EXPECT().Delete(gomock.Any(), gomock.Eq("1"), gomock.Eq("2")).Return(fmt.Errorf(""))
		}, args{"1", "2"}, true},
		{"RepoErrorNotExist", func(r *repo.MockEvent) {
			r.EXPECT().Delete(gomock.Any(), gomock.Eq(""), gomock.Eq("")).Return(repo.ErrNotExist)
		}, args{"", ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := repo.NewMockEvent(ctrl)
			tt.prepare(repo)
			e := eventV1{repo: repo}

			if err := e.Delete(ctx, tt.args.userID, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("eventV1.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
