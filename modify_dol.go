package powerpc

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/logrusorgru/aurora/v3"
)

var (
	ErrInconsistentPatch = errors.New("before and after data present within file are not the same size")
	ErrPatchOutOfRange   = errors.New("patch cannot be applied past binary size")
	ErrInvalidPatch      = errors.New("before data present within patch did not exist in file")
)

// Patch represents a patch applied to the main binary.
type Patch struct {
	// Name is an optional name for this patch.
	// If present, its name will be logged upon application.
	Name string

	// AtOffset is the offset within the file this patch should be applied at.
	// If not present, the patch will be recursively applied across the entire file.
	// Relying on this behavior is highly discouraged, as it may damage other parts of the binary
	// if gone unchecked.
	AtOffset int

	// Before is an array of the bytes to find for, i.e. present within the original file.
	Before []byte

	// After is an array of the bytes to replace with.
	After []byte
}

// PatchSet represents multiple patches available to be applied.
// This is most useful when you have a collection of related patches.
type PatchSet struct {
	// Name is an optional name for this patch.
	// If present, its name will be logged upon application.
	Name string

	// Patches is a slice of patches to apply to the given binary.
	Patches []Patch
}

// ApplyPatch applies the given patch to the given binary.
func ApplyPatch(patch Patch, binary []byte) ([]byte, error) {
	// Print name if present
	if patch.Name != "" {
		fmt.Println(" + Applying patch", aurora.Cyan(patch.Name))
	}

	// Ensure consistency
	if len(patch.Before) != len(patch.After) {
		return nil, ErrInconsistentPatch
	}
	if patch.AtOffset != 0 && patch.AtOffset > len(binary) {
		return nil, ErrPatchOutOfRange
	}

	// Either Before or After should return the same length.
	patchLen := len(patch.Before)

	// Determine our patching behavior.
	if patch.AtOffset != 0 {
		// Ensure original bytes are present
		originalBytes := binary[patch.AtOffset : patch.AtOffset+patchLen]
		if !bytes.Equal(originalBytes, patch.Before) {
			return nil, ErrInvalidPatch
		}

		// Apply patch at the specified offset
		copy(binary[patch.AtOffset:], patch.After)
	} else {
		// Recursively apply this patch.
		// We cannot verify if the original contents are present via this.
		binary = bytes.ReplaceAll(binary, patch.Before, patch.After)
	}

	return binary, nil
}

// ApplyPatchSet applies a set of patches to a binary, noting their name.
func ApplyPatchSet(set PatchSet, binary []byte) ([]byte, error) {
	if set.Name != "" {
		fmt.Printf("Handling patch set \"%s\":\n", aurora.Yellow(set.Name))
	}

	var err error
	for _, patch := range set.Patches {
		binary, err = ApplyPatch(patch, binary)
		if err != nil {
			return nil, err
		}
	}

	return binary, err
}

// ApplyPatchSets applies an array of patch sets.
func ApplyPatchSets(sets []PatchSet, binary []byte) ([]byte, error) {
	var err error
	for _, patch := range sets {
		binary, err = ApplyPatchSet(patch, binary)
		if err != nil {
			return nil, err
		}
	}

	return binary, nil
}

// EmptyBytes returns an empty byte array of the given length.
// It is useful when creating a patch with an original value of none.
func EmptyBytes(length int) []byte {
	return bytes.Repeat([]byte{0x00}, length)
}
