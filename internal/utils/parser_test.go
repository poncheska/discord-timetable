package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTimetable(t *testing.T) {
	tt, err := ParseTimetable("https://timetable.spbu.ru/GSOM/StudentGroupEvents/Primary/276304/2021-05-31")
	fmt.Println(tt.GetString())
	assert.Nil(t, err)
}