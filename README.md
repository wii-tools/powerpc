# powerpc
A simple library to handle PowerPC instruction encoding, and to enable patching ranges of a binary.
This allows simple reproduction of patches to a channel - perhaps ones downloaded via [NUS](https://github.com/wii-tools/GoNUSD).

Consider this example of a patch:
```go
example := PatchSet{
	// A name describing this patch.
	Name: "Nullify access check",
	// The offset of the function in the binary - i.e. the DOL.
	AtOffset: 57236,
	
	// Instructions present previously.
	// In this example, we have a generic function prolog.
    Before: Instructions{
        STWU(R1, R1, 0xFC10),
    }.Bytes(),
	// Instructions present afterwards.
	// They must match the same length as what they are replacing.
	// In this example, we immediately return.
    After: Instructions{
        BLR(),
    }.Bytes(),
}
patched, err := ApplyPatch(example, binary)
```

You additionally have the option to use `PatchSet`s, collections of patches that may be related to each other.
Consider the following scenario:
```go
example := PatchSet{
	Name: "Change domains",
	
	Patches: []Patch{
		Patch{
            Name: "Remove domain whitelist",
			// [...]
		},
		Patch{
		    Name: "Use custom domain",
			// [...]
        },
    },
}
patched, err := ApplyPatchSet(example, binary)
```

It is recommended to import this package without a prefix in order to permit easier usage of instruction types.
An example doing this is as follows:
```go
import (
	. "github.com/wii-tools/powerpc"
)
```

## Wishlist
 - Implement more instructions forms.
 - Use type constraints within structs to allow both byte and Instruction slices.

## Resources
 - [PowerPC mnemonics/instruction forms](https://www.nxp.com/docs/en/reference-manual/MPC82XINSET.pdf) via NXP
 - [WSC-Patcher](https://github.com/OpenShopChannel/WSC-Patcher), the original use case for encoding