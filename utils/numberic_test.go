package utils

import (
	"math"
	"testing"
)

func TestCheckFloat32Special(t *testing.T) {
	tests := []struct {
		value float32
		want  bool
	}{
		{float32(math.NaN()), true},
		{float32(math.Inf(1)), true},
		{float32(math.Inf(-1)), true},
		{1.0, false},
		{-1.0, false},
	}

	for _, tcase := range tests {
		got := CheckFloat32Special(tcase.value)
		if got != tcase.want {
			t.Errorf("CheckFloat32Special(%v) = %v; want %v", tcase.value, got, tcase.want)
		}
	}
}

func TestCheckFloat64Special(t *testing.T) {
	tests := []struct {
		value float64
		want  bool
	}{
		{math.NaN(), true},
		{math.Inf(1), true},
		{math.Inf(-1), true},
		{1.0, false},
		{-1.0, false},
	}

	for _, tcase := range tests {
		got := CheckFloat64Special(tcase.value)
		if got != tcase.want {
			t.Errorf("CheckFloat64Special(%v) = %v; want %v", tcase.value, got, tcase.want)
		}
	}
}

func TestInt16CheckOver(t *testing.T) {
	tests := []struct {
		value, delta int16
		want         bool
	}{
		{math.MaxInt16, 1, true},
		{math.MinInt16, -1, true},
		{100, 50, false},
		{-100, -50, false},
	}

	for _, tcase := range tests {
		got := Int16CheckOver(tcase.value, tcase.delta)
		if got != tcase.want {
			t.Errorf("Int16CheckOver(%v, %v) = %v; want %v", tcase.value, tcase.delta, got, tcase.want)
		}
	}
}

func TestInt32CheckOver(t *testing.T) {
	tests := []struct {
		value, delta int32
		want         bool
	}{
		{math.MaxInt32, 1, true},
		{math.MinInt32, -1, true},
		{100, 50, false},
		{-100, -50, false},
	}

	for _, tcase := range tests {
		got := Int32CheckOver(tcase.value, tcase.delta)
		if got != tcase.want {
			t.Errorf("Int32CheckOver(%v, %v) = %v; want %v", tcase.value, tcase.delta, got, tcase.want)
		}
	}
}

func TestInt64CheckOver(t *testing.T) {
	tests := []struct {
		value, delta int64
		want         bool
	}{
		{math.MaxInt64, 1, true},
		{math.MinInt64, -1, true},
		{100, 50, false},
		{-100, -50, false},
	}

	for _, tcase := range tests {
		got := Int64CheckOver(tcase.value, tcase.delta)
		if got != tcase.want {
			t.Errorf("Int64CheckOver(%v, %v) = %v; want %v", tcase.value, tcase.delta, got, tcase.want)
		}
	}
}

func TestUIntCheckUnderFlow(t *testing.T) {
	got := UintCheckUnderFlow[uint16](0, 1)
	if !got {
		t.Errorf("UintCheckUnderFlow(%v, %v) = %v; want %v", 0, 1, got, true)
	}
	got = UintCheckUnderFlow[uint64](100, 1)
	if got {
		t.Errorf("UintCheckUnderFlow(%v, %v) = %v; want %v", 100, 1, got, false)
	}
}

func TestUInt16CheckOverFlow(t *testing.T) {
	tests := []struct {
		value, delta uint16
		want         bool
	}{
		{math.MaxUint16, 1, true},
		{100, 50, false},
	}

	for _, tcase := range tests {
		got := UInt16CheckOverFlow(tcase.value, tcase.delta)
		if got != tcase.want {
			t.Errorf("UInt16CheckOverFlow(%v, %v) = %v; want %v", tcase.value, tcase.delta, got, tcase.want)
		}
	}
}

func TestUInt32CheckOverFlow(t *testing.T) {
	tests := []struct {
		value, delta uint32
		want         bool
	}{
		{math.MaxUint32, 1, true},
		{100, 50, false},
	}

	for _, tcase := range tests {
		got := UInt32CheckOverFlow(tcase.value, tcase.delta)
		if got != tcase.want {
			t.Errorf("UInt32CheckOverFlow(%v, %v) = %v; want %v", tcase.value, tcase.delta, got, tcase.want)
		}
	}
}

func TestUInt64CheckOverFlow(t *testing.T) {
	tests := []struct {
		value, delta uint64
		want         bool
	}{
		{math.MaxUint64, 1, true},
		{100, 50, false},
	}

	for _, tcase := range tests {
		got := UInt64CheckOverFlow(tcase.value, tcase.delta)
		if got != tcase.want {
			t.Errorf("UInt64CheckOverFlow(%v, %v) = %v; want %v", tcase.value, tcase.delta, got, tcase.want)
		}
	}
}

func TestFloat32CheckOver(t *testing.T) {
	tests := []struct {
		value, delta float32
		want         bool
	}{
		{math.MaxFloat32, 1, true},
		{-math.MaxFloat32, -1, true},
		{100.1, 50.5, false},
		{-1003, -50.8, false},
	}

	for _, tcase := range tests {
		got := Float32CheckOver(tcase.value, tcase.delta)
		if got != tcase.want {
			t.Errorf("Float32CheckOver(%v, %v) = %v; want %v", tcase.value, tcase.delta, got, tcase.want)
		}
	}
}

