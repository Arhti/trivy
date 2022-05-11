// Code generated by 'ccgo uuid/gen.c -crt-import-path "" -export-defines "" -export-enums "" -export-externs X -export-fields F -export-structs "" -export-typedefs "" -header -hide _OSSwapInt16,_OSSwapInt32,_OSSwapInt64 -o uuid/uuid_openbsd_amd64.go -pkgname uuid', DO NOT EDIT.

package uuid

import (
	"math"
	"reflect"
	"sync/atomic"
	"unsafe"
)

var _ = math.Pi
var _ reflect.Kind
var _ atomic.Value
var _ unsafe.Pointer

const (
	BIG_ENDIAN                 = 4321
	BYTE_ORDER                 = 1234
	LITTLE_ENDIAN              = 1234
	PDP_ENDIAN                 = 3412
	UUID_BUF_LEN               = 38
	UUID_STR_LEN               = 36
	X_BIG_ENDIAN               = 4321
	X_BYTE_ORDER               = 1234
	X_CLOCKID_T_DEFINED_       = 0
	X_CLOCK_T_DEFINED_         = 0
	X_FILE_OFFSET_BITS         = 64
	X_INT16_T_DEFINED_         = 0
	X_INT32_T_DEFINED_         = 0
	X_INT64_T_DEFINED_         = 0
	X_INT8_T_DEFINED_          = 0
	X_LITTLE_ENDIAN            = 1234
	X_LP64                     = 1
	X_MACHINE_CDEFS_H_         = 0
	X_MACHINE_ENDIAN_H_        = 0
	X_MACHINE__TYPES_H_        = 0
	X_MAX_PAGE_SHIFT           = 12
	X_OFF_T_DEFINED_           = 0
	X_PDP_ENDIAN               = 3412
	X_PID_T_DEFINED_           = 0
	X_QUAD_HIGHWORD            = 1
	X_QUAD_LOWWORD             = 0
	X_RET_PROTECTOR            = 1
	X_SIZE_T_DEFINED_          = 0
	X_SSIZE_T_DEFINED_         = 0
	X_STACKALIGNBYTES          = 15
	X_SYS_CDEFS_H_             = 0
	X_SYS_ENDIAN_H_            = 0
	X_SYS_TYPES_H_             = 0
	X_SYS_UUID_H_              = 0
	X_SYS__ENDIAN_H_           = 0
	X_SYS__TYPES_H_            = 0
	X_TIMER_T_DEFINED_         = 0
	X_TIME_T_DEFINED_          = 0
	X_UINT16_T_DEFINED_        = 0
	X_UINT32_T_DEFINED_        = 0
	X_UINT64_T_DEFINED_        = 0
	X_UINT8_T_DEFINED_         = 0
	X_UUID_BUF_LEN             = 38
	X_UUID_H_                  = 0
	X_UUID_NODE_LEN            = 6
	Unix                       = 1
	Uuid_s_bad_version         = 1
	Uuid_s_invalid_string_uuid = 2
	Uuid_s_no_memory           = 3
	Uuid_s_ok                  = 0
)

type Ptrdiff_t = int64 /* <builtin>:3:26 */

type Size_t = uint64 /* <builtin>:9:23 */

type Wchar_t = int32 /* <builtin>:15:24 */

type X__int128_t = struct {
	Flo int64
	Fhi int64
} /* <builtin>:21:43 */ // must match modernc.org/mathutil.Int128
type X__uint128_t = struct {
	Flo uint64
	Fhi uint64
} /* <builtin>:22:44 */ // must match modernc.org/mathutil.Int128

type X__builtin_va_list = uintptr /* <builtin>:46:14 */
type X__float128 = float64        /* <builtin>:47:21 */

//	$OpenBSD: uuid.h,v 1.1 2014/08/31 09:36:36 miod Exp $
//	$NetBSD: uuid.h,v 1.2 2008/04/23 07:52:32 plunky Exp $

// Copyright (c) 2002 Marcel Moolenaar
// Copyright (c) 2002 Hiten Mahesh Pandya
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS ``AS IS'' AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.
//
// $FreeBSD: src/include/uuid.h,v 1.2 2002/11/05 10:55:16 jmallett Exp $

//	$OpenBSD: types.h,v 1.48 2019/02/09 04:54:11 guenther Exp $
//	$NetBSD: types.h,v 1.29 1996/11/15 22:48:25 jtc Exp $

