# svncheckins
A utility that uses goroutines to query the status of a set of Subversion repositories.
This was written for a project that had approaching 20 separate Maven projects in the
build, and while we were still working on our release/build/CI processes it was often
useful to be able to tell what trunk and branch commits had recently been made across
these projects. Running a simple command script (such as the svncheckins.cmd listed
here) worked ok, but took over a minute to run. The Go version was able to query
information about the trunk, branches and POM for all the projects in about 5 seconds.
