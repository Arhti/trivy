package rpm

import (
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aquasecurity/trivy/pkg/fanal/types"
)

func TestParseRpmInfo(t *testing.T) {
	var tests = map[string]struct {
		path string
		pkgs []types.Package
	}{
		"Valid": {
			path: "./testdata/valid",
			// cp ./testdata/valid /path/to/testdir/Packages
			// rpm --dbpath /path/to/testdir -qa --qf "{Name: \"%{NAME}\", Epoch: %{EPOCHNUM}, Version: \"%{VERSION}\", Release: \"%{RELEASE}\", Arch: \"%{ARCH}\"\},\n"
			pkgs: []types.Package{
				{
					Name: "centos-release", Epoch: 0, Version: "7", Release: "1.1503.el7.centos.2.8", Arch: "x86_64",
					SrcName: "centos-release", SrcEpoch: 0, SrcVersion: "7", SrcRelease: "1.1503.el7.centos.2.8",
					Licenses: []string{"GPLv2"},
				},
				{
					Name: "filesystem", Epoch: 0, Version: "3.2", Release: "18.el7", Arch: "x86_64",
					SrcName: "filesystem", SrcEpoch: 0, SrcVersion: "3.2", SrcRelease: "18.el7",
					Licenses: []string{"Public Domain"},
				},
			},
		},
		"ValidBig": {
			path: "./testdata/valid_big",
			// $ cat rpmqa.py
			// import rpm
			// from rpmUtils.miscutils import splitFilename
			//
			//
			// rpm.addMacro('_dbpath', '/tmp/')
			// ts = rpm.TransactionSet()
			// mi = ts.dbMatch()
			// for h in mi:
			//     sname = sversion = srelease = ""
			//     if h[rpm.RPMTAG_SOURCERPM] != "(none)":
			//         sname, sversion, srelease, _, _ = splitFilename(h[rpm.RPMTAG_SOURCERPM])
			//     print "{Name: \"%s\", Epoch: %d, Version: \"%s\", Release: \"%s\", Arch: \"%s\", SrcName: \"%s\", SrcEpoch: %d, SrcVersion: \"%s\", SrcRelease: \"%s\"}," % (
			//         h[rpm.RPMTAG_NAME], h[rpm.RPMTAG_EPOCHNUM], h[rpm.RPMTAG_VERSION], h[rpm.RPMTAG_RELEASE], h[rpm.RPMTAG_ARCH],
			//         sname, h[rpm.RPMTAG_EPOCHNUM], sversion, srelease)
			pkgs: []types.Package{
				{
					Name: "publicsuffix-list-dafsa", Epoch: 0, Version: "20180514", Release: "1.fc28", Arch: "noarch",
					SrcName: "publicsuffix-list", SrcEpoch: 0, SrcVersion: "20180514", SrcRelease: "1.fc28",
					Modularitylabel: "", Licenses: []string{"MPLv2.0"},
				},
				{
					Name: "libreport-filesystem", Epoch: 0, Version: "2.9.5", Release: "1.fc28", Arch: "x86_64",
					SrcName: "libreport", SrcEpoch: 0, SrcVersion: "2.9.5", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+"},
				},
				{
					Name: "fedora-gpg-keys", Epoch: 0, Version: "28", Release: "5", Arch: "noarch",
					SrcName: "fedora-repos", SrcEpoch: 0, SrcVersion: "28", SrcRelease: "5", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "fedora-release", Epoch: 0, Version: "28", Release: "2", Arch: "noarch",
					SrcName: "fedora-release", SrcEpoch: 0, SrcVersion: "28", SrcRelease: "2", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "filesystem", Epoch: 0, Version: "3.8", Release: "2.fc28", Arch: "x86_64",
					SrcName: "filesystem", SrcEpoch: 0, SrcVersion: "3.8", SrcRelease: "2.fc28", Modularitylabel: "",
					Licenses: []string{"Public Domain"},
				},
				{
					Name: "tzdata", Epoch: 0, Version: "2018e", Release: "1.fc28", Arch: "noarch", SrcName: "tzdata",
					SrcEpoch: 0, SrcVersion: "2018e", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"Public Domain"},
				},
				{
					Name: "pcre2", Epoch: 0, Version: "10.31", Release: "10.fc28", Arch: "x86_64", SrcName: "pcre2",
					SrcEpoch: 0, SrcVersion: "10.31", SrcRelease: "10.fc28", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "glibc-minimal-langpack", Epoch: 0, Version: "2.27", Release: "32.fc28", Arch: "x86_64",
					SrcName: "glibc", SrcEpoch: 0, SrcVersion: "2.27", SrcRelease: "32.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+ and LGPLv2+ with exceptions and GPLv2+ and GPLv2+ with exceptions and BSD and Inner-Net and ISC and Public Domain and GFDL"},
				},
				{
					Name: "glibc-common", Epoch: 0, Version: "2.27", Release: "32.fc28", Arch: "x86_64",
					SrcName: "glibc", SrcEpoch: 0, SrcVersion: "2.27", SrcRelease: "32.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+ and LGPLv2+ with exceptions and GPLv2+ and GPLv2+ with exceptions and BSD and Inner-Net and ISC and Public Domain and GFDL"},
				},
				{
					Name: "bash", Epoch: 0, Version: "4.4.23", Release: "1.fc28", Arch: "x86_64", SrcName: "bash",
					SrcEpoch: 0, SrcVersion: "4.4.23", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "zlib", Epoch: 0, Version: "1.2.11", Release: "8.fc28", Arch: "x86_64", SrcName: "zlib",
					SrcEpoch: 0, SrcVersion: "1.2.11", SrcRelease: "8.fc28", Modularitylabel: "",
					Licenses: []string{"zlib and Boost"},
				},
				{
					Name: "bzip2-libs", Epoch: 0, Version: "1.0.6", Release: "26.fc28", Arch: "x86_64",
					SrcName: "bzip2", SrcEpoch: 0, SrcVersion: "1.0.6", SrcRelease: "26.fc28", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "libcap", Epoch: 0, Version: "2.25", Release: "9.fc28", Arch: "x86_64", SrcName: "libcap",
					SrcEpoch: 0, SrcVersion: "2.25", SrcRelease: "9.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2"},
				},
				{
					Name: "libgpg-error", Epoch: 0, Version: "1.31", Release: "1.fc28", Arch: "x86_64",
					SrcName: "libgpg-error", SrcEpoch: 0, SrcVersion: "1.31", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "libzstd", Epoch: 0, Version: "1.3.5", Release: "1.fc28", Arch: "x86_64", SrcName: "zstd",
					SrcEpoch: 0, SrcVersion: "1.3.5", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"BSD and GPLv2"},
				},
				{
					Name: "expat", Epoch: 0, Version: "2.2.5", Release: "3.fc28", Arch: "x86_64", SrcName: "expat",
					SrcEpoch: 0, SrcVersion: "2.2.5", SrcRelease: "3.fc28", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "nss-util", Epoch: 0, Version: "3.38.0", Release: "1.0.fc28", Arch: "x86_64",
					SrcName: "nss-util", SrcEpoch: 0, SrcVersion: "3.38.0", SrcRelease: "1.0.fc28", Modularitylabel: "",
					Licenses: []string{"MPLv2.0"},
				},
				{
					Name: "libcom_err", Epoch: 0, Version: "1.44.2", Release: "0.fc28", Arch: "x86_64",
					SrcName: "e2fsprogs", SrcEpoch: 0, SrcVersion: "1.44.2", SrcRelease: "0.fc28", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "libffi", Epoch: 0, Version: "3.1", Release: "16.fc28", Arch: "x86_64", SrcName: "libffi",
					SrcEpoch: 0, SrcVersion: "3.1", SrcRelease: "16.fc28", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "libgcrypt", Epoch: 0, Version: "1.8.3", Release: "1.fc28", Arch: "x86_64",
					SrcName: "libgcrypt", SrcEpoch: 0, SrcVersion: "1.8.3", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "libxml2", Epoch: 0, Version: "2.9.8", Release: "4.fc28", Arch: "x86_64", SrcName: "libxml2",
					SrcEpoch: 0, SrcVersion: "2.9.8", SrcRelease: "4.fc28", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "libacl", Epoch: 0, Version: "2.2.53", Release: "1.fc28", Arch: "x86_64", SrcName: "acl",
					SrcEpoch: 0, SrcVersion: "2.2.53", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "sed", Epoch: 0, Version: "4.5", Release: "1.fc28", Arch: "x86_64", SrcName: "sed",
					SrcEpoch: 0, SrcVersion: "4.5", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "libmount", Epoch: 0, Version: "2.32.1", Release: "1.fc28", Arch: "x86_64",
					SrcName: "util-linux", SrcEpoch: 0, SrcVersion: "2.32.1", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "p11-kit", Epoch: 0, Version: "0.23.12", Release: "1.fc28", Arch: "x86_64",
					SrcName: "p11-kit", SrcEpoch: 0, SrcVersion: "0.23.12", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "libidn2", Epoch: 0, Version: "2.0.5", Release: "1.fc28", Arch: "x86_64", SrcName: "libidn2",
					SrcEpoch: 0, SrcVersion: "2.0.5", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"(GPLv2+ or LGPLv3+) and GPLv3+"},
				},
				{
					Name: "libcap-ng", Epoch: 0, Version: "0.7.9", Release: "4.fc28", Arch: "x86_64",
					SrcName: "libcap-ng", SrcEpoch: 0, SrcVersion: "0.7.9", SrcRelease: "4.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "lz4-libs", Epoch: 0, Version: "1.8.1.2", Release: "4.fc28", Arch: "x86_64", SrcName: "lz4",
					SrcEpoch: 0, SrcVersion: "1.8.1.2", SrcRelease: "4.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and BSD"},
				},
				{
					Name: "libassuan", Epoch: 0, Version: "2.5.1", Release: "3.fc28", Arch: "x86_64",
					SrcName: "libassuan", SrcEpoch: 0, SrcVersion: "2.5.1", SrcRelease: "3.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+ and GPLv3+"},
				},
				{
					Name: "keyutils-libs", Epoch: 0, Version: "1.5.10", Release: "6.fc28", Arch: "x86_64",
					SrcName: "keyutils", SrcEpoch: 0, SrcVersion: "1.5.10", SrcRelease: "6.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and LGPLv2+"},
				},
				{
					Name: "glib2", Epoch: 0, Version: "2.56.1", Release: "4.fc28", Arch: "x86_64", SrcName: "glib2",
					SrcEpoch: 0, SrcVersion: "2.56.1", SrcRelease: "4.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "systemd-libs", Epoch: 0, Version: "238", Release: "9.git0e0aa59.fc28", Arch: "x86_64",
					SrcName: "systemd", SrcEpoch: 0, SrcVersion: "238", SrcRelease: "9.git0e0aa59.fc28",
					Modularitylabel: "", Licenses: []string{"LGPLv2+ and MIT"},
				},
				{
					Name: "dbus-libs", Epoch: 1, Version: "1.12.10", Release: "1.fc28", Arch: "x86_64", SrcName: "dbus",
					SrcEpoch: 1, SrcVersion: "1.12.10", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"(GPLv2+ or AFL) and GPLv2+"},
				},
				{
					Name: "libtasn1", Epoch: 0, Version: "4.13", Release: "2.fc28", Arch: "x86_64", SrcName: "libtasn1",
					SrcEpoch: 0, SrcVersion: "4.13", SrcRelease: "2.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv3+ and LGPLv2+"},
				},
				{
					Name: "ca-certificates", Epoch: 0, Version: "2018.2.24", Release: "1.0.fc28", Arch: "noarch",
					SrcName: "ca-certificates", SrcEpoch: 0, SrcVersion: "2018.2.24", SrcRelease: "1.0.fc28",
					Modularitylabel: "", Licenses: []string{"Public Domain"},
				},
				{
					Name: "libarchive", Epoch: 0, Version: "3.3.1", Release: "4.fc28", Arch: "x86_64",
					SrcName: "libarchive", SrcEpoch: 0, SrcVersion: "3.3.1", SrcRelease: "4.fc28", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "openssl", Epoch: 1, Version: "1.1.0h", Release: "3.fc28", Arch: "x86_64", SrcName: "openssl",
					SrcEpoch: 1, SrcVersion: "1.1.0h", SrcRelease: "3.fc28", Modularitylabel: "",
					Licenses: []string{"OpenSSL"},
				},
				{
					Name: "libusbx", Epoch: 0, Version: "1.0.22", Release: "1.fc28", Arch: "x86_64", SrcName: "libusbx",
					SrcEpoch: 0, SrcVersion: "1.0.22", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "libsemanage", Epoch: 0, Version: "2.8", Release: "2.fc28", Arch: "x86_64",
					SrcName: "libsemanage", SrcEpoch: 0, SrcVersion: "2.8", SrcRelease: "2.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "libutempter", Epoch: 0, Version: "1.1.6", Release: "14.fc28", Arch: "x86_64",
					SrcName: "libutempter", SrcEpoch: 0, SrcVersion: "1.1.6", SrcRelease: "14.fc28",
					Modularitylabel: "", Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "mpfr", Epoch: 0, Version: "3.1.6", Release: "1.fc28", Arch: "x86_64", SrcName: "mpfr",
					SrcEpoch: 0, SrcVersion: "3.1.6", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv3+ and GPLv3+ and GFDL"},
				},
				{
					Name: "gnutls", Epoch: 0, Version: "3.6.3", Release: "4.fc28", Arch: "x86_64", SrcName: "gnutls",
					SrcEpoch: 0, SrcVersion: "3.6.3", SrcRelease: "4.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv3+ and LGPLv2+"},
				},
				{
					Name: "gzip", Epoch: 0, Version: "1.9", Release: "3.fc28", Arch: "x86_64", SrcName: "gzip",
					SrcEpoch: 0, SrcVersion: "1.9", SrcRelease: "3.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv3+ and GFDL"},
				},
				{
					Name: "acl", Epoch: 0, Version: "2.2.53", Release: "1.fc28", Arch: "x86_64", SrcName: "acl",
					SrcEpoch: 0, SrcVersion: "2.2.53", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+"},
				},
				{
					Name: "nss-softokn-freebl", Epoch: 0, Version: "3.38.0", Release: "1.0.fc28", Arch: "x86_64",
					SrcName: "nss-softokn", SrcEpoch: 0, SrcVersion: "3.38.0", SrcRelease: "1.0.fc28",
					Modularitylabel: "", Licenses: []string{"MPLv2.0"},
				},
				{
					Name: "nss", Epoch: 0, Version: "3.38.0", Release: "1.0.fc28", Arch: "x86_64", SrcName: "nss",
					SrcEpoch: 0, SrcVersion: "3.38.0", SrcRelease: "1.0.fc28", Modularitylabel: "",
					Licenses: []string{"MPLv2.0"},
				},
				{
					Name: "libmetalink", Epoch: 0, Version: "0.1.3", Release: "6.fc28", Arch: "x86_64",
					SrcName: "libmetalink", SrcEpoch: 0, SrcVersion: "0.1.3", SrcRelease: "6.fc28", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "libdb-utils", Epoch: 0, Version: "5.3.28", Release: "30.fc28", Arch: "x86_64",
					SrcName: "libdb", SrcEpoch: 0, SrcVersion: "5.3.28", SrcRelease: "30.fc28", Modularitylabel: "",
					Licenses: []string{"BSD and LGPLv2 and Sleepycat"},
				},
				{
					Name: "file-libs", Epoch: 0, Version: "5.33", Release: "7.fc28", Arch: "x86_64", SrcName: "file",
					SrcEpoch: 0, SrcVersion: "5.33", SrcRelease: "7.fc28", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "libsss_idmap", Epoch: 0, Version: "1.16.3", Release: "2.fc28", Arch: "x86_64",
					SrcName: "sssd", SrcEpoch: 0, SrcVersion: "1.16.3", SrcRelease: "2.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv3+"},
				},
				{
					Name: "libsigsegv", Epoch: 0, Version: "2.11", Release: "5.fc28", Arch: "x86_64",
					SrcName: "libsigsegv", SrcEpoch: 0, SrcVersion: "2.11", SrcRelease: "5.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+"},
				},
				{
					Name: "krb5-libs", Epoch: 0, Version: "1.16.1", Release: "13.fc28", Arch: "x86_64", SrcName: "krb5",
					SrcEpoch: 0, SrcVersion: "1.16.1", SrcRelease: "13.fc28", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "libnsl2", Epoch: 0, Version: "1.2.0", Release: "2.20180605git4a062cf.fc28", Arch: "x86_64",
					SrcName: "libnsl2", SrcEpoch: 0, SrcVersion: "1.2.0", SrcRelease: "2.20180605git4a062cf.fc28",
					Modularitylabel: "", Licenses: []string{"BSD and LGPLv2+"},
				},
				{
					Name: "python3-pip", Epoch: 0, Version: "9.0.3", Release: "2.fc28", Arch: "noarch",
					SrcName: "python-pip", SrcEpoch: 0, SrcVersion: "9.0.3", SrcRelease: "2.fc28", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "python3", Epoch: 0, Version: "3.6.6", Release: "1.fc28", Arch: "x86_64", SrcName: "python3",
					SrcEpoch: 0, SrcVersion: "3.6.6", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"Python"},
				},
				{
					Name: "pam", Epoch: 0, Version: "1.3.1", Release: "1.fc28", Arch: "x86_64", SrcName: "pam",
					SrcEpoch: 0, SrcVersion: "1.3.1", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"BSD and GPLv2+"},
				},
				{
					Name: "python3-gobject-base", Epoch: 0, Version: "3.28.3", Release: "1.fc28", Arch: "x86_64",
					SrcName: "pygobject3", SrcEpoch: 0, SrcVersion: "3.28.3", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+ and MIT"},
				},
				{
					Name: "python3-smartcols", Epoch: 0, Version: "0.3.0", Release: "2.fc28", Arch: "x86_64",
					SrcName: "python-smartcols", SrcEpoch: 0, SrcVersion: "0.3.0", SrcRelease: "2.fc28",
					Modularitylabel: "", Licenses: []string{"GPLv3+"},
				},
				{
					Name: "python3-iniparse", Epoch: 0, Version: "0.4", Release: "30.fc28", Arch: "noarch",
					SrcName: "python-iniparse", SrcEpoch: 0, SrcVersion: "0.4", SrcRelease: "30.fc28",
					Modularitylabel: "", Licenses: []string{"MIT and Python"},
				},
				{
					Name: "openldap", Epoch: 0, Version: "2.4.46", Release: "3.fc28", Arch: "x86_64",
					SrcName: "openldap", SrcEpoch: 0, SrcVersion: "2.4.46", SrcRelease: "3.fc28", Modularitylabel: "",
					Licenses: []string{"OpenLDAP"},
				},
				{
					Name: "libseccomp", Epoch: 0, Version: "2.3.3", Release: "2.fc28", Arch: "x86_64",
					SrcName: "libseccomp", SrcEpoch: 0, SrcVersion: "2.3.3", SrcRelease: "2.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2"},
				},
				{
					Name: "npth", Epoch: 0, Version: "1.5", Release: "4.fc28", Arch: "x86_64", SrcName: "npth",
					SrcEpoch: 0, SrcVersion: "1.5", SrcRelease: "4.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "gpgme", Epoch: 0, Version: "1.10.0", Release: "4.fc28", Arch: "x86_64", SrcName: "gpgme",
					SrcEpoch: 0, SrcVersion: "1.10.0", SrcRelease: "4.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "json-c", Epoch: 0, Version: "0.13.1", Release: "2.fc28", Arch: "x86_64", SrcName: "json-c",
					SrcEpoch: 0, SrcVersion: "0.13.1", SrcRelease: "2.fc28", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "libyaml", Epoch: 0, Version: "0.1.7", Release: "5.fc28", Arch: "x86_64", SrcName: "libyaml",
					SrcEpoch: 0, SrcVersion: "0.1.7", SrcRelease: "5.fc28", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "libpkgconf", Epoch: 0, Version: "1.4.2", Release: "1.fc28", Arch: "x86_64",
					SrcName: "pkgconf", SrcEpoch: 0, SrcVersion: "1.4.2", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"ISC"},
				},
				{
					Name: "pkgconf-pkg-config", Epoch: 0, Version: "1.4.2", Release: "1.fc28", Arch: "x86_64",
					SrcName: "pkgconf", SrcEpoch: 0, SrcVersion: "1.4.2", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"ISC"},
				},
				{
					Name: "iptables-libs", Epoch: 0, Version: "1.6.2", Release: "3.fc28", Arch: "x86_64",
					SrcName: "iptables", SrcEpoch: 0, SrcVersion: "1.6.2", SrcRelease: "3.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2 and Artistic Licence 2.0 and ISC"},
				},
				{
					Name: "device-mapper-libs", Epoch: 0, Version: "1.02.146", Release: "5.fc28", Arch: "x86_64",
					SrcName: "lvm2", SrcEpoch: 0, SrcVersion: "2.02.177", SrcRelease: "5.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2"},
				},
				{
					Name: "systemd-pam", Epoch: 0, Version: "238", Release: "9.git0e0aa59.fc28", Arch: "x86_64",
					SrcName: "systemd", SrcEpoch: 0, SrcVersion: "238", SrcRelease: "9.git0e0aa59.fc28",
					Modularitylabel: "", Licenses: []string{"LGPLv2+ and MIT and GPLv2+"},
				},
				{
					Name: "systemd", Epoch: 0, Version: "238", Release: "9.git0e0aa59.fc28", Arch: "x86_64",
					SrcName: "systemd", SrcEpoch: 0, SrcVersion: "238", SrcRelease: "9.git0e0aa59.fc28",
					Modularitylabel: "", Licenses: []string{"LGPLv2+ and MIT and GPLv2+"},
				},
				{
					Name: "elfutils-default-yama-scope", Epoch: 0, Version: "0.173", Release: "1.fc28", Arch: "noarch",
					SrcName: "elfutils", SrcEpoch: 0, SrcVersion: "0.173", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+ or LGPLv3+"},
				},
				{
					Name: "libcurl", Epoch: 0, Version: "7.59.0", Release: "6.fc28", Arch: "x86_64", SrcName: "curl",
					SrcEpoch: 0, SrcVersion: "7.59.0", SrcRelease: "6.fc28", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "python3-librepo", Epoch: 0, Version: "1.8.1", Release: "7.fc28", Arch: "x86_64",
					SrcName: "librepo", SrcEpoch: 0, SrcVersion: "1.8.1", SrcRelease: "7.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "rpm-plugin-selinux", Epoch: 0, Version: "4.14.1", Release: "9.fc28", Arch: "x86_64",
					SrcName: "rpm", SrcEpoch: 0, SrcVersion: "4.14.1", SrcRelease: "9.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+"},
				},
				{
					Name: "rpm", Epoch: 0, Version: "4.14.1", Release: "9.fc28", Arch: "x86_64", SrcName: "rpm",
					SrcEpoch: 0, SrcVersion: "4.14.1", SrcRelease: "9.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+"},
				},
				{
					Name: "libdnf", Epoch: 0, Version: "0.11.1", Release: "3.fc28", Arch: "x86_64", SrcName: "libdnf",
					SrcEpoch: 0, SrcVersion: "0.11.1", SrcRelease: "3.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "rpm-build-libs", Epoch: 0, Version: "4.14.1", Release: "9.fc28", Arch: "x86_64",
					SrcName: "rpm", SrcEpoch: 0, SrcVersion: "4.14.1", SrcRelease: "9.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and LGPLv2+ with exceptions"},
				},
				{
					Name: "python3-rpm", Epoch: 0, Version: "4.14.1", Release: "9.fc28", Arch: "x86_64", SrcName: "rpm",
					SrcEpoch: 0, SrcVersion: "4.14.1", SrcRelease: "9.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+"},
				},
				{
					Name: "dnf", Epoch: 0, Version: "2.7.5", Release: "12.fc28", Arch: "noarch", SrcName: "dnf",
					SrcEpoch: 0, SrcVersion: "2.7.5", SrcRelease: "12.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and GPLv2 and GPL"},
				},
				{
					Name: "deltarpm", Epoch: 0, Version: "3.6", Release: "25.fc28", Arch: "x86_64", SrcName: "deltarpm",
					SrcEpoch: 0, SrcVersion: "3.6", SrcRelease: "25.fc28", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "sssd-client", Epoch: 0, Version: "1.16.3", Release: "2.fc28", Arch: "x86_64",
					SrcName: "sssd", SrcEpoch: 0, SrcVersion: "1.16.3", SrcRelease: "2.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv3+"},
				},
				{
					Name: "cracklib-dicts", Epoch: 0, Version: "2.9.6", Release: "13.fc28", Arch: "x86_64",
					SrcName: "cracklib", SrcEpoch: 0, SrcVersion: "2.9.6", SrcRelease: "13.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "tar", Epoch: 2, Version: "1.30", Release: "3.fc28", Arch: "x86_64", SrcName: "tar",
					SrcEpoch: 2, SrcVersion: "1.30", SrcRelease: "3.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "diffutils", Epoch: 0, Version: "3.6", Release: "4.fc28", Arch: "x86_64",
					SrcName: "diffutils", SrcEpoch: 0, SrcVersion: "3.6", SrcRelease: "4.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "langpacks-en", Epoch: 0, Version: "1.0", Release: "12.fc28", Arch: "noarch",
					SrcName: "langpacks", SrcEpoch: 0, SrcVersion: "1.0", SrcRelease: "12.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+"},
				},
				{
					Name: "gpg-pubkey", Epoch: 0, Version: "9db62fb1", Release: "59920156", Arch: "None", SrcName: "",
					SrcEpoch: 0, SrcVersion: "", SrcRelease: "", Modularitylabel: "", Licenses: []string{"pubkey"},
				},
				{
					Name: "libgcc", Epoch: 0, Version: "8.1.1", Release: "5.fc28", Arch: "x86_64", SrcName: "gcc",
					SrcEpoch: 0, SrcVersion: "8.1.1", SrcRelease: "5.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv3+ and GPLv3+ with exceptions and GPLv2+ with exceptions and LGPLv2+ and BSD"},
				},
				{
					Name: "pkgconf-m4", Epoch: 0, Version: "1.4.2", Release: "1.fc28", Arch: "noarch",
					SrcName: "pkgconf", SrcEpoch: 0, SrcVersion: "1.4.2", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+ with exceptions"},
				},
				{
					Name: "dnf-conf", Epoch: 0, Version: "2.7.5", Release: "12.fc28", Arch: "noarch", SrcName: "dnf",
					SrcEpoch: 0, SrcVersion: "2.7.5", SrcRelease: "12.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and GPLv2 and GPL"},
				},
				{
					Name: "fedora-repos", Epoch: 0, Version: "28", Release: "5", Arch: "noarch",
					SrcName: "fedora-repos", SrcEpoch: 0, SrcVersion: "28", SrcRelease: "5", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "setup", Epoch: 0, Version: "2.11.4", Release: "1.fc28", Arch: "noarch", SrcName: "setup",
					SrcEpoch: 0, SrcVersion: "2.11.4", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"Public Domain"},
				},
				{
					Name: "basesystem", Epoch: 0, Version: "11", Release: "5.fc28", Arch: "noarch",
					SrcName: "basesystem", SrcEpoch: 0, SrcVersion: "11", SrcRelease: "5.fc28", Modularitylabel: "",
					Licenses: []string{"Public Domain"},
				},
				{
					Name: "ncurses-base", Epoch: 0, Version: "6.1", Release: "5.20180224.fc28", Arch: "noarch",
					SrcName: "ncurses", SrcEpoch: 0, SrcVersion: "6.1", SrcRelease: "5.20180224.fc28",
					Modularitylabel: "", Licenses: []string{"MIT"},
				},
				{
					Name: "libselinux", Epoch: 0, Version: "2.8", Release: "1.fc28", Arch: "x86_64",
					SrcName: "libselinux", SrcEpoch: 0, SrcVersion: "2.8", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"Public Domain"},
				},
				{
					Name: "ncurses-libs", Epoch: 0, Version: "6.1", Release: "5.20180224.fc28", Arch: "x86_64",
					SrcName: "ncurses", SrcEpoch: 0, SrcVersion: "6.1", SrcRelease: "5.20180224.fc28",
					Modularitylabel: "", Licenses: []string{"MIT"},
				},
				{
					Name: "glibc", Epoch: 0, Version: "2.27", Release: "32.fc28", Arch: "x86_64", SrcName: "glibc",
					SrcEpoch: 0, SrcVersion: "2.27", SrcRelease: "32.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+ and LGPLv2+ with exceptions and GPLv2+ and GPLv2+ with exceptions and BSD and Inner-Net and ISC and Public Domain and GFDL"},
				},
				{
					Name: "libsepol", Epoch: 0, Version: "2.8", Release: "1.fc28", Arch: "x86_64", SrcName: "libsepol",
					SrcEpoch: 0, SrcVersion: "2.8", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "xz-libs", Epoch: 0, Version: "5.2.4", Release: "2.fc28", Arch: "x86_64", SrcName: "xz",
					SrcEpoch: 0, SrcVersion: "5.2.4", SrcRelease: "2.fc28", Modularitylabel: "",
					Licenses: []string{"Public Domain"},
				},
				{
					Name: "info", Epoch: 0, Version: "6.5", Release: "4.fc28", Arch: "x86_64", SrcName: "texinfo",
					SrcEpoch: 0, SrcVersion: "6.5", SrcRelease: "4.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "libdb", Epoch: 0, Version: "5.3.28", Release: "30.fc28", Arch: "x86_64", SrcName: "libdb",
					SrcEpoch: 0, SrcVersion: "5.3.28", SrcRelease: "30.fc28", Modularitylabel: "",
					Licenses: []string{"BSD and LGPLv2 and Sleepycat"},
				},
				{
					Name: "elfutils-libelf", Epoch: 0, Version: "0.173", Release: "1.fc28", Arch: "x86_64",
					SrcName: "elfutils", SrcEpoch: 0, SrcVersion: "0.173", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+ or LGPLv3+"},
				},
				{
					Name: "popt", Epoch: 0, Version: "1.16", Release: "14.fc28", Arch: "x86_64", SrcName: "popt",
					SrcEpoch: 0, SrcVersion: "1.16", SrcRelease: "14.fc28", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "nspr", Epoch: 0, Version: "4.19.0", Release: "1.fc28", Arch: "x86_64", SrcName: "nspr",
					SrcEpoch: 0, SrcVersion: "4.19.0", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"MPLv2.0"},
				},
				{
					Name: "libxcrypt", Epoch: 0, Version: "4.1.2", Release: "1.fc28", Arch: "x86_64",
					SrcName: "libxcrypt", SrcEpoch: 0, SrcVersion: "4.1.2", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+ and BSD and Public Domain"},
				},
				{
					Name: "lua-libs", Epoch: 0, Version: "5.3.4", Release: "10.fc28", Arch: "x86_64", SrcName: "lua",
					SrcEpoch: 0, SrcVersion: "5.3.4", SrcRelease: "10.fc28", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "libuuid", Epoch: 0, Version: "2.32.1", Release: "1.fc28", Arch: "x86_64",
					SrcName: "util-linux", SrcEpoch: 0, SrcVersion: "2.32.1", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "readline", Epoch: 0, Version: "7.0", Release: "11.fc28", Arch: "x86_64", SrcName: "readline",
					SrcEpoch: 0, SrcVersion: "7.0", SrcRelease: "11.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "libattr", Epoch: 0, Version: "2.4.48", Release: "3.fc28", Arch: "x86_64", SrcName: "attr",
					SrcEpoch: 0, SrcVersion: "2.4.48", SrcRelease: "3.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "coreutils-single", Epoch: 0, Version: "8.29", Release: "7.fc28", Arch: "x86_64",
					SrcName: "coreutils", SrcEpoch: 0, SrcVersion: "8.29", SrcRelease: "7.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "libblkid", Epoch: 0, Version: "2.32.1", Release: "1.fc28", Arch: "x86_64",
					SrcName: "util-linux", SrcEpoch: 0, SrcVersion: "2.32.1", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "gmp", Epoch: 1, Version: "6.1.2", Release: "7.fc28", Arch: "x86_64", SrcName: "gmp",
					SrcEpoch: 1, SrcVersion: "6.1.2", SrcRelease: "7.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv3+ or GPLv2+"},
				},
				{
					Name: "libunistring", Epoch: 0, Version: "0.9.10", Release: "1.fc28", Arch: "x86_64",
					SrcName: "libunistring", SrcEpoch: 0, SrcVersion: "0.9.10", SrcRelease: "1.fc28",
					Modularitylabel: "", Licenses: []string{"GPLV2+ or LGPLv3+"},
				},
				{
					Name: "sqlite-libs", Epoch: 0, Version: "3.22.0", Release: "4.fc28", Arch: "x86_64",
					SrcName: "sqlite", SrcEpoch: 0, SrcVersion: "3.22.0", SrcRelease: "4.fc28", Modularitylabel: "",
					Licenses: []string{"Public Domain"},
				},
				{
					Name: "audit-libs", Epoch: 0, Version: "2.8.4", Release: "2.fc28", Arch: "x86_64", SrcName: "audit",
					SrcEpoch: 0, SrcVersion: "2.8.4", SrcRelease: "2.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "chkconfig", Epoch: 0, Version: "1.10", Release: "4.fc28", Arch: "x86_64",
					SrcName: "chkconfig", SrcEpoch: 0, SrcVersion: "1.10", SrcRelease: "4.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2"},
				},
				{
					Name: "libsmartcols", Epoch: 0, Version: "2.32.1", Release: "1.fc28", Arch: "x86_64",
					SrcName: "util-linux", SrcEpoch: 0, SrcVersion: "2.32.1", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "pcre", Epoch: 0, Version: "8.42", Release: "3.fc28", Arch: "x86_64", SrcName: "pcre",
					SrcEpoch: 0, SrcVersion: "8.42", SrcRelease: "3.fc28", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "grep", Epoch: 0, Version: "3.1", Release: "5.fc28", Arch: "x86_64", SrcName: "grep",
					SrcEpoch: 0, SrcVersion: "3.1", SrcRelease: "5.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "crypto-policies", Epoch: 0, Version: "20180425", Release: "5.git6ad4018.fc28",
					Arch: "noarch", SrcName: "crypto-policies", SrcEpoch: 0, SrcVersion: "20180425",
					SrcRelease: "5.git6ad4018.fc28", Modularitylabel: "", Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "gdbm-libs", Epoch: 1, Version: "1.14.1", Release: "4.fc28", Arch: "x86_64", SrcName: "gdbm",
					SrcEpoch: 1, SrcVersion: "1.14.1", SrcRelease: "4.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "p11-kit-trust", Epoch: 0, Version: "0.23.12", Release: "1.fc28", Arch: "x86_64",
					SrcName: "p11-kit", SrcEpoch: 0, SrcVersion: "0.23.12", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "openssl-libs", Epoch: 1, Version: "1.1.0h", Release: "3.fc28", Arch: "x86_64",
					SrcName: "openssl", SrcEpoch: 1, SrcVersion: "1.1.0h", SrcRelease: "3.fc28", Modularitylabel: "",
					Licenses: []string{"OpenSSL"},
				},
				{
					Name: "ima-evm-utils", Epoch: 0, Version: "1.1", Release: "2.fc28", Arch: "x86_64",
					SrcName: "ima-evm-utils", SrcEpoch: 0, SrcVersion: "1.1", SrcRelease: "2.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2"},
				},
				{
					Name: "gdbm", Epoch: 1, Version: "1.14.1", Release: "4.fc28", Arch: "x86_64", SrcName: "gdbm",
					SrcEpoch: 1, SrcVersion: "1.14.1", SrcRelease: "4.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "gobject-introspection", Epoch: 0, Version: "1.56.1", Release: "1.fc28", Arch: "x86_64",
					SrcName: "gobject-introspection", SrcEpoch: 0, SrcVersion: "1.56.1", SrcRelease: "1.fc28",
					Modularitylabel: "", Licenses: []string{"GPLv2+, LGPLv2+, MIT"},
				},
				{
					Name: "shadow-utils", Epoch: 2, Version: "4.6", Release: "1.fc28", Arch: "x86_64",
					SrcName: "shadow-utils", SrcEpoch: 2, SrcVersion: "4.6", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"BSD and GPLv2+"},
				},
				{
					Name: "libpsl", Epoch: 0, Version: "0.20.2", Release: "2.fc28", Arch: "x86_64", SrcName: "libpsl",
					SrcEpoch: 0, SrcVersion: "0.20.2", SrcRelease: "2.fc28", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "nettle", Epoch: 0, Version: "3.4", Release: "2.fc28", Arch: "x86_64", SrcName: "nettle",
					SrcEpoch: 0, SrcVersion: "3.4", SrcRelease: "2.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv3+ or GPLv2+"},
				},
				{
					Name: "libfdisk", Epoch: 0, Version: "2.32.1", Release: "1.fc28", Arch: "x86_64",
					SrcName: "util-linux", SrcEpoch: 0, SrcVersion: "2.32.1", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "cracklib", Epoch: 0, Version: "2.9.6", Release: "13.fc28", Arch: "x86_64",
					SrcName: "cracklib", SrcEpoch: 0, SrcVersion: "2.9.6", SrcRelease: "13.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "libcomps", Epoch: 0, Version: "0.1.8", Release: "11.fc28", Arch: "x86_64",
					SrcName: "libcomps", SrcEpoch: 0, SrcVersion: "0.1.8", SrcRelease: "11.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+"},
				},
				{
					Name: "nss-softokn", Epoch: 0, Version: "3.38.0", Release: "1.0.fc28", Arch: "x86_64",
					SrcName: "nss-softokn", SrcEpoch: 0, SrcVersion: "3.38.0", SrcRelease: "1.0.fc28",
					Modularitylabel: "", Licenses: []string{"MPLv2.0"},
				},
				{
					Name: "nss-sysinit", Epoch: 0, Version: "3.38.0", Release: "1.0.fc28", Arch: "x86_64",
					SrcName: "nss", SrcEpoch: 0, SrcVersion: "3.38.0", SrcRelease: "1.0.fc28", Modularitylabel: "",
					Licenses: []string{"MPLv2.0"},
				},
				{
					Name: "libksba", Epoch: 0, Version: "1.3.5", Release: "7.fc28", Arch: "x86_64", SrcName: "libksba",
					SrcEpoch: 0, SrcVersion: "1.3.5", SrcRelease: "7.fc28", Modularitylabel: "",
					Licenses: []string{"(LGPLv3+ or GPLv2+) and GPLv3+"},
				},
				{
					Name: "kmod-libs", Epoch: 0, Version: "25", Release: "2.fc28", Arch: "x86_64", SrcName: "kmod",
					SrcEpoch: 0, SrcVersion: "25", SrcRelease: "2.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "libsss_nss_idmap", Epoch: 0, Version: "1.16.3", Release: "2.fc28", Arch: "x86_64",
					SrcName: "sssd", SrcEpoch: 0, SrcVersion: "1.16.3", SrcRelease: "2.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv3+"},
				},
				{
					Name: "libverto", Epoch: 0, Version: "0.3.0", Release: "5.fc28", Arch: "x86_64",
					SrcName: "libverto", SrcEpoch: 0, SrcVersion: "0.3.0", SrcRelease: "5.fc28", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "gawk", Epoch: 0, Version: "4.2.1", Release: "1.fc28", Arch: "x86_64", SrcName: "gawk",
					SrcEpoch: 0, SrcVersion: "4.2.1", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv3+ and GPLv2+ and LGPLv2+ and BSD"},
				},
				{
					Name: "libtirpc", Epoch: 0, Version: "1.0.3", Release: "3.rc2.fc28", Arch: "x86_64",
					SrcName: "libtirpc", SrcEpoch: 0, SrcVersion: "1.0.3", SrcRelease: "3.rc2.fc28",
					Modularitylabel: "", Licenses: []string{"SISSL and BSD"},
				},
				{
					Name: "python3-libs", Epoch: 0, Version: "3.6.6", Release: "1.fc28", Arch: "x86_64",
					SrcName: "python3", SrcEpoch: 0, SrcVersion: "3.6.6", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"Python"},
				},
				{
					Name: "python3-setuptools", Epoch: 0, Version: "39.2.0", Release: "6.fc28", Arch: "noarch",
					SrcName: "python-setuptools", SrcEpoch: 0, SrcVersion: "39.2.0", SrcRelease: "6.fc28",
					Modularitylabel: "", Licenses: []string{"MIT"},
				},
				{
					Name: "libpwquality", Epoch: 0, Version: "1.4.0", Release: "7.fc28", Arch: "x86_64",
					SrcName: "libpwquality", SrcEpoch: 0, SrcVersion: "1.4.0", SrcRelease: "7.fc28",
					Modularitylabel: "", Licenses: []string{"BSD or GPLv2+"},
				},
				{
					Name: "util-linux", Epoch: 0, Version: "2.32.1", Release: "1.fc28", Arch: "x86_64",
					SrcName: "util-linux", SrcEpoch: 0, SrcVersion: "2.32.1", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2 and GPLv2+ and LGPLv2+ and BSD with advertising and Public Domain"},
				},
				{
					Name: "python3-libcomps", Epoch: 0, Version: "0.1.8", Release: "11.fc28", Arch: "x86_64",
					SrcName: "libcomps", SrcEpoch: 0, SrcVersion: "0.1.8", SrcRelease: "11.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+"},
				},
				{
					Name: "python3-six", Epoch: 0, Version: "1.11.0", Release: "3.fc28", Arch: "noarch",
					SrcName: "python-six", SrcEpoch: 0, SrcVersion: "1.11.0", SrcRelease: "3.fc28", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "cyrus-sasl-lib", Epoch: 0, Version: "2.1.27", Release: "0.2rc7.fc28", Arch: "x86_64",
					SrcName: "cyrus-sasl", SrcEpoch: 0, SrcVersion: "2.1.27", SrcRelease: "0.2rc7.fc28",
					Modularitylabel: "", Licenses: []string{"BSD with advertising"},
				},
				{
					Name: "libssh", Epoch: 0, Version: "0.8.2", Release: "1.fc28", Arch: "x86_64", SrcName: "libssh",
					SrcEpoch: 0, SrcVersion: "0.8.2", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "qrencode-libs", Epoch: 0, Version: "3.4.4", Release: "5.fc28", Arch: "x86_64",
					SrcName: "qrencode", SrcEpoch: 0, SrcVersion: "3.4.4", SrcRelease: "5.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "gnupg2", Epoch: 0, Version: "2.2.8", Release: "1.fc28", Arch: "x86_64", SrcName: "gnupg2",
					SrcEpoch: 0, SrcVersion: "2.2.8", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "python3-gpg", Epoch: 0, Version: "1.10.0", Release: "4.fc28", Arch: "x86_64",
					SrcName: "gpgme", SrcEpoch: 0, SrcVersion: "1.10.0", SrcRelease: "4.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "libargon2", Epoch: 0, Version: "20161029", Release: "5.fc28", Arch: "x86_64",
					SrcName: "argon2", SrcEpoch: 0, SrcVersion: "20161029", SrcRelease: "5.fc28", Modularitylabel: "",
					Licenses: []string{"Public Domain or ASL 2.0"},
				},
				{
					Name: "libmodulemd", Epoch: 0, Version: "1.6.2", Release: "2.fc28", Arch: "x86_64",
					SrcName: "libmodulemd", SrcEpoch: 0, SrcVersion: "1.6.2", SrcRelease: "2.fc28", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "pkgconf", Epoch: 0, Version: "1.4.2", Release: "1.fc28", Arch: "x86_64", SrcName: "pkgconf",
					SrcEpoch: 0, SrcVersion: "1.4.2", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"ISC"},
				},
				{
					Name: "libpcap", Epoch: 14, Version: "1.9.0", Release: "1.fc28", Arch: "x86_64", SrcName: "libpcap",
					SrcEpoch: 14, SrcVersion: "1.9.0", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"BSD with advertising"},
				},
				{
					Name: "device-mapper", Epoch: 0, Version: "1.02.146", Release: "5.fc28", Arch: "x86_64",
					SrcName: "lvm2", SrcEpoch: 0, SrcVersion: "2.02.177", SrcRelease: "5.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2"},
				},
				{
					Name: "cryptsetup-libs", Epoch: 0, Version: "2.0.4", Release: "1.fc28", Arch: "x86_64",
					SrcName: "cryptsetup", SrcEpoch: 0, SrcVersion: "2.0.4", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and LGPLv2+"},
				},
				{
					Name: "elfutils-libs", Epoch: 0, Version: "0.173", Release: "1.fc28", Arch: "x86_64",
					SrcName: "elfutils", SrcEpoch: 0, SrcVersion: "0.173", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+ or LGPLv3+"},
				},
				{
					Name: "dbus", Epoch: 1, Version: "1.12.10", Release: "1.fc28", Arch: "x86_64", SrcName: "dbus",
					SrcEpoch: 1, SrcVersion: "1.12.10", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"(GPLv2+ or AFL) and GPLv2+"},
				},
				{
					Name: "libnghttp2", Epoch: 0, Version: "1.32.1", Release: "1.fc28", Arch: "x86_64",
					SrcName: "nghttp2", SrcEpoch: 0, SrcVersion: "1.32.1", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "librepo", Epoch: 0, Version: "1.8.1", Release: "7.fc28", Arch: "x86_64", SrcName: "librepo",
					SrcEpoch: 0, SrcVersion: "1.8.1", SrcRelease: "7.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "curl", Epoch: 0, Version: "7.59.0", Release: "6.fc28", Arch: "x86_64", SrcName: "curl",
					SrcEpoch: 0, SrcVersion: "7.59.0", SrcRelease: "6.fc28", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "rpm-libs", Epoch: 0, Version: "4.14.1", Release: "9.fc28", Arch: "x86_64", SrcName: "rpm",
					SrcEpoch: 0, SrcVersion: "4.14.1", SrcRelease: "9.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and LGPLv2+ with exceptions"},
				},
				{
					Name: "libsolv", Epoch: 0, Version: "0.6.35", Release: "1.fc28", Arch: "x86_64", SrcName: "libsolv",
					SrcEpoch: 0, SrcVersion: "0.6.35", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "python3-hawkey", Epoch: 0, Version: "0.11.1", Release: "3.fc28", Arch: "x86_64",
					SrcName: "libdnf", SrcEpoch: 0, SrcVersion: "0.11.1", SrcRelease: "3.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "rpm-sign-libs", Epoch: 0, Version: "4.14.1", Release: "9.fc28", Arch: "x86_64",
					SrcName: "rpm", SrcEpoch: 0, SrcVersion: "4.14.1", SrcRelease: "9.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and LGPLv2+ with exceptions"},
				},
				{
					Name: "python3-dnf", Epoch: 0, Version: "2.7.5", Release: "12.fc28", Arch: "noarch", SrcName: "dnf",
					SrcEpoch: 0, SrcVersion: "2.7.5", SrcRelease: "12.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and GPLv2 and GPL"},
				},
				{
					Name: "dnf-yum", Epoch: 0, Version: "2.7.5", Release: "12.fc28", Arch: "noarch", SrcName: "dnf",
					SrcEpoch: 0, SrcVersion: "2.7.5", SrcRelease: "12.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and GPLv2 and GPL"},
				},
				{
					Name: "rpm-plugin-systemd-inhibit", Epoch: 0, Version: "4.14.1", Release: "9.fc28", Arch: "x86_64",
					SrcName: "rpm", SrcEpoch: 0, SrcVersion: "4.14.1", SrcRelease: "9.fc28", Modularitylabel: "",
					Licenses: []string{"GPLv2+"},
				},
				{
					Name: "nss-tools", Epoch: 0, Version: "3.38.0", Release: "1.0.fc28", Arch: "x86_64", SrcName: "nss",
					SrcEpoch: 0, SrcVersion: "3.38.0", SrcRelease: "1.0.fc28", Modularitylabel: "",
					Licenses: []string{"MPLv2.0"},
				},
				{
					Name: "openssl-pkcs11", Epoch: 0, Version: "0.4.8", Release: "1.fc28", Arch: "x86_64",
					SrcName: "openssl-pkcs11", SrcEpoch: 0, SrcVersion: "0.4.8", SrcRelease: "1.fc28",
					Modularitylabel: "", Licenses: []string{"LGPLv2+ and BSD"},
				},
				{
					Name: "vim-minimal", Epoch: 2, Version: "8.1.328", Release: "1.fc28", Arch: "x86_64",
					SrcName: "vim", SrcEpoch: 2, SrcVersion: "8.1.328", SrcRelease: "1.fc28", Modularitylabel: "",
					Licenses: []string{"Vim and MIT"},
				},
				{
					Name: "glibc-langpack-en", Epoch: 0, Version: "2.27", Release: "32.fc28", Arch: "x86_64",
					SrcName: "glibc", SrcEpoch: 0, SrcVersion: "2.27", SrcRelease: "32.fc28", Modularitylabel: "",
					Licenses: []string{"LGPLv2+ and LGPLv2+ with exceptions and GPLv2+ and GPLv2+ with exceptions and BSD and Inner-Net and ISC and Public Domain and GFDL"},
				},
				{
					Name: "rootfiles", Epoch: 0, Version: "8.1", Release: "22.fc28", Arch: "noarch",
					SrcName: "rootfiles", SrcEpoch: 0, SrcVersion: "8.1", SrcRelease: "22.fc28", Modularitylabel: "",
					Licenses: []string{"Public Domain"},
				},
			},
		},
		"ValidWithModularitylabel": {
			path: "./testdata/valid_with_modularitylabel",
			// docker run --name centos -it --rm  centos:8 /bin/bash
			// docker cp ./testdata/valid_with_modularitylabel centos:/tmp/Packages
			//
			// $ cat rpmqa.py
			// #!/bin/python3
			//
			// import rpm
			//
			// def splitFilename(filename):
			//     # sourcerpm spec: https://github.com/rpm-software-management/dnf/blob/4.2.23/dnf/package.py#L116-L120
			//     srcname = rtrim(filename, ".src.rpm")
			//     sname, sversion, srelease = srcname.rsplit('-', 2)
			//     return sname, sversion, srelease
			//
			// # ref. https://github.com/rpm-software-management/dnf/blob/4.2.23/dnf/util.py#L122
			// def rtrim(s, r):
			//     if s.endswith(r):
			//         s = s[:-len(r)]
			//     return s
			//
			// def license_format(s):
			//     return s.replace(" and ",",").replace(" or ",",")
			//
			// rpm.addMacro('_dbpath', '/tmp/')
			// ts = rpm.TransactionSet()
			// mi = ts.dbMatch()
			// for h in mi:
			//     sname = sversion = srelease = ""
			//     if h[rpm.RPMTAG_SOURCERPM] != "(none)":
			//         sname, sversion, srelease = splitFilename(h[rpm.RPMTAG_SOURCERPM])
			//
			//     mlabel = h[rpm.RPMTAG_MODULARITYLABEL] if h[rpm.RPMTAG_MODULARITYLABEL] is not None else  ""
			//     print("{{Name: \"{0}\", Epoch: {1}, Version: \"{2}\", Release: \"{3}\", Arch: \"{4}\", SrcName: \"{5}\", SrcEpoch: {6}, SrcVersion: \"{7}\", SrcRelease: \"{8}\", Modularitylabel: \"{9}\",Licenses: []string{\"{10}\"}},".format(h[rpm.RPMTAG_NAME], h[rpm.RPMTAG_EPOCHNUM], h[rpm.RPMTAG_VERSION], h[rpm.RPMTAG_RELEASE], h[rpm.RPMTAG_ARCH], sname, h[rpm.RPMTAG_EPOCHNUM], sversion, srelease, mlabel},license_format(h[rpm.RPMTAG_LICENSE])))
			pkgs: []types.Package{
				{
					Name: "perl-podlators", Epoch: 0, Version: "4.11", Release: "1.el8", Arch: "noarch",
					SrcName: "perl-podlators", SrcEpoch: 0, SrcVersion: "4.11", SrcRelease: "1.el8",
					Modularitylabel: "", Licenses: []string{"(GPL+ or Artistic) and FSFAP"},
				},
				{
					Name: "python3-setuptools-wheel", Epoch: 0, Version: "39.2.0", Release: "5.el8", Arch: "noarch",
					SrcName: "python-setuptools", SrcEpoch: 0, SrcVersion: "39.2.0", SrcRelease: "5.el8",
					Modularitylabel: "", Licenses: []string{"MIT"},
				},
				{
					Name: "perl-Pod-Perldoc", Epoch: 0, Version: "3.28", Release: "396.el8", Arch: "noarch",
					SrcName: "perl-Pod-Perldoc", SrcEpoch: 0, SrcVersion: "3.28", SrcRelease: "396.el8",
					Modularitylabel: "", Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "perl-IO-Socket-SSL", Epoch: 0, Version: "2.066", Release: "4.el8", Arch: "noarch",
					SrcName: "perl-IO-Socket-SSL", SrcEpoch: 0, SrcVersion: "2.066", SrcRelease: "4.el8",
					Modularitylabel: "", Licenses: []string{"(GPL+ or Artistic) and MPLv2.0"},
				},
				{
					Name: "perl-URI", Epoch: 0, Version: "1.73", Release: "3.el8", Arch: "noarch", SrcName: "perl-URI",
					SrcEpoch: 0, SrcVersion: "1.73", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "filesystem", Epoch: 0, Version: "3.8", Release: "2.el8", Arch: "x86_64",
					SrcName: "filesystem", SrcEpoch: 0, SrcVersion: "3.8", SrcRelease: "2.el8", Modularitylabel: "",
					Licenses: []string{"Public Domain"},
				},
				{
					Name: "emacs-filesystem", Epoch: 1, Version: "26.1", Release: "5.el8", Arch: "noarch",
					SrcName: "emacs", SrcEpoch: 1, SrcVersion: "26.1", SrcRelease: "5.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+ and CC0-1.0"},
				},
				{
					Name: "git", Epoch: 0, Version: "2.18.4", Release: "2.el8_2", Arch: "x86_64", SrcName: "git",
					SrcEpoch: 0, SrcVersion: "2.18.4", SrcRelease: "2.el8_2", Modularitylabel: "",
					Licenses: []string{"GPLv2"},
				},
				{
					Name: "pcre2", Epoch: 0, Version: "10.32", Release: "1.el8", Arch: "x86_64", SrcName: "pcre2",
					SrcEpoch: 0, SrcVersion: "10.32", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "vim-common", Epoch: 2, Version: "8.0.1763", Release: "13.el8", Arch: "x86_64",
					SrcName: "vim", SrcEpoch: 2, SrcVersion: "8.0.1763", SrcRelease: "13.el8", Modularitylabel: "",
					Licenses: []string{"Vim and MIT"},
				},
				{
					Name: "ncurses-libs", Epoch: 0, Version: "6.1", Release: "7.20180224.el8", Arch: "x86_64",
					SrcName: "ncurses", SrcEpoch: 0, SrcVersion: "6.1", SrcRelease: "7.20180224.el8",
					Modularitylabel: "", Licenses: []string{"MIT"},
				},
				{
					Name: "vim-enhanced", Epoch: 2, Version: "8.0.1763", Release: "13.el8", Arch: "x86_64",
					SrcName: "vim", SrcEpoch: 2, SrcVersion: "8.0.1763", SrcRelease: "13.el8", Modularitylabel: "",
					Licenses: []string{"Vim and MIT"},
				},
				{
					Name: "glibc-common", Epoch: 0, Version: "2.28", Release: "101.el8", Arch: "x86_64",
					SrcName: "glibc", SrcEpoch: 0, SrcVersion: "2.28", SrcRelease: "101.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+ and LGPLv2+ with exceptions and GPLv2+ and GPLv2+ with exceptions and BSD and Inner-Net and ISC and Public Domain and GFDL"},
				},
				{
					Name: "openssl-devel", Epoch: 1, Version: "1.1.1c", Release: "15.el8", Arch: "x86_64",
					SrcName: "openssl", SrcEpoch: 1, SrcVersion: "1.1.1c", SrcRelease: "15.el8", Modularitylabel: "",
					Licenses: []string{"OpenSSL"},
				},
				{
					Name: "bash", Epoch: 0, Version: "4.4.19", Release: "10.el8", Arch: "x86_64", SrcName: "bash",
					SrcEpoch: 0, SrcVersion: "4.4.19", SrcRelease: "10.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "popt-devel", Epoch: 0, Version: "1.16", Release: "14.el8", Arch: "x86_64", SrcName: "popt",
					SrcEpoch: 0, SrcVersion: "1.16", SrcRelease: "14.el8", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "libarchive-devel", Epoch: 0, Version: "3.3.2", Release: "8.el8_1", Arch: "x86_64",
					SrcName: "libarchive", SrcEpoch: 0, SrcVersion: "3.3.2", SrcRelease: "8.el8_1", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "bzip2-libs", Epoch: 0, Version: "1.0.6", Release: "26.el8", Arch: "x86_64", SrcName: "bzip2",
					SrcEpoch: 0, SrcVersion: "1.0.6", SrcRelease: "26.el8", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "xz-lzma-compat", Epoch: 0, Version: "5.2.4", Release: "3.el8", Arch: "x86_64", SrcName: "xz",
					SrcEpoch: 0, SrcVersion: "5.2.4", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"Public Domain"},
				},
				{
					Name: "libgpg-error", Epoch: 0, Version: "1.31", Release: "1.el8", Arch: "x86_64",
					SrcName: "libgpg-error", SrcEpoch: 0, SrcVersion: "1.31", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "libdb-devel", Epoch: 0, Version: "5.3.28", Release: "37.el8", Arch: "x86_64",
					SrcName: "libdb", SrcEpoch: 0, SrcVersion: "5.3.28", SrcRelease: "37.el8", Modularitylabel: "",
					Licenses: []string{"BSD and LGPLv2 and Sleepycat"},
				},
				{
					Name: "elfutils-libelf", Epoch: 0, Version: "0.178", Release: "7.el8", Arch: "x86_64",
					SrcName: "elfutils", SrcEpoch: 0, SrcVersion: "0.178", SrcRelease: "7.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+ or LGPLv3+"},
				},
				{
					Name: "libgomp", Epoch: 0, Version: "8.3.1", Release: "5.el8.0.2", Arch: "x86_64", SrcName: "gcc",
					SrcEpoch: 0, SrcVersion: "8.3.1", SrcRelease: "5.el8.0.2", Modularitylabel: "",
					Licenses: []string{"GPLv3+ and GPLv3+ with exceptions and GPLv2+ with exceptions and LGPLv2+ and BSD"},
				},
				{
					Name: "libxcrypt", Epoch: 0, Version: "4.1.1", Release: "4.el8", Arch: "x86_64",
					SrcName: "libxcrypt", SrcEpoch: 0, SrcVersion: "4.1.1", SrcRelease: "4.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+ and BSD and Public Domain"},
				},
				{
					Name: "gettext-libs", Epoch: 0, Version: "0.19.8.1", Release: "17.el8", Arch: "x86_64",
					SrcName: "gettext", SrcEpoch: 0, SrcVersion: "0.19.8.1", SrcRelease: "17.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+ and GPLv3+"},
				},
				{
					Name: "sqlite-libs", Epoch: 0, Version: "3.26.0", Release: "6.el8", Arch: "x86_64",
					SrcName: "sqlite", SrcEpoch: 0, SrcVersion: "3.26.0", SrcRelease: "6.el8", Modularitylabel: "",
					Licenses: []string{"Public Domain"},
				},
				{
					Name: "cpp", Epoch: 0, Version: "8.3.1", Release: "5.el8.0.2", Arch: "x86_64", SrcName: "gcc",
					SrcEpoch: 0, SrcVersion: "8.3.1", SrcRelease: "5.el8.0.2", Modularitylabel: "",
					Licenses: []string{"GPLv3+ and GPLv3+ with exceptions and GPLv2+ with exceptions and LGPLv2+ and BSD"},
				},
				{
					Name: "libstdc++", Epoch: 0, Version: "8.3.1", Release: "5.el8.0.2", Arch: "x86_64", SrcName: "gcc",
					SrcEpoch: 0, SrcVersion: "8.3.1", SrcRelease: "5.el8.0.2", Modularitylabel: "",
					Licenses: []string{"GPLv3+ and GPLv3+ with exceptions and GPLv2+ with exceptions and LGPLv2+ and BSD"},
				},
				{
					Name: "m4", Epoch: 0, Version: "1.4.18", Release: "7.el8", Arch: "x86_64", SrcName: "m4",
					SrcEpoch: 0, SrcVersion: "1.4.18", SrcRelease: "7.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "popt", Epoch: 0, Version: "1.16", Release: "14.el8", Arch: "x86_64", SrcName: "popt",
					SrcEpoch: 0, SrcVersion: "1.16", SrcRelease: "14.el8", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "libgpg-error-devel", Epoch: 0, Version: "1.31", Release: "1.el8", Arch: "x86_64",
					SrcName: "libgpg-error", SrcEpoch: 0, SrcVersion: "1.31", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "readline", Epoch: 0, Version: "7.0", Release: "10.el8", Arch: "x86_64", SrcName: "readline",
					SrcEpoch: 0, SrcVersion: "7.0", SrcRelease: "10.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "glibc-headers", Epoch: 0, Version: "2.28", Release: "101.el8", Arch: "x86_64",
					SrcName: "glibc", SrcEpoch: 0, SrcVersion: "2.28", SrcRelease: "101.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+ and LGPLv2+ with exceptions and GPLv2+ and GPLv2+ with exceptions and BSD and Inner-Net and ISC and Public Domain and GFDL"},
				},
				{
					Name: "json-c", Epoch: 0, Version: "0.13.1", Release: "0.2.el8", Arch: "x86_64", SrcName: "json-c",
					SrcEpoch: 0, SrcVersion: "0.13.1", SrcRelease: "0.2.el8", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "glibc-devel", Epoch: 0, Version: "2.28", Release: "101.el8", Arch: "x86_64",
					SrcName: "glibc", SrcEpoch: 0, SrcVersion: "2.28", SrcRelease: "101.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+ and LGPLv2+ with exceptions and GPLv2+ and GPLv2+ with exceptions and BSD and Inner-Net and ISC and Public Domain and GFDL"},
				},
				{
					Name: "libacl", Epoch: 0, Version: "2.2.53", Release: "1.el8", Arch: "x86_64", SrcName: "acl",
					SrcEpoch: 0, SrcVersion: "2.2.53", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "perl-Thread-Queue", Epoch: 0, Version: "3.13", Release: "1.el8", Arch: "noarch",
					SrcName: "perl-Thread-Queue", SrcEpoch: 0, SrcVersion: "3.13", SrcRelease: "1.el8",
					Modularitylabel: "", Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "sed", Epoch: 0, Version: "4.5", Release: "1.el8", Arch: "x86_64", SrcName: "sed",
					SrcEpoch: 0, SrcVersion: "4.5", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "isl", Epoch: 0, Version: "0.16.1", Release: "6.el8", Arch: "x86_64", SrcName: "isl",
					SrcEpoch: 0, SrcVersion: "0.16.1", SrcRelease: "6.el8", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "libmount", Epoch: 0, Version: "2.32.1", Release: "22.el8", Arch: "x86_64",
					SrcName: "util-linux", SrcEpoch: 0, SrcVersion: "2.32.1", SrcRelease: "22.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "libtool", Epoch: 0, Version: "2.4.6", Release: "25.el8", Arch: "x86_64", SrcName: "libtool",
					SrcEpoch: 0, SrcVersion: "2.4.6", SrcRelease: "25.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and LGPLv2+ and GFDL"},
				},
				{
					Name: "audit-libs", Epoch: 0, Version: "3.0", Release: "0.17.20191104git1c2f876.el8",
					Arch: "x86_64", SrcName: "audit", SrcEpoch: 0, SrcVersion: "3.0",
					SrcRelease: "0.17.20191104git1c2f876.el8", Modularitylabel: "", Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "libgcrypt-devel", Epoch: 0, Version: "1.8.3", Release: "4.el8", Arch: "x86_64",
					SrcName: "libgcrypt", SrcEpoch: 0, SrcVersion: "1.8.3", SrcRelease: "4.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+ and GPLv2+"},
				},
				{
					Name: "libsmartcols", Epoch: 0, Version: "2.32.1", Release: "22.el8", Arch: "x86_64",
					SrcName: "util-linux", SrcEpoch: 0, SrcVersion: "2.32.1", SrcRelease: "22.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "nodejs-full-i18n", Epoch: 1, Version: "10.21.0", Release: "3.module_el8.2.0+391+8da3adc6",
					Arch: "x86_64", SrcName: "nodejs", SrcEpoch: 1, SrcVersion: "10.21.0",
					SrcRelease:      "3.module_el8.2.0+391+8da3adc6",
					Modularitylabel: "nodejs:10:8020020200707141642:6a468ee4",
					Licenses:        []string{"MIT and ASL 2.0 and ISC and BSD"},
				},
				{
					Name: "lua-libs", Epoch: 0, Version: "5.3.4", Release: "11.el8", Arch: "x86_64", SrcName: "lua",
					SrcEpoch: 0, SrcVersion: "5.3.4", SrcRelease: "11.el8", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "nodejs", Epoch: 1, Version: "10.21.0", Release: "3.module_el8.2.0+391+8da3adc6",
					Arch: "x86_64", SrcName: "nodejs", SrcEpoch: 1, SrcVersion: "10.21.0",
					SrcRelease:      "3.module_el8.2.0+391+8da3adc6",
					Modularitylabel: "nodejs:10:8020020200707141642:6a468ee4",
					Licenses:        []string{"MIT and ASL 2.0 and ISC and BSD"},
				},
				{
					Name: "p11-kit", Epoch: 0, Version: "0.23.14", Release: "5.el8_0", Arch: "x86_64",
					SrcName: "p11-kit", SrcEpoch: 0, SrcVersion: "0.23.14", SrcRelease: "5.el8_0", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "libbabeltrace", Epoch: 0, Version: "1.5.4", Release: "3.el8", Arch: "x86_64",
					SrcName: "babeltrace", SrcEpoch: 0, SrcVersion: "1.5.4", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"MIT and GPLv2"},
				},
				{
					Name: "gzip", Epoch: 0, Version: "1.9", Release: "9.el8", Arch: "x86_64", SrcName: "gzip",
					SrcEpoch: 0, SrcVersion: "1.9", SrcRelease: "9.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+ and GFDL"},
				},
				{
					Name: "libatomic_ops", Epoch: 0, Version: "7.6.2", Release: "3.el8", Arch: "x86_64",
					SrcName: "libatomic_ops", SrcEpoch: 0, SrcVersion: "7.6.2", SrcRelease: "3.el8",
					Modularitylabel: "", Licenses: []string{"GPLv2 and MIT"},
				},
				{
					Name: "libunistring", Epoch: 0, Version: "0.9.9", Release: "3.el8", Arch: "x86_64",
					SrcName: "libunistring", SrcEpoch: 0, SrcVersion: "0.9.9", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+ or LGPLv3+"},
				},
				{
					Name: "guile", Epoch: 5, Version: "2.0.14", Release: "7.el8", Arch: "x86_64", SrcName: "guile",
					SrcEpoch: 5, SrcVersion: "2.0.14", SrcRelease: "7.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv3+"},
				},
				{
					Name: "libassuan", Epoch: 0, Version: "2.5.1", Release: "3.el8", Arch: "x86_64",
					SrcName: "libassuan", SrcEpoch: 0, SrcVersion: "2.5.1", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+ and GPLv3+"},
				},
				{
					Name: "gdb", Epoch: 0, Version: "8.2", Release: "12.el8", Arch: "x86_64", SrcName: "gdb",
					SrcEpoch: 0, SrcVersion: "8.2", SrcRelease: "12.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+ and GPLv3+ with exceptions and GPLv2+ and GPLv2+ with exceptions and GPL+ and LGPLv2+ and LGPLv3+ and BSD and Public Domain and GFDL"},
				},
				{
					Name: "gdbm-libs", Epoch: 1, Version: "1.18", Release: "1.el8", Arch: "x86_64", SrcName: "gdbm",
					SrcEpoch: 1, SrcVersion: "1.18", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "platform-python-setuptools", Epoch: 0, Version: "39.2.0", Release: "6.el8", Arch: "noarch",
					SrcName: "python-setuptools", SrcEpoch: 0, SrcVersion: "39.2.0", SrcRelease: "6.el8",
					Modularitylabel: "", Licenses: []string{"MIT"},
				},
				{
					Name: "libtasn1", Epoch: 0, Version: "4.13", Release: "3.el8", Arch: "x86_64", SrcName: "libtasn1",
					SrcEpoch: 0, SrcVersion: "4.13", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+ and LGPLv2+"},
				},
				{
					Name: "python3-setuptools", Epoch: 0, Version: "39.2.0", Release: "6.el8", Arch: "noarch",
					SrcName: "python-setuptools", SrcEpoch: 0, SrcVersion: "39.2.0", SrcRelease: "6.el8",
					Modularitylabel: "", Licenses: []string{"MIT"},
				},
				{
					Name: "lzo", Epoch: 0, Version: "2.08", Release: "14.el8", Arch: "x86_64", SrcName: "lzo",
					SrcEpoch: 0, SrcVersion: "2.08", SrcRelease: "14.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+"},
				},
				{
					Name: "python3-pip", Epoch: 0, Version: "9.0.3", Release: "18.el8", Arch: "noarch",
					SrcName: "python-pip", SrcEpoch: 0, SrcVersion: "9.0.3", SrcRelease: "18.el8", Modularitylabel: "",
					Licenses: []string{"MIT and Python and ASL 2.0 and BSD and ISC and LGPLv2 and MPLv2.0 and (ASL 2.0 or BSD)"},
				},
				{
					Name: "grep", Epoch: 0, Version: "3.1", Release: "6.el8", Arch: "x86_64", SrcName: "grep",
					SrcEpoch: 0, SrcVersion: "3.1", SrcRelease: "6.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "python2-pip-wheel", Epoch: 0, Version: "9.0.3", Release: "18.module_el8.3.0+478+7570e00c",
					Arch: "noarch", SrcName: "python2-pip", SrcEpoch: 0, SrcVersion: "9.0.3",
					SrcRelease:      "18.module_el8.3.0+478+7570e00c",
					Modularitylabel: "python27:2.7:8030020200831201838:851f4228",
					Licenses:        []string{"MIT and Python and ASL 2.0 and BSD and ISC and LGPLv2 and MPLv2.0 and (ASL 2.0 or BSD)"},
				},
				{
					Name: "dbus-libs", Epoch: 1, Version: "1.12.8", Release: "10.el8_2", Arch: "x86_64",
					SrcName: "dbus", SrcEpoch: 1, SrcVersion: "1.12.8", SrcRelease: "10.el8_2", Modularitylabel: "",
					Licenses: []string{"(GPLv2+ or AFL) and GPLv2+"},
				},
				{
					Name: "python2-pip", Epoch: 0, Version: "9.0.3", Release: "18.module_el8.3.0+478+7570e00c",
					Arch: "noarch", SrcName: "python2-pip", SrcEpoch: 0, SrcVersion: "9.0.3",
					SrcRelease:      "18.module_el8.3.0+478+7570e00c",
					Modularitylabel: "python27:2.7:8030020200831201838:851f4228",
					Licenses:        []string{"MIT and Python and ASL 2.0 and BSD and ISC and LGPLv2 and MPLv2.0 and (ASL 2.0 or BSD)"},
				},
				{
					Name: "dhcp-libs", Epoch: 12, Version: "4.3.6", Release: "40.el8", Arch: "x86_64", SrcName: "dhcp",
					SrcEpoch: 12, SrcVersion: "4.3.6", SrcRelease: "40.el8", Modularitylabel: "",
					Licenses: []string{"ISC"},
				},
				{
					Name: "python2", Epoch: 0, Version: "2.7.17", Release: "2.module_el8.3.0+478+7570e00c",
					Arch: "x86_64", SrcName: "python2", SrcEpoch: 0, SrcVersion: "2.7.17",
					SrcRelease:      "2.module_el8.3.0+478+7570e00c",
					Modularitylabel: "python27:2.7:8030020200831201838:851f4228", Licenses: []string{"Python"},
				},
				{
					Name: "procps-ng", Epoch: 0, Version: "3.3.15", Release: "1.el8", Arch: "x86_64",
					SrcName: "procps-ng", SrcEpoch: 0, SrcVersion: "3.3.15", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"GPL+ and GPLv2 and GPLv2+ and GPLv3+ and LGPLv2+"},
				},
				{
					Name: "python2-rpmUtils", Epoch: 0, Version: "0.1", Release: "1.el8", Arch: "noarch",
					SrcName: "python-rpmUtils", SrcEpoch: 0, SrcVersion: "0.1", SrcRelease: "1.el8",
					Modularitylabel: "", Licenses: []string{"GPLv2+"},
				},
				{
					Name: "xz", Epoch: 0, Version: "5.2.4", Release: "3.el8", Arch: "x86_64", SrcName: "xz",
					SrcEpoch: 0, SrcVersion: "5.2.4", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and Public Domain"},
				},
				{
					Name: "rpm", Epoch: 0, Version: "4.14.3", Release: "4.el8", Arch: "x86_64", SrcName: "rpm",
					SrcEpoch: 0, SrcVersion: "4.14.3", SrcRelease: "4.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+"},
				},
				{
					Name: "gdbm", Epoch: 1, Version: "1.18", Release: "1.el8", Arch: "x86_64", SrcName: "gdbm",
					SrcEpoch: 1, SrcVersion: "1.18", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "python3-rpm", Epoch: 0, Version: "4.14.3", Release: "4.el8", Arch: "x86_64", SrcName: "rpm",
					SrcEpoch: 0, SrcVersion: "4.14.3", SrcRelease: "4.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+"},
				},
				{
					Name: "shadow-utils", Epoch: 2, Version: "4.6", Release: "8.el8", Arch: "x86_64",
					SrcName: "shadow-utils", SrcEpoch: 2, SrcVersion: "4.6", SrcRelease: "8.el8", Modularitylabel: "",
					Licenses: []string{"BSD and GPLv2+"},
				},
				{
					Name: "libfdisk", Epoch: 0, Version: "2.32.1", Release: "22.el8", Arch: "x86_64",
					SrcName: "util-linux", SrcEpoch: 0, SrcVersion: "2.32.1", SrcRelease: "22.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "mpfr", Epoch: 0, Version: "3.1.6", Release: "1.el8", Arch: "x86_64", SrcName: "mpfr",
					SrcEpoch: 0, SrcVersion: "3.1.6", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv3+ and GPLv3+ and GFDL"},
				},
				{
					Name: "snappy", Epoch: 0, Version: "1.1.7", Release: "5.el8", Arch: "x86_64", SrcName: "snappy",
					SrcEpoch: 0, SrcVersion: "1.1.7", SrcRelease: "5.el8", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "libmetalink", Epoch: 0, Version: "0.1.3", Release: "7.el8", Arch: "x86_64",
					SrcName: "libmetalink", SrcEpoch: 0, SrcVersion: "0.1.3", SrcRelease: "7.el8", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "libksba", Epoch: 0, Version: "1.3.5", Release: "7.el8", Arch: "x86_64", SrcName: "libksba",
					SrcEpoch: 0, SrcVersion: "1.3.5", SrcRelease: "7.el8", Modularitylabel: "",
					Licenses: []string{"(LGPLv3+ or GPLv2+) and GPLv3+"},
				},
				{
					Name: "ethtool", Epoch: 2, Version: "5.0", Release: "2.el8", Arch: "x86_64", SrcName: "ethtool",
					SrcEpoch: 2, SrcVersion: "5.0", SrcRelease: "2.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2"},
				},
				{
					Name: "libmnl", Epoch: 0, Version: "1.0.4", Release: "6.el8", Arch: "x86_64", SrcName: "libmnl",
					SrcEpoch: 0, SrcVersion: "1.0.4", SrcRelease: "6.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "libpcap", Epoch: 14, Version: "1.9.0", Release: "3.el8", Arch: "x86_64", SrcName: "libpcap",
					SrcEpoch: 14, SrcVersion: "1.9.0", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"BSD with advertising"},
				},
				{
					Name: "libseccomp", Epoch: 0, Version: "2.4.1", Release: "1.el8", Arch: "x86_64",
					SrcName: "libseccomp", SrcEpoch: 0, SrcVersion: "2.4.1", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2"},
				},
				{
					Name: "gawk", Epoch: 0, Version: "4.2.1", Release: "1.el8", Arch: "x86_64", SrcName: "gawk",
					SrcEpoch: 0, SrcVersion: "4.2.1", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+ and GPLv2+ and LGPLv2+ and BSD"},
				},
				{
					Name: "libnsl2", Epoch: 0, Version: "1.2.0", Release: "2.20180605git4a062cf.el8", Arch: "x86_64",
					SrcName: "libnsl2", SrcEpoch: 0, SrcVersion: "1.2.0", SrcRelease: "2.20180605git4a062cf.el8",
					Modularitylabel: "", Licenses: []string{"BSD and LGPLv2+"},
				},
				{
					Name: "krb5-libs", Epoch: 0, Version: "1.17", Release: "18.el8", Arch: "x86_64", SrcName: "krb5",
					SrcEpoch: 0, SrcVersion: "1.17", SrcRelease: "18.el8", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "crypto-policies", Epoch: 0, Version: "20191128", Release: "2.git23e1bf1.el8", Arch: "noarch",
					SrcName: "crypto-policies", SrcEpoch: 0, SrcVersion: "20191128", SrcRelease: "2.git23e1bf1.el8",
					Modularitylabel: "", Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "platform-python", Epoch: 0, Version: "3.6.8", Release: "23.el8", Arch: "x86_64",
					SrcName: "python3", SrcEpoch: 0, SrcVersion: "3.6.8", SrcRelease: "23.el8", Modularitylabel: "",
					Licenses: []string{"Python"},
				},
				{
					Name: "libdb", Epoch: 0, Version: "5.3.28", Release: "37.el8", Arch: "x86_64", SrcName: "libdb",
					SrcEpoch: 0, SrcVersion: "5.3.28", SrcRelease: "37.el8", Modularitylabel: "",
					Licenses: []string{"BSD and LGPLv2 and Sleepycat"},
				},
				{
					Name: "pam", Epoch: 0, Version: "1.3.1", Release: "8.el8", Arch: "x86_64", SrcName: "pam",
					SrcEpoch: 0, SrcVersion: "1.3.1", SrcRelease: "8.el8", Modularitylabel: "",
					Licenses: []string{"BSD and GPLv2+"},
				},
				{
					Name: "gnutls", Epoch: 0, Version: "3.6.8", Release: "11.el8_2", Arch: "x86_64", SrcName: "gnutls",
					SrcEpoch: 0, SrcVersion: "3.6.8", SrcRelease: "11.el8_2", Modularitylabel: "",
					Licenses: []string{"GPLv3+ and LGPLv2+"},
				},
				{
					Name: "kmod-libs", Epoch: 0, Version: "25", Release: "16.el8", Arch: "x86_64", SrcName: "kmod",
					SrcEpoch: 0, SrcVersion: "25", SrcRelease: "16.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "ima-evm-utils", Epoch: 0, Version: "1.1", Release: "5.el8", Arch: "x86_64",
					SrcName: "ima-evm-utils", SrcEpoch: 0, SrcVersion: "1.1", SrcRelease: "5.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2"},
				},
				{
					Name: "libcurl-minimal", Epoch: 0, Version: "7.61.1", Release: "12.el8", Arch: "x86_64",
					SrcName: "curl", SrcEpoch: 0, SrcVersion: "7.61.1", SrcRelease: "12.el8", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "cyrus-sasl-lib", Epoch: 0, Version: "2.1.27", Release: "1.el8", Arch: "x86_64",
					SrcName: "cyrus-sasl", SrcEpoch: 0, SrcVersion: "2.1.27", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"BSD with advertising"},
				},
				{
					Name: "libdb-utils", Epoch: 0, Version: "5.3.28", Release: "37.el8", Arch: "x86_64",
					SrcName: "libdb", SrcEpoch: 0, SrcVersion: "5.3.28", SrcRelease: "37.el8", Modularitylabel: "",
					Licenses: []string{"BSD and LGPLv2 and Sleepycat"},
				},
				{
					Name: "libsolv", Epoch: 0, Version: "0.7.7", Release: "1.el8", Arch: "x86_64", SrcName: "libsolv",
					SrcEpoch: 0, SrcVersion: "0.7.7", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "libmodulemd1", Epoch: 0, Version: "1.8.16", Release: "0.2.8.2.1", Arch: "x86_64",
					SrcName: "libmodulemd", SrcEpoch: 0, SrcVersion: "2.8.2", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "gnupg2", Epoch: 0, Version: "2.2.9", Release: "1.el8", Arch: "x86_64", SrcName: "gnupg2",
					SrcEpoch: 0, SrcVersion: "2.2.9", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "python3-libdnf", Epoch: 0, Version: "0.39.1", Release: "6.el8_2", Arch: "x86_64",
					SrcName: "libdnf", SrcEpoch: 0, SrcVersion: "0.39.1", SrcRelease: "6.el8_2", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "python3-gpg", Epoch: 0, Version: "1.10.0", Release: "6.el8.0.1", Arch: "x86_64",
					SrcName: "gpgme", SrcEpoch: 0, SrcVersion: "1.10.0", SrcRelease: "6.el8.0.1", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "dnf-data", Epoch: 0, Version: "4.2.17", Release: "7.el8_2", Arch: "noarch", SrcName: "dnf",
					SrcEpoch: 0, SrcVersion: "4.2.17", SrcRelease: "7.el8_2", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and GPLv2 and GPL"},
				},
				{
					Name: "dbus-common", Epoch: 1, Version: "1.12.8", Release: "10.el8_2", Arch: "noarch",
					SrcName: "dbus", SrcEpoch: 1, SrcVersion: "1.12.8", SrcRelease: "10.el8_2", Modularitylabel: "",
					Licenses: []string{"(GPLv2+ or AFL) and GPLv2+"},
				},
				{
					Name: "device-mapper", Epoch: 8, Version: "1.02.169", Release: "3.el8", Arch: "x86_64",
					SrcName: "lvm2", SrcEpoch: 8, SrcVersion: "2.03.08", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2"},
				},
				{
					Name: "cryptsetup-libs", Epoch: 0, Version: "2.2.2", Release: "1.el8", Arch: "x86_64",
					SrcName: "cryptsetup", SrcEpoch: 0, SrcVersion: "2.2.2", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and LGPLv2+"},
				},
				{
					Name: "elfutils-libs", Epoch: 0, Version: "0.178", Release: "7.el8", Arch: "x86_64",
					SrcName: "elfutils", SrcEpoch: 0, SrcVersion: "0.178", SrcRelease: "7.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+ or LGPLv3+"},
				},
				{
					Name: "systemd", Epoch: 0, Version: "239", Release: "31.el8_2.2", Arch: "x86_64",
					SrcName: "systemd", SrcEpoch: 0, SrcVersion: "239", SrcRelease: "31.el8_2.2", Modularitylabel: "",
					Licenses: []string{"LGPLv2+ and MIT and GPLv2+"},
				},
				{
					Name: "iputils", Epoch: 0, Version: "20180629", Release: "2.el8", Arch: "x86_64",
					SrcName: "iputils", SrcEpoch: 0, SrcVersion: "20180629", SrcRelease: "2.el8", Modularitylabel: "",
					Licenses: []string{"BSD and GPLv2+"},
				},
				{
					Name: "libkcapi", Epoch: 0, Version: "1.1.1", Release: "16_1.el8", Arch: "x86_64",
					SrcName: "libkcapi", SrcEpoch: 0, SrcVersion: "1.1.1", SrcRelease: "16_1.el8", Modularitylabel: "",
					Licenses: []string{"BSD or GPLv2"},
				},
				{
					Name: "systemd-udev", Epoch: 0, Version: "239", Release: "31.el8_2.2", Arch: "x86_64",
					SrcName: "systemd", SrcEpoch: 0, SrcVersion: "239", SrcRelease: "31.el8_2.2", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "dracut-network", Epoch: 0, Version: "049", Release: "70.git20200228.el8", Arch: "x86_64",
					SrcName: "dracut", SrcEpoch: 0, SrcVersion: "049", SrcRelease: "70.git20200228.el8",
					Modularitylabel: "", Licenses: []string{"GPLv2+ and LGPLv2+"},
				},
				{
					Name: "python3-dnf", Epoch: 0, Version: "4.2.17", Release: "7.el8_2", Arch: "noarch",
					SrcName: "dnf", SrcEpoch: 0, SrcVersion: "4.2.17", SrcRelease: "7.el8_2", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and GPLv2 and GPL"},
				},
				{
					Name: "yum", Epoch: 0, Version: "4.2.17", Release: "7.el8_2", Arch: "noarch", SrcName: "dnf",
					SrcEpoch: 0, SrcVersion: "4.2.17", SrcRelease: "7.el8_2", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and GPLv2 and GPL"},
				},
				{
					Name: "binutils", Epoch: 0, Version: "2.30", Release: "73.el8", Arch: "x86_64", SrcName: "binutils",
					SrcEpoch: 0, SrcVersion: "2.30", SrcRelease: "73.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "vim-minimal", Epoch: 2, Version: "8.0.1763", Release: "13.el8", Arch: "x86_64",
					SrcName: "vim", SrcEpoch: 2, SrcVersion: "8.0.1763", SrcRelease: "13.el8", Modularitylabel: "",
					Licenses: []string{"Vim and MIT"},
				},
				{
					Name: "less", Epoch: 0, Version: "530", Release: "1.el8", Arch: "x86_64", SrcName: "less",
					SrcEpoch: 0, SrcVersion: "530", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+ or BSD"},
				},
				{
					Name: "rootfiles", Epoch: 0, Version: "8.1", Release: "22.el8", Arch: "noarch",
					SrcName: "rootfiles", SrcEpoch: 0, SrcVersion: "8.1", SrcRelease: "22.el8", Modularitylabel: "",
					Licenses: []string{"Public Domain"},
				},
				{
					Name: "centos-gpg-keys", Epoch: 0, Version: "8.2", Release: "2.2004.0.2.el8", Arch: "noarch",
					SrcName: "centos-release", SrcEpoch: 0, SrcVersion: "8.2", SrcRelease: "2.2004.0.2.el8",
					Modularitylabel: "", Licenses: []string{"GPLv2"},
				},
				{
					Name: "centos-repos", Epoch: 0, Version: "8.2", Release: "2.2004.0.2.el8", Arch: "x86_64",
					SrcName: "centos-release", SrcEpoch: 0, SrcVersion: "8.2", SrcRelease: "2.2004.0.2.el8",
					Modularitylabel: "", Licenses: []string{"GPLv2"},
				},
				{
					Name: "tzdata", Epoch: 0, Version: "2020d", Release: "1.el8", Arch: "noarch", SrcName: "tzdata",
					SrcEpoch: 0, SrcVersion: "2020d", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"Public Domain"},
				},
				{
					Name: "ca-certificates", Epoch: 0, Version: "2020.2.41", Release: "80.0.el8_2", Arch: "noarch",
					SrcName: "ca-certificates", SrcEpoch: 0, SrcVersion: "2020.2.41", SrcRelease: "80.0.el8_2",
					Modularitylabel: "", Licenses: []string{"Public Domain"},
				},
				{
					Name: "perl-Exporter", Epoch: 0, Version: "5.72", Release: "396.el8", Arch: "noarch",
					SrcName: "perl-Exporter", SrcEpoch: 0, SrcVersion: "5.72", SrcRelease: "396.el8",
					Modularitylabel: "", Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "perl-Carp", Epoch: 0, Version: "1.42", Release: "396.el8", Arch: "noarch",
					SrcName: "perl-Carp", SrcEpoch: 0, SrcVersion: "1.42", SrcRelease: "396.el8", Modularitylabel: "",
					Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "perl-parent", Epoch: 1, Version: "0.237", Release: "1.el8", Arch: "noarch",
					SrcName: "perl-parent", SrcEpoch: 1, SrcVersion: "0.237", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "nss-util", Epoch: 0, Version: "3.53.1", Release: "11.el8_2", Arch: "x86_64", SrcName: "nss",
					SrcEpoch: 0, SrcVersion: "3.53.1", SrcRelease: "11.el8_2", Modularitylabel: "",
					Licenses: []string{"MPLv2.0"},
				},
				{
					Name: "nss-softokn", Epoch: 0, Version: "3.53.1", Release: "11.el8_2", Arch: "x86_64",
					SrcName: "nss", SrcEpoch: 0, SrcVersion: "3.53.1", SrcRelease: "11.el8_2", Modularitylabel: "",
					Licenses: []string{"MPLv2.0"},
				},
				{
					Name: "nss-sysinit", Epoch: 0, Version: "3.53.1", Release: "11.el8_2", Arch: "x86_64",
					SrcName: "nss", SrcEpoch: 0, SrcVersion: "3.53.1", SrcRelease: "11.el8_2", Modularitylabel: "",
					Licenses: []string{"MPLv2.0"},
				},
				{
					Name: "nss-softokn-freebl-devel", Epoch: 0, Version: "3.53.1", Release: "11.el8_2", Arch: "x86_64",
					SrcName: "nss", SrcEpoch: 0, SrcVersion: "3.53.1", SrcRelease: "11.el8_2", Modularitylabel: "",
					Licenses: []string{"MPLv2.0"},
				},
				{
					Name: "perl-macros", Epoch: 4, Version: "5.26.3", Release: "416.el8", Arch: "x86_64",
					SrcName: "perl", SrcEpoch: 4, SrcVersion: "5.26.3", SrcRelease: "416.el8", Modularitylabel: "",
					Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "perl-Socket", Epoch: 4, Version: "2.027", Release: "3.el8", Arch: "x86_64",
					SrcName: "perl-Socket", SrcEpoch: 4, SrcVersion: "2.027", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "perl-Unicode-Normalize", Epoch: 0, Version: "1.25", Release: "396.el8", Arch: "x86_64",
					SrcName: "perl-Unicode-Normalize", SrcEpoch: 0, SrcVersion: "1.25", SrcRelease: "396.el8",
					Modularitylabel: "", Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "perl-IO", Epoch: 0, Version: "1.38", Release: "416.el8", Arch: "x86_64", SrcName: "perl",
					SrcEpoch: 0, SrcVersion: "5.26.3", SrcRelease: "416.el8", Modularitylabel: "",
					Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "perl-constant", Epoch: 0, Version: "1.33", Release: "396.el8", Arch: "noarch",
					SrcName: "perl-constant", SrcEpoch: 0, SrcVersion: "1.33", SrcRelease: "396.el8",
					Modularitylabel: "", Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "perl-threads-shared", Epoch: 0, Version: "1.58", Release: "2.el8", Arch: "x86_64",
					SrcName: "perl-threads-shared", SrcEpoch: 0, SrcVersion: "1.58", SrcRelease: "2.el8",
					Modularitylabel: "", Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "perl-MIME-Base64", Epoch: 0, Version: "3.15", Release: "396.el8", Arch: "x86_64",
					SrcName: "perl-MIME-Base64", SrcEpoch: 0, SrcVersion: "3.15", SrcRelease: "396.el8",
					Modularitylabel: "", Licenses: []string{"(GPL+ or Artistic) and MIT"},
				},
				{
					Name: "perl-Time-Local", Epoch: 1, Version: "1.280", Release: "1.el8", Arch: "noarch",
					SrcName: "perl-Time-Local", SrcEpoch: 1, SrcVersion: "1.280", SrcRelease: "1.el8",
					Modularitylabel: "", Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "perl-Digest", Epoch: 0, Version: "1.17", Release: "395.el8", Arch: "noarch",
					SrcName: "perl-Digest", SrcEpoch: 0, SrcVersion: "1.17", SrcRelease: "395.el8", Modularitylabel: "",
					Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "perl-Net-SSLeay", Epoch: 0, Version: "1.88", Release: "1.el8", Arch: "x86_64",
					SrcName: "perl-Net-SSLeay", SrcEpoch: 0, SrcVersion: "1.88", SrcRelease: "1.el8",
					Modularitylabel: "", Licenses: []string{"Artistic 2.0"},
				},
				{
					Name: "perl-TermReadKey", Epoch: 0, Version: "2.37", Release: "7.el8", Arch: "x86_64",
					SrcName: "perl-TermReadKey", SrcEpoch: 0, SrcVersion: "2.37", SrcRelease: "7.el8",
					Modularitylabel: "", Licenses: []string{"(Copyright only) and (Artistic or GPL+)"},
				},
				{
					Name: "perl-Pod-Escapes", Epoch: 1, Version: "1.07", Release: "395.el8", Arch: "noarch",
					SrcName: "perl-Pod-Escapes", SrcEpoch: 1, SrcVersion: "1.07", SrcRelease: "395.el8",
					Modularitylabel: "", Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "perl-Mozilla-CA", Epoch: 0, Version: "20160104", Release: "7.el8", Arch: "noarch",
					SrcName: "perl-Mozilla-CA", SrcEpoch: 0, SrcVersion: "20160104", SrcRelease: "7.el8",
					Modularitylabel: "", Licenses: []string{"MPLv2.0"},
				},
				{
					Name: "fipscheck", Epoch: 0, Version: "1.5.0", Release: "4.el8", Arch: "x86_64",
					SrcName: "fipscheck", SrcEpoch: 0, SrcVersion: "1.5.0", SrcRelease: "4.el8", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "which", Epoch: 0, Version: "2.21", Release: "12.el8", Arch: "x86_64", SrcName: "which",
					SrcEpoch: 0, SrcVersion: "2.21", SrcRelease: "12.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3"},
				},
				{
					Name: "libpsl", Epoch: 0, Version: "0.20.2", Release: "5.el8", Arch: "x86_64", SrcName: "libpsl",
					SrcEpoch: 0, SrcVersion: "0.20.2", SrcRelease: "5.el8", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "pcre2-utf32", Epoch: 0, Version: "10.32", Release: "1.el8", Arch: "x86_64", SrcName: "pcre2",
					SrcEpoch: 0, SrcVersion: "10.32", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "openssl", Epoch: 1, Version: "1.1.1c", Release: "15.el8", Arch: "x86_64", SrcName: "openssl",
					SrcEpoch: 1, SrcVersion: "1.1.1c", SrcRelease: "15.el8", Modularitylabel: "",
					Licenses: []string{"OpenSSL"},
				},
				{
					Name: "perl-Term-Cap", Epoch: 0, Version: "1.17", Release: "395.el8", Arch: "noarch",
					SrcName: "perl-Term-Cap", SrcEpoch: 0, SrcVersion: "1.17", SrcRelease: "395.el8",
					Modularitylabel: "", Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "libpkgconf", Epoch: 0, Version: "1.4.2", Release: "1.el8", Arch: "x86_64",
					SrcName: "pkgconf", SrcEpoch: 0, SrcVersion: "1.4.2", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"ISC"},
				},
				{
					Name: "pkgconf-pkg-config", Epoch: 0, Version: "1.4.2", Release: "1.el8", Arch: "x86_64",
					SrcName: "pkgconf", SrcEpoch: 0, SrcVersion: "1.4.2", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"ISC"},
				},
				{
					Name: "nss-util-devel", Epoch: 0, Version: "3.53.1", Release: "11.el8_2", Arch: "x86_64",
					SrcName: "nss", SrcEpoch: 0, SrcVersion: "3.53.1", SrcRelease: "11.el8_2", Modularitylabel: "",
					Licenses: []string{"MPLv2.0"},
				},
				{
					Name: "libcom_err-devel", Epoch: 0, Version: "1.45.4", Release: "3.el8", Arch: "x86_64",
					SrcName: "e2fsprogs", SrcEpoch: 0, SrcVersion: "1.45.4", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "libverto-devel", Epoch: 0, Version: "0.3.0", Release: "5.el8", Arch: "x86_64",
					SrcName: "libverto", SrcEpoch: 0, SrcVersion: "0.3.0", SrcRelease: "5.el8", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "libselinux-devel", Epoch: 0, Version: "2.9", Release: "3.el8", Arch: "x86_64",
					SrcName: "libselinux", SrcEpoch: 0, SrcVersion: "2.9", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"Public Domain"},
				},
				{
					Name: "libkadm5", Epoch: 0, Version: "1.17", Release: "18.el8", Arch: "x86_64", SrcName: "krb5",
					SrcEpoch: 0, SrcVersion: "1.17", SrcRelease: "18.el8", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "openssh-clients", Epoch: 0, Version: "8.0p1", Release: "4.el8_1", Arch: "x86_64",
					SrcName: "openssh", SrcEpoch: 0, SrcVersion: "8.0p1", SrcRelease: "4.el8_1", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "git-core-doc", Epoch: 0, Version: "2.18.4", Release: "2.el8_2", Arch: "noarch",
					SrcName: "git", SrcEpoch: 0, SrcVersion: "2.18.4", SrcRelease: "2.el8_2", Modularitylabel: "",
					Licenses: []string{"GPLv2"},
				},
				{
					Name: "krb5-devel", Epoch: 0, Version: "1.17", Release: "18.el8", Arch: "x86_64", SrcName: "krb5",
					SrcEpoch: 0, SrcVersion: "1.17", SrcRelease: "18.el8", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "perl-Encode", Epoch: 4, Version: "2.97", Release: "3.el8", Arch: "x86_64",
					SrcName: "perl-Encode", SrcEpoch: 4, SrcVersion: "2.97", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"(GPL+ or Artistic) and Artistic 2.0 and UCD"},
				},
				{
					Name: "perl-Getopt-Long", Epoch: 1, Version: "2.50", Release: "4.el8", Arch: "noarch",
					SrcName: "perl-Getopt-Long", SrcEpoch: 1, SrcVersion: "2.50", SrcRelease: "4.el8",
					Modularitylabel: "", Licenses: []string{"GPLv2+ or Artistic"},
				},
				{
					Name: "libgcc", Epoch: 0, Version: "8.3.1", Release: "5.el8.0.2", Arch: "x86_64", SrcName: "gcc",
					SrcEpoch: 0, SrcVersion: "8.3.1", SrcRelease: "5.el8.0.2", Modularitylabel: "",
					Licenses: []string{"GPLv3+ and GPLv3+ with exceptions and GPLv2+ with exceptions and LGPLv2+ and BSD"},
				},
				{
					Name: "perl-Pod-Usage", Epoch: 4, Version: "1.69", Release: "395.el8", Arch: "noarch",
					SrcName: "perl-Pod-Usage", SrcEpoch: 4, SrcVersion: "1.69", SrcRelease: "395.el8",
					Modularitylabel: "", Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "python3-pip-wheel", Epoch: 0, Version: "9.0.3", Release: "16.el8", Arch: "noarch",
					SrcName: "python-pip", SrcEpoch: 0, SrcVersion: "9.0.3", SrcRelease: "16.el8", Modularitylabel: "",
					Licenses: []string{"MIT and Python and ASL 2.0 and BSD and ISC and LGPLv2 and MPLv2.0 and (ASL 2.0 or BSD)"},
				},
				{
					Name: "perl-HTTP-Tiny", Epoch: 0, Version: "0.074", Release: "1.el8", Arch: "noarch",
					SrcName: "perl-HTTP-Tiny", SrcEpoch: 0, SrcVersion: "0.074", SrcRelease: "1.el8",
					Modularitylabel: "", Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "perl-libnet", Epoch: 0, Version: "3.11", Release: "3.el8", Arch: "noarch",
					SrcName: "perl-libnet", SrcEpoch: 0, SrcVersion: "3.11", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "setup", Epoch: 0, Version: "2.12.2", Release: "5.el8", Arch: "noarch", SrcName: "setup",
					SrcEpoch: 0, SrcVersion: "2.12.2", SrcRelease: "5.el8", Modularitylabel: "",
					Licenses: []string{"Public Domain"},
				},
				{
					Name: "file", Epoch: 0, Version: "5.33", Release: "13.el8", Arch: "x86_64", SrcName: "file",
					SrcEpoch: 0, SrcVersion: "5.33", SrcRelease: "13.el8", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "basesystem", Epoch: 0, Version: "11", Release: "5.el8", Arch: "noarch",
					SrcName: "basesystem", SrcEpoch: 0, SrcVersion: "11", SrcRelease: "5.el8", Modularitylabel: "",
					Licenses: []string{"Public Domain"},
				},
				{
					Name: "perl-Git", Epoch: 0, Version: "2.18.4", Release: "2.el8_2", Arch: "noarch", SrcName: "git",
					SrcEpoch: 0, SrcVersion: "2.18.4", SrcRelease: "2.el8_2", Modularitylabel: "",
					Licenses: []string{"GPLv2"},
				},
				{
					Name: "ncurses-base", Epoch: 0, Version: "6.1", Release: "7.20180224.el8", Arch: "noarch",
					SrcName: "ncurses", SrcEpoch: 0, SrcVersion: "6.1", SrcRelease: "7.20180224.el8",
					Modularitylabel: "", Licenses: []string{"MIT"},
				},
				{
					Name: "vim-filesystem", Epoch: 2, Version: "8.0.1763", Release: "13.el8", Arch: "noarch",
					SrcName: "vim", SrcEpoch: 2, SrcVersion: "8.0.1763", SrcRelease: "13.el8", Modularitylabel: "",
					Licenses: []string{"Vim and MIT"},
				},
				{
					Name: "libselinux", Epoch: 0, Version: "2.9", Release: "3.el8", Arch: "x86_64",
					SrcName: "libselinux", SrcEpoch: 0, SrcVersion: "2.9", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"Public Domain"},
				},
				{
					Name: "gpm-libs", Epoch: 0, Version: "1.20.7", Release: "15.el8", Arch: "x86_64", SrcName: "gpm",
					SrcEpoch: 0, SrcVersion: "1.20.7", SrcRelease: "15.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2 and GPLv2+ with exceptions and GPLv3+ and Verbatim and Copyright only"},
				},
				{
					Name: "glibc-minimal-langpack", Epoch: 0, Version: "2.28", Release: "101.el8", Arch: "x86_64",
					SrcName: "glibc", SrcEpoch: 0, SrcVersion: "2.28", SrcRelease: "101.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+ and LGPLv2+ with exceptions and GPLv2+ and GPLv2+ with exceptions and BSD and Inner-Net and ISC and Public Domain and GFDL"},
				},
				{
					Name: "file-devel", Epoch: 0, Version: "5.33", Release: "13.el8", Arch: "x86_64", SrcName: "file",
					SrcEpoch: 0, SrcVersion: "5.33", SrcRelease: "13.el8", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "glibc", Epoch: 0, Version: "2.28", Release: "101.el8", Arch: "x86_64", SrcName: "glibc",
					SrcEpoch: 0, SrcVersion: "2.28", SrcRelease: "101.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+ and LGPLv2+ with exceptions and GPLv2+ and GPLv2+ with exceptions and BSD and Inner-Net and ISC and Public Domain and GFDL"},
				},
				{
					Name: "nss-devel", Epoch: 0, Version: "3.53.1", Release: "11.el8_2", Arch: "x86_64", SrcName: "nss",
					SrcEpoch: 0, SrcVersion: "3.53.1", SrcRelease: "11.el8_2", Modularitylabel: "",
					Licenses: []string{"MPLv2.0"},
				},
				{
					Name: "libsepol", Epoch: 0, Version: "2.9", Release: "1.el8", Arch: "x86_64", SrcName: "libsepol",
					SrcEpoch: 0, SrcVersion: "2.9", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "xz-devel", Epoch: 0, Version: "5.2.4", Release: "3.el8", Arch: "x86_64", SrcName: "xz",
					SrcEpoch: 0, SrcVersion: "5.2.4", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"Public Domain"},
				},
				{
					Name: "xz-libs", Epoch: 0, Version: "5.2.4", Release: "3.el8", Arch: "x86_64", SrcName: "xz",
					SrcEpoch: 0, SrcVersion: "5.2.4", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"Public Domain"},
				},
				{
					Name: "wget", Epoch: 0, Version: "1.19.5", Release: "8.el8_1.1", Arch: "x86_64", SrcName: "wget",
					SrcEpoch: 0, SrcVersion: "1.19.5", SrcRelease: "8.el8_1.1", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "libcap", Epoch: 0, Version: "2.26", Release: "3.el8", Arch: "x86_64", SrcName: "libcap",
					SrcEpoch: 0, SrcVersion: "2.26", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2"},
				},
				{
					Name: "strace", Epoch: 0, Version: "4.24", Release: "9.el8", Arch: "x86_64", SrcName: "strace",
					SrcEpoch: 0, SrcVersion: "4.24", SrcRelease: "9.el8", Modularitylabel: "",
					Licenses: []string{"LGPL-2.1+ and GPL-2.0+"},
				},
				{
					Name: "info", Epoch: 0, Version: "6.5", Release: "6.el8", Arch: "x86_64", SrcName: "texinfo",
					SrcEpoch: 0, SrcVersion: "6.5", SrcRelease: "6.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "gdb-gdbserver", Epoch: 0, Version: "8.2", Release: "11.el8", Arch: "x86_64", SrcName: "gdb",
					SrcEpoch: 0, SrcVersion: "8.2", SrcRelease: "11.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+ and GPLv3+ with exceptions and GPLv2+ and GPLv2+ with exceptions and GPL+ and LGPLv2+ and LGPLv3+ and BSD and Public Domain and GFDL"},
				},
				{
					Name: "libcom_err", Epoch: 0, Version: "1.45.4", Release: "3.el8", Arch: "x86_64",
					SrcName: "e2fsprogs", SrcEpoch: 0, SrcVersion: "1.45.4", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "libcroco", Epoch: 0, Version: "0.6.12", Release: "4.el8_2.1", Arch: "x86_64",
					SrcName: "libcroco", SrcEpoch: 0, SrcVersion: "0.6.12", SrcRelease: "4.el8_2.1",
					Modularitylabel: "", Licenses: []string{"LGPLv2"},
				},
				{
					Name: "libxml2", Epoch: 0, Version: "2.9.7", Release: "7.el8", Arch: "x86_64", SrcName: "libxml2",
					SrcEpoch: 0, SrcVersion: "2.9.7", SrcRelease: "7.el8", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "libmpc", Epoch: 0, Version: "1.0.2", Release: "9.el8", Arch: "x86_64", SrcName: "libmpc",
					SrcEpoch: 0, SrcVersion: "1.0.2", SrcRelease: "9.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv3+ and GFDL"},
				},
				{
					Name: "expat", Epoch: 0, Version: "2.2.5", Release: "3.el8", Arch: "x86_64", SrcName: "expat",
					SrcEpoch: 0, SrcVersion: "2.2.5", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "gettext", Epoch: 0, Version: "0.19.8.1", Release: "17.el8", Arch: "x86_64",
					SrcName: "gettext", SrcEpoch: 0, SrcVersion: "0.19.8.1", SrcRelease: "17.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+ and LGPLv2+"},
				},
				{
					Name: "libuuid", Epoch: 0, Version: "2.32.1", Release: "22.el8", Arch: "x86_64",
					SrcName: "util-linux", SrcEpoch: 0, SrcVersion: "2.32.1", SrcRelease: "22.el8", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "autoconf", Epoch: 0, Version: "2.69", Release: "27.el8", Arch: "noarch", SrcName: "autoconf",
					SrcEpoch: 0, SrcVersion: "2.69", SrcRelease: "27.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and GFDL"},
				},
				{
					Name: "chkconfig", Epoch: 0, Version: "1.11", Release: "1.el8", Arch: "x86_64",
					SrcName: "chkconfig", SrcEpoch: 0, SrcVersion: "1.11", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2"},
				},
				{
					Name: "kernel-headers", Epoch: 0, Version: "4.18.0", Release: "193.28.1.el8_2", Arch: "x86_64",
					SrcName: "kernel", SrcEpoch: 0, SrcVersion: "4.18.0", SrcRelease: "193.28.1.el8_2",
					Modularitylabel: "", Licenses: []string{"GPLv2 and Redistributable, no modification permitted"},
				},
				{
					Name: "gmp", Epoch: 1, Version: "6.1.2", Release: "10.el8", Arch: "x86_64", SrcName: "gmp",
					SrcEpoch: 1, SrcVersion: "6.1.2", SrcRelease: "10.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv3+ or GPLv2+"},
				},
				{
					Name: "libxcrypt-devel", Epoch: 0, Version: "4.1.1", Release: "4.el8", Arch: "x86_64",
					SrcName: "libxcrypt", SrcEpoch: 0, SrcVersion: "4.1.1", SrcRelease: "4.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+ and BSD and Public Domain"},
				},
				{
					Name: "libattr", Epoch: 0, Version: "2.4.48", Release: "3.el8", Arch: "x86_64", SrcName: "attr",
					SrcEpoch: 0, SrcVersion: "2.4.48", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "gettext-common-devel", Epoch: 0, Version: "0.19.8.1", Release: "17.el8", Arch: "noarch",
					SrcName: "gettext", SrcEpoch: 0, SrcVersion: "0.19.8.1", SrcRelease: "17.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "coreutils-single", Epoch: 0, Version: "8.30", Release: "7.el8_2.1", Arch: "x86_64",
					SrcName: "coreutils", SrcEpoch: 0, SrcVersion: "8.30", SrcRelease: "7.el8_2.1", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "automake", Epoch: 0, Version: "1.16.1", Release: "6.el8", Arch: "noarch",
					SrcName: "automake", SrcEpoch: 0, SrcVersion: "1.16.1", SrcRelease: "6.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and GFDL and Public Domain and MIT"},
				},
				{
					Name: "libblkid", Epoch: 0, Version: "2.32.1", Release: "22.el8", Arch: "x86_64",
					SrcName: "util-linux", SrcEpoch: 0, SrcVersion: "2.32.1", SrcRelease: "22.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "gcc", Epoch: 0, Version: "8.3.1", Release: "5.el8.0.2", Arch: "x86_64", SrcName: "gcc",
					SrcEpoch: 0, SrcVersion: "8.3.1", SrcRelease: "5.el8.0.2", Modularitylabel: "",
					Licenses: []string{"GPLv3+ and GPLv3+ with exceptions and GPLv2+ with exceptions and LGPLv2+ and BSD"},
				},
				{
					Name: "libcap-ng", Epoch: 0, Version: "0.7.9", Release: "5.el8", Arch: "x86_64",
					SrcName: "libcap-ng", SrcEpoch: 0, SrcVersion: "0.7.9", SrcRelease: "5.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "gettext-devel", Epoch: 0, Version: "0.19.8.1", Release: "17.el8", Arch: "x86_64",
					SrcName: "gettext", SrcEpoch: 0, SrcVersion: "0.19.8.1", SrcRelease: "17.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+ and GPLv3+"},
				},
				{
					Name: "libffi", Epoch: 0, Version: "3.1", Release: "21.el8", Arch: "x86_64", SrcName: "libffi",
					SrcEpoch: 0, SrcVersion: "3.1", SrcRelease: "21.el8", Modularitylabel: "", Licenses: []string{"MIT"},
				},
				{
					Name: "make", Epoch: 1, Version: "4.2.1", Release: "10.el8", Arch: "x86_64", SrcName: "make",
					SrcEpoch: 1, SrcVersion: "4.2.1", SrcRelease: "10.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "libzstd", Epoch: 0, Version: "1.4.2", Release: "2.el8", Arch: "x86_64", SrcName: "zstd",
					SrcEpoch: 0, SrcVersion: "1.4.2", SrcRelease: "2.el8", Modularitylabel: "",
					Licenses: []string{"BSD and GPLv2"},
				},
				{
					Name: "npm", Epoch: 1, Version: "6.14.4", Release: "1.10.21.0.3.module_el8.2.0+391+8da3adc6",
					Arch: "x86_64", SrcName: "nodejs", SrcEpoch: 1, SrcVersion: "10.21.0",
					SrcRelease:      "3.module_el8.2.0+391+8da3adc6",
					Modularitylabel: "nodejs:10:8020020200707141642:6a468ee4",
					Licenses:        []string{"MIT and ASL 2.0 and ISC and BSD"},
				},
				{
					Name: "lz4-libs", Epoch: 0, Version: "1.8.1.2", Release: "4.el8", Arch: "x86_64", SrcName: "lz4",
					SrcEpoch: 0, SrcVersion: "1.8.1.2", SrcRelease: "4.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and BSD"},
				},
				{
					Name: "libtool-ltdl", Epoch: 0, Version: "2.4.6", Release: "25.el8", Arch: "x86_64",
					SrcName: "libtool", SrcEpoch: 0, SrcVersion: "2.4.6", SrcRelease: "25.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "libgcrypt", Epoch: 0, Version: "1.8.3", Release: "4.el8", Arch: "x86_64",
					SrcName: "libgcrypt", SrcEpoch: 0, SrcVersion: "1.8.3", SrcRelease: "4.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "libipt", Epoch: 0, Version: "1.6.1", Release: "8.el8", Arch: "x86_64", SrcName: "libipt",
					SrcEpoch: 0, SrcVersion: "1.6.1", SrcRelease: "8.el8", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "cracklib", Epoch: 0, Version: "2.9.6", Release: "15.el8", Arch: "x86_64",
					SrcName: "cracklib", SrcEpoch: 0, SrcVersion: "2.9.6", SrcRelease: "15.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "gc", Epoch: 0, Version: "7.6.4", Release: "3.el8", Arch: "x86_64", SrcName: "gc",
					SrcEpoch: 0, SrcVersion: "7.6.4", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "libidn2", Epoch: 0, Version: "2.2.0", Release: "1.el8", Arch: "x86_64", SrcName: "libidn2",
					SrcEpoch: 0, SrcVersion: "2.2.0", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"(GPLv2+ or LGPLv3+) and GPLv3+"},
				},
				{
					Name: "gdb-headless", Epoch: 0, Version: "8.2", Release: "12.el8", Arch: "x86_64", SrcName: "gdb",
					SrcEpoch: 0, SrcVersion: "8.2", SrcRelease: "12.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+ and GPLv3+ with exceptions and GPLv2+ and GPLv2+ with exceptions and GPL+ and LGPLv2+ and LGPLv3+ and BSD and Public Domain and GFDL"},
				},
				{
					Name: "file-libs", Epoch: 0, Version: "5.33", Release: "13.el8", Arch: "x86_64", SrcName: "file",
					SrcEpoch: 0, SrcVersion: "5.33", SrcRelease: "13.el8", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "epel-release", Epoch: 0, Version: "8", Release: "8.el8", Arch: "noarch",
					SrcName: "epel-release", SrcEpoch: 0, SrcVersion: "8", SrcRelease: "8.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2"},
				},
				{
					Name: "keyutils-libs", Epoch: 0, Version: "1.5.10", Release: "6.el8", Arch: "x86_64",
					SrcName: "keyutils", SrcEpoch: 0, SrcVersion: "1.5.10", SrcRelease: "6.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and LGPLv2+"},
				},
				{
					Name: "platform-python-pip", Epoch: 0, Version: "9.0.3", Release: "18.el8", Arch: "noarch",
					SrcName: "python-pip", SrcEpoch: 0, SrcVersion: "9.0.3", SrcRelease: "18.el8", Modularitylabel: "",
					Licenses: []string{"MIT and Python and ASL 2.0 and BSD and ISC and LGPLv2 and MPLv2.0 and (ASL 2.0 or BSD)"},
				},
				{
					Name: "p11-kit-trust", Epoch: 0, Version: "0.23.14", Release: "5.el8_0", Arch: "x86_64",
					SrcName: "p11-kit", SrcEpoch: 0, SrcVersion: "0.23.14", SrcRelease: "5.el8_0", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "python36", Epoch: 0, Version: "3.6.8", Release: "2.module_el8.3.0+562+e162826a",
					Arch: "x86_64", SrcName: "python36", SrcEpoch: 0, SrcVersion: "3.6.8",
					SrcRelease:      "2.module_el8.3.0+562+e162826a",
					Modularitylabel: "python36:3.6:8030020201104034153:24f1489c", Licenses: []string{"Python"},
				},
				{
					Name: "pcre", Epoch: 0, Version: "8.42", Release: "4.el8", Arch: "x86_64", SrcName: "pcre",
					SrcEpoch: 0, SrcVersion: "8.42", SrcRelease: "4.el8", Modularitylabel: "", Licenses: []string{"BSD"},
				},
				{
					Name: "python2-setuptools-wheel", Epoch: 0, Version: "39.0.1",
					Release: "12.module_el8.3.0+478+7570e00c", Arch: "noarch", SrcName: "python2-setuptools",
					SrcEpoch: 0, SrcVersion: "39.0.1", SrcRelease: "12.module_el8.3.0+478+7570e00c",
					Modularitylabel: "python27:2.7:8030020200831201838:851f4228", Licenses: []string{"MIT"},
				},
				{
					Name: "systemd-libs", Epoch: 0, Version: "239", Release: "31.el8_2.2", Arch: "x86_64",
					SrcName: "systemd", SrcEpoch: 0, SrcVersion: "239", SrcRelease: "31.el8_2.2", Modularitylabel: "",
					Licenses: []string{"LGPLv2+ and MIT"},
				},
				{
					Name: "python2-libs", Epoch: 0, Version: "2.7.17", Release: "2.module_el8.3.0+478+7570e00c",
					Arch: "x86_64", SrcName: "python2", SrcEpoch: 0, SrcVersion: "2.7.17",
					SrcRelease:      "2.module_el8.3.0+478+7570e00c",
					Modularitylabel: "python27:2.7:8030020200831201838:851f4228", Licenses: []string{"Python"},
				},
				{
					Name: "dbus-tools", Epoch: 1, Version: "1.12.8", Release: "10.el8_2", Arch: "x86_64",
					SrcName: "dbus", SrcEpoch: 1, SrcVersion: "1.12.8", SrcRelease: "10.el8_2", Modularitylabel: "",
					Licenses: []string{"(GPLv2+ or AFL) and GPLv2+"},
				},
				{
					Name: "python2-setuptools", Epoch: 0, Version: "39.0.1", Release: "12.module_el8.3.0+478+7570e00c",
					Arch: "noarch", SrcName: "python2-setuptools", SrcEpoch: 0, SrcVersion: "39.0.1",
					SrcRelease:      "12.module_el8.3.0+478+7570e00c",
					Modularitylabel: "python27:2.7:8030020200831201838:851f4228", Licenses: []string{"MIT"},
				},
				{
					Name: "libusbx", Epoch: 0, Version: "1.0.22", Release: "1.el8", Arch: "x86_64", SrcName: "libusbx",
					SrcEpoch: 0, SrcVersion: "1.0.22", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "gpg-pubkey", Epoch: 0, Version: "ce977fe0", Release: "5db1f171", Arch: "None", SrcName: "",
					SrcEpoch: 0, SrcVersion: "", SrcRelease: "", Modularitylabel: "", Licenses: []string{"pubkey"},
				},
				{
					Name: "rpm-libs", Epoch: 0, Version: "4.14.3", Release: "4.el8", Arch: "x86_64", SrcName: "rpm",
					SrcEpoch: 0, SrcVersion: "4.14.3", SrcRelease: "4.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and LGPLv2+ with exceptions"},
				},
				{
					Name: "squashfs-tools", Epoch: 0, Version: "4.3", Release: "19.el8", Arch: "x86_64",
					SrcName: "squashfs-tools", SrcEpoch: 0, SrcVersion: "4.3", SrcRelease: "19.el8",
					Modularitylabel: "", Licenses: []string{"GPLv2+"},
				},
				{
					Name: "rpm-build-libs", Epoch: 0, Version: "4.14.3", Release: "4.el8", Arch: "x86_64",
					SrcName: "rpm", SrcEpoch: 0, SrcVersion: "4.14.3", SrcRelease: "4.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and LGPLv2+ with exceptions"},
				},
				{
					Name: "libsemanage", Epoch: 0, Version: "2.9", Release: "2.el8", Arch: "x86_64",
					SrcName: "libsemanage", SrcEpoch: 0, SrcVersion: "2.9", SrcRelease: "2.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "libutempter", Epoch: 0, Version: "1.1.6", Release: "14.el8", Arch: "x86_64",
					SrcName: "libutempter", SrcEpoch: 0, SrcVersion: "1.1.6", SrcRelease: "14.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "acl", Epoch: 0, Version: "2.2.53", Release: "1.el8", Arch: "x86_64", SrcName: "acl",
					SrcEpoch: 0, SrcVersion: "2.2.53", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+"},
				},
				{
					Name: "nettle", Epoch: 0, Version: "3.4.1", Release: "1.el8", Arch: "x86_64", SrcName: "nettle",
					SrcEpoch: 0, SrcVersion: "3.4.1", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv3+ or GPLv2+"},
				},
				{
					Name: "libcomps", Epoch: 0, Version: "0.1.11", Release: "4.el8", Arch: "x86_64",
					SrcName: "libcomps", SrcEpoch: 0, SrcVersion: "0.1.11", SrcRelease: "4.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+"},
				},
				{
					Name: "findutils", Epoch: 1, Version: "4.6.0", Release: "20.el8", Arch: "x86_64",
					SrcName: "findutils", SrcEpoch: 1, SrcVersion: "4.6.0", SrcRelease: "20.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "cpio", Epoch: 0, Version: "2.12", Release: "8.el8", Arch: "x86_64", SrcName: "cpio",
					SrcEpoch: 0, SrcVersion: "2.12", SrcRelease: "8.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "ipcalc", Epoch: 0, Version: "0.2.4", Release: "4.el8", Arch: "x86_64", SrcName: "ipcalc",
					SrcEpoch: 0, SrcVersion: "0.2.4", SrcRelease: "4.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+"},
				},
				{
					Name: "libnghttp2", Epoch: 0, Version: "1.33.0", Release: "3.el8_2.1", Arch: "x86_64",
					SrcName: "nghttp2", SrcEpoch: 0, SrcVersion: "1.33.0", SrcRelease: "3.el8_2.1", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "iptables-libs", Epoch: 0, Version: "1.8.4", Release: "10.el8_2.1", Arch: "x86_64",
					SrcName: "iptables", SrcEpoch: 0, SrcVersion: "1.8.4", SrcRelease: "10.el8_2.1",
					Modularitylabel: "", Licenses: []string{"GPLv2 and Artistic 2.0 and ISC"},
				},
				{
					Name: "libsigsegv", Epoch: 0, Version: "2.11", Release: "5.el8", Arch: "x86_64",
					SrcName: "libsigsegv", SrcEpoch: 0, SrcVersion: "2.11", SrcRelease: "5.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+"},
				},
				{
					Name: "libverto", Epoch: 0, Version: "0.3.0", Release: "5.el8", Arch: "x86_64", SrcName: "libverto",
					SrcEpoch: 0, SrcVersion: "0.3.0", SrcRelease: "5.el8", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "libtirpc", Epoch: 0, Version: "1.1.4", Release: "4.el8", Arch: "x86_64", SrcName: "libtirpc",
					SrcEpoch: 0, SrcVersion: "1.1.4", SrcRelease: "4.el8", Modularitylabel: "",
					Licenses: []string{"SISSL and BSD"},
				},
				{
					Name: "openssl-libs", Epoch: 1, Version: "1.1.1c", Release: "15.el8", Arch: "x86_64",
					SrcName: "openssl", SrcEpoch: 1, SrcVersion: "1.1.1c", SrcRelease: "15.el8", Modularitylabel: "",
					Licenses: []string{"OpenSSL"},
				},
				{
					Name: "python3-libs", Epoch: 0, Version: "3.6.8", Release: "23.el8", Arch: "x86_64",
					SrcName: "python3", SrcEpoch: 0, SrcVersion: "3.6.8", SrcRelease: "23.el8", Modularitylabel: "",
					Licenses: []string{"Python"},
				},
				{
					Name: "libpwquality", Epoch: 0, Version: "1.4.0", Release: "9.el8", Arch: "x86_64",
					SrcName: "libpwquality", SrcEpoch: 0, SrcVersion: "1.4.0", SrcRelease: "9.el8", Modularitylabel: "",
					Licenses: []string{"BSD or GPLv2+"},
				},
				{
					Name: "util-linux", Epoch: 0, Version: "2.32.1", Release: "22.el8", Arch: "x86_64",
					SrcName: "util-linux", SrcEpoch: 0, SrcVersion: "2.32.1", SrcRelease: "22.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2 and GPLv2+ and LGPLv2+ and BSD with advertising and Public Domain"},
				},
				{
					Name: "glib2", Epoch: 0, Version: "2.56.4", Release: "8.el8", Arch: "x86_64", SrcName: "glib2",
					SrcEpoch: 0, SrcVersion: "2.56.4", SrcRelease: "8.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "iproute", Epoch: 0, Version: "5.3.0", Release: "1.el8", Arch: "x86_64", SrcName: "iproute",
					SrcEpoch: 0, SrcVersion: "5.3.0", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and Public Domain"},
				},
				{
					Name: "kmod", Epoch: 0, Version: "25", Release: "16.el8", Arch: "x86_64", SrcName: "kmod",
					SrcEpoch: 0, SrcVersion: "25", SrcRelease: "16.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+"},
				},
				{
					Name: "curl", Epoch: 0, Version: "7.61.1", Release: "12.el8", Arch: "x86_64", SrcName: "curl",
					SrcEpoch: 0, SrcVersion: "7.61.1", SrcRelease: "12.el8", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "openldap", Epoch: 0, Version: "2.4.46", Release: "11.el8_1", Arch: "x86_64",
					SrcName: "openldap", SrcEpoch: 0, SrcVersion: "2.4.46", SrcRelease: "11.el8_1", Modularitylabel: "",
					Licenses: []string{"OpenLDAP"},
				},
				{
					Name: "python3-libcomps", Epoch: 0, Version: "0.1.11", Release: "4.el8", Arch: "x86_64",
					SrcName: "libcomps", SrcEpoch: 0, SrcVersion: "0.1.11", SrcRelease: "4.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+"},
				},
				{
					Name: "libarchive", Epoch: 0, Version: "3.3.2", Release: "8.el8_1", Arch: "x86_64",
					SrcName: "libarchive", SrcEpoch: 0, SrcVersion: "3.3.2", SrcRelease: "8.el8_1", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "libyaml", Epoch: 0, Version: "0.1.7", Release: "5.el8", Arch: "x86_64", SrcName: "libyaml",
					SrcEpoch: 0, SrcVersion: "0.1.7", SrcRelease: "5.el8", Modularitylabel: "",
					Licenses: []string{"MIT"},
				},
				{
					Name: "npth", Epoch: 0, Version: "1.5", Release: "4.el8", Arch: "x86_64", SrcName: "npth",
					SrcEpoch: 0, SrcVersion: "1.5", SrcRelease: "4.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "gpgme", Epoch: 0, Version: "1.10.0", Release: "6.el8.0.1", Arch: "x86_64", SrcName: "gpgme",
					SrcEpoch: 0, SrcVersion: "1.10.0", SrcRelease: "6.el8.0.1", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "libdnf", Epoch: 0, Version: "0.39.1", Release: "6.el8_2", Arch: "x86_64", SrcName: "libdnf",
					SrcEpoch: 0, SrcVersion: "0.39.1", SrcRelease: "6.el8_2", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "python3-hawkey", Epoch: 0, Version: "0.39.1", Release: "6.el8_2", Arch: "x86_64",
					SrcName: "libdnf", SrcEpoch: 0, SrcVersion: "0.39.1", SrcRelease: "6.el8_2", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "libreport-filesystem", Epoch: 0, Version: "2.9.5", Release: "10.el8", Arch: "x86_64",
					SrcName: "libreport", SrcEpoch: 0, SrcVersion: "2.9.5", SrcRelease: "10.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+"},
				},
				{
					Name: "dhcp-common", Epoch: 12, Version: "4.3.6", Release: "40.el8", Arch: "noarch",
					SrcName: "dhcp", SrcEpoch: 12, SrcVersion: "4.3.6", SrcRelease: "40.el8", Modularitylabel: "",
					Licenses: []string{"ISC"},
				},
				{
					Name: "dbus-daemon", Epoch: 1, Version: "1.12.8", Release: "10.el8_2", Arch: "x86_64",
					SrcName: "dbus", SrcEpoch: 1, SrcVersion: "1.12.8", SrcRelease: "10.el8_2", Modularitylabel: "",
					Licenses: []string{"(GPLv2+ or AFL) and GPLv2+"},
				},
				{
					Name: "device-mapper-libs", Epoch: 8, Version: "1.02.169", Release: "3.el8", Arch: "x86_64",
					SrcName: "lvm2", SrcEpoch: 8, SrcVersion: "2.03.08", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2"},
				},
				{
					Name: "elfutils-default-yama-scope", Epoch: 0, Version: "0.178", Release: "7.el8", Arch: "noarch",
					SrcName: "elfutils", SrcEpoch: 0, SrcVersion: "0.178", SrcRelease: "7.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+ or LGPLv3+"},
				},
				{
					Name: "systemd-pam", Epoch: 0, Version: "239", Release: "31.el8_2.2", Arch: "x86_64",
					SrcName: "systemd", SrcEpoch: 0, SrcVersion: "239", SrcRelease: "31.el8_2.2", Modularitylabel: "",
					Licenses: []string{"LGPLv2+ and MIT and GPLv2+"},
				},
				{
					Name: "dbus", Epoch: 1, Version: "1.12.8", Release: "10.el8_2", Arch: "x86_64", SrcName: "dbus",
					SrcEpoch: 1, SrcVersion: "1.12.8", SrcRelease: "10.el8_2", Modularitylabel: "",
					Licenses: []string{"(GPLv2+ or AFL) and GPLv2+"},
				},
				{
					Name: "dhcp-client", Epoch: 12, Version: "4.3.6", Release: "40.el8", Arch: "x86_64",
					SrcName: "dhcp", SrcEpoch: 12, SrcVersion: "4.3.6", SrcRelease: "40.el8", Modularitylabel: "",
					Licenses: []string{"ISC"},
				},
				{
					Name: "libkcapi-hmaccalc", Epoch: 0, Version: "1.1.1", Release: "16_1.el8", Arch: "x86_64",
					SrcName: "libkcapi", SrcEpoch: 0, SrcVersion: "1.1.1", SrcRelease: "16_1.el8", Modularitylabel: "",
					Licenses: []string{"BSD or GPLv2"},
				},
				{
					Name: "dracut", Epoch: 0, Version: "049", Release: "70.git20200228.el8", Arch: "x86_64",
					SrcName: "dracut", SrcEpoch: 0, SrcVersion: "049", SrcRelease: "70.git20200228.el8",
					Modularitylabel: "", Licenses: []string{"GPLv2+ and LGPLv2+"},
				},
				{
					Name: "dracut-squash", Epoch: 0, Version: "049", Release: "70.git20200228.el8", Arch: "x86_64",
					SrcName: "dracut", SrcEpoch: 0, SrcVersion: "049", SrcRelease: "70.git20200228.el8",
					Modularitylabel: "", Licenses: []string{"GPLv2+ and LGPLv2+"},
				},
				{
					Name: "dnf", Epoch: 0, Version: "4.2.17", Release: "7.el8_2", Arch: "noarch", SrcName: "dnf",
					SrcEpoch: 0, SrcVersion: "4.2.17", SrcRelease: "7.el8_2", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and GPLv2 and GPL"},
				},
				{
					Name: "kexec-tools", Epoch: 0, Version: "2.0.20", Release: "14.el8", Arch: "x86_64",
					SrcName: "kexec-tools", SrcEpoch: 0, SrcVersion: "2.0.20", SrcRelease: "14.el8",
					Modularitylabel: "", Licenses: []string{"GPLv2"},
				},
				{
					Name: "tar", Epoch: 2, Version: "1.30", Release: "4.el8", Arch: "x86_64", SrcName: "tar",
					SrcEpoch: 2, SrcVersion: "1.30", SrcRelease: "4.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+"},
				},
				{
					Name: "hostname", Epoch: 0, Version: "3.20", Release: "6.el8", Arch: "x86_64", SrcName: "hostname",
					SrcEpoch: 0, SrcVersion: "3.20", SrcRelease: "6.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+"},
				},
				{
					Name: "langpacks-en", Epoch: 0, Version: "1.0", Release: "12.el8", Arch: "noarch",
					SrcName: "langpacks", SrcEpoch: 0, SrcVersion: "1.0", SrcRelease: "12.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+"},
				},
				{
					Name: "gpg-pubkey", Epoch: 0, Version: "8483c65d", Release: "5ccc5b19", Arch: "None", SrcName: "",
					SrcEpoch: 0, SrcVersion: "", SrcRelease: "", Modularitylabel: "", Licenses: []string{"pubkey"},
				},
				{
					Name: "centos-release", Epoch: 0, Version: "8.2", Release: "2.2004.0.2.el8", Arch: "x86_64",
					SrcName: "centos-release", SrcEpoch: 0, SrcVersion: "8.2", SrcRelease: "2.2004.0.2.el8",
					Modularitylabel: "", Licenses: []string{"GPLv2"},
				},
				{
					Name: "zlib", Epoch: 0, Version: "1.2.11", Release: "16.el8_2", Arch: "x86_64", SrcName: "zlib",
					SrcEpoch: 0, SrcVersion: "1.2.11", SrcRelease: "16.el8_2", Modularitylabel: "",
					Licenses: []string{"zlib and Boost"},
				},
				{
					Name: "librepo", Epoch: 0, Version: "1.11.0", Release: "3.el8_2", Arch: "x86_64",
					SrcName: "librepo", SrcEpoch: 0, SrcVersion: "1.11.0", SrcRelease: "3.el8_2", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "bind-export-libs", Epoch: 32, Version: "9.11.13", Release: "6.el8_2.1", Arch: "x86_64",
					SrcName: "bind", SrcEpoch: 32, SrcVersion: "9.11.13", SrcRelease: "6.el8_2.1", Modularitylabel: "",
					Licenses: []string{"MPLv2.0"},
				},
				{
					Name: "perl-libs", Epoch: 4, Version: "5.26.3", Release: "416.el8", Arch: "x86_64", SrcName: "perl",
					SrcEpoch: 4, SrcVersion: "5.26.3", SrcRelease: "416.el8", Modularitylabel: "",
					Licenses: []string{"(GPL+ or Artistic) and HSRL and MIT and UCD"},
				},
				{
					Name: "perl-Scalar-List-Utils", Epoch: 3, Version: "1.49", Release: "2.el8", Arch: "x86_64",
					SrcName: "perl-Scalar-List-Utils", SrcEpoch: 3, SrcVersion: "1.49", SrcRelease: "2.el8",
					Modularitylabel: "", Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "nspr", Epoch: 0, Version: "4.25.0", Release: "2.el8_2", Arch: "x86_64", SrcName: "nspr",
					SrcEpoch: 0, SrcVersion: "4.25.0", SrcRelease: "2.el8_2", Modularitylabel: "",
					Licenses: []string{"MPLv2.0"},
				},
				{
					Name: "nss-softokn-freebl", Epoch: 0, Version: "3.53.1", Release: "11.el8_2", Arch: "x86_64",
					SrcName: "nss", SrcEpoch: 0, SrcVersion: "3.53.1", SrcRelease: "11.el8_2", Modularitylabel: "",
					Licenses: []string{"MPLv2.0"},
				},
				{
					Name: "nss", Epoch: 0, Version: "3.53.1", Release: "11.el8_2", Arch: "x86_64", SrcName: "nss",
					SrcEpoch: 0, SrcVersion: "3.53.1", SrcRelease: "11.el8_2", Modularitylabel: "",
					Licenses: []string{"MPLv2.0"},
				},
				{
					Name: "perl-Text-ParseWords", Epoch: 0, Version: "3.30", Release: "395.el8", Arch: "noarch",
					SrcName: "perl-Text-ParseWords", SrcEpoch: 0, SrcVersion: "3.30", SrcRelease: "395.el8",
					Modularitylabel: "", Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "perl-Term-ANSIColor", Epoch: 0, Version: "4.06", Release: "396.el8", Arch: "noarch",
					SrcName: "perl-Term-ANSIColor", SrcEpoch: 0, SrcVersion: "4.06", SrcRelease: "396.el8",
					Modularitylabel: "", Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "perl-Errno", Epoch: 0, Version: "1.28", Release: "416.el8", Arch: "x86_64", SrcName: "perl",
					SrcEpoch: 0, SrcVersion: "5.26.3", SrcRelease: "416.el8", Modularitylabel: "",
					Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "perl-Text-Tabs+Wrap", Epoch: 0, Version: "2013.0523", Release: "395.el8", Arch: "noarch",
					SrcName: "perl-Text-Tabs+Wrap", SrcEpoch: 0, SrcVersion: "2013.0523", SrcRelease: "395.el8",
					Modularitylabel: "", Licenses: []string{"TTWL"},
				},
				{
					Name: "perl-File-Path", Epoch: 0, Version: "2.15", Release: "2.el8", Arch: "noarch",
					SrcName: "perl-File-Path", SrcEpoch: 0, SrcVersion: "2.15", SrcRelease: "2.el8",
					Modularitylabel: "", Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "perl-PathTools", Epoch: 0, Version: "3.74", Release: "1.el8", Arch: "x86_64",
					SrcName: "perl-PathTools", SrcEpoch: 0, SrcVersion: "3.74", SrcRelease: "1.el8",
					Modularitylabel: "", Licenses: []string{"(GPL+ or Artistic) and BSD"},
				},
				{
					Name: "perl-threads", Epoch: 1, Version: "2.21", Release: "2.el8", Arch: "x86_64",
					SrcName: "perl-threads", SrcEpoch: 1, SrcVersion: "2.21", SrcRelease: "2.el8", Modularitylabel: "",
					Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "perl-interpreter", Epoch: 4, Version: "5.26.3", Release: "416.el8", Arch: "x86_64",
					SrcName: "perl", SrcEpoch: 4, SrcVersion: "5.26.3", SrcRelease: "416.el8", Modularitylabel: "",
					Licenses: []string{"(GPL+ or Artistic) and (GPLv2+ or Artistic) and BSD and Public Domain and UCD"},
				},
				{
					Name: "perl-IO-Socket-IP", Epoch: 0, Version: "0.39", Release: "5.el8", Arch: "noarch",
					SrcName: "perl-IO-Socket-IP", SrcEpoch: 0, SrcVersion: "0.39", SrcRelease: "5.el8",
					Modularitylabel: "", Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "perl-File-Temp", Epoch: 0, Version: "0.230.600", Release: "1.el8", Arch: "noarch",
					SrcName: "perl-File-Temp", SrcEpoch: 0, SrcVersion: "0.230.600", SrcRelease: "1.el8",
					Modularitylabel: "", Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "perl-Digest-MD5", Epoch: 0, Version: "2.55", Release: "396.el8", Arch: "x86_64",
					SrcName: "perl-Digest-MD5", SrcEpoch: 0, SrcVersion: "2.55", SrcRelease: "396.el8",
					Modularitylabel: "", Licenses: []string{"(GPL+ or Artistic) and BSD"},
				},
				{
					Name: "perl-Error", Epoch: 1, Version: "0.17025", Release: "2.el8", Arch: "noarch",
					SrcName: "perl-Error", SrcEpoch: 1, SrcVersion: "0.17025", SrcRelease: "2.el8", Modularitylabel: "",
					Licenses: []string{"(GPL+ or Artistic) and MIT"},
				},
				{
					Name: "perl-Data-Dumper", Epoch: 0, Version: "2.167", Release: "399.el8", Arch: "x86_64",
					SrcName: "perl-Data-Dumper", SrcEpoch: 0, SrcVersion: "2.167", SrcRelease: "399.el8",
					Modularitylabel: "", Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "perl-Storable", Epoch: 1, Version: "3.11", Release: "3.el8", Arch: "x86_64",
					SrcName: "perl-Storable", SrcEpoch: 1, SrcVersion: "3.11", SrcRelease: "3.el8", Modularitylabel: "",
					Licenses: []string{"GPL+ or Artistic"},
				},
				{
					Name: "fipscheck-lib", Epoch: 0, Version: "1.5.0", Release: "4.el8", Arch: "x86_64",
					SrcName: "fipscheck", SrcEpoch: 0, SrcVersion: "1.5.0", SrcRelease: "4.el8", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "openssh", Epoch: 0, Version: "8.0p1", Release: "4.el8_1", Arch: "x86_64", SrcName: "openssh",
					SrcEpoch: 0, SrcVersion: "8.0p1", SrcRelease: "4.el8_1", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "publicsuffix-list-dafsa", Epoch: 0, Version: "20180723", Release: "1.el8", Arch: "noarch",
					SrcName: "publicsuffix-list", SrcEpoch: 0, SrcVersion: "20180723", SrcRelease: "1.el8",
					Modularitylabel: "", Licenses: []string{"MPLv2.0"},
				},
				{
					Name: "pkgconf-m4", Epoch: 0, Version: "1.4.2", Release: "1.el8", Arch: "noarch",
					SrcName: "pkgconf", SrcEpoch: 0, SrcVersion: "1.4.2", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+ with exceptions"},
				},
				{
					Name: "pcre2-utf16", Epoch: 0, Version: "10.32", Release: "1.el8", Arch: "x86_64", SrcName: "pcre2",
					SrcEpoch: 0, SrcVersion: "10.32", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "ncurses", Epoch: 0, Version: "6.1", Release: "7.20180224.el8", Arch: "x86_64",
					SrcName: "ncurses", SrcEpoch: 0, SrcVersion: "6.1", SrcRelease: "7.20180224.el8",
					Modularitylabel: "", Licenses: []string{"MIT"},
				},
				{
					Name: "libsecret", Epoch: 0, Version: "0.18.6", Release: "1.el8", Arch: "x86_64",
					SrcName: "libsecret", SrcEpoch: 0, SrcVersion: "0.18.6", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "pkgconf", Epoch: 0, Version: "1.4.2", Release: "1.el8", Arch: "x86_64", SrcName: "pkgconf",
					SrcEpoch: 0, SrcVersion: "1.4.2", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"ISC"},
				},
				{
					Name: "nspr-devel", Epoch: 0, Version: "4.25.0", Release: "2.el8_2", Arch: "x86_64",
					SrcName: "nspr", SrcEpoch: 0, SrcVersion: "4.25.0", SrcRelease: "2.el8_2", Modularitylabel: "",
					Licenses: []string{"MPLv2.0"},
				},
				{
					Name: "nss-softokn-devel", Epoch: 0, Version: "3.53.1", Release: "11.el8_2", Arch: "x86_64",
					SrcName: "nss", SrcEpoch: 0, SrcVersion: "3.53.1", SrcRelease: "11.el8_2", Modularitylabel: "",
					Licenses: []string{"MPLv2.0"},
				},
				{
					Name: "libsepol-devel", Epoch: 0, Version: "2.9", Release: "1.el8", Arch: "x86_64",
					SrcName: "libsepol", SrcEpoch: 0, SrcVersion: "2.9", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"LGPLv2+"},
				},
				{
					Name: "pcre2-devel", Epoch: 0, Version: "10.32", Release: "1.el8", Arch: "x86_64", SrcName: "pcre2",
					SrcEpoch: 0, SrcVersion: "10.32", SrcRelease: "1.el8", Modularitylabel: "",
					Licenses: []string{"BSD"},
				},
				{
					Name: "zlib-devel", Epoch: 0, Version: "1.2.11", Release: "16.el8_2", Arch: "x86_64",
					SrcName: "zlib", SrcEpoch: 0, SrcVersion: "1.2.11", SrcRelease: "16.el8_2", Modularitylabel: "",
					Licenses: []string{"zlib and Boost"},
				},
				{
					Name: "libedit", Epoch: 0, Version: "3.1", Release: "23.20170329cvs.el8", Arch: "x86_64",
					SrcName: "libedit", SrcEpoch: 0, SrcVersion: "3.1", SrcRelease: "23.20170329cvs.el8",
					Modularitylabel: "", Licenses: []string{"BSD"},
				},
				{
					Name: "git-core", Epoch: 0, Version: "2.18.4", Release: "2.el8_2", Arch: "x86_64", SrcName: "git",
					SrcEpoch: 0, SrcVersion: "2.18.4", SrcRelease: "2.el8_2", Modularitylabel: "",
					Licenses: []string{"GPLv2"},
				},
				{
					Name: "keyutils-libs-devel", Epoch: 0, Version: "1.5.10", Release: "6.el8", Arch: "x86_64",
					SrcName: "keyutils", SrcEpoch: 0, SrcVersion: "1.5.10", SrcRelease: "6.el8", Modularitylabel: "",
					Licenses: []string{"GPLv2+ and LGPLv2+"},
				},
				{
					Name: "groff-base", Epoch: 0, Version: "1.22.3", Release: "18.el8", Arch: "x86_64",
					SrcName: "groff", SrcEpoch: 0, SrcVersion: "1.22.3", SrcRelease: "18.el8", Modularitylabel: "",
					Licenses: []string{"GPLv3+ and GFDL and BSD and MIT"},
				},
				{
					Name: "perl-Pod-Simple", Epoch: 1, Version: "3.35", Release: "395.el8", Arch: "noarch",
					SrcName: "perl-Pod-Simple", SrcEpoch: 1, SrcVersion: "3.35", SrcRelease: "395.el8",
					Modularitylabel: "", Licenses: []string{"GPL+ or Artistic"},
				},
			},
		},
	}
	a := rpmPkgAnalyzer{}
	for testname, tc := range tests {
		t.Run(testname, func(t *testing.T) {
			f, err := os.Open(tc.path)
			require.NoError(t, err)
			defer f.Close()

			pkgs, _, err := a.parsePkgInfo(f)
			require.NoError(t, err)

			sort.Slice(tc.pkgs, func(i, j int) bool {
				return tc.pkgs[i].Name < tc.pkgs[j].Name
			})
			sort.Slice(pkgs, func(i, j int) bool {
				return pkgs[i].Name < pkgs[j].Name
			})

			assert.Equal(t, tc.pkgs, pkgs)
		})
	}
}

func Test_splitFileName(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantName string
		wantVer  string
		wantRel  string
		wantErr  bool
	}{
		{
			name:     "valid name",
			filename: "glibc-2.17-307.el7.1.src.rpm",
			wantName: "glibc",
			wantVer:  "2.17",
			wantRel:  "307.el7.1",
			wantErr:  false,
		},
		{
			name:     "invalid name",
			filename: "elasticsearch-5.6.16-1-src.rpm",
			wantName: "",
			wantVer:  "",
			wantRel:  "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotVer, gotRel, err := splitFileName(tt.filename)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantName, gotName)
			assert.Equal(t, tt.wantVer, gotVer)
			assert.Equal(t, tt.wantRel, gotRel)
		})
	}
}