func TestFloat64CheckOver(t *testing.T) {
	tests := []struct {
		value, delta float64
		want         bool
	}{
		{math.MaxFloat64, 1, true},
		{-math.MaxFloat64, -1, true},
		{100.1, 50.5, false},
		{-1003, -50.8, false},
	}

	for _, tcase := range tests {
		got := Float64CheckOver(tcase.value, tcase.delta)
		if got != tcase.want {
			t.Errorf("Float64CheckOver(%v, %v) = %v; want %v", tcase.value, tcase.delta, got, tcase.want)
		}
	}
}

func TestBinary2Int16(t *testing.T) {
	tests := []struct {
		data    []byte
		errwant bool
	}{
		{[]byte{0xFF}, true},
		{[]byte{0xFF, 0xFF, 0xFF}, true},
		{[]byte{0xFF, 0xFF}, false},
	}
	for _, tcase := range tests {
		got, err := Binary2Int16(tcase.data)
		if tcase.errwant {
			if err == nil {
				t.Fatalf("expected error for bytes length")
			}
		} else {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tcase.errwant && got != -1 {
				t.Errorf("Binary2Int16(%v) = %v; want %v", tcase.data, got, -1)
			}
		}
	}
}

func TestInt16toBinary(t *testing.T) {
	wants := []int16{-100, 0, 57, math.MaxInt16, math.MinInt16}

	for i, want := range wants {
		binary := Int16toBinary(want)

		got, err := Binary2Int16(binary)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
			continue
		}

		if got != want {
			t.Errorf("Test case %d: conversion failed\n  original: %d\n  binary: %v\n  restored: %d",
				i, want, binary, got)
		}
	}
}

func TestBinary2Int32(t *testing.T) {
	tests := []struct {
		data    []byte
		errwant bool
	}{
		{[]byte{0xFF, 0xFF, 0xFF}, true},
		{[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, true},
		{[]byte{0xFE, 0xFF, 0xFF, 0xFF}, false},
	}
	for _, tcase := range tests {
		got, err := Binary2Int32(tcase.data)
		if tcase.errwant {
			if err == nil {
				t.Fatalf("expected error for bytes length")
			}
		} else {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tcase.errwant && got != -2 {
				t.Errorf("Binary2Int32(%v) = %v; want %v", tcase.data, got, -2)
			}
		}
	}
}

func TestInt32toBinary(t *testing.T) {
	wants := []int32{-197, 0, 157, math.MaxInt32, math.MinInt32}

	for i, want := range wants {
		binary := Int32toBinary(want)

		got, err := Binary2Int32(binary)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
			continue
		}

		if got != want {
			t.Errorf("Test case %d: conversion failed\n  original: %d\n  binary: %v\n  restored: %d",
				i, want, binary, got)
		}
	}
}

func TestBinary2Int64(t *testing.T) {
	tests := []struct {
		data    []byte
		errwant bool
	}{
		{[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, true},
		{[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, true},
		{[]byte{0xFD, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, false},
	}
	for _, tcase := range tests {
		got, err := Binary2Int64(tcase.data)
		if tcase.errwant {
			if err == nil {
				t.Fatalf("expected error for bytes length")
			}
		} else {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tcase.errwant && got != -3 {
				t.Errorf("Binary2Int64(%v) = %v; want %v", tcase.data, got, -3)
			}
		}
	}
}

func TestInt64toBinary(t *testing.T) {
	wants := []int64{-100, 0, 989898765, math.MaxInt64, math.MinInt64}

	for i, want := range wants {
		binary := Int64toBinary(want)

		got, err := Binary2Int64(binary)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
			continue
		}

		if got != want {
			t.Errorf("Test case %d: conversion failed\n  original: %d\n  binary: %v\n  restored: %d",
				i, want, binary, got)
		}
	}
}

func TestBinary2UInt16(t *testing.T) {
	tests := []struct {
		data    []byte
		errwant bool
	}{
		{[]byte{0x01}, true},
		{[]byte{0x01, 0x00, 0x00}, true},
		{[]byte{0x01, 0x00}, false},
	}
	for _, tcase := range tests {
		got, err := Binary2UInt16(tcase.data)
		if tcase.errwant {
			if err == nil {
				t.Fatalf("expected error for bytes length")
			}
		} else {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tcase.errwant && got != 1 {
				t.Errorf("Binary2UInt16(%v) = %v; want %v", tcase.data, got, 1)
			}
		}
	}
}

func TestUInt16toBinary(t *testing.T) {
	wants := []uint16{0, 57, math.MaxUint16}

	for i, want := range wants {
		binary := UInt16toBinary(want)

		got, err := Binary2UInt16(binary)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
			continue
		}

		if got != want {
			t.Errorf("Test case %d: conversion failed\n  original: %d\n  binary: %v\n  restored: %d",
				i, want, binary, got)
		}
	}
}

func TestBinary2UInt32(t *testing.T) {
	tests := []struct {
		data    []byte
		errwant bool
	}{
		{[]byte{0x02, 0x00, 0x00}, true},
		{[]byte{0x02, 0x00, 0x00, 0x00, 0x00}, true},
		{[]byte{0x02, 0x00, 0x00, 0x00}, false},
	}
	for _, tcase := range tests {
		got, err := Binary2UInt32(tcase.data)
		if tcase.errwant {
			if err == nil {
				t.Fatalf("expected error for bytes length")
			}
		} else {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tcase.errwant && got != 2 {
				t.Errorf("Binary2UInt32(%v) = %v; want %v", tcase.data, got, 2)
			}
		}
	}
}

func TestUInt32toBinary(t *testing.T) {
	wants := []uint32{0, 5007, math.MaxUint16, math.MaxUint32}

	for i, want := range wants {
		binary := UInt32toBinary(want)

		got, err := Binary2UInt32(binary)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
			continue
		}

		if got != want {
			t.Errorf("Test case %d: conversion failed\n  original: %d\n  binary: %v\n  restored: %d",
				i, want, binary, got)
		}
	}
}