// -
// Copyright (c) 1982, 1986, 1991, 1993
//	The Regents of the University of California.  All rights reserved.
// (c) UNIX System Laboratories, Inc.
// All or some portions of this file are derived from material licensed
// to the University of California by American Telephone and Telegraph
// Co. or Unix System Laboratories, Inc. and are reproduced herein with
// the permission of UNIX System Laboratories, Inc.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
// 3. Neither the name of the University nor the names of its contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE REGENTS AND CONTRIBUTORS ``AS IS'' AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.
//
//	@(#)types.h	8.4 (Berkeley) 1/21/94

//	$OpenBSD: cdefs.h,v 1.43 2018/10/29 17:10:40 guenther Exp $
//	$NetBSD: cdefs.h,v 1.16 1996/04/03 20:46:39 christos Exp $

// Copyright (c) 1991, 1993
//	The Regents of the University of California.  All rights reserved.
//
// This code is derived from software contributed to Berkeley by
// Berkeley Software Design, Inc.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
// 3. Neither the name of the University nor the names of its contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE REGENTS AND CONTRIBUTORS ``AS IS'' AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.
//
//	@(#)cdefs.h	8.7 (Berkeley) 1/21/94

//	$OpenBSD: cdefs.h,v 1.3 2013/03/28 17:30:45 martynas Exp $

// Written by J.T. Conklin <jtc@wimsey.com> 01/17/95.
// Public domain.

// Macro to test if we're using a specific version of gcc or later.

// The __CONCAT macro is used to concatenate parts of symbol names, e.g.
// with "#define OLD(foo) __CONCAT(old,foo)", OLD(foo) produces oldfoo.
// The __CONCAT macro is a bit tricky -- make sure you don't put spaces
// in between its arguments.  Do not use __CONCAT on double-quoted strings,
// such as those from the __STRING macro: to concatenate strings just put
// them next to each other.

// GCC1 and some versions of GCC2 declare dead (non-returning) and
// pure (no side effects) functions using "volatile" and "const";
// unfortunately, these then cause warnings under "-ansi -pedantic".
// GCC >= 2.5 uses the __attribute__((attrs)) style.  All of these
// work for GNU C++ (modulo a slight glitch in the C++ grammar in
// the distribution version of 2.5.5).

// __returns_twice makes the compiler not assume the function
// only returns once.  This affects registerisation of variables:
// even local variables need to be in memory across such a call.
// Example: setjmp()

// __only_inline makes the compiler only use this function definition
// for inlining; references that can't be inlined will be left as
// external references instead of generating a local copy.  The
// matching library should include a simple extern definition for
// the function to handle those references.  c.f. ctype.h

// GNU C version 2.96 adds explicit branch prediction so that
// the CPU back-end can hint the processor and also so that
// code blocks can be reordered such that the predicted path
// sees a more linear flow, thus improving cache behavior, etc.
//
// The following two macros provide us with a way to utilize this
// compiler feature.  Use __predict_true() if you expect the expression
// to evaluate to true, and __predict_false() if you expect the
// expression to evaluate to false.
//
// A few notes about usage:
//
//	* Generally, __predict_false() error condition checks (unless
//	  you have some _strong_ reason to do otherwise, in which case
//	  document it), and/or __predict_true() `no-error' condition
//	  checks, assuming you want to optimize for the no-error case.
//
//	* Other than that, if you don't know the likelihood of a test
//	  succeeding from empirical or other `hard' evidence, don't
//	  make predictions.
//
//	* These are meant to be used in places that are run `a lot'.
//	  It is wasteful to make predictions in code that is run
//	  seldomly (e.g. at subsystem initialization time) as the
//	  basic block reordering that this affects can often generate
//	  larger code.

// Delete pseudo-keywords wherever they are not available or needed.

// The __packed macro indicates that a variable or structure members
// should have the smallest possible alignment, despite any host CPU
// alignment requirements.
//
// The __aligned(x) macro specifies the minimum alignment of a
// variable or structure.
//
// These macros together are useful for describing the layout and
// alignment of messages exchanged with hardware or other systems.

