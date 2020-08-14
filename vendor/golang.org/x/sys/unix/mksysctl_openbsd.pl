#!/usr/bin/env perl

# Copyright 2011 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

#
# Parse the header files for OpenBSD and generate a Go usable sysctl MIB.
#
# Build a MIB with each entry being an array containing the level, type and
# a hash that will contain additional entries if the current entry is a node.
# We then walk this MIB and create a flattened sysctl name to OID hash.
#

use strict;

if($ENV***REMOVED***'GOARCH'***REMOVED*** eq "" || $ENV***REMOVED***'GOOS'***REMOVED*** eq "") ***REMOVED***
	print STDERR "GOARCH or GOOS not defined in environment\n";
	exit 1;
***REMOVED***

my $debug = 0;
my %ctls = ();

my @headers = qw (
	sys/sysctl.h
	sys/socket.h
	sys/tty.h
	sys/malloc.h
	sys/mount.h
	sys/namei.h
	sys/sem.h
	sys/shm.h
	sys/vmmeter.h
	uvm/uvm_param.h
	uvm/uvm_swap_encrypt.h
	ddb/db_var.h
	net/if.h
	net/if_pfsync.h
	net/pipex.h
	netinet/in.h
	netinet/icmp_var.h
	netinet/igmp_var.h
	netinet/ip_ah.h
	netinet/ip_carp.h
	netinet/ip_divert.h
	netinet/ip_esp.h
	netinet/ip_ether.h
	netinet/ip_gre.h
	netinet/ip_ipcomp.h
	netinet/ip_ipip.h
	netinet/pim_var.h
	netinet/tcp_var.h
	netinet/udp_var.h
	netinet6/in6.h
	netinet6/ip6_divert.h
	netinet6/pim6_var.h
	netinet/icmp6.h
	netmpls/mpls.h
);

my @ctls = qw (
	kern
	vm
	fs
	net
	#debug				# Special handling required
	hw
	#machdep			# Arch specific
	user
	ddb
	#vfs				# Special handling required
	fs.posix
	kern.forkstat
	kern.intrcnt
	kern.malloc
	kern.nchstats
	kern.seminfo
	kern.shminfo
	kern.timecounter
	kern.tty
	kern.watchdog
	net.bpf
	net.ifq
	net.inet
	net.inet.ah
	net.inet.carp
	net.inet.divert
	net.inet.esp
	net.inet.etherip
	net.inet.gre
	net.inet.icmp
	net.inet.igmp
	net.inet.ip
	net.inet.ip.ifq
	net.inet.ipcomp
	net.inet.ipip
	net.inet.mobileip
	net.inet.pfsync
	net.inet.pim
	net.inet.tcp
	net.inet.udp
	net.inet6
	net.inet6.divert
	net.inet6.ip6
	net.inet6.icmp6
	net.inet6.pim6
	net.inet6.tcp6
	net.inet6.udp6
	net.mpls
	net.mpls.ifq
	net.key
	net.pflow
	net.pfsync
	net.pipex
	net.rt
	vm.swapencrypt
	#vfsgenctl			# Special handling required
);

# Node name "fixups"
my %ctl_map = (
	"ipproto" => "net.inet",
	"net.inet.ipproto" => "net.inet",
	"net.inet6.ipv6proto" => "net.inet6",
	"net.inet6.ipv6" => "net.inet6.ip6",
	"net.inet.icmpv6" => "net.inet6.icmp6",
	"net.inet6.divert6" => "net.inet6.divert",
	"net.inet6.tcp6" => "net.inet.tcp",
	"net.inet6.udp6" => "net.inet.udp",
	"mpls" => "net.mpls",
	"swpenc" => "vm.swapencrypt"
);

# Node mappings
my %node_map = (
	"net.inet.ip.ifq" => "net.ifq",
	"net.inet.pfsync" => "net.pfsync",
	"net.mpls.ifq" => "net.ifq"
);

my $ctlname;
my %mib = ();
my %sysctl = ();
my $node;

sub debug() ***REMOVED***
	print STDERR "$_[0]\n" if $debug;
***REMOVED***

