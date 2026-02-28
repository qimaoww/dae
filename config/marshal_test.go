/*
 * SPDX-License-Identifier: AGPL-3.0-only
 * Copyright (c) 2022-2025, daeuniverse Organization <dae@v2raya.org>
 */

package config

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestMarshal(t *testing.T) {
	abs, err := filepath.Abs("../example.dae")
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(abs)
	if err != nil {
		t.Fatal(err)
	}
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "source.dae")
	if err = os.WriteFile(src, b, 0640); err != nil {
		t.Fatal(err)
	}
	merger := NewMerger(src)
	sections, _, err := merger.Merge()
	if err != nil {
		t.Fatal(err)
	}
	conf1, err := New(sections)
	if err != nil {
		t.Fatal(err)
	}
	b, err = conf1.Marshal(2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
	// Read it again.
	dst := filepath.Join(tmpDir, "roundtrip.dae")
	if err = os.WriteFile(dst, b, 0640); err != nil {
		t.Fatal(err)
	}
	sections, _, err = NewMerger(dst).Merge()
	if err != nil {
		t.Fatal(err)
	}
	conf2, err := New(sections)
	if err != nil {
		t.Fatal(err)
	}
	b2, err := conf2.Marshal(2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, b2) {
		t.Fatal("roundtrip marshal mismatch")
	}
}