// "The nice thing about standards is that there are so many to choose from."
// There are a number of "feature test macros" specified by (different)
// standards that determine which interfaces and types the header files
// should expose.
//
// Because of inconsistencies in these macros, we define our own
// set in the private name space that end in _VISIBLE.  These are
// always defined and so headers can test their values easily.
// Things can get tricky when multiple feature macros are defined.
// We try to take the union of all the features requested.
//
// The following macros are guaranteed to have a value after cdefs.h
// has been included:
//	__POSIX_VISIBLE
//	__XPG_VISIBLE
//	__ISO_C_VISIBLE
//	__BSD_VISIBLE

// X/Open Portability Guides and Single Unix Specifications.
// _XOPEN_SOURCE				XPG3
// _XOPEN_SOURCE && _XOPEN_VERSION = 4		XPG4
// _XOPEN_SOURCE && _XOPEN_SOURCE_EXTENDED = 1	XPG4v2
// _XOPEN_SOURCE == 500				XPG5
// _XOPEN_SOURCE == 520				XPG5v2
// _XOPEN_SOURCE == 600				POSIX 1003.1-2001 with XSI
// _XOPEN_SOURCE == 700				POSIX 1003.1-2008 with XSI
//
// The XPG spec implies a specific value for _POSIX_C_SOURCE.

// POSIX macros, these checks must follow the XOPEN ones above.
//
// _POSIX_SOURCE == 1		1003.1-1988 (superseded by _POSIX_C_SOURCE)
// _POSIX_C_SOURCE == 1		1003.1-1990
// _POSIX_C_SOURCE == 2		1003.2-1992
// _POSIX_C_SOURCE == 199309L	1003.1b-1993
// _POSIX_C_SOURCE == 199506L   1003.1c-1995, 1003.1i-1995,
//				and the omnibus ISO/IEC 9945-1:1996
// _POSIX_C_SOURCE == 200112L   1003.1-2001
// _POSIX_C_SOURCE == 200809L   1003.1-2008
//
// The POSIX spec implies a specific value for __ISO_C_VISIBLE, though
// this may be overridden by the _ISOC99_SOURCE macro later.

// _ANSI_SOURCE means to expose ANSI C89 interfaces only.
// If the user defines it in addition to one of the POSIX or XOPEN
// macros, assume the POSIX/XOPEN macro(s) should take precedence.

// _ISOC99_SOURCE, _ISOC11_SOURCE, __STDC_VERSION__, and __cplusplus
// override any of the other macros since they are non-exclusive.

// Finally deal with BSD-specific interfaces that are not covered
// by any standards.  We expose these when none of the POSIX or XPG
// macros is defined or if the user explicitly asks for them.

// Default values.

//	$OpenBSD: endian.h,v 1.25 2014/12/21 04:49:00 guenther Exp $

// -
// Copyright (c) 1997 Niklas Hallqvist.  All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR
// IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
// THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Public definitions for little- and big-endian systems.
// This file should be included as <endian.h> in userspace and as
// <sys/endian.h> in the kernel.
//
// System headers that need endian information but that can't or don't
// want to export the public names here should include <sys/_endian.h>
// and use the internal names: _BYTE_ORDER, _*_ENDIAN, etc.

//	$OpenBSD: cdefs.h,v 1.43 2018/10/29 17:10:40 guenther Exp $
//	$NetBSD: cdefs.h,v 1.16 1996/04/03 20:46:39 christos Exp $

// Copyright (c) 1991, 1993
//	The Regents of the University of California.  All rights reserved.
//
// This code is derived from software contributed to Berkeley by
// Berkeley Software Design, Inc.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
// 3. Neither the name of the University nor the names of its contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE REGENTS AND CONTRIBUTORS ``AS IS'' AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.
//
//	@(#)cdefs.h	8.7 (Berkeley) 1/21/94

//	$OpenBSD: _endian.h,v 1.8 2018/01/11 23:13:37 dlg Exp $

// -
// Copyright (c) 1997 Niklas Hallqvist.  All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR
// IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
// THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Internal endianness macros.  This pulls in <machine/endian.h> to
// get the correct setting direction for the platform and sets internal
// ('__' prefix) macros appropriately.

//	$OpenBSD: _types.h,v 1.9 2014/08/22 23:05:15 krw Exp $

// -
// Copyright (c) 1990, 1993
//	The Regents of the University of California.  All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
// 3. Neither the name of the University nor the names of its contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE REGENTS AND CONTRIBUTORS ``AS IS'' AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.
//
//	@(#)types.h	8.3 (Berkeley) 1/5/94

//	$OpenBSD: _types.h,v 1.17 2018/03/05 01:15:25 deraadt Exp $

