# v0.4.0 (2018-03-03)

* added functionality to list s3 objects
* added functionality to download a s3 object
* added functionality to download s3 objects in batch
* minor performance enhancements

# v0.3.0 (2018-02-12)

* added functionality to list based on filters
* added functionality to start or stop multiple instances
* added ec2 subcommand. Closes #4
* added unit tests for store
* added mocks for testing
* chore: added prerun hooks
* Fixes #5: added functionality to list s3 buckets
* fixed weird whitespace caused by word wrap function
* fixed: silence usage on errors. Closes #6
* fixed verb for bool type
* fixed import errors
* refactored store to be unit testable
* refactored functions to accommodate unit testing
* refactored: directory structure

# v0.2.1 (2018-01-26)

* added word wrap utility function. Closes #2
* added custom spinner and spinner colors
* added tablewriter for rendering lists
* added docstring for godoc generation
* added unit test for utils
* added Makefile
* added travis build
* added build and godoc badges
* refactored: go fmt
* refactored: changed vendor in build script

# v0.2.0 (2018-01-07)

* added upgrade and rds list feature
* added dependencies
* refactored installation script: display download progress
* refactored: use a version constant
* refactored: rename sqlite filename
* refactored: changed project directory structure
* fixed badge link in README

# v0.1.0 (2017-12-31)

* Initial release