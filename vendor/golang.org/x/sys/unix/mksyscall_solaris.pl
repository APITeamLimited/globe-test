#!/usr/bin/env perl
# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# This program reads a file containing function prototypes
# (like syscall_solaris.go) and generates system call bodies.
# The prototypes are marked by lines beginning with "//sys"
# and read like func declarations if //sys is replaced by func, but:
#	* The parameter lists must give a name for each argument.
#	  This includes return parameters.
#	* The parameter lists must give a type for each argument:
#	  the (x, y, z int) shorthand is not allowed.
#	* If the return parameter is an error number, it must be named err.
#	* If go func name needs to be different than its libc name,
#	* or the function is not in libc, name could be specified
#	* at the end, after "=" sign, like
#	  //sys getsockopt(s int, level int, name int, val uintptr, vallen *_Socklen) (err error) = libsocket.getsockopt

use strict;

my $cmdline = "mksyscall_solaris.pl " . join(' ', @ARGV);
my $errors = 0;
my $_32bit = "";
my $tags = "";  # build tags

binmode STDOUT;

if($ARGV[0] eq "-b32") ***REMOVED***
	$_32bit = "big-endian";
	shift;
***REMOVED*** elsif($ARGV[0] eq "-l32") ***REMOVED***
	$_32bit = "little-endian";
	shift;
***REMOVED***
if($ARGV[0] eq "-tags") ***REMOVED***
	shift;
	$tags = $ARGV[0];
	shift;
***REMOVED***

if($ARGV[0] =~ /^-/) ***REMOVED***
	print STDERR "usage: mksyscall_solaris.pl [-b32 | -l32] [-tags x,y] [file ...]\n";
	exit 1;
***REMOVED***

sub parseparamlist($) ***REMOVED***
	my ($list) = @_;
	$list =~ s/^\s*//;
	$list =~ s/\s*$//;
	if($list eq "") ***REMOVED***
		return ();
	***REMOVED***
	return split(/\s*,\s*/, $list);
***REMOVED***

sub parseparam($) ***REMOVED***
	my ($p) = @_;
	if($p !~ /^(\S*) (\S*)$/) ***REMOVED***
		print STDERR "$ARGV:$.: malformed parameter: $p\n";
		$errors = 1;
		return ("xx", "int");
	***REMOVED***
	return ($1, $2);
***REMOVED***