# Walk the MIB and build a sysctl name to OID mapping.
sub build_sysctl() ***REMOVED***
	my ($node, $name, $oid) = @_;
	my %node = %***REMOVED***$node***REMOVED***;
	my @oid = @***REMOVED***$oid***REMOVED***;

	foreach my $key (sort keys %node) ***REMOVED***
		my @node = @***REMOVED***$node***REMOVED***$key***REMOVED******REMOVED***;
		my $nodename = $name.($name ne '' ? '.' : '').$key;
		my @nodeoid = (@oid, $node[0]);
		if ($node[1] eq 'CTLTYPE_NODE') ***REMOVED***
			if (exists $node_map***REMOVED***$nodename***REMOVED***) ***REMOVED***
				$node = \%mib;
				$ctlname = $node_map***REMOVED***$nodename***REMOVED***;
				foreach my $part (split /\./, $ctlname) ***REMOVED***
					$node = \%***REMOVED***@***REMOVED***$$node***REMOVED***$part***REMOVED******REMOVED***[2]***REMOVED***;
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				$node = $node[2];
			***REMOVED***
			&build_sysctl($node, $nodename, \@nodeoid);
		***REMOVED*** elsif ($node[1] ne '') ***REMOVED***
			$sysctl***REMOVED***$nodename***REMOVED*** = \@nodeoid;
		***REMOVED***
	***REMOVED***
***REMOVED***

foreach my $ctl (@ctls) ***REMOVED***
	$ctls***REMOVED***$ctl***REMOVED*** = $ctl;
***REMOVED***

# Build MIB
foreach my $header (@headers) ***REMOVED***
	&debug("Processing $header...");
	open HEADER, "/usr/include/$header" ||
	    print STDERR "Failed to open $header\n";
	while (<HEADER>) ***REMOVED***
		if ($_ =~ /^#define\s+(CTL_NAMES)\s+***REMOVED***/ ||
		    $_ =~ /^#define\s+(CTL_(.*)_NAMES)\s+***REMOVED***/ ||
		    $_ =~ /^#define\s+((.*)CTL_NAMES)\s+***REMOVED***/) ***REMOVED***
			if ($1 eq 'CTL_NAMES') ***REMOVED***
				# Top level.
				$node = \%mib;
			***REMOVED*** else ***REMOVED***
				# Node.
				my $nodename = lc($2);
				if ($header =~ /^netinet\//) ***REMOVED***
					$ctlname = "net.inet.$nodename";
				***REMOVED*** elsif ($header =~ /^netinet6\//) ***REMOVED***
					$ctlname = "net.inet6.$nodename";
				***REMOVED*** elsif ($header =~ /^net\//) ***REMOVED***
					$ctlname = "net.$nodename";
				***REMOVED*** else ***REMOVED***
					$ctlname = "$nodename";
					$ctlname =~ s/^(fs|net|kern)_/$1\./;
				***REMOVED***
				if (exists $ctl_map***REMOVED***$ctlname***REMOVED***) ***REMOVED***
					$ctlname = $ctl_map***REMOVED***$ctlname***REMOVED***;
				***REMOVED***
				if (not exists $ctls***REMOVED***$ctlname***REMOVED***) ***REMOVED***
					&debug("Ignoring $ctlname...");
					next;
				***REMOVED***

				# Walk down from the top of the MIB.
				$node = \%mib;
				foreach my $part (split /\./, $ctlname) ***REMOVED***
					if (not exists $$node***REMOVED***$part***REMOVED***) ***REMOVED***
						&debug("Missing node $part");
						$$node***REMOVED***$part***REMOVED*** = [ 0, '', ***REMOVED******REMOVED*** ];
					***REMOVED***
					$node = \%***REMOVED***@***REMOVED***$$node***REMOVED***$part***REMOVED******REMOVED***[2]***REMOVED***;
				***REMOVED***
			***REMOVED***

			# Populate current node with entries.
			my $i = -1;
			while (defined($_) && $_ !~ /^***REMOVED***/) ***REMOVED***
				$_ = <HEADER>;
				$i++ if $_ =~ /***REMOVED***.****REMOVED***/;
				next if $_ !~ /***REMOVED***\s+"(\w+)",\s+(CTLTYPE_[A-Z]+)\s+***REMOVED***/;
				$$node***REMOVED***$1***REMOVED*** = [ $i, $2, ***REMOVED******REMOVED*** ];
			***REMOVED***
		***REMOVED***
	***REMOVED***
	close HEADER;
***REMOVED***

&build_sysctl(\%mib, "", []);

print <<EOF;
// mksysctl_openbsd.pl
// MACHINE GENERATED BY THE ABOVE COMMAND; DO NOT EDIT

// +build $ENV***REMOVED***'GOARCH'***REMOVED***,$ENV***REMOVED***'GOOS'***REMOVED***

package unix;

type mibentry struct ***REMOVED***
	ctlname string
	ctloid []_C_int
***REMOVED***

var sysctlMib = []mibentry ***REMOVED***
EOF

foreach my $name (sort keys %sysctl) ***REMOVED***
	my @oid = @***REMOVED***$sysctl***REMOVED***$name***REMOVED******REMOVED***;
	print "\t***REMOVED*** \"$name\", []_C_int***REMOVED*** ", join(', ', @oid), " ***REMOVED*** ***REMOVED***, \n";
***REMOVED***

print <<EOF;
***REMOVED***
EOF
