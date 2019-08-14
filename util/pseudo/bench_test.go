package pseudo

import (
	"testing"

	"github.com/corestoreio/pkg/storage/null"
)

func BenchmarkFakerDataNOTTagged(b *testing.B) {
	s := MustNewService(0, nil)
	for i := 0; i < b.N; i++ {
		a := NotTaggedStruct{}
		err := s.FakeData(&a)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFakerDataTagged(b *testing.B) {
	s := MustNewService(0, nil)
	for i := 0; i < b.N; i++ {
		a := TaggedStruct{}
		err := s.FakeData(&a)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCustomerEntity(b *testing.B) {
	s := MustNewService(0, nil)

	b.Run("tagged", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a := CustomerEntityTagged{}
			err := s.FakeData(&a)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("untagged", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a := CustomerEntityUnTagged{}
			err := s.FakeData(&a)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

type CustomerEntityTagged struct {
	EntityID     uint32      `faker:"id" max_len:"10"`
	WebsiteID    null.Uint32 `max_len:"5"`
	Email        null.String `faker:"email" max_len:"255"`
	GroupID      uint32      `max_len:"5"`
	Prefix       null.String `faker:"prefix" max_len:"40"`
	Firstname    null.String `faker:"first_name" max_len:"255"`
	Middlename   null.String `faker:"username" max_len:"255"`
	Lastname     null.String `max_len:"255"`
	Dob          null.Time
	PasswordHash null.String `max_len:"128"`
}

type CustomerEntityUnTagged struct {
	EntityID     uint32
	WebsiteID    null.Uint32
	Email        null.String
	GroupID      uint32
	Prefix       null.String
	Firstname    null.String
	Middlename   null.String
	Lastname     null.String
	Dob          null.Time
	PasswordHash null.String
}