// -
// Copyright (c) 1990, 1993
//	The Regents of the University of California.  All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
// 3. Neither the name of the University nor the names of its contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE REGENTS AND CONTRIBUTORS ``AS IS'' AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.
//
//	@(#)types.h	8.3 (Berkeley) 1/5/94
//	@(#)ansi.h	8.2 (Berkeley) 1/4/94

// _ALIGN(p) rounds p (pointer or byte index) up to a correctly-aligned
// value for all data types (int, long, ...).   The result is an
// unsigned long and must be cast to any desired pointer type.
//
// _ALIGNED_POINTER is a boolean macro that checks whether an address
// is valid to fetch data elements of type t from on this architecture.
// This does not reflect the optimal alignment, just the possibility
// (within reasonable limits).

// 7.18.1.1 Exact-width integer types
type X__int8_t = int8     /* _types.h:61:22 */
type X__uint8_t = uint8   /* _types.h:62:24 */
type X__int16_t = int16   /* _types.h:63:17 */
type X__uint16_t = uint16 /* _types.h:64:25 */
type X__int32_t = int32   /* _types.h:65:15 */
type X__uint32_t = uint32 /* _types.h:66:23 */
type X__int64_t = int64   /* _types.h:67:20 */
type X__uint64_t = uint64 /* _types.h:68:28 */

// 7.18.1.2 Minimum-width integer types
type X__int_least8_t = X__int8_t     /* _types.h:71:19 */
type X__uint_least8_t = X__uint8_t   /* _types.h:72:20 */
type X__int_least16_t = X__int16_t   /* _types.h:73:20 */
type X__uint_least16_t = X__uint16_t /* _types.h:74:21 */
type X__int_least32_t = X__int32_t   /* _types.h:75:20 */
type X__uint_least32_t = X__uint32_t /* _types.h:76:21 */
type X__int_least64_t = X__int64_t   /* _types.h:77:20 */
type X__uint_least64_t = X__uint64_t /* _types.h:78:21 */

// 7.18.1.3 Fastest minimum-width integer types
type X__int_fast8_t = X__int32_t    /* _types.h:81:20 */
type X__uint_fast8_t = X__uint32_t  /* _types.h:82:21 */
type X__int_fast16_t = X__int32_t   /* _types.h:83:20 */
type X__uint_fast16_t = X__uint32_t /* _types.h:84:21 */
type X__int_fast32_t = X__int32_t   /* _types.h:85:20 */
type X__uint_fast32_t = X__uint32_t /* _types.h:86:21 */
type X__int_fast64_t = X__int64_t   /* _types.h:87:20 */
type X__uint_fast64_t = X__uint64_t /* _types.h:88:21 */

// 7.18.1.4 Integer types capable of holding object pointers
type X__intptr_t = int64   /* _types.h:103:16 */
type X__uintptr_t = uint64 /* _types.h:104:24 */

// 7.18.1.5 Greatest-width integer types
type X__intmax_t = X__int64_t   /* _types.h:107:20 */
type X__uintmax_t = X__uint64_t /* _types.h:108:21 */

// Register size
type X__register_t = int64 /* _types.h:111:16 */

// VM system types
type X__vaddr_t = uint64 /* _types.h:114:24 */
type X__paddr_t = uint64 /* _types.h:115:24 */
type X__vsize_t = uint64 /* _types.h:116:24 */
type X__psize_t = uint64 /* _types.h:117:24 */

// Standard system types
type X__double_t = float64           /* _types.h:120:18 */
type X__float_t = float32            /* _types.h:121:17 */
type X__ptrdiff_t = int64            /* _types.h:122:16 */
type X__size_t = uint64              /* _types.h:123:24 */
type X__ssize_t = int64              /* _types.h:124:16 */
type X__va_list = X__builtin_va_list /* _types.h:126:27 */

// Wide character support types
type X__wchar_t = int32     /* _types.h:133:15 */
type X__wint_t = int32      /* _types.h:135:15 */
type X__rune_t = int32      /* _types.h:136:15 */
type X__wctrans_t = uintptr /* _types.h:137:14 */
type X__wctype_t = uintptr  /* _types.h:138:14 */

