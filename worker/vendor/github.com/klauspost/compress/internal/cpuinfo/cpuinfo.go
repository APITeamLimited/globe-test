// Package cpuinfo gives runtime info about the current CPU.
//
// This is a very limited module meant for use internally
// in this project. For more versatile solution check
// https://github.com/klauspost/cpuid.
package cpuinfo

// HasBMI1 checks whether an x86 CPU supports the BMI1 extension.
func HasBMI1() bool ***REMOVED***
	return hasBMI1
***REMOVED***

// HasBMI2 checks whether an x86 CPU supports the BMI2 extension.
func HasBMI2() bool ***REMOVED***
	return hasBMI2
***REMOVED***

// DisableBMI2 will disable BMI2, for testing purposes.
// Call returned function to restore previous state.
func DisableBMI2() func() ***REMOVED***
	old := hasBMI2
	hasBMI2 = false
	return func() ***REMOVED***
		hasBMI2 = old
	***REMOVED***
***REMOVED***

// HasBMI checks whether an x86 CPU supports both BMI1 and BMI2 extensions.
func HasBMI() bool ***REMOVED***
	return HasBMI1() && HasBMI2()
***REMOVED***

var hasBMI1 bool
var hasBMI2 bool
