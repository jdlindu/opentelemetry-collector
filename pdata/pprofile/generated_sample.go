// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated by "pdata/internal/cmd/pdatagen/main.go". DO NOT EDIT.
// To regenerate this file run "make genpdata".

package pprofile

import (
	"go.opentelemetry.io/collector/pdata/internal"
	otlpprofiles "go.opentelemetry.io/collector/pdata/internal/data/protogen/profiles/v1development"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

// Sample represents each record value encountered within a profiled program.
//
// This is a reference type, if passed by value and callee modifies it the
// caller will see the modification.
//
// Must use NewSample function to create new instances.
// Important: zero-initialized instance is not valid for use.
type Sample struct {
	orig  *otlpprofiles.Sample
	state *internal.State
}

func newSample(orig *otlpprofiles.Sample, state *internal.State) Sample {
	return Sample{orig: orig, state: state}
}

// NewSample creates a new empty Sample.
//
// This must be used only in testing code. Users should use "AppendEmpty" when part of a Slice,
// OR directly access the member if this is embedded in another struct.
func NewSample() Sample {
	state := internal.StateMutable
	return newSample(&otlpprofiles.Sample{}, &state)
}

// MoveTo moves all properties from the current struct overriding the destination and
// resetting the current instance to its zero value
func (ms Sample) MoveTo(dest Sample) {
	ms.state.AssertMutable()
	dest.state.AssertMutable()
	*dest.orig = *ms.orig
	*ms.orig = otlpprofiles.Sample{}
}

// LocationsStartIndex returns the locationsstartindex associated with this Sample.
func (ms Sample) LocationsStartIndex() int32 {
	return ms.orig.LocationsStartIndex
}

// SetLocationsStartIndex replaces the locationsstartindex associated with this Sample.
func (ms Sample) SetLocationsStartIndex(v int32) {
	ms.state.AssertMutable()
	ms.orig.LocationsStartIndex = v
}

// LocationsLength returns the locationslength associated with this Sample.
func (ms Sample) LocationsLength() int32 {
	return ms.orig.LocationsLength
}

// SetLocationsLength replaces the locationslength associated with this Sample.
func (ms Sample) SetLocationsLength(v int32) {
	ms.state.AssertMutable()
	ms.orig.LocationsLength = v
}

// Value returns the Value associated with this Sample.
func (ms Sample) Value() pcommon.Int64Slice {
	return pcommon.Int64Slice(internal.NewInt64Slice(&ms.orig.Value, ms.state))
}

// AttributeIndices returns the AttributeIndices associated with this Sample.
func (ms Sample) AttributeIndices() pcommon.Int32Slice {
	return pcommon.Int32Slice(internal.NewInt32Slice(&ms.orig.AttributeIndices, ms.state))
}

// LinkIndex returns the linkindex associated with this Sample.
func (ms Sample) LinkIndex() int32 {
	return ms.orig.GetLinkIndex()
}

// HasLinkIndex returns true if the Sample contains a
// LinkIndex value, false otherwise.
func (ms Sample) HasLinkIndex() bool {
	return ms.orig.LinkIndex_ != nil
}

// SetLinkIndex replaces the linkindex associated with this Sample.
func (ms Sample) SetLinkIndex(v int32) {
	ms.state.AssertMutable()
	ms.orig.LinkIndex_ = &otlpprofiles.Sample_LinkIndex{LinkIndex: v}
}

// RemoveLinkIndex removes the linkindex associated with this Sample.
func (ms Sample) RemoveLinkIndex() {
	ms.state.AssertMutable()
	ms.orig.LinkIndex_ = nil
}

// TimestampsUnixNano returns the TimestampsUnixNano associated with this Sample.
func (ms Sample) TimestampsUnixNano() pcommon.UInt64Slice {
	return pcommon.UInt64Slice(internal.NewUInt64Slice(&ms.orig.TimestampsUnixNano, ms.state))
}

// CopyTo copies all properties from the current struct overriding the destination.
func (ms Sample) CopyTo(dest Sample) {
	dest.state.AssertMutable()
	dest.SetLocationsStartIndex(ms.LocationsStartIndex())
	dest.SetLocationsLength(ms.LocationsLength())
	ms.Value().CopyTo(dest.Value())
	ms.AttributeIndices().CopyTo(dest.AttributeIndices())
	if ms.HasLinkIndex() {
		dest.SetLinkIndex(ms.LinkIndex())
	}

	ms.TimestampsUnixNano().CopyTo(dest.TimestampsUnixNano())
}