func TestBinary2UInt64(t *testing.T) {
	tests := []struct {
		data    []byte
		errwant bool
	}{
		{[]byte{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, true},
		{[]byte{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, true},
		{[]byte{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, false},
	}
	for _, tcase := range tests {
		got, err := Binary2UInt64(tcase.data)
		if tcase.errwant {
			if err == nil {
				t.Fatalf("expected error for bytes length")
			}
		} else {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tcase.errwant && got != 3 {
				t.Errorf("Binary2UInt64(%v) = %v; want %v", tcase.data, got, 3)
			}
		}
	}
}

func TestUInt64toBinary(t *testing.T) {
	wants := []uint64{0, 25237, math.MaxUint16, math.MaxUint32, math.MaxUint64}

	for i, want := range wants {
		binary := UInt64toBinary(want)

		got, err := Binary2UInt64(binary)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
			continue
		}

		if got != want {
			t.Errorf("Test case %d: conversion failed\n  original: %d\n  binary: %v\n  restored: %d",
				i, want, binary, got)
		}
	}
}

func TestBinary2Float32(t *testing.T) {
	tests := []struct {
		data    []byte
		errwant bool
	}{
		{[]byte{0xC3, 0xF5, 0x48}, true},
		{[]byte{0xC3, 0xF5, 0x48, 0xC0, 0xE3}, true},
		{[]byte{0xC3, 0xF5, 0x48, 0xC0}, false},
	}
	for _, tcase := range tests {
		got, err := Binary2Float32(tcase.data)
		if tcase.errwant {
			if err == nil {
				t.Fatalf("expected error for bytes length")
			}
		} else {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tcase.errwant && got != -3.14 {
				t.Errorf("Binary2Float32(%v) = %v; want %v", tcase.data, got, -3.14)
			}
		}
	}
}

func TestFloat32toBinary(t *testing.T) {
	wants := []float32{-5.5, 25237.174, 0.0, math.MaxFloat32, -math.MaxFloat32}

	for i, want := range wants {
		binary := Float32toBinary(want)

		got, err := Binary2Float32(binary)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
			continue
		}

		if got != want {
			t.Errorf("Test case %d: conversion failed\n  original: %f\n  binary: %v\n  restored: %f",
				i, want, binary, got)
		}
	}
}

func TestBinary2Float64(t *testing.T) {
	tests := []struct {
		data    []byte
		errwant bool
	}{
		{[]byte{0x18, 0x2D, 0x44, 0x54, 0xFB, 0x21, 0x09}, true},
		{[]byte{0x18, 0x2D, 0x44, 0x54, 0xFB, 0x21, 0x09, 0xC0, 0x3a}, true},
		{[]byte{0x18, 0x2D, 0x44, 0x54, 0xFB, 0x21, 0x09, 0xC0}, false},
	}
	for _, tcase := range tests {
		got, err := Binary2Float64(tcase.data)
		if tcase.errwant {
			if err == nil {
				t.Fatalf("expected error for bytes length")
			}
		} else {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tcase.errwant && got != -math.Pi {
				t.Errorf("Binary2Float64(%v) = %v; want %v", tcase.data, got, -math.Pi)
			}
		}
	}
}

func TestFloat64toBinary(t *testing.T) {
	t.Logf("%v", Float64toBinary(-math.Pi))

	wants := []float64{-5.5, 8640.1747517238, 0.0, math.MaxFloat32, -math.MaxFloat32, math.MaxFloat64, -math.MaxFloat64, math.Pi}

	for i, want := range wants {
		binary := Float64toBinary(want)

		got, err := Binary2Float64(binary)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
			continue
		}

		if got != want {
			t.Errorf("Test case %d: conversion failed\n  original: %f\n  binary: %v\n  restored: %f",
				i, want, binary, got)
		}
	}
}