my $package = "";
my $text = "";
my $dynimports = "";
my $linknames = "";
my @vars = ();
while(<>) ***REMOVED***
	chomp;
	s/\s+/ /g;
	s/^\s+//;
	s/\s+$//;
	$package = $1 if !$package && /^package (\S+)$/;
	my $nonblock = /^\/\/sysnb /;
	next if !/^\/\/sys / && !$nonblock;

	# Line must be of the form
	#	func Open(path string, mode int, perm int) (fd int, err error)
	# Split into name, in params, out params.
	if(!/^\/\/sys(nb)? (\w+)\(([^()]*)\)\s*(?:\(([^()]+)\))?\s*(?:=\s*(?:(\w*)\.)?(\w*))?$/) ***REMOVED***
		print STDERR "$ARGV:$.: malformed //sys declaration\n";
		$errors = 1;
		next;
	***REMOVED***
	my ($nb, $func, $in, $out, $modname, $sysname) = ($1, $2, $3, $4, $5, $6);

	# Split argument lists on comma.
	my @in = parseparamlist($in);
	my @out = parseparamlist($out);

	# So file name.
	if($modname eq "") ***REMOVED***
		$modname = "libc";
	***REMOVED***

	# System call name.
	if($sysname eq "") ***REMOVED***
		$sysname = "$func";
	***REMOVED***

	# System call pointer variable name.
	my $sysvarname = "proc$sysname";

	my $strconvfunc = "BytePtrFromString";
	my $strconvtype = "*byte";

	$sysname =~ y/A-Z/a-z/; # All libc functions are lowercase.

	# Runtime import of function to allow cross-platform builds.
	$dynimports .= "//go:cgo_import_dynamic libc_$***REMOVED***sysname***REMOVED*** $***REMOVED***sysname***REMOVED*** \"$modname.so\"\n";
	# Link symbol to proc address variable.
	$linknames .= "//go:linkname $***REMOVED***sysvarname***REMOVED*** libc_$***REMOVED***sysname***REMOVED***\n";
	# Library proc address variable.
	push @vars, $sysvarname;

	# Go function header.
	$out = join(', ', @out);
	if($out ne "") ***REMOVED***
		$out = " ($out)";
	***REMOVED***
	if($text ne "") ***REMOVED***
		$text .= "\n"
	***REMOVED***
	$text .= sprintf "func %s(%s)%s ***REMOVED***\n", $func, join(', ', @in), $out;

	# Check if err return available
	my $errvar = "";
	foreach my $p (@out) ***REMOVED***
		my ($name, $type) = parseparam($p);
		if($type eq "error") ***REMOVED***
			$errvar = $name;
			last;
		***REMOVED***
	***REMOVED***

	# Prepare arguments to Syscall.
	my @args = ();
	my $n = 0;
	foreach my $p (@in) ***REMOVED***
		my ($name, $type) = parseparam($p);
		if($type =~ /^\*/) ***REMOVED***
			push @args, "uintptr(unsafe.Pointer($name))";
		***REMOVED*** elsif($type eq "string" && $errvar ne "") ***REMOVED***
			$text .= "\tvar _p$n $strconvtype\n";
			$text .= "\t_p$n, $errvar = $strconvfunc($name)\n";
			$text .= "\tif $errvar != nil ***REMOVED***\n\t\treturn\n\t***REMOVED***\n";
			push @args, "uintptr(unsafe.Pointer(_p$n))";
			$n++;
		***REMOVED*** elsif($type eq "string") ***REMOVED***
			print STDERR "$ARGV:$.: $func uses string arguments, but has no error return\n";
			$text .= "\tvar _p$n $strconvtype\n";
			$text .= "\t_p$n, _ = $strconvfunc($name)\n";
			push @args, "uintptr(unsafe.Pointer(_p$n))";
			$n++;
		***REMOVED*** elsif($type =~ /^\[\](.*)/) ***REMOVED***
			# Convert slice into pointer, length.
			# Have to be careful not to take address of &a[0] if len == 0:
			# pass nil in that case.
			$text .= "\tvar _p$n *$1\n";
			$text .= "\tif len($name) > 0 ***REMOVED***\n\t\t_p$n = \&$name\[0]\n\t***REMOVED***\n";
			push @args, "uintptr(unsafe.Pointer(_p$n))", "uintptr(len($name))";
			$n++;
		***REMOVED*** elsif($type eq "int64" && $_32bit ne "") ***REMOVED***
			if($_32bit eq "big-endian") ***REMOVED***
				push @args, "uintptr($name >> 32)", "uintptr($name)";
			***REMOVED*** else ***REMOVED***
				push @args, "uintptr($name)", "uintptr($name >> 32)";
			***REMOVED***
		***REMOVED*** elsif($type eq "bool") ***REMOVED***
 			$text .= "\tvar _p$n uint32\n";
			$text .= "\tif $name ***REMOVED***\n\t\t_p$n = 1\n\t***REMOVED*** else ***REMOVED***\n\t\t_p$n = 0\n\t***REMOVED***\n";
			push @args, "uintptr(_p$n)";
			$n++;
		***REMOVED*** else ***REMOVED***
			push @args, "uintptr($name)";
		***REMOVED***
	***REMOVED***
	my $nargs = @args;

	# Determine which form to use; pad args with zeros.
	my $asm = "sysvicall6";
	if ($nonblock) ***REMOVED***
		$asm = "rawSysvicall6";
	***REMOVED***
	if(@args <= 6) ***REMOVED***
		while(@args < 6) ***REMOVED***
			push @args, "0";
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		print STDERR "$ARGV:$.: too many arguments to system call\n";
	***REMOVED***

	# Actual call.
	my $args = join(', ', @args);
	my $call = "$asm(uintptr(unsafe.Pointer(&$sysvarname)), $nargs, $args)";

	# Assign return values.
	my $body = "";
	my $failexpr = "";
	my @ret = ("_", "_", "_");
	my @pout= ();
	my $do_errno = 0;
	for(my $i=0; $i<@out; $i++) ***REMOVED***
		my $p = $out[$i];
		my ($name, $type) = parseparam($p);
		my $reg = "";
		if($name eq "err") ***REMOVED***
			$reg = "e1";
			$ret[2] = $reg;
			$do_errno = 1;
		***REMOVED*** else ***REMOVED***
			$reg = sprintf("r%d", $i);
			$ret[$i] = $reg;
		***REMOVED***
		if($type eq "bool") ***REMOVED***
			$reg = "$reg != 0";
		***REMOVED***
		if($type eq "int64" && $_32bit ne "") ***REMOVED***
			# 64-bit number in r1:r0 or r0:r1.
			if($i+2 > @out) ***REMOVED***
				print STDERR "$ARGV:$.: not enough registers for int64 return\n";
			***REMOVED***
			if($_32bit eq "big-endian") ***REMOVED***
				$reg = sprintf("int64(r%d)<<32 | int64(r%d)", $i, $i+1);
			***REMOVED*** else ***REMOVED***
				$reg = sprintf("int64(r%d)<<32 | int64(r%d)", $i+1, $i);
			***REMOVED***
			$ret[$i] = sprintf("r%d", $i);
			$ret[$i+1] = sprintf("r%d", $i+1);
		***REMOVED***
		if($reg ne "e1") ***REMOVED***
			$body .= "\t$name = $type($reg)\n";
		***REMOVED***
	***REMOVED***
	if ($ret[0] eq "_" && $ret[1] eq "_" && $ret[2] eq "_") ***REMOVED***
		$text .= "\t$call\n";
	***REMOVED*** else ***REMOVED***
		$text .= "\t$ret[0], $ret[1], $ret[2] := $call\n";
	***REMOVED***
	$text .= $body;

	if ($do_errno) ***REMOVED***
		$text .= "\tif e1 != 0 ***REMOVED***\n";
		$text .= "\t\terr = e1\n";
		$text .= "\t***REMOVED***\n";
	***REMOVED***
	$text .= "\treturn\n";
	$text .= "***REMOVED***\n";
***REMOVED***

if($errors) ***REMOVED***
	exit 1;
***REMOVED***

print <<EOF;
// $cmdline
// Code generated by the command above; see README.md. DO NOT EDIT.

// +build $tags

package $package

import (
	"syscall"
	"unsafe"
)
EOF

print "import \"golang.org/x/sys/unix\"\n" if $package ne "unix";

my $vardecls = "\t" . join(",\n\t", @vars);
$vardecls .= " syscallFunc";

chomp($_=<<EOF);

$dynimports
$linknames
var (
$vardecls
)

$text
EOF
print $_;
exit 0;
