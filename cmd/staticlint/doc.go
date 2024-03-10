/*
Package main contains of analyzers inside multichecker

#Usage

Usage: staticlint [-flag] [package]

Core flags:

	-S1000
	      enable S1000 analysis
	-S1001
	      enable S1001 analysis
	-S1002
	      enable S1002 analysis
	-S1003
	      enable S1003 analysis
	-S1004
	      enable S1004 analysis
	-S1005
	      enable S1005 analysis
	-S1006
	      enable S1006 analysis
	-S1007
	      enable S1007 analysis
	-S1008
	      enable S1008 analysis
	-S1009
	      enable S1009 analysis
	-S1010
	      enable S1010 analysis
	-S1011
	      enable S1011 analysis
	-S1012
	      enable S1012 analysis
	-S1016
	      enable S1016 analysis
	-S1017
	      enable S1017 analysis
	-S1018
	      enable S1018 analysis
	-S1019
	      enable S1019 analysis
	-S1020
	      enable S1020 analysis
	-S1021
	      enable S1021 analysis
	-S1023
	      enable S1023 analysis
	-S1024
	      enable S1024 analysis
	-S1025
	      enable S1025 analysis
	-S1028
	      enable S1028 analysis
	-S1029
	      enable S1029 analysis
	-S1030
	      enable S1030 analysis
	-S1031
	      enable S1031 analysis
	-S1032
	      enable S1032 analysis
	-S1033
	      enable S1033 analysis
	-S1034
	      enable S1034 analysis
	-S1035
	      enable S1035 analysis
	-S1036
	      enable S1036 analysis
	-S1037
	      enable S1037 analysis
	-S1038
	      enable S1038 analysis
	-S1039
	      enable S1039 analysis
	-S1040
	      enable S1040 analysis
	-SA1000
	      enable SA1000 analysis
	-SA1001
	      enable SA1001 analysis
	-SA1002
	      enable SA1002 analysis
	-SA1003
	      enable SA1003 analysis
	-SA1004
	      enable SA1004 analysis
	-SA1005
	      enable SA1005 analysis
	-SA1006
	      enable SA1006 analysis
	-SA1007
	      enable SA1007 analysis
	-SA1008
	      enable SA1008 analysis
	-SA1010
	      enable SA1010 analysis
	-SA1011
	      enable SA1011 analysis
	-SA1012
	      enable SA1012 analysis
	-SA1013
	      enable SA1013 analysis
	-SA1014
	      enable SA1014 analysis
	-SA1015
	      enable SA1015 analysis
	-SA1016
	      enable SA1016 analysis
	-SA1017
	      enable SA1017 analysis
	-SA1018
	      enable SA1018 analysis
	-SA1019
	      enable SA1019 analysis
	-SA1020
	      enable SA1020 analysis
	-SA1021
	      enable SA1021 analysis
	-SA1023
	      enable SA1023 analysis
	-SA1024
	      enable SA1024 analysis
	-SA1025
	      enable SA1025 analysis
	-SA1026
	      enable SA1026 analysis
	-SA1027
	      enable SA1027 analysis
	-SA1028
	      enable SA1028 analysis
	-SA1029
	      enable SA1029 analysis
	-SA1030
	      enable SA1030 analysis
	-SA2000
	      enable SA2000 analysis
	-SA2001
	      enable SA2001 analysis
	-SA2002
	      enable SA2002 analysis
	-SA2003
	      enable SA2003 analysis
	-SA3000
	      enable SA3000 analysis
	-SA3001
	      enable SA3001 analysis
	-SA4000
	      enable SA4000 analysis
	-SA4001
	      enable SA4001 analysis
	-SA4003
	      enable SA4003 analysis
	-SA4004
	      enable SA4004 analysis
	-SA4005
	      enable SA4005 analysis
	-SA4006
	      enable SA4006 analysis
	-SA4008
	      enable SA4008 analysis
	-SA4009
	      enable SA4009 analysis
	-SA4010
	      enable SA4010 analysis
	-SA4011
	      enable SA4011 analysis
	-SA4012
	      enable SA4012 analysis
	-SA4013
	      enable SA4013 analysis
	-SA4014
	      enable SA4014 analysis
	-SA4015
	      enable SA4015 analysis
	-SA4016
	      enable SA4016 analysis
	-SA4017
	      enable SA4017 analysis
	-SA4018
	      enable SA4018 analysis
	-SA4019
	      enable SA4019 analysis
	-SA4020
	      enable SA4020 analysis
	-SA4021
	      enable SA4021 analysis
	-SA4022
	      enable SA4022 analysis
	-SA4023
	      enable SA4023 analysis
	-SA4024
	      enable SA4024 analysis
	-SA4025
	      enable SA4025 analysis
	-SA4026
	      enable SA4026 analysis
	-SA4027
	      enable SA4027 analysis
	-SA4028
	      enable SA4028 analysis
	-SA4029
	      enable SA4029 analysis
	-SA4030
	      enable SA4030 analysis
	-SA4031
	      enable SA4031 analysis
	-SA5000
	      enable SA5000 analysis
	-SA5001
	      enable SA5001 analysis
	-SA5002
	      enable SA5002 analysis
	-SA5003
	      enable SA5003 analysis
	-SA5004
	      enable SA5004 analysis
	-SA5005
	      enable SA5005 analysis
	-SA5007
	      enable SA5007 analysis
	-SA5008
	      enable SA5008 analysis
	-SA5009
	      enable SA5009 analysis
	-SA5010
	      enable SA5010 analysis
	-SA5011
	      enable SA5011 analysis
	-SA5012
	      enable SA5012 analysis
	-SA6000
	      enable SA6000 analysis
	-SA6001
	      enable SA6001 analysis
	-SA6002
	      enable SA6002 analysis
	-SA6003
	      enable SA6003 analysis
	-SA6005
	      enable SA6005 analysis
	-SA9001
	      enable SA9001 analysis
	-SA9002
	      enable SA9002 analysis
	-SA9003
	      enable SA9003 analysis
	-SA9004
	      enable SA9004 analysis
	-SA9005
	      enable SA9005 analysis
	-SA9006
	      enable SA9006 analysis
	-SA9007
	      enable SA9007 analysis
	-SA9008
	      enable SA9008 analysis
	-V    print version and exit
	-all
	      no effect (deprecated)
	-asmdecl
	      enable asmdecl analysis
	-assign
	      enable assign analysis
	-atomic
	      enable atomic analysis
	-atomicalign
	      enable atomicalign analysis
	-bool
	      deprecated alias for -bools
	-bools
	      enable bools analysis
	-buildssa
	      enable buildssa analysis
	-buildtag
	      enable buildtag analysis
	-buildtags
	      deprecated alias for -buildtag
	-c int
	      display offending line with this many lines of context (default -1)
	-cgocall
	      enable cgocall analysis
	-composites
	      enable composites analysis
	-compositewhitelist
	      deprecated alias for -composites.whitelist (default true)
	-copylocks
	      enable copylocks analysis
	-cpuprofile string
	      write CPU profile to this file
	-ctrlflow
	      enable ctrlflow analysis
	-debug string
	      debug flags, any subset of "fpstv"
	-deepequalerrors
	      enable deepequalerrors analysis
	-directive
	      enable directive analysis
	-errcheck
	      enable errcheck analysis
	-errorsas
	      enable errorsas analysis
	-fieldalignment
	      enable fieldalignment analysis
	-findcall
	      enable findcall analysis
	-fix
	      apply all suggested fixes
	-flags
	      print analyzer flags in JSON
	-framepointer
	      enable framepointer analysis
	-httpresponse
	      enable httpresponse analysis
	-ifaceassert
	      enable ifaceassert analysis
	-ineffassign
	      enable ineffassign analysis
	-inspect
	      enable inspect analysis
	-json
	      emit JSON output
	-loopclosure
	      enable loopclosure analysis
	-lostcancel
	      enable lostcancel analysis
	-mainosexit
	      enable mainosexit analysis
	-memprofile string
	      write memory profile to this file
	-methods
	      deprecated alias for -stdmethods
	-nilfunc
	      enable nilfunc analysis
	-nilness
	      enable nilness analysis
	-pkgfact
	      enable pkgfact analysis
	-printf
	      enable printf analysis
	-printfuncs value
	      deprecated alias for -printf.funcs (default (*log.Logger).Fatal,(*log.Logger).Fatalf,(*log.Logger).Fatalln,(*log.Logger).Panic,(*log.Logger).Panicf,(*log.Logger).Panicln,(*log.Logger).Print,(*log.Logger).Printf,(*log.Logger).Println,(*testing.common).Error,(*testing.common).Errorf,(*testing.common).Fatal,(*testing.common).Fatalf,(*testing.common).Log,(*testing.common).Logf,(*testing.common).Skip,(*testing.common).Skipf,(testing.TB).Error,(testing.TB).Errorf,(testing.TB).Fatal,(testing.TB).Fatalf,(testing.TB).Log,(testing.TB).Logf,(testing.TB).Skip,(testing.TB).Skipf,fmt.Append,fmt.Appendf,fmt.Appendln,fmt.Errorf,fmt.Fprint,fmt.Fprintf,fmt.Fprintln,fmt.Print,fmt.Printf,fmt.Println,fmt.Sprint,fmt.Sprintf,fmt.Sprintln,log.Fatal,log.Fatalf,log.Fatalln,log.Panic,log.Panicf,log.Panicln,log.Print,log.Printf,log.Println,runtime/trace.Logf)
	-rangeloops
	      deprecated alias for -loopclosure
	-reflectvaluecompare
	      enable reflectvaluecompare analysis
	-shadow
	      enable shadow analysis
	-shadowstrict
	      deprecated alias for -shadow.strict
	-shift
	      enable shift analysis
	-sigchanyzer
	      enable sigchanyzer analysis
	-slog
	      enable slog analysis
	-sortslice
	      enable sortslice analysis
	-source
	      no effect (deprecated)
	-stdmethods
	      enable stdmethods analysis
	-stringintconv
	      enable stringintconv analysis
	-structtag
	      enable structtag analysis
	-tags string
	      no effect (deprecated)
	-test
	      indicates whether test files should be analyzed, too (default true)
	-testinggoroutine
	      enable testinggoroutine analysis
	-tests
	      enable tests analysis
	-timeformat
	      enable timeformat analysis
	-trace string
	      write trace log to this file
	-unmarshal
	      enable unmarshal analysis
	-unreachable
	      enable unreachable analysis
	-unsafeptr
	      enable unsafeptr analysis
	-unusedfuncs value
	      deprecated alias for -unusedresult.funcs (default context.WithCancel,context.WithDeadline,context.WithTimeout,context.WithValue,errors.New,fmt.Errorf,fmt.Sprint,fmt.Sprintf,slices.Clip,slices.Compact,slices.CompactFunc,slices.Delete,slices.DeleteFunc,slices.Grow,slices.Insert,slices.Replace,sort.Reverse)
	-unusedresult
	      enable unusedresult analysis
	-unusedstringmethods value
	      deprecated alias for -unusedresult.stringmethods (default Error,String)
	-unusedwrite
	      enable unusedwrite analysis
	-usesgenerics
	      enable usesgenerics analysis
	-v    no effect (deprecated)

# Content
It contains:
  - analyzers from golang.org/x/tools/go/analysis/passes
  - analyzers from golang.org/x/tools/go/analysis/passes
  - static check analyzers from class SA
  - check analyzers for code simplicity
  - Add third party analyzers
  - Custom identifying call of os.Exit inside main func of main package analyzer

# analyzers from golang.org/x/tools/go/analysis/passes

  - asmdecl - defines an Analyzer that reports mismatches between assembly files and Go declarations.
  - assign - defines an Analyzer that detects useless assignments.
  - atomic - defines an Analyzer that checks for common mistakes using the sync/atomic package.
  - atomicalign - defines an Analyzer that checks for non-64-bit-aligned arguments to sync/atomic functions.
  - bools - defines an Analyzer that detects common mistakes involving boolean operators.
  - buildssa - defines an Analyzer that constructs the SSA representation of an error-free package and returns the set of all functions within it.
  - buildtag - defines an Analyzer that checks build tags.
  - cgocall - defines an Analyzer that detects some violations of the cgo pointer passing rules.
  - composite - defines an Analyzer that checks for unkeyed composite literals.
  - copylock - defines an Analyzer that checks for locks erroneously passed by value.
  - ctrlflow - is an analysis that provides a syntactic control-flow graph (CFG) for the body of a function.
  - deepequalerrors - defines an Analyzer that checks for the use of reflect.DeepEqual with error values.
  - directive - defines an Analyzer that checks known Go toolchain directives.
  - errorsas - defines an Analyzer that checks that the second argument to errors.As is a pointer to a type implementing error.
  - fieldalignment - defines an Analyzer that detects structs that would use less memory if their fields were sorted.
  - findcall - defines an Analyzer that serves as a trivial example and test of the Analysis API.
  - framepointer - defines an Analyzer that reports assembly code that clobbers the frame pointer before saving it.
  - httpresponse - defines an Analyzer that checks for mistakes using HTTP responses.
  - ifaceassert - defines an Analyzer that flags impossible interface-interface type assertions.
  - inspect - defines an Analyzer that provides an AST inspector (golang.org/x/tools/go/ast/inspector.Inspector) for the syntax trees of a package.
  - loopclosure - defines an Analyzer that checks for references to enclosing loop variables from within nested functions.
  - lostcancel - defines an Analyzer that checks for failure to call a context cancellation function.
  - nilfunc - defines an Analyzer that checks for useless comparisons against nil.
  - nilness - inspects the control-flow graph of an SSA function and reports errors such as nil pointer dereferences and degenerate nil pointer comparisons.
  - pkgfact - is a demonstration and test of the package fact mechanism.
  - printf - defines an Analyzer that checks consistency of Printf format strings and arguments.
  - reflectvaluecompare - defines an Analyzer that checks for accidentally using == or reflect.DeepEqual to compare reflect.Value values.
  - shadow - defines an Analyzer that checks for shadowed variables.
  - shift - defines an Analyzer that checks for shifts that exceed the width of an integer.
  - sigchanyzer - defines an Analyzer that detects misuse of unbuffered signal as argument to signal.Notify.
  - slog - defines an Analyzer that checks for mismatched key-value pairs in log/slog calls.
  - sortslice - defines an Analyzer that checks for calls to sort.Slice that do not use a slice type as first argument.
  - stdmethods - defines an Analyzer that checks for misspellings in the signatures of methods similar to well-known interfaces.
  - stringintconv - defines an Analyzer that flags type conversions from integers to strings.
  - structtag - defines an Analyzer that checks struct field tags are well formed.
  - testinggoroutine - defines an Analyzerfor detecting calls to Fatal from a test goroutine.
  - tests - defines an Analyzer that checks for common mistaken usages of tests and examples.
  - timeformat - defines an Analyzer that checks for the use of time.Format or time.Parse calls with a bad format.
  - unmarshal - package defines an Analyzer that checks for passing non-pointer or non-interface types to unmarshal and decode functions.
  - unreachable - defines an Analyzer that checks for unreachable code.
  - unsafeptr - defines an Analyzer that checks for invalid conversions of uintptr to unsafe.Pointer.
  - unusedresult - defines an analyzer that checks for unused results of calls to certain pure functions.
  - unusedwrite - checks for unused writes to the elements of a struct or array object.
  - usesgenerics - defines an Analyzer that checks for usage of generic features added in Go 1.18.

# static check analyzers from class SA
  - SA9005 Trying to marshal a struct with no public fields nor custom marshaling
  - SA1001 Invalid template
  - SA1012 A nil context.Context is being passed to a function, consider using context.TODO instead
  - SA1019 Using a deprecated function, variable, constant or field
  - SA4010 The result of append will never be observed anywhere
  - SA5010 Impossible type assertion
  - SA9001 Defers in range loops may not run when you expect them to
  - SA1028 sort.Slice can only be used on slices
  - SA4006 A value assigned to a variable is never read before being overwritten. Forgotten error check or dead code?
  - SA4018 Self-assignment of variables
  - SA5005 The finalizer references the finalized object, preventing garbage collection
  - SA6005 Inefficient string comparison with strings.ToLower or strings.ToUpper
  - SA9007 Deleting a directory that shouldn't be deleted
  - SA5008 Invalid struct tag
  - SA9002 Using a non-octal os.FileMode that looks like it was meant to be in octal.
  - SA1004 Suspiciously small untyped constant in time.Sleep
  - SA1016 Trapping a signal that cannot be trapped
  - SA1026 Cannot marshal channels or functions
  - SA4013 Negating a boolean twice (!!b) is the same as writing b. This is either redundant, or a typo.
  - SA4019 Multiple, identical build constraints in the same file
  - SA4028 x % 1 is always zero
  - SA9003 Empty body in an if or else branch
  - SA1000 Invalid regular expression
  - SA1008 Non-canonical key in http.Header map
  - SA1011 Various methods in the 'strings' package expect valid UTF-8, but invalid input is provided
  - SA1023 Modifying the buffer in an io.Writer implementation
  - SA4021 'x = append(y)' is equivalent to 'x = y'
  - SA9008 else branch of a type assertion is probably not reading the right value
  - SA1007 Invalid URL in net/url.Parse
  - SA2001 Empty critical section, did you mean to defer the unlock?
  - SA2002 Called testing.T.FailNow or SkipNow in a goroutine, which isn't allowed
  - SA4005 Field assignment that will never be observed. Did you mean to use a pointer receiver?
  - SA4011 Break statement with no effect. Did you mean to break out of an outer loop?
  - SA5003 Defers in infinite loops will never execute
  - SA1013 io.Seeker.Seek is being called with the whence constant as the first argument, but it should be the second
  - SA1014 Non-pointer value passed to Unmarshal or Decode
  - SA4012 Comparing a value against NaN even though no value is equal to NaN
  - SA4024 Checking for impossible return value from a builtin function
  - SA5001 Deferring Close before checking for a possible error
  - SA5002 The empty for loop ('for {}') spins and can block the scheduler
  - SA5009 Invalid Printf call
  - SA4025 Integer division of literals that results in zero
  - SA4030 Ineffective attempt at generating random number
  - SA1010 (*regexp.Regexp).FindAll called with n == 0, which will always return zero results
  - SA1018 strings.Replace called with n == 0, which does nothing
  - SA1027 Atomic access to 64-bit variable must be 64-bit aligned
  - SA1029 Inappropriate key in call to context.WithValue
  - SA4003 Comparing unsigned values against negative values is pointless
  - SA4009 A function argument is overwritten before its first use
  - SA6000 Using regexp.Match or related in a loop, should use regexp.Compile
  - SA1002 Invalid format in time.Parse
  - SA1030 Invalid argument in call to a strconv function
  - SA4014 An if/else if chain has repeated conditions and no side-effects; if the condition didn't match the first time, it won't match the second time, either
  - SA4004 The loop exits unconditionally after one iteration
  - SA4017 Discarding the return values of a function without side effects, making the call pointless
  - SA4027 (*net/url.URL).Query returns a copy, modifying it doesn't change the URL
  - SA5011 Possible nil pointer dereference
  - SA6002 Storing non-pointer values in sync.Pool allocates memory
  - SA6003 Converting a string to a slice of runes before ranging over it
  - SA9004 Only the first constant has an explicit type
  - SA1024 A string cutset contains duplicate characters
  - SA3000 TestMain doesn't call os.Exit, hiding test failures
  - SA4000 Binary operator has identical expressions on both sides
  - SA4001 &*x gets simplified to x, it does not copy x
  - SA4020 Unreachable case clause in a type switch
  - SA4026 Go constants cannot express negative zero
  - SA1017 Channels used with os/signal.Notify should be buffered
  - SA2003 Deferred Lock right after locking, likely meant to defer Unlock instead
  - SA4008 The variable in the loop condition never changes, are you incrementing the wrong variable?
  - SA1020 Using an invalid host:port pair with a net.Listen-related function
  - SA2000 sync.WaitGroup.Add called inside the goroutine, leading to a race condition
  - SA3001 Assigning to b.N in benchmarks distorts the results
  - SA4016 Certain bitwise operations, such as x ^ 0, do not do anything useful
  - SA9006 Dubious bit shifting of a fixed size integer value
  - SA1021 Using bytes.Equal to compare two net.IP
  - SA1025 It is not possible to use (*time.Timer).Reset's return value correctly
  - SA4015 Calling functions like math.Ceil on floats converted from integers doesn't do anything useful
  - SA1003 Unsupported argument to functions in encoding/binary
  - SA1015 Using time.Tick in a way that will leak. Consider using time.NewTicker, and only use time.Tick in tests, commands and endless functions
  - SA4029 Ineffective attempt at sorting slice
  - SA5004 'for { select { ...' with an empty default branch spins
  - SA5007 Infinite recursive call
  - SA5012 Passing odd-sized slice to function expecting even size
  - SA6001 Missing an optimization opportunity when indexing maps by byte slices
  - SA1005 Invalid first argument to exec.Command
  - SA1006 Printf with dynamic first argument and no further arguments
  - SA4022 Comparing the address of a variable against nil
  - SA4023 Impossible comparison of interface value with untyped nil
  - SA4031 Checking never-nil value against nil
  - SA5000 Assignment to nil map

# check analyzers for code simplicity
  - S1028 Simplify error construction with fmt.Errorf
  - S1031 Omit redundant nil check around loop
  - S1005 Drop unnecessary use of the blank identifier
  - S1012 Replace time.Now().Sub(x) with time.Since(x)
  - S1006 Use 'for { ... }' for infinite loops
  - S1021 Merge variable declaration and assignment
  - S1033 Unnecessary guard around call to 'delete'
  - S1000 Use plain channel send or receive instead of single-case select
  - S1002 Omit comparison with boolean constant
  - S1030 Use bytes.Buffer.String or bytes.Buffer.Bytes
  - S1010 Omit default slice index
  - S1019 Simplify 'make' call by omitting redundant arguments
  - S1035 Redundant call to net/http.CanonicalHeaderKey in method call on net/http.Header
  - S1038 Unnecessarily complex way of printing formatted string
  - S1024 Replace x.Sub(time.Now()) with time.Until(x)
  - S1032 Use sort.Ints(x), sort.Float64s(x), and sort.Strings(x)
  - S1009 Omit redundant nil check on slices
  - S1025 Don't use fmt.Sprintf("%s", x) unnecessarily
  - S1034 Use result of type assertion to simplify cases
  - S1036 Unnecessary guard around map access
  - S1037 Elaborate way of sleeping
  - S1001 Replace for loop with call to copy
  - S1007 Simplify regular expression by using raw string literal
  - S1020 Omit redundant nil check in type assertion
  - S1029 Range over the string directly
  - S1039 Unnecessary use of fmt.Sprint
  - S1011 Use a single append to concatenate two slices
  - S1016 Use a type conversion instead of manually copying struct fields
  - S1008 Simplify returning boolean expression
  - S1040 Type assertion to current type
  - S1017 Replace manual trimming with strings.TrimPrefix
  - S1018 Use 'copy' for sliding elements
  - S1003 Replace call to strings.Index with strings.Contains
  - S1004 Replace call to bytes.Compare with bytes.Equal

# Add third party analyzers
  - errcheck - check for unchecked errors
  - ineffassign - detect ineffectual assignments in Go code

# Custom identifying call of os.Exit inside main func of main package analyzer
  - mainosexit - analyzes if os.Exit called in main function of main package
*/
package main
