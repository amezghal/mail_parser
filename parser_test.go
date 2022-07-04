package mail_parser

import (
	"fmt"
	"regexp"
	"testing"
)

var tests = []string{
	`email@example.com`,
	`emailexample.com`,
	`emailsdasdadasdasdsaexample.com`,
	`firstname.lastname@example.com`,
	`email@subdomain.example.com`,
	`firstname+lastname@example.com`,
	`email@123.123.123.123`,
	`email@[18.123.123.root:root]`,
	`email@AA`,
	`email@[18.123.123.root`,
	`email@[18.123.123.root------asdasd`,
	`email@[18.123.123.root---a`,
	`email@[18.123.123.root:\root\]`,
	`email@[18.123.123.root:\rootà]`,
	`email@[18.123.123.AAA`,
	`email@[18.123.123]`,
	`email@[255.123.123]`,
	`email@[257.123.123]`,
	`email@[244.123.123]`,
	`email@[24a.123.123]`,
	`email@[14a.123.123]`,
	`email@[1aa.123.123]`,
	`email@[18333.123.123]`,
	`email@[18.123.123.123`,
	`email@example.com`,
	`1234567890@example.com`,
	`email@example-one.com`,
	`_______@example.com`,
	`email@example.name`,
	`email@example.museum`,
	`email-asd--asdas@example---asd---a.co.jp`,
	`firstname---lastname@example.com`,
	`amez.goo{{s.goo.com`,
	`much.”more\ unusual”@example.com`,
	`"helloqwewqw(0000000000@@"@gmail.com`,
	`very.unusual.”-”.unusual.com@example.com`,
	`very.”(),:;<>[]”.VERY.”very@\\ "very”.unusual@strange.example.com`,
	`""@strange.example.com`,
	`"\tes@t"@strange.example.com`,
	`a@a.c-m`,
}

func TestValidator_Validate(t *testing.T) {

	re := regexp.MustCompile(`^(?:[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])$`)

	for _, tt := range tests {
		want := re.MatchString(tt)
		t.Run(fmt.Sprintf("%s_%v", tt, want), func(t *testing.T) {
			instance := New()
			instance.Compile(tt)

			if got := instance.Validate(); got != want {
				t.Errorf("Validate() = %v, want %v", got, want)
			}
		})
	}
}

func BenchmarkValidator_Validate(b *testing.B) {

	for i := 0; i < b.N; i++ {
		for _, dd := range tests {
			instance := New()
			instance.Compile(dd)
			instance.Validate()
		}
	}
}

func BenchmarkValidator_ValidateGO(b *testing.B) {
	re := regexp.MustCompile(`^(?:[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])$`)
	for i := 0; i < b.N; i++ {
		for _, dd := range tests {
			re.MatchString(dd)
		}
	}
}