type X__blkcnt_t = X__int64_t    /* _types.h:39:19 */ // blocks allocated for file
type X__blksize_t = X__int32_t   /* _types.h:40:19 */ // optimal blocksize for I/O
type X__clock_t = X__int64_t     /* _types.h:41:19 */ // ticks in CLOCKS_PER_SEC
type X__clockid_t = X__int32_t   /* _types.h:42:19 */ // CLOCK_* identifiers
type X__cpuid_t = uint64         /* _types.h:43:23 */ // CPU id
type X__dev_t = X__int32_t       /* _types.h:44:19 */ // device number
type X__fixpt_t = X__uint32_t    /* _types.h:45:20 */ // fixed point number
type X__fsblkcnt_t = X__uint64_t /* _types.h:46:20 */ // file system block count
type X__fsfilcnt_t = X__uint64_t /* _types.h:47:20 */ // file system file count
type X__gid_t = X__uint32_t      /* _types.h:48:20 */ // group id
type X__id_t = X__uint32_t       /* _types.h:49:20 */ // may contain pid, uid or gid
type X__in_addr_t = X__uint32_t  /* _types.h:50:20 */ // base type for internet address
type X__in_port_t = X__uint16_t  /* _types.h:51:20 */ // IP port type
type X__ino_t = X__uint64_t      /* _types.h:52:20 */ // inode number
type X__key_t = int64            /* _types.h:53:15 */ // IPC key (for Sys V IPC)
type X__mode_t = X__uint32_t     /* _types.h:54:20 */ // permissions
type X__nlink_t = X__uint32_t    /* _types.h:55:20 */ // link count
type X__off_t = X__int64_t       /* _types.h:56:19 */ // file offset or size
type X__pid_t = X__int32_t       /* _types.h:57:19 */ // process id
type X__rlim_t = X__uint64_t     /* _types.h:58:20 */ // resource limit
type X__sa_family_t = X__uint8_t /* _types.h:59:19 */ // sockaddr address family type
type X__segsz_t = X__int32_t     /* _types.h:60:19 */ // segment size
type X__socklen_t = X__uint32_t  /* _types.h:61:20 */ // length type for network syscalls
type X__suseconds_t = int64      /* _types.h:62:15 */ // microseconds (signed)
type X__swblk_t = X__int32_t     /* _types.h:63:19 */ // swap offset
type X__time_t = X__int64_t      /* _types.h:64:19 */ // epoch time
type X__timer_t = X__int32_t     /* _types.h:65:19 */ // POSIX timer identifiers
type X__uid_t = X__uint32_t      /* _types.h:66:20 */ // user id
type X__useconds_t = X__uint32_t /* _types.h:67:20 */ // microseconds

// mbstate_t is an opaque object to keep conversion state, during multibyte
// stream conversions. The content must not be referenced by user programs.
type X__mbstate_t = struct {
	F__ccgo_pad1 [0]uint64
	F__mbstate8  [128]int8
} /* _types.h:76:3 */

// Tell sys/endian.h we have MD variants of the swap macros.

// Note that these macros evaluate their arguments several times.

// Public names

// These are specified to be function-like macros to match the spec

// POSIX names

// original BSD names

// these were exposed here before

// ancient stuff

type U_char = uint8   /* types.h:51:23 */
type U_short = uint16 /* types.h:52:24 */
type U_int = uint32   /* types.h:53:22 */
type U_long = uint64  /* types.h:54:23 */

type Unchar = uint8  /* types.h:56:23 */ // Sys V compatibility
type Ushort = uint16 /* types.h:57:24 */ // Sys V compatibility
type Uint = uint32   /* types.h:58:22 */ // Sys V compatibility
type Ulong = uint64  /* types.h:59:23 */ // Sys V compatibility

type Cpuid_t = X__cpuid_t       /* types.h:61:19 */ // CPU id
type Register_t = X__register_t /* types.h:62:22 */ // register-sized type

// XXX The exact-width bit types should only be exposed if __BSD_VISIBLE
//     but the rest of the includes are not ready for that yet.

type Int8_t = X__int8_t /* types.h:75:19 */

type Uint8_t = X__uint8_t /* types.h:80:20 */

type Int16_t = X__int16_t /* types.h:85:20 */

type Uint16_t = X__uint16_t /* types.h:90:21 */

type Int32_t = X__int32_t /* types.h:95:20 */

type Uint32_t = X__uint32_t /* types.h:100:21 */

type Int64_t = X__int64_t /* types.h:105:20 */

type Uint64_t = X__uint64_t /* types.h:110:21 */

