Name:           k6-rpm
Version:        0.0.1
Release:        1
Summary:        k6 RPM Repository Configuration
Group:          System Environment/Base
License:        AGPL-3.0
URL:            https://dl.k6.io
Source0:        RPM-GPG-KEY-k6-io
Source1:        k6-io.repo
BuildRoot:      %***REMOVED***_builddir***REMOVED***/%***REMOVED***name***REMOVED***-%***REMOVED***version***REMOVED***-rpmroot
BuildArch:      noarch

%description
This package installs the repository GPG and repo files for the k6 software repository.

%prep
%setup -c -T

%build

%install
rm -rf $RPM_BUILD_ROOT

# gpg
install -Dpm 644 %***REMOVED***SOURCE0***REMOVED*** $RPM_BUILD_ROOT%***REMOVED***_sysconfdir***REMOVED***/pki/rpm-gpg/RPM-GPG-KEY-k6-io

# yum
install -dm 755 $RPM_BUILD_ROOT%***REMOVED***_sysconfdir***REMOVED***/yum.repos.d
install -pm 644 %***REMOVED***SOURCE1***REMOVED*** $RPM_BUILD_ROOT%***REMOVED***_sysconfdir***REMOVED***/yum.repos.d

%clean
rm -rf $RPM_BUILD_ROOT

%files
%defattr(-,root,root,-)
%***REMOVED***_sysconfdir***REMOVED***/pki/rpm-gpg/RPM-GPG-KEY-k6-io
%config %***REMOVED***_sysconfdir***REMOVED***/yum.repos.d/k6-io.repo
