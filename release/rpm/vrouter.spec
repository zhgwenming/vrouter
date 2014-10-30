%global debug_package %{nil}

Name:		vrouter
Version:	0.1.1
Release:	1%{?dist}
Summary:	vrouter a general tool for distributed docker networking

License:	ASL 2.0
URL:		https://github.com/zhgwenming/%{name}/
Source0:	https://github.com/zhgwenming/%{name}/archive/v%{version}/%{name}-v%{version}.tar.gz
Source1:	vrouter.service
#Source2:	vrouter.socket

BuildRequires:	golang
#BuildRequires:	systemd
#BuildRequires:	golang(github.com/coreos/go-systemd/activation) = 2-1.el7

#Requires(post): systemd
#Requires(preun): systemd
#Requires(postun): systemd

%description
vrouter for distributed docker networking

%prep
%setup -q

%build
make

%install
install -D -p  build/bin/vrouter %{buildroot}%{_bindir}/vrouter
install -d -m 700 %{buildroot}%{_sysconfdir}/%{name}
#install -D -p -m 0644 %{SOURCE1} %{buildroot}%{_unitdir}/%{name}.service
#install -D -p -m 0644 %{SOURCE2} %{buildroot}%{_unitdir}/%{name}.socket

%post
#%systemd_post %{name}.service

%preun
#%systemd_preun %{name}.service

%postun
#%systemd_postun %{name}.service

%files
%{_bindir}/vrouter
%{_sysconfdir}/%{name}
#%{_unitdir}/%{name}.service
#%{_unitdir}/%{name}.socket
%doc README.md

%changelog
* Fri Oct 31 2014 Albert Zhang <zhgwenming@gmail.com> - 0.1
- initial version