// BSD-style unsigned bits types
type U_int8_t = X__uint8_t   /* types.h:114:19 */
type U_int16_t = X__uint16_t /* types.h:115:20 */
type U_int32_t = X__uint32_t /* types.h:116:20 */
type U_int64_t = X__uint64_t /* types.h:117:20 */

// quads, deprecated in favor of 64 bit int types
type Quad_t = X__int64_t    /* types.h:120:19 */
type U_quad_t = X__uint64_t /* types.h:121:20 */

// VM system types
type Vaddr_t = X__vaddr_t /* types.h:125:19 */
type Paddr_t = X__paddr_t /* types.h:126:19 */
type Vsize_t = X__vsize_t /* types.h:127:19 */
type Psize_t = X__psize_t /* types.h:128:19 */

// Standard system types
type Blkcnt_t = X__blkcnt_t       /* types.h:132:20 */ // blocks allocated for file
type Blksize_t = X__blksize_t     /* types.h:133:21 */ // optimal blocksize for I/O
type Caddr_t = uintptr            /* types.h:134:14 */ // core address
type Daddr32_t = X__int32_t       /* types.h:135:19 */ // 32-bit disk address
type Daddr_t = X__int64_t         /* types.h:136:19 */ // 64-bit disk address
type Dev_t = X__dev_t             /* types.h:137:18 */ // device number
type Fixpt_t = X__fixpt_t         /* types.h:138:19 */ // fixed point number
type Gid_t = X__gid_t             /* types.h:139:18 */ // group id
type Id_t = X__id_t               /* types.h:140:17 */ // may contain pid, uid or gid
type Ino_t = X__ino_t             /* types.h:141:18 */ // inode number
type Key_t = X__key_t             /* types.h:142:18 */ // IPC key (for Sys V IPC)
type Mode_t = X__mode_t           /* types.h:143:18 */ // permissions
type Nlink_t = X__nlink_t         /* types.h:144:19 */ // link count
type Rlim_t = X__rlim_t           /* types.h:145:18 */ // resource limit
type Segsz_t = X__segsz_t         /* types.h:146:19 */ // segment size
type Swblk_t = X__swblk_t         /* types.h:147:19 */ // swap offset
type Uid_t = X__uid_t             /* types.h:148:18 */ // user id
type Useconds_t = X__useconds_t   /* types.h:149:22 */ // microseconds
type Suseconds_t = X__suseconds_t /* types.h:150:23 */ // microseconds (signed)
type Fsblkcnt_t = X__fsblkcnt_t   /* types.h:151:22 */ // file system block count
type Fsfilcnt_t = X__fsfilcnt_t   /* types.h:152:22 */ // file system file count

// The following types may be defined in multiple header files.
type Clock_t = X__clock_t /* types.h:159:19 */

type Clockid_t = X__clockid_t /* types.h:164:21 */

type Pid_t = X__pid_t /* types.h:169:18 */

type Ssize_t = X__ssize_t /* types.h:179:19 */

type Time_t = X__time_t /* types.h:184:18 */

type Timer_t = X__timer_t /* types.h:189:19 */

type Off_t = X__off_t /* types.h:194:18 */

// Major, minor numbers, dev_t's.

//	$OpenBSD: uuid.h,v 1.3 2014/08/31 09:36:39 miod Exp $
//	$NetBSD: uuid.h,v 1.5 2008/11/18 14:01:03 joerg Exp $

// Copyright (c) 2002 Marcel Moolenaar
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR
// IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
// THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// $FreeBSD: /repoman/r/ncvs/src/sys/sys/uuid.h,v 1.3 2003/05/31 16:47:07 phk Exp $

// Length of a node address (an IEEE 802 address).

// Length of a printed UUID.

// See also:
//      http://www.opengroup.org/dce/info/draft-leach-uuids-guids-01.txt
//      http://www.opengroup.org/onlinepubs/009629399/apdxa.htm
//
// A DCE 1.1 compatible source representation of UUIDs.
type Uuid = struct {
	Ftime_low                  Uint32_t
	Ftime_mid                  Uint16_t
	Ftime_hi_and_version       Uint16_t
	Fclock_seq_hi_and_reserved Uint8_t
	Fclock_seq_low             Uint8_t
	Fnode                      [6]Uint8_t
} /* uuid.h:48:1 */

type Uuid_t = Uuid /* uuid.h:71:21 */

var _ int8 /* gen.c:2:13: */
